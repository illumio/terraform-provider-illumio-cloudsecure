// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var (
	applicationResource = Resource{
		TypeName: "application",
		Schema: resource_schema.Schema{
			Version:     1,
			Description: "Manages an application in CloudSecure.",
			Attributes: map[string]resource_schema.Attribute{
				IDFieldName: idAttribute,
				"deployment_id": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "ID of the CloudSecure deployment.",
						Required:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					attributeWithMode: attributeWithMode{
						Mode: ImmutableAttributeMode,
					},
				},
				"description": resource_schema.StringAttribute{
					Description: "Description of the CloudSecure application.",
					Optional:    true,
				},
				"name": resource_schema.StringAttribute{
					Description: "Display name for the CloudSecure application.",
					Required:    true,
				},
			},
		},
	}
)
