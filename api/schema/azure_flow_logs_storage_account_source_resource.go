// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// NOTE: The exact Azure field set is provisional and pending the per-service
// requirements sync. Confirm before backend work begins.
var (
	azureFlowLogsStorageAccountSourceResource = Resource{
		TypeName: "azure_flow_logs_storage_account_source",
		Schema: resource_schema.Schema{
			Version:     1,
			Description: "Registers a customer-owned Azure Storage Account (by name, container, and optional path prefix) as a source of flow logs for CloudSecure to ingest.",
			Attributes: map[string]resource_schema.Attribute{
				IDFieldName: idAttribute,
				"subscription_id": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						MarkdownDescription: "Azure subscription ID that owns the Storage Account.",
						Required:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					attributeWithMode: attributeWithMode{
						Mode: ImmutableAttributeMode,
					},
				},
				"storage_account_name": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "Name of the customer-specified Azure Storage Account where flow logs are delivered.",
						Required:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					attributeWithMode: attributeWithMode{
						Mode: ImmutableAttributeMode,
					},
				},
				"container_name": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "Optional blob container within the Storage Account where flow log blobs are written.",
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					attributeWithMode: attributeWithMode{
						Mode: ImmutableAttributeMode,
					},
				},
				"path_prefix": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "Optional path/prefix within the container under which flow log blobs are written.",
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					attributeWithMode: attributeWithMode{
						Mode: ImmutableAttributeMode,
					},
				},
			},
		},
	}
)
