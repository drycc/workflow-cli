/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"time"

	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/parser"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

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
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cmdr = commands.DryccCmd{ConfigFile: flags.config, WOut: os.Stdout, WErr: os.Stderr, WIn: os.Stdin, Location: time.Local}
		},
	}

	rootCmd.PersistentFlags().StringVarP(&flags.config, "config", "c", "~/.drycc/client.json", i18n.T("Path to configuration file"))
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
	rootCmd.AddCommand(parser.NewHealthchecksCommand(&cmdr))
	rootCmd.AddCommand(parser.NewKeysCommand(&cmdr))
	rootCmd.AddCommand(parser.NewLabelsCommand(&cmdr))
	rootCmd.AddCommand(parser.NewLimitsCommand(&cmdr))
	rootCmd.AddCommand(parser.NewPermsCommand(&cmdr))
	rootCmd.AddCommand(parser.NewPsCommand(&cmdr))
	rootCmd.AddCommand(parser.NewPtsCommand(&cmdr))
	rootCmd.AddCommand(parser.NewRegistryCommand(&cmdr))
	rootCmd.AddCommand(parser.NewReleasesCommand(&cmdr))
	rootCmd.AddCommand(parser.NewResourcesCommand(&cmdr))
	rootCmd.AddCommand(parser.NewRoutesCommand(&cmdr))
	rootCmd.AddCommand(parser.NewRoutingCommand(&cmdr))
	rootCmd.AddCommand(parser.NewServicesCommand(&cmdr))
	rootCmd.AddCommand(parser.NewTagsCommand(&cmdr))
	rootCmd.AddCommand(parser.NewTimeoutsCommand(&cmdr))
	rootCmd.AddCommand(parser.NewTLSCommand(&cmdr))
	rootCmd.AddCommand(parser.NewTokensCommand(&cmdr))
	rootCmd.AddCommand(parser.NewUpdateCommand(&cmdr))
	rootCmd.AddCommand(parser.NewUsersCommand(&cmdr))
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
