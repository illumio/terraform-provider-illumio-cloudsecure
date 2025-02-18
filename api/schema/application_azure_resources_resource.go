// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	applicationAzureResourcesResource = Resource{
		TypeName: "application_azure_resources",
		Schema: resource_schema.Schema{
			Version:     1,
			Description: "Manages a set of Azure resources belonging to a single Azure subscription that are associated with a CloudSecure application.",
			Attributes: map[string]resource_schema.Attribute{
				IDFieldName: idAttribute,
				"application_id": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "ID of the CloudSecure application.",
						Required:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					attributeWithMode: attributeWithMode{
						Mode: CreatableIDAttributeMode,
					},
				},
				"application_resource_ids": ListResourceAttributeWithMode{
					ListAttribute: resource_schema.ListAttribute{
						ElementType: types.StringType,
						Description: "CloudSecure IDs of the resources in the CloudSecure application",
						Computed:    true,
					},
					attributeWithMode: attributeWithMode{
						Mode: CreatableIDAttributeMode,
					},
				},
				"resource_ids": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "IDs of Azure resources to associate with the CloudSecure application",
					Optional:    true,
				},
				"subscription_id": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "ID of the Azure subscription the Azure resources belong to.",
						Required:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					attributeWithMode: attributeWithMode{
						Mode: IDAttributeMode,
					},
				},
			},
		},
	}
)
