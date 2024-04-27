package cmd

import (
	"fmt"
	"strings"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/config"
)

// RegistryList lists an app's registry information.
func (d *DryccCmd) RegistryList(appID string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	config, err := config.List(s.Client, appID)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	if len(config.Registry) == 0 {
		d.Println(fmt.Sprintf("No registrys found in %s app.", appID))
	} else {
		table := d.getDefaultFormatTable([]string{"KEY", "VALUE"})
		for _, key := range *sortKeys(config.Registry) {
			table.Append([]string{
				key,
				fmt.Sprintf("%v", config.Registry[key]),
			})
		}
		table.Render()
	}
	return nil
}

// RegistrySet sets an app's registry information.
func (d *DryccCmd) RegistrySet(appID string, item []string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	registryMap, err := parseInfos(item)
	if err != nil {
		return err
	}

	d.Print("Applying registry information... ")

	quit := progress(d.WOut)
	configObj := api.Config{}
	configObj.Registry = registryMap

	_, err = config.Set(s.Client, appID, configObj)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Print("done\n\n")

	return d.RegistryList(appID)
}

// RegistryUnset removes an app's registry information.
func (d *DryccCmd) RegistryUnset(appID string, items []string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Print("Applying registry information... ")

	quit := progress(d.WOut)

	configObj := api.Config{}

	registryMap := make(map[string]interface{})

	for _, key := range items {
		registryMap[key] = nil
	}

	configObj.Registry = registryMap

	_, err = config.Set(s.Client, appID, configObj)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Print("done\n\n")

	return d.RegistryList(appID)
}

func parseInfos(items []string) (map[string]interface{}, error) {
	registryMap := make(map[string]interface{})

	for _, item := range items {
		key, value, err := parseInfo(item)

		if err != nil {
			return nil, err
		}

		registryMap[key] = value
	}

	return registryMap, nil
}

func parseInfo(item string) (string, string, error) {
	parts := strings.SplitN(item, "=", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf(`%s is invalid. Must be in format key=value
Examples: username=bob password=s3cur3pw1`, item)
	}

	if parts[0] != "username" && parts[0] != "password" {
		return "", "", fmt.Errorf(`%s is invalid. Valid keys are "username" or "password"`, parts[0])
	}

	return parts[0], parts[1], nil
}
