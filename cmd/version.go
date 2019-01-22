package cmd

import (
	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/workflow-cli/settings"
	"github.com/drycc/workflow-cli/version"
)

// Version prints the various CLI versions.
func (d *DryccCmd) Version(all bool) error {
	if !all {
		d.Println(version.Version)
		return nil
	}

	d.Printf("Workflow CLI Version:            %s\n", version.Version)
	d.Printf("Workflow CLI API Version:        %s\n", drycc.APIVersion)

	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	// retrive version information from drycc controller
	err = s.Client.Healthcheck()

	if err != nil && err != drycc.ErrAPIMismatch {
		return err
	}

	d.Printf("Workflow Controller API Version: %s\n", s.Client.ControllerAPIVersion)
	return nil
}
