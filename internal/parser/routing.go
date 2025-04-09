package parser

import (
	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

func NewRoutingCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "routing",
		Short: i18n.T("Manage routability of an application"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.RoutingInfo(app)
		},
	}

	cmd.PersistentFlags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)

	cmd.AddCommand(routingInfoCommand(cmdr))
	cmd.AddCommand(routingEnableCommand(cmdr))
	cmd.AddCommand(routingDisableCommand(cmdr))
	return cmd
}

func routingInfoCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: i18n.T("View application's routability information"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.RoutingInfo(app)
		},
	}

	return cmd
}

func routingEnableCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable",
		Short: i18n.T("Enable routing for an application"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.RoutingEnable(app)
		},
	}

	return cmd
}

func routingDisableCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable",
		Short: i18n.T("Disable routing for an application"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.RoutingDisable(app)
		},
	}

	return cmd
}
