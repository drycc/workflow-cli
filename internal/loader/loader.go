// Package loader provides common utility functions and helper methods for the workflow CLI.
package loader

import (
	"fmt"

	"github.com/drycc/workflow-cli/pkg/git"
	"github.com/drycc/workflow-cli/pkg/settings"
)

// LoadAppSettings loads settings file, validates workspace, and looks up the app name
func LoadAppSettings(cf string, appID string) (string, *settings.Settings, error) {
	s, err := settings.Load(cf)
	if err != nil {
		return "", nil, err
	}

	if s.Workspace == "" {
		return "", nil, fmt.Errorf("no workspace specified, set a default workspace with 'drycc workspaces switch'")
	}

	if appID == "" {
		appID, err = git.DetectAppName(git.DefaultCmd, s.Client.ControllerURL.Host)
		if err != nil {
			return "", nil, err
		}
	}

	return appID, s, nil
}

// LoadWorkspace resolves the workspace name from the default workspace in settings.
// If no default workspace is set, it returns an error prompting the user to use "drycc workspaces switch".
func LoadWorkspace(cf string) (string, *settings.Settings, error) {
	s, err := settings.Load(cf)
	if err != nil {
		return "", nil, err
	}

	if s.Workspace == "" {
		return "", nil, fmt.Errorf("no workspace specified, set a default workspace with 'drycc workspaces switch'")
	}

	return s.Workspace, s, nil
}
