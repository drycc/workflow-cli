package settings

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var filepathRegex = regexp.MustCompile(`^.*[/\\].+\.json$`)

func locateSettingsFile(cf string) string {
	if strings.HasPrefix(cf, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		cf = homeDir + cf[1:]
	}

	cf = os.ExpandEnv(cf)

	if cf == "" {
		if v, ok := os.LookupEnv("DRYCC_PROFILE"); ok {
			cf = v
		} else {
			cf = "client"
		}
	}

	cf = filepath.Clean(cf)

	if filepathRegex.MatchString(cf) {
		return cf
	}

	return filepath.Join(DryccHome(), cf+".json")
}
