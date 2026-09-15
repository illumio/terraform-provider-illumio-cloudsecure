// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var (
	awsFlowLogsS3BucketSourceResource = Resource{
		TypeName: "aws_flow_logs_s3_bucket_source",
		Schema: resource_schema.Schema{
			Version:     1,
			Description: "Registers a customer-owned AWS S3 bucket (by name and optional path prefix) as a source of VPC flow logs for CloudSecure to ingest.",
			Attributes: map[string]resource_schema.Attribute{
				IDFieldName: idAttribute,
				"account_id": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						MarkdownDescription: "AWS account ID that owns the S3 bucket.",
						Required:            true,
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
						Description: "Name of the customer-specified AWS S3 bucket where VPC flow logs are delivered.",
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
						Description: "Optional path/prefix within the S3 bucket under which flow log objects are written.",
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
