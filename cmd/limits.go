package cmd

import (
	"fmt"
	"regexp"

	"github.com/drycc/pkg/prettyprint"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/config"
)

// LimitsList lists an app's limits.
func (d *DryccCmd) LimitsList(appID string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	config, err := config.List(s.Client, appID)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Printf("=== %s Limits\n\n", appID)

	d.Println("--- Memory")
	if len(config.Memory) == 0 {
		d.Println("Unlimited")
	} else {
		memoryMap := make(map[string]string)

		for key, value := range config.Memory {
			memoryMap[key] = fmt.Sprintf("%v", value)
		}

		d.Print(prettyprint.PrettyTabs(memoryMap, 5))
	}

	d.Println("\n--- CPU")
	if len(config.CPU) == 0 {
		d.Println("Unlimited")
	} else {
		cpuMap := make(map[string]string)

		for key, value := range config.CPU {
			cpuMap[key] = value.(string)
		}

		d.Print(prettyprint.PrettyTabs(cpuMap, 5))
	}

	return nil
}

// LimitsSet sets an app's limits.
func (d *DryccCmd) LimitsSet(appID string, limits []string, limitType string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	limitsMap, err := parseLimits(limits)
	if err != nil {
		return err
	}

	d.Print("Applying limits... ")

	quit := progress(d.WOut)
	configObj := api.Config{}

	if limitType == "cpu" {
		configObj.CPU = limitsMap
	} else {
		configObj.Memory = limitsMap
	}

	_, err = config.Set(s.Client, appID, configObj)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Print("done\n\n")

	return d.LimitsList(appID)
}

// LimitsUnset removes an app's limits.
func (d *DryccCmd) LimitsUnset(appID string, limits []string, limitType string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Print("Applying limits... ")

	quit := progress(d.WOut)

	configObj := api.Config{}

	valuesMap := make(map[string]interface{})

	for _, limit := range limits {
		valuesMap[limit] = nil
	}

	if limitType == "cpu" {
		configObj.CPU = valuesMap
	} else {
		configObj.Memory = valuesMap
	}

	_, err = config.Set(s.Client, appID, configObj)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Print("done\n\n")

	return d.LimitsList(appID)
}

func parseLimits(limits []string) (map[string]interface{}, error) {
	limitsMap := make(map[string]interface{})

	for _, limit := range limits {
		key, value, err := parseLimit(limit)

		if err != nil {
			return nil, err
		}

		limitsMap[key] = value
	}

	return limitsMap, nil
}

func parseLimit(limit string) (string, string, error) {
	regex := regexp.MustCompile("^([a-z0-9]+(?:-[a-z0-9]+)*)=(([0-9]+[bkmgBKMG]{1,2}|[0-9.]{1,5}|[0-9.]{1,5}m?))$")

	if !regex.MatchString(limit) {
		return "", "", fmt.Errorf(`%s doesn't fit format type=#unit or type=#
Examples: web=2G worker=500M db=1G`, limit)
	}

	capture := regex.FindStringSubmatch(limit)

	return capture[1], capture[2], nil
}
