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
			MarkdownDescription: "Maps AWS account tags to CloudSecure labels.",
			Attributes: map[string]resource_schema.Attribute{
				IDFieldName: idAttribute,
				"cloud_tags": resource_schema.ListAttribute{
					MarkdownDescription: "List of tags to map to CloudSecure labels with the specified key. The values of the created labels correspond to the values of the tags. The cloud field for each tag must be \"aws\" or \"azure\".",
					Required:            true,
					ElementType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"key":   types.StringType,
							"cloud": types.StringType,
						},
					},
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
