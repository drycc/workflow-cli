package coder

import (
	"strings"
	"testing"

	"github.com/drycc/controller-sdk-go/api"
)

func TestServiceCoderDecode(t *testing.T) {
	t.Parallel()

	input := `
ptype = "worker"

[[ports]]
port = 5000
protocol = "TCP"
target_port = 5000

[[ports]]
port = 80
protocol = "UDP"
target_port = 8080
`

	c := &ServiceCoder{}
	if err := c.Decode([]byte(input)); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if c.Request.Ptype != "worker" {
		t.Errorf("expected Ptype=worker, got %s", c.Request.Ptype)
	}
	if len(c.Request.Ports) != 2 {
		t.Fatalf("expected 2 Ports, got %d", len(c.Request.Ports))
	}
	if c.Request.Ports[0].Port != 5000 || c.Request.Ports[0].TargetPort != 5000 {
		t.Errorf("unexpected port[0]: %+v", c.Request.Ports[0])
	}
	if c.Request.Ports[1].Protocol != "UDP" {
		t.Errorf("unexpected port[1]: %+v", c.Request.Ports[1])
	}
}

func TestServiceCoderDecodeInvalid(t *testing.T) {
	t.Parallel()

	c := &ServiceCoder{}
	if err := c.Decode([]byte(`{invalid`)); err == nil {
		t.Error("expected error for invalid TOML, got nil")
	}
}

func TestServiceCoderEncode(t *testing.T) {
	t.Parallel()

	c := &ServiceCoder{
		Info: api.ServiceInfo{
			Domain: "example-go-worker.example-go.svc",
			Ptype:  "worker",
			Ports: []api.ServicePort{
				{Port: 5000, Protocol: "TCP", TargetPort: 5000},
			},
		},
	}

	data, err := c.Encode()
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
	toml := string(data)
	if !strings.Contains(toml, `domain = "example-go-worker.example-go.svc"`) {
		t.Errorf("expected domain in output, got:\n%s", toml)
	}
	if !strings.Contains(toml, `ptype = "worker"`) {
		t.Errorf("expected ptype in output, got:\n%s", toml)
	}
	if !strings.Contains(toml, "[[ports]]") || !strings.Contains(toml, "target_port = 5000") {
		t.Errorf("expected ports in output, got:\n%s", toml)
	}
	if strings.Contains(toml, "5000.0") {
		t.Errorf("integer should not be rendered as float, got:\n%s", toml)
	}
	if strings.Contains(toml, "name =") || strings.Contains(toml, "tcp-5000") {
		t.Errorf("names should not be rendered, got:\n%s", toml)
	}
}

func TestServiceCoderRoundTrip(t *testing.T) {
	t.Parallel()

	original := &ServiceCoder{
		Info: api.ServiceInfo{
			Domain: "example-go-worker.example-go.svc",
			Ptype:  "worker",
			Ports: []api.ServicePort{
				{Port: 5000, Protocol: "TCP", TargetPort: 5000},
			},
		},
	}
	encoded, err := original.Encode()
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	decoded := &ServiceCoder{}
	if err := decoded.Decode(encoded); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if decoded.Request.Ptype != "worker" {
		t.Errorf("expected Ptype=worker, got %s", decoded.Request.Ptype)
	}
	if len(decoded.Request.Ports) != 1 || decoded.Request.Ports[0].TargetPort != 5000 {
		t.Errorf("Ports mismatch after round trip: %+v", decoded.Request.Ports)
	}
}
