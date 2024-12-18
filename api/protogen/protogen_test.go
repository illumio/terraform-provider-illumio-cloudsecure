// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"bytes"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/illumio/terraform-provider-illumio-cloudsecure/util"
	"github.com/stretchr/testify/suite"
)

type GenerateTestSuite struct {
	suite.Suite
}

func TestGenerateTestSuite(t *testing.T) {
	suite.Run(t, new(GenerateTestSuite))
}

func (suite *GenerateTestSuite) TestTerraformAttributeTypeToProtoType() {
	fieldName := "the_field"
	tests := map[string]struct {
		tfType           attr.Type
		expectedRepeated bool
		expectedType     string
		expectedMessage  *message
		expectedErr      error
	}{
		"bool": {
			tfType:       types.BoolType,
			expectedType: "bool",
		},
		"float64": {
			tfType:       types.Float64Type,
			expectedType: "double",
		},
		"int64": {
			tfType:       types.Int64Type,
			expectedType: "int64",
		},
		"string": {
			tfType:       types.StringType,
			expectedType: "string",
		},
		"repeated-list-string": {
			tfType: types.ListType{
				ElemType: types.StringType,
			},
			expectedRepeated: true,
			expectedType:     "string",
		},
		"repeated-set-string": {
			tfType: types.SetType{
				ElemType: types.StringType,
			},
			expectedRepeated: true,
			expectedType:     "string",
		},
		"repeated-set-set-string": {
			tfType: types.SetType{
				ElemType: types.SetType{
					ElemType: types.StringType,
				},
			},
			expectedRepeated: true,
			expectedType:     "TheField",
			expectedMessage: &message{
				Name: "TheField",
				Fields: []field{
					{
						Repeated: true,
						Type:     "string",
						Name:     "the_field",
						Tag:      1,
					},
				},
			},
		},
		"repeated-set-list-set-string": {
			tfType: types.SetType{
				ElemType: types.ListType{
					ElemType: types.SetType{
						ElemType: types.StringType,
					},
				},
			},
			expectedRepeated: true,
			expectedType:     "TheField",
			expectedMessage: &message{
				Name: "TheField",
				Fields: []field{
					{
						Repeated: true,
						Type:     "TheField",
						Name:     "the_field",
						Tag:      1,
					},
				},
				Messages: []message{
					{
						Name: "TheField",
						Fields: []field{
							{
								Repeated: true,
								Type:     "string",
								Name:     "the_field",
								Tag:      1,
							},
						},
					},
				},
			},
		},
		"object": {
			tfType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"nested": types.StringType,
				},
			},
			expectedType: "TheField",
			expectedMessage: &message{
				Name: "TheField",
				Fields: []field{
					{
						Type: "string",
						Name: "nested",
						Tag:  1,
					},
				},
			},
		},
		"object-with-nested-object": {
			tfType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"name": types.StringType,
					"phone_numbers": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"office": types.StringType,
							"mobile": types.StringType,
						},
					},
				},
			},
			expectedType: "TheField",
			expectedMessage: &message{
				Name: "TheField",
				Messages: []message{
					{
						Name: "PhoneNumbers",
						Fields: []field{
							{
								Type: "string",
								Name: "mobile",
								Tag:  1,
							},
							{
								Type: "string",
								Name: "office",
								Tag:  2,
							},
						},
					},
				},
				Fields: []field{
					{
						Type: "string",
						Name: "name",
						Tag:  1,
					},
					{
						Type: "PhoneNumbers",
						Name: "phone_numbers",
						Tag:  2,
					},
				},
			},
		},
		"list-nested-objects": {
			tfType: types.ListType{
				ElemType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"email": types.StringType,
						"name":  types.StringType,
					},
				},
			},
			expectedRepeated: true,
			expectedType:     "TheField",
			expectedMessage: &message{
				Name:     "TheField",
				Messages: nil,
				Fields: []field{
					{
						Type: "string",
						Name: "email",
						Tag:  1,
					},
					{
						Type: "string",
						Name: "name",
						Tag:  2,
					},
				},
			},
		},
	}

	for name, tc := range tests {
		suite.Run(name, func() {
			gotRepeated, gotType, gotMessages, gotErr := terraformAttributeTypeToProtoType(fieldName, tc.tfType)
			suite.Equal(tc.expectedRepeated, gotRepeated, "Protocol Buffer repeated flag should match")
			suite.Equal(tc.expectedType, gotType, "Protocol Buffer type should match")
			suite.Equal(tc.expectedMessage, gotMessages, "Messages should match")
			suite.Equal(tc.expectedErr, gotErr, "Error should match")
		})
	}
}

func (suite *GenerateTestSuite) TestGRPCAPISpecTemplateMessage() { //nolint:maintidx
	tests := map[string]struct {
		message message
		output  string
	}{
		"all_field_modifiers": {
			message: message{
				Name:     "FooBar",
				Messages: nil,
				Fields: []field{
					{
						Repeated: false,
						Optional: false,
						Type:     "string",
						Name:     "not_repeated_not_optional",
						Tag:      1,
					},
					{
						Repeated: true,
						Optional: false,
						Type:     "string",
						Name:     "repeated_not_optional",
						Tag:      2,
					},
					{
						Repeated: false,
						Optional: true,
						Type:     "string",
						Name:     "not_repeated_optional",
						Tag:      3,
					},
					{
						Repeated: true,
						// Optional is ignored when Repeated is true, since all repeated fields are implicitly optional in Protocol Buffers.
						// So this case is equivalent to Repeated: true and Optional: false.
						Optional: true,
						Type:     "string",
						Name:     "repeated_optional",
						Tag:      4,
					},
				},
			},
			output: `
				message FooBar {
					string not_repeated_not_optional = 1;
					repeated string repeated_not_optional = 2;
					optional string not_repeated_optional = 3;
					repeated string repeated_optional = 4;
				}`,
		},
		"nested_message": {
			message: message{
				Name: "TopLevel",
				Messages: []message{
					{
						Name: "Nested",
						Fields: []field{
							{
								Type: "string",
								Name: "data",
								Tag:  1,
							},
						},
					},
				},
				Fields: []field{
					{
						Type: "Nested",
						Name: "nested",
						Tag:  1,
					},
				},
			},
			output: `
				message TopLevel {
					message Nested {
						string data = 1;
					}
					Nested nested = 1;
				}`,
		},
		"nested_nested_message": {
			message: message{
				Name: "TopLevel",
				Messages: []message{
					{
						Name: "Nested",
						Messages: []message{
							{
								Name: "Bottom",
								Fields: []field{
									{
										Type: "string",
										Name: "data",
										Tag:  1,
									},
								},
							},
						},
						Fields: []field{
							{
								Type: "Bottom",
								Name: "bottom",
								Tag:  1,
							},
						},
					},
				},
				Fields: []field{
					{
						Type: "Nested",
						Name: "nested",
						Tag:  1,
					},
				},
			},
			output: `
				message TopLevel {
					message Nested {
						message Bottom {
							string data = 1;
						}
						Bottom bottom = 1;
					}
					Nested nested = 1;
				}`,
		},
		"multiple_nested_messages": {
			message: message{
				Name: "TopLevel",
				Messages: []message{
					{
						Name: "A",
						Messages: []message{
							{
								Name: "A1",
								Fields: []field{
									{
										Type: "string",
										Name: "data1",
										Tag:  1,
									},
								},
							},
							{
								Name: "A2",
								Fields: []field{
									{
										Type: "string",
										Name: "data2",
										Tag:  1,
									},
								},
							},
						},
						Fields: []field{
							{
								Type: "A1",
								Name: "a1",
								Tag:  1,
							},
							{
								Type: "A2",
								Name: "a2",
								Tag:  2,
							}},
					},
					{
						Name: "B",
						Messages: []message{
							{
								Name: "B1",
								Fields: []field{
									{
										Type: "string",
										Name: "data1",
										Tag:  1,
									},
								},
							},
							{
								Name: "B2",
								Fields: []field{
									{
										Type: "string",
										Name: "data2",
										Tag:  1,
									},
								},
							},
						},
						Fields: []field{
							{
								Type: "B1",
								Name: "b1",
								Tag:  1,
							},
							{
								Type: "B2",
								Name: "b2",
								Tag:  2,
							}},
					},
				},
				Fields: []field{
					{
						Type: "A",
						Name: "a",
						Tag:  1,
					},
					{
						Type: "B",
						Name: "b",
						Tag:  2,
					},
				},
			},
			output: `
				message TopLevel {
					message A {
						message A1 {
							string data1 = 1;
						}
						message A2 {
							string data2 = 1;
						}
						A1 a1 = 1;
						A2 a2 = 2;
					}
					message B {
						message B1 {
							string data1 = 1;
						}
						message B2 {
							string data2 = 1;
						}
						B1 b1 = 1;
						B2 b2 = 2;
					}
					A a = 1;
					B b = 2;
				}`,
		},
		"nested_message_with_repeated_field": {
			message: message{
				Name: "TopLevel",
				Messages: []message{
					{
						Name: "Nested",
						Messages: []message{
							{
								Name: "Address",
								Fields: []field{
									{
										Type: "string",
										Name: "street",
										Tag:  1,
									},
									{
										Type: "string",
										Name: "city",
										Tag:  2,
									},
									{
										Type: "string",
										Name: "state",
										Tag:  3,
									},
								},
							},
						},
						Fields: []field{
							{
								Repeated: false,
								Type:     "string",
								Name:     "address_name",
								Tag:      1,
							},
							{
								Repeated: true,
								Name:     "addresses",
								Type:     "Address",
								Tag:      2,
							},
						},
					},
				},
			},
			output: `
				message TopLevel {
					message Nested {
						message Address {
							string street = 1;
							string city = 2;
							string state = 3;
						}
						string address_name = 1;
						repeated Address addresses = 2;
					}
				}
			`,
		},
		"list_of_nested_messages": {
			message: message{
				Name: "TopLevel",
				Messages: []message{
					{
						Name: "Address",
						Fields: []field{
							{
								Type: "string",
								Name: "state",
								Tag:  1,
							},
							{
								Type: "string",
								Name: "city",
								Tag:  2,
							},
						},
					},
				},
				Fields: []field{
					{
						Repeated: true,
						Type:     "Address",
						Optional: true,
						Name:     "addresses",
						Tag:      1,
					},
				},
			},
			output: `
				message TopLevel {
					message Address {
						string state = 1;
						string city = 2;
					}
					repeated Address addresses = 1;
				}
			`,
		},
	}

	for name, tc := range tests {
		suite.Run(name, func() {
			var buf bytes.Buffer
			err := grpcAPISpecTemplate.ExecuteTemplate(&buf, "message", tc.message)

			if suite.NoError(err, "template execution failed") {
				got := buf.String()
				suite.Equal(util.TrimEmptyLinesAndSpaces(tc.output), util.TrimEmptyLinesAndSpaces(got), "generated text must match")
			}
		})
	}
}

func (suite *GenerateTestSuite) TestGRPCAPISpecTemplate() {
	tests := map[string]struct {
		data   grpcAPISpecTemplateData
		output string
	}{
		"only_rpcs": {
			data: grpcAPISpecTemplateData{
				Version: "1.2.3",
				RPCs: []rpc{
					{
						Name:                "DoSomething1",
						RequestMessageName:  "RequestMessage1",
						ResponseMessageName: "ResponseMessage1",
					},
					{
						Name:                "DoSomething2",
						RequestMessageName:  "RequestMessage2",
						ResponseMessageName: "ResponseMessage2",
					},
				},
			},
			output: `
				// Copyright (c) Illumio, Inc.
				// SPDX-License-Identifier: MPL-2.0
				syntax = "proto3";
				package illumio.cloud.config.1.2.3;

				import "google/protobuf/empty.proto";
				import "google/protobuf/field_mask.proto";

				service ConfigService {
					rpc DoSomething1(RequestMessage1) returns (ResponseMessage1);
					rpc DoSomething2(RequestMessage2) returns (ResponseMessage2);
				}`,
		},
		"rpcs_and_messages": {
			data: grpcAPISpecTemplateData{
				Version: "1.2.3",
				RPCs: []rpc{
					{
						Name:                "DoSomething1",
						RequestMessageName:  "RequestMessage1",
						ResponseMessageName: "ResponseMessage1",
					},
					{
						Name:                "DoSomething2",
						RequestMessageName:  "RequestMessage2",
						ResponseMessageName: "ResponseMessage2",
					},
				},
				Messages: []message{
					{
						Name: "RequestMessage1",
						Messages: []message{
							{
								Name: "Nested",
								Fields: []field{
									{
										Repeated: true,
										Type:     "string",
										Name:     "strings",
										Tag:      1,
									},
								},
							},
						},
						Fields: []field{
							{
								Type: "string",
								Name: "id",
								Tag:  1,
							},
							{
								Repeated: true,
								Type:     "Nested",
								Name:     "list_of_list_of_strings",
								Tag:      2,
							},
						},
					},
					{
						Name: "ResponseMessage1",
						Fields: []field{
							{
								Type: "string",
								Name: "id",
								Tag:  1,
							},
							{
								Optional: true,
								Type:     "string",
								Name:     "optional_string",
								Tag:      2,
							},
						},
					},
					{
						Name: "RequestMessage2",
						Fields: []field{
							{
								Type: "string",
								Name: "id",
								Tag:  1,
							},
							{
								Repeated: true,
								Type:     "string",
								Name:     "list_of_strings",
								Tag:      2,
							},
						},
					},
					{
						Name: "ResponseMessage2",
						Fields: []field{
							{
								Type: "string",
								Name: "id",
								Tag:  1,
							},
						},
					},
					{
						Name: "NestedMessageParent",
						Messages: []message{
							{
								Name: "NestedMessageChild",
								Fields: []field{
									{
										Type: "string",
										Name: "child_id",
										Tag:  1,
									},
								},
							},
						},
						Fields: []field{
							{
								Type: "string",
								Name: "id",
								Tag:  1,
							},
						},
					},
				},
			},
			output: `
				// Copyright (c) Illumio, Inc.
				// SPDX-License-Identifier: MPL-2.0
				syntax = "proto3";
				package illumio.cloud.config.1.2.3;
				import "google/protobuf/empty.proto";
				import "google/protobuf/field_mask.proto";
				service ConfigService {
					rpc DoSomething1(RequestMessage1) returns (ResponseMessage1);
					rpc DoSomething2(RequestMessage2) returns (ResponseMessage2);
				}
				message RequestMessage1 {
				message Nested {
					repeated string strings = 1;
				}
					string id = 1;
					repeated Nested list_of_list_of_strings = 2;
				}
				message ResponseMessage1 {
					string id = 1;
					optional string optional_string = 2;
				}
				message RequestMessage2 {
					string id = 1;
					repeated string list_of_strings = 2;
				}
				message ResponseMessage2 {
					string id = 1;
				}
				message NestedMessageParent {
					message NestedMessageChild {
						string child_id = 1;
					}
					string id = 1;
				}`,
		},
	}

	for name, tc := range tests {
		suite.Run(name, func() {
			var buf bytes.Buffer
			err := grpcAPISpecTemplate.Execute(&buf, tc.data)

			if suite.NoError(err, "template execution failed") {
				got := buf.String()
				suite.Equal(util.TrimEmptyLinesAndSpaces(tc.output), util.TrimEmptyLinesAndSpaces(got), "generated text must match")
			}
		})
	}
}
