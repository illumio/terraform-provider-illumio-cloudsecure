// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"io"
	"text/template" // nosemgrep: go.lang.security.audit.xss.import-text-template.import-text-template

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/illumio/terraform-provider-illumio-cloudsecure/api/schema"
)

var (
	// grpcAPISpecTemplate is the template of the gRPC spec for the Illumio CloudSecure Config API.
	//
	// TODO: Add comments into the generated spec from the schema's descriptions.
	// TODO: Make go_package configurable.
	grpcAPISpecTemplate = template.Must(template.New("apispec").Parse(`// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0
syntax = "proto3";
package illumio.cloud.config.{{.Version}};
import "google/protobuf/empty.proto";
import "google/protobuf/field_mask.proto";
service ConfigService {
	{{- range $rpc := .RPCs}}
	rpc {{$rpc.Name}}({{$rpc.RequestMessageName}}) returns ({{$rpc.ResponseMessageName}});
	{{- end}}
}
{{- define "message"}}
message {{.Name}} {
	{{- range $message := .Messages}}
	{{- template "message" $message}}
	{{- end}}
	{{- range $field := .Fields}}
	{{if $field.Repeated}}repeated {{else}}{{if $field.Optional}}optional {{end}}{{end}}{{$field.Type}} {{$field.Name}} = {{$field.Tag}};
	{{- end}}
}
{{- end}}
{{- range $message := .Messages}}
{{- template "message" $message}}
{{- end}}
`))
)

// grpcAPISpecTemplateData contains all the data needed to execute the grpcAPISpecTemplate to generate the API spec.
type grpcAPISpecTemplateData struct {
	// Version is the API version.
	Version string

	// RPCs is the list of RPCs in the API spec.
	RPCs []rpc

	// Messages is the list of messages in the API spec.
	Messages []message
}

// rpc is the definition of a gRPC RPC.
type rpc struct {
	// Name is the RPC name.
	Name string

	// RequestMessageName is the name of the RPC's request Protocol Buffer message.
	RequestMessageName string

	// ResponseMessageName is the name of the RPC's response Protocol Buffer message.
	ResponseMessageName string
}

// message is the definition of a Protocol Buffer message.
type message struct {
	// Name is the message name.
	Name string

	// Messages is the list of messages nested within the message.
	Messages []message

	// Fields is the list of fields in the message.
	Fields []field
}

// field is the definition of a field in a Protocol Buffer message.
type field struct {
	// Repeated is true is the field is repeated.
	Repeated bool

	// Optional is true is the field is optional.
	Optional bool

	// Type is the Protocol Buffer type of the field.
	Type string

	// Name is the field name.
	Name string

	// Tag is the Protocol Buffer tag of the field.
	Tag int
}

const MessageNameSeperator = "_"

// GenerateGRPCAPISpec generates the Protocol Buffer and gRPC spec for the Illumio CloudSecure Config API, and outputs it into the given Writer.
func GenerateGRPCAPISpec(dst io.Writer, src schema.Schema, tagger *apiSpecTagger) error {
	data := grpcAPISpecTemplateData{
		Version: src.Version(),
		// TODO: Support data sources.
		RPCs:     make([]rpc, 0, len(src.Resources())*4),
		Messages: make([]message, 0, len(src.Resources())*7),
	}

	for _, resource := range src.Resources() {
		resourceName := resource.TypeName
		resourceMessageName := schema.ProtoMessageName(resourceName)
		numFields := len(resource.Schema.Attributes)

		createRequestMessage := message{
			Name:   schema.ProtoMessageNameForCreateRequest(resourceMessageName),
			Fields: make([]field, 0, numFields-1),
		}

		createResponseMessage := message{
			Name:   schema.ProtoMessageNameForCreateResponse(resourceMessageName),
			Fields: make([]field, 0, numFields),
		}

		readRequestMessage := message{
			Name:   schema.ProtoMessageNameForReadRequest(resourceMessageName),
			Fields: make([]field, 0, 1),
		}

		readResponseMessage := message{
			Name:   schema.ProtoMessageNameForReadResponse(resourceMessageName),
			Fields: make([]field, 0, numFields),
		}

		updateRequestMessage := message{
			Name:   schema.ProtoMessageNameForUpdateRequest(resourceMessageName),
			Fields: make([]field, 0, numFields+1),
		}

		updateResponseMessage := message{
			Name:   schema.ProtoMessageNameForUpdateResponse(resourceMessageName),
			Fields: make([]field, 0, numFields),
		}

		deleteRequestMessage := message{
			Name:   schema.ProtoMessageNameForDeleteRequest(resourceMessageName),
			Fields: make([]field, 0, 1),
		}

		attrNames := schema.SortResourceAttributes(resource.Schema.Attributes)

		for _, attrName := range attrNames {
			attrSchema := resource.Schema.Attributes[attrName]

			repeated, t, msg, err := terraformAttributeTypeToProtoType(resourceName, attrName, attrSchema.GetType(), tagger)
			if err != nil {
				return fmt.Errorf("failed to parse field %s in resource %s: %w", attrName, resourceMessageName, err)
			}

			if msg != nil {
				// TODO: This should be nested into each message to prevent name collisions, or the message name must be prefixed to become unique.
				// The message used as this field's type must either:
				//   1. be defined globally, so it can be shared among the request/response messages for the resource, or
				//   2. be defined as a nested messaged within each request/response message for the resource, which provides more encapsulation.
				//
				// Nesting each message (2.) would cause too much overhead as it would require generating duplicate code for handling those duplicate messages.
				// Therefore, define it globally, but prefix the message name with the resource name to make it globally unique.
				msg.Name = resourceMessageName + MessageNameSeperator + msg.Name
				t = msg.Name

				data.Messages = append(data.Messages, *msg)
			}

			f := field{
				Repeated: repeated,
				Optional: schema.AttributeIsOptional(attrSchema),
				Type:     t,
				Name:     attrName,
				Tag:      tagger.AssignTag("resource/"+resourceName, attrName),
			}

			attrMode := schema.GetResourceAttributeMode(attrSchema)

			if attrMode.InCreateRequest {
				createRequestMessage.Fields = append(createRequestMessage.Fields, f)
			}

			if attrMode.InCreateResponse {
				createResponseMessage.Fields = append(createResponseMessage.Fields, f)
			}

			if attrMode.InReadRequest {
				readRequestMessage.Fields = append(readRequestMessage.Fields, f)
			}

			if attrMode.InReadResponse {
				readResponseMessage.Fields = append(readResponseMessage.Fields, f)
			}

			if attrMode.InUpdateRequest {
				updateRequestMessage.Fields = append(updateRequestMessage.Fields, f)
			}

			if attrMode.InUpdateResponse {
				updateResponseMessage.Fields = append(updateResponseMessage.Fields, f)
			}

			if attrMode.InDeleteRequest {
				deleteRequestMessage.Fields = append(deleteRequestMessage.Fields, f)
			}
		}

		updateRequestMessage.Fields = append(updateRequestMessage.Fields, field{
			Type: "google.protobuf.FieldMask",
			Name: schema.UpdateMaskFieldName,
			Tag:  tagger.AssignTag("resource/"+resourceName, schema.UpdateMaskFieldName),
		})

		data.RPCs = append(data.RPCs,
			rpc{
				Name:                schema.RPCNameForCreate(resourceMessageName),
				RequestMessageName:  createRequestMessage.Name,
				ResponseMessageName: createResponseMessage.Name,
			},
			rpc{
				Name:                schema.RPCNameForRead(resourceMessageName),
				RequestMessageName:  readRequestMessage.Name,
				ResponseMessageName: readResponseMessage.Name,
			},
			rpc{
				Name:                schema.RPCNameForUpdate(resourceMessageName),
				RequestMessageName:  updateRequestMessage.Name,
				ResponseMessageName: updateResponseMessage.Name,
			},
			rpc{
				Name:                schema.RPCNameForDelete(resourceMessageName),
				RequestMessageName:  deleteRequestMessage.Name,
				ResponseMessageName: "google.protobuf.Empty",
			},
		)

		data.Messages = append(data.Messages,
			createRequestMessage,
			createResponseMessage,
			readRequestMessage,
			readResponseMessage,
			updateRequestMessage,
			updateResponseMessage,
			deleteRequestMessage,
		)
	}

	return grpcAPISpecTemplate.Execute(dst, &data)
}

// terraformAttributeTypeToProtoType converts a Terraform attribute type into the corresponding Protocol Buffer type, and optionally additional Protocol Buffer messages that represent nested types.
func terraformAttributeTypeToProtoType(messageNamePrefix, attrName string, attrType attr.Type, tagger *apiSpecTagger) (repeated bool, protoType string, nestedMessage *message, err error) {
	switch v := attrType.(type) {
	case basetypes.BoolType:
		return false, "bool", nil, nil
	case basetypes.Float64Type:
		return false, "double", nil, nil
	case basetypes.Int64Type:
		return false, "int64", nil, nil
	case basetypes.StringType:
		return false, "string", nil, nil
	case types.ListType:
		return terraformRepeatedAttributeTypeToProtoType(messageNamePrefix, attrName, v.ElementType(), tagger)
	case types.SetType:
		return terraformRepeatedAttributeTypeToProtoType(messageNamePrefix, attrName, v.ElementType(), tagger)
	case types.ObjectType:
		return terraformObjectAttributeTypeToProtoType(messageNamePrefix, attrName, v, tagger)

	default:
		return false, "", nil, fmt.Errorf("unsupported Terraform type: %s", attrType.String())
	}
}

// terraformObjectAttributeTypeToProtoType converts a Terraform object attribute into a Protocol Buffer message type.
func terraformObjectAttributeTypeToProtoType(messageNamePrefix, attrName string, obj types.ObjectType, tagger *apiSpecTagger) (repeated bool, protoType string, nestedMessage *message, err error) {
	messageName := schema.ProtoMessageName(attrName)
	newMessage := &message{
		Name:   messageName,
		Fields: []field{},
	}
	wrappedMessageName := messageNamePrefix + MessageNameSeperator + messageName

	attrs := schema.SortObjectAttributes(obj.AttrTypes)

	for _, name := range attrs {
		attrType := obj.AttrTypes[name]
		isRepeated, t, msg, err := terraformAttributeTypeToProtoType(wrappedMessageName, name, attrType, tagger)

		if err != nil {
			return false, "", nil, fmt.Errorf("failed to convert field %s in object %s: %w", name, attrName, err)
		}

		newMessage.Fields = append(newMessage.Fields, field{
			Repeated: isRepeated,
			Type:     t,
			Name:     name,
			Tag:      tagger.AssignTag("resource/"+wrappedMessageName, name),
			Optional: false,
		})

		if msg != nil {
			newMessage.Messages = append(newMessage.Messages, *msg)
		}
	}

	return false, messageName, newMessage, nil
}

// terraformRepeatedAttributeTypeToProtoType converts a Terraform repeated attribute type into the corresponding Protocol Buffer type, and optionally additional Protocol Buffer messages that represent nested types.
func terraformRepeatedAttributeTypeToProtoType(messageNamePrefix, attrName string, elementType attr.Type, tagger *apiSpecTagger) (repeated bool, protoType string, nestedMessage *message, err error) {
	elemRepeated, elemProtoType, elemMessage, err := terraformAttributeTypeToProtoType(messageNamePrefix, attrName, elementType, tagger)

	switch {
	case err != nil:
		return false, "", nil, fmt.Errorf("unsupported element type %s: %w", elementType.String(), err)

	case elemRepeated: // The element type itself is repeated.
		// The attribute is a set of lists or a set of sets. This must be modeled in Protocol Buffer as a repeated field of a message type, which itself contains a repeated field.
		// In case an extra message is created for a nested field type, it will be named with the CamelCased attribute name.
		wrapperMessageName := schema.ProtoMessageName(attrName)

		wrapperMessage := &message{
			Name: wrapperMessageName,
			Fields: []field{
				{
					Repeated: true,
					Type:     elemProtoType,
					Name:     attrName,
					Tag:      1,
				},
			},
		}
		if elemMessage != nil {
			wrapperMessage.Messages = []message{*elemMessage}
		}

		return true, wrapperMessageName, wrapperMessage, nil

	default: // The element type is not repeated. Normal case.
		return true, elemProtoType, elemMessage, nil
	}
}
