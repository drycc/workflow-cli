package coder

import (
	"strings"
	"testing"

	"github.com/drycc/controller-sdk-go/api"
)

func TestGatewayCoderDecode(t *testing.T) {
	t.Parallel()

	input := `
name = "my-gateway"

[[ports]]
port = 80
protocol = "HTTP"

[[ports]]
port = 443
protocol = "HTTPS"
`

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

	input := `name = "simple"`

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

func TestGatewayCoderDecodeInvalid(t *testing.T) {
	t.Parallel()

	c := &GatewayCoder{}
	if err := c.Decode([]byte(`{invalid`)); err == nil {
		t.Error("expected error for invalid TOML, got nil")
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

	toml := string(data)

	if !strings.Contains(toml, `name = "my-gateway"`) {
		t.Errorf("expected name in output, got:\n%s", toml)
	}
	if !strings.Contains(toml, "[[ports]]") {
		t.Errorf("expected [[ports]] in output, got:\n%s", toml)
	}
	if !strings.Contains(toml, "port = 80") {
		t.Errorf("expected port = 80 in output, got:\n%s", toml)
	}
	if strings.Contains(toml, "80.0") {
		t.Errorf("integer should not be rendered as float, got:\n%s", toml)
	}
	if !strings.Contains(toml, `protocol = "HTTP"`) {
		t.Errorf("expected protocol in output, got:\n%s", toml)
	}
	if !strings.Contains(toml, "[[addresses]]") {
		t.Errorf("expected [[addresses]] in output, got:\n%s", toml)
	}
	if !strings.Contains(toml, "192.168.1.1") {
		t.Errorf("expected address value in output, got:\n%s", toml)
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

	if strings.Contains(string(data), "addresses") {
		t.Errorf("addresses should not appear when empty, got:\n%s", string(data))
	}
}

func TestGatewayCoderRoundTrip(t *testing.T) {
	t.Parallel()

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

	encoded, err := original.Encode()
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// info output feeds back into apply directly
	decoded := &GatewayCoder{}
	if err := decoded.Decode(encoded); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if decoded.Request.Name != "round-trip" {
		t.Errorf("expected Name=round-trip, got %s", decoded.Request.Name)
	}
	if len(decoded.Request.Ports) != 1 || decoded.Request.Ports[0].Port != 8080 {
		t.Error("Ports mismatch after round trip")
	}
}
