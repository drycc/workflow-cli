package cmd

import (
	"fmt"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/appsettings"
)

// AutoscaleList tells the informations about app's autoscale status
func (d *DryccCmd) AutoscaleList(appID string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	appSettings, err := appsettings.List(s.Client, appID)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	if appSettings.Autoscale == nil {
		d.Println("No autoscale rules found.")
	} else {
		table := d.getDefaultFormatTable([]string{"UUID", "PTYPE", "PERCENT", "MIN", "MAX"})
		for process, kv := range appSettings.Autoscale {
			table.Append([]string{
				appSettings.UUID,
				process,
				fmt.Sprintf("%d", (*kv).CPUPercent),
				fmt.Sprintf("%d", (*kv).Min),
				fmt.Sprintf("%d", (*kv).Max),
			})
		}
		table.Render()
	}

	return nil
}

// AutoscaleSet sets autoscale options for the app.
func (d *DryccCmd) AutoscaleSet(appID string, processType string, min int, max int, CPUPercent int) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Applying autoscale settings for process type %s on %s... ", processType, appID)

	quit := progress(d.WOut)
	data := map[string]*api.Autoscale{
		processType: {
			Min:        min,
			Max:        max,
			CPUPercent: CPUPercent,
		},
	}
	_, err = appsettings.Set(s.Client, appID, api.AppSettings{Autoscale: data})

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	d.Println("done")
	return nil
}

// AutoscaleUnset removes autoscale for the app.
func (d *DryccCmd) AutoscaleUnset(appID string, processType string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Removing autoscale for process type %s on %s... ", processType, appID)

	quit := progress(d.WOut)
	data := map[string]*api.Autoscale{
		processType: nil,
	}
	_, err = appsettings.Set(s.Client, appID, api.AppSettings{Autoscale: data})

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	d.Println("done")
	return nil
}
