// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"fmt"

	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// AttributeIsOptional returns whether an attribute should be optional in gRPC messages, i.e., whether it is optional with no default value.
func AttributeIsOptional(attribute any) bool { //nolint: gocyclo
	if !attribute.(resource_schema.Attribute).IsOptional() { //nolint: forcetypeassert
		return false
	}

	// Only consider a resource attribute to be optional if Optional is true
	// AND it has no default value.
	switch a := attribute.(type) {
	case datasource_schema.BoolAttribute:
		return true
	case datasource_schema.Float32Attribute:
		return true
	case datasource_schema.Float64Attribute:
		return true
	case datasource_schema.Int32Attribute:
		return true
	case datasource_schema.Int64Attribute:
		return true
	case datasource_schema.ListAttribute:
		return true
	case datasource_schema.MapAttribute:
		return true
	case datasource_schema.NumberAttribute:
		return true
	case datasource_schema.ObjectAttribute:
		return true
	case datasource_schema.SetAttribute:
		return true
	case datasource_schema.StringAttribute:
		return true
	case datasource_schema.DynamicAttribute:
		return true
	case resource_schema.BoolAttribute:
		return a.Default == nil
	case resource_schema.Float32Attribute:
		return a.Default == nil
	case resource_schema.Float64Attribute:
		return a.Default == nil
	case resource_schema.Int32Attribute:
		return a.Default == nil
	case resource_schema.Int64Attribute:
		return a.Default == nil
	case resource_schema.ListAttribute:
		return a.Default == nil
	case resource_schema.MapAttribute:
		return a.Default == nil
	case resource_schema.NumberAttribute:
		return a.Default == nil
	case resource_schema.ObjectAttribute:
		return a.Default == nil
	case resource_schema.SetAttribute:
		return a.Default == nil
	case resource_schema.StringAttribute:
		return a.Default == nil
	case resource_schema.DynamicAttribute:
		return a.Default == nil
	case resource_schema.ListNestedAttribute:
		return a.Default == nil
	case resource_schema.MapNestedAttribute:
		return a.Default == nil
	case resource_schema.SingleNestedAttribute:
		return a.Default == nil
	case resource_schema.SetNestedAttribute:
		return a.Default == nil
	case BoolResourceAttributeWithMode:
		return a.Default == nil
	case Float64ResourceAttributeWithMode:
		return a.Default == nil
	case Int64ResourceAttributeWithMode:
		return a.Default == nil
	case ListResourceAttributeWithMode:
		return a.Default == nil
	case MapResourceAttributeWithMode:
		return a.Default == nil
	case NumberResourceAttributeWithMode:
		return a.Default == nil
	case ObjectResourceAttributeWithMode:
		return a.Default == nil
	case SetResourceAttributeWithMode:
		return a.Default == nil
	case StringResourceAttributeWithMode:
		return a.Default == nil
	default:
		panic(fmt.Sprintf("unknown attribute type: %T", attribute))
	}
}

// NestedAttributes returns the attributes of the object nested in attribute, keyed by name, or nil
// if attribute has no nested schema. Attributes declared with an object type rather than nested
// attributes, e.g. a ListAttribute of an ObjectType, have no nested schema.
func NestedAttributes(attribute any) map[string]any {
	nested, ok := attribute.(resource_schema.NestedAttribute)
	if !ok {
		return nil
	}

	attributes := nested.GetNestedObject().GetAttributes()
	out := make(map[string]any, len(attributes))

	for name, a := range attributes {
		out[name] = a
	}

	return out
}

// NestedScalarIsOptional returns whether a scalar field of a nested object should be optional in
// gRPC messages, so that an unset value round-trips as null rather than as its zero value. attrSchema
// is the field's schema, or nil if the object has no nested schema.
func NestedScalarIsOptional(attrSchema any) bool {
	a, ok := attrSchema.(resource_schema.Attribute)
	if !ok {
		return false
	}

	switch a.GetType().(type) {
	case basetypes.BoolType, basetypes.Float64Type, basetypes.Int64Type, basetypes.StringType:
		return AttributeIsOptional(a)
	default:
		return false
	}
}
