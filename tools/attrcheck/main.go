// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

// attrcheck is a tiny linter that checks that Terraform schema attributes and
// object AttrTypes map literals are ordered alphabetically, with the special
// case that the "id" field (referenced either as the string literal "id" or
// the identifier IDFieldName) may appear first.
package main

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

func main() {
	flag.Usage = func() {
		_, err := fmt.Fprintf(os.Stderr, "Usage: %s [path ...]\n", filepath.Base(os.Args[0]))
		if err != nil {
			return
		}

		_, err = fmt.Fprintf(os.Stderr, "Scans Go files under the given paths (files or directories) and verifies ordering of keys in map literals assigned to 'Attributes' or 'AttrTypes'. 'id' may appear first.\n")
		if err != nil {
			return
		}
	}

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	}

	var hadErr bool

	for _, arg := range args {
		err := processArg(arg, &hadErr)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "attrcheck: %v\n", err)
			hadErr = true
		}
	}

	if hadErr {
		os.Exit(1)
	}
}

func processArg(path string, hadErr *bool) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}

	if fi.IsDir() {
		return walkDir(path, hadErr)
	}

	if filepath.Ext(path) == ".go" {
		err = checkFile(path)
		if err != nil {
			_, werr := fmt.Fprintln(os.Stderr, err.Error())
			if werr != nil {
				return werr
			}

			*hadErr = true
		}
	}

	return nil
}

func walkDir(root string, hadErr *bool) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			// Skip vendor and hidden directories
			name := d.Name()
			if name == "vendor" || (len(name) > 0 && name[0] == '.') {
				return filepath.SkipDir
			}

			return nil
		}

		if filepath.Ext(path) != ".go" {
			return nil
		}

		err = checkFile(path)
		if err != nil {
			_, werr := fmt.Fprintln(os.Stderr, err.Error())
			if werr != nil {
				return werr
			}

			*hadErr = true
		}

		return nil
	})
}

func checkFile(filename string) error {
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("%s: parse error: %w", filename, err)
	}

	var errs []error

	ast.Inspect(file, func(n ast.Node) bool {
		kv, ok := n.(*ast.KeyValueExpr)
		if !ok {
			return true
		}
		// We care only for keys named Attributes or AttrTypes
		ident, ok := kv.Key.(*ast.Ident)
		if !ok {
			return true
		}

		if ident.Name != "Attributes" && ident.Name != "AttrTypes" {
			return true
		}

		cl, ok := kv.Value.(*ast.CompositeLit)
		if !ok {
			return true
		}
		// Ensure it's a map literal
		if _, ok := cl.Type.(*ast.MapType); !ok {
			return true
		}

		keys, positions := extractMapKeys(cl)
		if len(keys) == 0 {
			return true
		}

		err = validateKeyOrder(keys)
		if err != nil {
			// Attach file/line of the map literal start
			pos := fset.Position(cl.Lbrace)
			errs = append(errs, fmt.Errorf("%s:%d:%d: %s (%s)", filename, pos.Line, pos.Column, err.Error(), ident.Name))
			// Also print the encountered order to help fixing
			errs = append(errs, fmt.Errorf("%s: encountered order: %v", filename, keys))
			// Optionally, report the first out-of-order element location
			if i := firstOutOfOrderIndex(keys); i >= 0 && i < len(positions) {
				pp := fset.Position(positions[i])
				errs = append(errs, fmt.Errorf("%s:%d:%d: first out-of-order key here: %q", filename, pp.Line, pp.Column, keys[i]))
			}
		}

		return true
	})

	return errors.Join(errs...)
}

func extractMapKeys(cl *ast.CompositeLit) ([]string, []token.Pos) {
	keys := make([]string, 0, len(cl.Elts))

	positions := make([]token.Pos, 0, len(cl.Elts))
	for _, elt := range cl.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		name := mapKeyName(kv.Key)
		if name == "" {
			continue
		}

		keys = append(keys, name)
		positions = append(positions, kv.Key.Pos())
	}

	return keys, positions
}

func mapKeyName(e ast.Expr) string {
	switch v := e.(type) {
	case *ast.Ident:
		if v.Name == "IDFieldName" {
			return "id"
		}

		return v.Name // unlikely for map keys, but keep for completeness
	case *ast.BasicLit:
		if v.Kind == token.STRING {
			unq, err := strconv.Unquote(v.Value)
			if err == nil {
				return unq
			}

			return v.Value
		}
	}

	return ""
}

func validateKeyOrder(keys []string) error {
	// Allow id first if present
	start := 0
	if len(keys) > 0 && keys[0] == "id" {
		start = 1
	}

	if start >= len(keys) {
		return nil
	}

	got := append([]string(nil), keys[start:]...)
	expected := append([]string(nil), got...)
	sort.Strings(expected)

	for i := range got {
		if got[i] != expected[i] {
			return errors.New("map keys must be alphabetically ordered after 'id' (if present)")
		}
	}

	return nil
}

func firstOutOfOrderIndex(keys []string) int {
	start := 0
	if len(keys) > 0 && keys[0] == "id" {
		start = 1
	}

	prev := ""
	for i := start; i < len(keys); i++ {
		if i > start && keys[i] < prev {
			return i
		}

		prev = keys[i]
	}

	return -1
}
