package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	configv1 "github.com/illumio/terraform-provider-illumio-cloudsecure/api/illumio/cloud/config/v1"

	"github.com/stretchr/testify/suite"
)

type GenerateProviderSuite struct {
	suite.Suite
}

func TestGenerateProviderSuite(t *testing.T) {
	suite.Run(t, new(GenerateProviderSuite))
}

func (suite *GenerateProviderSuite) TestGenerateProvider() {
	data := AwsTagToLabelResourceModel{
		Id:  types.StringValue("id"),
		Key: types.StringValue("key"),
		Icon: types.ObjectValueMust(
			map[string]attr.Type{
				"name":             types.StringType,
				"background_color": types.StringType,
				"foreground_color": types.StringType,
			},
			map[string]attr.Value{
				"name":             types.StringValue("demo"),
				"background_color": types.StringValue("Beckford Close"),
				"foreground_color": types.StringValue("Gotham"),
			},
		),
		CloudTags: types.ListValueMust(
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"key":   types.StringType,
					"cloud": types.StringType,
				},
			},
			[]attr.Value{
				types.ObjectValueMust(
					map[string]attr.Type{
						"key":      types.StringType,
						"cloud":    types.StringType,
						"cloud_id": types.StringType,
					},
					map[string]attr.Value{
						"key":   types.StringValue("key"),
						"cloud": types.StringValue("cloud"),
					},
				),
			},
		),
		Name: types.StringValue("name"),
	}

	resp := NewCreateAwsTagToLabelRequest(&data)

	suite.NotNil(resp)
}

func (suite *GenerateProviderSuite) TestUpdateProvider() {
	prev := AwsTagToLabelResourceModel{
		Id:  types.StringValue("id"),
		Key: types.StringValue("key"),
		Icon: types.ObjectValueMust(
			map[string]attr.Type{
				"name":             types.StringType,
				"background_color": types.StringType,
				"foreground_color": types.StringType,
			},
			map[string]attr.Value{
				"name":             types.StringValue("demo"),
				"background_color": types.StringValue("Beckford Close"),
				"foreground_color": types.StringValue("Gotham"),
			},
		),
		CloudTags: types.ListValueMust(
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"key":   types.StringType,
					"cloud": types.StringType,
				},
			},
			[]attr.Value{
				types.ObjectValueMust(
					map[string]attr.Type{
						"key":   types.StringType,
						"cloud": types.StringType,
					},
					map[string]attr.Value{
						"key":   types.StringValue("key"),
						"cloud": types.StringValue("cloud"),
					},
				),
			},
		),
		Name: types.StringValue("name"),
	}

	after := AwsTagToLabelResourceModel{
		Id:  types.StringValue("id"),
		Key: types.StringValue("key"),
		Icon: types.ObjectValueMust(
			map[string]attr.Type{
				"name":             types.StringType,
				"background_color": types.StringType,
				"foreground_color": types.StringType,
			},
			map[string]attr.Value{
				"name":             types.StringValue("demo"),
				"background_color": types.StringValue("Beckford Close"),
				"foreground_color": types.StringValue("Gotham"),
			},
		),
		CloudTags: types.ListValueMust(
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"key":   types.StringType,
					"cloud": types.StringType,
				},
			},
			[]attr.Value{
				types.ObjectValueMust(
					map[string]attr.Type{
						"key":   types.StringType,
						"cloud": types.StringType,
					},
					map[string]attr.Value{
						"key":   types.StringValue("key1"),
						"cloud": types.StringValue("cloud1"),
					},
				),
				types.ObjectValueMust(
					map[string]attr.Type{
						"key":   types.StringType,
						"cloud": types.StringType,
					},
					map[string]attr.Value{
						"key":   types.StringValue("key2"),
						"cloud": types.StringValue("cloud2"),
					},
				),
			},
		),
		Name: types.StringValue("name1"),
	}

	resp := NewUpdateAwsTagToLabelRequest(&prev, &after)
	suite.NotNil(resp)
	suite.Require().Len(resp.GetUpdateMask().GetPaths(), 2)
	suite.Contains(resp.GetUpdateMask().GetPaths(), "name")
	suite.Contains(resp.GetUpdateMask().GetPaths(), "cloud_tags")
}

type AwsTagToLabelIcon1 struct {
	Name            string `tfsdk:"name"`
	BackgroundColor string `tfsdk:"background_color"`
	ForegroundColor string `tfsdk:"foreground_color"`
}

func (suite *GenerateProviderSuite) TestGenerateProvider1() {

	rsp := AwsTagToLabelIcon1{}
	obj := types.ObjectValueMust(
		map[string]attr.Type{
			"name":             types.StringType,
			"background_color": types.StringType,
			"foreground_color": types.StringType,
		},
		map[string]attr.Value{
			"name":             types.StringValue("demo"),
			"background_color": types.StringValue("Beckford Close"),
			"foreground_color": types.StringValue("Gotham"),
		},
	)
	diags := tfsdk.ValueAs(context.Background(), obj, &rsp)
	suite.Nil(diags)

	suite.NotNil(rsp)
}

func (suite *GenerateProviderSuite) TestGenerateProvider2() {
	dst := &AwsTagToLabelResourceModel{}
	CopyCreateAwsTagToLabelResponse(dst, &configv1.CreateAwsTagToLabelResponse{
		Id: "id1",
		CloudTags: []*configv1.AwsTagToLabelCloudTagsInstance{
			{
				Key:   "key1",
				Cloud: "cloud1",
			},
		},
		Icon: &configv1.AwsTagToLabelIcon{
			Name:            "name1",
			BackgroundColor: "#000000",
			ForegroundColor: "#FFFFFF",
		},
		Key:  "key1_main",
		Name: "name_main",
	})
	suite.Equal("id1", dst.Id.ValueString())
	suite.Equal("key1_main", dst.Key.ValueString())
	suite.Equal("name_main", dst.Name.ValueString())
}
