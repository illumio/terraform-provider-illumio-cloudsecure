// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"context"
	"fmt"
	"maps"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	policyAPICilium = "cilium"
	policyAPIAwsVpc = "aws_vpc"

	// awsVpcMaxPriority is the highest priority the aws_vpc policy API supports.
	awsVpcMaxPriority = 1000

	tierAdmin    = "Admin"
	tierBaseline = "Baseline"
)

// clusterPolicyAPIAttribute is policy_api on a clusters entry: the Kubernetes network policy API to
// render for the clusters the entry selects.
//
// It defaults to "", meaning unset. The generated converters send an omitted value as "" and read
// it back as "", so without the default, apply would store "" where the plan had null and fail
// with an inconsistent result.
var clusterPolicyAPIAttribute = resource_schema.StringAttribute{
	MarkdownDescription: "The Kubernetes network policy API required from the clusters this entry selects: " +
		"`cilium` (Cilium's `CiliumClusterwideNetworkPolicy` API) or `aws_vpc` (AWS VPC CNI plugin's `ClusterNetworkPolicy` API). " +
		"When set, every selected cluster must provide it. When unset, any cluster can be used. " +
		"Some rule attributes require it; see each attribute.",
	Optional: true,
	Computed: true,
	Default:  stringdefault.StaticString(""),
	Validators: []validator.String{
		// "" is the default, so writing it explicitly is the same as omitting it.
		stringvalidator.OneOf("", policyAPICilium, policyAPIAwsVpc),
	},
}

// k8sPolicyAttributes are the attributes shared by source.k8s and destination.k8s that control
// how that side of a rule is rendered into Kubernetes network policies.
//
// tier and priority default to "Admin" and 0. The generated converters cannot tell an unset
// nested scalar from its zero value, so the default is planned up front, which keeps apply
// consistent with the plan and makes removing the attribute from config a change. (Optional and
// Computed alone would instead keep the prior value, so removing tier would plan no change.)
var k8sPolicyAttributes = map[string]resource_schema.Attribute{
	"tier": resource_schema.StringAttribute{
		MarkdownDescription: "Policy tier for this side of the rule: `Admin` or `Baseline`. Defaults to `Admin`. " +
			"`Baseline` requires `policy_api = \"aws_vpc\"` on every clusters entry of this side, except for a default deny: " +
			"a `Deny` to or from `any` with `port_ranges` exactly TCP 0-0 and UDP 0-0. " +
			"Both APIs support the egress default deny (to destination `any`); the ingress one (from source `any`) " +
			"requires `policy_api = \"cilium\"`, like all ingress from source `any`.",
		Optional: true,
		Computed: true,
		Default:  stringdefault.StaticString(tierAdmin),
		Validators: []validator.String{
			stringvalidator.OneOf(tierAdmin, tierBaseline),
		},
	},
	"priority": resource_schema.Int64Attribute{
		MarkdownDescription: "Priority of this side of the rule within its tier, between 0 and 1000, inclusive. " +
			"A non-zero priority requires `policy_api = \"aws_vpc\"` on every clusters entry of this side.",
		Optional: true,
		Computed: true,
		Default:  int64default.StaticInt64(0),
		Validators: []validator.Int64{
			int64validator.Between(0, awsVpcMaxPriority),
		},
	},
}

// withK8sPolicyAttributes returns a copy of attributes with k8sPolicyAttributes added.
func withK8sPolicyAttributes(attributes map[string]resource_schema.Attribute) map[string]resource_schema.Attribute {
	m := make(map[string]resource_schema.Attribute, len(attributes)+len(k8sPolicyAttributes))
	maps.Copy(m, attributes)
	maps.Copy(m, k8sPolicyAttributes)

	return m
}

// policyVersionRuleValidator checks a policy version rule against the constraints that can only
// be expressed across attributes:
//
//   - source.any and destination.any cannot both be set.
//   - A k8s side that uses a feature only one policy API supports must set that policy_api on
//     every clusters entry.
//
// Values that are unknown at plan time are skipped; the backend validates them.
type policyVersionRuleValidator struct{}

var _ validator.Object = policyVersionRuleValidator{}

func (policyVersionRuleValidator) Description(_ context.Context) string {
	return "rule must not use `any` on both sides, and must use only features supported by the policy_api of every clusters entry"
}

func (v policyVersionRuleValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (policyVersionRuleValidator) ValidateObject(_ context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	rule, ok := knownAttributes(req.ConfigValue)
	if !ok {
		return
	}

	r := ruleCheck{
		name:  ruleName(req.Path),
		path:  req.Path,
		rule:  rule,
		diags: resp,
	}

	r.checkAny()
	r.checkSide("source")
	r.checkSide("destination")
}

// ruleCheck holds one rule while it is being validated.
type ruleCheck struct {
	name  string
	path  path.Path
	rule  map[string]attr.Value
	diags *validator.ObjectResponse
}

func (r ruleCheck) addError(attrPath path.Path, side, reason string) {
	r.diags.Diagnostics.AddAttributeError(attrPath, "Invalid Policy Version Rule",
		fmt.Sprintf("Rule %q, %s: %s", r.name, side, reason))
}

// peerOf returns the side of the rule opposite side.
func peerOf(side string) string {
	if side == "source" {
		return "destination"
	}

	return "source"
}

// checkAny rejects any on both sides.
func (r ruleCheck) checkAny() {
	source, okSource := knownAttributes(r.rule["source"])
	destination, okDestination := knownAttributes(r.rule["destination"])

	if okSource && okDestination && isSet(source["any"]) && isSet(destination["any"]) {
		r.addError(r.path.AtName("destination").AtName("any"), "destination",
			"`any` is not valid when the source is also `any`.")
	}
}

// checkSide checks that every feature one k8s side uses is supported by the policy_api of each of
// its clusters entries.
func (r ruleCheck) checkSide(side string) {
	sideAttributes, ok := knownAttributes(r.rule[side])
	if !ok {
		return
	}

	k8s, ok := knownAttributes(sideAttributes["k8s"])
	if !ok {
		return
	}

	k8sPath := r.path.AtName(side).AtName("k8s")
	sideName := side + ".k8s"
	clusters := k8s["clusters"]

	if priority, ok := knownInt64(k8s["priority"]); ok && priority > 0 {
		r.requireAPI(k8sPath.AtName("priority"), sideName, "a non-zero priority", clusters, policyAPIAwsVpc)
	}

	if tier, ok := knownString(k8s["tier"]); ok && tier == tierBaseline && !r.isDefaultDeny(side) {
		r.requireAPI(k8sPath.AtName("tier"), sideName,
			fmt.Sprintf("tier %q, other than for a default deny (a Deny to or from `any` with port_ranges exactly TCP 0-0 and UDP 0-0),", tierBaseline),
			clusters, policyAPIAwsVpc)
	}

	if hasServiceAccounts(k8s) {
		r.requireAPI(k8sPath.AtName("workload_selector").AtName("service_accounts"), sideName,
			"workload_selector.service_accounts", clusters, policyAPICilium)
	}

	// AWS ClusterNetworkPolicy has no CIDR ingress peer, so ingress from IP lists or from any
	// source needs Cilium.
	if source, ok := knownAttributes(r.rule["source"]); ok && side == "destination" {
		for _, name := range []string{"any", "ip_lists"} {
			if isSet(source[name]) {
				r.requireAPI(r.path.AtName("source").AtName(name), sideName, "ingress from source."+name, clusters, policyAPICilium)
			}
		}
	}
}

// requireAPI reports an error unless every clusters entry sets policy_api to api. Entries whose
// policy_api is unknown at plan time are skipped.
func (r ruleCheck) requireAPI(attrPath path.Path, sideName, feature string, clusters attr.Value, api string) {
	list, ok := clusters.(types.List)
	if !ok || !isSet(clusters) {
		return
	}

	for _, element := range list.Elements() {
		entry, ok := knownAttributes(element)
		if !ok || entry["policy_api"].IsUnknown() {
			continue
		}

		// A null policy_api is unset, so it matches no API.
		if value, _ := knownString(entry["policy_api"]); value != api {
			r.addError(attrPath, sideName, fmt.Sprintf("%s requires policy_api %q on every clusters entry.", feature, api))

			return
		}
	}
}

// isDefaultDeny reports whether the rule denies all TCP and UDP traffic between the k8s side and
// any peer: a Deny to destination any (egress) or from source any (ingress), with port_ranges
// exactly TCP 0-0 and UDP 0-0. Both policy APIs render this shape in the Baseline tier, so it is
// exempt from the Baseline requirement; the ingress form still needs cilium for its source.any
// peer, which checkSide checks separately. Unknown values count as a match, leaving them to the
// backend.
func (r ruleCheck) isDefaultDeny(side string) bool {
	if action, ok := knownString(r.rule["action"]); ok && action != "Deny" {
		return false
	}

	if peer, ok := knownAttributes(r.rule[peerOf(side)]); ok && isNull(peer["any"]) {
		return false
	}

	if allPorts, known := isAllPortsTCPAndUDP(r.rule["port_ranges"]); known && !allPorts {
		return false
	}

	return true
}

// hasServiceAccounts reports whether the k8s side's workload_selector lists service accounts.
func hasServiceAccounts(k8s map[string]attr.Value) bool {
	workloadSelector, ok := knownAttributes(k8s["workload_selector"])
	if !ok {
		return false
	}

	serviceAccounts, ok := workloadSelector["service_accounts"].(types.List)

	return ok && isSet(serviceAccounts) && len(serviceAccounts.Elements()) > 0
}

// isAllPortsTCPAndUDP reports whether the port ranges are exactly "all TCP and UDP ports" in the
// form a Baseline default deny uses: every range is 0-0, which the backend renders as all ports,
// and both TCP and UDP are present. known is false when any value needed to decide is unknown.
func isAllPortsTCPAndUDP(portRanges attr.Value) (allPorts, known bool) {
	if portRanges == nil || portRanges.IsUnknown() {
		return false, false
	}

	list, ok := portRanges.(types.List)
	if !ok || portRanges.IsNull() {
		return false, true
	}

	protocols := map[string]bool{}

	for _, element := range list.Elements() {
		portRange, ok := knownAttributes(element)
		if !ok {
			return false, false
		}

		protocol, okProtocol := knownString(portRange["protocol"])
		from, okFrom := knownInt64(portRange["from_port"])
		to, okTo := knownInt64(portRange["to_port"])

		if !okProtocol || !okFrom || !okTo {
			return false, false
		}

		if from != 0 || to != 0 {
			return false, true
		}

		protocols[strings.ToUpper(protocol)] = true
	}

	return protocols["TCP"] && protocols["UDP"], true
}

// ruleName returns the map key of the rule at p, or p itself if it has none.
func ruleName(p path.Path) string {
	steps := p.Steps()
	if len(steps) > 0 {
		if key, ok := steps[len(steps)-1].(path.PathStepElementKeyString); ok {
			return string(key)
		}
	}

	return p.String()
}

// isNull reports whether v is known to be null.
func isNull(v attr.Value) bool {
	return v != nil && !v.IsUnknown() && v.IsNull()
}

// isSet reports whether v is known and not null.
func isSet(v attr.Value) bool {
	return v != nil && !v.IsUnknown() && !v.IsNull()
}

// knownAttributes returns the attributes of v if it is a known, non-null object.
func knownAttributes(v attr.Value) (map[string]attr.Value, bool) {
	if !isSet(v) {
		return nil, false
	}

	object, ok := v.(types.Object)
	if !ok {
		return nil, false
	}

	return object.Attributes(), true
}

// knownString returns the value of v if it is a known, non-null string.
func knownString(v attr.Value) (string, bool) {
	if !isSet(v) {
		return "", false
	}

	s, ok := v.(types.String)
	if !ok {
		return "", false
	}

	return s.ValueString(), true
}

// knownInt64 returns the value of v if it is a known, non-null integer.
func knownInt64(v attr.Value) (int64, bool) {
	if !isSet(v) {
		return 0, false
	}

	i, ok := v.(types.Int64)
	if !ok {
		return 0, false
	}

	return i.ValueInt64(), true
}
