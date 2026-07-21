package parser

import (
	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/internal/template"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

// NewAddonsCommand creates a command for managing addon instances attached to applications.
// The `classes` subcommand queries the global addon catalog and is not app-scoped.
func NewAddonsCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addons",
		Short: i18n.T("Manage addons for your applications"),
		RunE: func(_ *cobra.Command, _ []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.AddonsList(app, results)
		},
	}

	cmd.PersistentFlags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))
	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)

	cmd.AddCommand(addonsListCommand(cmdr))
	cmd.AddCommand(addonsInfoCommand(cmdr))
	cmd.AddCommand(addonsConnectionCommand(cmdr))
	cmd.AddCommand(addonsApplyCommand(cmdr))
	cmd.AddCommand(addonsRemoveCommand(cmdr))
	cmd.AddCommand(addonsClassesCommand(cmdr))
	return cmd
}

func addonsListCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: i18n.T("List addons attached to an application"),
		RunE: func(_ *cobra.Command, _ []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.AddonsList(app, results)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))

	return cmd
}

func addonsInfoCommand(cmdr *commands.DryccCmd) *cobra.Command {
	addonCompletion := completion.AddonCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "info <name>",
		Args: cobra.ExactArgs(1),
		Example: template.CustomExample(
			"drycc addons info my-addon",
			map[string]string{
				"<name>": i18n.T("The addon name"),
			},
		),
		Short:             i18n.T("Show information about an addon"),
		Long:              i18n.T("Shows detailed information about an addon attached to an application"),
		ValidArgsFunction: addonCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			return cmdr.AddonsInfo(app, args[0])
		},
	}

	return cmd
}

func addonsConnectionCommand(cmdr *commands.DryccCmd) *cobra.Command {
	addonCompletion := completion.AddonCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "conn <name>",
		Args: cobra.ExactArgs(1),
		Example: template.CustomExample(
			"drycc addons conn my-addon",
			map[string]string{
				"<name>": i18n.T("The addon name"),
			},
		),
		Short:             i18n.T("Show connection details for an addon"),
		Long:              i18n.T("Shows the connection details (such as host, port and credentials) for an addon attached to an application"),
		ValidArgsFunction: addonCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			return cmdr.AddonsConnection(app, args[0])
		},
	}

	return cmd
}

func addonsApplyCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "apply <file>",
		Args: cobra.ExactArgs(1),
		Example: template.CustomExample(
			"drycc addons apply addon.yaml",
			map[string]string{
				"<file>": i18n.T("Path to the YAML configuration file"),
			},
		),
		Short: i18n.T("Apply addon configuration"),
		Long:  i18n.T("Applies addon configuration from a YAML file using the raw API payload in the file"),
		RunE: func(_ *cobra.Command, args []string) error {
			return cmdr.AddonsApply(app, args[0])
		},
	}

	return cmd
}

func addonsRemoveCommand(cmdr *commands.DryccCmd) *cobra.Command {
	addonCompletion := completion.AddonCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "remove <name>",
		Args: cobra.ExactArgs(1),
		Example: template.CustomExample(
			"drycc addons remove my-addon",
			map[string]string{
				"<name>": i18n.T("The addon name to be removed"),
			},
		),
		Short:             i18n.T("Remove an addon from an application"),
		Long:              i18n.T("Removes an addon from an application"),
		ValidArgsFunction: addonCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			return cmdr.AddonsRemove(app, args[0])
		},
	}

	return cmd
}

func addonsClassesCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "classes [name]",
		Short: i18n.T("List available addon classes, or show details for one class"),
		Long:  i18n.T("Without a name argument, lists all addon classes in the catalog. With a name, shows detailed information about that class including its plans."),
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) == 1 {
				return cmdr.AddonsClassesInfo(args[0])
			}
			results, _ := commands.ResponseLimit(limit)
			return cmdr.AddonsClassesList(results)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))

	return cmd
}
