// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	awsTagToLabelResource = Resource{
		TypeName: "aws_tag_to_label",
		Schema: resource_schema.Schema{
			Version:     1,
			Description: "Maps AWS account tags to CloudSecure labels.",
			Attributes: map[string]resource_schema.Attribute{
				IDFieldName: idAttribute,
				"key": resource_schema.StringAttribute{
					MarkdownDescription: "CloudSecure label key.",
					Required:            true,
				},
				"name": resource_schema.StringAttribute{
					MarkdownDescription: "CloudSecure label display name.",
					Required:            true,
				},
				"icon": resource_schema.ObjectAttribute{
					AttributeTypes: map[string]attr.Type{
						"name":             types.StringType,
						"background_color": types.StringType,
						"foreground_color": types.StringType,
					},
					Required:    true,
					Description: "Icon details.",
				},
				"cloud_tags": resource_schema.ListAttribute{
					Required:    true,
					Description: "List of AWS account tags to map to the CloudSecure label.",
					ElementType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"key":   types.StringType,
							"cloud": types.StringType,
						},
					},
				},
			},
		},
	}
)
