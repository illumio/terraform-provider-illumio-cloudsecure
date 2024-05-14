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

	// IdAttributeMode is the mode of "id" attributes.
	IdAttributeMode AttributeMode = AttributeMode{
		InCreateResponse: true,
		InReadRequest:    true,
		InReadResponse:   true,
		InUpdateRequest:  true,
		InUpdateResponse: true,
		InDeleteRequest:  true,
	}

	// ReadWriteAttributeMode is the mode of read-write attributes. This is the default mode.
	ReadWriteAttributeMode AttributeMode = AttributeMode{
		InCreateRequest:  true,
		InCreateResponse: true,
		InReadResponse:   true,
		InUpdateRequest:  true,
		InUpdateResponse: true,
	}

	// ImmutableAttributeMode is the mode of immutable attributes, which can be set only at creation and can be read afterwards. This is the default mode.
	// Attributes with this mode should also have the RequiresReplace plan modifier.
	ImmutableAttributeMode AttributeMode = AttributeMode{
		InCreateRequest:  true,
		InCreateResponse: true,
		InReadResponse:   true,
		InUpdateResponse: true,
	}

	// ReadOnlyAttributeMode is the mode of read-only attributes, which are returned by every create, read, and update operation.
	ReadOnlyAttributeMode AttributeMode = AttributeMode{
		InCreateResponse: true,
		InReadResponse:   true,
		InUpdateResponse: true,
	}

	// ReadOnlyOnceAttributeMode is the mode of read-only-once attributes, which are returned only by create operations.
	// Attributes with this mode should also have the UseStateForUnknown plan modifier.
	ReadOnlyOnceAttributeMode AttributeMode = AttributeMode{
		InCreateResponse: true,
	}

	// WriteOnlyAttributeMode is the mode of write-only attributes, which are sent in every create and update operation, and never in any response.
	// Attribute with this mode should also have the UseStateForUnknown plan modifier.
	WriteOnlyAttributeMode AttributeMode = AttributeMode{
		InCreateRequest: true,
		InUpdateRequest: true,
	}

	// WriteOnlyOnceAttributeMode is the mode of write-only-once attributes, which are sent in only in create operations.
	// Attribute with this mode should also have the UseStateForUnknown plan modifier.
	WriteOnlyOnceAttributeMode AttributeMode = AttributeMode{
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

// BoolAttributeWithMode is a BoolAttribute with an explicit attribute mode.
type BoolAttributeWithMode struct {
	resource_schema.BoolAttribute
	attributeWithMode
}

// Float64AttributeWithMode is a Float64Attribute with an explicit attribute mode.
type Float64AttributeWithMode struct {
	resource_schema.Float64Attribute
	attributeWithMode
}

// Int64AttributeWithMode is a Int64Attribute with an explicit attribute mode.
type Int64AttributeWithMode struct {
	resource_schema.Int64Attribute
	attributeWithMode
}

// ListAttributeWithMode is a ListAttribute with an explicit attribute mode.
type ListAttributeWithMode struct {
	resource_schema.ListAttribute
	attributeWithMode
}

// MapAttributeWithMode is a MapAttribute with an explicit attribute mode.
type MapAttributeWithMode struct {
	resource_schema.MapAttribute
	attributeWithMode
}

// NumberAttributeWithMode is a NumberAttribute with an explicit attribute mode.
type NumberAttributeWithMode struct {
	resource_schema.NumberAttribute
	attributeWithMode
}

// ObjectAttributeWithMode is a ObjectAttribute with an explicit attribute mode.
type ObjectAttributeWithMode struct {
	resource_schema.ObjectAttribute
	attributeWithMode
}

// SetAttributeWithMode is a SetAttribute with an explicit attribute mode.
type SetAttributeWithMode struct {
	resource_schema.SetAttribute
	attributeWithMode
}

// StringAttributeWithMode is a StringAttribute with an explicit attribute mode.
type StringAttributeWithMode struct {
	resource_schema.StringAttribute
	attributeWithMode
}

var (
	_ AttributeWithMode = attributeWithMode{}
	_ AttributeWithMode = BoolAttributeWithMode{}
	_ AttributeWithMode = Float64AttributeWithMode{}
	_ AttributeWithMode = Int64AttributeWithMode{}
	_ AttributeWithMode = ListAttributeWithMode{}
	_ AttributeWithMode = MapAttributeWithMode{}
	_ AttributeWithMode = NumberAttributeWithMode{}
	_ AttributeWithMode = ObjectAttributeWithMode{}
	_ AttributeWithMode = SetAttributeWithMode{}
	_ AttributeWithMode = StringAttributeWithMode{}
)

// GetAttributeMode returns the mode for a resource attribute.
func GetAttributeMode(attrSchema resource_schema.Attribute) AttributeMode {
	if m, ok := attrSchema.(AttributeWithMode); ok {
		return m.GetMode()
	} else {
		return ReadWriteAttributeMode
	}
}
