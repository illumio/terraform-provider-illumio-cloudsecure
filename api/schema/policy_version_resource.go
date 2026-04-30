// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"maps"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// labelSelectorAttributes is reused for namespace_selector and workload_selector.
var labelSelectorAttributes = map[string]resource_schema.Attribute{
	"match_labels": resource_schema.ListNestedAttribute{
		Description: "List of label key-value pairs to match. All labels must match (AND logic).",
		Optional:    true,
		NestedObject: resource_schema.NestedAttributeObject{
			Attributes: map[string]resource_schema.Attribute{
				"key": resource_schema.StringAttribute{
					Description: "Label key.",
					Required:    true,
				},
				"value": resource_schema.StringAttribute{
					Description: "Label value.",
					Required:    true,
				},
			},
		},
	},
	"match_expressions": resource_schema.ListNestedAttribute{
		Description: "List of label selector requirements. All requirements must match (AND logic).",
		Optional:    true,
		NestedObject: resource_schema.NestedAttributeObject{
			Attributes: map[string]resource_schema.Attribute{
				"key": resource_schema.StringAttribute{
					Description: "Label key to match.",
					Required:    true,
				},
				"operator": resource_schema.StringAttribute{
					Description: "Operator for the label selector. Must be one of: In, NotIn, Exists, DoesNotExist.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.OneOf("In", "NotIn", "Exists", "DoesNotExist"),
					},
				},
				"values": resource_schema.ListAttribute{
					Description: "List of values for the label selector. Required for In and NotIn operators.",
					Optional:    true,
					ElementType: types.StringType,
				},
			},
		},
	},
}

// k8sSelectorAttributes is reused for source.k8s and destination.k8s.
var k8sSelectorAttributes = map[string]resource_schema.Attribute{
	"clusters": resource_schema.ListNestedAttribute{
		Description: "List of K8s clusters. Each entry identifies one cluster. Any cluster matches (OR logic).",
		Required:    true,
		NestedObject: resource_schema.NestedAttributeObject{
			Attributes: map[string]resource_schema.Attribute{
				"aws": resource_schema.SingleNestedAttribute{
					Description: "AWS EKS cluster. Mutually exclusive with gcp, azure, and oci.",
					Optional:    true,
					Attributes: map[string]resource_schema.Attribute{
						"account_id": resource_schema.StringAttribute{
							Description: "AWS account ID.",
							Required:    true,
						},
						"region": resource_schema.StringAttribute{
							Description: "AWS region (e.g., us-east-1).",
							Required:    true,
						},
						"cluster_name": resource_schema.StringAttribute{
							Description: "EKS cluster name.",
							Required:    true,
						},
					},
				},
				"gcp": resource_schema.SingleNestedAttribute{
					Description: "GCP GKE cluster. Mutually exclusive with aws, azure, and oci.",
					Optional:    true,
					Attributes: map[string]resource_schema.Attribute{
						"project_id": resource_schema.StringAttribute{
							Description: "GCP project ID.",
							Required:    true,
						},
						"location": resource_schema.StringAttribute{
							Description: "GKE cluster location (region or zone).",
							Required:    true,
						},
						"cluster_name": resource_schema.StringAttribute{
							Description: "GKE cluster name.",
							Required:    true,
						},
					},
				},
				"azure": resource_schema.SingleNestedAttribute{
					Description: "Azure AKS cluster. Mutually exclusive with aws, gcp, and oci.",
					Optional:    true,
					Attributes: map[string]resource_schema.Attribute{
						"subscription_id": resource_schema.StringAttribute{
							Description: "Azure subscription ID.",
							Required:    true,
						},
						"resource_group": resource_schema.StringAttribute{
							Description: "Azure resource group name.",
							Required:    true,
						},
						"cluster_name": resource_schema.StringAttribute{
							Description: "AKS cluster name.",
							Required:    true,
						},
					},
				},
				"oci": resource_schema.SingleNestedAttribute{
					Description: "OCI OKE cluster. Mutually exclusive with aws, gcp, and azure.",
					Optional:    true,
					Attributes: map[string]resource_schema.Attribute{
						"compartment_id": resource_schema.StringAttribute{
							Description: "OCI compartment OCID.",
							Required:    true,
						},
						"region": resource_schema.StringAttribute{
							Description: "OCI region (e.g., us-ashburn-1).",
							Required:    true,
						},
						"cluster_name": resource_schema.StringAttribute{
							Description: "OKE cluster name.",
							Required:    true,
						},
					},
				},
			},
		},
	},
	"namespace_selector": resource_schema.SingleNestedAttribute{
		Description: "Label selector for K8s namespaces. Use empty {} to match all namespaces.",
		Required:    true,
		Attributes:  labelSelectorAttributes,
	},
	"workload_selector": resource_schema.SingleNestedAttribute{
		Description: "Label selector for K8s workloads (pods). Use empty {} to match all pods.",
		Required:    true,
		Attributes:  workloadSelectorAttributes,
	},
}

// ipSelectorAttributes is reused for source.ip_list and destination.ip_list.
var ipSelectorAttributes = map[string]resource_schema.Attribute{
	"ids": resource_schema.ListAttribute{
		Description: "List of ip_list resource IDs. Traffic matching any list is selected (OR logic).",
		Required:    true,
		ElementType: types.StringType,
	},
}

// workloadSelectorAttributes extends labelSelectorAttributes with service_accounts.
var workloadSelectorAttributes = func() map[string]resource_schema.Attribute {
	m := maps.Clone(labelSelectorAttributes)
	m["service_accounts"] = resource_schema.ListAttribute{
		Description: "List of K8s service account names. Workload must run with one of these service accounts (OR logic).",
		Optional:    true,
		ElementType: types.StringType,
	}

	return m
}()

var (
	policyVersionResource = Resource{
		TypeName: "policy_version",
		Schema: resource_schema.Schema{
			Version:             1,
			MarkdownDescription: "Manages an immutable policy version on CloudSecure. Each version contains a set of rules that define allowed or denied traffic between sources and destinations.",
			Attributes: map[string]resource_schema.Attribute{
				IDFieldName: idAttribute,
				"description": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "Description of the policy version.",
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					attributeWithMode: attributeWithMode{
						Mode: ImmutableAttributeMode,
					},
				},
				"policy_id": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "ID of the CloudSecure policy this version belongs to.",
						Required:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					attributeWithMode: attributeWithMode{
						Mode: ImmutableAttributeMode,
					},
				},
				"version_number": resource_schema.Int64Attribute{
					Description: "Sequential version number within the policy, for display purposes.",
					Computed:    true,
				},
				"rules": resource_schema.ListNestedAttribute{
					MarkdownDescription: "List of rules in this policy version. Each rule specifies an action (Allow/Deny), source, destination, and port ranges. Multiple rules use OR logic for allows; deny takes precedence.",
					Required:            true,
					PlanModifiers: []planmodifier.List{
						listplanmodifier.RequiresReplace(),
					},
					NestedObject: resource_schema.NestedAttributeObject{
						Attributes: map[string]resource_schema.Attribute{
							"action": resource_schema.StringAttribute{
								Description: "Action to take: Allow or Deny.",
								Required:    true,
								Validators: []validator.String{
									stringvalidator.OneOf("Allow", "Deny"),
								},
							},
							"source": resource_schema.SingleNestedAttribute{
								Description: "Traffic source selector.",
								Required:    true,
								Attributes: map[string]resource_schema.Attribute{
									"k8s": resource_schema.SingleNestedAttribute{
										Description: "K8s workload selector.",
										Optional:    true,
										Attributes:  k8sSelectorAttributes,
										Validators: []validator.Object{
											objectvalidator.ExactlyOneOf(
												path.MatchRelative().AtParent().AtName("k8s"),
												path.MatchRelative().AtParent().AtName("ip_list"),
												path.MatchRelative().AtParent().AtName("cloud"),
												path.MatchRelative().AtParent().AtName("illumio_labels"),
											),
										},
									},
									"ip_list": resource_schema.SingleNestedAttribute{
										Description: "IP list reference selector.",
										Optional:    true,
										Attributes:  ipSelectorAttributes,
										Validators: []validator.Object{
											objectvalidator.ExactlyOneOf(
												path.MatchRelative().AtParent().AtName("k8s"),
												path.MatchRelative().AtParent().AtName("ip_list"),
												path.MatchRelative().AtParent().AtName("cloud"),
												path.MatchRelative().AtParent().AtName("illumio_labels"),
											),
										},
									},
									"cloud": resource_schema.SingleNestedAttribute{
										Description: "Cloud resource selector (e.g., VMs, instances). Not yet implemented.",
										Optional:    true,
										Attributes:  map[string]resource_schema.Attribute{},
										Validators: []validator.Object{
											objectvalidator.ExactlyOneOf(
												path.MatchRelative().AtParent().AtName("k8s"),
												path.MatchRelative().AtParent().AtName("ip_list"),
												path.MatchRelative().AtParent().AtName("cloud"),
												path.MatchRelative().AtParent().AtName("illumio_labels"),
											),
										},
									},
									"illumio_labels": resource_schema.SingleNestedAttribute{
										Description: "Illumio label selector. Not yet implemented.",
										Optional:    true,
										Attributes:  map[string]resource_schema.Attribute{},
										Validators: []validator.Object{
											objectvalidator.ExactlyOneOf(
												path.MatchRelative().AtParent().AtName("k8s"),
												path.MatchRelative().AtParent().AtName("ip_list"),
												path.MatchRelative().AtParent().AtName("cloud"),
												path.MatchRelative().AtParent().AtName("illumio_labels"),
											),
										},
									},
								},
							},
							"destination": resource_schema.SingleNestedAttribute{
								Description: "Traffic destination selector.",
								Required:    true,
								Attributes: map[string]resource_schema.Attribute{
									"k8s": resource_schema.SingleNestedAttribute{
										Description: "K8s workload selector.",
										Optional:    true,
										Attributes:  k8sSelectorAttributes,
										Validators: []validator.Object{
											objectvalidator.ExactlyOneOf(
												path.MatchRelative().AtParent().AtName("k8s"),
												path.MatchRelative().AtParent().AtName("ip_list"),
												path.MatchRelative().AtParent().AtName("fqdns"),
												path.MatchRelative().AtParent().AtName("cloud"),
												path.MatchRelative().AtParent().AtName("illumio_labels"),
											),
										},
									},
									"ip_list": resource_schema.SingleNestedAttribute{
										Description: "IP list reference selector.",
										Optional:    true,
										Attributes:  ipSelectorAttributes,
										Validators: []validator.Object{
											objectvalidator.ExactlyOneOf(
												path.MatchRelative().AtParent().AtName("k8s"),
												path.MatchRelative().AtParent().AtName("ip_list"),
												path.MatchRelative().AtParent().AtName("fqdns"),
												path.MatchRelative().AtParent().AtName("cloud"),
												path.MatchRelative().AtParent().AtName("illumio_labels"),
											),
										},
									},
									"fqdns": resource_schema.SingleNestedAttribute{
										Description: "FQDN selector.",
										Optional:    true,
										Attributes: map[string]resource_schema.Attribute{
											"names": resource_schema.ListAttribute{
												Description: "List of FQDNs (e.g., api.example.com, *.googleapis.com).",
												Required:    true,
												ElementType: types.StringType,
											},
										},
										Validators: []validator.Object{
											objectvalidator.ExactlyOneOf(
												path.MatchRelative().AtParent().AtName("k8s"),
												path.MatchRelative().AtParent().AtName("ip_list"),
												path.MatchRelative().AtParent().AtName("fqdns"),
												path.MatchRelative().AtParent().AtName("cloud"),
												path.MatchRelative().AtParent().AtName("illumio_labels"),
											),
										},
									},
									"cloud": resource_schema.SingleNestedAttribute{
										Description: "Cloud resource selector (e.g., VMs, instances). Not yet implemented.",
										Optional:    true,
										Attributes:  map[string]resource_schema.Attribute{},
										Validators: []validator.Object{
											objectvalidator.ExactlyOneOf(
												path.MatchRelative().AtParent().AtName("k8s"),
												path.MatchRelative().AtParent().AtName("ip_list"),
												path.MatchRelative().AtParent().AtName("fqdns"),
												path.MatchRelative().AtParent().AtName("cloud"),
												path.MatchRelative().AtParent().AtName("illumio_labels"),
											),
										},
									},
									"illumio_labels": resource_schema.SingleNestedAttribute{
										Description: "Illumio label selector. Not yet implemented.",
										Optional:    true,
										Attributes:  map[string]resource_schema.Attribute{},
										Validators: []validator.Object{
											objectvalidator.ExactlyOneOf(
												path.MatchRelative().AtParent().AtName("k8s"),
												path.MatchRelative().AtParent().AtName("ip_list"),
												path.MatchRelative().AtParent().AtName("fqdns"),
												path.MatchRelative().AtParent().AtName("cloud"),
												path.MatchRelative().AtParent().AtName("illumio_labels"),
											),
										},
									},
								},
							},
							"port_ranges": resource_schema.ListNestedAttribute{
								Description: "List of port ranges for the rule.",
								Required:    true,
								NestedObject: resource_schema.NestedAttributeObject{
									Attributes: map[string]resource_schema.Attribute{
										"protocol": resource_schema.StringAttribute{
											Description: "Transport protocol: TCP or UDP.",
											Required:    true,
											Validators: []validator.String{
												stringvalidator.OneOf("TCP", "UDP"),
											},
										},
										"from_port": resource_schema.Int64Attribute{
											Description: "Start port number.",
											Required:    true,
										},
										"to_port": resource_schema.Int64Attribute{
											Description: "End port number.",
											Required:    true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
)
