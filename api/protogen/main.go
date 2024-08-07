// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/illumio/terraform-provider-illumio-cloudsecure/api/schema"
)

var (
	// outfile is the output filename, as set by the --outfile configuration flag.
	outfile string

	// tagsfile is the tags filename, as set by the --tagsfile configuration flag.
	tagsfile string
)

func init() {
	flag.StringVar(&outfile, "outfile", "config.proto", "name of the output *.proto file")
	flag.StringVar(&tagsfile, "tagsfile", "tags.json", "name of the JSON file containing previously assigned attribute tags")
}

func main() {
	flag.Parse()

	// Read the previously assigned attribute tags, so that adding a new attribute
	// doesn't modify the tags of existing attributes.
	var tagger *apiSpecTagger

	tagsJSONBytes, err := os.ReadFile(tagsfile)
	if err == nil {
		err = json.Unmarshal(tagsJSONBytes, &tagger)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse tags file %q: %s (%s)", tagsfile, err, string(tagsJSONBytes))
			os.Exit(1)
		}
	}

	if tagger == nil {
		tagger = newApiSpecTagger()
	}

	f, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open output file %q: %s", outfile, err)
		os.Exit(1)
	}

	err = GenerateGRPCAPISpec(f, schema.CloudSecure(), tagger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate and write API into output file %q: %s", outfile, err)
		os.Exit(1)
	}

	err = f.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to close output file %q: %s", outfile, err)
		os.Exit(1)
	}

	// Save the attribute tags for the next run.
	tagsJSONBytes, err = json.Marshal(tagger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to encode tags: %s", err)
		os.Exit(1)
	}

	err = os.WriteFile(tagsfile, tagsJSONBytes, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to write tags into tags file %q: %s", tagsfile, err)
		os.Exit(1)
	}
}
