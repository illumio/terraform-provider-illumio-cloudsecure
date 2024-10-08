// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func (suite *SchemaTestSuite) TestProtoMessageName() {
	tests := map[string]struct {
		tfName string
		output string
	}{
		"aws_account": {
			tfName: "aws_account",
			output: "AwsAccount",
		},
		"id": {
			tfName: "id",
			output: "Id",
		},
		"account_id": {
			tfName: "account_id",
			output: "AccountId",
		},
		"service_account_id": {
			tfName: "service_account_id",
			output: "ServiceAccountId",
		},
		"k8s_cluster": { // gRPC considers a digit to be a delimiter, so "k8s" is capitalized as "K8S".
			tfName: "k8s_cluster",
			output: "K8SCluster",
		},
	}

	for name, tc := range tests {
		suite.Run(name, func() {
			got := ProtoMessageName(tc.tfName)
			suite.Equal(tc.output, got, "Converted name should match")
		})
	}
}

func (suite *SchemaTestSuite) TestSortResourceAttributes() {
	tests := map[string]struct {
		attrs  map[string]resource_schema.Attribute
		output []string
	}{
		"id-only": {
			attrs: map[string]resource_schema.Attribute{
				IDFieldName: resource_schema.StringAttribute{},
			},
			output: []string{IDFieldName},
		},
		"3-attributes": {
			attrs: map[string]resource_schema.Attribute{
				IDFieldName: resource_schema.StringAttribute{},
				"a":         resource_schema.StringAttribute{},
				"z":         resource_schema.StringAttribute{},
			},
			output: []string{IDFieldName, "a", "z"},
		},
	}

	for name, tc := range tests {
		suite.Run(name, func() {
			got := SortResourceAttributes(tc.attrs)
			suite.Equal(tc.output, got, "Attributes should be sorted")
		})
	}
}

func (suite *SchemaTestSuite) TestSortDataSourceAttributes() {
	tests := map[string]struct {
		attrs  map[string]datasource_schema.Attribute
		output []string
	}{
		"id-only": {
			attrs: map[string]datasource_schema.Attribute{
				IDFieldName: datasource_schema.StringAttribute{},
			},
			output: []string{IDFieldName},
		},
		"3-attributes": {
			attrs: map[string]datasource_schema.Attribute{
				IDFieldName: datasource_schema.StringAttribute{},
				"a":         datasource_schema.StringAttribute{},
				"z":         datasource_schema.StringAttribute{},
			},
			output: []string{IDFieldName, "a", "z"},
		},
	}

	for name, tc := range tests {
		suite.Run(name, func() {
			got := SortDataSourceAttributes(tc.attrs)
			suite.Equal(tc.output, got, "Attributes should be sorted")
		})
	}
}
