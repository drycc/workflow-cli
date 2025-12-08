package commands

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/events"
	"github.com/drycc/controller-sdk-go/pts"
	"github.com/drycc/workflow-cli/internal/loader"
)

// PtsList lists an app's processes.
func (d *DryccCmd) PtsList(appID string, results int) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
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
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
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
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
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
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
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
		"ptypes": ptypes,
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

// PtsClean cleans process types that are not used.
func (d *DryccCmd) PtsClean(appID string, targets []string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}
	d.Printf("Cleaning process types... but first, %s!\n", drinkOfChoice())
	startTime := time.Now()
	quit := progress(d.WOut)
	ptypes := strings.Join(targets, ",")
	targetMap := map[string]string{
		"ptypes": ptypes,
	}
	err = pts.Clean(s.Client, appID, targetMap)
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
		table := d.getDefaultFormatTable([]string{"NAME", "RELEASE", "READY", "UP-TO-DATE", "AVAILABLE", "STARTED", "GARBAGE"})
		for _, pt := range pts {
			table.Append([]string{
				pt.Name,
				pt.Release,
				pt.Ready,
				fmt.Sprintf("%v", pt.UpToDate),
				fmt.Sprintf("%v", pt.AvailableReplicas),
				d.formatTime(pt.Started),
				fmt.Sprintf("%t", pt.Garbage),
			})
		}
		table.Render()
	}
}

func printProcessTypeDetail(d *DryccCmd, ptypeStates api.PtypeStates, events api.AppEvents) {
	// table process type
	tpt := d.getDefaultFormatTable([]string{})
	for _, ptypeState := range ptypeStates {
		tpt.Append([]string{"Container:", ptypeState.Container})
		tpt.Append([]string{"Image:", ptypeState.Image})
		if len(ptypeState.Command) != 0 {
			tpt.Append([]string{"Command:"})
			for _, command := range ptypeState.Command {
				tpt.Append([]string{"", fmt.Sprintf("- %v", command)})
			}
		}
		if len(ptypeState.Args) != 0 {
			tpt.Append([]string{"Args:"})
			for _, arg := range ptypeState.Args {
				tpt.Append([]string{"", fmt.Sprintf("- %v", arg)})
			}
		}
		if len(ptypeState.Limits) != 0 {
			tpt.Append([]string{"Limits:"})
			for r, q := range ptypeState.Limits {
				tpt.Append([]string{"", fmt.Sprintf("%s %s", r, q)})
			}
		}
		if len(ptypeState.VolumeMounts) != 0 {
			tpt.Append([]string{"Mounts:"})
			for _, mount := range ptypeState.VolumeMounts {
				tpt.Append([]string{"", fmt.Sprintf("%s from %s", mount.MountPath, mount.Name)})
			}
		}
		sp := getContainerProbeString("", "", &ptypeState.StartupProbe)
		if sp != "" {
			tpt.Append([]string{"Startup:", strings.TrimSpace(sp)})
		}
		lp := getContainerProbeString("", "", &ptypeState.LivenessProbe)
		if lp != "" {
			tpt.Append([]string{"Liveness:", strings.TrimSpace(lp)})
		}
		rp := getContainerProbeString("", "", &ptypeState.ReadinessProbe)
		if rp != "" {
			tpt.Append([]string{"Readiness:", strings.TrimSpace(rp)})
		}
		if len(ptypeState.NodeSelector) != 0 {
			tpt.Append([]string{"Node-Selectors:"})
			for k, v := range ptypeState.NodeSelector {
				tpt.Append([]string{"", fmt.Sprintf("%s=%s", k, v)})
			}
		}

	}
	tpt.Render()

	if len(events) != 0 {
		// table event
		te := d.getDefaultFormatTable([]string{})
		te.Append([]string{"Events:"})
		te.Append([]string{d.indentString("REASON", 2), "MESSAGE", "CREATED"})
		for _, ev := range events {
			te.Append([]string{
				d.indentString(ev.Reason, 2),
				ev.Message,
				d.formatTime(ev.Created),
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
			return nil, fmt.Errorf("'%s' does not match the pattern 'ptype=num', ex: web=2", target)
		}
	}

	return targetMap, nil
}
