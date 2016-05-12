package client

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/deis/workflow-cli/version"
)

func locateSettingsFile() string {
	filename := os.Getenv("DEIS_PROFILE")

	if filename == "" {
		filename = "client"
	}

	return filepath.Join(FindHome(), ".deis", filename+".json")
}

func checkAPICompatibility(serverAPIVersion string) {
	if serverAPIVersion != version.APIVersion {
		fmt.Printf(`!    WARNING: Client and server API versions do not match. Please consider upgrading.
!    Client version: %s
!    Server version: %s
`, version.APIVersion, serverAPIVersion)
	}
}
