package coder

import (
	"encoding/json"

	"github.com/drycc/controller-sdk-go/api"
)

// RouteCoder implements Coder for HTTPRoute resources.
//
// Decode field mapping:
//
//	Metadata.name         → Name
//	kind                  → Kind
//	spec.parentRefs       → Parents
//	spec.rules[].backends → Rules[].backendRefs
//
// Encode field mapping:
//
//	Rules[].backendRefs → spec.rules[].backends
//	Routable            → status.routable
type RouteCoder struct {
	Request api.RouteUpdateRequest
	Info    api.RouteInfo
}

// Decode unmarshals JSON data into the Route update request.
func (c *RouteCoder) Decode(data []byte) error {
	var env Manifest
	if err := json.Unmarshal(data, &env); err != nil {
		return err
	}

	c.Request = api.RouteUpdateRequest{
		Name: env.Metadata.Name,
		Kind: env.Kind,
	}

	if parents, ok := env.Spec["parents"]; ok {
		parentsJSON, err := json.Marshal(parents)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(parentsJSON, &c.Request.ParentRefs); err != nil {
			return err
		}
	}

	if rules, ok := env.Spec["rules"]; ok {
		if rulesList, ok := rules.([]any); ok {
			for _, r := range rulesList {
				if ruleMap, ok := r.(map[string]any); ok {
					routeRule := api.RouteRule{}
					for k, v := range ruleMap {
						if k == "backends" {
							routeRule["backendRefs"] = v
						} else {
							routeRule[k] = v
						}
					}
					c.Request.Rules = append(c.Request.Rules, routeRule)
				}
			}
		}
	}

	return nil
}

// Encode marshals the Route info into a YAML manifest.
func (c *RouteCoder) Encode() ([]byte, error) {
	var rules []map[string]any
	for _, r := range c.Info.Rules {
		rule := map[string]any{}
		if backends, ok := r["backendRefs"]; ok {
			rule["backends"] = backends
		}
		rules = append(rules, rule)
	}

	spec := make(map[string]any)
	spec["parents"] = normalize(c.Info.ParentRefs)
	spec["rules"] = rules

	status := make(map[string]any)
	if c.Info.Routable != nil {
		status["routable"] = c.Info.Routable
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
