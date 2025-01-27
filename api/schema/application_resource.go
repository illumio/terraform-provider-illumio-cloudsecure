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
				"name": resource_schema.StringAttribute{
					MarkdownDescription: "Display name of the created CloudSecure application.",
					Required:            true,
				},
				"deployment_id": resource_schema.StringAttribute{
					MarkdownDescription: "ID of the Cloudsecure deployment.",
					Required:            true,
				},
				"description": resource_schema.StringAttribute{
					MarkdownDescription: "Description of the created application.",
					Optional:            true,
				},
			},
		},
	}
)
