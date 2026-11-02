package coder

import (
	"github.com/drycc/controller-sdk-go/api"
)

type GatewayCoder struct {
	Request api.GatewayUpdateRequest
	Info    api.GatewayInfo
}

func (c *GatewayCoder) Decode(data []byte) error {
	return UnmarshalTOML(data, &c.Request)
}

func (c *GatewayCoder) Encode() ([]byte, error) {
	return MarshalTOML(c.Info)
}
