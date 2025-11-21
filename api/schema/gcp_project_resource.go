// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	gcpProjectResource = Resource{
		TypeName: "gcp_project",
		Schema: resource_schema.Schema{
			Version:     1,
			Description: "Manages a GCP project in CloudSecure.",
			Attributes: map[string]resource_schema.Attribute{
				IDFieldName: idAttribute,
				"account_id": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "GCP account ID.",
						Required:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					attributeWithMode: attributeWithMode{
						Mode: ImmutableAttributeMode,
					},
				},
				"mode": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "Access mode, must be `\"ReadWrite\"` or `\"Read\"`.",
						Required:    true,
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
				"organization_id": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "ID of the GCP organization.",
						Required:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					attributeWithMode: attributeWithMode{
						Mode: ImmutableAttributeMode,
					},
				},
				"name": resource_schema.StringAttribute{
					Description: "Display name for the GCP project.",
					Required:    true,
				},
				"service_account_email": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "Service account email that Illumio will impersonate.",
						Required:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					attributeWithMode: attributeWithMode{
						Mode: ImmutableAttributeMode,
					},
				},
				"type": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "Onboarding type, fixed to `\"GcpProject\"`.",
						Computed:    true,
					},
					attributeWithMode: attributeWithMode{
						Mode: ReadOnlyAttributeMode,
					},
				},
				"enable_projects": BoolResourceAttributeWithMode{
					BoolAttribute: resource_schema.BoolAttribute{
						Description: "Whether child project onboarding is enabled, fixed to true.",
						Computed:    true,
					},
					attributeWithMode: attributeWithMode{
						Mode: ReadOnlyAttributeMode,
					},
				},
			},
		},
	}
)
