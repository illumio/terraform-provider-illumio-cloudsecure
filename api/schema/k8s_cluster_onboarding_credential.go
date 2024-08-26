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
	k8sClusterOnboardingCredential = Resource{
		TypeName: "k8s_cluster_onboarding_credential",
		Schema: resource_schema.Schema{
			Version:     1,
			Description: "Manages a credential (`client_id`/`client_secret` pair) that can be used to onboard one or more k8s clusters onto CloudSecure within a specific Illumio Region.",
			Attributes: map[string]resource_schema.Attribute{
				IDFieldName: idAttribute,
				"client_id": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "Client identifier used by CloudSecure's k8s operator to authenticate to CloudSecure for onboarding, in combination with `client_secret`. Identical to `id`.",
						Computed:    true,
					},
					attributeWithMode: attributeWithMode{
						Mode: ReadOnlyAttributeMode,
					},
				},
				"client_secret": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "Client secret used by CloudSecure's k8s operator to authenticate to CloudSecure for onboarding, in combination with `client_id`.",
						Computed:    true,
						Sensitive:   true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					attributeWithMode: attributeWithMode{
						Mode: ReadOnlyOnceAttributeMode,
					},
				},
				"created_at": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "Timestamp of the creation of the onboarding credential, in RFC 3339 format.",
						Computed:    true,
					},
					attributeWithMode: attributeWithMode{
						Mode: ReadOnlyAttributeMode,
					},
				},
				"description": resource_schema.StringAttribute{
					Description: "Description of the onboarding credential.",
					Optional:    true,
				},
				"illumio_region": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						MarkdownDescription: "Illumio Region where the k8s cluster can be onboarded using this credential. " +
							"An Illumio Region is a designated cloud region where the CloudSecure k8s operators in onboarded k8s clusters connect after onboarding. " +
							"Choose the Illumio Region nearest to each cluster to maximize performance and security. " +
							"Must be one of: `aws-ap-southeast-2`, `aws-eu-west-2`, `aws-us-west-2`.",
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf(
								"aws-ap-southeast-2",
								"aws-eu-west-2",
								"aws-us-west-2",
							),
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
					Description: "Display name for the onboarding credential.",
					Required:    true,
				},
			},
		},
	}
)
