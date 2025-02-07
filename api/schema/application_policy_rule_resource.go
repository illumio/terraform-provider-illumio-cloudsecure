// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
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

var Port = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"port_number": types.Int32Type,
		"protocol":    types.StringType,
	},
}

var (
	applicationPolicyRuleResource = Resource{
		TypeName: "application_policy_rule",
		Schema: resource_schema.Schema{
			Version:     1,
			Description: "Manages policy rules on Cloudsecure applications.",
			Attributes: map[string]resource_schema.Attribute{
				IDFieldName: idAttribute,
				"action": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "Must be `\"Allow\"` or `\"Deny\"`.",
						Validators: []validator.String{
							stringvalidator.OneOf("Allow", "Deny"),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
				"description": resource_schema.StringAttribute{
					Description: "Description of the Application policy rule.",
					Optional:    true,
				},
				"from_ip_list_ids": resource_schema.ListAttribute{
					Description: "List of IDs of IP Lists for source.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"from_labels": resource_schema.ListAttribute{
					Description: "List of Cloudsecure labels for source to be associated with this rule.",
					Optional:    true,
					ElementType: Label,
				},
				"to_ip_list_ids": resource_schema.ListAttribute{
					Description: "List of IDs of IP Lists for destination.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"to_labels": resource_schema.ListAttribute{
					Description: "List of Cloudsecure labels for destination to be associated with this rule.",
					Optional:    true,
					ElementType: Label,
				},
				"to_ports": resource_schema.ListAttribute{
					Description: "List of Ports to be associated with this rule.",
					Optional:    true,
					ElementType: Port,
				},
			},
		},
	}
)
