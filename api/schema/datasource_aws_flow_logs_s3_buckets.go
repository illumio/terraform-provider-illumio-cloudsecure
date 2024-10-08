// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var (
	awsFlowLogsS3BucketsDataSource = DataSource{
		TypeName: "aws_flow_logs_s3_buckets",
		Schema: datasource_schema.Schema{
			Description: "Lists the AWS S3 buckets that CloudSecure is configured to access for an AWS account.",
			Attributes: map[string]datasource_schema.Attribute{
				IDFieldName: idAttribute,
				"account_id": resource_schema.StringAttribute{
					Description: "AWS account ID.",
					Required:    true,
				},
				"s3_bucket_arns": resource_schema.StringAttribute{
					Description: "ARNs of the AWS S3 buckets containing flow logs that CloudSecure is configured to access for the AWS account.",
					Computed:    true,
				},
			},
		},
	}
)
