//go:build linux || darwin
// +build linux darwin

// Package settings provides functionality for managing Drycc CLI settings.
package settings

import (
	"os"
)

// FindHome returns the HOME directory of the current user
func FindHome() string {
	return os.Getenv("HOME")
}

// SetHome sets the HOME directory of the current user
func SetHome(path string) {
	os.Setenv("HOME", path)
}
