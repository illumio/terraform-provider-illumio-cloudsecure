// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	organizationPolicyRuleResource = Resource{
		TypeName: "organization_policy_rule",
		Schema: resource_schema.Schema{
			Version:     1,
			Description: "Manages policy rules on CloudSecure organizations.",
			Attributes: map[string]resource_schema.Attribute{
				IDFieldName: idAttribute,
				"action": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						MarkdownDescription: "The action to take for flows matched by the organization policy rule. Must be `\"Allow\"`, `\"Deny\"` or `\"OverrideDeny\"`.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("Allow", "Deny", "OverrideDeny"),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					attributeWithMode: attributeWithMode{
						Mode: ReadWriteAttributeMode,
					},
				},
				"description": resource_schema.StringAttribute{
					Description: "Description of the organization policy rule.",
					Optional:    true,
				},
				"external_scope": BoolResourceAttributeWithMode{
					BoolAttribute: resource_schema.BoolAttribute{
						Description: "Specifies whether the organization policy rule can be applied outside of the scopes of the organization policy. Applicable only for `\"Allow\"` action.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					attributeWithMode: attributeWithMode{
						Mode: ReadWriteAttributeMode,
					},
				},
				"from_ip_list_ids": resource_schema.ListAttribute{
					Description: "List of IDs of IP lists to allow/deny traffic from.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"from_labels": resource_schema.ListAttribute{
					Description: "List of CloudSecure labels of sources to allow/deny traffic from.",
					Optional:    true,
					ElementType: Label,
				},
				"organization_policy_id": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "ID of the CloudSecure organization policy to contain this rule.",
						Required:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					attributeWithMode: attributeWithMode{
						Mode: ReadWriteAttributeMode,
					},
				},
				"to_ip_list_ids": resource_schema.ListAttribute{
					Description: "List of IDs of IP lists to allow/deny traffic to.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"to_labels": resource_schema.ListAttribute{
					Description: "List of CloudSecure labels of destinations to allow/deny traffic to.",
					Optional:    true,
					ElementType: Label,
				},
				"to_port_ranges": resource_schema.ListAttribute{
					MarkdownDescription: "List of transport protocol ports to allow/deny traffic to. The `protocol` for each port must be `\"TCP\"` or `\"UDP\"`.",
					Required:            true,
					ElementType:         PortRange,
				},
			},
		},
	}
)
