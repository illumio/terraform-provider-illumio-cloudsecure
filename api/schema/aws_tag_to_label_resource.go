// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var cloudTagsAttribute = resource_schema.ListAttribute{
	ElementType: types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"key":   types.StringType,
			"cloud": types.StringType,
		},
	},
	Required:    true,
	Description: "List of AWS account tags to map to the CloudSecure label.",
}

var setObjAttribute = resource_schema.SetAttribute{
	ElementType: types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"set_key": types.StringType,
			"set_val": types.StringType,
		},
	},
	Required:    true,
	Description: "List of AWS account tags to map to the CloudSecure label.",
}

var iconAttribute = resource_schema.ObjectAttribute{
	AttributeTypes: map[string]attr.Type{
		"name":             types.StringType,
		"background_color": types.StringType,
		"foreground_color": types.StringType,
	},
	Required:    true,
	Description: "Icon details.",
}

var objInObjAttribute = resource_schema.ObjectAttribute{
	AttributeTypes: map[string]attr.Type{
		"name": types.StringType,
		"child": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"name": types.StringType,
				"grand_child": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name": types.StringType,
					},
				},
			},
		},
	},
	Required:    true,
	Description: "Icon details.",
}

var (
	awsTagToLabelResource = Resource{
		TypeName: "nested_object_tester",
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
				"icon":       iconAttribute,
				"cloud_tags": cloudTagsAttribute,
				"set_obj":    setObjAttribute,
				"set_of_sets_string": resource_schema.SetAttribute{
					ElementType: types.SetType{
						ElemType: types.SetType{
							ElemType: types.StringType,
						},
					},
					Optional: true,
				},
				"set_of_strings": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Optional:    true,
				},
				"obj_in_obj": objInObjAttribute,
			},
		},
	}
)
