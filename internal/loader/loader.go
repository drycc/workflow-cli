// Package loader provides common utility functions and helper methods for the workflow CLI.
package loader

import (
	"github.com/drycc/workflow-cli/pkg/git"
	"github.com/drycc/workflow-cli/pkg/settings"
)

// LoadAppSettings loads settings file and looks up the app name
func LoadAppSettings(cf string, appID string) (string, *settings.Settings, error) {
	s, err := settings.Load(cf)
	if err != nil {
		return "", nil, err
	}

	if appID == "" {
		appID, err = git.DetectAppName(git.DefaultCmd, s.Client.ControllerURL.Host)
		if err != nil {
			return "", nil, err
		}
	}

	return appID, s, nil
}
