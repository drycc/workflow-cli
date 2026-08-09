package coder

import (
	"strings"
	"testing"

	"github.com/drycc/controller-sdk-go/api"
)

func TestAddonCoderDecode(t *testing.T) {
	t.Parallel()

	input := `{
		"kind": "Valkey",
		"plan": "micro",
		"parameters": {
			"shards": 1,
			"tls": {"enabled": true}
		}
	}`

	c := &AddonCoder{}
	if err := c.Decode([]byte(input)); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if c.Request.Kind != "Valkey" {
		t.Errorf("expected Kind=Valkey, got %s", c.Request.Kind)
	}
	if c.Request.Plan != "micro" {
		t.Errorf("expected Plan=micro, got %s", c.Request.Plan)
	}
	if c.Request.Parameters == nil {
		t.Fatal("expected Parameters to be non-nil")
	}
	if v, ok := c.Request.Parameters["shards"].(float64); !ok || v != 1 {
		t.Errorf("expected Parameters[shards]=1, got %v", c.Request.Parameters["shards"])
	}
}

func TestAddonCoderDecodeMinimal(t *testing.T) {
	t.Parallel()

	input := `{"kind": "Valkey"}`

	c := &AddonCoder{}
	if err := c.Decode([]byte(input)); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if c.Request.Kind != "Valkey" {
		t.Errorf("expected Kind=Valkey, got %s", c.Request.Kind)
	}
	if c.Request.Plan != "" {
		t.Errorf("expected empty Plan, got %s", c.Request.Plan)
	}
	if c.Request.Parameters != nil {
		t.Errorf("expected nil Parameters, got %v", c.Request.Parameters)
	}
}

func TestAddonCoderDecodeInvalidJSON(t *testing.T) {
	t.Parallel()

	c := &AddonCoder{}
	if err := c.Decode([]byte(`{invalid`)); err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestAddonCoderEncode(t *testing.T) {
	t.Parallel()

	c := &AddonCoder{
		Info: api.AddonInstance{
			Name:       "my-addon",
			Kind:       "Valkey",
			Plan:       "micro",
			Multiplier: 2,
			Parameters: map[string]any{
				"shards": 1,
			},
		},
	}

	data, err := c.Encode()
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	yaml := string(data)

	if !strings.Contains(yaml, "kind: Valkey") {
		t.Error("expected kind: Valkey in output")
	}
	if !strings.Contains(yaml, "name: my-addon") {
		t.Error("expected name: my-addon in output")
	}
	if !strings.Contains(yaml, "plan: micro") {
		t.Error("expected plan: micro in output")
	}
	if !strings.Contains(yaml, "shards:") {
		t.Error("expected shards: in parameters output")
	}
	if !strings.Contains(yaml, "multiplier: 2") {
		t.Error("expected multiplier: 2 in output")
	}
}

func TestAddonCoderEncodeEmpty(t *testing.T) {
	t.Parallel()

	c := &AddonCoder{
		Info: api.AddonInstance{
			Name: "empty-addon",
			Kind: "Valkey",
			Plan: "micro",
		},
	}

	data, err := c.Encode()
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	yaml := string(data)

	if strings.Contains(yaml, "multiplier") {
		t.Error("multiplier should be omitted when 0")
	}
}

func TestAddonCoderRoundTrip(t *testing.T) {
	t.Parallel()

	c := &AddonCoder{
		Info: api.AddonInstance{
			Name:       "round-trip",
			Kind:       "Valkey",
			Plan:       "small",
			Multiplier: 1,
			Parameters: map[string]any{"replicas": 1},
		},
	}

	data, err := c.Encode()
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	yaml := string(data)

	if !strings.Contains(yaml, "kind: Valkey") {
		t.Error("expected kind: Valkey in output")
	}
	if !strings.Contains(yaml, "plan: small") {
		t.Error("expected plan: small in output")
	}
	if !strings.Contains(yaml, "replicas:") {
		t.Error("expected replicas: in parameters output")
	}
}
