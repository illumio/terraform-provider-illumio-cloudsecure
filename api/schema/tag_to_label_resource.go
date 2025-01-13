// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	tagToLabelResource = Resource{
		TypeName: "tag_to_label",
		Schema: resource_schema.Schema{
			Version:             1,
			MarkdownDescription: "Maps cloud resource tags to CloudSecure labels.",
			Attributes: map[string]resource_schema.Attribute{
				IDFieldName: idAttribute,
				"aws_tag_keys": resource_schema.SetAttribute{
					Description: "Sets of keys of AWS resource tags to map to CloudSecure labels with the same keys. The values of the created labels correspond to the values of the tags.",
					Required:    true,
					ElementType: types.StringType,
				},
				"azure_tag_keys": resource_schema.SetAttribute{
					Description: "Set of keys of Azure resource tags to map to CloudSecure labels with the same keys. The values of the created labels correspond to the values of the tags.",
					Required:    true,
					ElementType: types.StringType,
				},
				"icon": resource_schema.ObjectAttribute{
					MarkdownDescription: "Icon of the created CloudSecure labels.",
					Required:            true,
					AttributeTypes: map[string]attr.Type{
						"name":             types.StringType,
						"background_color": types.StringType,
						"foreground_color": types.StringType,
					},
				},
				"key": resource_schema.StringAttribute{
					MarkdownDescription: "Key of the created CloudSecure labels.",
					Required:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				"name": resource_schema.StringAttribute{
					MarkdownDescription: "Display name of the created CloudSecure labels.",
					Required:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
		},
	}
)
