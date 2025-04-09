package commands

import (
	"fmt"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/appsettings"
	"github.com/drycc/workflow-cli/internal/utils"
)

// AutoscaleList tells the informations about app's autoscale status
func (d *DryccCmd) AutoscaleList(appID string) error {
	appID, s, err := utils.LoadAppSettings(d.ConfigFile, appID)

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
		table := d.getDefaultFormatTable([]string{"PTYPE", "PERCENT", "MIN", "MAX"})
		for process, kv := range appSettings.Autoscale {
			table.Append([]string{
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
func (d *DryccCmd) AutoscaleSet(appID string, ptype string, minCPU int, maxCPU int, CPUPercent int) error {
	appID, s, err := utils.LoadAppSettings(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Applying autoscale settings for process type %s on %s... ", ptype, appID)

	quit := progress(d.WOut)
	data := map[string]*api.Autoscale{
		ptype: {
			Min:        minCPU,
			Max:        maxCPU,
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
func (d *DryccCmd) AutoscaleUnset(appID string, ptype string) error {
	appID, s, err := utils.LoadAppSettings(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Removing autoscale for process type %s on %s... ", ptype, appID)

	quit := progress(d.WOut)
	data := map[string]*api.Autoscale{
		ptype: nil,
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
