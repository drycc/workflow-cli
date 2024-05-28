package cmd

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"

	"github.com/drycc/controller-sdk-go/builds"
)

// BuildsList lists an app's builds.
func (d *DryccCmd) BuildsList(appID string, results int) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}

	builds, count, err := builds.List(s.Client, appID, results)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	if count > 0 {
		table := d.getDefaultFormatTable([]string{"OWNER", "SHA", "CREATED"})
		for _, build := range builds {
			table.Append([]string{
				safeGetString(build.Owner),
				safeGetString(build.Sha),
				d.formatTime(build.Created),
			})
		}
		table.Render()
	} else {
		d.Println(fmt.Sprintf("No builds found in %s app.", appID))
	}
	return nil
}

// BuildsCreate creates a build for an app.
func (d *DryccCmd) BuildsCreate(appID, image, stack, procfile string, dryccfile string) error {
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
