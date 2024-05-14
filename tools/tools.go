// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build tools

package tools

import (
	// Documentation generation
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"

	// Buf tools to format and compile Protocol Buffer / gRPC
	_ "github.com/bufbuild/buf/cmd/buf"
)
