package cmd

import (
	"sort"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/config"
)

func (d *DryccCmd) printHealthCheck(healthcheck api.Healthchecks) {
	d.Println("--- Liveness")
	if livenessProbe, found := healthcheck["livenessProbe"]; found {
		d.Println(livenessProbe)
	} else {
		d.Println("No liveness probe configured.")
	}

	d.Println("\n--- Readiness")
	if readinessProbe, found := healthcheck["readinessProbe"]; found {
		d.Println(readinessProbe)
	} else {
		d.Println("No readiness probe configured.")
	}
}

// HealthchecksList lists an app's healthchecks.
func (d *DryccCmd) HealthchecksList(appID, procType string) error {
	s, appID, err := load(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	config, err := config.List(s.Client, appID)

	if err != nil {
		return err
	}

	d.Printf("=== %s Healthchecks\n", appID)
	if procType == "" {
		if len(config.Healthcheck) == 0 {
			d.Println("No health checks configured.")
			return nil
		}
		var keys []string
		for k := range config.Healthcheck {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, key := range keys {
			d.Printf("\n%s:\n", key)
			d.printHealthCheck(*config.Healthcheck[key])
		}
	} else {
		d.Printf("\n%s:\n", procType)
		if healthcheck, found := config.Healthcheck[procType]; found {
			d.printHealthCheck(*healthcheck)
		} else {
			d.printHealthCheck(api.Healthchecks{})
		}
	}

	return nil
}

// HealthchecksSet sets an app's healthchecks.
func (d *DryccCmd) HealthchecksSet(appID, healthcheckType, procType string, probe *api.Healthcheck) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Applying %s healthcheck... ", healthcheckType)

	quit := progress(d.WOut)

	healthcheckMap := make(api.Healthchecks)
	healthcheckMap[healthcheckType] = probe
	configObj := api.Config{Healthcheck: make(map[string]*api.Healthchecks)}
	configObj.Healthcheck[procType] = &healthcheckMap

	_, err = config.Set(s.Client, appID, configObj)

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	d.Print("done\n\n")

	return d.HealthchecksList(appID, procType)
}

// HealthchecksUnset removes an app's healthchecks.
func (d *DryccCmd) HealthchecksUnset(appID, procType string, healthchecks []string) error {
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
	healthchecksMap[procType] = &healthcheckMap

	configObj.Healthcheck = healthchecksMap

	_, err = config.Set(s.Client, appID, configObj)

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	d.Print("done\n\n")

	return d.HealthchecksList(appID, procType)
}
