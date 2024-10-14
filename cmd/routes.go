package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/routes"
	"sigs.k8s.io/yaml"
)

// RoutesCreate create a route to an app.
func (d *DryccCmd) RoutesCreate(appID, name string, kind string, backendRefs ...api.BackendRefRequest) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}
	d.Printf("Adding route %s to %s... ", name, appID)

	quit := progress(d.WOut)
	err = routes.New(s.Client, appID, name, kind, backendRefs...)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// RoutesList lists routes for the app
func (d *DryccCmd) RoutesList(appID string, results int) error {
	s, appID, err := load(d.ConfigFile, appID)

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
		table := d.getDefaultFormatTable([]string{"NAME", "OWNER", "KIND", "SERVICE", "GATEWAY"})
		for _, route := range routes {
			var services []string
			for _, rule := range route.Rules {
				if backends, ok := rule["backendRefs"].([]interface{}); ok {
					for _, backend := range backends {
						if service, ok := backend.(map[string]interface{}); ok {
							services = append(services, fmt.Sprintf("%v:%v", service["name"], service["port"]))
						}
					}
				}
			}
			var gateways []string
			for _, gateway := range route.ParentRefs {
				gateways = append(gateways, fmt.Sprintf("%s:%d", gateway.Name, gateway.Port))
			}
			table.Append([]string{
				route.Name,
				route.Owner,
				route.Kind,
				strings.Join(services, "\n"),
				strings.Join(gateways, "\n"),
			})
			fmt.Println("==============================", services, services)
		}
		table.Render()
	}
	return nil
}

// RoutesAttach bind a route to gateway.
func (d *DryccCmd) RoutesAttach(appID, name string, port int, gateway string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}
	d.Printf("Attaching route %s to gateway %s... ", name, gateway)

	quit := progress(d.WOut)
	err = routes.AttachGateway(s.Client, appID, name, port, gateway)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// RoutesDetach bind a route to gateway.
func (d *DryccCmd) RoutesDetach(appID, name string, port int, gateway string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}
	d.Printf("Detaching route %s to gateway %s... ", name, gateway)

	quit := progress(d.WOut)
	err = routes.DetachGateway(s.Client, appID, name, port, gateway)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// RoutesGet get rule of route for the app
func (d *DryccCmd) RoutesGet(appID string, name string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	route, err := routes.GetRule(s.Client, appID, name)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	var rules []byte
	rules, err = yaml.JSONToYAML([]byte(route))
	if err != nil {
		return err
	}
	d.Println(string(rules))
	return nil
}

// RoutesSet set rule of route for the app
func (d *DryccCmd) RoutesSet(appID string, name string, ruleFile string) error {
	s, appID, err := load(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	var contents []byte
	if _, err := os.Stat(ruleFile); err != nil {
		return err
	}
	contents, err = os.ReadFile(ruleFile)
	if err != nil {
		return err
	}
	rules, err := yaml.YAMLToJSON(contents)
	if err != nil {
		return err
	}
	d.Print("Applying rules... ")
	quit := progress(d.WOut)
	err = routes.SetRule(s.Client, appID, name, string(rules))
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
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}
	d.Printf("Removing route %s to %s... ", name, appID)

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
