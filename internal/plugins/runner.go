// Package plugins provides a plugin system for extending the Drycc CLI.
package plugins

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/drycc/workflow-cli/pkg/settings"
)

// Plugin represents a CLI plugin
type Plugin struct {
	Name string
	Path string
}

// LookupPlugin searches for a plugin executable in PATH
// Plugin naming convention: drycc-<name>
func LookupPlugin(name string) (string, bool) {
	pluginName := fmt.Sprintf("drycc-%s", name)
	path, err := exec.LookPath(pluginName)
	if err != nil {
		return "", false
	}
	return path, true
}

// ListPlugins returns all available plugins in PATH
func ListPlugins() []Plugin {
	var plugins []Plugin
	seen := make(map[string]bool)

	pathEnv := os.Getenv("PATH")
	paths := filepath.SplitList(pathEnv)

	for _, dir := range paths {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			name := entry.Name()
			if !strings.HasPrefix(name, "drycc-") {
				continue
			}

			// Extract plugin name (remove "drycc-" prefix)
			pluginName := strings.TrimPrefix(name, "drycc-")

			// On Windows, remove .exe suffix
			if runtime.GOOS == "windows" {
				pluginName = strings.TrimSuffix(pluginName, ".exe")
			}

			// Skip if we've already seen this plugin
			if seen[pluginName] {
				continue
			}

			fullPath := filepath.Join(dir, name)

			// Check if executable
			info, err := entry.Info()
			if err != nil {
				continue
			}

			if runtime.GOOS != "windows" {
				if info.Mode()&0111 == 0 {
					continue
				}
			}

			seen[pluginName] = true
			plugins = append(plugins, Plugin{
				Name: pluginName,
				Path: fullPath,
			})
		}
	}

	return plugins
}

// Run executes a plugin with the given arguments
func Run(pluginPath string, args []string, s *settings.Settings) error {
	cmd := exec.Command(pluginPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Pass environment variables
	cmd.Env = os.Environ()
	if s.Client != nil && s.Client.ControllerURL != nil {
		cmd.Env = append(cmd.Env, fmt.Sprintf("DRYCC_CONTROLLER_URL=%s", s.Client.ControllerURL.String()))
		cmd.Env = append(cmd.Env, fmt.Sprintf("DRYCC_TOKEN=%s", s.Client.Token))
		cmd.Env = append(cmd.Env, fmt.Sprintf("DRYCC_SSL_VERIFY=%t", s.Client.VerifySSL))
	}
	cmd.Env = append(cmd.Env, fmt.Sprintf("DRYCC_USERNAME=%s", s.Username))
	cmd.Env = append(cmd.Env, fmt.Sprintf("DRYCC_RESPONSE_LIMIT=%d", s.Limit))

	return cmd.Run()
}
