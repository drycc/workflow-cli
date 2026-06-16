// Package coder provides K8s-style Manifest encoding/decoding for
// Gateway API resources (HTTPRoute, Gateway) used by the Drycc CLI.
//
// A Coder converts between the flat API types used by controller-sdk-go
// and the K8s Manifest format (kind/Metadata/spec/status) used in CLI YAML files.
package coder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	drycc "github.com/drycc/controller-sdk-go"
	"gopkg.in/yaml.v3"
)

// APIVersion is the Drycc API version used in Kubernetes-style manifests.
var APIVersion string = fmt.Sprintf("controller.drycc.cc/v%s", strings.Split(drycc.APIVersion, ".")[0])

// marshalYAML encodes a value to YAML with 2-space indentation and K8s-style
// array formatting (array items are not indented relative to their parent key).
// A JSON round-trip normalizes the entire structure to pure map[string]any /
// []any so that yaml.v3 renders all arrays with consistent compact style.
func marshalYAML(v interface{}) ([]byte, error) {
	// Normalize the entire structure to pure map[string]any / []any
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var normalized any
	if err := json.Unmarshal(data, &normalized); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(normalized); err != nil {
		return nil, err
	}
	if err := enc.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// normalize converts a typed struct or slice into plain map[string]any /
// []any via a JSON round-trip so that yaml.v3 renders arrays with consistent
// compact style (no extra indentation for array items).
func normalize(v any) any {
	data, err := json.Marshal(v)
	if err != nil {
		return v
	}
	var result any
	if err := json.Unmarshal(data, &result); err != nil {
		return v
	}
	return result
}

// Manifest is the K8s-style wrapper used for YAML input/output.
type Manifest struct {
	APIVersion string         `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind       string         `json:"kind" yaml:"kind"`
	Metadata   Metadata       `json:"metadata" yaml:"metadata"`
	Spec       map[string]any `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status     map[string]any `json:"status,omitempty" yaml:"status,omitempty"`
}

// Metadata represents the Kubernetes-style metadata for a manifest resource.
type Metadata struct {
	Name string `json:"name" yaml:"name"`
}
