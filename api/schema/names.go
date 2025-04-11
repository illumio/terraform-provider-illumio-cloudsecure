// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"maps"
	"slices"

	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

const (
	// IDFieldName is the name of the "id" attribute/field that is used to identify any Illumio CloudSecure resource or data source.
	IDFieldName = "id"

	// UpdateMaskFieldName is the name of the "update_mask" field used in resource update requests to list the fields modified by the operation.
	UpdateMaskFieldName = "update_mask"
)

// ProtoMessageName converts a Terraform resource, data source, or field name (in snake_case) into a ProtocolBuffer message name (in UpperCamelCase).
func ProtoMessageName(tfName string) string {
	return generator.CamelCase(tfName)
}

// SortResourceAttributes returns the sorted names of a set of resource attributes.
func SortResourceAttributes(attrs map[string]resource_schema.Attribute) []string {
	names := slices.AppendSeq(make([]string, 0, len(attrs)), maps.Keys(attrs))
	sortAttributeNames(names)

	return names
}

func SortObjectAttributes(attrs map[string]attr.Type) []string {
	names := slices.AppendSeq(make([]string, 0, len(attrs)), maps.Keys(attrs))
	sortAttributeNames(names)

	return names
}

// SortDataSourceAttributes returns the sorted names of a set of data source attributes.
func SortDataSourceAttributes(attrs map[string]datasource_schema.Attribute) []string {
	names := slices.AppendSeq(make([]string, 0, len(attrs)), maps.Keys(attrs))
	sortAttributeNames(names)

	return names
}

func sortAttributeNames(names []string) {
	slices.SortFunc(names, func(a, b string) int {
		switch {
		case a == IDFieldName:
			return -1
		case b == IDFieldName:
			return 1
		case a < b:
			return -1
		case a > b:
			return 1
		default:
			return 0
		}
	})
}

// ProtoMessageNameForCreateRequest returns the name of the Protocol Buffer message for create requests for the given CamelCased resource name.
func ProtoMessageNameForCreateRequest(camelCaseResourcename string) string {
	return "Create" + camelCaseResourcename + "Request"
}

// ProtoMessageNameForCreatResponse returns the name of the Protocol Buffer message for create responses for the given CamelCased resource name.
func ProtoMessageNameForCreateResponse(camelCaseResourcename string) string {
	return "Create" + camelCaseResourcename + "Response"
}

// ProtoMessageNameForReadRequest returns the name of the Protocol Buffer message for read requests for the given CamelCased resource name.
func ProtoMessageNameForReadRequest(camelCaseResourcename string) string {
	return "Read" + camelCaseResourcename + "Request"
}

// ProtoMessageNameForReadResponse returns the name of the Protocol Buffer message for read responses for the given CamelCased resource name.
func ProtoMessageNameForReadResponse(camelCaseResourcename string) string {
	return "Read" + camelCaseResourcename + "Response"
}

// ProtoMessageNameForUpdateRequest returns the name of the Protocol Buffer message for update requests for the given CamelCased resource name.
func ProtoMessageNameForUpdateRequest(camelCaseResourcename string) string {
	return "Update" + camelCaseResourcename + "Request"
}

// ProtoMessageNameForUpdateResponse returns the name of the Protocol Buffer message for update responses for the given CamelCased resource name.
func ProtoMessageNameForUpdateResponse(camelCaseResourcename string) string {
	return "Update" + camelCaseResourcename + "Response"
}

// ProtoMessageNameForDeleteRequest returns the name of the Protocol Buffer message for delete requests for the given CamelCased resource name.
func ProtoMessageNameForDeleteRequest(camelCaseResourcename string) string {
	return "Delete" + camelCaseResourcename + "Request"
}

// RPCNameForCreate returns the name of the gRPC RPC used to create resources with the given CamelCased resource name.
func RPCNameForCreate(camelCaseResourcename string) string {
	return "Create" + camelCaseResourcename
}

// RPCNameForRead returns the name of the gRPC RPC used to read resources with the given CamelCased resource name.
func RPCNameForRead(camelCaseResourcename string) string {
	return "Read" + camelCaseResourcename
}

// RPCNameForUpdate returns the name of the gRPC RPC used to update resources with the given CamelCased resource name.
func RPCNameForUpdate(camelCaseResourcename string) string {
	return "Update" + camelCaseResourcename
}

// RPCNameForDelete returns the name of the gRPC RPC used to delete resources with the given CamelCased resource name.
func RPCNameForDelete(camelCaseResourcename string) string {
	return "Delete" + camelCaseResourcename
}
