// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var (
	k8sClusterResource = Resource{
		TypeName: "k8s_cluster",
		Schema: resource_schema.Schema{
			Version:             1,
			MarkdownDescription: "Manages the onboarding of a k8s cluster on CloudSecure in a specific Illumio Region.",
			Attributes: map[string]resource_schema.Attribute{
				IDFieldName: idAttribute,
				"client_id": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						MarkdownDescription: "Client identifier specific to the k8s cluster, used by CloudSecure's k8s operator to authenticate to CloudSecure, in combination with `client_secret`. Identical to `id`.",
						Computed:            true,
					},
					attributeWithMode: attributeWithMode{
						Mode: ReadOnlyAttributeMode,
					},
				},
				"client_secret": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						MarkdownDescription: "Client secret specific to the k8s cluster, used by CloudSecure's k8s operator to authenticate to CloudSecure, in combination with `client_id`.",
						Computed:            true,
						Sensitive:           true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					attributeWithMode: attributeWithMode{
						Mode: ReadOnlyOnceAttributeMode,
					},
				},
				"illumio_region": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						MarkdownDescription: "Illumio Region where the k8s cluster will be onboarded. " +
							"An Illumio Region is a designated cloud region where the CloudSecure cloud-operator deployed in the k8s cluster connects after onboarding. " +
							"Choose the Illumio Region nearest to the k8s cluster to maximize performance and security. " +
							"Must be one of: `aws-ap-southeast-2`, `aws-eu-west-2`, `aws-me-central-1`, `aws-us-west-1`, `aws-us-west-2`, `azure-australia-east`, `azure-germany-west-central`, `azure-japan-east`, `azure-north-europe`, `azure-uae-north`, `azure-us-east`, `azure-us-west-2`, `azure-us-west-3`.",
						Required: true,
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
