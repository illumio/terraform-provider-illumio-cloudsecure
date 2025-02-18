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
	applicationAwsResourcesResource = Resource{
		TypeName: "application_aws_resources",
		Schema: resource_schema.Schema{
			Version:     1,
			Description: "Manages a set of AWS resources belonging to a single AWS account that are associated with a CloudSecure application.",
			Attributes: map[string]resource_schema.Attribute{
				IDFieldName: idAttribute,
				"account_id": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "ID of the AWS account the AWS resources belong to.",
						Required:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					attributeWithMode: attributeWithMode{
						Mode: CreatableIDAttributeMode,
					},
				},
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
						Mode: IDAttributeMode,
					},
				},
				"arns": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "ARNs of AWS resources to associate with the CloudSecure application",
					Optional:    true,
				},
				"aws_customer_gateway_ids": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS customer gateways to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_dx_connection_ids": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS Direct Connect connections to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_dx_virtual_interface_ids": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS Direct Connect virtual interfaces (public/private/hosted/transit) to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_ebs_volume_ids": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS EBS volumes to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_ec2_instance_connect_endpoint_ids": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS EC2 Instance connect endpoints to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_ec2_transit_gateway_attachment_ids": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS EC2 transit gateway attachments (peering/vpc) to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_ec2_transit_gateway_ids": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS EC2 transit gateways to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_ec2_transit_gateway_multicast_domain_ids": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS EC2 transit gateway multicast domains to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_ec2_transit_gateway_route_table_ids": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS EC2 transit gateway route tables to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_egress_only_internet_gateway_ids": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS egress-only Internet gateways to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_eip_ids": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS Elastic IPs to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_flow_log_ids": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS Flow Logs to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_instances_ids": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS EC2 instances to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_internet_gateway_ids": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS Internet Gateways to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_nat_gateway_ids": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS NAT Gateways to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_network_acl_ids": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS network ACLs to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_network_interface_ids": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS Elastic Network Interfaces (ENI) to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_rds_cluster_ids": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS RDS database clusters to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_route_table_ids": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS VPC routing tables to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_security_group_ids": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS security groups to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_security_group_rule_ids": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS security group rules to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_spot_fleet_request_ids": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS spot fleet requests to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_spot_instance_request_ids": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS spot instance requests to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_subnet_ids": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS subnets to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_vpc_endpoint_ids": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS VPC endpoints to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_vpc_endpoint_service_ids": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS VPC endpoint services to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_vpc_ids": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS VPCs to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_vpc_peering_connection_ids": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS VPC peering connections to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_vpn_connection_ids": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS VPN connections to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_vpn_gateway_ids": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS VPN gateways to associate with the CloudSecure Application.",
					Optional:    true,
				},
			},
		},
	}
)
