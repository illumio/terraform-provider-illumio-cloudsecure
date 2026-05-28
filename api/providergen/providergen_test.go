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

func (suite *GenerateProviderTestSuite) TestConvertCollectionsOfObjectsSucceeds() {
	testResource := schema.Resource{
		TypeName: "dummy",
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
					ElementType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"key":   types.StringType,
							"cloud": types.StringType,
						},
					},
				},
				"labels": resource_schema.MapAttribute{
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

	err := AddResourceToProviderTemplateData(&testResource, &data, "Dummy", "Dummy")

	suite.Require().NoError(err, "AddResourceToProviderTemplateData should not return an error")
}

// TestListAndMapNestedAttribute verifies that ListNestedAttribute and MapNestedAttribute generate correct:
// 1. Proto type names for nested objects.
// 2. Converter code for List/Set/Map of objects.
func (suite *GenerateProviderTestSuite) TestListAndMapNestedAttribute() { //nolint:maintidx
	testResource := schema.Resource{
		TypeName: "policy_version",
		Schema: resource_schema.Schema{
			Version:     1,
			Description: "Manages a policy version.",
			Attributes: map[string]resource_schema.Attribute{
				"id": resource_schema.StringAttribute{
					Description: "Policy version ID.",
					Computed:    true,
				},
				"spec": resource_schema.SingleNestedAttribute{
					Description: "Policy specification.",
					Required:    true,
					Attributes: map[string]resource_schema.Attribute{
						"rules": resource_schema.ListNestedAttribute{
							Description: "List of rules.",
							Optional:    true,
							NestedObject: resource_schema.NestedAttributeObject{
								Attributes: map[string]resource_schema.Attribute{
									"action": resource_schema.StringAttribute{
										Description: "Action to take.",
										Required:    true,
									},
									"destination": resource_schema.SingleNestedAttribute{
										Description: "Traffic destination.",
										Optional:    true,
										Attributes: map[string]resource_schema.Attribute{
											"k8s": resource_schema.SingleNestedAttribute{
												Description: "K8s workload selector.",
												Optional:    true,
												Attributes: map[string]resource_schema.Attribute{
													"cluster_name": resource_schema.StringAttribute{
														Description: "Cluster name.",
														Optional:    true,
													},
												},
											},
										},
									},
								},
							},
						},
						"rules_by_name": resource_schema.MapNestedAttribute{
							Description: "Map",
							Optional:    true,
							NestedObject: resource_schema.NestedAttributeObject{
								Attributes: map[string]resource_schema.Attribute{
									"action": resource_schema.StringAttribute{
										Description: "Action to take.",
										Required:    true,
									},
									"destination": resource_schema.SingleNestedAttribute{
										Description: "Traffic destination.",
										Optional:    true,
										Attributes: map[string]resource_schema.Attribute{
											"k8s": resource_schema.SingleNestedAttribute{
												Description: "K8s workload selector.",
												Optional:    true,
												Attributes: map[string]resource_schema.Attribute{
													"cluster_name": resource_schema.StringAttribute{
														Description: "Cluster name.",
														Optional:    true,
													},
												},
											},
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

	data := providerTemplateData{
		Package:               "testpkg",
		ProviderTypeName:      "Provider",
		Models:                make([]model, 0),
		NewRequestFuncs:       make([]convertFunc, 0, 3),
		NewUpdateRequestFuncs: make([]convertFunc, 0, 1),
		CopyResponseFuncs:     make([]convertFunc, 0, 3),
		Resources:             make([]resourceData, 0, 1),
	}

	err := AddResourceToProviderTemplateData(&testResource, &data, "PolicyVersion", "PolicyVersion")
	suite.Require().NoError(err, "AddResourceToProviderTemplateData should not return an error")

	// Verify the model was created
	suite.Require().Len(data.Models, 1, "Should have one resource model")
	resourceModel := data.Models[0]

	// Find the "spec" field
	var specField *field

	for i := range resourceModel.Fields {
		if resourceModel.Fields[i].AttributeName == "spec" {
			specField = &resourceModel.Fields[i]

			break
		}
	}

	suite.Require().NotNil(specField, "Should have 'spec' field")
	suite.Require().NotNil(specField.Type.NestedModel, "spec should have NestedModel")
	specModel := specField.Type.NestedModel
	suite.Equal("PolicyVersion_Spec", specModel.Name, "spec model name")

	// Find the "rules" field inside spec
	var rulesField *field

	for i := range specModel.Fields {
		if specModel.Fields[i].AttributeName == "rules" {
			rulesField = &specModel.Fields[i]

			break
		}
	}

	suite.Require().NotNil(rulesField, "Should have 'rules' field in spec")

	// rules is a List, so check CollectionElementType for the nested object
	suite.Require().NotNil(rulesField.Type.CollectionElementType, "rules should have CollectionElementType")
	rulesElemType := rulesField.Type.CollectionElementType
	suite.Require().NotNil(rulesElemType.NestedModel, "rules element should have NestedModel")
	suite.Equal("PolicyVersion_Spec_Rules", rulesElemType.NestedModel.Name, "rules element model name")

	// Find "destination" inside the rules model
	rulesModel := rulesElemType.NestedModel

	var destField *field

	for i := range rulesModel.Fields {
		if rulesModel.Fields[i].AttributeName == "destination" {
			destField = &rulesModel.Fields[i]

			break
		}
	}

	suite.Require().NotNil(destField, "Should have 'destination' field in rules")
	suite.Require().NotNil(destField.Type.NestedModel, "destination should have NestedModel")
	suite.Equal("PolicyVersion_Spec_Rules_Destination", destField.Type.NestedModel.Name, "destination model name")

	// Find "k8s" inside the destination model
	destModel := destField.Type.NestedModel

	var k8sField *field

	for i := range destModel.Fields {
		if destModel.Fields[i].AttributeName == "k8s" {
			k8sField = &destModel.Fields[i]

			break
		}
	}

	suite.Require().NotNil(k8sField, "Should have 'k8s' field in destination")
	suite.Require().NotNil(k8sField.Type.NestedModel, "k8s should have NestedModel")
	suite.Equal("PolicyVersion_Spec_Rules_Destination_K8S", k8sField.Type.NestedModel.Name, "k8s model name")

	// Find the "rules_by_name" field inside spec
	var rulesByNameField *field

	for i := range specModel.Fields {
		if specModel.Fields[i].AttributeName == "rules_by_name" {
			rulesByNameField = &specModel.Fields[i]

			break
		}
	}

	suite.Require().NotNil(rulesByNameField, "Should have 'rules_by_name' field in spec")

	// rules_by_name is a Map, so check CollectionElementType for the nested object
	suite.Require().NotNil(rulesByNameField.Type.CollectionElementType, "rules_by_name should have CollectionElementType")
	rulesByNameElemType := rulesByNameField.Type.CollectionElementType
	suite.Require().NotNil(rulesByNameElemType.NestedModel, "rules_by_name element should have NestedModel")
	suite.Equal("PolicyVersion_Spec_RulesByName", rulesByNameElemType.NestedModel.Name, "rules_by_name element model name")

	// Find "destination" inside the rules_by_name model
	rulesByNameModel := rulesByNameElemType.NestedModel

	for i := range rulesByNameModel.Fields {
		if rulesByNameModel.Fields[i].AttributeName == "destination" {
			destField = &rulesByNameModel.Fields[i]

			break
		}
	}

	suite.Require().NotNil(destField, "Should have 'destination' field in rules")
	suite.Require().NotNil(destField.Type.NestedModel, "destination should have NestedModel")
	suite.Equal("PolicyVersion_Spec_RulesByName_Destination", destField.Type.NestedModel.Name, "destination model name")

	// Find "k8s" inside the destination model
	destModel = destField.Type.NestedModel

	for i := range destModel.Fields {
		if destModel.Fields[i].AttributeName == "k8s" {
			k8sField = &destModel.Fields[i]

			break
		}
	}

	suite.Require().NotNil(k8sField, "Should have 'k8s' field in destination")
	suite.Require().NotNil(k8sField.Type.NestedModel, "k8s should have NestedModel")
	suite.Equal("PolicyVersion_Spec_RulesByName_Destination_K8S", k8sField.Type.NestedModel.Name, "k8s model name")

	// Test converter generation
	dst := new(bytes.Buffer)
	err = ProviderConvertersTemplate.Execute(dst, &data)
	suite.Require().NoError(err, "ProviderConvertersTemplate.Execute should not return an error")

	output := dst.String()

	// Verify GetTypeAttrsFor generates correct List type with nested object element type
	expectedGetTypeAttrs := `
func GetTypeAttrsForPolicyVersion_Spec() map[string]attr.Type {
	return map[string]attr.Type{
		"rules": types.ListType{ElemType: types.ObjectType{
			AttrTypes: GetTypeAttrsForPolicyVersion_Spec_Rules(),
		}},
		"rules_by_name": types.MapType{ElemType: types.ObjectType{
			AttrTypes: GetTypeAttrsForPolicyVersion_Spec_RulesByName(),
		}},
	}
`

	suite.Contains(output, expectedGetTypeAttrs)

	// Verify ConvertToObjectValueFromProto iterates over list and calls nested converter for each element
	expectedConvertToObjectValue := `
func ConvertPolicyVersion_SpecToObjectValueFromProto(proto *configv1.PolicyVersion_Spec) basetypes.ObjectValue  {
	if proto == nil {
		return types.ObjectNull(GetTypeAttrsForPolicyVersion_Spec())
	}
	elementsInRules := make([]attr.Value, 0, len(proto.Rules))
	for _, item := range proto.Rules {
		elementsInRules = append(elementsInRules, ConvertPolicyVersion_Spec_RulesToObjectValueFromProto(item))
	}
	elementsInRulesByName := make(map[string]attr.Value, len(proto.RulesByName))
	for key, item := range proto.RulesByName {
		elementsInRulesByName[key] = ConvertPolicyVersion_Spec_RulesByNameToObjectValueFromProto(item)
	}
	return types.ObjectValueMust(
		GetTypeAttrsForPolicyVersion_Spec(),
		map[string]attr.Value{
			"rules": func() basetypes.ListValue { if proto.Rules == nil { return types.ListNull(types.ObjectType{AttrTypes: GetTypeAttrsForPolicyVersion_Spec_Rules()}) }; return types.ListValueMust(types.ObjectType{AttrTypes: GetTypeAttrsForPolicyVersion_Spec_Rules()}, elementsInRules) }(),
			"rules_by_name": func() basetypes.MapValue { if proto.RulesByName == nil { return types.MapNull(types.ObjectType{AttrTypes: GetTypeAttrsForPolicyVersion_Spec_RulesByName()}) }; return types.MapValueMust(types.ObjectType{AttrTypes: GetTypeAttrsForPolicyVersion_Spec_RulesByName()}, elementsInRulesByName) }(),
		},
	)
}
`
	suite.Contains(output, expectedConvertToObjectValue)

	// Verify ConvertDataValueToProto uses Elements() and iterates to convert each element back to proto
	expectedConvertDataValue := `
func ConvertDataValueToPolicyVersion_SpecProto(ctx context.Context, dataValue attr.Value) (*configv1.PolicyVersion_Spec, diag.Diagnostics) {
	if dataValue.IsNull() || dataValue.IsUnknown() {
		return nil, nil
	}
	pv := PolicyVersion_Spec{}
	diags := tfsdk.ValueAs(ctx, dataValue, &pv)
	if diags.HasError() {
		return nil, diags
	}
	proto := &configv1.PolicyVersion_Spec{}
	pvElemModelRules := pv.Rules.Elements()
	proto.Rules = make([]*configv1.PolicyVersion_Spec_Rules, 0, len(pvElemModelRules))
	for _, elem := range pvElemModelRules {
		pvModel, dvDiags := ConvertDataValueToPolicyVersion_Spec_RulesProto(ctx, elem)
		diags.Append(dvDiags...)
		if diags.HasError() {
			return nil, diags
		}
		proto.Rules = append(proto.Rules, pvModel)
	}
	pvElemModelRulesByName := pv.RulesByName.Elements()
	proto.RulesByName = make(map[string]*configv1.PolicyVersion_Spec_RulesByName, len(pvElemModelRulesByName))
	for key, elem := range pvElemModelRulesByName {
		pvModel, dvDiags := ConvertDataValueToPolicyVersion_Spec_RulesByNameProto(ctx, elem)
		diags.Append(dvDiags...)
		if diags.HasError() {
			return nil, diags
		}
		proto.RulesByName[key] = pvModel
	}
	return proto, diags
}
`
	suite.Contains(output, expectedConvertDataValue)
}

// TestNestedCollectionsOfObjects verifies that nested collections (Set of Sets, List of Lists)
// containing objects generate correct proto type names with proper nesting depth.
// For Set of Sets of Objects, the naming uses "_elem" suffix to distinguish nesting levels:
//
//	Resource_Items (outer wrapper) -> Resource_ItemsElem (inner object)
func (suite *GenerateProviderTestSuite) TestNestedCollectionsOfObjects() {
	// Schema: Set of Sets of Objects
	testResource := schema.Resource{
		TypeName: "nested_collection_test",
		Schema: resource_schema.Schema{
			Version:     1,
			Description: "Test resource for nested collections of objects.",
			Attributes: map[string]resource_schema.Attribute{
				"id": resource_schema.StringAttribute{
					Description: "Resource ID.",
					Computed:    true,
				},
				// Set of Sets of Objects
				"items": resource_schema.SetAttribute{
					Description: "Nested set of sets of objects.",
					Required:    true,
					ElementType: types.SetType{
						ElemType: types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"name":  types.StringType,
								"value": types.Int64Type,
							},
						},
					},
				},
				// List of Lists of Objects (similar pattern)
				"records": resource_schema.ListAttribute{
					Description: "Nested list of lists of objects.",
					Optional:    true,
					ElementType: types.ListType{
						ElemType: types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"key": types.StringType,
							},
						},
					},
				},
			},
		},
	}

	data := providerTemplateData{
		Package:               "testpkg",
		ProviderTypeName:      "Provider",
		Models:                make([]model, 0),
		NewRequestFuncs:       make([]convertFunc, 0, 3),
		NewUpdateRequestFuncs: make([]convertFunc, 0, 1),
		CopyResponseFuncs:     make([]convertFunc, 0, 3),
		Resources:             make([]resourceData, 0, 1),
	}

	err := AddResourceToProviderTemplateData(&testResource, &data, "NestedCollectionTest", "NestedCollectionTest")
	suite.Require().NoError(err, "AddResourceToProviderTemplateData should not return an error")

	// Verify the model was created
	suite.Require().Len(data.Models, 1, "Should have one resource model")
	resourceModel := data.Models[0]

	// Find the "items" field (Set of Sets of Objects)
	var itemsField *field

	for i := range resourceModel.Fields {
		if resourceModel.Fields[i].AttributeName == "items" {
			itemsField = &resourceModel.Fields[i]

			break
		}
	}

	suite.Require().NotNil(itemsField, "Should have 'items' field")

	// For Set of Sets of Objects (attr "items"):
	// - Outer wrapper type: []*configv1.NestedCollectionTest_Items
	// - Object model: NestedCollectionTest_ItemsElem (uses _elem suffix for inner element)
	//
	// The _elem suffix distinguishes nested collection elements from the outer wrapper,
	// matching the naming convention used by protogen.

	// Check outer collection type name
	suite.Equal("[]*configv1.NestedCollectionTest_Items", itemsField.Type.ProtoTypeName,
		"Outer Set should have correct wrapper type name")

	// Check inner collection (first level of CollectionElementType)
	suite.Require().NotNil(itemsField.Type.CollectionElementType, "Should have CollectionElementType for outer Set")
	innerSetType := itemsField.Type.CollectionElementType
	suite.Equal("[]*configv1.NestedCollectionTest_ItemsElem", innerSetType.ProtoTypeName,
		"Inner Set should have _Elem suffix to distinguish from outer wrapper")

	// Check object (inside inner CollectionElementType)
	suite.Require().NotNil(innerSetType.CollectionElementType, "Should have CollectionElementType for inner Set")
	objectType := innerSetType.CollectionElementType
	suite.Require().NotNil(objectType.NestedModel, "Inner Set element should have NestedModel for the Object")
	suite.Equal("NestedCollectionTest_ItemsElem", objectType.NestedModel.Name,
		"Object should use _Elem suffix matching inner collection element name")

	// Find the "records" field (List of Lists of Objects)
	var recordsField *field

	for i := range resourceModel.Fields {
		if resourceModel.Fields[i].AttributeName == "records" {
			recordsField = &resourceModel.Fields[i]

			break
		}
	}

	suite.Require().NotNil(recordsField, "Should have 'records' field")

	// Check outer List type name
	suite.Equal("[]*configv1.NestedCollectionTest_Records", recordsField.Type.ProtoTypeName,
		"Outer List should have correct wrapper type name")

	// Check inner List
	suite.Require().NotNil(recordsField.Type.CollectionElementType, "Should have CollectionElementType for outer List")
	innerListType := recordsField.Type.CollectionElementType
	suite.Equal("[]*configv1.NestedCollectionTest_RecordsElem", innerListType.ProtoTypeName,
		"Inner List should have _Elem suffix")

	// Check object in inner List
	suite.Require().NotNil(innerListType.CollectionElementType, "Should have CollectionElementType for inner List")
	recordObjectType := innerListType.CollectionElementType
	suite.Require().NotNil(recordObjectType.NestedModel, "Inner List element should have NestedModel")
	suite.Equal("NestedCollectionTest_RecordsElem", recordObjectType.NestedModel.Name,
		"Object should use _Elem suffix")
}

// TestListPrimitiveAttribute verifies that ListAttribute with primitive element types
// (e.g., list of strings) generates correct converter code.
func (suite *GenerateProviderTestSuite) TestListPrimitiveAttribute() {
	testResource := schema.Resource{
		TypeName: "test_resource",
		Schema: resource_schema.Schema{
			Version:     1,
			Description: "Test resource.",
			Attributes: map[string]resource_schema.Attribute{
				"id": resource_schema.StringAttribute{
					Description: "Resource ID.",
					Computed:    true,
				},
				// Outer object containing a list of primitives
				"config": resource_schema.SingleNestedAttribute{
					Description: "Configuration.",
					Required:    true,
					Attributes: map[string]resource_schema.Attribute{
						"name": resource_schema.StringAttribute{
							Description: "Config name.",
							Required:    true,
						},
						// List of strings inside the nested object
						"tags": resource_schema.ListAttribute{
							Description: "List of tags.",
							Required:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}

	data := providerTemplateData{
		Package:               "testpkg",
		ProviderTypeName:      "Provider",
		Models:                make([]model, 0),
		NewRequestFuncs:       make([]convertFunc, 0, 3),
		NewUpdateRequestFuncs: make([]convertFunc, 0, 1),
		CopyResponseFuncs:     make([]convertFunc, 0, 3),
		Resources:             make([]resourceData, 0, 1),
	}

	err := AddResourceToProviderTemplateData(&testResource, &data, "TestResource", "TestResource")
	suite.Require().NoError(err, "AddResourceToProviderTemplateData should not return an error")

	// Generate the converters code
	dst := new(bytes.Buffer)
	err = ProviderConvertersTemplate.Execute(dst, &data)
	suite.Require().NoError(err, "ProviderConvertersTemplate.Execute should not return an error")

	output := dst.String()

	// Verify GetTypeAttrsFor generates correct List type with primitive element type
	expectedGetTypeAttrs := `
func GetTypeAttrsForTestResource_Config() map[string]attr.Type {
	return map[string]attr.Type{
		"name": types.StringType,
		"tags": types.ListType{ElemType: types.StringType},
	}
}
`
	suite.Contains(output, expectedGetTypeAttrs)

	// Verify ConvertToObjectValueFromProto iterates and converts each primitive
	expectedToProto := `
func ConvertTestResource_ConfigToObjectValueFromProto(proto *configv1.TestResource_Config) basetypes.ObjectValue  {
	if proto == nil {
		return types.ObjectNull(GetTypeAttrsForTestResource_Config())
	}
	elementsInTags := make([]attr.Value, 0, len(proto.Tags))
	for _, item := range proto.Tags {
		elementsInTags = append(elementsInTags, types.StringValue(item))
	}
	return types.ObjectValueMust(
		GetTypeAttrsForTestResource_Config(),
		map[string]attr.Value{
			"name": types.StringValue(proto.Name),
			"tags": func() basetypes.ListValue { if proto.Tags == nil { return types.ListNull(types.StringType) }; return types.ListValueMust(types.StringType, elementsInTags) }(),
		},
	)
}
`
	suite.Contains(output, expectedToProto)

	// Verify ConvertDataValueToProto uses ElementsAs for primitives
	expectedFromProto := `
func ConvertDataValueToTestResource_ConfigProto(ctx context.Context, dataValue attr.Value) (*configv1.TestResource_Config, diag.Diagnostics) {
	if dataValue.IsNull() || dataValue.IsUnknown() {
		return nil, nil
	}
	pv := TestResource_Config{}
	diags := tfsdk.ValueAs(ctx, dataValue, &pv)
	if diags.HasError() {
		return nil, diags
	}
	proto := &configv1.TestResource_Config{}
	proto.Name = pv.Name.ValueString()
	var pvModelTags []string
	dvDiagsTags := pv.Tags.ElementsAs(ctx, &pvModelTags, false)
	diags.Append(dvDiagsTags...)
	if diags.HasError() {
		return nil, diags
	}
	proto.Tags = pvModelTags
	return proto, diags
}
`
	suite.Contains(output, expectedFromProto)
}

// TestMapPrimitiveAttribute verifies that MapAttribute with primitive element types
// (e.g., map of strings) generates correct converter code.
func (suite *GenerateProviderTestSuite) TestMapPrimitiveAttribute() {
	testResource := schema.Resource{
		TypeName: "test_map_resource",
		Schema: resource_schema.Schema{
			Version:     1,
			Description: "Test resource with map attribute.",
			Attributes: map[string]resource_schema.Attribute{
				"id": resource_schema.StringAttribute{
					Description: "Resource ID.",
					Computed:    true,
				},
				// Outer object containing a map of strings
				"metadata": resource_schema.SingleNestedAttribute{
					Description: "Metadata configuration.",
					Required:    true,
					Attributes: map[string]resource_schema.Attribute{
						"name": resource_schema.StringAttribute{
							Description: "Metadata name.",
							Required:    true,
						},
						// Map of string to string
						"labels": resource_schema.MapAttribute{
							Description: "Key-value labels.",
							Required:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}

	data := providerTemplateData{
		Package:               "testpkg",
		ProviderTypeName:      "Provider",
		Models:                make([]model, 0),
		NewRequestFuncs:       make([]convertFunc, 0, 3),
		NewUpdateRequestFuncs: make([]convertFunc, 0, 1),
		CopyResponseFuncs:     make([]convertFunc, 0, 3),
		Resources:             make([]resourceData, 0, 1),
	}

	err := AddResourceToProviderTemplateData(&testResource, &data, "TestMapResource", "TestMapResource")
	suite.Require().NoError(err, "AddResourceToProviderTemplateData should not return an error")

	// Generate the converters code
	dst := new(bytes.Buffer)
	err = ProviderConvertersTemplate.Execute(dst, &data)
	suite.Require().NoError(err, "ProviderConvertersTemplate.Execute should not return an error")

	output := dst.String()

	// Verify GetTypeAttrsFor generates correct Map type with primitive element type
	expectedGetTypeAttrs := `
func GetTypeAttrsForTestMapResource_Metadata() map[string]attr.Type {
	return map[string]attr.Type{
		"labels": types.MapType{ElemType: types.StringType},
		"name": types.StringType,
	}
}
`
	suite.Contains(output, expectedGetTypeAttrs)

	// Verify ConvertToObjectValueFromProto iterates map and converts each value
	expectedToProto := `
func ConvertTestMapResource_MetadataToObjectValueFromProto(proto *configv1.TestMapResource_Metadata) basetypes.ObjectValue  {
	if proto == nil {
		return types.ObjectNull(GetTypeAttrsForTestMapResource_Metadata())
	}
	elementsInLabels := make(map[string]attr.Value, len(proto.Labels))
	for key, item := range proto.Labels {
		elementsInLabels[key] = types.StringValue(item)
	}
	return types.ObjectValueMust(
		GetTypeAttrsForTestMapResource_Metadata(),
		map[string]attr.Value{
			"labels": func() basetypes.MapValue { if proto.Labels == nil { return types.MapNull(types.StringType) }; return types.MapValueMust(types.StringType, elementsInLabels) }(),
			"name": types.StringValue(proto.Name),
		},
	)
}
`
	suite.Contains(output, expectedToProto)

	// Verify ConvertDataValueToProto uses ElementsAs for map
	expectedFromProto := `
func ConvertDataValueToTestMapResource_MetadataProto(ctx context.Context, dataValue attr.Value) (*configv1.TestMapResource_Metadata, diag.Diagnostics) {
	if dataValue.IsNull() || dataValue.IsUnknown() {
		return nil, nil
	}
	pv := TestMapResource_Metadata{}
	diags := tfsdk.ValueAs(ctx, dataValue, &pv)
	if diags.HasError() {
		return nil, diags
	}
	proto := &configv1.TestMapResource_Metadata{}
	var pvModelLabels map[string]string
	dvDiagsLabels := pv.Labels.ElementsAs(ctx, &pvModelLabels, false)
	diags.Append(dvDiagsLabels...)
	if diags.HasError() {
		return nil, diags
	}
	proto.Labels = pvModelLabels
	proto.Name = pv.Name.ValueString()
	return proto, diags
}
`
	suite.Contains(output, expectedFromProto)
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

// TestNilProtoHandling verifies the three null/nil handling fixes in generated converters:
// 1. Convert...ToObjectValueFromProto returns ObjectNull when proto is nil
// 2. ConvertDataValueTo...Proto returns nil when dataValue is null/unknown
// 3. Collection fields return typed null (not empty) when the proto field is nil
func (suite *GenerateProviderTestSuite) TestNilProtoHandling() {
	testResource := schema.Resource{
		TypeName: "nil_test",
		Schema: resource_schema.Schema{
			Version:     1,
			Description: "Test resource for nil/null handling.",
			Attributes: map[string]resource_schema.Attribute{
				"id": resource_schema.StringAttribute{
					Description: "Resource ID.",
					Computed:    true,
				},
				"nested": resource_schema.SingleNestedAttribute{
					Description: "Nested object with optional collections.",
					Required:    true,
					Attributes: map[string]resource_schema.Attribute{
						"name": resource_schema.StringAttribute{
							Description: "Name.",
							Required:    true,
						},
						"tags": resource_schema.ListAttribute{
							Description: "Optional list of strings.",
							Optional:    true,
							ElementType: types.StringType,
						},
						"items": resource_schema.ListNestedAttribute{
							Description: "Optional list of nested objects.",
							Optional:    true,
							NestedObject: resource_schema.NestedAttributeObject{
								Attributes: map[string]resource_schema.Attribute{
									"key": resource_schema.StringAttribute{
										Description: "Key.",
										Required:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	data := providerTemplateData{
		Package:               "testpkg",
		ProviderTypeName:      "Provider",
		Models:                make([]model, 0),
		NewRequestFuncs:       make([]convertFunc, 0, 3),
		NewUpdateRequestFuncs: make([]convertFunc, 0, 1),
		CopyResponseFuncs:     make([]convertFunc, 0, 3),
		Resources:             make([]resourceData, 0, 1),
	}

	err := AddResourceToProviderTemplateData(&testResource, &data, "NilTest", "NilTest")
	suite.Require().NoError(err)

	dst := new(bytes.Buffer)
	err = ProviderConvertersTemplate.Execute(dst, &data)
	suite.Require().NoError(err)

	output := dst.String()

	// Fix 1: Convert...ToObjectValueFromProto returns ObjectNull for nil proto
	suite.Contains(output, `func ConvertNilTest_NestedToObjectValueFromProto(proto *configv1.NilTest_Nested) basetypes.ObjectValue  {
	if proto == nil {
		return types.ObjectNull(GetTypeAttrsForNilTest_Nested())
	}`,
		"should return ObjectNull when proto is nil")

	// Fix 2: ConvertDataValueTo...Proto returns nil for null/unknown data values
	suite.Contains(output, `func ConvertDataValueToNilTest_NestedProto(ctx context.Context, dataValue attr.Value) (*configv1.NilTest_Nested, diag.Diagnostics) {
	if dataValue.IsNull() || dataValue.IsUnknown() {
		return nil, nil
	}`,
		"should return nil when dataValue is null or unknown")

	// Fix 3a: Primitive collection (tags) returns typed null when proto field is nil
	suite.Contains(output,
		`"tags": func() basetypes.ListValue { if proto.Tags == nil { return types.ListNull(types.StringType) }; return types.ListValueMust(types.StringType, elementsInTags) }()`,
		"primitive list should return ListNull when proto field is nil")

	// Fix 3b: Nested object collection (items) returns typed null when proto field is nil
	suite.Contains(output,
		`"items": func() basetypes.ListValue { if proto.Items == nil { return types.ListNull(types.ObjectType{AttrTypes: GetTypeAttrsForNilTest_Nested_Items()}) }; return types.ListValueMust(types.ObjectType{AttrTypes: GetTypeAttrsForNilTest_Nested_Items()}, elementsInItems) }()`,
		"nested object list should return ListNull when proto field is nil")
}
