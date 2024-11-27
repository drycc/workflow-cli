package cmd

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/builds"
)

// BuildsInfo lists an app's builds.
func (d *DryccCmd) BuildsInfo(appID string, version int) error {
	s, appID, err := load(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	build, err := builds.Get(s.Client, appID, version)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	table := d.getDefaultFormatTable([]string{})
	table.Append([]string{"App:", build.App})
	table.Append([]string{"Sha:", build.Sha})
	table.Append([]string{"UUID:", build.UUID})
	table.Append([]string{"Owner:", build.Owner})
	table.Append([]string{"Image:", build.Image})
	table.Append([]string{"Stack:", build.Stack})
	table.Append([]string{"Created:", build.Created})
	table.Append([]string{"Updated:", build.Updated})
	table.Render()

	if len(build.Dockerfile) != 0 {
		table = d.getDefaultFormatTable([]string{})
		table.Append([]string{"Dockerfile:"})
		table.Append([]string{d.indentString(build.Dockerfile, 2)})
		table.Render()
	}

	if len(build.Procfile) != 0 {
		table = d.getDefaultFormatTable([]string{})
		table.Append([]string{"Procfile:"})
		table.Append([]string{d.indentString(d.toYamlString(build.Procfile, 2), 2)})
		table.Render()
	}

	if len(build.Dryccfile) != 0 {
		table = d.getDefaultFormatTable([]string{})
		table.Append([]string{"Dryccfile:"})
		table.Append([]string{d.indentString(d.toYamlString(build.Dryccfile, 2), 2)})
		table.Render()
	}
	return nil
}

// BuildsCreate creates a build for an app.
func (d *DryccCmd) BuildsCreate(appID, image, stack, procfile, dryccfile, confirm string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	procfileMap := make(map[string]string)
	if _, err := os.Stat(procfile); err == nil {
		contents, err := os.ReadFile(procfile)
		if err != nil {
			return err
		}

		if procfileMap, err = parseProcfile(contents); err != nil {
			return err
		}
	}

	dryccfileMap := make(map[string]interface{})
	if _, err := os.Stat(dryccfile); err == nil {
		contents, err := os.ReadFile(dryccfile)
		if err != nil {
			return err
		}

		if dryccfileMap, err = parseDryccfile(contents); err != nil {
			return err
		}
	}
	// check procfileMap dryccfileMap stack is exist
	err = buildConfirmAction(s.Client, appID, procfileMap, dryccfileMap, confirm)
	if err != nil {
		return err
	}
	d.Print("Creating build... ")
	quit := progress(d.WOut)
	_, err = builds.New(s.Client, appID, image, stack, procfileMap, dryccfileMap)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")

	return nil
}

func parseProcfile(procfile []byte) (map[string]string, error) {
	procfileMap := make(map[string]string)
	return procfileMap, yaml.Unmarshal(procfile, &procfileMap)
}

func parseDryccfile(dryccfile []byte) (map[string]interface{}, error) {
	dryccfileMap := make(map[string]interface{})
	return dryccfileMap, yaml.Unmarshal(dryccfile, &dryccfileMap)
}

func buildConfirmAction(c *drycc.Client, appID string, procfileMap map[string]string,
	dryccfileMap map[string]interface{}, confirm string) error {

	build, _ := builds.Get(c, appID, 0)

	if ((len(build.Procfile) != 0 && len(procfileMap) == 0) || (len(build.Dryccfile) != 0 && len(dryccfileMap) == 0)) && (confirm == "" || confirm != "yes") {
		// hint
		fmt.Printf(` !    WARNING: Potentially Build Create Action
 !    The Procfile or drycc.yaml is empty, not last time
 !    To proceed, type "yes" !

> `)

		fmt.Scanln(&confirm)
		if confirm != "yes" {
			return fmt.Errorf("cancel the build create action")
		}
		return nil
	}
	return nil
}
