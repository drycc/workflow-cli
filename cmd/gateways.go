package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/drycc/controller-sdk-go/gateways"
	"github.com/olekukonko/tablewriter"
)

// GatewaysList lists gateways for the app
func (d *DryccCmd) GatewaysList(appID string, results int) error {
	s, appID, err := load(d.ConfigFile, appID)

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
	d.Printf("=== %s Gateways\n", appID)
	if count == 0 {
		d.Println("Could not find any gateway")
	} else {
		table := tablewriter.NewWriter(d.WOut)
		table.SetHeader([]string{"Name", "Lisenter", "Port", "Protocol", "allowedRoutes"})
		for _, gateway := range gateways {
			for _, listener := range gateway.Listeners {
				allowedRoutes, _ := json.Marshal(listener.AllowedRoutes)
				table.Append([]string{gateway.Name, listener.Name, fmt.Sprint(listener.Port), listener.Protocol, string(allowedRoutes)})
			}
		}
		table.SetAutoMergeCellsByColumnIndex([]int{0})
		table.SetRowLine(true)
		table.Render()
	}
	return nil
}

// GatewaysAdd adds a gateway to an app.
func (d *DryccCmd) GatewaysAdd(appID, name string, port int, protocol string) error {
	s, appID, err := load(d.ConfigFile, appID)

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
	s, appID, err := load(d.ConfigFile, appID)

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
