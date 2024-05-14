// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// durationStringValidator validates that a schema attribute's string value is a valid time duration.
type durationStringValidator struct {
}

func (v *durationStringValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v *durationStringValidator) MarkdownDescription(ctx context.Context) string {
	return "value must be a valid time duration"
}

// ValidateString implements validator.String.
func (v *durationStringValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	_, err := time.ParseDuration(value)
	if err != nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			v.Description(ctx),
			value,
		))
	}
}

var _ validator.String = &durationStringValidator{}

func Duration() validator.String {
	return &durationStringValidator{}
}
