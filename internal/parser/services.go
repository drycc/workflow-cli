package parser

import (
	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

// NewServicesCommand creates the services command
func NewServicesCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "services",
		Short: i18n.T("Manage services for your applications"),
		RunE: func(_ *cobra.Command, _ []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.ServicesList(app, results)
		},
	}

	cmd.PersistentFlags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))
	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)

	cmd.AddCommand(servicesList(cmdr))
	cmd.AddCommand(servicesInfo(cmdr))
	cmd.AddCommand(servicesApply(cmdr))
	cmd.AddCommand(servicesRemove(cmdr))
	return cmd
}

func servicesList(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: i18n.T("List application services"),
		Long:  i18n.T("Lists services for an application"),
		RunE: func(_ *cobra.Command, _ []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.ServicesList(app, results)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))

	return cmd
}

func servicesInfo(cmdr *commands.DryccCmd) *cobra.Command {
	serviceCompletion := completion.ServiceCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:               "info <ptype>",
		Args:              cobra.ExactArgs(1),
		Short:             i18n.T("Show service information"),
		Long:              i18n.T("Shows detailed information about a service"),
		ValidArgsFunction: serviceCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			ptype := args[0]
			return cmdr.ServicesInfo(app, ptype)
		},
	}

	return cmd
}

func servicesApply(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply <file>",
		Args:  cobra.ExactArgs(1),
		Short: i18n.T("Apply service configuration"),
		Long:  i18n.T("Applies service configuration from a TOML file using the raw API payload in the file"),
		RunE: func(_ *cobra.Command, args []string) error {
			return cmdr.ServicesApply(app, args[0])
		},
	}

	return cmd
}

func servicesRemove(cmdr *commands.DryccCmd) *cobra.Command {
	serviceCompletion := completion.ServiceCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:               "remove <ptype>",
		Args:              cobra.ExactArgs(1),
		Short:             i18n.T("Remove a service"),
		Long:              i18n.T("Removes a service from an application"),
		ValidArgsFunction: serviceCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			ptype := args[0]
			return cmdr.ServicesRemove(app, ptype)
		},
	}

	return cmd
}
