package cmd

import (
	deis "github.com/teamhephy/controller-sdk-go"
	"github.com/teamhephy/workflow-cli/settings"
	"github.com/teamhephy/workflow-cli/version"
)

// Version prints the various CLI versions.
func (d *DeisCmd) Version(all bool) error {
	if !all {
		d.Println(version.Version)
		return nil
	}

	d.Printf("Workflow CLI Version:            %s\n", version.Version)
	d.Printf("Workflow CLI API Version:        %s\n", deis.APIVersion)

	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	// retrive version information from deis controller
	err = s.Client.Healthcheck()

	if err != nil && err != deis.ErrAPIMismatch {
		return err
	}

	d.Printf("Workflow Controller API Version: %s\n", s.Client.ControllerAPIVersion)
	return nil
}
