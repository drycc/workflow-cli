package commands

import (
	"fmt"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/config"
	"github.com/drycc/workflow-cli/internal/loader"
)

// RegistryList lists an app's registry information.
func (d *DryccCmd) RegistryList(appID, ptype string, version int) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	config, err := config.List(s.Client, appID, version)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	if len(config.Registry) == 0 {
		d.Println(fmt.Sprintf("No registrys found in %s app.", appID))
		return nil
	}
	ptypes := []string{}
	if ptype != "" {
		ptypes = append(ptypes, ptype)
	} else {
		for ptype := range config.Registry {
			ptypes = append(ptypes, ptype)
		}
	}

	table := d.getDefaultFormatTable([]string{"PTYPE", "USERNAME", "PASSWORD"})
	for _, ptype := range sortPtypes(ptypes) {
		if config.Registry[ptype]["username"] != nil {
			table.Append([]string{
				ptype,
				fmt.Sprintf("%v", config.Registry[ptype]["username"]),
				fmt.Sprintf("%v", config.Registry[ptype]["password"]),
			})
		}
	}

	table.Render()

	return nil
}

// RegistrySet sets an app's registry information.
func (d *DryccCmd) RegistrySet(appID, ptype, username, password string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	d.Print("Applying registry information... ")

	quit := progress(d.WOut)

	configObj := api.Config{}
	registry := make(map[string]map[string]any)
	registry[ptype] = map[string]any{
		"username": username,
		"password": password,
	}
	configObj.Registry = registry

	_, err = config.Set(s.Client, appID, configObj, true)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Print("done\n\n")

	return d.RegistryList(appID, ptype, -1)
}

// RegistryUnset removes an app's registry information.
func (d *DryccCmd) RegistryUnset(appID, ptype string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	d.Print("Applying registry information... ")

	quit := progress(d.WOut)

	configObj := api.Config{}
	registry := make(map[string]map[string]any)
	registry[ptype] = map[string]any{
		"username": nil,
		"password": nil,
	}
	configObj.Registry = registry
	_, err = config.Set(s.Client, appID, configObj, true)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Print("done\n\n")

	return nil
}
