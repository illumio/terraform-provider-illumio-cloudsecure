// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// urlStringValidator validates that a schema attribute's string value is a valid RFC 3986 URL.
type urlStringValidator struct {
}

func (v *urlStringValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v *urlStringValidator) MarkdownDescription(ctx context.Context) string {
	return "value must be a valid RFC 3986 URL"
}

// ValidateString implements validator.String.
func (v *urlStringValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	_, err := url.Parse(value)
	if err != nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			v.Description(ctx),
			value,
		))
	}
}

var _ validator.String = &urlStringValidator{}

func URL() validator.String {
	return &urlStringValidator{}
}
