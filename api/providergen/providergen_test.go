// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"bytes"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/illumio/terraform-provider-illumio-cloudsecure/api/schema"
	"github.com/stretchr/testify/suite"
)

type GenerateProviderTestSuite struct {
	suite.Suite
}

func TestGenerateProviderSuite(t *testing.T) {
	suite.Run(t, new(GenerateProviderTestSuite))
}

func (suite *GenerateProviderTestSuite) TestGenerateProviderDataGenerator() {
	// Setup test schema
	testResource := schema.Resource{
		TypeName: "test_object",
		Schema: resource_schema.Schema{
			Version:     1,
			Description: "Manages an AWS account in CloudSecure.",
			Attributes: map[string]resource_schema.Attribute{
				"id": resource_schema.StringAttribute{
					Description: "CloudSecure ID.",
					Computed:    true,
				},
				"names": resource_schema.ListAttribute{
					ElementType: types.StringType,
					Description: "List of names.",
					Required:    true,
				},
				"address": resource_schema.ObjectAttribute{
					AttributeTypes: map[string]attr.Type{
						"city":    types.StringType,
						"state":   types.StringType,
						"pincode": types.Int64Type,
						"phone_numbers": types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"home":   types.StringType,
								"office": types.StringType,
							},
						},
					},
					Required:    true,
					Description: "Address attribute.",
				},
				"rules": resource_schema.ListAttribute{
					ElementType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"source":      types.StringType,
							"destination": types.StringType,
							"port":        types.Int64Type,
							"metadata": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"tag": types.StringType,
									"key": types.StringType,
								},
							},
						},
					},
					Required:    true,
					Description: "Rules attribute.",
				},
			},
		},
	}

	countCount := 1

	data := providerTemplateData{
		Package:               "testpkg",
		ProviderTypeName:      "Provider",
		Models:                make([]model, 0),
		NewRequestFuncs:       make([]convertFunc, 0, countCount*3),
		NewUpdateRequestFuncs: make([]convertFunc, 0, countCount),
		CopyResponseFuncs:     make([]convertFunc, 0, countCount*3),
		Resources:             make([]resourceData, 0, countCount),
	}

	err := AddResourceToProviderTemplateData(&testResource, &data, "TestObject", "TestObject")
	// Assert no error
	suite.Require().NoError(err, "AddResourceToProviderTemplateData should not return an error")

	var buffer bytes.Buffer
	// Assert the output is not empty
	err = providerTemplate.Execute(&buffer, &data)
	suite.Require().NoError(err, "providerTemplate.Execute should not return an error")

	output := buffer.String()
	suite.NotEmpty(output, "Generated provider output should not be empty")

	// Basic content check
	suite.Contains(output, "package testpkg", "Generated provider should include the correct package name")
	suite.Contains(output, "type TestObjectResource struct", "Generated provider should include the resource model definition")
	suite.Contains(output, "func NewTestObjectResource", "Generated provider should include the resource creation function")
	suite.Contains(output, "github.com/hashicorp/terraform-plugin-log/tflog", "Generated provider should include the resource creation function")

	// Check for the generated resource models
	suite.Len(data.Models, 1, "Number of models should match the number of resources")
	suite.Contains(data.Models[0].Name, "TestObjectResourceModel", "Generated provider should include the correct resource model name")
}

func (suite *GenerateProviderTestSuite) TestListOfObjects() {
	testResource := schema.Resource{
		TypeName: "aws_tag_to_label",
		Schema: resource_schema.Schema{
			Version:     1,
			Description: "Manages an AWS account in CloudSecure.",
			Attributes: map[string]resource_schema.Attribute{
				"id": resource_schema.StringAttribute{
					Description: "CloudSecure ID.",
				},
				"icon": resource_schema.ObjectAttribute{
					AttributeTypes: map[string]attr.Type{
						"name":             types.StringType,
						"background_color": types.StringType,
						"foreground_color": types.StringType,
					},
					Required:    true,
					Description: "Icon details.",
				},
				"cloud_tags": resource_schema.ListAttribute{
					Required:    true,
					Description: "List of AWS account tags to map to the CloudSecure label.",
					ElementType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"key":   types.StringType,
							"cloud": types.StringType,
						},
					},
				},
			},
		},
	}

	countCount := 1
	data := providerTemplateData{
		Package:               "testpkg",
		ProviderTypeName:      "Provider",
		Models:                make([]model, 0),
		NewRequestFuncs:       make([]convertFunc, 0, countCount*3),
		NewUpdateRequestFuncs: make([]convertFunc, 0, countCount),
		CopyResponseFuncs:     make([]convertFunc, 0, countCount*3),
		Resources:             make([]resourceData, 0, countCount),
	}

	err := AddResourceToProviderTemplateData(&testResource, &data, "AwsTagToLabel", "AwsTagToLabel")
	// Assert no error
	suite.Require().NoError(err, "AddResourceToProviderTemplateData should not return an error")
}

func (suite *GenerateProviderTestSuite) TestSetsAndMore() {
	var cloudTagsAttribute = resource_schema.ListAttribute{
		ElementType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"key":   types.StringType,
				"cloud": types.StringType,
			},
		},
		Required:    true,
		Description: "List of AWS account tags to map to the CloudSecure label.",
	}

	var setObjAttribute = resource_schema.SetAttribute{
		ElementType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"set_key": types.StringType,
				"set_val": types.StringType,
			},
		},
		Required:    true,
		Description: "List of AWS account tags to map to the CloudSecure label.",
	}

	var iconAttribute = resource_schema.ObjectAttribute{
		AttributeTypes: map[string]attr.Type{
			"name":             types.StringType,
			"background_color": types.StringType,
			"foreground_color": types.StringType,
		},
		Required:    true,
		Description: "Icon details.",
	}

	var objInObjAttribute = resource_schema.ObjectAttribute{
		AttributeTypes: map[string]attr.Type{
			"name": types.StringType,
			"child": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"name": types.StringType,
					"grand_child": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"name": types.StringType,
						},
					},
				},
			},
		},
		Required:    true,
		Description: "Icon details.",
	}

	testResource := schema.Resource{
		TypeName: "nested_object_tester",
		Schema: resource_schema.Schema{
			Version:     1,
			Description: "Maps AWS account tags to CloudSecure labels.",
			Attributes: map[string]resource_schema.Attribute{
				"key": resource_schema.StringAttribute{
					MarkdownDescription: "CloudSecure label key.",
					Required:            true,
				},
				"name": resource_schema.StringAttribute{
					MarkdownDescription: "CloudSecure label display name.",
					Required:            true,
				},
				"icon":       iconAttribute,
				"cloud_tags": cloudTagsAttribute,
				"set_obj":    setObjAttribute,
				"set_of_sets_string": resource_schema.SetAttribute{
					ElementType: types.SetType{
						ElemType: types.SetType{
							ElemType: types.StringType,
						},
					},
					Optional: true,
				},
				"set_of_strings": resource_schema.SetAttribute{
					ElementType: types.StringType,
					Optional:    true,
				},
				"obj_in_obj": objInObjAttribute,
			},
		},
	}

	countCount := 1

	data := providerTemplateData{
		Package:               "testpkg",
		ProviderTypeName:      "Provider",
		Models:                make([]model, 0),
		NewRequestFuncs:       make([]convertFunc, 0, countCount*3),
		NewUpdateRequestFuncs: make([]convertFunc, 0, countCount),
		CopyResponseFuncs:     make([]convertFunc, 0, countCount*3),
		Resources:             make([]resourceData, 0, countCount),
	}

	err := AddResourceToProviderTemplateData(&testResource, &data, "TestObject", "TestObject")
	// Assert no error
	suite.Require().NoError(err, "AddResourceToProviderTemplateData should not return an error")

	dst := new(bytes.Buffer)
	err = ProviderConvertersTemplate.Execute(dst, &data)
	suite.Require().NoError(err, "ProviderConvertersTemplate.Execute should not return an error")

	suite.NotEmpty(dst.String(), "Generated provider output should not be empty")
}
