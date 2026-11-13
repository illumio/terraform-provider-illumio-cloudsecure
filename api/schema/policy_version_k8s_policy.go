// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"context"
	"fmt"
	"maps"
	"slices"

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
	policyAPIAWSVPC = "aws_vpc"

	// awsVPCMaxPriority is the highest priority the aws_vpc policy API supports.
	awsVPCMaxPriority = 1000

	tierAdmin    = "Admin"
	tierBaseline = "Baseline"
)

// clusterPolicyAPIAttribute is policy_api on a clusters entry: it selects clusters by the Kubernetes
// network policy API they provide, together with the entry's other selectors.
//
// TODO: Omitting it currently fails apply with an inconsistent result: the generated converters
// don't make optional scalars in nested objects proto3 optional, so null is sent as "" and read
// back as "". Fix the generators rather than adding a default.
var clusterPolicyAPIAttribute = resource_schema.StringAttribute{
	MarkdownDescription: "Selects only clusters that provide this Kubernetes network policy API: " +
		"`cilium` (Cilium's `CiliumClusterwideNetworkPolicy` API) or `aws_vpc` (AWS VPC CNI plugin's `ClusterNetworkPolicy` API). " +
		"When unset, any cluster is selected. " +
		"Some rule attributes require this to be set to specific values; see the description for each attribute.",
	Optional: true,
	Validators: []validator.String{
		stringvalidator.OneOf(policyAPICilium, policyAPIAWSVPC),
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
			"a `Deny` to or from `any` with `port_ranges` not set (all ports and protocols). " +
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
			int64validator.Between(0, awsVPCMaxPriority),
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
//   - port_ranges can be empty (not set), meaning all ports and protocols, only where the backend
//     can enforce that.
//   - A k8s side that uses a feature only one policy API supports must set that policy_api on
//     every clusters entry.
//
// Values that are unknown at plan time are skipped; the backend validates them.
type policyVersionRuleValidator struct{}

var _ validator.Object = policyVersionRuleValidator{}

func (policyVersionRuleValidator) Description(_ context.Context) string {
	return "rule must not use `any` on both sides, must leave port_ranges unset only where supported, " +
		"and must use only features supported by the policy_api of every clusters entry"
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

	r.checkAnyOnBothSides()  // source.any and destination.any can't both be set
	r.checkEmptyPortRanges() // empty port_ranges isn't supported with cloud or services
	r.checkPriority()        // a non-zero priority needs aws_vpc
	r.checkBaselineTier()    // Baseline, other than a default deny, needs aws_vpc
	r.checkServiceAccounts() // service_accounts needs cilium
	r.checkIngressFromCIDR() // source any or ip_lists to k8s needs cilium on the destination
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

// ruleSides are the sides of a rule.
var ruleSides = []string{"source", "destination"}

// k8sOf returns the attributes of side's k8s object, if side is known to be k8s.
func (r ruleCheck) k8sOf(side string) (map[string]attr.Value, bool) {
	sideAttributes, ok := knownAttributes(r.rule[side])
	if !ok {
		return nil, false
	}

	return knownAttributes(sideAttributes["k8s"])
}

// checkAnyOnBothSides rejects any on both sides.
func (r ruleCheck) checkAnyOnBothSides() {
	source, okSource := knownAttributes(r.rule["source"])
	destination, okDestination := knownAttributes(r.rule["destination"])

	if okSource && okDestination && isSet(source["any"]) && isSet(destination["any"]) {
		r.addError(r.path.AtName("destination").AtName("any"), "destination",
			"`any` is not valid when the source is also `any`.")
	}
}

// checkEmptyPortRanges rejects empty (unset) port_ranges, meaning all ports and protocols, where the
// backend can't enforce it: cloud security rules need a protocol and ports, and on a rule targeting
// services it would cover every port of the services' pods rather than the services' own ports.
func (r ruleCheck) checkEmptyPortRanges() {
	if !isNull(r.rule["port_ranges"]) {
		return
	}

	portRangesPath := r.path.AtName("port_ranges")

	for _, side := range ruleSides {
		if sideAttributes, ok := knownAttributes(r.rule[side]); ok && isSet(sideAttributes["cloud"]) {
			r.addError(portRangesPath, side+".cloud",
				"port_ranges must be set for cloud resources.")
		}
	}

	if k8s, ok := r.k8sOf("destination"); ok && isSet(k8s["services"]) {
		r.addError(portRangesPath, "destination.k8s",
			"port_ranges must be set when targeting services.")
	}
}

// checkPriority requires aws_vpc on every clusters entry of a k8s side with a non-zero priority.
// Cilium has no priority.
func (r ruleCheck) checkPriority() {
	for _, side := range ruleSides {
		k8s, ok := r.k8sOf(side)
		if !ok {
			continue
		}

		if priority, ok := knownInt64(k8s["priority"]); ok && priority > 0 {
			r.requireAPI(r.path.AtName(side).AtName("k8s").AtName("priority"), side+".k8s",
				"a non-zero priority", k8s["clusters"], policyAPIAWSVPC)
		}
	}
}

// checkBaselineTier requires aws_vpc on every clusters entry of a k8s side in the Baseline tier,
// unless the rule is a default deny, which Cilium also supports.
func (r ruleCheck) checkBaselineTier() {
	for _, side := range ruleSides {
		k8s, ok := r.k8sOf(side)
		if !ok {
			continue
		}

		if tier, ok := knownString(k8s["tier"]); ok && tier == tierBaseline && !r.isDefaultDeny(side) {
			r.requireAPI(r.path.AtName(side).AtName("k8s").AtName("tier"), side+".k8s",
				fmt.Sprintf("tier %q, other than for a default deny (a Deny to or from `any` with port_ranges not set),", tierBaseline),
				k8s["clusters"], policyAPIAWSVPC)
		}
	}
}

// checkServiceAccounts requires cilium on every clusters entry of a k8s side that selects
// workloads by service account. AWS ClusterNetworkPolicy can't select by service account.
func (r ruleCheck) checkServiceAccounts() {
	for _, side := range ruleSides {
		k8s, ok := r.k8sOf(side)
		if !ok {
			continue
		}

		if hasServiceAccounts(k8s) {
			r.requireAPI(r.path.AtName(side).AtName("k8s").AtName("workload_selector").AtName("service_accounts"), side+".k8s",
				"workload_selector.service_accounts", k8s["clusters"], policyAPICilium)
		}
	}
}

// checkIngressFromCIDR checks a rule from source any or ip_lists to a k8s destination. The
// destination's clusters render that rule as ingress from a CIDR, which AWS ClusterNetworkPolicy
// cannot express, so every destination clusters entry must select Cilium clusters.
func (r ruleCheck) checkIngressFromCIDR() {
	source, ok := knownAttributes(r.rule["source"])
	if !ok {
		return
	}

	k8s, ok := r.k8sOf("destination")
	if !ok {
		return
	}

	for _, name := range []string{"any", "ip_lists"} {
		if isSet(source[name]) {
			r.requireAPI(r.path.AtName("source").AtName(name), "destination.k8s", "ingress from source."+name,
				k8s["clusters"], policyAPICilium)
		}
	}
}

// requireAPI reports an error unless every clusters entry sets policy_api to one of apis, the policy
// APIs that support feature. Entries whose policy_api is unknown at plan time are skipped.
func (r ruleCheck) requireAPI(attrPath path.Path, sideName, feature string, clusters attr.Value, apis ...string) {
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
		if value, _ := knownString(entry["policy_api"]); !slices.Contains(apis, value) {
			required := fmt.Sprintf("%q", apis[0])
			if len(apis) > 1 {
				required = fmt.Sprintf("to be one of %q", apis)
			}

			r.addError(attrPath, sideName, fmt.Sprintf("%s requires policy_api %s on every clusters entry.", feature, required))

			return
		}
	}
}

// isDefaultDeny reports whether the rule denies all traffic between the k8s side and any peer: a
// Deny to destination any (egress) or from source any (ingress), with port_ranges not set. Both
// policy APIs render this shape in the Baseline tier, so it is exempt from the Baseline
// requirement; the ingress form still needs cilium for its source.any peer, which
// checkIngressFromCIDR checks. Unknown values count as a match, leaving them to the backend.
func (r ruleCheck) isDefaultDeny(side string) bool {
	if action, ok := knownString(r.rule["action"]); ok && action != "Deny" {
		return false
	}

	if peer, ok := knownAttributes(r.rule[peerOf(side)]); ok && isNull(peer["any"]) {
		return false
	}

	// An empty list is rejected on its own by the port_ranges validators, so only listed port
	// ranges make this not a default deny.
	if portRanges, ok := r.rule["port_ranges"].(types.List); ok && isSet(portRanges) && len(portRanges.Elements()) > 0 {
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
