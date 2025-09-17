package commands

import (
	"fmt"
	"strings"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/gateways"
	"github.com/drycc/workflow-cli/internal/loader"
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
		table := d.getDefaultFormatTable([]string{"NAME", "LISENTER", "PORT", "PROTOCOL", "ADDRESSES"})
		for _, gateway := range gateways {
			addresesStr := parseAddress(gateway.Addresses)
			for _, listener := range gateway.Listeners {
				table.Append([]string{gateway.Name, listener.Name, fmt.Sprint(listener.Port), listener.Protocol, addresesStr})
			}
		}
		table.Render()
	}
	return nil
}

// GatewaysAdd adds a gateway to an app.
func (d *DryccCmd) GatewaysAdd(appID, name string, port int, protocol string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}
	d.Printf("Adding gateway %s to %s... ", name, appID)

	quit := progress(d.WOut)
	err = gateways.New(s.Client, appID, name, port, protocol)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// GatewaysRemove removes a gateway registered with an app.
func (d *DryccCmd) GatewaysRemove(appID, name string, port int, protocol string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}
	d.Printf("Removing gateway %s to %s... ", name, appID)

	quit := progress(d.WOut)
	err = gateways.Delete(s.Client, appID, name, port, protocol)
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
