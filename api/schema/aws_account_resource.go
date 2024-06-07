// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	awsAccountResource = Resource{
		TypeName: "aws_account",
		Schema: resource_schema.Schema{
			Version:     1,
			Description: "Manages an AWS account in CloudSecure.",
			Attributes: map[string]resource_schema.Attribute{
				IDFieldName: idAttribute,
				"name": resource_schema.StringAttribute{
					Description: "Display name.",
					Required:    true,
				},
				"disabled": resource_schema.BoolAttribute{
					Description: "If true, disables this account.",
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(false),
				},
				"account_id": resource_schema.StringAttribute{
					MarkdownDescription: "AWS account ID.",
					Required:            true,
				},
				"account_type": resource_schema.StringAttribute{
					Description: "AWS account type.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.OneOf("Account", "Organization"),
					},
				},
				"mode": resource_schema.StringAttribute{
					Description: "Access mode.",
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString("ReadWrite"),
					Validators: []validator.String{
						stringvalidator.OneOf("Read", "ReadWrite"),
					},
				},
				"service_account_id": resource_schema.StringAttribute{
					Description: "AWS service account ID.",
					Required:    true,
				},
				"excluded_regions": resource_schema.SetAttribute{
					Description: "Set of excluded AWS regions.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"excluded_vpc_ids": resource_schema.SetAttribute{
					Description: "Set of IDs of excluded AWS VPCs.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"excluded_subnet_ids": resource_schema.SetAttribute{
					Description: "Set of IDs of excluded AWS subnets.",
					Optional:    true,
					ElementType: types.StringType,
				},
			},
		},
	}
)
