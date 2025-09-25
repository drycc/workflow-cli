package commands

import (
	"fmt"
	"regexp"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/config"
	"github.com/drycc/controller-sdk-go/limits"
	"github.com/drycc/workflow-cli/internal/loader"
	"github.com/drycc/workflow-cli/pkg/settings"
)

// LimitsList lists an app's limits.
func (d *DryccCmd) LimitsList(appID string, version int) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	config, err := config.List(s.Client, appID, version)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	cached := make(map[string]api.LimitPlan)

	if len(config.Limits) > 0 {
		table := d.getDefaultFormatTable([]string{"PTYPE", "PLAN", "VCPUS", "MEMORY", "FEATURES"})
		for _, ptype := range *sortKeys(config.Limits) {
			limitPlanID := fmt.Sprintf("%v", config.Limits[ptype])
			if _, ok := cached[limitPlanID]; !ok {
				limitPlan, err := limits.GetPlan(s.Client, limitPlanID)
				if err != nil {
					return err
				}
				cached[limitPlanID] = limitPlan
			}
			limitPlan := cached[limitPlanID]
			gpuCount := limitPlan.Features["gpu"]
			gpuName := limitPlan.Spec.Features["gpu"].(map[string]any)["name"]
			gpuMemory := limitPlan.Spec.Features["gpu"].(map[string]any)["memory"].(map[string]any)["size"]
			table.Append([]string{
				ptype,
				limitPlanID,
				fmt.Sprintf("%v", limitPlan.CPU),
				fmt.Sprintf("%v GiB", limitPlan.Memory),
				fmt.Sprintf("%v %v * %v", gpuName, gpuMemory, gpuCount),
			})
		}

		table.Render()
	} else {
		d.Println(fmt.Sprintf("No limits found in %s app.", appID))
	}
	return nil
}

// LimitsSet sets an app's limits.
func (d *DryccCmd) LimitsSet(appID string, limits []string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	configObj := api.Config{}
	if len(limits) > 0 {
		limitsMap, err := parseLimits(limits)
		if err != nil {
			return err
		}
		configObj.Limits = limitsMap
	}

	d.Print("Applying limits... ")

	quit := progress(d.WOut)

	_, err = config.Set(s.Client, appID, configObj, true)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Print("done\n\n")

	return d.LimitsList(appID, -1)
}

// LimitsUnset removes an app's limits.
func (d *DryccCmd) LimitsUnset(appID string, limits []string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	d.Print("Applying limits... ")

	quit := progress(d.WOut)

	configObj := api.Config{}
	if len(limits) > 0 {
		limitsMap := make(map[string]any)
		for _, limit := range limits {
			limitsMap[limit] = nil
		}
		configObj.Limits = limitsMap
	}

	_, err = config.Set(s.Client, appID, configObj, true)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Print("done\n\n")

	return d.LimitsList(appID, -1)
}

// LimitsSpecs list limit spec
func (d *DryccCmd) LimitsSpecs(keywords string, results int) error {
	s, err := settings.Load(d.ConfigFile)
	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}
	limitSpecs, count, err := limits.Specs(s.Client, keywords, results)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	if count == 0 {
		d.Println("Could not find any limit spec.")
	} else {
		table := d.getDefaultFormatTable([]string{"ID", "CPU", "CLOCK", "BOOST", "CORES", "THREADS", "NETWORK", "FEATURES"})
		for _, limitSpec := range limitSpecs {
			gpuName := limitSpec.Features["gpu"].(map[string]any)["name"]
			gpuMemory := limitSpec.Features["gpu"].(map[string]any)["memory"].(map[string]any)["size"]
			table.Append([]string{
				limitSpec.ID,
				fmt.Sprintf("%v", limitSpec.CPU["name"]),
				fmt.Sprintf("%v", limitSpec.CPU["clock"]),
				fmt.Sprintf("%v", limitSpec.CPU["boost"]),
				fmt.Sprintf("%v", limitSpec.CPU["cores"]),
				fmt.Sprintf("%v", limitSpec.CPU["threads"]),
				fmt.Sprintf("%v", limitSpec.Features["network"]),
				fmt.Sprintf("%v %v", gpuName, gpuMemory),
			})
		}
		table.Render()
	}
	return nil
}

// LimitsPlans list limit plan
func (d *DryccCmd) LimitsPlans(specID string, cpu, memory, results int) error {
	s, err := settings.Load(d.ConfigFile)
	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}
	limitPlans, count, err := limits.Plans(s.Client, specID, cpu, memory, results)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	if count == 0 {
		d.Println("Could not find any limit spec.")
	} else {
		table := d.getDefaultFormatTable([]string{"ID", "SPEC", "CPU", "VCPUS", "MEMORY", "FEATURES"})
		for _, limitPlan := range limitPlans {
			gpuName := limitPlan.Spec.Features["gpu"].(map[string]any)["name"]
			gpuMemory := limitPlan.Spec.Features["gpu"].(map[string]any)["memory"].(map[string]any)["size"]
			table.Append([]string{
				limitPlan.ID,
				limitPlan.Spec.ID,
				fmt.Sprintf("%v", limitPlan.Spec.CPU["name"]),
				fmt.Sprintf("%v", limitPlan.CPU),
				fmt.Sprintf("%v GiB", limitPlan.Memory),
				fmt.Sprintf("%v %v", gpuName, gpuMemory),
			})
		}
		table.Render()
	}
	return nil
}

func parseLimits(limits []string) (map[string]any, error) {
	limitsMap := make(map[string]any)

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
	regex := regexp.MustCompile("^([a-z0-9]+(?:-[a-z0-9]+)*)=([-.a-zA-Z0-9]+)$")

	if !regex.MatchString(limit) {
		return "", "", fmt.Errorf(`%s doesn't fit format type=#unit or type=#
Examples: web=std1.large.c1m1`, limit)
	}

	capture := regex.FindStringSubmatch(limit)

	return capture[1], capture[2], nil
}
