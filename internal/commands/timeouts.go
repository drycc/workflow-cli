package commands

import (
	"fmt"
	"regexp"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/config"
	"github.com/drycc/workflow-cli/internal/loader"
)

// TimeoutsList lists an app's timeouts.
func (d *DryccCmd) TimeoutsList(appID string, version int) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	config, err := config.List(s.Client, appID, version)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	if len(config.Timeout) == 0 {
		d.Println("Default (30 sec) or controlled by drycc controller.")
	} else {
		table := d.getDefaultFormatTable([]string{"PTYPE", "TIMEOUT"})
		for _, key := range *sortKeys(config.Timeout) {
			table.Append([]string{
				key,
				fmt.Sprintf("%v", config.Timeout[key]),
			})
		}
		table.Render()
	}

	return nil
}

// TimeoutsSet sets an app's timeouts.
func (d *DryccCmd) TimeoutsSet(appID string, timeouts []string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	timeoutsMap, err := parseTimeouts(timeouts)
	if err != nil {
		return err
	}

	d.Print("Applying timeouts... ")

	quit := progress(d.WOut)
	configObj := api.Config{}

	configObj.Timeout = timeoutsMap

	_, err = config.Set(s.Client, appID, configObj)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Print("done\n\n")

	return d.TimeoutsList(appID, -1)
}

// TimeoutsUnset removes an app's timeouts.
func (d *DryccCmd) TimeoutsUnset(appID string, timeouts []string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	d.Print("Applying timeouts... ")

	quit := progress(d.WOut)

	configObj := api.Config{}

	valuesMap := make(map[string]any)

	for _, timeout := range timeouts {
		valuesMap[timeout] = nil
	}

	configObj.Timeout = valuesMap

	_, err = config.Set(s.Client, appID, configObj)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Print("done\n\n")

	return d.TimeoutsList(appID, -1)
}

func parseTimeouts(timeouts []string) (map[string]any, error) {
	timeoutsMap := make(map[string]any)

	for _, timeout := range timeouts {
		key, value, err := parseTimeout(timeout)
		if err != nil {
			return nil, err
		}

		timeoutsMap[key] = value
	}

	return timeoutsMap, nil
}

func parseTimeout(timeout string) (string, string, error) {
	regex := regexp.MustCompile("^([a-z0-9]+(?:-[a-z0-9]+)*)=([0-9]+)$")

	if !regex.MatchString(timeout) {
		return "", "", fmt.Errorf(`%s doesn't fit format type=#
Examples: web=30 worker=300`, timeout)
	}

	capture := regex.FindStringSubmatch(timeout)

	return capture[1], capture[2], nil
}
