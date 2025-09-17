package commands

import (
	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/appsettings"
	"github.com/drycc/workflow-cli/internal/loader"
)

// AutodeployInfo provides information about the status of app autodeploy.
func (d *DryccCmd) AutodeployInfo(appID string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	appSettings, err := appsettings.List(s.Client, appID)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	if appSettings.Autodeploy == nil || *appSettings.Autodeploy {
		d.Println("Autodeploy is enabled.")
	} else {
		d.Println("Autodeploy is disabled.")
	}
	return nil
}

// AutodeployEnable enables an app when deploy failed
func (d *DryccCmd) AutodeployEnable(appID string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	d.Printf("Enabling autodeploy for %s... ", appID)

	quit := progress(d.WOut)
	appSettings := api.AppSettings{Autodeploy: api.NewAutodeploy()}
	_, err = appsettings.Set(s.Client, appID, appSettings)
	quit <- true
	<-quit

	if err != nil {
		return err
	}

	d.Println("done")
	return nil
}

// AutodeployDisable disables an app when deploy failed
func (d *DryccCmd) AutodeployDisable(appID string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	d.Printf("Disabling autodeploy for %s... ", appID)

	quit := progress(d.WOut)
	appSettings := api.AppSettings{Autodeploy: api.NewAutodeploy()}
	*appSettings.Autodeploy = false
	_, err = appsettings.Set(s.Client, appID, appSettings)
	quit <- true
	<-quit

	if err != nil {
		return err
	}

	d.Println("done")
	return nil
}
