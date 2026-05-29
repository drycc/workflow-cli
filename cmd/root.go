// Package cmd provides the root command for the Drycc CLI.
/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"time"

	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/parser"
	"github.com/drycc/workflow-cli/internal/plugins"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/drycc/workflow-cli/pkg/settings"
	"github.com/spf13/cobra"
)

// NewDryccCommand creates the root command for the Drycc CLI.
func NewDryccCommand() *cobra.Command {
	var flags struct {
		config  string
		version bool
		help    bool
	}

	var cmdr commands.DryccCmd

	rootCmd := &cobra.Command{
		Use:   "drycc",
		Short: i18n.T("The Drycc command-line client issues API calls to a Drycc controller"),
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			cmdr = commands.DryccCmd{ConfigFile: flags.config, WOut: os.Stdout, WErr: os.Stderr, WIn: os.Stdin, Location: time.Local}
		},
	}
	config := "~/.drycc/client.json"
	if v, ok := os.LookupEnv("DRYCC_PROFILE"); ok {
		config = v
	}
	rootCmd.PersistentFlags().StringVarP(&flags.config, "config", "c", config, i18n.T("Path to configuration file"))
	rootCmd.PersistentFlags().BoolVarP(&flags.help, "help", "h", false, i18n.T("Display help information"))

	rootCmd.AddCommand(parser.NewAppsCommand(&cmdr))
	rootCmd.AddCommand(parser.NewAuthCommand(&cmdr))
	rootCmd.AddCommand(parser.NewAutodeployCommand(&cmdr))
	rootCmd.AddCommand(parser.NewAutorollbackCommand(&cmdr))
	rootCmd.AddCommand(parser.NewAutoscaleCommand(&cmdr))
	rootCmd.AddCommand(parser.NewBuildsCommand(&cmdr))
	rootCmd.AddCommand(parser.NewCertsCommand(&cmdr))
	rootCmd.AddCommand(parser.NewConfigCommand(&cmdr))
	rootCmd.AddCommand(parser.NewDomainsCommand(&cmdr))
	rootCmd.AddCommand(parser.NewGatewaysCommand(&cmdr))
	rootCmd.AddCommand(parser.NewGitCommand(&cmdr))
	rootCmd.AddCommand(parser.NewLifecyclesCommand(&cmdr))
	rootCmd.AddCommand(parser.NewHealthchecksCommand(&cmdr))
	rootCmd.AddCommand(parser.NewKeysCommand(&cmdr))
	rootCmd.AddCommand(parser.NewLabelsCommand(&cmdr))
	rootCmd.AddCommand(parser.NewLimitsCommand(&cmdr))
	rootCmd.AddCommand(parser.NewWorkspacesCommand(&cmdr))
	rootCmd.AddCommand(parser.NewPsCommand(&cmdr))
	rootCmd.AddCommand(parser.NewPtsCommand(&cmdr))
	rootCmd.AddCommand(parser.NewPluginsCommand(&cmdr))
	rootCmd.AddCommand(parser.NewRegistryCommand(&cmdr))
	rootCmd.AddCommand(parser.NewReleasesCommand(&cmdr))
	rootCmd.AddCommand(parser.NewRoutesCommand(&cmdr))
	rootCmd.AddCommand(parser.NewRoutingCommand(&cmdr))
	rootCmd.AddCommand(parser.NewServicesCommand(&cmdr))
	rootCmd.AddCommand(parser.NewTagsCommand(&cmdr))
	rootCmd.AddCommand(parser.NewTimeoutsCommand(&cmdr))
	rootCmd.AddCommand(parser.NewTLSCommand(&cmdr))
	rootCmd.AddCommand(parser.NewTokensCommand(&cmdr))
	rootCmd.AddCommand(parser.NewUpdateCommand(&cmdr))
	rootCmd.AddCommand(parser.NewVolumesCommand(&cmdr))
	rootCmd.AddCommand(parser.NewVersionCommand(&cmdr))
	// shortcuts
	rootCmd.AddGroup(&cobra.Group{ID: "shortcut", Title: i18n.T("Shortcut Commands:")})
	for _, shortcuts := range parser.SupportedShortcuts {
		for _, shortcut := range shortcuts.Create(&cmdr) {
			shortcut.GroupID = "shortcut"
			rootCmd.AddCommand(shortcut)
		}
	}
	rootCmd.SilenceUsage = true

	return rootCmd
}

// ExecuteWithPlugins runs the root command with plugin dispatch support
func ExecuteWithPlugins(rootCmd *cobra.Command, config string) error {
	// Try to find the command first
	cmd, _, err := rootCmd.Find(os.Args[1:])
	if err != nil || cmd == rootCmd {
		// Command not found or is root command, try plugin dispatch
		name := firstNonFlagArg(os.Args[1:])
		if name != "" {
			if path, ok := plugins.LookupPlugin(name); ok {
				restArgs := argsAfter(os.Args[1:], name)
				s, serr := settings.Load(config)
				if serr != nil {
					// If settings can't be loaded, use empty settings
					s = &settings.Settings{}
				}
				return plugins.Run(path, restArgs, s)
			}
		}
	}
	return rootCmd.Execute()
}

// firstNonFlagArg returns the first argument that doesn't start with "-"
func firstNonFlagArg(args []string) string {
	for _, arg := range args {
		if len(arg) > 0 && arg[0] != '-' {
			return arg
		}
	}
	return ""
}

// argsAfter returns all arguments after the specified name
func argsAfter(args []string, name string) []string {
	for i, arg := range args {
		if arg == name {
			return args[i+1:]
		}
	}
	return args
}
