// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

//go:generate go run ./protogen/ --outfile=illumio/cloud/config/v1/config.proto
//go:generate buf format -w illumio/cloud/config/v1/config.proto
//go:generate buf lint illumio/cloud/config/v1/config.proto
//go:generate buf generate illumio/cloud/config/v1/config.proto
