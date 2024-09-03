// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var (
	awsOrganizationAccountResource = Resource{
		TypeName: "aws_organization_account",
		Schema: resource_schema.Schema{
			Version:     1,
			Description: "Manages an AWS account in CloudSecure.",
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
				"disabled": resource_schema.BoolAttribute{
					Description: "If true, disables this account.",
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(false),
				},
				"role_arn": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "ARN of the AWS role to be assumed by CloudSecure to manage this account.",
						Required:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					attributeWithMode: attributeWithMode{
						Mode: WriteOnlyOnceAttributeMode,
					},
				},
				"role_external_id": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "External ID defined in the AWS role to authenticate CloudSecure when assuming that role.",
						Required:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
						Sensitive: true,
					},
					attributeWithMode: attributeWithMode{
						Mode: WriteOnlyOnceAttributeMode,
					},
				},
				"organization_master_account_id": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "ID of the master account of the AWS organization this account belongs to. If specified, should be the `master_account_id` of an `aws_organization`.",
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
