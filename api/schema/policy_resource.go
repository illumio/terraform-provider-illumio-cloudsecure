// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var (
	policyResource = Resource{
		TypeName: "policy",
		Schema: resource_schema.Schema{
			Version:     1,
			Description: "Manages a policy on CloudSecure.",
			Attributes: map[string]resource_schema.Attribute{
				IDFieldName: idAttribute,
				"description": resource_schema.StringAttribute{
					Description: "Description of the CloudSecure policy.",
					Optional:    true,
				},
				"name": resource_schema.StringAttribute{
					Description: "Display name for the CloudSecure policy.",
					Required:    true,
				},
			},
		},
	}
)
