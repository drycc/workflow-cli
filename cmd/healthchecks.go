package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/deis/controller-sdk-go/api"
	"github.com/deis/controller-sdk-go/config"
)

func printHealthCheck(out io.Writer, healthcheck api.Healthchecks) {
	fmt.Fprintln(out, "--- Liveness")
	if livenessProbe, found := healthcheck["livenessProbe"]; found {
		fmt.Fprintln(out, livenessProbe)
	} else {
		fmt.Fprintln(out, "No liveness probe configured.")
	}

	fmt.Fprintln(out, "\n--- Readiness")
	if readinessProbe, found := healthcheck["readinessProbe"]; found {
		fmt.Fprintln(out, readinessProbe)
	} else {
		fmt.Fprintln(out, "No readiness probe configured.")
	}
}

// HealthchecksList lists an app's healthchecks.
func HealthchecksList(appID, procType string) error {
	s, appID, err := load(appID)

	if err != nil {
		return err
	}

	config, err := config.List(s.Client, appID)

	if err != nil {
		return err
	}

	fmt.Printf("=== %s Healthchecks\n\n", appID)
	if procType == "" {
		for proc, healthcheck := range config.Healthcheck {
			fmt.Println(proc + ":")
			printHealthCheck(os.Stdout, *healthcheck)
		}
	} else {
		fmt.Println(procType + ":")
		if healthcheck, found := config.Healthcheck[procType]; found {
			printHealthCheck(os.Stdout, *healthcheck)
		} else {
			return errors.New(appID + " doesn't have proctype" + procType)
		}
	}

	return nil
}

// HealthchecksSet sets an app's healthchecks.
func HealthchecksSet(appID, healthcheckType, procType string, probe *api.Healthcheck) error {
	s, appID, err := load(appID)

	if err != nil {
		return err
	}

	fmt.Printf("Applying %s healthcheck... ", healthcheckType)

	quit := progress()

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

	fmt.Print("done\n\n")

	return HealthchecksList(appID, procType)
}

// HealthchecksUnset removes an app's healthchecks.
func HealthchecksUnset(appID, procType string, healthchecks []string) error {
	s, appID, err := load(appID)

	if err != nil {
		return err
	}

	fmt.Print("Removing healthchecks... ")

	quit := progress()

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

	fmt.Print("done\n\n")

	return HealthchecksList(appID, procType)
}
