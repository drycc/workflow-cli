package coder

import (
	"strings"
	"testing"

	"github.com/drycc/controller-sdk-go/api"
)

func TestRouteCoderDecode(t *testing.T) {
	t.Parallel()

	input := `
name = "my-route"
kind = "HTTPRoute"

[[parent_refs]]
name = "my-gateway"
port = 80

[[rules]]

[[rules.backend_refs]]
kind = "Service"
name = "my-svc"
port = 5000
weight = 100
`

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
	if _, ok := c.Request.Rules[0]["backend_refs"]; !ok {
		t.Error("expected backend_refs key in rule, got none")
	}
}

func TestRouteCoderDecodeDottedKeys(t *testing.T) {
	t.Parallel()

	// path 用点号键写法
	input := `
name = "my-route"
kind = "HTTPRoute"

[[parent_refs]]
name = "my-gateway"
port = 80

[[rules]]

[[rules.matches]]
path.type = "PathPrefix"
path.value = "/"

[[rules.backend_refs]]
name = "my-svc"
port = 5000
`

	c := &RouteCoder{}
	if err := c.Decode([]byte(input)); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if len(c.Request.Rules) != 1 {
		t.Fatalf("expected 1 Rule, got %d", len(c.Request.Rules))
	}
	matches, ok := c.Request.Rules[0]["matches"].([]any)
	if !ok || len(matches) != 1 {
		t.Fatalf("expected 1 match, got %v", c.Request.Rules[0]["matches"])
	}
	path := matches[0].(map[string]any)["path"].(map[string]any)
	if path["type"] != "PathPrefix" || path["value"] != "/" {
		t.Errorf("unexpected path: %v", path)
	}
}

func TestRouteCoderDecodeMinimal(t *testing.T) {
	t.Parallel()

	input := "name = \"simple\"\nkind = \"HTTPRoute\"\n"

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

func TestRouteCoderDecodeInvalid(t *testing.T) {
	t.Parallel()

	c := &RouteCoder{}
	if err := c.Decode([]byte(`{invalid`)); err == nil {
		t.Error("expected error for invalid TOML, got nil")
	}
}

func TestRouteCoderDecodeMissingRulesHeader(t *testing.T) {
	t.Parallel()

	// 忘记 [[rules]] 表头时 rules 会变成对象而不是数组，必须报错
	input := `
name = "my-route"
kind = "HTTPRoute"

[[rules.matches]]
path.type = "PathPrefix"
path.value = "/"
`

	c := &RouteCoder{}
	if err := c.Decode([]byte(input)); err == nil {
		t.Error("expected error when [[rules]] header is missing, got nil")
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
					"backend_refs": []map[string]any{
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

	toml := string(data)

	if !strings.Contains(toml, `kind = "HTTPRoute"`) {
		t.Errorf("expected kind in output, got:\n%s", toml)
	}
	if !strings.Contains(toml, `name = "my-route"`) {
		t.Errorf("expected name in output, got:\n%s", toml)
	}
	if !strings.Contains(toml, "[[parent_refs]]") {
		t.Errorf("expected [[parent_refs]] in output, got:\n%s", toml)
	}
	if !strings.Contains(toml, "[[rules.backend_refs]]") {
		t.Errorf("expected [[rules.backend_refs]] in output, got:\n%s", toml)
	}
	if !strings.Contains(toml, "port = 5000") {
		t.Errorf("expected port = 5000 in output, got:\n%s", toml)
	}
	if !strings.Contains(toml, "routable = true") {
		t.Errorf("expected routable = true in output, got:\n%s", toml)
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

	if strings.Contains(string(data), "routable") {
		t.Errorf("routable should not appear when nil, got:\n%s", string(data))
	}
}

func TestRouteCoderRoundTrip(t *testing.T) {
	t.Parallel()

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
					"backend_refs": []map[string]any{
						{"name": "svc", "port": 8080},
					},
				},
			},
			Routable: &routable,
		},
	}

	encoded, err := original.Encode()
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// info output feeds back into apply directly
	decoded := &RouteCoder{}
	if err := decoded.Decode(encoded); err != nil {
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
	if _, ok := decoded.Request.Rules[0]["backend_refs"]; !ok {
		t.Error("expected backend_refs in rule after round trip")
	}
}
