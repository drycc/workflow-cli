package coder

import (
	"encoding/json"

	"github.com/drycc/controller-sdk-go/api"
)

type AddonCoder struct {
	Request api.AddonInstanceUpsertRequest
	Info    api.AddonInstance
}

func (c *AddonCoder) Decode(data []byte) error {
	return json.Unmarshal(data, &c.Request)
}

func (c *AddonCoder) Encode() ([]byte, error) {
	return marshalYAML(c.Info)
}

func EncodeClass(class api.AddonClass) ([]byte, error) {
	return marshalYAML(class)
}
