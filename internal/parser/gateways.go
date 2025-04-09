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
	cmd.AddCommand(gatewaysAdd(cmdr))
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

func gatewaysAdd(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		name     string
		port     int
		protocol string
	}

	cmd := &cobra.Command{
		Use:   "add <name>",
		Args:  cobra.ExactArgs(1),
		Short: i18n.T("Create gateways for an application"),
		Long:  i18n.T("Creates gateways for an application and binds it to allow listener of the main app domain"),
		RunE: func(_ *cobra.Command, args []string) error {
			name := args[0]
			return cmdr.GatewaysAdd(app, name, flags.port, flags.protocol)
		},
	}

	cmd.Flags().IntVar(&flags.port, "port", 0, i18n.T("Port is the network port, the listener expects to receive"))
	cmd.Flags().StringVarP(&flags.protocol, "protocol", "", "", i18n.T("Protocol specifies the network protocol this listener expects to receive. Supports TCP, UDP, TLS, HTTP, and HTTPS"))

	mustFlags := []string{"port", "protocol"}
	for _, mustFlag := range mustFlags {
		cmd.MarkFlagRequired(mustFlag)
	}

	gatewayProtocolCompletion := completion.GatewayProtocolCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("protocol", gatewayProtocolCompletion.CompletionFunc)
	return cmd
}

func gatewaysRemove(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		name     string
		port     int
		protocol string
	}

	gatewayNameCompletion := completion.GatewayNameCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:               "remove <name> --port=<port> --protocol=<protocol>",
		Args:              cobra.ExactArgs(1),
		Short:             i18n.T("Remove gateways from an application"),
		Long:              i18n.T("Deletes specific gateway for application"),
		ValidArgsFunction: gatewayNameCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			name := args[0]
			return cmdr.GatewaysRemove(app, name, flags.port, flags.protocol)
		},
	}

	cmd.Flags().IntVar(&flags.port, "port", 0, i18n.T("Port is the network port, the listener received"))
	cmd.Flags().StringVarP(&flags.protocol, "protocol", "", "", i18n.T("Protocol specifies the network protocol this listener received"))

	mustFlags := []string{"port", "protocol"}
	for _, mustFlag := range mustFlags {
		cmd.MarkFlagRequired(mustFlag)
	}

	gatewayProtocolCompletion := completion.GatewayProtocolCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("protocol", gatewayProtocolCompletion.CompletionFunc)
	return cmd
}
