// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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

// k8sPolicyAttributes are the attributes shared by source.k8s and destination.k8s that control
// how that side of a rule is rendered into Kubernetes network policies.
//
// tier and priority default to "" and 0. The generated converters cannot tell an unset nested
// scalar from its zero value, so an omitted value is sent as "" or 0 and read back the same way.
// The default plans that zero value up front, which keeps apply consistent with the plan and makes
// removing the attribute from config a change. (Optional and Computed alone would instead keep
// the prior value, so removing tier would plan no change.)
//
// policy_apis needs no such treatment: a null list is sent as an empty repeated field and read
// back as null. An empty list would read back as null too, so it is rejected.
var k8sPolicyAttributes = map[string]resource_schema.Attribute{
	"policy_apis": resource_schema.ListAttribute{
		MarkdownDescription: "Kubernetes network policy APIs this side of the rule is expected to be rendered to: `cilium` or `aws_vpc`. " +
			"When set, the rule is checked at plan time against what every listed API supports. When omitted, the backend validates the rule when the policy is provisioned.",
		Optional:    true,
		ElementType: types.StringType,
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
			listvalidator.UniqueValues(),
			listvalidator.ValueStringsAre(stringvalidator.OneOf(policyAPICilium, policyAPIAwsVpc)),
		},
	},
	"tier": resource_schema.StringAttribute{
		MarkdownDescription: "Policy tier for this side of the rule: `Admin` or `Baseline`. When omitted, the side has no tier. " +
			"Only `cilium` supports no tier: when `policy_apis` includes `aws_vpc`, a tier must be set.",
		Optional: true,
		Computed: true,
		Default:  stringdefault.StaticString(""),
		Validators: []validator.String{
			stringvalidator.OneOf(tierAdmin, tierBaseline),
		},
	},
	"priority": resource_schema.Int64Attribute{
		MarkdownDescription: "Priority of this side of the rule within its tier. Must be a non-negative integer. When `policy_apis` includes `aws_vpc`, must be between 0 and 1000. Not supported by `cilium`; leave it unset or 0.",
		Optional:            true,
		Computed:            true,
		Default:             int64default.StaticInt64(0),
		Validators: []validator.Int64{
			int64validator.AtLeast(0),
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
//   - destination.any is only valid with a k8s source.
//   - Each k8s side that sets policy_apis is checked against what every listed API supports.
//     Sides are checked independently: the backend renders each side as its own object, at that
//     side's tier. An omitted tier means no tier, which only cilium supports.
//
// Values that are unknown at plan time are skipped; the backend validates them at provision.
type policyVersionRuleValidator struct{}

var _ validator.Object = policyVersionRuleValidator{}

func (policyVersionRuleValidator) Description(_ context.Context) string {
	return "rule must be supported by the policy APIs listed in source.k8s.policy_apis and destination.k8s.policy_apis"
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

	r.checkAnyDestination()
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

// checkAnyDestination rejects destination.any unless the source is k8s.
func (r ruleCheck) checkAnyDestination() {
	destination, ok := knownAttributes(r.rule["destination"])
	if !ok || !isSet(destination["any"]) {
		return
	}

	source, ok := knownAttributes(r.rule["source"])
	if !ok || !isNull(source["k8s"]) {
		return
	}

	r.addError(r.path.AtName("destination").AtName("any"), "destination",
		"`any` is only valid when the source is k8s.")
}

// checkSide checks one k8s side of the rule against every API in its policy_apis.
func (r ruleCheck) checkSide(side string) {
	sideAttributes, ok := knownAttributes(r.rule[side])
	if !ok {
		return
	}

	k8s, ok := knownAttributes(sideAttributes["k8s"])
	if !ok {
		return
	}

	apis, ok := knownStrings(k8s["policy_apis"])
	if !ok {
		return
	}

	k8sPath := r.path.AtName(side).AtName("k8s")
	sideName := side + ".k8s"

	// An omitted tier means no tier. In config it is null; an unknown tier is left to the backend.
	if tier := k8s["tier"]; slices.Contains(apis, policyAPIAwsVpc) && tier != nil && tier.IsNull() {
		r.addError(k8sPath.AtName("tier"), sideName,
			fmt.Sprintf("policy API %q requires a tier; set tier to %q or %q. Only %q supports no tier.",
				policyAPIAwsVpc, tierAdmin, tierBaseline, policyAPICilium))
	}

	if priority, ok := knownInt64(k8s["priority"]); ok && slices.Contains(apis, policyAPIAwsVpc) && priority > awsVpcMaxPriority {
		r.addError(k8sPath.AtName("priority"), sideName,
			fmt.Sprintf("policy API %q supports priority 0-%d, got: %d.", policyAPIAwsVpc, awsVpcMaxPriority, priority))
	}

	if !slices.Contains(apis, policyAPICilium) {
		return
	}

	if tier, ok := knownString(k8s["tier"]); ok && tier == tierBaseline {
		r.checkCiliumBaseline(k8sPath.AtName("tier"), side)
	}

	// 0 is what an omitted priority plans to, so it is treated as unset.
	if priority, ok := knownInt64(k8s["priority"]); ok && priority != 0 {
		r.addError(k8sPath.AtName("priority"), sideName,
			fmt.Sprintf("policy API %q does not support priority.", policyAPICilium))
	}
}

// checkCiliumBaseline checks that a Baseline tier k8s side is a deny of all TCP and UDP traffic
// from the source to any destination, which is the only Baseline rule Cilium can render. The k8s
// side must therefore be the source, and all ports must be written as TCP 0-0 and UDP 0-0.
func (r ruleCheck) checkCiliumBaseline(tierPath path.Path, side string) {
	sideName := side + ".k8s"
	prefix := fmt.Sprintf("policy API %q supports tier %q only on source.k8s, for a Deny to destination = { any = {} } on all TCP and UDP ports",
		policyAPICilium, tierBaseline)

	if side != "source" {
		r.addError(tierPath, sideName, prefix+".")

		return
	}

	if action, ok := knownString(r.rule["action"]); ok && action != "Deny" {
		r.addError(tierPath, sideName, fmt.Sprintf("%s; action is %q.", prefix, action))
	}

	if destination, ok := knownAttributes(r.rule["destination"]); ok && isNull(destination["any"]) {
		r.addError(tierPath, sideName, prefix+"; the destination is not `any`.")
	}

	if allPorts, known := isAllPortsTCPAndUDP(r.rule["port_ranges"]); known && !allPorts {
		r.addError(tierPath, sideName, prefix+`; port_ranges must be exactly TCP 0-0 and UDP 0-0 (all ports).`)
	}
}

// isAllPortsTCPAndUDP reports whether the port ranges are exactly "all TCP and UDP ports" in the
// form Cilium Baseline rules accept: every range is 0-0, which the backend renders as all ports,
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

// knownStrings returns the elements of v if it is a known, non-null list of known strings.
func knownStrings(v attr.Value) ([]string, bool) {
	if !isSet(v) {
		return nil, false
	}

	list, ok := v.(types.List)
	if !ok {
		return nil, false
	}

	out := make([]string, 0, len(list.Elements()))

	for _, element := range list.Elements() {
		s, ok := knownString(element)
		if !ok {
			return nil, false
		}

		out = append(out, s)
	}

	return out, true
}
