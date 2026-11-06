package commands

import (
	"fmt"
	"os"

	"github.com/drycc/controller-sdk-go/services"
	"github.com/drycc/workflow-cli/internal/loader"
	"github.com/drycc/workflow-cli/pkg/coder"
)

// ServicesList lists services for the app
func (d *DryccCmd) ServicesList(appID string, results int) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}
	if results == defaultLimit {
		results = s.Limit
	}

	services, count, err := services.List(s.Client, appID, results)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	if count == 0 {
		d.Println(fmt.Sprintf("No services found in %s app.", appID))
	} else {
		table := d.getDefaultFormatTable([]string{"PTYPE", "PORT", "PROTOCOL", "TARGET_PORT", "DOMAIN"})
		for _, service := range services {
			for _, port := range service.Ports {
				table.Append([]string{
					service.Ptype,
					fmt.Sprint(port.Port),
					port.Protocol,
					fmt.Sprint(port.TargetPort),
					service.Domain,
				})
			}
		}
		table.Render()
	}
	return nil
}

// ServicesInfo shows detailed information about a service.
func (d *DryccCmd) ServicesInfo(appID, ptype string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	info, err := services.Info(s.Client, appID, ptype)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	// app is implied by the command context and apply overrides it, so omit it
	info.App = ""

	c := &coder.ServiceCoder{Info: info}
	tomlBytes, err := c.Encode()
	if err != nil {
		return err
	}
	d.Println(string(tomlBytes))
	return nil
}

// ServicesApply applies service configuration from a TOML file.
func (d *DryccCmd) ServicesApply(appID, filePath string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	tomlData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	c := &coder.ServiceCoder{}
	if err := c.Decode(tomlData); err != nil {
		return fmt.Errorf("invalid service configuration: %w", err)
	}
	if c.Request.Ptype == "" {
		return fmt.Errorf("invalid service configuration: missing ptype")
	}
	req := c.Request

	d.Printf("Applying service %s to %s... ", req.Ptype, appID)

	quit := progress(d.WOut)
	_, err = services.Apply(s.Client, appID, req)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// ServicesRemove removes a service from an app.
func (d *DryccCmd) ServicesRemove(appID, ptype string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}
	d.Printf("Removing service %s from %s... ", ptype, appID)

	quit := progress(d.WOut)
	err = services.Delete(s.Client, appID, ptype)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}
