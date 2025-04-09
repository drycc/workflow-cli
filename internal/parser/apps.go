package parser

import (
	"strings"

	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/internal/template"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

// NewAppsCommand creates the apps command
func NewAppsCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apps",
		Short: i18n.T("Manage applications used to provide services"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.AppsList(limit)
		},
	}
	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("Maximum number of results to display"))

	cmd.AddCommand(appsCreate(cmdr))
	cmd.AddCommand(appsList(cmdr))
	cmd.AddCommand(appsInfo(cmdr))
	cmd.AddCommand(appsOpen(cmdr))
	cmd.AddCommand(appsLogs(cmdr))
	cmd.AddCommand(appsRun(cmdr))
	cmd.AddCommand(appsDestroy(cmdr))
	cmd.AddCommand(appsTransfer(cmdr))
	return cmd
}

// AppsCreate creates the apps create command
func appsCreate(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		remote   string
		noRemote bool
	}
	cmd := &cobra.Command{
		Use:   "create [<id>]",
		Args:  cobra.MaximumNArgs(1),
		Short: i18n.T("Create a new application"),
		Long:  i18n.T(`Creates a new application. If no <id> is provided, one will be generated automatically`),
		RunE: func(_ *cobra.Command, args []string) error {
			id := ""
			if len(args) > 0 {
				id = args[0]
			}
			return cmdr.AppCreate(id, flags.remote, flags.noRemote)
		},
	}

	cmd.Flags().StringVarP(&flags.remote, "remote", "r", "drycc", i18n.T("Name of remote to create"))
	cmd.Flags().BoolVar(&flags.noRemote, "no-remote", false, i18n.T("Do not create a 'drycc' git remote"))
	return cmd
}

// AppsList creates the apps list command
func appsList(cmdr *commands.DryccCmd) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "list",
		Short: i18n.T("List accessible applications"),
		Long:  i18n.T("Lists applications visible to the current user"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.AppsList(limit)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("Maximum number of results to display"))
	return cmd
}

// AppsInfo creates the apps info command
func appsInfo(cmdr *commands.DryccCmd) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "info",
		Short: i18n.T("View info about an application"),
		Long:  i18n.T("Prints info about the current application"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.AppInfo(app)
		},
	}

	cmd.Flags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)
	return cmd
}

// AppsOpen creates the apps open command
func appsOpen(cmdr *commands.DryccCmd) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "open",
		Short: i18n.T("Open the application in a browser"),
		Long:  i18n.T("Opens a URL to the application in the default browser"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.AppOpen(app)
		},
	}

	cmd.Flags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)
	return cmd
}

// AppsLogs creates the apps logs command
func appsLogs(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		app     string
		lines   int
		follow  bool
		timeout int
	}

	cmd := &cobra.Command{
		Use:   "logs",
		Short: i18n.T("Retrieve application log events"),
		Long:  i18n.T("Retrieves the most recent log events"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.AppLogs(app, flags.lines, flags.follow, flags.timeout)
		},
	}

	cmd.Flags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))
	cmd.Flags().IntVarP(&flags.lines, "lines", "n", 300, i18n.T("The number of lines to display"))
	cmd.Flags().BoolVarP(&flags.follow, "follow", "f", false, i18n.T("Specify if the logs should be streamed"))
	cmd.Flags().IntVarP(&flags.timeout, "timeout", "t", 300, i18n.T("The max seconds of followz the log stream"))

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)
	return cmd
}

// AppsRun creates the apps run command
func appsRun(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		app     string
		timeout uint32
		expires uint32
		mounts  []string
	}

	cmd := &cobra.Command{
		Use:  "run [--mount=<volume>:<path>...] -- <command>...",
		Args: cobra.MinimumNArgs(1),
		Example: template.CustomExample(
			"drycc apps run --mount=myvolume:/data -- 'echo hello'",
			map[string]string{
				"<volume>":  i18n.T("The volume name"),
				"<path>":    i18n.T("The filesystem path"),
				"<command>": i18n.T("The shell command to run inside the container"),
			},
		),
		Short: i18n.T("Run a command in an ephemeral app container"),
		Long:  i18n.T("Runs a command inside an ephemeral app container"),
		RunE: func(_ *cobra.Command, args []string) error {
			command := strings.Join(args, " ")
			return cmdr.AppRun(app, command, flags.mounts, flags.timeout, flags.expires)
		},
	}

	cmd.Flags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))
	cmd.Flags().StringSliceVarP(&flags.mounts, "mount", "m", nil, i18n.T("Volume mounts in format 'volume:path'"))
	cmd.Flags().Uint32VarP(&flags.timeout, "timeout", "t", 3600, i18n.T("Command execution timeout in seconds"))
	cmd.Flags().Uint32VarP(&flags.expires, "expires", "e", 3600, i18n.T("Retention time of records in seconds"))

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)
	return cmd
}

// AppsDestroy creates the apps destroy command
func appsDestroy(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		app     string
		confirm string
	}

	cmd := &cobra.Command{
		Use:     "destroy",
		Example: "drycc apps destroy -a <app> --confirm <app>",
		Short:   i18n.T("Destroy an application"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.AppDestroy(app, flags.confirm)
		},
	}

	cmd.Flags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))
	cmd.Flags().StringVar(&flags.confirm, "confirm", "", i18n.T("Skips the prompt for the application name. \n<app> is the uniquely identifiable name for the application"))

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)
	return cmd
}

// AppsTransfer creates the apps transfer command
func appsTransfer(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer <username>",
		Args:  cobra.ExactArgs(1),
		Short: i18n.T("Transfer app ownership to another user"),
		Long:  i18n.T("Transfer application ownership to another user."),
		RunE: func(_ *cobra.Command, args []string) error {
			user := args[0]
			return cmdr.AppTransfer(app, user)
		},
	}

	cmd.Flags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)
	return cmd
}
