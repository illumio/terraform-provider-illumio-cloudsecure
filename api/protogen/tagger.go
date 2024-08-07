// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import "encoding/json"

// attributeTagger assigns unique tags to the attributes of a resource or data source, which can be used as the tags of corresdponding fields in a Protocol Buffer message.
type attributeTagger struct {
	// AssignedTags maps names of attributes to the tags that are assigned to them.
	AssignedTags map[string]int

	// NextTag is the next tag in the resource or data source that has not been previously assigned.
	NextTag int
}

var (
	_ json.Marshaler   = &attributeTagger{}
	_ json.Unmarshaler = &attributeTagger{}
)

// newAttributeTagger creates a new attributeTagger with the given previously assigned tags.
// The values in preAssignedTags must be distinct.
func newAttributeTagger(preAssignedTags map[string]int) *attributeTagger {
	// Since all tags are assigned monotonically, the next tag immediately follows the highest previously assigned tag.
	// If no tags were previouysly assigned, start from 1.
	nextTag := 1

	if preAssignedTags == nil {
		preAssignedTags = make(map[string]int)
	}

	for _, tag := range preAssignedTags {
		if nextTag <= tag {
			nextTag = tag + 1
		}
	}

	return &attributeTagger{
		AssignedTags: preAssignedTags,
		NextTag:      nextTag,
	}
}

// AssignTag assigns a unique tag to the attribute with the given name.
func (t *attributeTagger) AssignTag(attrName string) int {
	// If the attribute was previously assigned a tag, assign that same tag.
	if tag, alreadyAssigned := t.AssignedTags[attrName]; alreadyAssigned {
		return tag
	}

	// Otherwise, assign the next smallest tag that was not previously assigned.
	tag := t.NextTag

	t.AssignedTags[attrName] = tag
	t.NextTag = tag + 1

	return tag
}

func (t *attributeTagger) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.AssignedTags)
}

func (t *attributeTagger) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &t.AssignedTags)
	if err != nil {
		return err
	}

	if t.AssignedTags == nil {
		t.AssignedTags = make(map[string]int)
	}

	t.NextTag = 1

	for _, tag := range t.AssignedTags {
		if t.NextTag <= tag {
			t.NextTag = tag + 1
		}
	}

	return nil
}

// apiSpecTagger assigns tags to the attributes of all the resources and data sources in an API spec.
type apiSpecTagger struct {
	// AttributeTaggers maps a namespace (resource or data source name) to the attributeTagger that can be used to tag its attributes.
	// Each namespace should be prefixed with "resource/" or "dataSource/" to make sure that each namespace is unique.
	AttributeTaggers map[string]*attributeTagger
}

var (
	_ json.Marshaler   = &apiSpecTagger{}
	_ json.Unmarshaler = &apiSpecTagger{}
)

// newAPISpecTagger creates a new empty apiSpecTagger.
func newAPISpecTagger() *apiSpecTagger {
	return &apiSpecTagger{
		AttributeTaggers: make(map[string]*attributeTagger),
	}
}

// AssignTag assigns a unique tag to the attribute with the given name in the given namespace.
// Each namespace should be "resource/" + a resource name, or "dataSource/" + a data source name.
func (t *apiSpecTagger) AssignTag(namespace, attrName string) int {
	var tagger *attributeTagger

	tagger, found := t.AttributeTaggers[namespace]
	if !found {
		// If the namespace is new, start with no pre-assigned tags.
		tagger = newAttributeTagger(nil)
		t.AttributeTaggers[namespace] = tagger
	}

	return tagger.AssignTag(attrName)
}

func (t *apiSpecTagger) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.AttributeTaggers)
}

func (t *apiSpecTagger) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.AttributeTaggers)
}
