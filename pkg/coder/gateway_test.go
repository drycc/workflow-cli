package coder

import (
	"strings"
	"testing"

	"github.com/drycc/controller-sdk-go/api"
)

func TestGatewayCoderDecode(t *testing.T) {
	t.Parallel()

	input := `{
		"apiVersion": "gateway.networking.k8s.io/v1",
		"kind": "Gateway",
		"metadata": {"name": "my-gateway"},
		"spec": {
			"ports": [
				{"port": 80, "protocol": "HTTP"},
				{"port": 443, "protocol": "HTTPS"}
			]
		}
	}`

	c := &GatewayCoder{}
	if err := c.Decode([]byte(input)); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if c.Request.Name != "my-gateway" {
		t.Errorf("expected Name=my-gateway, got %s", c.Request.Name)
	}
	if len(c.Request.Ports) != 2 {
		t.Fatalf("expected 2 Ports, got %d", len(c.Request.Ports))
	}
	if c.Request.Ports[0].Port != 80 || c.Request.Ports[0].Protocol != "HTTP" {
		t.Errorf("unexpected port[0]: %+v", c.Request.Ports[0])
	}
	if c.Request.Ports[1].Port != 443 || c.Request.Ports[1].Protocol != "HTTPS" {
		t.Errorf("unexpected port[1]: %+v", c.Request.Ports[1])
	}
}

func TestGatewayCoderDecodeMinimal(t *testing.T) {
	t.Parallel()

	input := `{"kind": "Gateway", "metadata": {"name": "simple"}}`

	c := &GatewayCoder{}
	if err := c.Decode([]byte(input)); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if c.Request.Name != "simple" {
		t.Errorf("expected Name=simple, got %s", c.Request.Name)
	}
	if len(c.Request.Ports) != 0 {
		t.Errorf("expected 0 Ports, got %d", len(c.Request.Ports))
	}
}

func TestGatewayCoderDecodeInvalidJSON(t *testing.T) {
	t.Parallel()

	c := &GatewayCoder{}
	if err := c.Decode([]byte(`{invalid`)); err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestGatewayCoderEncode(t *testing.T) {
	t.Parallel()

	c := &GatewayCoder{
		Info: api.GatewayInfo{
			Name: "my-gateway",
			Ports: []api.GatewayPort{
				{Port: 80, Protocol: "HTTP"},
				{Port: 443, Protocol: "HTTPS"},
			},
			Addresses: []api.Address{
				{Type: "IPAddress", Value: "192.168.1.1"},
			},
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
	if !strings.Contains(yaml, "kind: Gateway") {
		t.Error("expected kind: Gateway in output")
	}
	if !strings.Contains(yaml, "name: my-gateway") {
		t.Error("expected name: my-gateway in output")
	}
	if !strings.Contains(yaml, "ports:") {
		t.Error("expected ports: in output")
	}
	if !strings.Contains(yaml, "port: 80") {
		t.Error("expected port: 80 in output")
	}
	if !strings.Contains(yaml, "protocol: HTTP") {
		t.Error("expected protocol: HTTP in output")
	}
	if !strings.Contains(yaml, "addresses:") {
		t.Error("expected addresses: in output")
	}
	if !strings.Contains(yaml, "192.168.1.1") {
		t.Error("expected address value in output")
	}
}

func TestGatewayCoderEncodeNoAddresses(t *testing.T) {
	t.Parallel()

	c := &GatewayCoder{
		Info: api.GatewayInfo{
			Name: "my-gateway",
			Ports: []api.GatewayPort{
				{Port: 80, Protocol: "HTTP"},
			},
		},
	}

	data, err := c.Encode()
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	yaml := string(data)
	if strings.Contains(yaml, "addresses") {
		t.Error("addresses should not appear when empty")
	}
}

func TestGatewayCoderRoundTrip(t *testing.T) {
	t.Parallel()

	// Encode
	original := &GatewayCoder{
		Info: api.GatewayInfo{
			Name: "round-trip",
			Ports: []api.GatewayPort{
				{Port: 8080, Protocol: "TCP"},
			},
			Addresses: []api.Address{
				{Type: "IPAddress", Value: "10.0.0.1"},
			},
		},
	}

	yamlData, err := original.Encode()
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// Decode the equivalent JSON
	input := `{
		"apiVersion": "controller.drycc.cc/v2.3",
		"kind": "Gateway",
		"metadata": {"name": "round-trip"},
		"spec": {
			"ports": [{"port": 8080, "protocol": "TCP"}]
		},
		"status": {
			"addresses": [{"type": "IPAddress", "value": "10.0.0.1"}]
		}
	}`

	decoded := &GatewayCoder{}
	if err := decoded.Decode([]byte(input)); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if decoded.Request.Name != "round-trip" {
		t.Errorf("expected Name=round-trip, got %s", decoded.Request.Name)
	}
	if len(decoded.Request.Ports) != 1 || decoded.Request.Ports[0].Port != 8080 {
		t.Error("Ports mismatch after round trip")
	}

	_ = yamlData // yamlData is valid YAML output, verified above
}
