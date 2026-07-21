package coder

import (
	"encoding/json"

	"github.com/drycc/controller-sdk-go/api"
)

// AddonCoder implements Coder for addon instances.
//
// Decode field mapping (apply YAML → request):
//
//	metadata.name    → Name
//	kind             → Request.Kind
//	spec.plan        → Request.Plan
//	spec.parameters  → Request.Parameters
//
// Encode field mapping (info → YAML):
//
//	Info.Kind        → kind
//	Info.Name        → metadata.name
//	Info.Plan        → spec.plan
//	Info.Parameters  → spec.parameters
//	Info.Multiplier  → status.multiplier
type AddonCoder struct {
	Name    string
	Request api.AddonInstanceUpsertRequest
	Info    api.AddonInstance
}

// Decode unmarshals JSON manifest data into the Addon upsert request.
func (c *AddonCoder) Decode(data []byte) error {
	var env Manifest
	if err := json.Unmarshal(data, &env); err != nil {
		return err
	}

	c.Name = env.Metadata.Name
	c.Request = api.AddonInstanceUpsertRequest{
		Kind: env.Kind,
	}

	if plan, ok := env.Spec["plan"]; ok {
		if s, ok := plan.(string); ok {
			c.Request.Plan = s
		}
	}

	if params, ok := env.Spec["parameters"]; ok {
		if m, ok := params.(map[string]any); ok {
			c.Request.Parameters = m
		}
	}

	return nil
}

// Encode marshals the Addon info into a YAML manifest.
func (c *AddonCoder) Encode() ([]byte, error) {
	spec := make(map[string]any)
	spec["plan"] = c.Info.Plan
	if c.Info.Parameters != nil {
		spec["parameters"] = normalize(c.Info.Parameters)
	} else {
		spec["parameters"] = map[string]any{}
	}

	status := make(map[string]any)
	if c.Info.Multiplier != 0 {
		status["multiplier"] = c.Info.Multiplier
	}

	env := Manifest{
		APIVersion: APIVersion,
		Kind:       c.Info.Kind,
		Metadata:   Metadata{Name: c.Info.Name},
		Spec:       spec,
		Status:     status,
	}

	return marshalYAML(env)
}

// EncodeClass marshals an AddonClass catalog entry into a YAML manifest.
// This is used by the `addons classes info` command to display catalog
// details (description, storage model, plans) in a consistent format.
func EncodeClass(class api.AddonClass) ([]byte, error) {
	spec := map[string]any{
		"description":  class.Description,
		"storageModel": class.StorageModel,
		"plans":        normalize(class.Plans),
	}

	env := Manifest{
		APIVersion: APIVersion,
		Kind:       class.Kind,
		Metadata:   Metadata{Name: class.Name},
		Spec:       spec,
	}

	return marshalYAML(env)
}
