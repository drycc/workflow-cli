package coder

import (
	"encoding/json"

	"github.com/drycc/controller-sdk-go/api"
)

// GatewayCoder implements Coder for Gateway resources.
//
// Decode field mapping:
//
//	Metadata.name → Name
//	spec.ports    → Ports
//
// Encode field mapping:
//
//	Ports     → spec.ports
//	Addresses → status.addresses
type GatewayCoder struct {
	Request api.GatewayUpdateRequest
	Info    api.GatewayInfo
}

func (c *GatewayCoder) Decode(data []byte) error {
	var env Manifest
	if err := json.Unmarshal(data, &env); err != nil {
		return err
	}

	c.Request = api.GatewayUpdateRequest{
		Name: env.Metadata.Name,
	}

	if ports, ok := env.Spec["ports"]; ok {
		portsJSON, err := json.Marshal(ports)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(portsJSON, &c.Request.Ports); err != nil {
			return err
		}
	}

	return nil
}

func (c *GatewayCoder) Encode() ([]byte, error) {
	spec := map[string]any{"ports": normalize(c.Info.Ports)}

	status := make(map[string]any)
	if len(c.Info.Addresses) > 0 {
		status["addresses"] = normalize(c.Info.Addresses)
	}

	env := Manifest{
		APIVersion: APIVersion,
		Kind:       "Gateway",
		Metadata:   Metadata{Name: c.Info.Name},
		Spec:       spec,
		Status:     status,
	}

	return marshalYAML(env)
}
