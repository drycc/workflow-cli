package coder

import (
	"strings"
	"testing"

	"github.com/drycc/controller-sdk-go/api"
	"sigs.k8s.io/yaml"
)

func TestAddonCoderDecode(t *testing.T) {
	t.Parallel()

	input := `{
		"apiVersion": "controller.drycc.cc/v2",
		"kind": "Valkey",
		"metadata": {"name": "my-addon"},
		"spec": {
			"plan": "micro",
			"parameters": {
				"shards": 1,
				"tls": {"enabled": true}
			}
		}
	}`

	c := &AddonCoder{}
	if err := c.Decode([]byte(input)); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if c.Name != "my-addon" {
		t.Errorf("expected Name=my-addon, got %s", c.Name)
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

	input := `{"kind": "Valkey", "metadata": {"name": "simple"}}`

	c := &AddonCoder{}
	if err := c.Decode([]byte(input)); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if c.Name != "simple" {
		t.Errorf("expected Name=simple, got %s", c.Name)
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

	if !strings.HasPrefix(yaml, "apiVersion:") {
		t.Errorf("expected YAML to start with apiVersion:, got:\n%s", yaml)
	}
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
		t.Error("expected multiplier: 2 in status")
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

	// When Parameters is nil, output should include empty parameters map
	if !strings.Contains(yaml, "parameters: {}") {
		t.Error("expected parameters: {} in output when Parameters is nil")
	}
	// When Multiplier is 0, it should be omitted
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

	yamlData, err := c.Encode()
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// Convert YAML back to JSON for decode (mirrors what commands/ does)
	jsonData, err := yaml.YAMLToJSON(yamlData)
	if err != nil {
		t.Fatalf("YAMLToJSON failed: %v", err)
	}

	decoded := &AddonCoder{}
	if err := decoded.Decode(jsonData); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if decoded.Name != "round-trip" {
		t.Errorf("expected Name=round-trip, got %s", decoded.Name)
	}
	if decoded.Request.Kind != "Valkey" {
		t.Errorf("expected Kind=Valkey, got %s", decoded.Request.Kind)
	}
	if decoded.Request.Plan != "small" {
		t.Errorf("expected Plan=small, got %s", decoded.Request.Plan)
	}
	if decoded.Request.Parameters == nil {
		t.Fatal("expected Parameters to be non-nil")
	}

	_ = yamlData
}
