// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"fmt"
	"maps"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
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

// stringMatchAttributes is one cluster matching criterion expressed as an operator and
// values, in the style of match_expressions on the namespace and workload selectors.
//
// Unlike those, the operator set is limited to In and NotIn. The remaining two ask whether a
// key is present at all, which only means something for labels: a cluster's account, scope,
// and name are structural, so Exists would match every cluster in the provider and
// DoesNotExist none of them. Exists on a required criterion would also satisfy the
// ExactlyOneOf below while constraining nothing, reopening the tenant-wide widening that
// making the criterion required is meant to prevent.
//
// With only In and NotIn left, values is always required, so the rule is a plain per-attribute
// one the framework checks at plan time rather than a cross-field rule deferred to the server.
var stringMatchAttributes = map[string]resource_schema.Attribute{
	"operator": resource_schema.StringAttribute{
		Description: "Match operator. Must be one of: In, NotIn.",
		Required:    true,
		Validators: []validator.String{
			stringvalidator.OneOf("In", "NotIn"),
		},
	},
	"values": resource_schema.ListAttribute{
		Description: "Values to match against. Must contain at least one value.",
		Required:    true,
		ElementType: types.StringType,
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
	},
}

// clusterCriterion returns the two attributes that express one cluster matching criterion.
//
// Each criterion can be written two ways: a plain string for the common case of exact
// equality, or a "<name>_match" object for the operator form. Both are declared Optional
// because the requirement is not "this attribute is set" but "exactly one of the pair is",
// which only a validator can say. This mirrors how the cluster block already expresses the
// choice between id and the four cloud providers.
//
// A required criterion must have one of the two set. An optional one may have neither, which
// leaves that criterion unconstrained and matches any value.
//
// Note that validation of either kind is skipped while a value is unknown, so a criterion
// taking its value from another resource is only checked at apply time rather than at plan.
func clusterCriterion(name, description string, required bool) map[string]resource_schema.Attribute {
	matchName := name + "_match"

	paths := []path.Expression{
		path.MatchRelative().AtParent().AtName(name),
		path.MatchRelative().AtParent().AtName(matchName),
	}

	var (
		stringValidator validator.String
		objectValidator validator.Object
	)
	if required {
		stringValidator = stringvalidator.ExactlyOneOf(paths...)
		objectValidator = objectvalidator.ExactlyOneOf(paths...)
	} else {
		stringValidator = stringvalidator.ConflictsWith(paths...)
		objectValidator = objectvalidator.ConflictsWith(paths...)
	}

	return map[string]resource_schema.Attribute{
		name: resource_schema.StringAttribute{
			Description: fmt.Sprintf("%s Matched for equality. Mutually exclusive with %s.", description, matchName),
			Optional:    true,
			Validators:  []validator.String{stringValidator},
		},
		matchName: resource_schema.SingleNestedAttribute{
			Description: fmt.Sprintf("%s Matched by operator. Mutually exclusive with %s.", description, name),
			Optional:    true,
			Attributes:  stringMatchAttributes,
			Validators:  []validator.Object{objectValidator},
		},
	}
}

// clusterProviderAttributes assembles one cloud provider's criteria into its attribute map.
//
// The two locator criteria are required, since between them they bound a selector to a single
// account and scope. That keeps a selector from silently widening to every cluster in the
// tenant, and it keeps NotIn cheap to evaluate: it is only ever applied within an already
// narrow set rather than against everything.
//
// The cluster name is optional. Omitting it selects every cluster in that account and scope,
// which is the case the operators exist to serve.
func clusterProviderAttributes(accountName, accountDescription, scopeName, scopeDescription, clusterNameDescription string) map[string]resource_schema.Attribute {
	attributes := make(map[string]resource_schema.Attribute, 6)
	maps.Copy(attributes, clusterCriterion(accountName, accountDescription, true))
	maps.Copy(attributes, clusterCriterion(scopeName, scopeDescription, true))
	maps.Copy(attributes, clusterCriterion("cluster_name", clusterNameDescription, false))

	return attributes
}

// clusterSelectorAttributes is shared between source and destination k8s selectors.
var clusterSelectorAttributes = resource_schema.ListNestedAttribute{
	Description: "List of K8s cluster selectors. Each entry selects one or more clusters. A cluster matching any entry matches (OR logic).",
	Required:    true,
	NestedObject: resource_schema.NestedAttributeObject{
		Attributes: map[string]resource_schema.Attribute{
			"id": resource_schema.StringAttribute{
				Description: "Cluster ID (from k8s_cluster resource/data source). Mutually exclusive with aws, gcp, azure, and oci.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("id"),
						path.MatchRelative().AtParent().AtName("aws"),
						path.MatchRelative().AtParent().AtName("gcp"),
						path.MatchRelative().AtParent().AtName("azure"),
						path.MatchRelative().AtParent().AtName("oci"),
					),
				},
			},
			"aws": resource_schema.SingleNestedAttribute{
				Description: "AWS EKS clusters. Mutually exclusive with id, gcp, azure, and oci.",
				Optional:    true,
				Attributes: clusterProviderAttributes(
					"account_id", "AWS account ID.",
					"region", "AWS region (e.g., us-east-1).",
					"EKS cluster name.",
				),
				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("id"),
						path.MatchRelative().AtParent().AtName("aws"),
						path.MatchRelative().AtParent().AtName("gcp"),
						path.MatchRelative().AtParent().AtName("azure"),
						path.MatchRelative().AtParent().AtName("oci"),
					),
				},
			},
			"gcp": resource_schema.SingleNestedAttribute{
				Description: "GCP GKE clusters. Mutually exclusive with id, aws, azure, and oci.",
				Optional:    true,
				Attributes: clusterProviderAttributes(
					"project_id", "GCP project ID.",
					"location", "GKE cluster location (region or zone).",
					"GKE cluster name.",
				),
				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("id"),
						path.MatchRelative().AtParent().AtName("aws"),
						path.MatchRelative().AtParent().AtName("gcp"),
						path.MatchRelative().AtParent().AtName("azure"),
						path.MatchRelative().AtParent().AtName("oci"),
					),
				},
			},
			"azure": resource_schema.SingleNestedAttribute{
				Description: "Azure AKS clusters. Mutually exclusive with id, aws, gcp, and oci.",
				Optional:    true,
				Attributes: clusterProviderAttributes(
					"subscription_id", "Azure subscription ID.",
					"resource_group", "Azure resource group name.",
					"AKS cluster name.",
				),
				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("id"),
						path.MatchRelative().AtParent().AtName("aws"),
						path.MatchRelative().AtParent().AtName("gcp"),
						path.MatchRelative().AtParent().AtName("azure"),
						path.MatchRelative().AtParent().AtName("oci"),
					),
				},
			},
			"oci": resource_schema.SingleNestedAttribute{
				Description: "OCI OKE clusters. Mutually exclusive with id, aws, gcp, and azure.",
				Optional:    true,
				Attributes: clusterProviderAttributes(
					"compartment_id", "OCI compartment OCID.",
					"region", "OCI region (e.g., us-ashburn-1).",
					"OKE cluster name.",
				),
				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("id"),
						path.MatchRelative().AtParent().AtName("aws"),
						path.MatchRelative().AtParent().AtName("gcp"),
						path.MatchRelative().AtParent().AtName("azure"),
						path.MatchRelative().AtParent().AtName("oci"),
					),
				},
			},
		},
	},
}

// k8sSourceSelectorAttributes is used for source.k8s. Source always targets pods.
var k8sSourceSelectorAttributes = map[string]resource_schema.Attribute{
	"clusters": clusterSelectorAttributes,
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

// k8sDestinationSelectorAttributes is used for destination.k8s.
// Destination supports targeting pods, Services, Ingresses, or Gateways (exactly one).
var k8sDestinationSelectorAttributes = map[string]resource_schema.Attribute{
	"clusters": clusterSelectorAttributes,
	"namespace_selector": resource_schema.SingleNestedAttribute{
		Description: "Label selector for K8s namespaces. Use empty {} to match all namespaces.",
		Required:    true,
		Attributes:  labelSelectorAttributes,
	},
	"workload_selector": resource_schema.SingleNestedAttribute{
		Description: "Label selector for K8s workloads (pods). Use empty {} to match all pods. Mutually exclusive with services, ingresses, and gateways.",
		Optional:    true,
		Attributes:  workloadSelectorAttributes,
		Validators: []validator.Object{
			objectvalidator.ExactlyOneOf(
				path.MatchRelative().AtParent().AtName("workload_selector"),
				path.MatchRelative().AtParent().AtName("services"),
				path.MatchRelative().AtParent().AtName("ingresses"),
				path.MatchRelative().AtParent().AtName("gateways"),
			),
		},
	},
	"services": resource_schema.ListAttribute{
		Description: "List of K8s Service names to target. Mutually exclusive with workload_selector, ingresses, and gateways.",
		Optional:    true,
		ElementType: types.StringType,
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
	},
	"ingresses": resource_schema.ListAttribute{
		Description: "List of K8s Ingress names to target. Mutually exclusive with workload_selector, services, and gateways.",
		Optional:    true,
		ElementType: types.StringType,
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
	},
	"gateways": resource_schema.ListAttribute{
		Description: "List of K8s Gateway names to target. Mutually exclusive with workload_selector, services, and ingresses.",
		Optional:    true,
		ElementType: types.StringType,
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
	},
}

// ipSelectorAttributes is reused for source.ip_lists and destination.ip_lists.
var ipSelectorAttributes = map[string]resource_schema.Attribute{
	"ids": resource_schema.ListAttribute{
		Description: "List of illumio-cloudsecure_ip_list resource IDs. Traffic matching any list is selected (OR logic).",
		Required:    true,
		ElementType: types.StringType,
	},
}

// azureOrgSelectorAttributes defines the Azure organizational boundary selector.
var azureOrgSelectorAttributes = map[string]resource_schema.Attribute{
	"subscriptions": resource_schema.ListNestedAttribute{
		Description: "List of Azure subscriptions. Targets all resources within each subscription.",
		Required:    true,
		NestedObject: resource_schema.NestedAttributeObject{
			Attributes: map[string]resource_schema.Attribute{
				"id": resource_schema.StringAttribute{
					Description: "Azure subscription GUID.",
					Required:    true,
				},
			},
		},
	},
}

// azureVirtualNetworkAttributes defines the VNet selector.
var azureVirtualNetworkAttributes = map[string]resource_schema.Attribute{
	"subscription_id": resource_schema.StringAttribute{
		Description: "Azure subscription ID containing the VNet.",
		Required:    true,
	},
	"resource_group": resource_schema.StringAttribute{
		Description: "Azure resource group name containing the VNet.",
		Required:    true,
	},
	"id": resource_schema.StringAttribute{
		Description: "Azure VNet name.",
		Required:    true,
	},
}

// azureSubnetAttributes defines the Subnet selector.
var azureSubnetAttributes = map[string]resource_schema.Attribute{
	"id": resource_schema.StringAttribute{
		Description: "Full Azure resource ID for the Subnet (csp_id).",
		Required:    true,
	},
}

// azureNetworkAttributes defines the network resource selectors within Azure.
// At least one of vnets or subnets must be set.
var azureNetworkAttributes = map[string]resource_schema.Attribute{
	"vnets": resource_schema.ListNestedAttribute{
		Description: "Optional list of VNet selectors. At least one of vnets or subnets must be specified.",
		Optional:    true,
		NestedObject: resource_schema.NestedAttributeObject{
			Attributes: azureVirtualNetworkAttributes,
		},
		Validators: []validator.List{
			listvalidator.AtLeastOneOf(
				path.MatchRelative().AtParent().AtName("vnets"),
				path.MatchRelative().AtParent().AtName("subnets"),
			),
		},
	},
	"subnets": resource_schema.ListNestedAttribute{
		Description: "Optional list of Subnet selectors. At least one of vnets or subnets must be specified.",
		Optional:    true,
		NestedObject: resource_schema.NestedAttributeObject{
			Attributes: azureSubnetAttributes,
		},
		Validators: []validator.List{
			listvalidator.AtLeastOneOf(
				path.MatchRelative().AtParent().AtName("vnets"),
				path.MatchRelative().AtParent().AtName("subnets"),
			),
		},
	},
}

// cloudAzureSelectorAttributes defines the Azure cloud resource selector.
// Exactly one of org_selector or network must be set.
var cloudAzureSelectorAttributes = map[string]resource_schema.Attribute{
	"org_selector": resource_schema.SingleNestedAttribute{
		Description: "Targets all resources within the specified subscriptions. Mutually exclusive with network.",
		Optional:    true,
		Attributes:  azureOrgSelectorAttributes,
		Validators: []validator.Object{
			objectvalidator.ExactlyOneOf(
				path.MatchRelative().AtParent().AtName("org_selector"),
				path.MatchRelative().AtParent().AtName("network"),
			),
		},
	},
	"network": resource_schema.SingleNestedAttribute{
		Description: "Targets specific VNets or Subnets. Mutually exclusive with org_selector.",
		Optional:    true,
		Attributes:  azureNetworkAttributes,
		Validators: []validator.Object{
			objectvalidator.ExactlyOneOf(
				path.MatchRelative().AtParent().AtName("org_selector"),
				path.MatchRelative().AtParent().AtName("network"),
			),
		},
	},
}

// awsOrgSelectorAttributes defines the AWS organizational boundary selector.
var awsOrgSelectorAttributes = map[string]resource_schema.Attribute{
	"accounts": resource_schema.ListNestedAttribute{
		Description: "List of AWS accounts. Targets all resources within each account.",
		Required:    true,
		NestedObject: resource_schema.NestedAttributeObject{
			Attributes: map[string]resource_schema.Attribute{
				"id": resource_schema.StringAttribute{
					Description: "AWS account ID (12-digit number, e.g. 123456789012).",
					Required:    true,
				},
			},
		},
	},
}

// awsVpcAttributes defines the VPC selector.
var awsVpcAttributes = map[string]resource_schema.Attribute{
	"id": resource_schema.StringAttribute{
		Description: "VPC ID (e.g., vpc-0a1b2c3d4e5f67890).",
		Required:    true,
	},
	"account_id": resource_schema.StringAttribute{
		Description: "AWS account ID containing the VPC.",
		Required:    true,
	},
	"region": resource_schema.StringAttribute{
		Description: "AWS region containing the VPC (e.g., us-east-1).",
		Required:    true,
	},
}

// awsSubnetAttributes defines the AWS Subnet selector.
var awsSubnetAttributes = map[string]resource_schema.Attribute{
	"id": resource_schema.StringAttribute{
		Description: "Subnet ID (e.g., subnet-0a1b2c3d4e5f67890).",
		Required:    true,
	},
	"account_id": resource_schema.StringAttribute{
		Description: "AWS account ID containing the subnet.",
		Required:    true,
	},
	"region": resource_schema.StringAttribute{
		Description: "AWS region containing the subnet (e.g., us-east-1).",
		Required:    true,
	},
}

// awsNetworkAttributes defines the network resource selectors within AWS.
// At least one of vpcs or subnets must be set.
var awsNetworkAttributes = map[string]resource_schema.Attribute{
	"vpcs": resource_schema.ListNestedAttribute{
		Description: "Optional list of VPC selectors. At least one of vpcs or subnets must be specified.",
		Optional:    true,
		NestedObject: resource_schema.NestedAttributeObject{
			Attributes: awsVpcAttributes,
		},
		Validators: []validator.List{
			listvalidator.AtLeastOneOf(
				path.MatchRelative().AtParent().AtName("vpcs"),
				path.MatchRelative().AtParent().AtName("subnets"),
			),
		},
	},
	"subnets": resource_schema.ListNestedAttribute{
		Description: "Optional list of Subnet selectors. At least one of vpcs or subnets must be specified.",
		Optional:    true,
		NestedObject: resource_schema.NestedAttributeObject{
			Attributes: awsSubnetAttributes,
		},
		Validators: []validator.List{
			listvalidator.AtLeastOneOf(
				path.MatchRelative().AtParent().AtName("vpcs"),
				path.MatchRelative().AtParent().AtName("subnets"),
			),
		},
	},
}

// cloudAwsSelectorAttributes defines the AWS cloud resource selector.
// Exactly one of org_selector or network must be set.
var cloudAwsSelectorAttributes = map[string]resource_schema.Attribute{
	"org_selector": resource_schema.SingleNestedAttribute{
		Description: "Targets all resources within the specified accounts. Mutually exclusive with network.",
		Optional:    true,
		Attributes:  awsOrgSelectorAttributes,
		Validators: []validator.Object{
			objectvalidator.ExactlyOneOf(
				path.MatchRelative().AtParent().AtName("org_selector"),
				path.MatchRelative().AtParent().AtName("network"),
			),
		},
	},
	"network": resource_schema.SingleNestedAttribute{
		Description: "Targets specific VPCs or Subnets. Mutually exclusive with org_selector.",
		Optional:    true,
		Attributes:  awsNetworkAttributes,
		Validators: []validator.Object{
			objectvalidator.ExactlyOneOf(
				path.MatchRelative().AtParent().AtName("org_selector"),
				path.MatchRelative().AtParent().AtName("network"),
			),
		},
	},
}

// cloudSelectorAttributes is reused for source.cloud and destination.cloud.
var cloudSelectorAttributes = map[string]resource_schema.Attribute{
	"azure": resource_schema.SingleNestedAttribute{
		Description: "Azure cloud resource selector. Mutually exclusive with aws.",
		Optional:    true,
		Attributes:  cloudAzureSelectorAttributes,
		Validators: []validator.Object{
			objectvalidator.ExactlyOneOf(
				path.MatchRelative().AtParent().AtName("azure"),
				path.MatchRelative().AtParent().AtName("aws"),
			),
		},
	},
	"aws": resource_schema.SingleNestedAttribute{
		Description: "AWS cloud resource selector. Mutually exclusive with azure.",
		Optional:    true,
		Attributes:  cloudAwsSelectorAttributes,
		Validators: []validator.Object{
			objectvalidator.ExactlyOneOf(
				path.MatchRelative().AtParent().AtName("azure"),
				path.MatchRelative().AtParent().AtName("aws"),
			),
		},
	},
}

// illumioLabelsSelectorAttributes is reused for source.illumio_labels and destination.illumio_labels.
//
// The selector is a restricted conjunctive normal form: values within an entry are
// OR'd, and entries are AND'd. So [{role, [web, api]}, {env, [prod]}] selects
// workloads that are (role=web OR role=api) AND env=prod. A clause spanning two
// keys, e.g. (role=web OR env=prod), is not expressible.
var illumioLabelsSelectorAttributes = map[string]resource_schema.Attribute{
	"labels": resource_schema.ListNestedAttribute{
		Description: "List of Illumio label requirements. Values within an entry are OR'd; entries are AND'd.",
		Required:    true,
		NestedObject: resource_schema.NestedAttributeObject{
			Attributes: map[string]resource_schema.Attribute{
				"key": resource_schema.StringAttribute{
					Description: "Illumio label key.",
					Required:    true,
				},
				"values": resource_schema.ListAttribute{
					Description: "Accepted values for the key. Any value matches (OR logic).",
					Required:    true,
					ElementType: types.StringType,
					Validators: []validator.List{
						listvalidator.SizeAtLeast(1),
						listvalidator.UniqueValues(),
					},
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
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
				"rules": MapNestedResourceAttributeWithMode{
					MapNestedAttribute: resource_schema.MapNestedAttribute{
						MarkdownDescription: "Map of rules in this policy version, keyed by rule name. Each rule specifies an action (Allow/Deny), source, destination, and port ranges. Multiple rules use OR logic for allows; deny takes precedence.",
						Required:            true,
						PlanModifiers: []planmodifier.Map{
							mapplanmodifier.RequiresReplace(),
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
											Attributes:  k8sSourceSelectorAttributes,
											Validators: []validator.Object{
												objectvalidator.ExactlyOneOf(
													path.MatchRelative().AtParent().AtName("k8s"),
													path.MatchRelative().AtParent().AtName("ip_lists"),
													path.MatchRelative().AtParent().AtName("cloud"),
													path.MatchRelative().AtParent().AtName("illumio_labels"),
												),
											},
										},
										"ip_lists": resource_schema.SingleNestedAttribute{
											Description: "IP list reference selector.",
											Optional:    true,
											Attributes:  ipSelectorAttributes,
											Validators: []validator.Object{
												objectvalidator.ExactlyOneOf(
													path.MatchRelative().AtParent().AtName("k8s"),
													path.MatchRelative().AtParent().AtName("ip_lists"),
													path.MatchRelative().AtParent().AtName("cloud"),
													path.MatchRelative().AtParent().AtName("illumio_labels"),
												),
											},
										},
										"cloud": resource_schema.SingleNestedAttribute{
											Description: "Cloud resource selector (Azure or AWS).",
											Optional:    true,
											Attributes:  cloudSelectorAttributes,
											Validators: []validator.Object{
												objectvalidator.ExactlyOneOf(
													path.MatchRelative().AtParent().AtName("k8s"),
													path.MatchRelative().AtParent().AtName("ip_lists"),
													path.MatchRelative().AtParent().AtName("cloud"),
													path.MatchRelative().AtParent().AtName("illumio_labels"),
												),
											},
										},
										"illumio_labels": resource_schema.SingleNestedAttribute{
											Description: "Illumio label selector.",
											Optional:    true,
											Attributes:  illumioLabelsSelectorAttributes,
											Validators: []validator.Object{
												objectvalidator.ExactlyOneOf(
													path.MatchRelative().AtParent().AtName("k8s"),
													path.MatchRelative().AtParent().AtName("ip_lists"),
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
											Description: "K8s destination selector. Target pods, Services, Ingresses, or Gateways.",
											Optional:    true,
											Attributes:  k8sDestinationSelectorAttributes,
											Validators: []validator.Object{
												objectvalidator.ExactlyOneOf(
													path.MatchRelative().AtParent().AtName("k8s"),
													path.MatchRelative().AtParent().AtName("ip_lists"),
													path.MatchRelative().AtParent().AtName("fqdns"),
													path.MatchRelative().AtParent().AtName("cloud"),
													path.MatchRelative().AtParent().AtName("illumio_labels"),
												),
											},
										},
										"ip_lists": resource_schema.SingleNestedAttribute{
											Description: "IP list reference selector.",
											Optional:    true,
											Attributes:  ipSelectorAttributes,
											Validators: []validator.Object{
												objectvalidator.ExactlyOneOf(
													path.MatchRelative().AtParent().AtName("k8s"),
													path.MatchRelative().AtParent().AtName("ip_lists"),
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
													path.MatchRelative().AtParent().AtName("ip_lists"),
													path.MatchRelative().AtParent().AtName("fqdns"),
													path.MatchRelative().AtParent().AtName("cloud"),
													path.MatchRelative().AtParent().AtName("illumio_labels"),
												),
											},
										},
										"cloud": resource_schema.SingleNestedAttribute{
											Description: "Cloud resource selector (Azure or AWS).",
											Optional:    true,
											Attributes:  cloudSelectorAttributes,
											Validators: []validator.Object{
												objectvalidator.ExactlyOneOf(
													path.MatchRelative().AtParent().AtName("k8s"),
													path.MatchRelative().AtParent().AtName("ip_lists"),
													path.MatchRelative().AtParent().AtName("fqdns"),
													path.MatchRelative().AtParent().AtName("cloud"),
													path.MatchRelative().AtParent().AtName("illumio_labels"),
												),
											},
										},
										"illumio_labels": resource_schema.SingleNestedAttribute{
											Description: "Illumio label selector.",
											Optional:    true,
											Attributes:  illumioLabelsSelectorAttributes,
											Validators: []validator.Object{
												objectvalidator.ExactlyOneOf(
													path.MatchRelative().AtParent().AtName("k8s"),
													path.MatchRelative().AtParent().AtName("ip_lists"),
													path.MatchRelative().AtParent().AtName("fqdns"),
													path.MatchRelative().AtParent().AtName("cloud"),
													path.MatchRelative().AtParent().AtName("illumio_labels"),
												),
											},
										},
									},
								},
								"port_ranges": resource_schema.ListAttribute{
									Description: "List of port ranges for the rule.",
									Required:    true,
									ElementType: PortRange,
								},
							},
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
