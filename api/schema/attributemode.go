// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// AttributeMode is the mode of a resource attribute, as the set of RPC messages it should be included in.
type AttributeMode struct {
	// InCreateRequest is true if the attribute should be included in create RPC request messages.
	InCreateRequest bool

	// InCreateResponse is true if the attribute should be included in create RPC response messages.
	InCreateResponse bool

	// InReadRequest is true if the attribute should be included in read RPC request messages.
	InReadRequest bool

	// InCreateRequest is true if the attribute should be included in read RPC response messages.
	InReadResponse bool

	// InUpdateRequest is true if the attribute should be included in update RPC request messages.
	InUpdateRequest bool

	// InUpdateResponse is true if the attribute should be included in update RPC response messages.
	InUpdateResponse bool

	// InDeleteRequest is true if the attribute should be included in delete RPC request messages.
	InDeleteRequest bool
}

var (
	// These variables should be considered constant and must not be modified.

	// KeyAttributeMode is the mode of attributes that are used to identify resource instances.
	KeyAttributeMode = AttributeMode{
		InCreateRequest:  true,
		InCreateResponse: true,
		InReadRequest:    true,
		InReadResponse:   true,
		InUpdateRequest:  true,
		InUpdateResponse: true,
		InDeleteRequest:  true,
	}

	// IDAttributeMode is the mode of "id" attributes.
	IDAttributeMode = AttributeMode{
		InCreateResponse: true,
		InReadRequest:    true,
		InReadResponse:   true,
		InUpdateRequest:  true,
		InUpdateResponse: true,
		InDeleteRequest:  true,
	}

	// ReadWriteAttributeMode is the mode of read-write attributes. This is the default mode.
	ReadWriteAttributeMode = AttributeMode{
		InCreateRequest:  true,
		InCreateResponse: true,
		InReadResponse:   true,
		InUpdateRequest:  true,
		InUpdateResponse: true,
	}

	// ImmutableAttributeMode is the mode of immutable attributes, which can be set only at creation and can be read afterwards.
	// Attributes with this mode should also have the RequiresReplace plan modifier.
	ImmutableAttributeMode = AttributeMode{
		InCreateRequest:  true,
		InCreateResponse: true,
		InReadResponse:   true,
		InUpdateResponse: true,
	}

	// ReadOnlyAttributeMode is the mode of read-only attributes, which are returned by every create, read, and update operation.
	ReadOnlyAttributeMode = AttributeMode{
		InCreateResponse: true,
		InReadResponse:   true,
		InUpdateResponse: true,
	}

	// ReadOnlyOnceAttributeMode is the mode of read-only-once attributes, which are returned only by create operations.
	// Attributes with this mode should also have the UseStateForUnknown plan modifier.
	ReadOnlyOnceAttributeMode = AttributeMode{
		InCreateResponse: true,
	}

	// WriteOnlyAttributeMode is the mode of write-only attributes, which are sent in every create and update operation, and never in any response.
	// Attribute with this mode should also have the UseStateForUnknown plan modifier.
	WriteOnlyAttributeMode = AttributeMode{
		InCreateRequest: true,
		InUpdateRequest: true,
	}

	// WriteOnlyOnceAttributeMode is the mode of write-only-once attributes, which are sent in only in create operations.
	// Attribute with this mode should also have the UseStateForUnknown plan modifier.
	WriteOnlyOnceAttributeMode = AttributeMode{
		InCreateRequest: true,
	}
)

// AttributeWithMode provides the mode of an attribute.
type AttributeWithMode interface {
	// GetMode() returns the mode of the attribute.
	GetMode() AttributeMode
}

// attributeWithMode implements AttributeWithMode.
type attributeWithMode struct {
	// Mode is the mode of the attribute.
	Mode AttributeMode
}

func (a attributeWithMode) GetMode() AttributeMode {
	return a.Mode
}

// BoolResourceAttributeWithMode is a BoolAttribute with an explicit attribute mode.
type BoolResourceAttributeWithMode struct {
	resource_schema.BoolAttribute
	attributeWithMode
}

// Float32ResourceAttributeWithMode is a Float32Attribute with an explicit attribute mode.
type Float32ResourceAttributeWithMode struct {
	resource_schema.Float32Attribute
	attributeWithMode
}

// Float64ResourceAttributeWithMode is a Float64Attribute with an explicit attribute mode.
type Float64ResourceAttributeWithMode struct {
	resource_schema.Float64Attribute
	attributeWithMode
}

// Int32ResourceAttributeWithMode is a Int32Attribute with an explicit attribute mode.
type Int32ResourceAttributeWithMode struct {
	resource_schema.Int32Attribute
	attributeWithMode
}

// Int64ResourceAttributeWithMode is a Int64Attribute with an explicit attribute mode.
type Int64ResourceAttributeWithMode struct {
	resource_schema.Int64Attribute
	attributeWithMode
}

// ListResourceAttributeWithMode is a ListAttribute with an explicit attribute mode.
type ListResourceAttributeWithMode struct {
	resource_schema.ListAttribute
	attributeWithMode
}

// MapResourceAttributeWithMode is a MapAttribute with an explicit attribute mode.
type MapResourceAttributeWithMode struct {
	resource_schema.MapAttribute
	attributeWithMode
}

// NumberResourceAttributeWithMode is a NumberAttribute with an explicit attribute mode.
type NumberResourceAttributeWithMode struct {
	resource_schema.NumberAttribute
	attributeWithMode
}

// ObjectResourceAttributeWithMode is a ObjectAttribute with an explicit attribute mode.
type ObjectResourceAttributeWithMode struct {
	resource_schema.ObjectAttribute
	attributeWithMode
}

// SetResourceAttributeWithMode is a SetAttribute with an explicit attribute mode.
type SetResourceAttributeWithMode struct {
	resource_schema.SetAttribute
	attributeWithMode
}

// StringResourceAttributeWithMode is a StringAttribute with an explicit attribute mode.
type StringResourceAttributeWithMode struct {
	resource_schema.StringAttribute
	attributeWithMode
}

var (
	_ AttributeWithMode = attributeWithMode{}
	_ AttributeWithMode = BoolResourceAttributeWithMode{}
	_ AttributeWithMode = Float32ResourceAttributeWithMode{}
	_ AttributeWithMode = Float64ResourceAttributeWithMode{}
	_ AttributeWithMode = Int32ResourceAttributeWithMode{}
	_ AttributeWithMode = Int64ResourceAttributeWithMode{}
	_ AttributeWithMode = ListResourceAttributeWithMode{}
	_ AttributeWithMode = MapResourceAttributeWithMode{}
	_ AttributeWithMode = NumberResourceAttributeWithMode{}
	_ AttributeWithMode = ObjectResourceAttributeWithMode{}
	_ AttributeWithMode = SetResourceAttributeWithMode{}
	_ AttributeWithMode = StringResourceAttributeWithMode{}
)

// GetResourceAttributeMode returns the mode for a resource attribute.
func GetResourceAttributeMode(attrSchema resource_schema.Attribute) AttributeMode {
	if m, ok := attrSchema.(AttributeWithMode); ok {
		return m.GetMode()
	}

	return ReadWriteAttributeMode
}
