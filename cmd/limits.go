package cmd

import (
	"fmt"
	"regexp"

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
	if len(config.CPU) > 0 || len(config.Memory) > 0 {
		table := d.getDefaultFormatTable([]string{"UUID", "OWNER", "TYPE", "DEVICE", "QUOTA"})
		for _, key := range *sortKeys(config.Memory) {
			table.Append([]string{
				config.UUID,
				config.Owner,
				key,
				"MEM",
				fmt.Sprintf("%v", config.Memory[key]),
			})
		}
		for _, key := range *sortKeys(config.CPU) {
			table.Append([]string{
				config.UUID,
				config.Owner,
				key,
				"CPU",
				fmt.Sprintf("%v", config.CPU[key]),
			})
		}
		table.Render()
	} else {
		d.Println(fmt.Sprintf("No limits found in %s app.", appID))
	}
	return nil
}

// LimitsSet sets an app's limits.
func (d *DryccCmd) LimitsSet(appID string, cpuLimits []string, memoryLimits []string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	configObj := api.Config{}
	if len(cpuLimits) > 0 {
		cpuLimitsMap, err := parseLimits(cpuLimits)
		if err != nil {
			return err
		}
		configObj.CPU = cpuLimitsMap
	}
	if len(memoryLimits) > 0 {
		memoryLimitsMap, err := parseLimits(memoryLimits)
		if err != nil {
			return err
		}
		configObj.Memory = memoryLimitsMap
	}

	d.Print("Applying limits... ")

	quit := progress(d.WOut)

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
func (d *DryccCmd) LimitsUnset(appID string, cpuLimits []string, memoryLimits []string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Print("Applying limits... ")

	quit := progress(d.WOut)

	configObj := api.Config{}
	if len(cpuLimits) > 0 {
		cpuMap := make(map[string]interface{})
		for _, limit := range cpuLimits {
			cpuMap[limit] = nil
		}
		configObj.CPU = cpuMap
	}
	if len(memoryLimits) > 0 {
		memoryMap := make(map[string]interface{})
		for _, limit := range memoryLimits {
			memoryMap[limit] = nil
		}
		configObj.Memory = memoryMap
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
	regex := regexp.MustCompile("^([a-z0-9]+(?:-[a-z0-9]+)*)=(([1-9][0-9]*[mgMG]|[1-9][0-9]*m?))$")

	if !regex.MatchString(limit) {
		return "", "", fmt.Errorf(`%s doesn't fit format type=#unit or type=#
Examples: web=2G worker=500M db=1G`, limit)
	}

	capture := regex.FindStringSubmatch(limit)

	return capture[1], capture[2], nil
}
