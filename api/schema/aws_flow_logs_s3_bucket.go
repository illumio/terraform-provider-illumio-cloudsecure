// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var (
	awsFlowLogsS3Bucket = Resource{
		TypeName: "aws_flow_logs_s3_bucket",
		Schema: resource_schema.Schema{
			Version:     1,
			Description: "Manages CloudSecure access to flow logs in an AWS S3 bucket.",
			Attributes: map[string]resource_schema.Attribute{
				IDFieldName: idAttribute,
				"account_id": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						MarkdownDescription: "AWS account ID.",
						Required:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					attributeWithMode: attributeWithMode{
						Mode: ImmutableAttributeMode,
					},
				},
				"s3_bucket_arn": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "ARN of the AWS S3 bucket containing flow logs.",
						Required:    true,
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
