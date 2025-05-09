// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var Label = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"key":   types.StringType,
		"value": types.StringType,
	},
}

var PortRange = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"from_port": types.Int64Type,
		"to_port":   types.Int64Type,
		"protocol":  types.StringType,
	},
}

var (
	applicationPolicyRuleResource = Resource{
		TypeName: "application_policy_rule",
		Schema: resource_schema.Schema{
			Version:     1,
			Description: "Manages policy rules on CloudSecure applications.",
			Attributes: map[string]resource_schema.Attribute{
				IDFieldName: idAttribute,
				"action": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						MarkdownDescription: "The action to take for flows matched by the application policy rule. Must be `\"Allow\"`, `\"Deny\"` or `\"OverrideDeny\"`.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("Allow", "Deny", "OverrideDeny"),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					attributeWithMode: attributeWithMode{
						// TODO: Change this to ReadWriteAttributeMode.
						Mode: KeyAttributeMode,
					},
				},
				"application_id": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "ID of the CloudSecure application on which this rule is based.",
						Required:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					attributeWithMode: attributeWithMode{
						// TODO: Change this to ReadWriteAttributeMode.
						Mode: KeyAttributeMode,
					},
				},
				"description": resource_schema.StringAttribute{
					Description: "Description of the application policy rule.",
					Optional:    true,
				},
				"external_scope": BoolResourceAttributeWithMode{
					BoolAttribute: resource_schema.BoolAttribute{
						MarkdownDescription: "Specifies whether the application policy allow rule can be applied to scope outside the CloudSecure application. Applicable only for `\"Allow\"` action.",
						Optional:            true,
						PlanModifiers: []planmodifier.Bool{
							// TODO: Remove this when moving to ReadWriteAttributeMode.
							boolplanmodifier.RequiresReplace(),
						},
					},
					attributeWithMode: attributeWithMode{
						// TODO: Change this to ReadWriteAttributeMode.
						Mode: WriteOnlyOnceAttributeMode,
					},
				},
				"from_ip_list_ids": resource_schema.ListAttribute{
					Description: "List of IDs of IP lists to allow/deny traffic from.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"from_labels": resource_schema.ListAttribute{
					Description: "List of Cloudsecure labels of sources to allow/deny traffic from.",
					Optional:    true,
					ElementType: Label,
				},
				"to_ip_list_ids": resource_schema.ListAttribute{
					Description: "List of IDs of IP lists to allow/deny traffic to.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"to_labels": resource_schema.ListAttribute{
					Description: "List of Cloudsecure labels of destinations to allow/deny traffic to.",
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
