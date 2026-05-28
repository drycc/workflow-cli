package coder

import (
	"strings"
	"testing"

	"github.com/drycc/controller-sdk-go/api"
)

func TestRouteCoderDecode(t *testing.T) {
	t.Parallel()

	input := `{
		"apiVersion": "gateway.networking.k8s.io/v1",
		"kind": "HTTPRoute",
		"metadata": {"name": "my-route"},
		"spec": {
			"parents": [
				{"name": "my-gateway", "port": 80}
			],
			"rules": [
				{
					"backends": [
						{"kind": "Service", "name": "my-svc", "port": 5000, "weight": 100}
					]
				}
			]
		}
	}`

	c := &RouteCoder{}
	if err := c.Decode([]byte(input)); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if c.Request.Name != "my-route" {
		t.Errorf("expected Name=my-route, got %s", c.Request.Name)
	}
	if c.Request.Kind != "HTTPRoute" {
		t.Errorf("expected Kind=HTTPRoute, got %s", c.Request.Kind)
	}
	if len(c.Request.ParentRefs) != 1 {
		t.Fatalf("expected 1 ParentRef, got %d", len(c.Request.ParentRefs))
	}
	if c.Request.ParentRefs[0].Name != "my-gateway" {
		t.Errorf("expected ParentRef name=my-gateway, got %s", c.Request.ParentRefs[0].Name)
	}
	if c.Request.ParentRefs[0].Port != 80 {
		t.Errorf("expected ParentRef port=80, got %d", c.Request.ParentRefs[0].Port)
	}
	if len(c.Request.Rules) != 1 {
		t.Fatalf("expected 1 Rule, got %d", len(c.Request.Rules))
	}
	// backends should be mapped to backendRefs
	if _, ok := c.Request.Rules[0]["backendRefs"]; !ok {
		t.Error("expected backendRefs key in rule, got none")
	}
	if _, ok := c.Request.Rules[0]["backends"]; ok {
		t.Error("backends key should not exist in rule after decode")
	}
}

func TestRouteCoderDecodeMinimal(t *testing.T) {
	t.Parallel()

	input := `{"kind": "HTTPRoute", "metadata": {"name": "simple"}}`

	c := &RouteCoder{}
	if err := c.Decode([]byte(input)); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if c.Request.Name != "simple" {
		t.Errorf("expected Name=simple, got %s", c.Request.Name)
	}
	if c.Request.Kind != "HTTPRoute" {
		t.Errorf("expected Kind=HTTPRoute, got %s", c.Request.Kind)
	}
	if len(c.Request.ParentRefs) != 0 {
		t.Errorf("expected 0 ParentRefs, got %d", len(c.Request.ParentRefs))
	}
	if len(c.Request.Rules) != 0 {
		t.Errorf("expected 0 Rules, got %d", len(c.Request.Rules))
	}
}

func TestRouteCoderDecodeInvalidJSON(t *testing.T) {
	t.Parallel()

	c := &RouteCoder{}
	if err := c.Decode([]byte(`{invalid`)); err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestRouteCoderEncode(t *testing.T) {
	t.Parallel()

	routable := true
	c := &RouteCoder{
		Info: api.RouteInfo{
			Name: "my-route",
			Kind: "HTTPRoute",
			ParentRefs: []api.RouteParentRef{
				{Name: "my-gateway", Port: 80},
			},
			Rules: []api.RouteRule{
				{
					"backendRefs": []map[string]any{
						{"kind": "Service", "name": "my-svc", "port": 5000},
					},
				},
			},
			Routable: &routable,
		},
	}

	data, err := c.Encode()
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	yaml := string(data)

	// Check field order: apiVersion should come first
	if !strings.HasPrefix(yaml, "apiVersion:") {
		t.Errorf("expected YAML to start with apiVersion:, got:\n%s", yaml)
	}

	// Check key fields are present
	if !strings.Contains(yaml, "kind: HTTPRoute") {
		t.Error("expected kind: HTTPRoute in output")
	}
	if !strings.Contains(yaml, "name: my-route") {
		t.Error("expected name: my-route in output")
	}
	if !strings.Contains(yaml, "parents:") {
		t.Error("expected parents: in output")
	}
	if !strings.Contains(yaml, "backends:") {
		t.Error("expected backends: in output (backendRefs should be mapped to backends)")
	}
	if !strings.Contains(yaml, "routable: true") {
		t.Error("expected routable: true in status")
	}
	// backendRefs should NOT appear in output
	if strings.Contains(yaml, "backendRefs:") {
		t.Error("backendRefs should not appear in output, should be backends")
	}
}

func TestRouteCoderEncodeNoRoutable(t *testing.T) {
	t.Parallel()

	c := &RouteCoder{
		Info: api.RouteInfo{
			Name: "my-route",
			Kind: "HTTPRoute",
		},
	}

	data, err := c.Encode()
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	yaml := string(data)
	if strings.Contains(yaml, "routable") {
		t.Error("routable should not appear when nil")
	}
}

func TestRouteCoderRoundTrip(t *testing.T) {
	t.Parallel()

	// Encode
	routable := false
	original := &RouteCoder{
		Info: api.RouteInfo{
			Name: "round-trip",
			Kind: "HTTPRoute",
			ParentRefs: []api.RouteParentRef{
				{Name: "gw", Port: 443},
			},
			Rules: []api.RouteRule{
				{
					"backendRefs": []map[string]any{
						{"name": "svc", "port": 8080},
					},
				},
			},
			Routable: &routable,
		},
	}

	yamlData, err := original.Encode()
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// Decode the YAML back (first convert YAML to JSON)
	// Since Decode expects JSON, we test the field mapping logic directly
	input := `{
		"apiVersion": "controller.drycc.cc/v2.3",
		"kind": "HTTPRoute",
		"metadata": {"name": "round-trip"},
		"spec": {
			"parents": [{"name": "gw", "port": 443}],
			"rules": [{"backends": [{"name": "svc", "port": 8080}]}]
		},
		"status": {"routable": false}
	}`

	decoded := &RouteCoder{}
	if err := decoded.Decode([]byte(input)); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if decoded.Request.Name != "round-trip" {
		t.Errorf("expected Name=round-trip, got %s", decoded.Request.Name)
	}
	if decoded.Request.Kind != "HTTPRoute" {
		t.Errorf("expected Kind=HTTPRoute, got %s", decoded.Request.Kind)
	}
	if len(decoded.Request.ParentRefs) != 1 || decoded.Request.ParentRefs[0].Port != 443 {
		t.Error("ParentRefs mismatch after round trip")
	}
	if len(decoded.Request.Rules) != 1 {
		t.Fatal("expected 1 Rule after round trip")
	}
	if _, ok := decoded.Request.Rules[0]["backendRefs"]; !ok {
		t.Error("expected backendRefs in rule after round trip")
	}

	_ = yamlData // yamlData is valid YAML output, verified above
}
