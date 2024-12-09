package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var MetadataType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"cloud": types.StringType,
		"id":    types.StringType,
	},
}
var TagType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"key":   types.StringType,
		"value": types.StringType,
		"cloud": types.StringType,
	},
}

var (
	deploymentResource = Resource{
		TypeName: "deployment",
		Schema: resource_schema.Schema{
			Version:     1,
			Description: "Groups resources under Cloudsecure deployments",
			Attributes: map[string]resource_schema.Attribute{
				IDFieldName: idAttribute,
				"environment": StringResourceAttributeWithMode{
					StringAttribute: resource_schema.StringAttribute{
						Description: "Deployment Name",
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
				"description": resource_schema.StringAttribute{
					MarkdownDescription: "Description of Cloudsecure deployment",
					Optional:            true,
				},
				"accounts": resource_schema.ListAttribute{
					Optional:    true,
					Description: "cloud accounts that need be a part of this deployment",
					ElementType: MetadataType,
				},
				"regions": resource_schema.ListAttribute{
					Optional:    true,
					Description: "regions that need to be associated with this deployment",
					ElementType: MetadataType,
				},
				"vnets": resource_schema.ListAttribute{
					Optional:    true,
					Description: "virtual networks that need to be associated with this deployment",
					ElementType: MetadataType,
				},
				"subnets": resource_schema.ListAttribute{
					Optional:    true,
					Description: "subnets whose resources need to be associated with this deployment",
					ElementType: MetadataType,
				},
				"tags": resource_schema.ListAttribute{
					Optional:    true,
					Description: "tags of resources that need to be grouped under this deployment",
					ElementType: TagType,
				},
			},
		},
	}
)
