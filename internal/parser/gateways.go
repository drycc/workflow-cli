package parser

import (
	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

// NewGatewaysCommand creates the gateways command
func NewGatewaysCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gateways",
		Short: i18n.T("Manage gateways for your applications"),
		RunE: func(_ *cobra.Command, _ []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.GatewaysList(app, results)
		},
	}

	cmd.PersistentFlags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))
	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)

	cmd.AddCommand(gatewaysList(cmdr))
	cmd.AddCommand(gatewaysInfo(cmdr))
	cmd.AddCommand(gatewaysApply(cmdr))
	cmd.AddCommand(gatewaysRemove(cmdr))
	return cmd
}

func gatewaysList(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: i18n.T("List application gateways"),
		Long:  i18n.T("Lists gateways for an application"),
		RunE: func(_ *cobra.Command, _ []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.GatewaysList(app, results)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))

	return cmd
}

func gatewaysInfo(cmdr *commands.DryccCmd) *cobra.Command {
	gatewayNameCompletion := completion.GatewayNameCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:               "info <name>",
		Args:              cobra.ExactArgs(1),
		Short:             i18n.T("Show gateway information"),
		Long:              i18n.T("Shows detailed information about a gateway"),
		ValidArgsFunction: gatewayNameCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			name := args[0]
			return cmdr.GatewaysInfo(app, name)
		},
	}

	return cmd
}

func gatewaysApply(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply <file>",
		Args:  cobra.ExactArgs(1),
		Short: i18n.T("Apply gateway configuration"),
		Long:  i18n.T("Applies gateway configuration from a YAML file using the raw API payload in the file"),
		RunE: func(_ *cobra.Command, args []string) error {
			return cmdr.GatewaysApply(app, args[0])
		},
	}

	return cmd
}

func gatewaysRemove(cmdr *commands.DryccCmd) *cobra.Command {
	gatewayNameCompletion := completion.GatewayNameCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:               "remove <name>",
		Args:              cobra.ExactArgs(1),
		Short:             i18n.T("Remove a gateway"),
		Long:              i18n.T("Removes a gateway from an application"),
		ValidArgsFunction: gatewayNameCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			name := args[0]
			return cmdr.GatewaysRemove(app, name)
		},
	}

	return cmd
}
