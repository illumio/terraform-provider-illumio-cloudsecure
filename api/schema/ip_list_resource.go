// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var IPRange = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"exclusion": types.BoolType,
		"from_ip":   types.StringType,
		"to_ip":     types.StringType,
	},
}

var IPAddress = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"exclusion": types.BoolType,
		"ip":        types.StringType,
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
				},
				"ip_addresses": resource_schema.ListAttribute{
					Optional:    true,
					Description: "List of IP addresses.",
					ElementType: IPAddress,
				},
				"ip_ranges": resource_schema.ListAttribute{
					Optional:    true,
					Description: "List of IP address ranges.",
					ElementType: IPRange,
				},
				"name": resource_schema.StringAttribute{
					Description: "Name of the IP list.",
					Required:    true,
				},
			},
		},
	}
)
