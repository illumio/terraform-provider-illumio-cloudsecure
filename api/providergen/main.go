// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/illumio/terraform-provider-illumio-cloudsecure/api/schema"
)

var (
	// goPackage is the Golang package name, as set by the --go-package configuration flag.
	goPackage string

	// outfile is the output filename, as set by the --outfile configuration flag.
	outfile string
)

func init() {
	flag.StringVar(&goPackage, "go-package", "provider", "Golang package name")
	flag.StringVar(&outfile, "outfile", "provider_impl.go", "name of the output *.go file")
}

func main() {
	flag.Parse()

	f, err := os.OpenFile(filepath.Clean(outfile), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file %q: %s", outfile, err)
		os.Exit(1)
	}

	err = GenerateProvider(f, goPackage, schema.CloudSecure())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate and write provider implementation into file %q: %s", outfile, err)
		os.Exit(1)
	}

	err = f.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to close file %q: %s", outfile, err)
		os.Exit(1)
	}
}
