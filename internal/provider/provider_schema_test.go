// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	provider_schema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProviderSchemaAccessTokenLength(t *testing.T) {
	t.Parallel()

	schemaResp := &provider.SchemaResponse{}
	New("test")().Schema(t.Context(), provider.SchemaRequest{}, schemaResp)
	require.False(t, schemaResp.Diagnostics.HasError(), "schema diagnostics: %v", schemaResp.Diagnostics)

	attr, ok := schemaResp.Schema.Attributes["access_token"].(provider_schema.StringAttribute)
	require.True(t, ok, "access_token must be a string attribute")

	tests := map[string]struct {
		length  int
		wantErr bool
	}{
		"empty":         {length: 0, wantErr: true},
		"one character": {length: 1, wantErr: false},
		// A CloudSecure service account token (RS256, 4096-bit key) was 1172 characters long.
		"1172 characters": {length: 1172, wantErr: false},
		"16 KiB":          {length: 16 * 1024, wantErr: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := validator.StringRequest{
				Path:        path.Root("access_token"),
				ConfigValue: types.StringValue(strings.Repeat("a", tc.length)),
			}

			var diags diag.Diagnostics

			for _, v := range attr.StringValidators() {
				resp := &validator.StringResponse{}
				v.ValidateString(t.Context(), req, resp)
				diags.Append(resp.Diagnostics...)
			}

			assert.Equal(t, tc.wantErr, diags.HasError(), "diagnostics: %v", diags)
		})
	}
}
