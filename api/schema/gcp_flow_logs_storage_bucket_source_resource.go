// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// NOTE: The exact GCP field set is provisional and pending the per-service
// requirements sync. Confirm before backend work begins.
var (
	gcpFlowLogsStorageBucketSourceResource = Resource{
		TypeName: "gcp_flow_logs_storage_bucket_source",
		Schema: resource_schema.Schema{
			Version:     1,
			Description: "Registers a customer-owned GCP Cloud Storage bucket (by name and optional path prefix) as a source of flow logs for CloudSecure to ingest.",
			Attributes: map[string]resource_schema.Attribute{
				IDFieldName: idAttribute,
				"project_id": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "ID of the GCP project that owns the Cloud Storage bucket.",
						Required:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					attributeWithMode: attributeWithMode{
						Mode: ImmutableAttributeMode,
					},
				},
				"bucket_name": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "Name of the customer-specified GCP Cloud Storage bucket where flow logs are delivered.",
						Required:    true,
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
						Description: "Optional path/prefix within the bucket under which flow log objects are written.",
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
