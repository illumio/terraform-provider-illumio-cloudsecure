// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/providervalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	provider_schema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	configv1 "github.com/illumio/terraform-provider-illumio-cloudsecure/api/illumio/cloud/config/v1"
	api_schema "github.com/illumio/terraform-provider-illumio-cloudsecure/api/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// Provider defines the provider implementation.
type Provider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string

	// schema contains the schemas of CloudSecure resources and data sources.
	schema api_schema.Schema
}

var _ provider.ProviderWithConfigValidators = &Provider{}

// ProviderModel describes the provider data model.
type ProviderModel struct {
	APIEndpoint    types.String `tfsdk:"api_endpoint"`
	TokenEndpoint  types.String `tfsdk:"token_endpoint"`
	ClientID       types.String `tfsdk:"client_id"`
	ClientSecret   types.String `tfsdk:"client_secret"`
	AccessToken    types.String `tfsdk:"access_token"`
	RequestTimeout types.String `tfsdk:"request_timeout"`
	DisableTLS     types.Bool   `tfsdk:"disable_tls"`
}

// providerData is the
type providerData struct {
	// client is the CloudSecure Config API client.
	client configv1.ConfigServiceClient

	// accessToken is the OAuth 2 access token to use to authenticate against the CloudSecure Config API.
	accessToken string

	// requestTimeout is the maximum duration of each API request.
	requestTimeout time.Duration
}

var _ ProviderData = &providerData{}

func (d *providerData) Client() configv1.ConfigServiceClient {
	return d.client
}

func (d *providerData) AccessToken() string {
	return d.accessToken
}

func (d *providerData) RequestTimeout() time.Duration {
	return d.requestTimeout
}

func (p *Provider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "illumio-cloudsecure"
	resp.Version = p.version
}

const (
	// DefaultAPIEndpoint is the default CloudSecure Config API endpoint.
	// TODO: This is a placeholder. Update to use the correct endpoint.
	DefaultAPIEndpoint = "dns:///cloud.illum.io:443"

	// DefaultTokenEndpoint is the default CloudSecure OAuth 2 Token endpoint.
	// TODO: This is a placeholder. Update to use the correct endpoint.
	DefaultTokenEndpoint = "https://cloud.illum.io/token"

	// DefaultRequestTimeout is the default CloudConfig Config API request timeout.
	DefaultRequestTimeout = "10s"
)

func (p *Provider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = provider_schema.Schema{
		MarkdownDescription: "A Provider for managing Illumio CloudSecure.",
		Attributes: map[string]provider_schema.Attribute{
			"api_endpoint": provider_schema.StringAttribute{
				MarkdownDescription: "CloudSecure Config API endpoint, defaults to " + DefaultAPIEndpoint + ".",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 1024),
					URL(),
				},
			},
			"token_endpoint": provider_schema.StringAttribute{
				MarkdownDescription: "CloudSecure OAuth 2 Token endpoint, defaults to " + DefaultTokenEndpoint + ".",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 1024),
					URL(),
				},
			},
			"client_id": provider_schema.StringAttribute{
				MarkdownDescription: "OAuth 2 client identifier used to authenticate against the CloudSecure Config API. Either client_id+client_secret or access_token must be specified.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 1024),
				},
			},
			"client_secret": provider_schema.StringAttribute{
				MarkdownDescription: "OAuth 2 client secret used to authenticate against the CloudSecure Config API. Either client_id+client_secret or access_token must be specified.",
				Optional:            true,
				Sensitive:           true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 1024),
				},
			},
			"access_token": provider_schema.StringAttribute{
				MarkdownDescription: "OAuth 2 access token used to authenticate against the CloudSecure Config API. Either client_id+client_secret or access_token must be specified.",
				Optional:            true,
				Sensitive:           true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 1024),
				},
			},
			"request_timeout": provider_schema.StringAttribute{
				MarkdownDescription: "Maximum duration of each API request, defaults to " + DefaultRequestTimeout + ".",
				Optional:            true,
				Validators: []validator.String{
					Duration(),
				},
			},
			"disable_tls": provider_schema.BoolAttribute{
				MarkdownDescription: "Disables TLS for all all requests to the CloudSecure Token and Config API endpoints. TLS is enabled by default. Should only be used for testing the provider.",
				Optional:            true,
			},
		},
	}
}

func (p *Provider) ConfigValidators(ctx context.Context) []provider.ConfigValidator {
	return []provider.ConfigValidator{
		providervalidator.RequiredTogether(
			path.MatchRoot("client_id"),
			path.MatchRoot("client_secret"),
		),
		providervalidator.ExactlyOneOf(
			path.MatchRoot("client_secret"),
			path.MatchRoot("access_token"),
		),
	}
}

func (p *Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.APIEndpoint.IsUnknown() || data.APIEndpoint.IsNull() {
		data.APIEndpoint = types.StringValue(DefaultAPIEndpoint)
	}
	if data.TokenEndpoint.IsUnknown() || data.TokenEndpoint.IsNull() {
		data.TokenEndpoint = types.StringValue(DefaultTokenEndpoint)
	}
	if data.RequestTimeout.IsUnknown() || data.RequestTimeout.IsNull() {
		data.RequestTimeout = types.StringValue(DefaultRequestTimeout)
	}

	requestTimeout, err := time.ParseDuration(data.RequestTimeout.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Invalid Request Timeout",
			fmt.Sprintf("Failed to parse request_timeout: %s. Please report this issue to the provider developers.", data.RequestTimeout.ValueString()),
		)
	}

	var creds credentials.TransportCredentials
	if data.DisableTLS.ValueBool() {
		creds = insecure.NewCredentials()
	} else {
		tlsConfig := &tls.Config{}
		creds = credentials.NewTLS(tlsConfig)
	}

	conn, err := grpc.NewClient(data.APIEndpoint.ValueString(), grpc.WithTransportCredentials(creds))
	if err != nil {
		resp.Diagnostics.AddError("Config API Error", fmt.Sprintf("Unable to create Config API client, got error: %s", err))
	}

	client := configv1.NewConfigServiceClient(conn)

	providerData := &providerData{
		client:         client,
		accessToken:    "", // TODO: Implement the OAuth 2 client credentials flow, and modify the generated resource code to pass this token as metadata in every RPC.
		requestTimeout: requestTimeout,
	}

	resp.DataSourceData = providerData
	resp.ResourceData = providerData
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &Provider{
			version: version,
			schema:  api_schema.CloudSecure(),
		}
	}
}
