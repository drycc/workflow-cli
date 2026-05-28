package parser

import (
	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/internal/template"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

// NewRoutesCommand creates a command for managing application routes.
func NewRoutesCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "routes",
		Short: i18n.T("Manage routes for your applications"),
		RunE: func(_ *cobra.Command, _ []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.RoutesList(app, results)
		},
	}

	cmd.PersistentFlags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))
	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)

	cmd.AddCommand(routesListCommand(cmdr))
	cmd.AddCommand(routesInfoCommand(cmdr))
	cmd.AddCommand(routesApplyCommand(cmdr))
	cmd.AddCommand(routesRemoveCommand(cmdr))
	return cmd
}

func routesInfoCommand(cmdr *commands.DryccCmd) *cobra.Command {
	routeCompletion := completion.RouteCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "info <name>",
		Args: cobra.ExactArgs(1),
		Example: template.CustomExample(
			"drycc routes info myroute",
			map[string]string{
				"<name>": i18n.T("The unique name of the route"),
			},
		),
		Short:             i18n.T("Show route information"),
		Long:              i18n.T("Shows detailed information about a route"),
		ValidArgsFunction: routeCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			name := args[0]
			return cmdr.RoutesInfo(app, name)
		},
	}

	return cmd
}

func routesApplyCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "apply <file>",
		Args: cobra.ExactArgs(1),
		Example: template.CustomExample(
			"drycc routes apply route.yaml",
			map[string]string{
				"<file>": i18n.T("Path to the YAML configuration file"),
			},
		),
		Short: i18n.T("Apply route configuration"),
		Long:  i18n.T("Applies route configuration from a YAML file using the raw API payload in the file"),
		RunE: func(_ *cobra.Command, args []string) error {
			return cmdr.RoutesApply(app, args[0])
		},
	}

	return cmd
}

func routesListCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: i18n.T("List application routes"),
		RunE: func(_ *cobra.Command, _ []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.RoutesList(app, results)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))

	return cmd
}

func routesRemoveCommand(cmdr *commands.DryccCmd) *cobra.Command {
	routeCompletion := completion.RouteCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "remove <name>",
		Args: cobra.ExactArgs(1),
		Example: template.CustomExample(
			"drycc routes remove myroute",
			map[string]string{
				"<name>": i18n.T("The unique name of the route"),
			},
		),
		Short:             i18n.T("Remove a route from an application"),
		ValidArgsFunction: routeCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			name := args[0]
			return cmdr.RoutesRemove(app, name)
		},
	}

	return cmd
}
