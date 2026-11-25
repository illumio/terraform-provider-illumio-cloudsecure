// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// validateIPList runs v on a list of elements and returns the summary and detail of the first error.
func (suite *SchemaTestSuite) validateIPList(v validator.List, name string, elemType types.ObjectType, elements []any) (string, string) {
	list, ok := toAttrValue(types.ListType{ElemType: elemType}, elements).(types.List)
	suite.Require().True(ok)

	resp := &validator.ListResponse{}
	v.ValidateList(context.Background(), validator.ListRequest{Path: path.Root(name), ConfigValue: list}, resp)

	if !resp.Diagnostics.HasError() {
		return "", ""
	}

	return resp.Diagnostics.Errors()[0].Summary(), resp.Diagnostics.Errors()[0].Detail()
}

func (suite *SchemaTestSuite) validateIPAddress(ip string) (string, string) {
	return suite.validateIPList(ipAddressesValidator{}, "ip_addresses", IPAddress, []any{
		map[string]any{"ip_address": ip, "exclusion": false},
	})
}

func (suite *SchemaTestSuite) validateIPRange(from, to any) (string, string) {
	return suite.validateIPList(ipRangesValidator{}, "ip_ranges", IPRange, []any{
		map[string]any{"from_ip_address": from, "to_ip_address": to, "exclusion": false},
	})
}

// TestIPAddressesValidatorAcceptsCanonicalForms covers the forms CloudSecure stores as written.
func (suite *SchemaTestSuite) TestIPAddressesValidatorAcceptsCanonicalForms() {
	for _, ip := range []string{"10.20.96.91", "10.20.0.0/16", "0.0.0.0/0", "2001:db8::1", "2001:db8::/64", "::/0", "::ffff:10.0.0.1"} {
		summary, detail := suite.validateIPAddress(ip)
		suite.Empty(summary+detail, "ip_address %q", ip)
	}
}

// TestIPAddressesValidatorRejectsNonCanonicalForms covers the forms CloudSecure rewrites, as
// observed on dev. Each error states the rule, the reason, and the value to write.
func (suite *SchemaTestSuite) TestIPAddressesValidatorRejectsNonCanonicalForms() {
	for ip, want := range map[string]struct{ canonical, rule string }{
		"10.20.96.91/32":  {"10.20.96.91", "single IPs have no prefix length."},
		"2001:db8::1/128": {"2001:db8::1", "single IPs have no prefix length."},
		"10.20.0.1/16":    {"10.20.0.0/16", "CIDRs use their network address."},
		"2001:0db8:0000:0000:0000:0000:0000:0002": {"2001:db8::2", "IPv6 is compressed and lowercase."},
		"2001:DB8::1":   {"2001:db8::1", "IPv6 is compressed and lowercase."},
		"2001:DB8::/64": {"2001:db8::/64", "IPv6 is compressed and lowercase."},
	} {
		summary, detail := suite.validateIPAddress(ip)
		suite.Equal("Non-Canonical IP Address", summary, "ip_address %q", ip)
		suite.Equal(`Write "`+want.canonical+`" instead of "`+ip+`": `+want.rule, detail, "ip_address %q", ip)
	}
}

func (suite *SchemaTestSuite) TestIPAddressesValidatorRejectsInvalidValues() {
	for _, ip := range []string{"not-an-ip", "", "10.0.0.256", "010.0.0.1", " 10.0.0.1", "10.0.0.0/33", "10.0.0.1/", "fe80::1%eth0", "fe80::/64%eth0"} {
		summary, detail := suite.validateIPAddress(ip)
		suite.Equal("Invalid IP Address", summary, "ip_address %q", ip)
		suite.Equal(`"`+ip+`" is not a valid IP address or CIDR.`, detail, "ip_address %q", ip)
	}
}

func (suite *SchemaTestSuite) TestIPRangesValidatorAcceptsCanonicalRanges() {
	for _, r := range [][2]string{{"10.0.0.1", "10.0.0.9"}, {"10.0.0.1", "10.0.0.1"}, {"2001:db8::1", "2001:db8::5"}} {
		summary, detail := suite.validateIPRange(r[0], r[1])
		suite.Empty(summary+detail, "range %q", r)
	}
}

func (suite *SchemaTestSuite) TestIPRangesValidatorRejectsInvalidRanges() {
	for _, tc := range []struct {
		from, to string
		summary  string
		detail   string
	}{
		// Rewritten by CloudSecure, as observed on dev.
		{"10.0.0.1/32", "10.0.0.9", "Non-Canonical IP Address", `Write "10.0.0.1" instead of "10.0.0.1/32"`},
		{"2001:0db8:0000:0000:0000:0000:0000:0001", "2001:db8::5", "Non-Canonical IP Address", `Write "2001:db8::1" instead of`},
		{"2001:db8::1", "2001:DB8::5", "Non-Canonical IP Address", `Write "2001:db8::5" instead of "2001:DB8::5"`},
		// A range end is a single IP.
		{"10.0.0.0/24", "10.0.1.9", "CIDR as IP Range End",
			`Range ends must be single IPs. Write "10.0.0.0" instead of "10.0.0.0/24", or put "10.0.0.0/24" in ip_addresses.`},
		{"10.0.0.1", "10.0.1.0/24", "CIDR as IP Range End",
			`Range ends must be single IPs. Write "10.0.1.255" instead of "10.0.1.0/24", or put "10.0.1.0/24" in ip_addresses.`},
		{"10.0.0.1/24", "10.0.1.9", "CIDR as IP Range End", `put "10.0.0.0/24" in ip_addresses`},
		{"2001:db8::/64", "2001:db8:0:1::", "CIDR as IP Range End", `Write "2001:db8::" instead of "2001:db8::/64"`},
		{"2001:db8::1", "2001:db8:0:1::/64", "CIDR as IP Range End", `Write "2001:db8:0:1:ffff:ffff:ffff:ffff" instead of`},
		{"not-an-ip", "10.0.0.9", "Invalid IP Address", `"not-an-ip" is not a valid IP address.`},
		{"fe80::1%eth0", "fe80::9", "Invalid IP Address", `"fe80::1%eth0" is not a valid IP address.`},
		{"10.0.0.1", "2001:db8::1", "Mixed IP Versions in IP Range", "must both be IPv4 or both IPv6"},
		{"10.0.0.9", "10.0.0.1", "Invalid IP Range",
			"from_ip_address must not be greater than to_ip_address."},
	} {
		summary, detail := suite.validateIPRange(tc.from, tc.to)
		suite.Equal(tc.summary, summary, "range %q-%q", tc.from, tc.to)
		suite.Contains(detail, tc.detail, "range %q-%q", tc.from, tc.to)
	}
}

func (suite *SchemaTestSuite) TestIPListValidatorsRejectEmptyLists() {
	summary, detail := suite.validateIPList(ipAddressesValidator{}, "ip_addresses", IPAddress, []any{})
	suite.Equal("Empty IP List Attribute", summary)
	suite.Equal("Omit ip_addresses instead of setting it to [].", detail)

	summary, detail = suite.validateIPList(ipRangesValidator{}, "ip_ranges", IPRange, []any{})
	suite.Equal("Empty IP List Attribute", summary)
	suite.Contains(detail, "Omit ip_ranges instead of setting it to [].")
}

func (suite *SchemaTestSuite) TestIPListValidatorsSkipUnsetValues() {
	summary, detail := suite.validateIPList(ipAddressesValidator{}, "ip_addresses", IPAddress, nil)
	suite.Empty(summary + detail)

	summary, detail = suite.validateIPList(ipAddressesValidator{}, "ip_addresses", IPAddress, []any{
		map[string]any{"ip_address": nil, "exclusion": false},
	})
	suite.Empty(summary + detail)

	summary, detail = suite.validateIPRange("10.0.0.9", nil)
	suite.Empty(summary + detail)
}
