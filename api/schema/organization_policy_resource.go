// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
)

var (
	organizationPolicyResource = Resource{
		TypeName: "organization_policy",
		Schema: resource_schema.Schema{
			Version:     1,
			Description: "Manages policy rules on CloudSecure organizations.",
			Attributes: map[string]resource_schema.Attribute{
				IDFieldName: idAttribute,
				"description": resource_schema.StringAttribute{
					Description: "Description of the CloudSecure application.",
					Optional:    true,
				},
				"name": resource_schema.StringAttribute{
					Description: "Display name for the CloudSecure application.",
					Required:    true,
				},
				"enabled": resource_schema.BoolAttribute{
					Description: "Indicates whether the organization policy is enabled.",
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(true),
				},
			},
		},
	}
)
