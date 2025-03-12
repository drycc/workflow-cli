package commands

import (
	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/appsettings"
	"github.com/drycc/workflow-cli/internal/utils"
)

// AutorollbackInfo provides information about the status of app autorollback.
func (d *DryccCmd) AutorollbackInfo(appID string) error {
	appID, s, err := utils.LoadAppSettings(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	appSettings, err := appsettings.List(s.Client, appID)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	if appSettings.Autorollback == nil || *appSettings.Autorollback {
		d.Println("Autorollback is enabled.")
	} else {
		d.Println("Autorollback is disabled.")
	}
	return nil
}

// AutorollbackEnable enables an app when deploy failed
func (d *DryccCmd) AutorollbackEnable(appID string) error {
	appID, s, err := utils.LoadAppSettings(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Enabling autorollback for %s... ", appID)

	quit := progress(d.WOut)
	appSettings := api.AppSettings{Autorollback: api.NewAutorollback()}
	_, err = appsettings.Set(s.Client, appID, appSettings)
	quit <- true
	<-quit

	if err != nil {
		return err
	}

	d.Println("done")
	return nil
}

// AutorollbackDisable disables an app when deploy failed
func (d *DryccCmd) AutorollbackDisable(appID string) error {
	appID, s, err := utils.LoadAppSettings(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Disabling autorollback for %s... ", appID)

	quit := progress(d.WOut)
	appSettings := api.AppSettings{Autorollback: api.NewAutorollback()}
	*appSettings.Autorollback = false
	_, err = appsettings.Set(s.Client, appID, appSettings)
	quit <- true
	<-quit

	if err != nil {
		return err
	}

	d.Println("done")
	return nil
}
