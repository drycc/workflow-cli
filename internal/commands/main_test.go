package commands

import (
	"os"
	"testing"
)

// TestMain sets the locale to en_US so that i18n.T() returns English output
// in all command tests, making assertions deterministic regardless of the
// system's LANG environment variable.
func TestMain(m *testing.M) {
	os.Setenv("LANG", "en_US.UTF-8")
	os.Setenv("LC_ALL", "en_US.UTF-8")
	os.Setenv("LC_MESSAGES", "en_US.UTF-8")
	code := m.Run()
	os.Exit(code)
}
