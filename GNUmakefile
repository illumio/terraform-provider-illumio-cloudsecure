# Copyright (c) Illumio, Inc.
# SPDX-License-Identifier: MPL-2.0

BINARY_NAME=terraform-provider-illumio-cloudsecure
VERSION=dev
PLUGIN_DIR=~/.terraform.d/plugins/registry.terraform.io/illumio/illumio-cloudsecure
OS ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH ?= $(shell uname -m)
ifeq ($(ARCH),x86_64)
    ARCH=amd64
else ifeq ($(ARCH),arm64)
    ARCH=arm64
else ifeq ($(ARCH),aarch64)
    ARCH=arm64
endif
TF_PLUGIN_DIR = $(OS)_$(ARCH)

default: testacc

# Run acceptance tests
.PHONY: testacc clean-build clean generate build run reset lint

testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m


clean: clean-build
	find . -name "*.gen.go" -exec rm -f {} \;
	find . -name "*.pb.go" -exec rm -f {} \;
	find . -name "*config.proto" -exec rm -f {} \;
	rm -rf ./docs
	rm -rf ./terraform-provider-illumio-cloudsecure

clean-build:
	rm -rf ~/.terraform.d/plugins/registry.terraform.io/illumio/illumio-cloudsecure/

generate:
	(cd api ; go generate) && \
    (cd internal/provider ; go generate) && \
    (cd fakeserver ; go generate) && \
    (go generate ./...)

build:
	# Create the appropriate directory and copy the built binary
	mkdir -p $(PLUGIN_DIR)/$(VERSION)/$(TF_PLUGIN_DIR)
	# Build the Go project
	go build -o $(PLUGIN_DIR)/$(VERSION)/$(TF_PLUGIN_DIR)/$(BINARY_NAME)


run:
	go run ./fakeserver/ -apiEndpoint 0.0.0.0:50123 -tokenEndpoint 0.0.0.0:50124

reset: clean generate build

lint: 
	golangci-lint run