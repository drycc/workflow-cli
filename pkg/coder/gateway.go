package coder

import (
	"encoding/json"

	"github.com/drycc/controller-sdk-go/api"
)

type GatewayCoder struct {
	Request api.GatewayUpdateRequest
	Info    api.GatewayInfo
}

func (c *GatewayCoder) Decode(data []byte) error {
	return json.Unmarshal(data, &c.Request)
}

func (c *GatewayCoder) Encode() ([]byte, error) {
	return marshalYAML(c.Info)
}
