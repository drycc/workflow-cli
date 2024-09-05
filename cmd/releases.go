package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/drycc/controller-sdk-go/releases"
)

// ReleasesList lists an app's releases.
func (d *DryccCmd) ReleasesList(appID string, results int) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}

	releases, count, err := releases.List(s.Client, appID, results)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	if count == 0 {
		d.Println(fmt.Sprintf("No releases found in %s app.", appID))
	} else {
		table := d.getDefaultFormatTable([]string{"OWNER", "STATE", "VERSION", "CREATED", "SUMMARY"})
		for _, r := range releases {
			summary := r.Summary
			if len(summary) > 64 {
				summary = fmt.Sprintf("%s[...]", summary[:64])
			}
			table.Append([]string{r.Owner, r.State, fmt.Sprintf("v%d", r.Version), d.formatTime(r.Created), summary})
		}
		table.Render()
	}
	return nil
}

// ReleasesInfo prints info about a specific release.
func (d *DryccCmd) ReleasesInfo(appID string, version int) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	r, err := releases.Get(s.Client, appID, version)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	table := d.getDefaultFormatTable([]string{})
	table.Append([]string{"App:", r.App})
	table.Append([]string{"UUID:", r.UUID})
	table.Append([]string{"State:", r.State})
	table.Append([]string{"Owner:", r.Owner})
	table.Append([]string{"Build:", r.Build})
	table.Append([]string{"Config:", r.Config})
	table.Append([]string{"Created:", d.formatTime(r.Created)})
	table.Append([]string{"Updated:", d.formatTime(r.Updated)})
	table.Append([]string{"Summary:", d.wrapString(r.Summary)})
	if r.Exception != "" {
		table.Append([]string{"Exception:", d.wrapString(r.Exception)})
	}
	table.Append([]string{"Version:", fmt.Sprintf("v%v", r.Version)})
	table.Render()
	// Conditions
	if len(r.Conditions) != 0 {
		// table event
		te := d.getDefaultFormatTable([]string{})
		te.Append([]string{"Conditions:"})
		for _, c := range r.Conditions {
			ptypes := "<none>"
			if len(c.Ptypes) != 0 {
				ptypes = strings.Join(c.Ptypes, ",")
			}
			exception := "<none>"
			if c.Exception != "" {
				exception = c.Exception
			}

			te.Append([]string{"  - created:", d.formatTime(c.Created)})
			te.Append([]string{"    state:", c.State})
			te.Append([]string{"    action:", c.Action})
			te.Append([]string{"    ptypes:", d.wrapString(ptypes)})
			te.Append([]string{"    exception:", d.wrapString(exception)})
		}
		te.Render()
	}
	return nil
}

// ReleasesDeploy force deploy lastest release.
func (d *DryccCmd) ReleasesDeploy(appID string, targets []string, force bool, confirm string) error {
	s, appID, err := load(d.ConfigFile, appID)
	if err != nil {
		return err
	}
	if len(targets) == 0 && (confirm == "" || confirm != "yes") {
		d.Printf(` !    WARNING: Potentially Deploy Action
 !    This command will deploy all processes of the application ptype
 !    To proceed, type "yes" !

> `)

		fmt.Scanln(&confirm)
		if confirm != "yes" {
			return fmt.Errorf("cancel the deploy action")
		}
	}
	d.Printf("Deploying process types... but first, %s!\n", drinkOfChoice())
	startTime := time.Now()
	quit := progress(d.WOut)
	ptypes := strings.Join(targets, ",")
	targetMap := map[string]interface{}{
		"ptypes": ptypes,
		"force":  force,
	}
	err = releases.Deploy(s.Client, appID, targetMap)
	quit <- true
	<-quit
	if err != nil {
		return err
	}

	d.Printf("done in %ds\n", int(time.Since(startTime).Seconds()))
	return nil
}

// ReleasesRollback rolls an app back to a previous release.
func (d *DryccCmd) ReleasesRollback(appID string, targets []string, version int) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}
	ptypes := strings.Join(targets, ",")
	if version == -1 {
		d.Print("Rolling back one release... ")
	} else {
		d.Printf("Rolling back to v%d... ", version)
	}

	quit := progress(d.WOut)
	newVersion, err := releases.Rollback(s.Client, appID, ptypes, version)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Printf("done, v%d\n", newVersion)

	return nil
}
