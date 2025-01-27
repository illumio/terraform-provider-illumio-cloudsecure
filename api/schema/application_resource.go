// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var (
	applicationResource = Resource{
		TypeName: "application",
		Schema: resource_schema.Schema{
			Version:             1,
			MarkdownDescription: "An application in the Illumio CloudSecure platform.",
			Attributes: map[string]resource_schema.Attribute{
				IDFieldName: idAttribute,
				"deployment_id": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						MarkdownDescription: "ID of the Cloudsecure deployment.",
						Required:            true,
					},
					attributeWithMode: attributeWithMode{
						Mode: ImmutableAttributeMode,
					},
				},
				"description": resource_schema.StringAttribute{
					MarkdownDescription: "Description of the created application.",
					Optional:            true,
				},
				"name": resource_schema.StringAttribute{
					MarkdownDescription: "Display name of the created CloudSecure application.",
					Required:            true,
				},
			},
		},
	}
)
