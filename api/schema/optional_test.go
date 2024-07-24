// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/dynamicdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/numberdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func (suite *SchemaTestSuite) TestAttributeIsOptional() {
	tests := map[string]struct {
		attribute any
		optional  bool
	}{
		"datasource_bool_optional": {
			attribute: datasource_schema.BoolAttribute{Optional: true},
			optional:  true,
		},
		"datasource_bool_nonoptional": {
			attribute: datasource_schema.BoolAttribute{Optional: false},
			optional:  false,
		},
		"datasource_float32_optional": {
			attribute: datasource_schema.Float32Attribute{Optional: true},
			optional:  true,
		},
		"datasource_float32_nonoptional": {
			attribute: datasource_schema.Float32Attribute{Optional: false},
			optional:  false,
		},
		"datasource_float64_optional": {
			attribute: datasource_schema.Float64Attribute{Optional: true},
			optional:  true,
		},
		"datasource_float64_nonoptional": {
			attribute: datasource_schema.Float64Attribute{Optional: false},
			optional:  false,
		},
		"datasource_int32_optional": {
			attribute: datasource_schema.Int32Attribute{Optional: true},
			optional:  true,
		},
		"datasource_int32_nonoptional": {
			attribute: datasource_schema.Int32Attribute{Optional: false},
			optional:  false,
		},
		"datasource_int64_optional": {
			attribute: datasource_schema.Int64Attribute{Optional: true},
			optional:  true,
		},
		"datasource_int64_nonoptional": {
			attribute: datasource_schema.Int64Attribute{Optional: false},
			optional:  false,
		},
		"datasource_list_optional": {
			attribute: datasource_schema.ListAttribute{Optional: true},
			optional:  true,
		},
		"datasource_list_nonoptional": {
			attribute: datasource_schema.ListAttribute{Optional: false},
			optional:  false,
		},
		"datasource_map_optional": {
			attribute: datasource_schema.MapAttribute{Optional: true},
			optional:  true,
		},
		"datasource_map_nonoptional": {
			attribute: datasource_schema.MapAttribute{Optional: false},
			optional:  false,
		},
		"datasource_number_optional": {
			attribute: datasource_schema.NumberAttribute{Optional: true},
			optional:  true,
		},
		"datasource_number_nonoptional": {
			attribute: datasource_schema.NumberAttribute{Optional: false},
			optional:  false,
		},
		"datasource_object_optional": {
			attribute: datasource_schema.ObjectAttribute{Optional: true},
			optional:  true,
		},
		"datasource_object_nonoptional": {
			attribute: datasource_schema.ObjectAttribute{Optional: false},
			optional:  false,
		},
		"datasource_set_optional": {
			attribute: datasource_schema.SetAttribute{Optional: true},
			optional:  true,
		},
		"datasource_set_nonoptional": {
			attribute: datasource_schema.SetAttribute{Optional: false},
			optional:  false,
		},
		"datasource_string_optional": {
			attribute: datasource_schema.StringAttribute{Optional: true},
			optional:  true,
		},
		"datasource_string_nonoptional": {
			attribute: datasource_schema.StringAttribute{Optional: false},
			optional:  false,
		},
		"datasource_dynamic_optional": {
			attribute: datasource_schema.DynamicAttribute{Optional: true},
			optional:  true,
		},
		"datasource_dynamic_nonoptional": {
			attribute: datasource_schema.DynamicAttribute{Optional: false},
			optional:  false,
		},
		"resource_bool_optional": {
			attribute: resource_schema.BoolAttribute{Optional: true},
			optional:  true,
		},
		"resource_bool_optional_default": {
			attribute: resource_schema.BoolAttribute{Optional: true, Default: booldefault.StaticBool(false)},
			optional:  false,
		},
		"resource_bool_nonoptional": {
			attribute: resource_schema.BoolAttribute{Optional: false},
			optional:  false,
		},
		"resource_float32_optional": {
			attribute: resource_schema.Float32Attribute{Optional: true},
			optional:  true,
		},
		"resource_float32_optional_default": {
			attribute: resource_schema.Float32Attribute{Optional: true, Default: float32default.StaticFloat32(0.0)},
			optional:  false,
		},
		"resource_float32_nonoptional": {
			attribute: resource_schema.Float32Attribute{Optional: false},
			optional:  false,
		},
		"resource_float64_optional": {
			attribute: resource_schema.Float64Attribute{Optional: true},
			optional:  true,
		},
		"resource_float64_optional_default": {
			attribute: resource_schema.Float64Attribute{Optional: true, Default: float64default.StaticFloat64(0.0)},
			optional:  false,
		},
		"resource_float64_nonoptional": {
			attribute: resource_schema.Float64Attribute{Optional: false},
			optional:  false,
		},
		"resource_int32_optional": {
			attribute: resource_schema.Int32Attribute{Optional: true},
			optional:  true,
		},
		"resource_int32_optional_default": {
			attribute: resource_schema.Int32Attribute{Optional: true, Default: int32default.StaticInt32(0)},
			optional:  false,
		},
		"resource_int32_nonoptional": {
			attribute: resource_schema.Int32Attribute{Optional: false},
			optional:  false,
		},
		"resource_int64_optional": {
			attribute: resource_schema.Int64Attribute{Optional: true},
			optional:  true,
		},
		"resource_int64_optional_default": {
			attribute: resource_schema.Int64Attribute{Optional: true, Default: int64default.StaticInt64(0)},
			optional:  false,
		},
		"resource_int64_nonoptional": {
			attribute: resource_schema.Int64Attribute{Optional: false},
			optional:  false,
		},
		"resource_list_optional": {
			attribute: resource_schema.ListAttribute{Optional: true},
			optional:  true,
		},
		"resource_list_optional_default": {
			attribute: resource_schema.ListAttribute{Optional: true, Default: listdefault.StaticValue(basetypes.ListValue{})},
			optional:  false,
		},
		"resource_list_nonoptional": {
			attribute: resource_schema.ListAttribute{Optional: false},
			optional:  false,
		},
		"resource_map_optional": {
			attribute: resource_schema.MapAttribute{Optional: true},
			optional:  true,
		},
		"resource_map_optional_default": {
			attribute: resource_schema.MapAttribute{Optional: true, Default: mapdefault.StaticValue(basetypes.MapValue{})},
			optional:  false,
		},
		"resource_map_nonoptional": {
			attribute: resource_schema.MapAttribute{Optional: false},
			optional:  false,
		},
		"resource_number_optional": {
			attribute: resource_schema.NumberAttribute{Optional: true},
			optional:  true,
		},
		"resource_number_optional_default": {
			attribute: resource_schema.NumberAttribute{Optional: true, Default: numberdefault.StaticBigFloat(nil)},
			optional:  false,
		},
		"resource_number_nonoptional": {
			attribute: resource_schema.NumberAttribute{Optional: false},
			optional:  false,
		},
		"resource_object_optional": {
			attribute: resource_schema.ObjectAttribute{Optional: true},
			optional:  true,
		},
		"resource_object_optional_default": {
			attribute: resource_schema.ObjectAttribute{Optional: true, Default: objectdefault.StaticValue(basetypes.ObjectValue{})},
			optional:  false,
		},
		"resource_object_nonoptional": {
			attribute: resource_schema.ObjectAttribute{Optional: false},
			optional:  false,
		},
		"resource_set_optional": {
			attribute: resource_schema.SetAttribute{Optional: true},
			optional:  true,
		},
		"resource_set_optional_default": {
			attribute: resource_schema.SetAttribute{Optional: true, Default: setdefault.StaticValue(basetypes.SetValue{})},
			optional:  false,
		},
		"resource_set_nonoptional": {
			attribute: resource_schema.SetAttribute{Optional: false},
			optional:  false,
		},
		"resource_string_optional": {
			attribute: resource_schema.StringAttribute{Optional: true},
			optional:  true,
		},
		"resource_string_optional_default": {
			attribute: resource_schema.StringAttribute{Optional: true, Default: stringdefault.StaticString("")},
			optional:  false,
		},
		"resource_string_nonoptional": {
			attribute: resource_schema.StringAttribute{Optional: false},
			optional:  false,
		},
		"resource_dynamic_optional": {
			attribute: resource_schema.DynamicAttribute{Optional: true},
			optional:  true,
		},
		"resource_dynamic_optional_default": {
			attribute: resource_schema.DynamicAttribute{Optional: true, Default: dynamicdefault.StaticValue(basetypes.DynamicValue{})},
			optional:  false,
		},
		"resource_dynamic_nonoptional": {
			attribute: resource_schema.DynamicAttribute{Optional: false},
			optional:  false,
		},
	}

	for name, tc := range tests {
		suite.Run(name, func() {
			gotOptional := AttributeIsOptional(tc.attribute)
			suite.Equal(tc.optional, gotOptional, "Optional boolean should match")
		})
	}
}
