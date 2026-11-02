package coder

import (
	"strings"
	"testing"

	"github.com/drycc/controller-sdk-go/api"
)

func TestAddonCoderDecode(t *testing.T) {
	t.Parallel()

	input := `
name = "my-addon"
kind = "Valkey"
plan = "micro"

[parameters]
shards = 1

[parameters.tls]
enabled = true
`

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
	tls, ok := c.Request.Parameters["tls"].(map[string]any)
	if !ok || tls["enabled"] != true {
		t.Errorf("expected Parameters[tls][enabled]=true, got %v", c.Request.Parameters["tls"])
	}
}

func TestAddonCoderDecodeMinimal(t *testing.T) {
	t.Parallel()

	input := `kind = "Valkey"`

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

func TestAddonCoderDecodeInvalid(t *testing.T) {
	t.Parallel()

	c := &AddonCoder{}
	if err := c.Decode([]byte(`{invalid`)); err == nil {
		t.Error("expected error for invalid TOML, got nil")
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

	toml := string(data)

	if !strings.Contains(toml, `kind = "Valkey"`) {
		t.Errorf("expected kind in output, got:\n%s", toml)
	}
	if !strings.Contains(toml, `name = "my-addon"`) {
		t.Errorf("expected name in output, got:\n%s", toml)
	}
	if !strings.Contains(toml, `plan = "micro"`) {
		t.Errorf("expected plan in output, got:\n%s", toml)
	}
	if !strings.Contains(toml, "shards = 1") {
		t.Errorf("expected shards in parameters output, got:\n%s", toml)
	}
	if !strings.Contains(toml, "multiplier = 2") {
		t.Errorf("expected multiplier in output, got:\n%s", toml)
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

	if strings.Contains(string(data), "multiplier") {
		t.Errorf("multiplier should be omitted when 0, got:\n%s", string(data))
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

	encoded, err := c.Encode()
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	toml := string(encoded)
	if !strings.Contains(toml, `kind = "Valkey"`) {
		t.Errorf("expected kind in output, got:\n%s", toml)
	}
	if !strings.Contains(toml, `plan = "small"`) {
		t.Errorf("expected plan in output, got:\n%s", toml)
	}
	if !strings.Contains(toml, "replicas = 1") {
		t.Errorf("expected replicas in parameters output, got:\n%s", toml)
	}

	// info output feeds back into apply directly
	decoded := &AddonCoder{}
	if err := decoded.Decode(encoded); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if decoded.Name != "round-trip" || decoded.Request.Kind != "Valkey" || decoded.Request.Plan != "small" {
		t.Errorf("mismatch after round trip: name=%s kind=%s plan=%s",
			decoded.Name, decoded.Request.Kind, decoded.Request.Plan)
	}
}
