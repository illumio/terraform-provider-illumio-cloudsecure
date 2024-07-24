// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var (
	idAttribute = StringResourceAttributeWithMode{
		StringAttribute: schema.StringAttribute{
			Description: "CloudSecure ID.",
			Computed:    true,
		},
		attributeWithMode: attributeWithMode{
			Mode: IDAttributeMode,
		},
	}
)
