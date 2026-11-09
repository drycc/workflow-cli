package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/drycc/controller-sdk-go/services"
	"github.com/drycc/workflow-cli/internal/loader"
	"github.com/drycc/workflow-cli/pkg/coder"
)

var errNameAddressedServicesUnsupported = errors.New(
	"controller does not support name-addressed services; upgrade controller")

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
		for _, service := range services {
			if service.Name == "" {
				return errNameAddressedServicesUnsupported
			}
		}
		table := d.getDefaultFormatTable([]string{"NAME", "PTYPE", "PORT", "PROTOCOL", "TARGET_PORT", "DOMAIN"})
		for _, service := range services {
			for _, port := range service.Ports {
				table.Append([]string{
					service.Name,
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
func (d *DryccCmd) ServicesInfo(appID, name string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	info, err := services.Info(s.Client, appID, name)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	if info.Name != name {
		return errNameAddressedServicesUnsupported
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
	if c.Request.Name == "" {
		return fmt.Errorf("invalid service configuration: missing name")
	}
	if c.Request.Ptype == "" {
		return fmt.Errorf("invalid service configuration: missing ptype")
	}
	req := c.Request

	d.Printf("Applying service %s to %s... ", req.Name, appID)

	quit := progress(d.WOut)
	info, err := services.Apply(s.Client, appID, req)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	if info.Name != req.Name {
		return errNameAddressedServicesUnsupported
	}

	d.Println("done")
	return nil
}

// ServicesRemove removes a service from an app.
func (d *DryccCmd) ServicesRemove(appID, name string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}
	info, err := services.Info(s.Client, appID, name)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	if info.Name != name {
		return errNameAddressedServicesUnsupported
	}
	d.Printf("Removing service %s from %s... ", name, appID)

	quit := progress(d.WOut)
	err = services.Delete(s.Client, appID, name)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}
