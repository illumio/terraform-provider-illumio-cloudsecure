// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// policyVersionRuleType returns the object type of one policy version rule.
func policyVersionRuleType() types.ObjectType {
	rules, ok := policyVersionResource.Schema.Attributes["rules"].(MapNestedResourceAttributeWithMode)
	if !ok {
		panic("rules is not a MapNestedResourceAttributeWithMode")
	}

	ruleType, ok := rules.NestedObject.Type().(types.ObjectType)
	if !ok {
		panic("rules element is not an ObjectType")
	}

	return ruleType
}

// toAttrValue builds a value of type t from plain Go values. Object attributes missing from a
// map are null, as Terraform does for omitted arguments.
func toAttrValue(t attr.Type, v any) attr.Value {
	switch t := t.(type) {
	case types.ObjectType:
		if v == nil {
			return types.ObjectNull(t.AttrTypes)
		}

		m, _ := v.(map[string]any)
		attrs := make(map[string]attr.Value, len(t.AttrTypes))

		for name, attrType := range t.AttrTypes {
			attrs[name] = toAttrValue(attrType, m[name])
		}

		return types.ObjectValueMust(t.AttrTypes, attrs)
	case types.ListType:
		if v == nil {
			return types.ListNull(t.ElemType)
		}

		l, _ := v.([]any)
		elems := make([]attr.Value, 0, len(l))

		for _, elem := range l {
			elems = append(elems, toAttrValue(t.ElemType, elem))
		}

		return types.ListValueMust(t.ElemType, elems)
	case basetypes.StringType:
		if v == nil {
			return types.StringNull()
		}

		s, _ := v.(string)

		return types.StringValue(s)
	case basetypes.Int64Type:
		if v == nil {
			return types.Int64Null()
		}

		i, _ := v.(int)

		return types.Int64Value(int64(i))
	default:
		panic(fmt.Sprintf("unsupported type %s", t))
	}
}

func diagnosticDetails(diags diag.Diagnostics) []string {
	details := make([]string, 0, len(diags))

	for _, d := range diags.Errors() {
		details = append(details, d.Detail())
	}

	return details
}

// testK8sSide returns a k8s side for one cluster with extra attributes merged in.
// policy_apis in extra is set on the first clusters entry, where the schema has it.
func testK8sSide(extra map[string]any) map[string]any {
	cluster := map[string]any{"id": "7f8f6b1e-5f0e-4d8c-9a64-3a8f2b6f1c11"}
	k8s := map[string]any{
		"clusters":           []any{cluster},
		"namespace_selector": map[string]any{},
		"workload_selector":  map[string]any{},
	}

	maps.Copy(k8s, extra)

	if apis, ok := k8s["policy_apis"]; ok {
		cluster["policy_apis"] = apis

		delete(k8s, "policy_apis")
	}

	return map[string]any{"k8s": k8s}
}

var (
	testAny      = map[string]any{"any": map[string]any{}}
	testAllPorts = []any{map[string]any{"protocol": "TCP", "from_port": 0, "to_port": 0}, map[string]any{"protocol": "UDP", "from_port": 0, "to_port": 0}}
	testTCP443   = []any{map[string]any{"protocol": "TCP", "from_port": 443, "to_port": 443}}
	testCilium   = map[string]any{"policy_apis": []any{"cilium"}, "tier": "Baseline"}
)

func (suite *SchemaTestSuite) TestPolicyVersionRuleValidator() {
	ruleType := policyVersionRuleType()

	baselineMessage := `policy API "cilium" supports tier "Baseline" only for a Deny %s %s = { any = {} } on all TCP and UDP ports; %s`
	egressBaseline := func(reason string) string { return fmt.Sprintf(baselineMessage, "to", "destination", reason) }
	ingressBaseline := func(reason string) string { return fmt.Sprintf(baselineMessage, "from", "source", reason) }
	allPortsReason := "port_ranges must be exactly TCP 0-0 and UDP 0-0 (all ports)."

	tests := map[string]struct {
		action      string
		source      map[string]any
		destination map[string]any
		portRanges  []any
		wantErrors  []string // substrings; nil means no errors
	}{
		// Cilium Baseline.
		"cilium Baseline Deny to any on all ports passes": {
			action: "Deny", source: testK8sSide(testCilium), destination: testAny, portRanges: testAllPorts,
		},
		"cilium Baseline Allow fails": {
			action: "Allow", source: testK8sSide(testCilium), destination: testAny, portRanges: testAllPorts,
			wantErrors: []string{`Rule "r", source.k8s: ` + egressBaseline(`action is "Allow".`)},
		},
		"cilium Baseline to a non-any destination fails": {
			action: "Deny", source: testK8sSide(testCilium), destination: testK8sSide(nil), portRanges: testAllPorts,
			wantErrors: []string{egressBaseline("the destination is not `any`.")},
		},
		"cilium Baseline on 1-65535 fails": {
			action: "Deny", source: testK8sSide(testCilium), destination: testAny,
			portRanges: []any{map[string]any{"protocol": "TCP", "from_port": 1, "to_port": 65535}, map[string]any{"protocol": "UDP", "from_port": 1, "to_port": 65535}},
			wantErrors: []string{egressBaseline(allPortsReason)},
		},
		"cilium Baseline on TCP 0-0 only fails": {
			action: "Deny", source: testK8sSide(testCilium), destination: testAny,
			portRanges: []any{map[string]any{"protocol": "TCP", "from_port": 0, "to_port": 0}},
			wantErrors: []string{egressBaseline(allPortsReason)},
		},
		"cilium Baseline with no port ranges fails": {
			action: "Deny", source: testK8sSide(testCilium), destination: testAny, portRanges: []any{},
			wantErrors: []string{egressBaseline(allPortsReason)},
		},
		"cilium Baseline accepts lowercase protocols": {
			action: "Deny", source: testK8sSide(testCilium), destination: testAny,
			portRanges: []any{map[string]any{"protocol": "tcp", "from_port": 0, "to_port": 0}, map[string]any{"protocol": "udp", "from_port": 0, "to_port": 0}},
		},
		"the strictest listed API wins": {
			action: "Allow", source: testK8sSide(map[string]any{"policy_apis": []any{"aws_vpc", "cilium"}, "tier": "Baseline"}),
			destination: testAny, portRanges: testAllPorts,
			wantErrors: []string{egressBaseline(`action is "Allow".`)},
		},

		"cilium Baseline Deny from any on all ports passes": {
			action: "Deny", source: testAny, destination: testK8sSide(testCilium), portRanges: testAllPorts,
		},
		"cilium Baseline on the destination needs source any": {
			action: "Deny", source: testK8sSide(nil), destination: testK8sSide(testCilium), portRanges: testAllPorts,
			wantErrors: []string{`Rule "r", destination.k8s: ` + ingressBaseline("the source is not `any`.")},
		},
		"policy_apis is the union across clusters entries": {
			action: "Allow",
			source: map[string]any{"k8s": map[string]any{
				"clusters": []any{
					map[string]any{"id": "7f8f6b1e-5f0e-4d8c-9a64-3a8f2b6f1c11", "policy_apis": []any{"aws_vpc"}},
					map[string]any{"id": "0c7d3b5a-9e41-4f2a-8b6d-1a2e3f4c5d6e", "policy_apis": []any{"cilium"}},
				},
				"namespace_selector": map[string]any{},
				"workload_selector":  map[string]any{},
				"tier":               "Baseline",
			}},
			destination: testAny, portRanges: testAllPorts,
			wantErrors: []string{egressBaseline(`action is "Allow".`)},
		},

		// Sides are checked independently.
		"K8S to K8S with different tiers passes": {
			action:      "Allow",
			source:      testK8sSide(map[string]any{"policy_apis": []any{"cilium"}, "tier": "Admin"}),
			destination: testK8sSide(map[string]any{"policy_apis": []any{"aws_vpc"}, "tier": "Baseline"}),
			portRanges:  testTCP443,
		},
		"without policy_apis a Baseline Allow is left to the backend": {
			action: "Allow", source: testK8sSide(map[string]any{"tier": "Baseline"}),
			destination: testK8sSide(map[string]any{"tier": "Baseline"}), portRanges: testTCP443,
		},
		"aws_vpc allows a Baseline Allow": {
			action: "Allow", source: testK8sSide(map[string]any{"policy_apis": []any{"aws_vpc"}, "tier": "Baseline"}),
			destination: testAny, portRanges: testTCP443,
		},

		// Tier.
		"cilium allows no tier": {
			action: "Allow", source: testK8sSide(map[string]any{"policy_apis": []any{"cilium"}}),
			destination: testK8sSide(nil), portRanges: testTCP443,
		},
		"aws_vpc requires a tier": {
			action: "Allow", source: testK8sSide(nil),
			destination: testK8sSide(map[string]any{"policy_apis": []any{"aws_vpc"}}), portRanges: testTCP443,
			wantErrors: []string{`Rule "r", destination.k8s: policy API "aws_vpc" requires a tier; set tier to "Admin" or "Baseline". Only "cilium" supports no tier.`},
		},
		"aws_vpc with cilium still requires a tier": {
			action: "Allow", source: testK8sSide(map[string]any{"policy_apis": []any{"cilium", "aws_vpc"}}),
			destination: testK8sSide(nil), portRanges: testTCP443,
			wantErrors: []string{`Rule "r", source.k8s: policy API "aws_vpc" requires a tier`},
		},

		// Priority.
		"aws_vpc accepts priority 1000": {
			action: "Allow", source: testK8sSide(map[string]any{"policy_apis": []any{"aws_vpc"}, "tier": "Admin", "priority": 1000}),
			destination: testK8sSide(nil), portRanges: testTCP443,
		},
		"aws_vpc rejects priority 1001": {
			action: "Allow", source: testK8sSide(map[string]any{"policy_apis": []any{"aws_vpc"}, "tier": "Admin", "priority": 1001}),
			destination: testK8sSide(nil), portRanges: testTCP443,
			wantErrors: []string{`Rule "r", source.k8s: policy API "aws_vpc" supports priority 0-1000, got: 1001.`},
		},
		"without policy_apis priority above 1000 is left to the backend": {
			action: "Allow", source: testK8sSide(map[string]any{"tier": "Admin", "priority": 5000}),
			destination: testK8sSide(nil), portRanges: testTCP443,
		},
		"cilium rejects priority": {
			action: "Allow", source: testK8sSide(map[string]any{"policy_apis": []any{"cilium"}, "priority": 10}),
			destination: testK8sSide(nil), portRanges: testTCP443,
			wantErrors: []string{`Rule "r", source.k8s: policy API "cilium" does not support priority.`},
		},

		"cilium treats priority 0 as unset": {
			action: "Allow", source: testK8sSide(map[string]any{"policy_apis": []any{"cilium"}, "priority": 0}),
			destination: testK8sSide(nil), portRanges: testTCP443,
		},

		// aws_vpc has no CIDR ingress peer.
		"aws_vpc rejects ingress from source any": {
			action: "Deny", source: testAny,
			destination: testK8sSide(map[string]any{"policy_apis": []any{"aws_vpc"}, "tier": "Admin"}), portRanges: testTCP443,
			wantErrors: []string{`Rule "r", destination.k8s: policy API "aws_vpc" cannot enforce ingress from source.any: AWS ClusterNetworkPolicy has no CIDR ingress peer.`},
		},
		"aws_vpc rejects ingress from source ip_lists": {
			action: "Allow", source: map[string]any{"ip_lists": map[string]any{"ids": []any{"ipl-1"}}},
			destination: testK8sSide(map[string]any{"policy_apis": []any{"aws_vpc"}, "tier": "Admin"}), portRanges: testTCP443,
			wantErrors: []string{`Rule "r", destination.k8s: policy API "aws_vpc" cannot enforce ingress from source.ip_lists: AWS ClusterNetworkPolicy has no CIDR ingress peer.`},
		},
		"aws_vpc allows egress to any": {
			action: "Deny", source: testK8sSide(map[string]any{"policy_apis": []any{"aws_vpc"}, "tier": "Admin"}),
			destination: testAny, portRanges: testTCP443,
		},

		// any placement.
		"destination any requires a k8s source": {
			action: "Deny", source: map[string]any{"ip_lists": map[string]any{"ids": []any{"ipl-1"}}},
			destination: testAny, portRanges: testTCP443,
			wantErrors: []string{"Rule \"r\", destination: `any` is only valid when the source is k8s."},
		},
		"source any requires a k8s destination": {
			action: "Deny", source: testAny,
			destination: map[string]any{"fqdns": map[string]any{"names": []any{"api.example.com"}}}, portRanges: testTCP443,
			wantErrors: []string{"Rule \"r\", source: `any` is only valid when the destination is k8s."},
		},
		"any to any fails on both sides": {
			action: "Deny", source: testAny, destination: testAny, portRanges: testTCP443,
			wantErrors: []string{
				"Rule \"r\", source: `any` is only valid when the destination is k8s.",
				"Rule \"r\", destination: `any` is only valid when the source is k8s.",
			},
		},
	}

	for name, tc := range tests {
		suite.Run(name, func() {
			rule := toAttrValue(ruleType, map[string]any{
				"action":      tc.action,
				"source":      tc.source,
				"destination": tc.destination,
				"port_ranges": tc.portRanges,
			})

			object, ok := rule.(types.Object)
			suite.Require().True(ok)

			resp := &validator.ObjectResponse{}
			policyVersionRuleValidator{}.ValidateObject(context.Background(), validator.ObjectRequest{
				Path:        path.Root("rules").AtMapKey("r"),
				ConfigValue: object,
			}, resp)

			details := diagnosticDetails(resp.Diagnostics)

			if tc.wantErrors == nil {
				suite.Empty(details)

				return
			}

			for _, want := range tc.wantErrors {
				suite.True(slices.ContainsFunc(details, func(d string) bool { return strings.Contains(d, want) }),
					"missing error %q in %q", want, details)
			}
		})
	}
}
