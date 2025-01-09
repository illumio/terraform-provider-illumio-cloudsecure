package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var TagType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"key":   types.StringType,
		"value": types.StringType,
	},
}

var (
	deploymentResource = Resource{
		TypeName: "deployment",
		Schema: resource_schema.Schema{
			Version:     1,
			Description: "Groups resources under Cloudsecure deployments.",
			Attributes: map[string]resource_schema.Attribute{
				IDFieldName: idAttribute,
				"aws_account_ids": resource_schema.ListAttribute{
					Optional:    true,
					Description: "IDs of AWS accounts included in this deployment.",
					ElementType: types.StringType,
				},
				"aws_regions": resource_schema.ListAttribute{
					Optional:    true,
					Description: "AWS regions included in this deployment.",
					ElementType: types.StringType,
				},
				"aws_subnet_ids": resource_schema.ListAttribute{
					Optional:    true,
					Description: "IDs of AWS subnets included in this deployment.",
					ElementType: types.StringType,
				},
				"aws_tags": resource_schema.ListAttribute{
					Optional:    true,
					Description: "AWS tags included in this deployment.",
					ElementType: TagType,
				},
				"aws_vpc_ids": resource_schema.ListAttribute{
					Optional:    true,
					Description: "IDs of AWS VPCs included in this deployment.",
					ElementType: types.StringType,
				},
				"azure_regions": resource_schema.ListAttribute{
					Optional:    true,
					Description: "Azure regions included in this deployment.",
					ElementType: types.StringType,
				},
				"azure_subnet_ids": resource_schema.ListAttribute{
					Optional:    true,
					Description: "IDs of Azure subnets included in this deployment.",
					ElementType: types.StringType,
				},
				"azure_subscription_ids": resource_schema.ListAttribute{
					Optional:    true,
					Description: "IDs of Azure subscriptions included in this deployment.",
					ElementType: types.StringType,
				},
				"azure_tags": resource_schema.ListAttribute{
					Optional:    true,
					Description: "Azure tags included in this deployment.",
					ElementType: TagType,
				},
				"azure_vnet_ids": resource_schema.ListAttribute{
					Optional:    true,
					Description: "IDs of Azure VNets included in this deployment.",
					ElementType: types.StringType,
				},
				"description": resource_schema.StringAttribute{
					Description: "Description of Cloudsecure deployment.",
					Optional:    true,
				},
				"name": resource_schema.StringAttribute{
					Description: "Name of the CloudSecure deployment.",
					Required:    true,
				},
			},
		},
	}
)
