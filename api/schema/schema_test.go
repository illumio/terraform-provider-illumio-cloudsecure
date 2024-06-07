// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"sort"
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
