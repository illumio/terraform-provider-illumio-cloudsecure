// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import "encoding/json"

func (suite *GenerateTestSuite) TestAttributeTaggerAssignTag() {
	attrName := "new_field"
	tests := map[string]struct {
		preAssignedTags map[string]int
		expectedTag     int
	}{
		"no_preassigned_tags": {
			preAssignedTags: nil,
			expectedTag:     1,
		},
		"one_preassigned_tag_different": {
			preAssignedTags: map[string]int{
				"field1": 1,
			},
			expectedTag: 2,
		},
		"several_preassigned_tags_different": {
			preAssignedTags: map[string]int{
				"field1": 1,
				"field2": 2,
				"field3": 3,
			},
			expectedTag: 4,
		},
		"one_preassigned_tag_including_new_attr": {
			preAssignedTags: map[string]int{
				"new_field": 42,
			},
			expectedTag: 42,
		},
		"several_preassigned_tags_including_new_attr": {
			preAssignedTags: map[string]int{
				"field1":    1,
				"field2":    2,
				"new_field": 42,
			},
			expectedTag: 42,
		},
	}

	for name, tc := range tests {
		suite.Run(name, func() {
			tagger := newAttributeTagger(tc.preAssignedTags)
			gotTag := tagger.AssignTag(attrName)
			suite.Equal(tc.expectedTag, gotTag, "attribute tag should match")
		})
	}
}

func (suite *GenerateTestSuite) TestAttributeTaggerMarshalJSON() {
	tests := map[string]struct {
		preAssignedTags map[string]int
		expectedJSON    string
	}{
		"no_preassigned_tags": {
			preAssignedTags: nil,
			expectedJSON:    `{}`,
		},
		"one_preassigned_tag": {
			preAssignedTags: map[string]int{
				"field1": 1,
			},
			expectedJSON: `{"field1":1}`,
		},
		"two_preassigned_tags": {
			preAssignedTags: map[string]int{
				"field1": 1,
				"field2": 2,
			},
			expectedJSON: `{"field1":1,"field2":2}`,
		},
	}

	for name, tc := range tests {
		suite.Run(name, func() {
			tagger := newAttributeTagger(tc.preAssignedTags)
			gotJSONBytes, gotErr := json.Marshal(tagger)
			suite.Require().NoError(gotErr, "JSON marshaling failed")
			suite.JSONEq(tc.expectedJSON, string(gotJSONBytes), "JSON value should match")
		})
	}
}

func (suite *GenerateTestSuite) TestAttributeTaggerUnmarshalJSON() {
	tests := map[string]struct {
		json                 string
		expectedAssignedTags map[string]int
		expectedNextTag      int
	}{
		"no_preassigned_tags": {
			json:                 `{}`,
			expectedAssignedTags: map[string]int{},
			expectedNextTag:      1,
		},
		"one_preassigned_tag": {
			json: `{"field1":1}`,
			expectedAssignedTags: map[string]int{
				"field1": 1,
			},
			expectedNextTag: 2,
		},
		"two_preassigned_tags": {
			json: `{"field1":1,"field2":2}`,
			expectedAssignedTags: map[string]int{
				"field1": 1,
				"field2": 2,
			},
			expectedNextTag: 3,
		},
	}

	for name, tc := range tests {
		suite.Run(name, func() {
			tagger := newAttributeTagger(nil)
			gotErr := json.Unmarshal([]byte(tc.json), &tagger)
			suite.Require().NoError(gotErr, "JSON unmarshaling failed")
			suite.Equal(tc.expectedAssignedTags, tagger.AssignedTags, "assigned tags should match")
			suite.Equal(tc.expectedNextTag, tagger.NextTag, "next tag should match")
		})
	}
}

func (suite *GenerateTestSuite) TestAPISpecTaggerAssignTag() {
	var gotTag int

	tagger := newAPISpecTagger()

	// The first attribute in a new namespace should have tag 1.
	gotTag = tagger.AssignTag("resource/res1", "field1")
	suite.Equal(1, gotTag, "attribute tag should match")

	// The second attribute should have tag 2, etc.
	gotTag = tagger.AssignTag("resource/res1", "field2")
	suite.Equal(2, gotTag, "attribute tag should match")

	// The same tag should be assigned to an already assigned attribute.
	gotTag = tagger.AssignTag("resource/res1", "field1")
	suite.Equal(1, gotTag, "attribute tag should match")

	// The first attribute in a new namespace should have tag 1.
	gotTag = tagger.AssignTag("resource/res2", "field1")
	suite.Equal(1, gotTag, "attribute tag should match")

	// The second attribute should have tag 2, etc.
	gotTag = tagger.AssignTag("resource/res2", "field2")
	suite.Equal(2, gotTag, "attribute tag should match")
}

func (suite *GenerateTestSuite) TestAPISpecTaggerJSONMarshaling() {
	var gotJSONBytes []byte

	var gotErr error

	var gotTag int

	tagger := newAPISpecTagger()

	_ = tagger.AssignTag("resource/res1", "field1")
	_ = tagger.AssignTag("resource/res1", "field2")
	_ = tagger.AssignTag("resource/res2", "field1")
	_ = tagger.AssignTag("resource/res2", "field2")

	// Test that marshaling the tagger generates the correct JSON value.
	gotJSONBytes, gotErr = json.Marshal(tagger)
	suite.Require().NoError(gotErr, "JSON marshaling failed")
	suite.JSONEq(`{"resource/res1":{"field1":1,"field2":2},"resource/res2":{"field1":1,"field2":2}}`, string(gotJSONBytes), "JSON value should match")

	// Test that unmarshaling and marshaling again generates the same JSON value.
	tagger = nil
	gotErr = json.Unmarshal(gotJSONBytes, &tagger)
	suite.Require().NoError(gotErr, "JSON unmarshaling failed")
	gotJSONBytes, gotErr = json.Marshal(tagger)
	suite.Require().NoError(gotErr, "JSON marshaling failed")
	suite.JSONEq(`{"resource/res1":{"field1":1,"field2":2},"resource/res2":{"field1":1,"field2":2}}`, string(gotJSONBytes), "JSON value should match")

	// Test that the next tags were calculated correctly on unmarshaling.
	gotTag = tagger.AssignTag("resource/res1", "field3")
	suite.Equal(3, gotTag, "attribute tag should match")

	// Test that adding a new namespace works.
	_ = tagger.AssignTag("resource/res3", "field1")
	gotJSONBytes, gotErr = json.Marshal(tagger)
	suite.Require().NoError(gotErr, "JSON marshaling failed")
	suite.JSONEq(`{"resource/res1":{"field1":1,"field2":2,"field3":3},"resource/res2":{"field1":1,"field2":2},"resource/res3":{"field1":1}}`, string(gotJSONBytes), "JSON value should match")
}
