package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/gateways"
	"github.com/drycc/workflow-cli/internal/loader"
	"github.com/drycc/workflow-cli/pkg/coder"
	"sigs.k8s.io/yaml"
)

// GatewaysList lists gateways for the app
func (d *DryccCmd) GatewaysList(appID string, results int) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}
	if results == defaultLimit {
		results = s.Limit
	}

	gateways, count, err := gateways.List(s.Client, appID, results)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	if count == 0 {
		d.Println(fmt.Sprintf("No gateways found in %s app.", appID))
	} else {
		table := d.getDefaultFormatTable([]string{"NAME", "PORT", "PROTOCOL", "ADDRESSES"})
		for _, gateway := range gateways {
			addresesStr := parseAddress(gateway.Addresses)
			for _, port := range gateway.Ports {
				table.Append([]string{gateway.Name, fmt.Sprint(port.Port), port.Protocol, addresesStr})
			}
		}
		table.Render()
	}
	return nil
}

// GatewaysInfo shows detailed information about a gateway.
func (d *DryccCmd) GatewaysInfo(appID, name string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	info, err := gateways.Info(s.Client, appID, name)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	c := &coder.GatewayCoder{Info: info}
	yamlBytes, err := c.Encode()
	if err != nil {
		return err
	}
	d.Println(string(yamlBytes))
	return nil
}

// GatewaysApply applies gateway configuration from a YAML file.
func (d *DryccCmd) GatewaysApply(appID, filePath string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	yamlData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	jsonData, err := yaml.YAMLToJSON(yamlData)
	if err != nil {
		return err
	}

	c := &coder.GatewayCoder{}
	if err := c.Decode(jsonData); err != nil {
		return fmt.Errorf("invalid gateway configuration: %w", err)
	}
	if c.Request.Name == "" {
		return fmt.Errorf("invalid gateway configuration: missing metadata.name")
	}
	req := c.Request

	d.Printf("Applying gateway %s to %s... ", req.Name, appID)

	quit := progress(d.WOut)
	_, err = gateways.Apply(s.Client, appID, req)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// GatewaysRemove removes a gateway from an app.
func (d *DryccCmd) GatewaysRemove(appID, name string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}
	d.Printf("Removing gateway %s from %s... ", name, appID)

	quit := progress(d.WOut)
	err = gateways.Delete(s.Client, appID, name)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

func parseAddress(addresses []api.Address) string {
	var addresList []string
	for _, address := range addresses {
		addresList = append(addresList, address.Value)
	}
	addresStr := strings.Join(addresList, ",")
	return addresStr
}
