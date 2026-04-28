// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var (
	policyProvisionResource = Resource{
		TypeName: "policy_provision",
		Schema: resource_schema.Schema{
			Version:     1,
			Description: "Provisions a policy version on CloudSecure, making it active. Only one version can be provisioned per policy.",
			Attributes: map[string]resource_schema.Attribute{
				IDFieldName: idAttribute,
				"policy_version_id": resource_schema.StringAttribute{
					Description: "ID of the policy version to provision.",
					Required:    true,
				},
			},
		},
	}
)
