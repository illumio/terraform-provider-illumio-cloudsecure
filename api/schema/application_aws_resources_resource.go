// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	applicationAwsResourcesResource = Resource{
		TypeName: "application_aws_resources",
		Schema: resource_schema.Schema{
			Version:             1,
			MarkdownDescription: "Resources associated with an application in the Illumio CloudSecure platform.",
			Attributes: map[string]resource_schema.Attribute{
				IDFieldName: idAttribute,
				"application_id": resource_schema.StringAttribute{
					MarkdownDescription: "ID of the application.",
					Required:            true,
				},
				"account_id": resource_schema.StringAttribute{
					MarkdownDescription: "ID of the AWS account.",
					Required:            true,
				},
				"aws_arns": resource_schema.SetAttribute{
					ElementType:         types.StringType,
					MarkdownDescription: "ARNs of AWS resources to associate with the Cloudsecure application",
					Optional:            true,
				},
				"aws_directconnect_connections": resource_schema.SetAttribute{
					ElementType:         types.StringType,
					MarkdownDescription: "IDs of aws directconnect connections to associate with Cloudsecure Application.",
					Optional:            true,
				},
				"aws_directconnect_virtualinterfaces": resource_schema.SetAttribute{
					ElementType:         types.StringType,
					MarkdownDescription: "IDs of aws directconnect virtualinterfaces to associate with Cloudsecure Application.",
					Optional:            true,
				},
				"aws_ec2_customergateways": resource_schema.SetAttribute{
					ElementType:         types.StringType,
					MarkdownDescription: "IDs of aws ec2 customergateways to associate with Cloudsecure Application.",
					Optional:            true,
				},
				"aws_ec2_egressonlyinternetgateways": resource_schema.SetAttribute{
					ElementType:         types.StringType,
					MarkdownDescription: "IDs of aws ec2 egressonlyinternetgateways to associate with Cloudsecure Application.",
					Optional:            true,
				},
				"aws_ec2_eips": resource_schema.SetAttribute{
					ElementType:         types.StringType,
					MarkdownDescription: "IDs of aws ec2 eips to associate with Cloudsecure Application.",
					Optional:            true,
				},
				"aws_ec2_flowlogs": resource_schema.SetAttribute{
					ElementType:         types.StringType,
					MarkdownDescription: "IDs of aws ec2 flowlogs to associate with Cloudsecure Application.",
					Optional:            true,
				},
				"aws_ec2_instances": resource_schema.SetAttribute{
					ElementType:         types.StringType,
					MarkdownDescription: "IDs of aws ec2 instances to associate with Cloudsecure Application.",
					Optional:            true,
				},
				"aws_ec2_instanceconnectendpoints": resource_schema.SetAttribute{
					ElementType:         types.StringType,
					MarkdownDescription: "IDs of aws ec2 instanceconnectendpoints to associate with Cloudsecure Application.",
					Optional:            true,
				},
				"aws_ec2_internetgateways": resource_schema.SetAttribute{
					ElementType:         types.StringType,
					MarkdownDescription: "IDs of aws ec2 internetgateways to associate with Cloudsecure Application.",
					Optional:            true,
				},
				"aws_ec2_natgateways": resource_schema.SetAttribute{
					ElementType:         types.StringType,
					MarkdownDescription: "IDs of aws ec2 natgateways to associate with Cloudsecure Application.",
					Optional:            true,
				},
				"aws_ec2_networkacls": resource_schema.SetAttribute{
					ElementType:         types.StringType,
					MarkdownDescription: "IDs of aws ec2 networkacls to associate with Cloudsecure Application.",
					Optional:            true,
				},
				"aws_ec2_networkinterfaces": resource_schema.SetAttribute{
					ElementType:         types.StringType,
					MarkdownDescription: "IDs of aws ec2 networkinterfaces to associate with Cloudsecure Application.",
					Optional:            true,
				},
				"aws_ec2_routetabless": resource_schema.SetAttribute{
					ElementType:         types.StringType,
					MarkdownDescription: "IDs of aws ec2 routetabless to associate with Cloudsecure Application.",
					Optional:            true,
				},
				"aws_ec2_securitygroups": resource_schema.SetAttribute{
					ElementType:         types.StringType,
					MarkdownDescription: "IDs of aws ec2 securitygroups to associate with Cloudsecure Application.",
					Optional:            true,
				},
				"aws_ec2_securitygrouprules": resource_schema.SetAttribute{
					ElementType:         types.StringType,
					MarkdownDescription: "IDs of aws ec2 securitygrouprules to associate with Cloudsecure Application.",
					Optional:            true,
				},
				"aws_ec2_spotfleetrequests": resource_schema.SetAttribute{
					ElementType:         types.StringType,
					MarkdownDescription: "IDs of aws ec2 spotfleetrequests to associate with Cloudsecure Application.",
					Optional:            true,
				},
				"aws_ec2_spotinstancerequests": resource_schema.SetAttribute{
					ElementType:         types.StringType,
					MarkdownDescription: "IDs of aws ec2 spotinstancerequests to associate with Cloudsecure Application.",
					Optional:            true,
				},
				"aws_ec2_subnets": resource_schema.SetAttribute{
					ElementType:         types.StringType,
					MarkdownDescription: "IDs of aws ec2 subnets to associate with Cloudsecure Application.",
					Optional:            true,
				},
				"aws_ec2_transitgateways": resource_schema.SetAttribute{
					ElementType:         types.StringType,
					MarkdownDescription: "IDs of aws ec2 transitgateways to associate with Cloudsecure Application.",
					Optional:            true,
				},
				"aws_ec2_transitgatewayattachments": resource_schema.SetAttribute{
					ElementType:         types.StringType,
					MarkdownDescription: "IDs of aws ec2 transitgatewayattachments to associate with Cloudsecure Application.",
					Optional:            true,
				},
				"aws_ec2_transitgatewaymulticastdomains": resource_schema.SetAttribute{
					ElementType:         types.StringType,
					MarkdownDescription: "IDs of aws ec2 transitgatewaymulticastdomains to associate with Cloudsecure Application.",
					Optional:            true,
				},
				"aws_ec2_transitgatewayroutetables": resource_schema.SetAttribute{
					ElementType:         types.StringType,
					MarkdownDescription: "IDs of aws ec2 transitgatewayroutetables to associate with Cloudsecure Application.",
					Optional:            true,
				},
				"aws_ec2_volumes": resource_schema.SetAttribute{
					ElementType:         types.StringType,
					MarkdownDescription: "IDs of aws ec2 volumes to associate with Cloudsecure Application.",
					Optional:            true,
				},
				"aws_ec2_vpcs": resource_schema.SetAttribute{
					ElementType:         types.StringType,
					MarkdownDescription: "IDs of aws ec2 vpcs to associate with Cloudsecure Application.",
					Optional:            true,
				},
				"aws_ec2_vpcendpoints": resource_schema.SetAttribute{
					ElementType:         types.StringType,
					MarkdownDescription: "IDs of aws ec2 vpcendpoints to associate with Cloudsecure Application.",
					Optional:            true,
				},
				"aws_ec2_vpcendpointservices": resource_schema.SetAttribute{
					ElementType:         types.StringType,
					MarkdownDescription: "IDs of aws ec2 vpcendpointservices to associate with Cloudsecure Application.",
					Optional:            true,
				},
				"aws_ec2_vpcpeerings": resource_schema.SetAttribute{
					ElementType:         types.StringType,
					MarkdownDescription: "IDs of aws ec2 vpcpeerings to associate with Cloudsecure Application.",
					Optional:            true,
				},
				"aws_ec2_vpnconnections": resource_schema.SetAttribute{
					ElementType:         types.StringType,
					MarkdownDescription: "IDs of aws ec2 vpnconnections to associate with Cloudsecure Application.",
					Optional:            true,
				},
				"aws_ec2_vpngateways": resource_schema.SetAttribute{
					ElementType:         types.StringType,
					MarkdownDescription: "IDs of aws ec2 vpngateways to associate with Cloudsecure Application.",
					Optional:            true,
				},
				"aws_rds_dbclusters": resource_schema.SetAttribute{
					ElementType:         types.StringType,
					MarkdownDescription: "IDs of aws rds dbclusters to associate with Cloudsecure Application.",
					Optional:            true,
				},
			},
		},
	}
)
