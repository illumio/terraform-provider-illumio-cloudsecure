// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type SchemaTestSuite struct {
	suite.Suite

	// Schema is the cloudSecureSchema to test.
	Schema Schema
}

func (suite *SchemaTestSuite) SetupTest() {
	suite.Schema = CloudSecure()
}

func TestCloudSecureSchemaTestSuite(t *testing.T) {
	suite.Run(t, new(SchemaTestSuite))
}
