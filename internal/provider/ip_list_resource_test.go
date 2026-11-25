// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net"
	"net/netip"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	configv1 "github.com/illumio/terraform-provider-illumio-cloudsecure/api/illumio/cloud/config/v1"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

const ipListTypeName = "illumio-cloudsecure_ip_list"

// ipListConfigServer is an in-memory Config API that stores IP lists the way dev does:
// addresses are canonicalized, and the description is always returned, empty when unset.
type ipListConfigServer struct {
	configv1.UnimplementedConfigServiceServer

	mu    sync.Mutex
	lists map[string]*configv1.CreateIpListRequest
}

// canonicalize mirrors dev: /32 and /128 are stripped, CIDRs are masked, IPv6 is compressed and lowercased.
func canonicalize(s string) string {
	if strings.Contains(s, "/") {
		prefix, err := netip.ParsePrefix(s)
		if err != nil {
			return s
		}

		if prefix.IsSingleIP() {
			return prefix.Addr().String()
		}

		return prefix.Masked().String()
	}

	addr, err := netip.ParseAddr(s)
	if err != nil {
		return s
	}

	return addr.String()
}

func (s *ipListConfigServer) CreateIpList(_ context.Context, req *configv1.CreateIpListRequest) (*configv1.CreateIpListResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stored, ok := proto.Clone(req).(*configv1.CreateIpListRequest)
	if !ok {
		panic("proto.Clone returned an unexpected type")
	}

	for _, a := range stored.GetIpAddresses() {
		a.IpAddress = canonicalize(a.GetIpAddress())
	}

	for _, r := range stored.GetIpRanges() {
		r.FromIpAddress = canonicalize(r.GetFromIpAddress())
		r.ToIpAddress = canonicalize(r.GetToIpAddress())
	}

	id := uuid.NewString()
	s.lists[id] = stored

	return &configv1.CreateIpListResponse{
		Id: id, Description: proto.String(stored.GetDescription()), IpAddresses: stored.GetIpAddresses(),
		IpRanges: stored.GetIpRanges(), Name: stored.GetName(),
	}, nil
}

func (s *ipListConfigServer) ReadIpList(_ context.Context, req *configv1.ReadIpListRequest) (*configv1.ReadIpListResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stored, ok := s.lists[req.GetId()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "no ip_list with id %s", req.GetId())
	}

	return &configv1.ReadIpListResponse{
		Id: req.GetId(), Description: proto.String(stored.GetDescription()), IpAddresses: stored.GetIpAddresses(),
		IpRanges: stored.GetIpRanges(), Name: stored.GetName(),
	}, nil
}

func (s *ipListConfigServer) DeleteIpList(_ context.Context, req *configv1.DeleteIpListRequest) (*emptypb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.lists, req.GetId())

	return &emptypb.Empty{}, nil
}

// IPListResourceTestSuite drives the ip_list resource through Terraform's protocol calls
// against an in-memory Config API that behaves like dev.
type IPListResourceTestSuite struct {
	suite.Suite

	server       tfprotov6.ProviderServer
	grpcServer   *grpc.Server
	resourceType tftypes.Type
}

func TestIPListResourceTestSuite(t *testing.T) {
	suite.Run(t, new(IPListResourceTestSuite))
}

func (s *IPListResourceTestSuite) SetupTest() {
	ctx := s.T().Context()

	tlsConfig, err := selfSignedTLSConfig()
	s.Require().NoError(err)

	listener, err := (&net.ListenConfig{}).Listen(ctx, "tcp", "127.0.0.1:0")
	s.Require().NoError(err)

	s.grpcServer = grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)))
	configv1.RegisterConfigServiceServer(s.grpcServer, &ipListConfigServer{lists: make(map[string]*configv1.CreateIpListRequest)})

	go func() { _ = s.grpcServer.Serve(listener) }()

	s.server, err = providerserver.NewProtocol6WithError(New("test")())()
	s.Require().NoError(err)

	schemaResp, err := s.server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	s.Require().NoError(err)
	s.Require().Empty(errorSummaries(schemaResp.Diagnostics))

	s.resourceType = schemaResp.ResourceSchemas[ipListTypeName].ValueType()

	providerConfig, err := toTerraformValue(schemaResp.Provider.ValueType(), map[string]any{
		"api_endpoint": "dns:///" + listener.Addr().String(),
		"access_token": "test-token",
		"insecure_tls": true,
	})
	s.Require().NoError(err)

	configureResp, err := s.server.ConfigureProvider(ctx, &tfprotov6.ConfigureProviderRequest{
		TerraformVersion: "1.9.0",
		Config:           s.dynamicValue(providerConfig),
	})
	s.Require().NoError(err)
	s.Require().Empty(errorSummaries(configureResp.Diagnostics))
}

func (s *IPListResourceTestSuite) TearDownTest() {
	s.grpcServer.Stop()
}

// TestCanonicalConfigsApplyCleanly checks that configs in canonical form, with and without a
// description, apply consistently and then plan no changes.
func (s *IPListResourceTestSuite) TestCanonicalConfigsApplyCleanly() {
	for name, config := range map[string]map[string]any{
		"no description": {
			"name":         "no-description",
			"ip_addresses": []any{ipAddress("10.20.96.91")},
		},
		"canonical forms": {
			"name":        "canonical",
			"description": "office",
			"ip_addresses": []any{
				ipAddress("10.20.96.91"), ipAddress("10.20.0.0/16"),
				ipAddress("2001:db8::1"), ipAddress("2001:db8::/64"),
			},
			"ip_ranges": []any{map[string]any{"from_ip_address": "10.0.0.1", "to_ip_address": "10.0.0.9", "exclusion": true}},
		},
	} {
		s.Run(name, func() {
			configValue := s.resourceValue(config)
			s.Require().Empty(errorSummaries(s.validate(configValue)))

			null := tftypes.NewValue(s.resourceType, nil)
			planResp := s.plan(null, configValue, configValue)
			s.Require().Empty(errorSummaries(planResp.Diagnostics))
			planned := s.value(planResp.PlannedState)

			applyResp, err := s.server.ApplyResourceChange(s.T().Context(), &tfprotov6.ApplyResourceChangeRequest{
				TypeName: ipListTypeName, PriorState: s.dynamicValue(null),
				PlannedState: s.dynamicValue(planned), Config: s.dynamicValue(configValue),
			})
			s.Require().NoError(err)
			s.Require().Empty(errorSummaries(applyResp.Diagnostics))

			state := s.value(applyResp.NewState)
			s.Require().NoError(assertAppliedMatchesPlan(tftypes.NewAttributePath(), planned, state))

			readResp, err := s.server.ReadResource(s.T().Context(), &tfprotov6.ReadResourceRequest{
				TypeName: ipListTypeName, CurrentState: s.dynamicValue(state),
			})
			s.Require().NoError(err)
			s.Require().Empty(errorSummaries(readResp.Diagnostics))
			s.Require().True(s.value(readResp.NewState).Equal(state), "read returned a different state")

			// Terraform proposes the config with computed values from state; the plan must be a no-op.
			var stateAttrs map[string]tftypes.Value
			s.Require().NoError(state.As(&stateAttrs))

			proposedConfig := make(map[string]any, len(config)+2)
			for k, v := range config {
				proposedConfig[k] = v
			}

			proposedConfig["id"] = stateAttrs["id"]
			proposedConfig["description"] = stateAttrs["description"]

			replan := s.plan(state, s.resourceValue(proposedConfig), configValue)
			s.Require().Empty(errorSummaries(replan.Diagnostics))
			s.True(s.value(replan.PlannedState).Equal(state), "follow-up plan is not a no-op")
		})
	}
}

// TestRemovingDescriptionClearsIt checks that removing description from config plans "", not the
// prior value, which Terraform proposes for an omitted Optional+Computed attribute.
func (s *IPListResourceTestSuite) TestRemovingDescriptionClearsIt() {
	prior := s.resourceValue(map[string]any{
		"id":           "0b3f6a52-2c4e-4d1a-8f43-5e7a1d3b9c20",
		"name":         "office",
		"description":  "old",
		"ip_addresses": []any{ipAddress("10.20.96.91")},
	})
	config := map[string]any{"name": "office", "ip_addresses": []any{ipAddress("10.20.96.91")}}
	proposed := map[string]any{
		"id":           "0b3f6a52-2c4e-4d1a-8f43-5e7a1d3b9c20",
		"name":         "office",
		"description":  "old",
		"ip_addresses": []any{ipAddress("10.20.96.91")},
	}

	planResp := s.plan(prior, s.resourceValue(proposed), s.resourceValue(config))
	s.Require().Empty(errorSummaries(planResp.Diagnostics))

	description, _, err := tftypes.WalkAttributePath(s.value(planResp.PlannedState), tftypes.NewAttributePath().WithAttributeName("description"))
	s.Require().NoError(err)

	value, ok := description.(tftypes.Value)
	s.Require().True(ok)
	s.True(tftypes.NewValue(tftypes.String, "").Equal(value), "planned description: %v", value)
}

// TestIPListValidatorsAreWired checks that the address validators are attached to the schema.
// Their rules are tested directly in api/schema.
func (s *IPListResourceTestSuite) TestIPListValidatorsAreWired() {
	diags := errorSummaries(s.validate(s.resourceValue(map[string]any{
		"name":         "rejected",
		"ip_addresses": []any{ipAddress("10.20.96.91/32")},
		"ip_ranges":    []any{ipRange("10.0.0.9", "10.0.0.1")},
	})))
	s.Contains(diags, `Write "10.20.96.91" instead of "10.20.96.91/32"`)
	s.Contains(diags, "from_ip_address must not be greater than to_ip_address.")
}

func ipRange(from, to string) map[string]any {
	return map[string]any{"from_ip_address": from, "to_ip_address": to, "exclusion": false}
}

func ipAddress(ip string) map[string]any {
	return map[string]any{"ip_address": ip, "exclusion": false}
}

func (s *IPListResourceTestSuite) dynamicValue(v tftypes.Value) *tfprotov6.DynamicValue {
	dv, err := tfprotov6.NewDynamicValue(v.Type(), v)
	s.Require().NoError(err)

	return &dv
}

func (s *IPListResourceTestSuite) value(dv *tfprotov6.DynamicValue) tftypes.Value {
	v, err := dv.Unmarshal(s.resourceType)
	s.Require().NoError(err)

	return v
}

func (s *IPListResourceTestSuite) resourceValue(v map[string]any) tftypes.Value {
	tv, err := toTerraformValue(s.resourceType, v)
	s.Require().NoError(err)

	return tv
}

func (s *IPListResourceTestSuite) validate(config tftypes.Value) []*tfprotov6.Diagnostic {
	resp, err := s.server.ValidateResourceConfig(s.T().Context(), &tfprotov6.ValidateResourceConfigRequest{
		TypeName: ipListTypeName,
		Config:   s.dynamicValue(config),
	})
	s.Require().NoError(err)

	return resp.Diagnostics
}

func (s *IPListResourceTestSuite) plan(prior, proposed, config tftypes.Value) *tfprotov6.PlanResourceChangeResponse {
	resp, err := s.server.PlanResourceChange(s.T().Context(), &tfprotov6.PlanResourceChangeRequest{
		TypeName:         ipListTypeName,
		PriorState:       s.dynamicValue(prior),
		ProposedNewState: s.dynamicValue(proposed),
		Config:           s.dynamicValue(config),
	})
	s.Require().NoError(err)

	return resp
}
