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
			Version:     1,
			Description: "Resources associated with an application in the Illumio CloudSecure platform.",
			Attributes: map[string]resource_schema.Attribute{
				IDFieldName: idAttribute,
				"application_id": resource_schema.StringAttribute{
					Description: "ID of the application.",
					Required:    true,
				},
				"account_id": resource_schema.StringAttribute{
					Description: "ID of the AWS account.",
					Required:    true,
				},
				"aws_arns": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "ARNs of AWS resources to associate with the Cloudsecure application",
					Optional:    true,
				},
				"aws_directconnect_connections": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS direct connect connections to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_directconnect_virtualinterfaces": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS direct connect virtual interfaces to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_ec2_customergateways": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 customer gateways to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_ec2_egressonlyinternetgateways": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 egress only internet gateways to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_ec2_eips": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 eips to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_ec2_flowlogs": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 flow logs to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_ec2_instances": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 instances to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_ec2_instanceconnectendpoints": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 instance connect endpoints to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_ec2_internetgateways": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 internet gateways to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_ec2_natgateways": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 nat gateways to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_ec2_networkacls": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 network acls to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_ec2_networkinterfaces": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 network interfaces to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_ec2_routetabless": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 route tabless to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_ec2_securitygroups": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 security groups to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_ec2_securitygrouprules": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 security group rules to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_ec2_spotfleetrequests": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 spot fleet requests to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_ec2_spotinstancerequests": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 spot instance requests to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_ec2_subnets": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 subnets to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_ec2_transitgateways": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 transit gateways to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_ec2_transitgatewayattachments": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 transit gateway attachments to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_ec2_transitgatewaymulticastdomains": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 transit gateway multicast domains to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_ec2_transitgatewayroutetables": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 transit gateway route tables to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_ec2_volumes": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 volumes to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_ec2_vpcs": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 VPCs to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_ec2_vpcendpoints": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 VPC endpoints to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_ec2_vpcendpointservices": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 VPC endpoint services to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_ec2_vpcpeerings": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 VPC peerings to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_ec2_vpnconnections": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 VPN connections to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_ec2_vpngateways": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 VPN gateways to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_rds_dbclusters": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS rds db clusters to associate with Cloudsecure Application.",
					Optional:    true,
				},
			},
		},
	}
)
