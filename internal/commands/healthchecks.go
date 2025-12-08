package commands

import (
	"fmt"
	"sort"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/config"
	"github.com/drycc/workflow-cli/internal/loader"
)

func getContainerProbeString(ptype, probeType string, containerProbe *api.ContainerProbe) string {
	params := fmt.Sprintf(
		"delay=%ds timeout=%ds period=%ds #success=%d #failure=%d",
		containerProbe.InitialDelaySeconds,
		containerProbe.TimeoutSeconds,
		containerProbe.PeriodSeconds,
		containerProbe.SuccessThreshold,
		containerProbe.FailureThreshold,
	)

	if containerProbe.Exec != nil {
		return fmt.Sprintf("%s %s exec %v %s", probeType, ptype, containerProbe.Exec.Command, params)
	} else if containerProbe.TCPSocket != nil {
		return fmt.Sprintf("%s %s tcp-socket port=%v %s", probeType, ptype, containerProbe.TCPSocket.Port, params)
	} else if containerProbe.HTTPGet != nil {
		return fmt.Sprintf(
			"%s %s http-get headers=%v path=%s port=%d %s",
			probeType,
			ptype,
			containerProbe.HTTPGet.HTTPHeaders,
			containerProbe.HTTPGet.Path,
			containerProbe.HTTPGet.Port,
			params,
		)
	}
	return ""
}

func getHealthchecksStrings(ptype string, healthcheck *api.Healthcheck) []string {
	var containerProbes []string
	if healthcheck.StartupProbe != nil {
		containerProbes = append(containerProbes, getContainerProbeString(ptype, "startupProbe", *healthcheck.StartupProbe))
	} else if healthcheck.LivenessProbe != nil {
		containerProbes = append(containerProbes, getContainerProbeString(ptype, "livenessProbe", *healthcheck.LivenessProbe))
	} else if healthcheck.ReadinessProbe != nil {
		containerProbes = append(containerProbes, getContainerProbeString(ptype, "readinessProbe", *healthcheck.ReadinessProbe))
	}
	return containerProbes
}

// HealthchecksList lists an app's healthchecks.
func (d *DryccCmd) HealthchecksList(appID, ptype string, version int) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	config, err := config.List(s.Client, appID, version)
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
func (d *DryccCmd) HealthchecksSet(appID, healthcheckType, ptype string, probe *api.ContainerProbe) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	d.Printf("Applying %s healthcheck... ", healthcheckType)

	quit := progress(d.WOut)
	var healthcheck api.Healthcheck
	switch healthcheckType {
	case "livenessProbe":
		healthcheck.LivenessProbe = &probe
	case "readinessProbe":
		healthcheck.ReadinessProbe = &probe
	case "startupProbe":
		healthcheck.StartupProbe = &probe
	default:
		return fmt.Errorf("unknown healthcheck type: %s", healthcheckType)
	}
	configObj := api.Config{Healthcheck: make(map[string]*api.Healthcheck)}
	configObj.Healthcheck[ptype] = &healthcheck

	_, err = config.Set(s.Client, appID, configObj, true)

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	d.Print("done\n\n")

	return d.HealthchecksList(appID, ptype, -1)
}

// HealthchecksUnset removes an app's healthchecks.
func (d *DryccCmd) HealthchecksUnset(appID, ptype string, containerProbeTypes []string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	d.Print("Removing healthchecks... ")

	quit := progress(d.WOut)

	configObj := api.Config{}

	healthcheckMap := make(map[string]*api.Healthcheck)
	var nullContainerProbe *api.ContainerProbe = nil
	for _, containerProbeType := range containerProbeTypes {
		healthcheck := &api.Healthcheck{}
		switch containerProbeType {
		case "livenessProbe":
			healthcheck.LivenessProbe = &nullContainerProbe
		case "readinessProbe":
			healthcheck.ReadinessProbe = &nullContainerProbe
		case "startupProbe":
			healthcheck.StartupProbe = &nullContainerProbe
		default:
			return fmt.Errorf("unknown container probe type: %s", containerProbeType)
		}
		healthcheckMap[containerProbeType] = healthcheck
	}
	configObj.Healthcheck = healthcheckMap
	_, err = config.Set(s.Client, appID, configObj, true)

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	d.Print("done\n\n")

	return d.HealthchecksList(appID, ptype, -1)
}
