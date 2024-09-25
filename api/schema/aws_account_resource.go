// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	awsAccountResource = Resource{
		TypeName: "aws_account",
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
				"access_mode": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "Access mode, must be `\"ReadWrite\"` (default) or `\"Read\"`.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("ReadWrite"),
						Validators: []validator.String{
							stringvalidator.OneOf("ReadWrite", "Read"),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					attributeWithMode: attributeWithMode{
						Mode: ImmutableAttributeMode,
					},
				},
				"name": resource_schema.StringAttribute{
					Description: "Display name for the AWS account.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
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
						Mode: WriteOnlyOnceAttributeMode,
					},
				},
				"organization_id": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "ID of the AWS organization.",
						Required:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					attributeWithMode: attributeWithMode{
						Mode: WriteOnlyOnceAttributeMode,
					},
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
						Sensitive:   true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					attributeWithMode: attributeWithMode{
						Mode: WriteOnlyOnceAttributeMode,
					},
				},
			},
		},
	}
)
