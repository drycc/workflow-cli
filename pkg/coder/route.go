package coder

import (
	"github.com/drycc/controller-sdk-go/api"
)

type RouteCoder struct {
	Request api.RouteUpdateRequest
	Info    api.RouteInfo
}

func (c *RouteCoder) Decode(data []byte) error {
	return UnmarshalTOML(data, &c.Request)
}

func (c *RouteCoder) Encode() ([]byte, error) {
	return MarshalTOML(c.Info)
}
