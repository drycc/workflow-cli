package parser

import (
	"fmt"

	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/internal/template"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

func NewResourcesCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resources",
		Short: i18n.T("Manage resources for your applications"),
		RunE: func(cmd *cobra.Command, args []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.ResourcesList(app, results)
		},
	}

	cmd.PersistentFlags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))
	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)

	cmd.AddCommand(resourcesServicesCommand(cmdr))
	cmd.AddCommand(resourcesPlansCommand(cmdr))
	cmd.AddCommand(resourcesCreateCommand(cmdr))
	cmd.AddCommand(resourcesListCommand(cmdr))
	cmd.AddCommand(resourcesDescribeCommand(cmdr))
	cmd.AddCommand(resourcesUpdateCommand(cmdr))
	cmd.AddCommand(resourcesBindCommand(cmdr))
	cmd.AddCommand(resourcesUnbindCommand(cmdr))
	cmd.AddCommand(resourcesDestroyCommand(cmdr))
	return cmd
}

func resourcesServicesCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "services",
		Short: i18n.T("List all available resource services"),
		RunE: func(cmd *cobra.Command, args []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.ResourcesServices(results)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))

	return cmd
}

func resourcesPlansCommand(cmdr *commands.DryccCmd) *cobra.Command {
	resourceServiceCompletion := completion.ResourceServiceCompletion{ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "plans <service>",
		Args: cobra.ExactArgs(1),
		Example: template.CustomExample(
			"drycc resources plans redis",
			map[string]string{
				"<service>": i18n.T("The service name for plans"),
			},
		),
		Short:             i18n.T("List all available plans for a resource service"),
		ValidArgsFunction: resourceServiceCompletion.CompletionFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			service := args[0]
			results, _ := commands.ResponseLimit(limit)
			return cmdr.ResourcesPlans(service, results)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))

	return cmd
}

func resourcesCreateCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		values  string
		path    string
		confirm string
	}
	resourceCreateCompletion := completion.ResourceCreateCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "create <name> <service> <plan> [<param>=<value>...]",
		Args: cobra.MinimumNArgs(3),
		Example: template.CustomExample(
			"drycc resources create myredis redis standard-128 -f file.yaml",
			map[string]string{
				"<name>":    i18n.T("This resource instance alias"),
				"<service>": i18n.T("The resource's service"),
				"<plan>":    i18n.T("The service's plan"),
				"<param>":   i18n.T("The resource instance parameters key"),
				"<value>":   i18n.T("The resource instance parameters value"),
			},
		),
		Short:             i18n.T("Create a resource for the application"),
		ValidArgsFunction: resourceCreateCompletion.CompletionFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			plan := fmt.Sprintf("%s:%s", args[1], args[2])
			params := args[3:]

			if flags.values != "" {
				params = nil
			}
			return cmdr.ResourcesCreate(app, plan, name, params, flags.values)
		},
	}

	cmd.Flags().StringVarP(&flags.values, "values", "f", "", i18n.T("Specify values in a YAML file. If set, params will be discard"))

	return cmd
}

func resourcesListCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: i18n.T("List resources in the application"),
		RunE: func(cmd *cobra.Command, args []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.ResourcesList(app, results)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))

	return cmd
}

func resourcesDescribeCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		app     string
		name    string
		details bool
	}
	resourceCompletion := completion.ResourceCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "describe <name>",
		Args: cobra.ExactArgs(1),
		Example: template.CustomExample(
			"drycc resources describe myredis",
			map[string]string{
				"<name>": i18n.T("This resource instance alias"),
			},
		),
		Short:             i18n.T("Get a resource's detail in the application"),
		ValidArgsFunction: resourceCompletion.CompletionFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			return cmdr.ResourceGet(app, name, flags.details)
		},
	}

	cmd.Flags().BoolVar(&flags.details, "details", false, i18n.T("Show extra details provided to resource"))

	return cmd
}

func resourcesUpdateCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		app    string
		values string
		path   string
	}
	resourceUpdateCompletion := completion.ResourceUpdateCompletion{AppID: &app, ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "update <name> <service> <plan> [<param>=<value>...]",
		Args: cobra.MinimumNArgs(3),
		Example: template.CustomExample(
			"drycc resources update myredis redis standard-128 -f file.yaml",
			map[string]string{
				"<name>":    i18n.T("This resource instance alias"),
				"<service>": i18n.T("The resource's service"),
				"<plan>":    i18n.T("The service's plan"),
				"<param>":   i18n.T("The resource instance parameters key"),
				"<value>":   i18n.T("The resource instance parameters value"),
			},
		),
		Short:             i18n.T("Update a resource from the application"),
		ValidArgsFunction: resourceUpdateCompletion.CompletionFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			plan := fmt.Sprintf("%s:%s", args[1], args[2])
			params := args[2:]

			if flags.values != "" {
				params = nil
			}

			return cmdr.ResourcePut(app, plan, name, params, flags.values)
		},
	}

	cmd.Flags().StringVarP(&flags.values, "values", "f", "", i18n.T("Specify values in a YAML file. If set, params will be discard"))

	return cmd
}

func resourcesDestroyCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		app     string
		name    string
		confirm string
	}
	resourceCompletion := completion.ResourceCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "destroy <name> --confirm=<resource>",
		Args: cobra.ExactArgs(1),
		Example: template.CustomExample(
			"drycc resources destroy myredis",
			map[string]string{
				"<name>": i18n.T("The resource instance alias name"),
			},
		),
		Short:             i18n.T("Delete a resource from the application"),
		ValidArgsFunction: resourceCompletion.CompletionFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			return cmdr.ResourceDelete(app, name, flags.confirm)
		},
	}

	cmd.Flags().StringVar(&flags.confirm, "confirm", "", i18n.T(`skips the prompt for the resource name. <resource> is the uniquely identifiable
name for the resource`))

	return cmd
}

func resourcesBindCommand(cmdr *commands.DryccCmd) *cobra.Command {
	resourceCompletion := completion.ResourceCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "bind <name>",
		Args: cobra.ExactArgs(1),
		Example: template.CustomExample(
			"drycc resources bind myredis",
			map[string]string{
				"<name>": i18n.T("The resource instance alias name"),
			},
		),
		Short:             i18n.T("Bind a resource for an application"),
		ValidArgsFunction: resourceCompletion.CompletionFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			return cmdr.ResourceBind(app, name)
		},
	}

	return cmd
}

func resourcesUnbindCommand(cmdr *commands.DryccCmd) *cobra.Command {
	resourceCompletion := completion.ResourceCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "unbind <name>",
		Args: cobra.ExactArgs(1),
		Example: template.CustomExample(
			"drycc resources unbind myredis",
			map[string]string{
				"<name>": i18n.T("The resource instance alias name"),
			},
		),
		Short:             i18n.T("unbind a resources for an application"),
		ValidArgsFunction: resourceCompletion.CompletionFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			return cmdr.ResourceUnbind(app, name)
		},
	}

	return cmd
}
