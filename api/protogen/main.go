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
	// outfile is the output filename, as set by the --outfile configuration flag.
	outfile string
)

func init() {
	flag.StringVar(&outfile, "outfile", "config.proto", "name of the output *.proto file")
}

func main() {
	flag.Parse()

	f, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file %q: %s", outfile, err)
		os.Exit(1)
	}

	err = GenerateGRPCAPISpec(f, schema.CloudSecure(), newApiSpecTagger())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate and write API into file %q: %s", outfile, err)
		os.Exit(1)
	}

	err = f.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to close file %q: %s", outfile, err)
		os.Exit(1)
	}
}
