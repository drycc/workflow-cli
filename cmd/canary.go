package cmd

import (
	"fmt"
	"strings"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/appsettings"
)

// CanaryList tells the informations about app's autoscale status
func (d *DryccCmd) CanaryInfo(appID string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	appSettings, err := appsettings.List(s.Client, appID)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	if len(appSettings.Canaries) > 0 {
		table := d.getDefaultFormatTable([]string{"OWNER", "PTYPE", "CREATED", "UPDATED"})
		for _, procType := range appSettings.Canaries {
			table.Append([]string{
				appSettings.Owner,
				procType,
				d.formatTime(appSettings.Created),
				d.formatTime(appSettings.Updated),
			})
		}
		table.Render()
	} else {
		d.Println(fmt.Sprintf("No canaries found in %s app.", appID))
	}
	return nil
}

// CanaryCreate sets canary options for the app proc type.
func (d *DryccCmd) CanaryCreate(appID string, processType []string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Applying canary settings for process type %s on %s... ", strings.Join(processType, ","), appID)

	quit := progress(d.WOut)
	_, err = appsettings.Set(s.Client, appID, api.AppSettings{Canaries: processType})

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	d.Println("done")
	return nil
}

// CanaryRemove remove canary for the app proc type.
func (d *DryccCmd) CanaryRemove(appID string, processType []string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Removing canary for process type %s on %s... ", strings.Join(processType, ","), appID)

	quit := progress(d.WOut)

	err = appsettings.CanaryRemove(s.Client, appID, api.AppSettings{Canaries: processType})

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	d.Println("done")
	return nil
}

// CanaryRelease release canary for the app.
func (d *DryccCmd) CanaryRelease(appID string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Release canary for %s... ", appID)

	quit := progress(d.WOut)

	err = appsettings.CanaryRelease(s.Client, appID)

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	d.Println("done")
	return nil
}

// CanaryRollback rollback canary for the app.
func (d *DryccCmd) CanaryRollback(appID string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Rollback canary for %s... ", appID)

	quit := progress(d.WOut)

	err = appsettings.CanaryRollback(s.Client, appID)

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	d.Println("done")
	return nil
}
