package coder

import (
	"github.com/drycc/controller-sdk-go/api"
)

type ServiceCoder struct {
	Request api.ServiceUpdateRequest
	Info    api.ServiceInfo
}

func (c *ServiceCoder) Decode(data []byte) error {
	return UnmarshalTOML(data, &c.Request)
}

func (c *ServiceCoder) Encode() ([]byte, error) {
	return MarshalTOML(c.Info)
}
