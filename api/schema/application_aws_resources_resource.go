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
			Description: "Manages application resources in CloudSecure.",
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
						Mode: ImmutableAttributeMode,
					},
				},
				"account_id": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "ID of the AWS account.",
						Required:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					attributeWithMode: attributeWithMode{
						Mode: ImmutableAttributeMode,
					},
				},
				"aws_arns": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "ARNs of AWS resources to associate with the CloudSecure application",
					Optional:    true,
				},
				"aws_customer_gateway_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS customer gateways to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_dx_connection_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS direct connect connections to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_dx_virtual_interface_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS direct connect virtual interfaces (public/private/hosted/transit) to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_ebs_volume_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ebs volumes to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_ec2_instance_connect_endpoint_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 instance connect endpoints to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_ec2_transit_gateway_attachments": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 transit gateway attachments (peering/vpc) to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_ec2_transit_gateway_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 transit gateways to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_ec2_transit_gateway_multicast_domain_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 transit gateway multicast domains to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_ec2_transit_gateway_route_table_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 transit gateway route tables to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_egress_only_internet_gateway_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS egress only internet gateways to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_eip_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS elastic IPs to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_flow_log_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS flow logs to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_instances_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS ec2 instances to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_internet_gateway_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS internet gateways to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_nat_gateway_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS nat gateways to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_network_acl": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS network acls to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_network_interface_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS network interfaces to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_rds_cluster_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS rds db clusters to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_route_table_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS route tabless to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_security_group_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS security groups to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_security_group_rule_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS security group rules to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_spot_fleet_request_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS spot fleet requests to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_spot_instance_request_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS spot instance requests to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_subnet_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS subnets to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_vpc_endpoint_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS VPC endpoints to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_vpc_endpoint_service_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS VPC endpoint services to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_vpc_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS VPCs to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_vpc_peering_connection_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS  VPC peering connections to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_vpn_connection_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS VPN connections to associate with the CloudSecure Application.",
					Optional:    true,
				},
				"aws_vpn_gateway_ids": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Description: "IDs of AWS VPN gateways to associate with the CloudSecure Application.",
					Optional:    true,
				},
			},
		},
	}
)
