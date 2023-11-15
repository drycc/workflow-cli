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
		table := d.getDefaultFormatTable([]string{"UUID", "OWNER", "SHA", "CREATED"})
		for _, build := range builds {
			table.Append([]string{
				build.UUID,
				safeGetString(build.Owner),
				safeGetString(build.Sha),
				build.Created,
			})
		}
		table.Render()
	} else {
		d.Println(fmt.Sprintf("No builds found in %s app.", appID))
	}
	return nil
}

// BuildsCreate creates a build for an app.
func (d *DryccCmd) BuildsCreate(appID, image, stack, procfile string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	procfileMap := make(map[string]string)

	if procfile != "" {
		if procfileMap, err = parseProcfile([]byte(procfile)); err != nil {
			return err
		}
	} else if _, err := os.Stat("Procfile"); err == nil {
		contents, err := os.ReadFile("Procfile")
		if err != nil {
			return err
		}

		if procfileMap, err = parseProcfile(contents); err != nil {
			return err
		}
	}

	d.Print("Creating build... ")
	quit := progress(d.WOut)
	_, err = builds.New(s.Client, appID, image, stack, procfileMap)
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
