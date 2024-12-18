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

type GenerateProviderSuite struct {
	suite.Suite
}

func TestGenerateProviderSuite(t *testing.T) {
	suite.Run(t, new(GenerateProviderSuite))
}

func (suite *GenerateProviderSuite) TestGenerateProvider() {
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

func (suite *GenerateProviderSuite) TestListOfObjects() {
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

func (suite *GenerateProviderSuite) TestSetsAndMore() {
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

	expectedVal := `
type NestedObjectTester_CloudTags struct {
	Cloud types.String ` + "`" + `tfsdk:"cloud"` + "`" + `
	Key types.String ` + "`" + `tfsdk:"key"` + "`" + `
}
func GetTypeAttrsForNestedObjectTester_CloudTags() map[string]attr.Type {
	return map[string]attr.Type{
		"cloud": types.StringType,
		"key": types.StringType,
	}
}


func ConvertNestedObjectTester_CloudTagsToObjectValueFromProto(proto *configv1.NestedObjectTester_CloudTags) basetypes.ObjectValue  {
	return types.ObjectValueMust(
		GetTypeAttrsForNestedObjectTester_CloudTags(),
		map[string]attr.Value{
			"cloud": types.StringValue(proto.Cloud),
			"key": types.StringValue(proto.Key),
		},
	)
}
func ConvertDataValueToNestedObjectTester_CloudTagsProto(ctx context.Context, dataValue attr.Value) *configv1.NestedObjectTester_CloudTags {
	pv := NestedObjectTester_CloudTags{}
	diags := tfsdk.ValueAs(context.Background(), dataValue, &pv)
	if len(diags) > 0 {
		tflog.Error(ctx, "Unexpected diagnostics", map[string]any{"diags": diags})
	}
	proto := &configv1.NestedObjectTester_CloudTags{}
	proto.Cloud = pv.Cloud.ValueString()
	proto.Key = pv.Key.ValueString()
	return proto
}
type NestedObjectTester_Icon struct {
	BackgroundColor types.String ` + "`" + `tfsdk:"background_color"` + "`" + `
	ForegroundColor types.String ` + "`" + `tfsdk:"foreground_color"` + "`" + `
	Name types.String ` + "`" + `tfsdk:"name"` + "`" + `
}
func GetTypeAttrsForNestedObjectTester_Icon() map[string]attr.Type {
	return map[string]attr.Type{
		"background_color": types.StringType,
		"foreground_color": types.StringType,
		"name": types.StringType,
	}
}


func ConvertNestedObjectTester_IconToObjectValueFromProto(proto *configv1.NestedObjectTester_Icon) basetypes.ObjectValue  {
	return types.ObjectValueMust(
		GetTypeAttrsForNestedObjectTester_Icon(),
		map[string]attr.Value{
			"background_color": types.StringValue(proto.BackgroundColor),
			"foreground_color": types.StringValue(proto.ForegroundColor),
			"name": types.StringValue(proto.Name),
		},
	)
}
func ConvertDataValueToNestedObjectTester_IconProto(ctx context.Context, dataValue attr.Value) *configv1.NestedObjectTester_Icon {
	pv := NestedObjectTester_Icon{}
	diags := tfsdk.ValueAs(context.Background(), dataValue, &pv)
	if len(diags) > 0 {
		tflog.Error(ctx, "Unexpected diagnostics", map[string]any{"diags": diags})
	}
	proto := &configv1.NestedObjectTester_Icon{}
	proto.BackgroundColor = pv.BackgroundColor.ValueString()
	proto.ForegroundColor = pv.ForegroundColor.ValueString()
	proto.Name = pv.Name.ValueString()
	return proto
}
type NestedObjectTester_ObjInObj struct {
	Child types.Object ` + "`" + `tfsdk:"child"` + "`" + `
	Name types.String ` + "`" + `tfsdk:"name"` + "`" + `
}
func GetTypeAttrsForNestedObjectTester_ObjInObj() map[string]attr.Type {
	return map[string]attr.Type{
		"child": types.ObjectType{
			AttrTypes: GetTypeAttrsForNestedObjectTester_ObjInObj_Child(),
		},
		"name": types.StringType,
	}
}


func ConvertNestedObjectTester_ObjInObjToObjectValueFromProto(proto *configv1.NestedObjectTester_ObjInObj) basetypes.ObjectValue  {
	return types.ObjectValueMust(
		GetTypeAttrsForNestedObjectTester_ObjInObj(),
		map[string]attr.Value{
			"child": ConvertNestedObjectTester_ObjInObj_ChildToObjectValueFromProto(proto.Child),
			"name": types.StringValue(proto.Name),
		},
	)
}
func ConvertDataValueToNestedObjectTester_ObjInObjProto(ctx context.Context, dataValue attr.Value) *configv1.NestedObjectTester_ObjInObj {
	pv := NestedObjectTester_ObjInObj{}
	diags := tfsdk.ValueAs(context.Background(), dataValue, &pv)
	if len(diags) > 0 {
		tflog.Error(ctx, "Unexpected diagnostics", map[string]any{"diags": diags})
	}
	proto := &configv1.NestedObjectTester_ObjInObj{}
	proto.Child = ConvertDataValueToNestedObjectTester_ObjInObj_ChildProto(ctx, pv.Child)
	proto.Name = pv.Name.ValueString()
	return proto
}
type NestedObjectTester_ObjInObj_Child struct {
	GrandChild types.Object ` + "`" + `tfsdk:"grand_child"` + "`" + `
	Name types.String ` + "`" + `tfsdk:"name"` + "`" + `
}
func GetTypeAttrsForNestedObjectTester_ObjInObj_Child() map[string]attr.Type {
	return map[string]attr.Type{
		"grand_child": types.ObjectType{
			AttrTypes: GetTypeAttrsForNestedObjectTester_ObjInObj_Child_GrandChild(),
		},
		"name": types.StringType,
	}
}


func ConvertNestedObjectTester_ObjInObj_ChildToObjectValueFromProto(proto *configv1.NestedObjectTester_ObjInObj_Child) basetypes.ObjectValue  {
	return types.ObjectValueMust(
		GetTypeAttrsForNestedObjectTester_ObjInObj_Child(),
		map[string]attr.Value{
			"grand_child": ConvertNestedObjectTester_ObjInObj_Child_GrandChildToObjectValueFromProto(proto.GrandChild),
			"name": types.StringValue(proto.Name),
		},
	)
}
func ConvertDataValueToNestedObjectTester_ObjInObj_ChildProto(ctx context.Context, dataValue attr.Value) *configv1.NestedObjectTester_ObjInObj_Child {
	pv := NestedObjectTester_ObjInObj_Child{}
	diags := tfsdk.ValueAs(context.Background(), dataValue, &pv)
	if len(diags) > 0 {
		tflog.Error(ctx, "Unexpected diagnostics", map[string]any{"diags": diags})
	}
	proto := &configv1.NestedObjectTester_ObjInObj_Child{}
	proto.GrandChild = ConvertDataValueToNestedObjectTester_ObjInObj_Child_GrandChildProto(ctx, pv.GrandChild)
	proto.Name = pv.Name.ValueString()
	return proto
}
type NestedObjectTester_ObjInObj_Child_GrandChild struct {
	Name types.String ` + "`" + `tfsdk:"name"` + "`" + `
}
func GetTypeAttrsForNestedObjectTester_ObjInObj_Child_GrandChild() map[string]attr.Type {
	return map[string]attr.Type{
		"name": types.StringType,
	}
}


func ConvertNestedObjectTester_ObjInObj_Child_GrandChildToObjectValueFromProto(proto *configv1.NestedObjectTester_ObjInObj_Child_GrandChild) basetypes.ObjectValue  {
	return types.ObjectValueMust(
		GetTypeAttrsForNestedObjectTester_ObjInObj_Child_GrandChild(),
		map[string]attr.Value{
			"name": types.StringValue(proto.Name),
		},
	)
}
func ConvertDataValueToNestedObjectTester_ObjInObj_Child_GrandChildProto(ctx context.Context, dataValue attr.Value) *configv1.NestedObjectTester_ObjInObj_Child_GrandChild {
	pv := NestedObjectTester_ObjInObj_Child_GrandChild{}
	diags := tfsdk.ValueAs(context.Background(), dataValue, &pv)
	if len(diags) > 0 {
		tflog.Error(ctx, "Unexpected diagnostics", map[string]any{"diags": diags})
	}
	proto := &configv1.NestedObjectTester_ObjInObj_Child_GrandChild{}
	proto.Name = pv.Name.ValueString()
	return proto
}
type NestedObjectTester_SetObj struct {
	SetKey types.String ` + "`" + `tfsdk:"set_key"` + "`" + `
	SetVal types.String ` + "`" + `tfsdk:"set_val"` + "`" + `
}
func GetTypeAttrsForNestedObjectTester_SetObj() map[string]attr.Type {
	return map[string]attr.Type{
		"set_key": types.StringType,
		"set_val": types.StringType,
	}
}


func ConvertNestedObjectTester_SetObjToObjectValueFromProto(proto *configv1.NestedObjectTester_SetObj) basetypes.ObjectValue  {
	return types.ObjectValueMust(
		GetTypeAttrsForNestedObjectTester_SetObj(),
		map[string]attr.Value{
			"set_key": types.StringValue(proto.SetKey),
			"set_val": types.StringValue(proto.SetVal),
		},
	)
}
func ConvertDataValueToNestedObjectTester_SetObjProto(ctx context.Context, dataValue attr.Value) *configv1.NestedObjectTester_SetObj {
	pv := NestedObjectTester_SetObj{}
	diags := tfsdk.ValueAs(context.Background(), dataValue, &pv)
	if len(diags) > 0 {
		tflog.Error(ctx, "Unexpected diagnostics", map[string]any{"diags": diags})
	}
	proto := &configv1.NestedObjectTester_SetObj{}
	proto.SetKey = pv.SetKey.ValueString()
	proto.SetVal = pv.SetVal.ValueString()
	return proto
}
`
	suite.Equal(expectedVal, dst.String(), "Generated provider output should match the expected output")

}
