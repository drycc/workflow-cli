package cmd

import (
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"

	"github.com/drycc/workflow-cli/version"
	"github.com/minio/selfupdate"
)

const WorkflowCliURL = "https://www.drycc.cc/workflow-cli.txt"

func (d *DryccCmd) latestVersion() (string, string, error) {
	d.Print("Get the latest version of workflow cli... ")
	quit := progress(d.WOut)
	resp, err := http.Get(WorkflowCliURL)
	quit <- true
	<-quit
	if err != nil {
		return "", "", err
	}
	d.Println("done")

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	prefix := "drycc-"
	suffix := fmt.Sprintf("-%s-%s", runtime.GOOS, runtime.GOARCH)
	for _, url := range strings.Split(string(body), "\n") {
		if strings.HasSuffix(url, suffix) {
			names := strings.Split(url, "/")
			name := names[len(names)-1]
			version := strings.ReplaceAll(strings.ReplaceAll(name, suffix, ""), prefix, "")
			return version, url, nil
		}
	}
	return "", "", fmt.Errorf("unable to obtain version: %s, %s", runtime.GOOS, runtime.GOARCH)
}

// Update workflow-cli to latest release
func (d *DryccCmd) Update(dryRun bool) error {
	latestVersion, downloadURL, err := d.latestVersion()
	if err != nil {
		return err
	}
	if latestVersion != version.Version {
		d.Printf("Update workflow cli from %s to %s... ", version.Version, latestVersion)
		if dryRun {
			d.Println("skip")
		} else {
			resp, err := http.Get(downloadURL)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			quit := progress(d.WOut)
			err = selfupdate.Apply(resp.Body, selfupdate.Options{})
			quit <- true
			<-quit
			if err != nil {
				return err
			}
			d.Println("done")
		}
	} else {
		d.Println("You are already running the most recent version.")
	}
	return nil
}
