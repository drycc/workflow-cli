package parser

import (
	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/internal/template"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

func NewPsCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ps",
		Short: i18n.T("Manage processes inside an app"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.PsList(app, 1000)
		},
	}

	cmd.PersistentFlags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)

	cmd.AddCommand(psListCommand(cmdr))
	cmd.AddCommand(psLogsCommand(cmdr))
	cmd.AddCommand(psExecCommand(cmdr))
	cmd.AddCommand(psDescribeCommand(cmdr))
	cmd.AddCommand(psDeleteCommand(cmdr))
	return cmd
}

func psListCommand(cmdr *commands.DryccCmd) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "list",
		Short: i18n.T("List application pods"),
		Long:  i18n.T("Lists processes servicing an application"),
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.PsList(app, 1000)
		},
	}

	return cmd
}

func psLogsCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		app       string
		lines     int
		follow    bool
		container string
		previous  bool
	}

	psCompletion := completion.PsCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}

	cmd := &cobra.Command{
		Use:               "logs <pod>",
		Args:              cobra.ExactArgs(1),
		Example:           "drycc ps logs my-pod",
		Short:             i18n.T("Print the logs for a container"),
		Long:              i18n.T("Print the logs for a container in a pod or specified resource"),
		ValidArgsFunction: psCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			podID := args[0]
			if flags.lines < 0 {
				flags.lines = -1
			}
			return cmdr.PsLogs(app, podID, flags.lines, flags.follow, flags.container, flags.previous)
		},
	}

	cmd.Flags().IntVarP(&flags.lines, "lines", "l", 300, i18n.T("The number of lines to display, -1 showing all log lines"))
	cmd.Flags().BoolVarP(&flags.follow, "follow", "f", false, i18n.T("Specify if the logs should be streamed"))
	cmd.Flags().StringVar(&flags.container, "container", "", i18n.T("Print the logs of this container"))
	cmd.Flags().BoolVarP(&flags.previous, "previous", "p", false, i18n.T("Print the logs for the previous instance of the container in a pod if it exists"))
	cmd.Flags().SortFlags = false

	return cmd
}

func psExecCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		app     string
		pod     string
		command []string
		tty     bool
		stdin   bool
	}
	psCompletion := completion.PsCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "exec <pod> [flags] -- <command>...",
		Args: cobra.MinimumNArgs(1),
		Example: template.CustomExample(
			"drycc ps exec my-pod -it -- bash",
			map[string]string{
				"<pod>": i18n.T("The pod name for the application"),
			},
		),
		Short:             i18n.T("Execute a command in a container"),
		ValidArgsFunction: psCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			flags.pod = args[0]
			flags.command = args[1:]
			return cmdr.PsExec(app, flags.pod, flags.tty, flags.stdin, flags.command)
		},
	}

	// shortcuts exec has not app flag
	cmd.Flags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))
	cmd.Flags().BoolVarP(&flags.tty, "tty", "t", false, i18n.T("Stdin is a TTY"))
	cmd.Flags().BoolVarP(&flags.stdin, "stdin", "i", false, i18n.T("Pass stdin to the container"))

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)

	return cmd
}

func psDescribeCommand(cmdr *commands.DryccCmd) *cobra.Command {
	psCompletion := completion.PsCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "describe <pod>",
		Args: cobra.ExactArgs(1),
		Example: template.CustomExample(
			"drycc ps describe my-pod",
			map[string]string{
				"<pod>": i18n.T("The pod name for the application"),
			},
		),
		Short:             i18n.T("Print a detailed description of the selected process"),
		ValidArgsFunction: psCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			podID := args[0]
			return cmdr.PsDescribe(app, podID)
		},
	}

	return cmd
}

func psDeleteCommand(cmdr *commands.DryccCmd) *cobra.Command {
	psCompletion := completion.PsCompletion{AppID: &app, ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "delete <pod>...",
		Args: cobra.MinimumNArgs(1),
		Example: template.CustomExample(
			"drycc ps delete my-pod another-pod",
			map[string]string{
				"<pod>": i18n.T("The pod name for the application"),
			},
		),
		Short:             i18n.T("Delete the selected processes"),
		ValidArgsFunction: psCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			return cmdr.PsDelete(app, args)
		},
	}

	return cmd
}
