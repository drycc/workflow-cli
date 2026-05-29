package parser

import (
	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/plugins"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

// NewPluginsCommand creates a command for managing plugins
func NewPluginsCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugins",
		Short: i18n.T("Provides utilities for interacting with plugins"),
		Long: i18n.T(`Provides utilities for interacting with plugins.

Plugins provide extended functionality that is not part of the major command-line distribution.

The easiest way to use plugins is to place executables named 'drycc-<name>' in your PATH.
When you run 'drycc <name>', the CLI will automatically invoke the 'drycc-<name>' plugin if it exists.`),
		Example: i18n.T(`  # List all available plugins
  drycc plugins list`),
	}

	cmd.AddCommand(pluginsListCommand(cmdr))

	return cmd
}

func pluginsListCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: i18n.T("List all visible plugin executables on a user's PATH"),
		Long: i18n.T(`List all available plugin files on a user's PATH.

Available plugin files are those that are:
- executable
- anywhere on the user's PATH
- begin with "drycc-"`),
		Example: i18n.T(`  # List all available plugins
  drycc plugins list`),
		RunE: func(_ *cobra.Command, _ []string) error {
			plugins := plugins.ListPlugins()

			if len(plugins) == 0 {
				cmdr.Println("Unable to find any drycc plugins in your PATH")
				return nil
			}

			cmdr.Println("The following compatible plugins are available:")
			cmdr.Println()

			for _, plugin := range plugins {
				cmdr.Println(plugin.Path)
			}

			return nil
		},
	}

	return cmd
}
