// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"bytes"
	"strings"
)

// TrimEmptyLinesAndSpaces returns the given string with all empty lines and leading and trailing spaces (and tabs, etc.) removed.
func TrimEmptyLinesAndSpaces(s string) string {
	lines := strings.Split(s, "\n")

	var buf bytes.Buffer

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			buf.WriteString(line)
			buf.WriteRune('\n')
		}
	}

	return buf.String()
}
