package cmd

import (
	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/appsettings"
)

// MaintenanceInfo tells the informations about app's maintenance status
func (d *DryccCmd) MaintenanceInfo(appID string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	appSettings, err := appsettings.List(s.Client, appID)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	if appSettings.Maintenance == nil || !*appSettings.Maintenance {
		d.Println("Maintenance mode is off.")
	} else {
		d.Println("Maintenance mode is on.")
	}
	return nil
}

// MaintenanceEnable turns on the maintenance for the app.
func (d *DryccCmd) MaintenanceEnable(appID string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Enabling maintenance mode for %s... ", appID)

	quit := progress(d.WOut)
	b := true
	_, err = appsettings.Set(s.Client, appID, api.AppSettings{Maintenance: &b})

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	d.Println("done")
	return nil
}

// MaintenanceDisable turns off the maintenance for the app.
func (d *DryccCmd) MaintenanceDisable(appID string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Disabling maintenance mode for %s... ", appID)

	quit := progress(d.WOut)
	b := false
	_, err = appsettings.Set(s.Client, appID, api.AppSettings{Maintenance: &b})

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	d.Println("done")
	return nil
}
