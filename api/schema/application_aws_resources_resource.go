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
					Description: "ID of the CloudSecure application.",
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
				"aws_dx_connection_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS direct connect connections to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_dx_virtual_interface_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS direct connect virtual interfaces (public/private/hosted/transit) to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_customer_gateway_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS customer gateways to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_egress_only_internet_gateway_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS egress only internet gateways to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_eip_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 eips to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_flow_log_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 flow logs to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_instances_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 instances to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_ec2_instance_connect_endpoint_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 instance connect endpoints to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_internet_gateway_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 internet gateways to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_nat_gateway_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 nat gateways to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_network_acl": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 network acls to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_network_interface_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 network interfaces to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_route_table_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 route tabless to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_security_group_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 security groups to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_security_group_rule_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 security group rules to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_spot_fleet_request_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 spot fleet requests to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_spot_instance_request_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 spot instance requests to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_subnet_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 subnets to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_ec2_transit_gateway_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 transit gateways to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_ec2_transit_gateway_attachments": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS transit gateway attachments (peering/vpc) to associate with Cloudsecure Application.",
					Optional:    true,
				},
				"aws_ec2_transit_gateway_multicast_domain_ids": resource_schema.SetAttribute{
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
