// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"maps"
	"math/big"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

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

const policyVersionTypeName = "illumio-cloudsecure_policy_version"

// storedPolicyVersion is a policy version held by policyVersionConfigServer.
type storedPolicyVersion struct {
	request       *configv1.CreatePolicyVersionRequest
	versionNumber int64
}

// policyVersionConfigServer is an in-memory Config API that stores policy versions as sent,
// so that every value read back has gone over the wire in both directions.
type policyVersionConfigServer struct {
	configv1.UnimplementedConfigServiceServer

	mu         sync.Mutex
	versions   map[string]storedPolicyVersion
	lastCreate *configv1.CreatePolicyVersionRequest
	created    int64
}

func cloneCreateRequest(req *configv1.CreatePolicyVersionRequest) *configv1.CreatePolicyVersionRequest {
	clone, ok := proto.Clone(req).(*configv1.CreatePolicyVersionRequest)
	if !ok {
		panic("proto.Clone returned an unexpected type")
	}

	return clone
}

func (s *policyVersionConfigServer) CreatePolicyVersion(_ context.Context, req *configv1.CreatePolicyVersionRequest) (*configv1.CreatePolicyVersionResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.created++
	id := uuid.NewString()
	stored := storedPolicyVersion{request: cloneCreateRequest(req), versionNumber: s.created}
	s.versions[id] = stored
	s.lastCreate = cloneCreateRequest(req)

	return &configv1.CreatePolicyVersionResponse{
		Id:            id,
		Description:   stored.request.Description,
		PolicyId:      stored.request.GetPolicyId(),
		Rules:         stored.request.GetRules(),
		VersionNumber: stored.versionNumber,
	}, nil
}

func (s *policyVersionConfigServer) ReadPolicyVersion(_ context.Context, req *configv1.ReadPolicyVersionRequest) (*configv1.ReadPolicyVersionResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stored, ok := s.versions[req.GetId()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "no policy_version with id %s", req.GetId())
	}

	return &configv1.ReadPolicyVersionResponse{
		Id:            req.GetId(),
		Description:   stored.request.Description,
		PolicyId:      stored.request.GetPolicyId(),
		Rules:         stored.request.GetRules(),
		VersionNumber: stored.versionNumber,
	}, nil
}

func (s *policyVersionConfigServer) DeletePolicyVersion(_ context.Context, req *configv1.DeletePolicyVersionRequest) (*emptypb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.versions[req.GetId()]; !ok {
		return nil, status.Errorf(codes.NotFound, "no policy_version with id %s", req.GetId())
	}

	delete(s.versions, req.GetId())

	return &emptypb.Empty{}, nil
}

// lastCreatedRules returns the rules of the most recent create request received.
func (s *policyVersionConfigServer) lastCreatedRules() map[string]*configv1.PolicyVersion_Rules {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.lastCreate.GetRules()
}

// selfSignedTLSConfig returns a server TLS config; the provider is run with insecure_tls.
func selfSignedTLSConfig() (*tls.Config, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "localhost"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}

	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{{Certificate: [][]byte{der}, PrivateKey: key}},
	}, nil
}

// toTerraformValue builds a tftypes.Value of type typ from plain Go values.
// Object attributes missing from a map are set to null, as Terraform does for omitted arguments.
func toTerraformValue(typ tftypes.Type, v any) (tftypes.Value, error) {
	switch v := v.(type) {
	case nil:
		return tftypes.NewValue(typ, nil), nil
	case tftypes.Value:
		return v, nil
	}

	switch t := typ.(type) {
	case tftypes.Object:
		return toTerraformObject(t, v)
	case tftypes.Map:
		return toTerraformMap(t, v)
	case tftypes.List:
		return toTerraformList(t, v)
	default:
		return tftypes.NewValue(typ, v), nil
	}
}

func toTerraformObject(typ tftypes.Object, v any) (tftypes.Value, error) {
	m, ok := v.(map[string]any)
	if !ok {
		return tftypes.Value{}, fmt.Errorf("expected map for object, got %T", v)
	}

	for name := range m {
		if _, ok := typ.AttributeTypes[name]; !ok {
			return tftypes.Value{}, fmt.Errorf("unsupported argument %q", name)
		}
	}

	attrs := make(map[string]tftypes.Value, len(typ.AttributeTypes))

	for name, attrType := range typ.AttributeTypes {
		attr, err := toTerraformValue(attrType, m[name])
		if err != nil {
			return tftypes.Value{}, fmt.Errorf("%s: %w", name, err)
		}

		attrs[name] = attr
	}

	return tftypes.NewValue(typ, attrs), nil
}

func toTerraformMap(typ tftypes.Map, v any) (tftypes.Value, error) {
	m, ok := v.(map[string]any)
	if !ok {
		return tftypes.Value{}, fmt.Errorf("expected map for map, got %T", v)
	}

	elems := make(map[string]tftypes.Value, len(m))

	for key, elem := range m {
		ev, err := toTerraformValue(typ.ElementType, elem)
		if err != nil {
			return tftypes.Value{}, fmt.Errorf("[%q]: %w", key, err)
		}

		elems[key] = ev
	}

	return tftypes.NewValue(typ, elems), nil
}

func toTerraformList(typ tftypes.List, v any) (tftypes.Value, error) {
	l, ok := v.([]any)
	if !ok {
		return tftypes.Value{}, fmt.Errorf("expected slice for list, got %T", v)
	}

	elems := make([]tftypes.Value, 0, len(l))

	for i, elem := range l {
		ev, err := toTerraformValue(typ.ElementType, elem)
		if err != nil {
			return tftypes.Value{}, fmt.Errorf("[%d]: %w", i, err)
		}

		elems = append(elems, ev)
	}

	return tftypes.NewValue(typ, elems), nil
}

// childValues returns the attributes or elements of a known, non-null object, map or list,
// keyed by path step. The boolean is false for primitive values, which have no children.
func childValues(v tftypes.Value) (map[tftypes.AttributePathStep]tftypes.Value, bool, error) {
	children := make(map[tftypes.AttributePathStep]tftypes.Value)

	switch v.Type().(type) {
	case tftypes.Object, tftypes.Map:
		var m map[string]tftypes.Value

		err := v.As(&m)
		if err != nil {
			return nil, false, err
		}

		_, isObject := v.Type().(tftypes.Object)

		for key, child := range m {
			if isObject {
				children[tftypes.AttributeName(key)] = child
			} else {
				children[tftypes.ElementKeyString(key)] = child
			}
		}
	case tftypes.List:
		var l []tftypes.Value

		err := v.As(&l)
		if err != nil {
			return nil, false, err
		}

		for i, child := range l {
			children[tftypes.ElementKeyInt(i)] = child
		}
	default:
		return nil, false, nil
	}

	return children, true, nil
}

// assertAppliedMatchesPlan mirrors Terraform core's check after apply: every value known in the
// plan must be returned unchanged, and every unknown value must be resolved.
func assertAppliedMatchesPlan(path *tftypes.AttributePath, planned, applied tftypes.Value) error {
	if !planned.IsKnown() {
		if !applied.IsFullyKnown() {
			return fmt.Errorf("%s: still unknown after apply", path)
		}

		return nil
	}

	if planned.IsNull() || applied.IsNull() {
		if planned.IsNull() != applied.IsNull() {
			return fmt.Errorf("%s: planned %s, applied %s", path, planned, applied)
		}

		return nil
	}

	plannedChildren, composite, err := childValues(planned)
	if err != nil {
		return err
	}

	if !composite {
		if !planned.Equal(applied) {
			return fmt.Errorf("%s: planned %s, applied %s", path, planned, applied)
		}

		return nil
	}

	appliedChildren, _, err := childValues(applied)
	if err != nil {
		return err
	}

	if len(plannedChildren) != len(appliedChildren) {
		return fmt.Errorf("%s: planned %d elements, applied %d", path, len(plannedChildren), len(appliedChildren))
	}

	for step, plannedChild := range plannedChildren {
		appliedChild, ok := appliedChildren[step]
		if !ok {
			return fmt.Errorf("%s: %v missing after apply", path, step)
		}

		childPath := tftypes.NewAttributePathWithSteps(append(path.Steps(), step))

		err := assertAppliedMatchesPlan(childPath, plannedChild, appliedChild)
		if err != nil {
			return err
		}
	}

	return nil
}

// k8sSource returns a source.k8s selector for one cluster, with extra attributes such as
// policy_apis and tier merged in.
func k8sSource(extra map[string]any) map[string]any {
	return map[string]any{"k8s": k8sSelector(extra)}
}

// k8sDestination returns a destination.k8s selector for the pods of one cluster, with extra
// attributes merged in.
func k8sDestination(extra map[string]any) map[string]any {
	return map[string]any{"k8s": k8sSelector(extra)}
}

func k8sSelector(extra map[string]any) map[string]any {
	k8s := map[string]any{
		"clusters":           []any{map[string]any{"id": "7f8f6b1e-5f0e-4d8c-9a64-3a8f2b6f1c11"}},
		"namespace_selector": map[string]any{},
		"workload_selector":  map[string]any{},
	}
	maps.Copy(k8s, extra)

	return k8s
}

func anyDestination() map[string]any {
	return map[string]any{"any": map[string]any{}}
}

func allTCP() []any {
	return []any{map[string]any{"protocol": "TCP", "from_port": 1, "to_port": 65535}}
}

// allPortsTCPAndUDP is all TCP and UDP ports in the 0-0 form Cilium Baseline rules require.
func allPortsTCPAndUDP() []any {
	return []any{
		map[string]any{"protocol": "TCP", "from_port": 0, "to_port": 0},
		map[string]any{"protocol": "UDP", "from_port": 0, "to_port": 0},
	}
}

func rule(action string, source, destination map[string]any, portRanges []any) map[string]any {
	return map[string]any{
		"action":      action,
		"source":      source,
		"destination": destination,
		"port_ranges": portRanges,
	}
}

func policyVersionConfig(rules map[string]any) map[string]any {
	return map[string]any{
		"policy_id":   "0b3f6a52-2c4e-4d1a-8f43-5e7a1d3b9c20",
		"description": "default deny",
		"rules":       rules,
	}
}

func errorSummaries(diags []*tfprotov6.Diagnostic) string {
	out := make([]string, 0, len(diags))

	for _, d := range diags {
		if d.Severity == tfprotov6.DiagnosticSeverityError {
			out = append(out, d.Summary+": "+d.Detail)
		}
	}

	return strings.Join(out, "\n")
}

// PolicyVersionResourceTestSuite drives the provider through the same protocol calls that
// Terraform makes (validate, plan, apply, read) against an in-memory Config API over TLS gRPC.
type PolicyVersionResourceTestSuite struct {
	suite.Suite

	server       tfprotov6.ProviderServer
	config       *policyVersionConfigServer
	grpcServer   *grpc.Server
	resourceType tftypes.Type
}

func TestPolicyVersionResourceTestSuite(t *testing.T) {
	suite.Run(t, new(PolicyVersionResourceTestSuite))
}

func (s *PolicyVersionResourceTestSuite) SetupTest() {
	ctx := s.T().Context()

	tlsConfig, err := selfSignedTLSConfig()
	s.Require().NoError(err)

	listener, err := (&net.ListenConfig{}).Listen(ctx, "tcp", "127.0.0.1:0")
	s.Require().NoError(err)

	s.config = &policyVersionConfigServer{versions: make(map[string]storedPolicyVersion)}
	s.grpcServer = grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)))
	configv1.RegisterConfigServiceServer(s.grpcServer, s.config)

	go func() { _ = s.grpcServer.Serve(listener) }()

	server, err := providerserver.NewProtocol6WithError(New("test")())()
	s.Require().NoError(err)

	s.server = server

	schemaResp, err := s.server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	s.Require().NoError(err)
	s.requireNoErrorDiags(schemaResp.Diagnostics)

	s.resourceType = schemaResp.ResourceSchemas[policyVersionTypeName].ValueType()

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
	s.requireNoErrorDiags(configureResp.Diagnostics)
}

func (s *PolicyVersionResourceTestSuite) TearDownTest() {
	s.grpcServer.Stop()
}

func (s *PolicyVersionResourceTestSuite) TestCreateReadRoundTripsNewFields() {
	config := s.resourceValue(policyVersionConfig(map[string]any{
		"baseline-deny": rule("Deny",
			k8sSource(map[string]any{"policy_apis": []any{"cilium", "aws_vpc"}, "tier": "Baseline"}),
			anyDestination(),
			allPortsTCPAndUDP()),
		"k8s-to-k8s": rule("Allow",
			k8sSource(map[string]any{"policy_apis": []any{"aws_vpc"}, "tier": "Admin", "priority": 100}),
			k8sDestination(map[string]any{"policy_apis": []any{"cilium"}, "tier": "Admin"}),
			allTCP()),
	}))

	state := s.create(config)

	sent := s.config.lastCreatedRules()

	baseline := sent["baseline-deny"]
	s.Equal([]string{"cilium", "aws_vpc"}, baseline.GetSource().GetK8S().GetPolicyApis())
	s.Equal("Baseline", baseline.GetSource().GetK8S().GetTier())
	s.NotNil(baseline.GetDestination().GetAny(), "destination.any must be sent as a present, empty message")
	s.Nil(baseline.GetDestination().GetK8S())

	k8sToK8s := sent["k8s-to-k8s"]
	s.Equal([]string{"aws_vpc"}, k8sToK8s.GetSource().GetK8S().GetPolicyApis())
	s.Equal(int64(100), k8sToK8s.GetSource().GetK8S().GetPriority())
	s.Equal([]string{"cilium"}, k8sToK8s.GetDestination().GetK8S().GetPolicyApis())
	s.Equal("Admin", k8sToK8s.GetDestination().GetK8S().GetTier())
	s.Nil(k8sToK8s.GetDestination().GetAny())

	// An omitted priority is stored as 0, and policy_apis keeps its order.
	sourceK8s := tftypes.NewAttributePath().
		WithAttributeName("rules").WithElementKeyString("baseline-deny").
		WithAttributeName("source").WithAttributeName("k8s")
	s.True(tftypes.NewValue(tftypes.Number, 0).Equal(s.walk(state, sourceK8s.WithAttributeName("priority"))))
	s.True(
		tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, []tftypes.Value{
			tftypes.NewValue(tftypes.String, "cilium"),
			tftypes.NewValue(tftypes.String, "aws_vpc"),
		}).Equal(s.walk(state, sourceK8s.WithAttributeName("policy_apis"))))
}

func (s *PolicyVersionResourceTestSuite) TestCreateReadWithOmittedFields() {
	config := s.resourceValue(policyVersionConfig(map[string]any{
		"plain": rule("Allow", k8sSource(nil), k8sDestination(nil), allTCP()),
	}))

	state := s.create(config)

	sent := s.config.lastCreatedRules()["plain"]
	s.Nil(sent.GetSource().GetK8S().GetPolicyApis())
	s.Empty(sent.GetSource().GetK8S().GetTier(), "an omitted tier is sent empty, meaning no tier")
	s.Zero(sent.GetSource().GetK8S().GetPriority())
	s.Nil(sent.GetDestination().GetAny())

	destinationK8s := tftypes.NewAttributePath().
		WithAttributeName("rules").WithElementKeyString("plain").
		WithAttributeName("destination").WithAttributeName("k8s")
	s.True(s.walk(state, destinationK8s.WithAttributeName("policy_apis")).IsNull(), "an omitted policy_apis stays null")
	s.True(tftypes.NewValue(tftypes.String, "").Equal(s.walk(state, destinationK8s.WithAttributeName("tier"))))
}

func (s *PolicyVersionResourceTestSuite) TestUpdateReplacesVersion() {
	before := policyVersionConfig(map[string]any{
		"rule": rule("Deny",
			k8sSource(map[string]any{"tier": "Admin"}),
			map[string]any{"fqdns": map[string]any{"names": []any{"api.example.com"}}},
			allTCP()),
	})
	prior := s.create(s.resourceValue(before))

	var priorAttrs map[string]tftypes.Value
	s.Require().NoError(prior.As(&priorAttrs))

	after := policyVersionConfig(map[string]any{
		"rule": rule("Deny",
			k8sSource(map[string]any{"policy_apis": []any{"cilium"}, "tier": "Baseline"}),
			anyDestination(),
			allPortsTCPAndUDP()),
	})
	config := s.resourceValue(after)
	s.requireNoErrorDiags(s.validate(config))

	// Terraform proposes the config, with computed attributes carried over from prior state.
	after["id"] = priorAttrs["id"]
	after["version_number"] = priorAttrs["version_number"]
	proposed := s.resourceValue(after)

	planResp := s.plan(prior, proposed, config)
	s.Contains(planResp.RequiresReplace, tftypes.NewAttributePath().WithAttributeName("rules"),
		"policy versions are immutable, so changing a rule replaces the version")

	// A replacement is a delete followed by a create.
	null := tftypes.NewValue(s.resourceType, nil)
	s.apply(prior, null, null)
	s.create(config)

	sent := s.config.lastCreatedRules()["rule"]
	s.Equal([]string{"cilium"}, sent.GetSource().GetK8S().GetPolicyApis())
	s.Equal("Baseline", sent.GetSource().GetK8S().GetTier())
	s.NotNil(sent.GetDestination().GetAny())
	s.Nil(sent.GetDestination().GetFqdns())
}

// TestRemovingTierAndPriorityReplacesVersion checks that removing tier or priority from config is
// a change. Terraform proposes the prior value for an omitted Optional+Computed attribute, so
// without a default the plan would show no change and the old tier would stay in effect.
func (s *PolicyVersionResourceTestSuite) TestRemovingTierAndPriorityReplacesVersion() {
	withTier := map[string]any{"policy_apis": []any{"aws_vpc"}, "tier": "Baseline", "priority": 5}
	prior := s.create(s.resourceValue(policyVersionConfig(map[string]any{
		"rule": rule("Allow", k8sSource(withTier), k8sDestination(nil), allTCP()),
	})))

	var priorAttrs map[string]tftypes.Value
	s.Require().NoError(prior.As(&priorAttrs))

	config := s.resourceValue(policyVersionConfig(map[string]any{
		"rule": rule("Allow", k8sSource(nil), k8sDestination(nil), allTCP()),
	}))

	proposedRules := policyVersionConfig(map[string]any{
		"rule": rule("Allow", k8sSource(map[string]any{"tier": "Baseline", "priority": 5}), k8sDestination(nil), allTCP()),
	})
	proposedRules["id"] = priorAttrs["id"]
	proposedRules["version_number"] = priorAttrs["version_number"]

	planResp := s.plan(prior, s.resourceValue(proposedRules), config)
	s.Contains(planResp.RequiresReplace, tftypes.NewAttributePath().WithAttributeName("rules"))

	sourceK8s := tftypes.NewAttributePath().
		WithAttributeName("rules").WithElementKeyString("rule").
		WithAttributeName("source").WithAttributeName("k8s")
	planned := s.value(planResp.PlannedState)
	s.True(tftypes.NewValue(tftypes.String, "").Equal(s.walk(planned, sourceK8s.WithAttributeName("tier"))))
	s.True(tftypes.NewValue(tftypes.Number, 0).Equal(s.walk(planned, sourceK8s.WithAttributeName("priority"))))
}

// TestRuleValidatorIsWired checks that the rule validator is attached to the schema. Its rules
// are tested directly in api/schema.
func (s *PolicyVersionResourceTestSuite) TestRuleValidatorIsWired() {
	diags := s.validate(s.resourceValue(policyVersionConfig(map[string]any{
		"r": rule("Allow",
			k8sSource(map[string]any{"policy_apis": []any{"cilium"}, "tier": "Baseline"}),
			anyDestination(), allPortsTCPAndUDP()),
	})))
	s.Contains(errorSummaries(diags), `Rule "r", source.k8s: policy API "cilium" supports tier "Baseline" only on source.k8s`)
}

func (s *PolicyVersionResourceTestSuite) walk(v tftypes.Value, p *tftypes.AttributePath) tftypes.Value {
	found, _, err := tftypes.WalkAttributePath(v, p)
	s.Require().NoError(err)

	value, ok := found.(tftypes.Value)
	s.Require().True(ok)

	return value
}

func (s *PolicyVersionResourceTestSuite) dynamicValue(v tftypes.Value) *tfprotov6.DynamicValue {
	dv, err := tfprotov6.NewDynamicValue(v.Type(), v)
	s.Require().NoError(err)

	return &dv
}

func (s *PolicyVersionResourceTestSuite) value(dv *tfprotov6.DynamicValue) tftypes.Value {
	v, err := dv.Unmarshal(s.resourceType)
	s.Require().NoError(err)

	return v
}

func (s *PolicyVersionResourceTestSuite) resourceValue(v map[string]any) tftypes.Value {
	tv, err := toTerraformValue(s.resourceType, v)
	s.Require().NoError(err)

	return tv
}

func (s *PolicyVersionResourceTestSuite) requireNoErrorDiags(diags []*tfprotov6.Diagnostic) {
	s.Require().Empty(errorSummaries(diags), "unexpected error diagnostics")
}

func (s *PolicyVersionResourceTestSuite) validate(config tftypes.Value) []*tfprotov6.Diagnostic {
	resp, err := s.server.ValidateResourceConfig(s.T().Context(), &tfprotov6.ValidateResourceConfigRequest{
		TypeName: policyVersionTypeName,
		Config:   s.dynamicValue(config),
	})
	s.Require().NoError(err)

	return resp.Diagnostics
}

func (s *PolicyVersionResourceTestSuite) plan(prior, proposed, config tftypes.Value) *tfprotov6.PlanResourceChangeResponse {
	resp, err := s.server.PlanResourceChange(s.T().Context(), &tfprotov6.PlanResourceChangeRequest{
		TypeName:         policyVersionTypeName,
		PriorState:       s.dynamicValue(prior),
		ProposedNewState: s.dynamicValue(proposed),
		Config:           s.dynamicValue(config),
	})
	s.Require().NoError(err)
	s.requireNoErrorDiags(resp.Diagnostics)

	return resp
}

func (s *PolicyVersionResourceTestSuite) apply(prior, planned, config tftypes.Value) tftypes.Value {
	resp, err := s.server.ApplyResourceChange(s.T().Context(), &tfprotov6.ApplyResourceChangeRequest{
		TypeName:     policyVersionTypeName,
		PriorState:   s.dynamicValue(prior),
		PlannedState: s.dynamicValue(planned),
		Config:       s.dynamicValue(config),
	})
	s.Require().NoError(err)
	s.requireNoErrorDiags(resp.Diagnostics)

	return s.value(resp.NewState)
}

func (s *PolicyVersionResourceTestSuite) read(state tftypes.Value) tftypes.Value {
	resp, err := s.server.ReadResource(s.T().Context(), &tfprotov6.ReadResourceRequest{
		TypeName:     policyVersionTypeName,
		CurrentState: s.dynamicValue(state),
	})
	s.Require().NoError(err)
	s.requireNoErrorDiags(resp.Diagnostics)

	return s.value(resp.NewState)
}

// create runs validate, plan and apply for a new resource, checks that the applied state is
// consistent with the plan, and checks that a read returns exactly the applied state.
func (s *PolicyVersionResourceTestSuite) create(config tftypes.Value) tftypes.Value {
	s.requireNoErrorDiags(s.validate(config))

	null := tftypes.NewValue(s.resourceType, nil)
	planned := s.value(s.plan(null, config, config).PlannedState)

	state := s.apply(null, planned, config)
	s.Require().NoError(assertAppliedMatchesPlan(tftypes.NewAttributePath(), planned, state))

	readState := s.read(state)
	s.Require().True(readState.Equal(state), "read returned a different state:\nread:    %s\napplied: %s", readState, state)

	return state
}
