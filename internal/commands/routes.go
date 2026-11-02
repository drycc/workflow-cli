package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/drycc/controller-sdk-go/routes"
	"github.com/drycc/workflow-cli/internal/loader"
	"github.com/drycc/workflow-cli/pkg/coder"
)

// RoutesList lists routes for the app
func (d *DryccCmd) RoutesList(appID string, results int) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}
	if results == defaultLimit {
		results = s.Limit
	}

	routes, count, err := routes.List(s.Client, appID, results)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	if count == 0 {
		d.Println(fmt.Sprintf("No routes found in %s app.", appID))
	} else {
		table := d.getDefaultFormatTable([]string{"NAME", "KIND", "GATEWAYS", "SERVICES"})
		for _, route := range routes {
			var services []string
			for _, rule := range route.Rules {
				if backends, ok := rule["backend_refs"].([]any); ok {
					for _, backend := range backends {
						if service, ok := backend.(map[string]any); ok {
							services = append(services, fmt.Sprintf("%v:%v", service["name"], service["port"]))
						}
					}
				}
			}
			var gateways []string
			for _, gateway := range route.ParentRefs {
				gateways = append(gateways, fmt.Sprintf("%s:%d", gateway.Name, gateway.Port))
			}
			gatewaysBytes, err := json.Marshal(gateways)
			if err != nil {
				return err
			}
			servicesBytes, err := json.Marshal(services)
			if err != nil {
				return err
			}
			table.Append([]string{route.Name, route.Kind, string(gatewaysBytes), string(servicesBytes)})
		}
		table.Render()
	}
	return nil
}

// RoutesInfo shows detailed information about a route.
func (d *DryccCmd) RoutesInfo(appID, name string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	info, err := routes.Info(s.Client, appID, name)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	// app is implied by the command context and apply overrides it, so omit it
	info.App = ""

	c := &coder.RouteCoder{Info: info}
	tomlBytes, err := c.Encode()
	if err != nil {
		return err
	}
	d.Println(string(tomlBytes))
	return nil
}

// RoutesApply applies route configuration from a TOML file.
func (d *DryccCmd) RoutesApply(appID, filePath string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	tomlData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	c := &coder.RouteCoder{}
	if err := c.Decode(tomlData); err != nil {
		if strings.Contains(err.Error(), "rules") {
			return fmt.Errorf("invalid route configuration: %w\n(hint: rules must be a list, did you forget the [[rules]] header?)", err)
		}
		return fmt.Errorf("invalid route configuration: %w", err)
	}
	if c.Request.Name == "" {
		return fmt.Errorf("invalid route configuration: missing name")
	}
	req := c.Request

	d.Printf("Applying route %s to %s... ", req.Name, appID)

	quit := progress(d.WOut)
	_, err = routes.Apply(s.Client, appID, req)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// RoutesRemove removes a route registered with an app.
func (d *DryccCmd) RoutesRemove(appID, name string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}
	d.Printf("Removing route %s from %s... ", name, appID)

	quit := progress(d.WOut)
	err = routes.Delete(s.Client, appID, name)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}
