// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// Resource is the complete definition of a resource type.
type Resource struct {
	// TypeName is the name of this resource type.
	TypeName string

	// Schema is the schema of this resource type.
	Schema resource_schema.Schema
}

// Resources is a list of Resource elements.
// The TypeName of each Resource must be unique.
type Resources []Resource

func (r Resources) Len() int           { return len(r) }
func (r Resources) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r Resources) Less(i, j int) bool { return r[i].TypeName < r[j].TypeName }

// DataSource is the complete definition of a data source type.
type DataSource struct {
	// TypeName is the name of this data source type.
	TypeName string

	// Schema is the schema of this data source type.
	Schema datasource_schema.Schema
}

// Resources is a list of Resource elements.
// The TypeName of each Resource must be unique.
type DataSources []DataSource

func (r DataSources) Len() int           { return len(r) }
func (r DataSources) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r DataSources) Less(i, j int) bool { return r[i].TypeName < r[j].TypeName }

// Schema contains the definitions of all the resources and data sources provided by a provider.
type Schema interface {
	// Version returns the version of this schema: "v1", "v2", etc.
	Version() string

	// Resources returns the list of resources provided by the provider sorted by TypeName.
	Resources() Resources

	// DataSources returns the list of data sources provided by the provider sorted by TypeName.
	DataSources() DataSources
}
