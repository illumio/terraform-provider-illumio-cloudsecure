// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var (
	gcpFlowLogsPubsubTopicResource = Resource{
		TypeName: "gcp_flow_logs_pubsub_topic",
		Schema: resource_schema.Schema{
			Version:     1,
			Description: "Manages CloudSecure access to flow logs in a GCP Pub/Sub topic.",
			Attributes: map[string]resource_schema.Attribute{
				IDFieldName: idAttribute,
				"project_id": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "ID of the GCP project.",
						Required:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					attributeWithMode: attributeWithMode{
						Mode: ImmutableAttributeMode,
					},
				},
				"pubsub_topic_id": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "Resource ID of the GCP Pub/Sub topic containing flow logs.",
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
