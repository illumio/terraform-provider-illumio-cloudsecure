// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
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
				"name": resource_schema.StringAttribute{
					Description: "Display name.",
					Required:    true,
				},
				"account_id": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						MarkdownDescription: "AWS account ID.",
						Required:            true,
					},
					attributeWithMode: attributeWithMode{
						Mode: ImmutableAttributeMode,
					},
				},
				"account_type": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						MarkdownDescription: "AWS account type, must be `\"Account\"` or `\"Organization\"`.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("Account", "Organization"),
						},
					},
					attributeWithMode: attributeWithMode{
						Mode: ImmutableAttributeMode,
					},
				},
				"mode": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "Access mode, must be `\"ReadWrite\"` (default) or `\"Read\"`.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("ReadWrite"),
						Validators: []validator.String{
							stringvalidator.OneOf("ReadWrite", "Read"),
						},
					},
					attributeWithMode: attributeWithMode{
						Mode: ImmutableAttributeMode,
					},
				},
				"role_external_id": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "External ID defined in the AWS role to authenticate CloudSecure when assuming that role.",
						Required:    true,
					},
					attributeWithMode: attributeWithMode{
						Mode: WriteOnlyOnceAttributeMode,
					},
				},
				"role_arn": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "ARN of the AWS role to be assumed by CloudSecure to manage this account.",
						Required:    true,
					},
					attributeWithMode: attributeWithMode{
						Mode: WriteOnlyOnceAttributeMode,
					},
				},
				"management_account_id": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						MarkdownDescription: "AWS organization management account ID. If specified, `organization_id` must also be specified. Required if `account_type` is `\"Organization\"`.",
						Optional:            true,
					},
					attributeWithMode: attributeWithMode{
						Mode: WriteOnlyOnceAttributeMode,
					},
				},
				"organization_id": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "AWS organization ID. If specified, the whole AWS organization is onboarded instead of just the AWS account. If specified, `management_account_id` must also be specified. Required if `account_type` is `\"Organization\"`.",
						Optional:    true,
					},
					attributeWithMode: attributeWithMode{
						Mode: WriteOnlyOnceAttributeMode,
					},
				},
			},
		},
	}
)
