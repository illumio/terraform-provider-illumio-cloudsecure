// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/illumio/terraform-provider-illumio-cloudsecure/api/schema"
)

var (
	// goPackage is the Golang package name, as set by the --go-package configuration flag.
	goPackage string

	// outfile is the output filename, as set by the --outfile configuration flag.
	outfile string
)

func init() {
	flag.StringVar(&goPackage, "go-package", "main", "Golang package name")
	flag.StringVar(&outfile, "outfile", "fakeserver_impl.go", "name of the output *.go file")
}

func main() {
	flag.Parse()

	file, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file %q: %s", outfile, err)
		os.Exit(1)
	}

	err = GenerateFakeServer(file, goPackage, schema.CloudSecure())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate and write fake server implementation into file %q: %s", outfile, err)
		os.Exit(1)
	}

	err = file.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to close file %q: %s", outfile, err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "Completed fake server implementation")
}
