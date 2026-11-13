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
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
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

// testUnknown stands for a string that is unknown at plan time.
type testUnknown struct{}

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

		if _, ok := v.(testUnknown); ok {
			return types.StringUnknown()
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
// policy_api in extra is set on the clusters entry, where the schema has it.
func testK8sSide(extra map[string]any) map[string]any {
	rest := maps.Clone(extra)
	delete(rest, "policy_api")

	side := testK8sEntries(rest, "")
	if api, ok := extra["policy_api"]; ok {
		k8s, _ := side["k8s"].(map[string]any)
		clusters, _ := k8s["clusters"].([]any)
		cluster, _ := clusters[0].(map[string]any)
		cluster["policy_api"] = api
	}

	return side
}

// testK8sEntries returns a k8s side with extra attributes merged in and one clusters entry per
// policy_api; "" leaves it unset.
func testK8sEntries(extra map[string]any, apis ...string) map[string]any {
	clusters := make([]any, 0, len(apis))

	for i, api := range apis {
		cluster := map[string]any{"id": fmt.Sprintf("7f8f6b1e-5f0e-4d8c-9a64-3a8f2b6f1c1%d", i)}
		if api != "" {
			cluster["policy_api"] = api
		}

		clusters = append(clusters, cluster)
	}

	k8s := map[string]any{
		"clusters":           clusters,
		"namespace_selector": map[string]any{},
		"workload_selector":  map[string]any{},
	}
	maps.Copy(k8s, extra)

	return map[string]any{"k8s": k8s}
}

var (
	testAny      = map[string]any{"any": map[string]any{}}
	testAllPorts = []any{map[string]any{"protocol": "TCP", "from_port": 0, "to_port": 0}, map[string]any{"protocol": "UDP", "from_port": 0, "to_port": 0}}
	testTCP443   = []any{map[string]any{"protocol": "TCP", "from_port": 443, "to_port": 443}}
	testCilium   = map[string]any{"policy_api": "cilium", "tier": "Baseline"}
)

func (suite *SchemaTestSuite) TestPolicyVersionRuleValidator() {
	ruleType := policyVersionRuleType()

	baselineNeedsAWSVPC := `tier "Baseline", other than for a default deny (a Deny to or from ` + "`any`" +
		` with empty port_ranges, or TCP 0-0 and UDP 0-0), requires policy_api "aws_vpc" on every clusters entry.`
	priorityNeedsAWSVPC := `a non-zero priority requires policy_api "aws_vpc" on every clusters entry.`
	needsCilium := func(feature string) string {
		return feature + ` requires policy_api "cilium" on every clusters entry.`
	}

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
		"Baseline Allow with cilium fails": {
			action: "Allow", source: testK8sSide(testCilium), destination: testAny, portRanges: testAllPorts,
			wantErrors: []string{`Rule "r", source.k8s: ` + baselineNeedsAWSVPC},
		},
		"cilium Baseline to a non-any destination fails": {
			action: "Deny", source: testK8sSide(testCilium), destination: testK8sSide(nil), portRanges: testAllPorts,
			wantErrors: []string{baselineNeedsAWSVPC},
		},
		"cilium Baseline on 1-65535 fails": {
			action: "Deny", source: testK8sSide(testCilium), destination: testAny,
			portRanges: []any{map[string]any{"protocol": "TCP", "from_port": 1, "to_port": 65535}, map[string]any{"protocol": "UDP", "from_port": 1, "to_port": 65535}},
			wantErrors: []string{baselineNeedsAWSVPC},
		},
		"cilium Baseline on TCP 0-0 only fails": {
			action: "Deny", source: testK8sSide(testCilium), destination: testAny,
			portRanges: []any{map[string]any{"protocol": "TCP", "from_port": 0, "to_port": 0}},
			wantErrors: []string{baselineNeedsAWSVPC},
		},
		"cilium Baseline with empty port_ranges passes": {
			action: "Deny", source: testK8sSide(testCilium), destination: testAny, portRanges: []any{},
		},
		"cilium Baseline accepts lowercase protocols": {
			action: "Deny", source: testK8sSide(testCilium), destination: testAny,
			portRanges: []any{map[string]any{"protocol": "tcp", "from_port": 0, "to_port": 0}, map[string]any{"protocol": "udp", "from_port": 0, "to_port": 0}},
		},

		// Ingress Baseline and multiple clusters entries.
		"cilium Baseline Deny from any on all ports passes": {
			action: "Deny", source: testAny, destination: testK8sSide(testCilium), portRanges: testAllPorts,
		},
		"cilium Baseline on the destination needs source any": {
			action: "Deny", source: testK8sSide(nil), destination: testK8sSide(testCilium), portRanges: testAllPorts,
			wantErrors: []string{`Rule "r", destination.k8s: ` + baselineNeedsAWSVPC},
		},
		"Baseline Allow fails if any clusters entry is cilium": {
			action: "Allow", source: testK8sEntries(map[string]any{"tier": "Baseline"}, "aws_vpc", "cilium"),
			destination: testAny, portRanges: testAllPorts,
			wantErrors: []string{baselineNeedsAWSVPC},
		},
		"Baseline default deny passes with both APIs": {
			action: "Deny", source: testK8sEntries(map[string]any{"tier": "Baseline"}, "cilium", "aws_vpc"),
			destination: testAny, portRanges: testAllPorts,
		},

		// Sides are checked independently.
		"K8S to K8S with different tiers passes": {
			action:      "Allow",
			source:      testK8sSide(map[string]any{"policy_api": "cilium", "tier": "Admin"}),
			destination: testK8sSide(map[string]any{"policy_api": "aws_vpc", "tier": "Baseline"}),
			portRanges:  testTCP443,
		},
		"aws_vpc allows a Baseline Allow": {
			action: "Allow", source: testK8sSide(map[string]any{"policy_api": "aws_vpc", "tier": "Baseline"}),
			destination: testAny, portRanges: testTCP443,
		},

		// Tier.
		"an omitted tier passes with both APIs": {
			action: "Allow", source: testK8sEntries(nil, "cilium", "aws_vpc"),
			destination: testK8sSide(nil), portRanges: testTCP443,
		},

		// Priority.
		"aws_vpc accepts priority 1000": {
			action: "Allow", source: testK8sSide(map[string]any{"policy_api": "aws_vpc", "tier": "Admin", "priority": 1000}),
			destination: testK8sSide(nil), portRanges: testTCP443,
		},
		"priority fails if any clusters entry is cilium": {
			action: "Allow", source: testK8sEntries(map[string]any{"priority": 5}, "cilium", "aws_vpc"),
			destination: testK8sSide(nil), portRanges: testTCP443,
			wantErrors: []string{`Rule "r", source.k8s: ` + priorityNeedsAWSVPC},
		},
		"priority with cilium fails": {
			action: "Allow", source: testK8sSide(map[string]any{"policy_api": "cilium", "priority": 10}),
			destination: testK8sSide(nil), portRanges: testTCP443,
			wantErrors: []string{`Rule "r", source.k8s: ` + priorityNeedsAWSVPC},
		},

		"priority 0 passes with cilium": {
			action: "Allow", source: testK8sSide(map[string]any{"policy_api": "cilium", "priority": 0}),
			destination: testK8sSide(nil), portRanges: testTCP443,
		},

		// policy_api is required for features only one API supports.
		"priority without policy_api fails": {
			action: "Allow", source: testK8sSide(map[string]any{"priority": 5}),
			destination: testK8sSide(nil), portRanges: testTCP443,
			wantErrors: []string{`Rule "r", source.k8s: ` + priorityNeedsAWSVPC},
		},
		"Baseline Allow without policy_api fails": {
			action: "Allow", source: testK8sSide(map[string]any{"tier": "Baseline"}),
			destination: testAny, portRanges: testTCP443,
			wantErrors: []string{`Rule "r", source.k8s: ` + baselineNeedsAWSVPC},
		},
		"Baseline default deny without policy_api passes": {
			action: "Deny", source: testK8sSide(map[string]any{"tier": "Baseline"}),
			destination: testAny, portRanges: testAllPorts,
		},
		"ingress from ip_lists without policy_api fails": {
			action: "Allow", source: map[string]any{"ip_lists": map[string]any{"ids": []any{"ipl-1"}}},
			destination: testK8sSide(nil), portRanges: testTCP443,
			wantErrors: []string{`Rule "r", destination.k8s: ` + needsCilium("ingress from source.ip_lists")},
		},
		"priority fails if any clusters entry has no policy_api": {
			action: "Allow", source: testK8sEntries(map[string]any{"priority": 5}, "aws_vpc", ""),
			destination: testK8sSide(nil), portRanges: testTCP443,
			wantErrors: []string{`Rule "r", source.k8s: ` + priorityNeedsAWSVPC},
		},
		"service_accounts without policy_api fails": {
			action: "Allow", source: testK8sSide(map[string]any{"workload_selector": map[string]any{"service_accounts": []any{"api"}}}),
			destination: testK8sSide(nil), portRanges: testTCP443,
			wantErrors: []string{`Rule "r", source.k8s: ` + needsCilium("workload_selector.service_accounts")},
		},
		"service_accounts with aws_vpc fails": {
			action: "Allow", source: testK8sSide(nil),
			destination: testK8sSide(map[string]any{"policy_api": "aws_vpc", "workload_selector": map[string]any{"service_accounts": []any{"api"}}}),
			portRanges:  testTCP443,
			wantErrors:  []string{`Rule "r", destination.k8s: ` + needsCilium("workload_selector.service_accounts")},
		},

		"ingress Baseline default deny with aws_vpc fails": {
			action: "Deny", source: testAny,
			destination: testK8sSide(map[string]any{"policy_api": "aws_vpc", "tier": "Baseline"}), portRanges: testAllPorts,
			wantErrors: []string{`Rule "r", destination.k8s: ` + needsCilium("ingress from source.any")},
		},
		"ingress Baseline default deny without policy_api fails": {
			action: "Deny", source: testAny,
			destination: testK8sSide(map[string]any{"tier": "Baseline"}), portRanges: testAllPorts,
			wantErrors: []string{`Rule "r", destination.k8s: ` + needsCilium("ingress from source.any")},
		},
		"an unknown policy_api is left to the backend": {
			action: "Allow", source: testK8sSide(map[string]any{"policy_api": testUnknown{}, "priority": 5}),
			destination: testK8sSide(nil), portRanges: testTCP443,
		},

		// aws_vpc has no CIDR ingress peer.
		"ingress from source any with aws_vpc fails": {
			action: "Deny", source: testAny,
			destination: testK8sSide(map[string]any{"policy_api": "aws_vpc", "tier": "Admin"}), portRanges: testTCP443,
			wantErrors: []string{`Rule "r", destination.k8s: ` + needsCilium("ingress from source.any")},
		},
		"ingress from source ip_lists with aws_vpc fails": {
			action: "Allow", source: map[string]any{"ip_lists": map[string]any{"ids": []any{"ipl-1"}}},
			destination: testK8sSide(map[string]any{"policy_api": "aws_vpc", "tier": "Admin"}), portRanges: testTCP443,
			wantErrors: []string{`Rule "r", destination.k8s: ` + needsCilium("ingress from source.ip_lists")},
		},
		"aws_vpc allows egress to any": {
			action: "Deny", source: testK8sSide(map[string]any{"policy_api": "aws_vpc", "tier": "Admin"}),
			destination: testAny, portRanges: testTCP443,
		},

		// Empty port_ranges means all ports and protocols.
		"empty port_ranges between k8s workloads passes": {
			action: "Allow", source: testK8sSide(nil), destination: testK8sSide(nil), portRanges: []any{},
		},
		"empty port_ranges with a cloud side fails": {
			action: "Allow", source: testK8sSide(nil),
			destination: map[string]any{"cloud": map[string]any{"aws": map[string]any{}}}, portRanges: []any{},
			wantErrors: []string{`Rule "r", destination.cloud: empty port_ranges (all ports and protocols) is not supported for cloud resources`},
		},
		"empty port_ranges targeting services fails": {
			action: "Allow", source: testK8sSide(nil),
			destination: testK8sSide(map[string]any{"services": []any{"api"}, "workload_selector": nil}), portRanges: []any{},
			wantErrors: []string{`Rule "r", destination.k8s: empty port_ranges (all ports and protocols) is not supported when targeting services`},
		},

		// any placement.
		"destination any with an ip_lists source passes": {
			action: "Deny", source: map[string]any{"ip_lists": map[string]any{"ids": []any{"ipl-1"}}},
			destination: testAny, portRanges: testTCP443,
		},
		"source any with an fqdns destination passes": {
			action: "Deny", source: testAny,
			destination: map[string]any{"fqdns": map[string]any{"names": []any{"api.example.com"}}}, portRanges: testTCP443,
		},
		"any to any fails": {
			action: "Deny", source: testAny, destination: testAny, portRanges: testTCP443,
			wantErrors: []string{"Rule \"r\", destination: `any` is not valid when the source is also `any`."},
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

func (suite *SchemaTestSuite) TestPriorityRange() {
	priority, ok := k8sPolicyAttributes["priority"].(resource_schema.Int64Attribute)
	suite.Require().True(ok)

	for value, wantError := range map[int64]bool{0: false, 1000: false, 1001: true, -1: true} {
		resp := &validator.Int64Response{}
		for _, v := range priority.Validators {
			v.ValidateInt64(context.Background(), validator.Int64Request{ConfigValue: types.Int64Value(value)}, resp)
		}

		suite.Equal(wantError, resp.Diagnostics.HasError(), "priority %d", value)
	}
}

func (suite *SchemaTestSuite) TestPortRangeProtocolValidator() {
	for protocol, wantError := range map[string]bool{"TCP": false, "udp": false, "ICMP": true, "": true} {
		ports := toAttrValue(types.ListType{ElemType: PortRange}, []any{
			map[string]any{"protocol": protocol, "from_port": 443, "to_port": 443},
		})

		list, ok := ports.(types.List)
		suite.Require().True(ok)

		resp := &validator.ListResponse{}
		portRangeProtocolValidator{}.ValidateList(context.Background(), validator.ListRequest{
			Path:        path.Root("port_ranges"),
			ConfigValue: list,
		}, resp)

		suite.Equal(wantError, resp.Diagnostics.HasError(), "protocol %q", protocol)
	}
}
