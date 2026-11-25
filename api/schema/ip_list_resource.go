// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"context"
	"fmt"
	"net/netip"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var IPRange = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"exclusion":       types.BoolType,
		"from_ip_address": types.StringType,
		"to_ip_address":   types.StringType,
	},
}

var IPAddress = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"exclusion":  types.BoolType,
		"ip_address": types.StringType,
	},
}

var (
	ipListResource = Resource{
		TypeName: "ip_list",
		Schema: resource_schema.Schema{
			Version:     1,
			Description: "Manages a list of IP address ranges to define CloudSecure policies.",
			Attributes: map[string]resource_schema.Attribute{
				IDFieldName: idAttribute,
				"description": resource_schema.StringAttribute{
					Description: "Description of the IP list.",
					Optional:    true,
					Computed:    true,
					// The Config API returns "" when no description is set.
					PlanModifiers: []planmodifier.String{
						defaultString(""),
					},
				},
				"ip_addresses": resource_schema.ListAttribute{
					Optional:    true,
					Description: "List of IP addresses.",
					ElementType: IPAddress,
					Validators: []validator.List{
						ipAddressesValidator{},
					},
				},
				"ip_ranges": resource_schema.ListAttribute{
					Optional:    true,
					Description: "List of IP address ranges.",
					ElementType: IPRange,
					Validators: []validator.List{
						ipRangesValidator{},
					},
				},
				"name": resource_schema.StringAttribute{
					Description: "Name of the IP list.",
					Required:    true,
				},
			},
		},
	}
)

// defaultString returns a plan modifier that plans value for an omitted string attribute, which
// must be Optional and Computed. It works like Default: stringdefault.StaticString(value), but a
// Default makes the field non-optional in the generated proto, which changes its Go type from
// *string to string and breaks clients that compile against it.
func defaultString(value string) planmodifier.String {
	return defaultStringModifier{value: value}
}

type defaultStringModifier struct {
	value string
}

func (m defaultStringModifier) Description(_ context.Context) string {
	return fmt.Sprintf("defaults to %q", m.value)
}

func (m defaultStringModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m defaultStringModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.ConfigValue.IsNull() {
		resp.PlanValue = types.StringValue(m.value)
	}
}

// ipFormatError is a violation of the address format, with a summary and a detail that states
// what to write instead or which rule is broken.
type ipFormatError struct {
	summary string
	detail  string
}

// checkIPAddress checks an ip_address: a single IP or a CIDR, in canonical form.
func checkIPAddress(s string) *ipFormatError {
	if !strings.Contains(s, "/") {
		addr, err := netip.ParseAddr(s)
		if err != nil || addr.Zone() != "" { // CloudSecure rejects IPv6 zones, e.g. fe80::1%eth0.
			return &ipFormatError{"Invalid IP Address", fmt.Sprintf("%q is not a valid IP address or CIDR.", s)}
		}

		return nonCanonical(s, addr.String(), "IPv6 is compressed and lowercase.")
	}

	prefix, err := netip.ParsePrefix(s)
	if err != nil {
		return &ipFormatError{"Invalid IP Address", fmt.Sprintf("%q is not a valid IP address or CIDR.", s)}
	}

	switch {
	case prefix.IsSingleIP():
		return nonCanonical(s, prefix.Addr().String(), "single IPs have no prefix length.")
	case prefix.Masked() != prefix:
		return nonCanonical(s, prefix.Masked().String(), "CIDRs use their network address.")
	default:
		return nonCanonical(s, prefix.String(), "IPv6 is compressed and lowercase.")
	}
}

// checkRangeEnd checks an end of an IP range: a single IP in canonical form. isFrom tells which end
// it is, to suggest the CIDR's first or last address. It returns the address when s is a valid
// single IP, canonical or not.
func checkRangeEnd(s string, isFrom bool) (netip.Addr, *ipFormatError) {
	if !strings.Contains(s, "/") {
		addr, err := netip.ParseAddr(s)
		if err != nil || addr.Zone() != "" { // CloudSecure rejects IPv6 zones, e.g. fe80::1%eth0.
			return netip.Addr{}, &ipFormatError{"Invalid IP Address", fmt.Sprintf("%q is not a valid IP address.", s)}
		}

		return addr, nonCanonical(s, addr.String(), "IPv6 is compressed and lowercase.")
	}

	prefix, err := netip.ParsePrefix(s)
	if err != nil {
		return netip.Addr{}, &ipFormatError{"Invalid IP Address", fmt.Sprintf("%q is not a valid IP address.", s)}
	}

	if prefix.IsSingleIP() {
		return prefix.Addr(), nonCanonical(s, prefix.Addr().String(), "single IPs have no prefix length.")
	}

	masked := prefix.Masked()

	suggestion := lastAddr(masked)
	if isFrom {
		suggestion = masked.Addr()
	}

	return netip.Addr{}, &ipFormatError{"CIDR as IP Range End", fmt.Sprintf(
		"Range ends must be single IPs. Write %q instead of %q, or put %q in ip_addresses.",
		suggestion.String(), s, masked.String())}
}

// nonCanonical returns an error if s differs from its canonical form, with the rule it breaks.
func nonCanonical(s, canonical, rule string) *ipFormatError {
	if s == canonical {
		return nil
	}

	return &ipFormatError{"Non-Canonical IP Address", fmt.Sprintf(
		"Write %q instead of %q: %s", canonical, s, rule)}
}

// lastAddr returns the last address of a masked prefix.
func lastAddr(prefix netip.Prefix) netip.Addr {
	b := prefix.Addr().As16()
	offset := 0

	if prefix.Addr().Is4() {
		offset = 96
	}

	for bit := offset + prefix.Bits(); bit < 128; bit++ {
		b[bit/8] |= 1 << (7 - bit%8)
	}

	addr := netip.AddrFrom16(b)
	if prefix.Addr().Is4() {
		addr = addr.Unmap()
	}

	return addr
}

// addEmptyListError adds an error if the list is set to [], which CloudSecure stores as no list.
func addEmptyListError(req validator.ListRequest, resp *validator.ListResponse) bool {
	if len(req.ConfigValue.Elements()) > 0 {
		return false
	}

	resp.Diagnostics.AddAttributeError(req.Path, "Empty IP List Attribute", fmt.Sprintf(
		"Omit %s instead of setting it to [].", req.Path))

	return true
}

// ipAddressesValidator enforces the address format on ip_addresses.
type ipAddressesValidator struct{}

var _ validator.List = ipAddressesValidator{}

func (ipAddressesValidator) Description(_ context.Context) string {
	return "must not be empty, and each ip_address must be a single IP or CIDR in canonical form"
}

func (v ipAddressesValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (ipAddressesValidator) ValidateList(_ context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() || addEmptyListError(req, resp) {
		return
	}

	for i, element := range req.ConfigValue.Elements() {
		ipAddress, ok := knownAttributes(element)
		if !ok {
			continue
		}

		s, ok := knownString(ipAddress["ip_address"])
		if !ok {
			continue
		}

		if e := checkIPAddress(s); e != nil {
			resp.Diagnostics.AddAttributeError(req.Path.AtListIndex(i).AtName("ip_address"), e.summary, e.detail)
		}
	}
}

// ipRangesValidator enforces the address format on ip_ranges.
type ipRangesValidator struct{}

var _ validator.List = ipRangesValidator{}

func (ipRangesValidator) Description(_ context.Context) string {
	return "must not be empty, and each range must be between two single IPs of the same version in canonical form, in ascending order"
}

func (v ipRangesValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (ipRangesValidator) ValidateList(_ context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() || addEmptyListError(req, resp) {
		return
	}

	for i, element := range req.ConfigValue.Elements() {
		ipRange, ok := knownAttributes(element)
		if !ok {
			continue
		}

		ends := make(map[string]netip.Addr, 2)

		for _, name := range []string{"from_ip_address", "to_ip_address"} {
			s, ok := knownString(ipRange[name])
			if !ok {
				continue
			}

			addr, e := checkRangeEnd(s, name == "from_ip_address")
			if e != nil {
				resp.Diagnostics.AddAttributeError(req.Path.AtListIndex(i).AtName(name), e.summary, e.detail)
			}

			if addr.IsValid() {
				ends[name] = addr
			}
		}

		from, fromOK := ends["from_ip_address"]
		to, toOK := ends["to_ip_address"]

		if !fromOK || !toOK {
			continue
		}

		switch {
		case from.Is4() != to.Is4():
			resp.Diagnostics.AddAttributeError(req.Path.AtListIndex(i), "Mixed IP Versions in IP Range",
				"from_ip_address and to_ip_address must both be IPv4 or both IPv6.")
		case from.Compare(to) > 0:
			resp.Diagnostics.AddAttributeError(req.Path.AtListIndex(i), "Invalid IP Range",
				"from_ip_address must not be greater than to_ip_address.")
		}
	}
}
