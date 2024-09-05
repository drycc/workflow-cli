package cmd

import (
	"fmt"
	"sort"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/config"
)

func getHealthcheckString(ptype, probeType string, healthcheck *api.Healthcheck) string {
	params := fmt.Sprintf(
		"delay=%ds timeout=%ds period=%ds #success=%d #failure=%d",
		healthcheck.InitialDelaySeconds,
		healthcheck.TimeoutSeconds,
		healthcheck.PeriodSeconds,
		healthcheck.SuccessThreshold,
		healthcheck.FailureThreshold,
	)

	if healthcheck.Exec != nil {
		return fmt.Sprintf("%s %s exec %v %s", probeType, ptype, healthcheck.Exec.Command, params)
	} else if healthcheck.TCPSocket != nil {
		return fmt.Sprintf("%s %s tcp-socket port=%v %s", probeType, ptype, healthcheck.TCPSocket.Port, params)
	} else if healthcheck.HTTPGet != nil {
		return fmt.Sprintf(
			"%s %s http-get headers=%v path=%s port=%d %s",
			probeType,
			ptype,
			healthcheck.HTTPGet.HTTPHeaders,
			healthcheck.HTTPGet.Path,
			healthcheck.HTTPGet.Port,
			params,
		)
	}
	return ""
}

func getHealthchecksStrings(ptype string, healthchecks *api.Healthchecks) []string {
	var probes []string
	for key := range *healthchecks {
		probes = append(probes, getHealthcheckString(ptype, key, (*healthchecks)[key]))
	}
	return probes
}

// HealthchecksList lists an app's healthchecks.
func (d *DryccCmd) HealthchecksList(appID, ptype string) error {
	s, appID, err := load(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	config, err := config.List(s.Client, appID)

	if err != nil {
		return err
	}

	if ptype == "" {
		if len(config.Healthcheck) == 0 {
			d.Println("No health checks configured.")
		} else {
			var keys []string
			for k := range config.Healthcheck {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			table := d.getDefaultFormatTable([]string{})
			table.Append([]string{"App:", config.App})
			table.Append([]string{"UUID:", config.UUID})
			table.Append([]string{"Owner:", config.Owner})
			table.Append([]string{"Created:", d.formatTime(config.Created)})
			table.Append([]string{"Updated:", d.formatTime(config.Updated)})
			table.Append([]string{"Healthchecks:"})
			for _, key := range keys {
				for _, probe := range getHealthchecksStrings(key, config.Healthcheck[key]) {
					if probe != "" {
						table.Append([]string{"", probe})
					}
				}
			}
			table.Render()
		}
	} else {
		if healthcheck, found := config.Healthcheck[ptype]; found {
			table := d.getDefaultFormatTable([]string{})
			table.Append([]string{"App:", config.App})
			table.Append([]string{"UUID:", config.UUID})
			table.Append([]string{"Owner:", config.Owner})
			table.Append([]string{"Created:", d.formatTime(config.Created)})
			table.Append([]string{"Updated:", d.formatTime(config.Updated)})
			table.Append([]string{"Healthchecks:"})
			for _, probe := range getHealthchecksStrings(ptype, healthcheck) {
				if probe != "" {
					table.Append([]string{"", probe})
				}
			}
			table.Render()
		} else {
			d.Println("No health checks configured.")
		}
	}
	return nil
}

// HealthchecksSet sets an app's healthchecks.
func (d *DryccCmd) HealthchecksSet(appID, healthcheckType, ptype string, probe *api.Healthcheck) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Applying %s healthcheck... ", healthcheckType)

	quit := progress(d.WOut)

	healthcheckMap := make(api.Healthchecks)
	healthcheckMap[healthcheckType] = probe
	configObj := api.Config{Healthcheck: make(map[string]*api.Healthchecks)}
	configObj.Healthcheck[ptype] = &healthcheckMap

	_, err = config.Set(s.Client, appID, configObj)

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	d.Print("done\n\n")

	return d.HealthchecksList(appID, ptype)
}

// HealthchecksUnset removes an app's healthchecks.
func (d *DryccCmd) HealthchecksUnset(appID, ptype string, healthchecks []string) error {
	s, appID, err := load(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	d.Print("Removing healthchecks... ")

	quit := progress(d.WOut)

	configObj := api.Config{}

	healthchecksMap := make(map[string]*api.Healthchecks)
	healthcheckMap := make(api.Healthchecks)

	for _, healthcheck := range healthchecks {
		healthcheckMap[healthcheck] = nil
	}
	healthchecksMap[ptype] = &healthcheckMap

	configObj.Healthcheck = healthchecksMap

	_, err = config.Set(s.Client, appID, configObj)

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	d.Print("done\n\n")

	return d.HealthchecksList(appID, ptype)
}
