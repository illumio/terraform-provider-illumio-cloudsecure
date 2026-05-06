// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"sort"

	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func (suite *SchemaTestSuite) TestResourcesAreSorted() {
	suite.True(sort.IsSorted(suite.Schema.Resources()), "Resources should be sorted")
}

func (suite *SchemaTestSuite) TestResourcesHaveUniqueTypeNames() {
	resources := suite.Schema.Resources()
	typeNames := make(map[string]struct{}, len(resources))

	for i, resource := range resources {
		typeName := resource.TypeName
		_, found := typeNames[typeName]

		suite.False(found, "Resource type name %q is duplicated at index %d", typeName, i)

		typeNames[typeName] = struct{}{}
	}
}

func (suite *SchemaTestSuite) TestEachResourceHasIdAttribute() {
	for i, resource := range suite.Schema.Resources() {
		attributes := resource.Schema.Attributes

		var idAttributeFound bool

		for name, attribute := range attributes {
			if name == IDFieldName {
				idAttributeFound = true

				suite.Equal(idAttribute, attribute, "Attribute %q in resource %q at index %d should have the expected schema", IDFieldName, resource.TypeName, i)
			}
		}

		suite.True(idAttributeFound, "Resource %q at index %d should have an %q attribute", resource.TypeName, i, IDFieldName)
	}
}

func (suite *SchemaTestSuite) TestNoResourceHasUpdateMaskAttribute() {
	for i, resource := range suite.Schema.Resources() {
		attributes := resource.Schema.Attributes
		for name := range attributes {
			suite.NotEqual(UpdateMaskFieldName, name, "Resource %q at index %d should not define attribute with name %q", resource.TypeName, i, UpdateMaskFieldName)
		}
	}
}

func (suite *SchemaTestSuite) TestDataSourcesAreSorted() {
	suite.True(sort.IsSorted(suite.Schema.DataSources()), "Data sources should be sorted")
}

func (suite *SchemaTestSuite) TestDataSourcesHaveUniqueTypeNames() {
	dataSources := suite.Schema.DataSources()
	typeNames := make(map[string]struct{}, len(dataSources))

	for i, dataSource := range dataSources {
		typeName := dataSource.TypeName
		_, found := typeNames[typeName]

		suite.False(found, "Data source type name %q is duplicated at index %d", typeName, i)

		typeNames[typeName] = struct{}{}
	}
}

func (suite *SchemaTestSuite) TestGetResourceAttributeMode() {
	tests := []struct {
		name         string
		attr         resource_schema.Attribute
		expectedMode AttributeMode
	}{
		{
			name: "ListNestedAttribute defaults to ReadWrite",
			attr: resource_schema.ListNestedAttribute{
				Required: true,
			},
			expectedMode: ReadWriteAttributeMode,
		},
		{
			name: "ListNestedResourceAttributeWithMode returns configured mode",
			attr: ListNestedResourceAttributeWithMode{
				ListNestedAttribute: resource_schema.ListNestedAttribute{
					Required: true,
				},
				attributeWithMode: attributeWithMode{
					Mode: ImmutableAttributeMode,
				},
			},
			expectedMode: ImmutableAttributeMode,
		},
		{
			name: "StringAttribute defaults to ReadWrite",
			attr: resource_schema.StringAttribute{
				Required: true,
			},
			expectedMode: ReadWriteAttributeMode,
		},
		{
			name: "StringResourceAttributeWithMode returns configured mode",
			attr: StringResourceAttributeWithMode{
				StringAttribute: resource_schema.StringAttribute{
					Required: true,
				},
				attributeWithMode: attributeWithMode{
					Mode: ImmutableAttributeMode,
				},
			},
			expectedMode: ImmutableAttributeMode,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			mode := GetResourceAttributeMode(tc.attr)
			suite.Equal(tc.expectedMode, mode)
		})
	}
}
