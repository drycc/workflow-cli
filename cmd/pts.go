package cmd

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/events"
	"github.com/drycc/controller-sdk-go/pts"
)

// PtsList lists an app's processes.
func (d *DryccCmd) PtsList(appID string, results int) error {
	s, appID, err := load(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}

	ptypes, _, err := pts.List(s.Client, appID, results)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	printProcessTypes(d, appID, ptypes)

	return nil
}

// PtsDescribe describe an app's processes.
func (d *DryccCmd) PtsDescribe(appID, ptype string) error {
	s, appID, err := load(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	ptypeStates, _, err := pts.Describe(s.Client, appID, ptype, 1000)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	events, _, err := events.ListPtypeEvents(s.Client, appID, ptype, 1000)
	if err != nil {
		return err
	}
	printProcessTypeDetail(d, ptypeStates, events)
	return nil
}

// PtsScale scales an app's processes.
func (d *DryccCmd) PtsScale(appID string, targets []string) error {
	s, appID, err := load(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	targetMap, err := parsePtsTargets(targets)
	if err != nil {
		return err
	}

	d.Printf("Scaling process types... but first, %s!\n", drinkOfChoice())
	startTime := time.Now()
	quit := progress(d.WOut)

	err = pts.Scale(s.Client, appID, targetMap)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Printf("done in %ds\n", int(time.Since(startTime).Seconds()))

	return nil
}

// PtsRestart restarts an app's processes.
func (d *DryccCmd) PtsRestart(appID string, targets []string, confirm string) error {
	s, appID, err := load(d.ConfigFile, appID)
	if err != nil {
		return err
	}
	if len(targets) == 0 && (confirm == "" || confirm != "yes") {
		d.Printf(` !    WARNING: Potentially Restart Action
 !    This command will restart all processes of the application ptype
 !    To proceed, type "yes" !

> `)

		fmt.Scanln(&confirm)
		if confirm != "yes" {
			return fmt.Errorf("cancel the restart action")
		}
	}
	d.Printf("Restarting process types... but first, %s!\n", drinkOfChoice())
	startTime := time.Now()
	quit := progress(d.WOut)
	ptypes := strings.Join(targets, ",")
	targetMap := map[string]string{
		"types": ptypes,
	}
	err = pts.Restart(s.Client, appID, targetMap)
	quit <- true
	<-quit
	if err != nil {
		return err
	}

	d.Printf("done in %ds\n", int(time.Since(startTime).Seconds()))
	return nil
}

func printProcessTypes(d *DryccCmd, appID string, ptypes api.Ptypes) {
	pts := pts.ByType(ptypes)

	if len(pts) == 0 {
		d.Println(fmt.Sprintf("No processes found in %s app.", appID))
	} else {
		table := d.getDefaultFormatTable([]string{"NAME", "RELEASE", "READY", "UP-TO-DATE", "AVAILABLE", "STARTED"})
		for _, pt := range pts {
			table.Append([]string{
				pt.Name,
				pt.Release,
				pt.Ready,
				fmt.Sprintf("%v", pt.UpToDate),
				fmt.Sprintf("%v", pt.AvailableReplicas),
				pt.Started.Format("2006-01-02T15:04:05MST"),
			})
		}
		table.Render()
	}
}

func printProcessTypeDetail(d *DryccCmd, ptypeStates api.PtypeStates, events api.AppEvents) {
	// table process type
	tpt := d.getDefaultFormatTable([]string{})
	for _, containerState := range ptypeStates {
		// tpt.Append([]string{"Container: " + containerState.Container})
		// tpt.Append([]string{"Image: " + containerState.Image})
		tpt.Append([]string{"Container:", containerState.Container})
		tpt.Append([]string{"Image:", containerState.Image})
		if len(containerState.Command) != 0 {
			tpt.Append([]string{"Command:"})
			for _, command := range containerState.Command {
				tpt.Append([]string{"", fmt.Sprintf("- %v", command)})
			}
		}
		if len(containerState.Args) != 0 {
			tpt.Append([]string{"Args:"})
			for _, arg := range containerState.Args {
				tpt.Append([]string{"", fmt.Sprintf("- %v", arg)})
			}
		}
		if containerState.Limits != nil {
			tpt.Append([]string{"Limits:"})
			for r, q := range containerState.Limits {
				tpt.Append([]string{"", fmt.Sprintf("%s %s", r, q)})
			}
		}
		if len(containerState.VolumeMounts) != 0 {
			tpt.Append([]string{"Mounts:"})
			for _, mount := range containerState.VolumeMounts {
				tpt.Append([]string{"", fmt.Sprintf("%s from %s", mount.MountPath, mount.Name)})
			}
		}
		lp := getHealthcheckString("", "", &containerState.LivenessProbe)
		if lp != "" {
			tpt.Append([]string{"Liveness:", strings.TrimSpace(lp)})
		}
		rp := getHealthcheckString("", "", &containerState.ReadinessProbe)
		if rp != "" {
			tpt.Append([]string{"Readiness:", strings.TrimSpace(rp)})
		}

	}
	tpt.Render()

	if len(events) != 0 {
		// table event
		te := d.getDefaultFormatTable([]string{})
		te.Append([]string{"Events:"})
		te.Append([]string{"  REASON", "MESSAGE", "CREATED"})
		for _, ev := range events {
			te.Append([]string{
				fmt.Sprintf("  %s", ev.Reason),
				ev.Message,
				ev.Created.Format("2006-01-02T15:04:05MST"),
			})
		}
		te.Render()
	}
}

func parsePtsTargets(targets []string) (map[string]int, error) {
	targetMap := make(map[string]int)
	regex := regexp.MustCompile(`^([a-z0-9]+(?:-[a-z0-9]+)*)=([0-9]+)$`)
	var err error

	for _, target := range targets {
		if regex.MatchString(target) {
			captures := regex.FindStringSubmatch(target)
			targetMap[captures[1]], err = strconv.Atoi(captures[2])

			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("'%s' does not match the pattern 'type=num', ex: web=2", target)
		}
	}

	return targetMap, nil
}
