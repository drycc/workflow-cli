package coder

import (
	"github.com/drycc/controller-sdk-go/api"
)

type AddonCoder struct {
	Name    string
	Request api.AddonInstanceUpsertRequest
	Info    api.AddonInstance
}

func (c *AddonCoder) Decode(data []byte) error {
	var combined struct {
		Name string `json:"name"`
		api.AddonInstanceUpsertRequest
	}
	if err := UnmarshalTOML(data, &combined); err != nil {
		return err
	}
	c.Name = combined.Name
	c.Request = combined.AddonInstanceUpsertRequest
	return nil
}

func (c *AddonCoder) Encode() ([]byte, error) {
	return MarshalTOML(c.Info)
}

func EncodeClass(class api.AddonClass) ([]byte, error) {
	return MarshalTOML(class)
}
