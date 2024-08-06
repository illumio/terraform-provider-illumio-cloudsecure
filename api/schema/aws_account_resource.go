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
						Description: "AWS account type.",
						Required:    true,
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
						Description: "Access mode.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("ReadWrite"),
						Validators: []validator.String{
							stringvalidator.OneOf("Read", "ReadWrite"),
						},
					},
					attributeWithMode: attributeWithMode{
						Mode: ImmutableAttributeMode,
					},
				},
				"service_account_id": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "AWS service account ID.",
						Required:    true,
					},
					attributeWithMode: attributeWithMode{
						Mode: ImmutableAttributeMode,
					},
				},
				"role_arn": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "Provider AWS role ARN.",
						Required:    true,
					},
					attributeWithMode: attributeWithMode{
						Mode: WriteOnlyOnceAttributeMode,
					},
				},
				"management_account_id": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "Management account ID.",
						Optional:    true,
					},
					attributeWithMode: attributeWithMode{
						Mode: WriteOnlyOnceAttributeMode,
					},
				},
				"organization_id": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "AWS organization ID.",
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
