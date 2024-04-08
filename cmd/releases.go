package cmd

import (
	"fmt"

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
		table := d.getDefaultFormatTable([]string{"UUID", "OWNER", "STATE", "VERSION", "CREATED", "SUMMARY"})
		for _, r := range releases {
			summary := r.Summary
			if len(summary) > 64 {
				summary = fmt.Sprintf("%s[...]", summary[:64])
			}
			table.Append([]string{r.UUID, r.Owner, r.State, fmt.Sprintf("v%d", r.Version), r.Created, summary})
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
	table.Append([]string{"Created:", r.Created})
	table.Append([]string{"Updated:", r.Updated})
	table.Append([]string{"Summary:", d.wrapString(r.Summary)})
	table.Append([]string{"Version:", fmt.Sprintf("v%v", r.Version)})
	table.Render()
	return nil
}

// ReleasesRollback rolls an app back to a previous release.
func (d *DryccCmd) ReleasesRollback(appID string, version int) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	if version == -1 {
		d.Print("Rolling back one release... ")
	} else {
		d.Printf("Rolling back to v%d... ", version)
	}

	quit := progress(d.WOut)
	newVersion, err := releases.Rollback(s.Client, appID, version)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Printf("done, v%d\n", newVersion)

	return nil
}
