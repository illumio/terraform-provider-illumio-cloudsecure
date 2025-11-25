// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	organizationPolicyResource = Resource{
		TypeName: "organization_policy",
		Schema: resource_schema.Schema{
			Version:     1,
			Description: "Manages a set of organization-wide policy rules on CloudSecure.",
			Attributes: map[string]resource_schema.Attribute{
				IDFieldName: idAttribute,
				"description": resource_schema.StringAttribute{
					Description: "Description of the CloudSecure organization policy.",
					Optional:    true,
				},
				"enabled": resource_schema.BoolAttribute{
					Description: "Indicates whether the organization policy is enabled.",
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(true),
				},
				"name": resource_schema.StringAttribute{
					Description: "Display name for the CloudSecure organization policy.",
					Required:    true,
				},
				"scopes": resource_schema.SetNestedAttribute{
					Description: "A list of lists of policy items.",
					Optional:    true,
					NestedObject: resource_schema.NestedAttributeObject{
						Attributes: map[string]resource_schema.Attribute{
							"scope": resource_schema.ListAttribute{
								Description: "Inner list of policy items.",
								Required:    true,
								ElementType: types.ObjectType{
									AttrTypes: map[string]attr.Type{
										"type":  types.StringType,
										"label": types.StringType,
									},
								},
							},
						},
					},
				},
			},
		},
	}
)
