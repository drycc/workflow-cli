package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/internal/template"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

var lifecycleFlags struct {
	ptype string
}

// NewLifecycleCommand creates a command for managing application lifecycle.
func NewLifecyclesCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lifecycles",
		Short: i18n.T("Manage application lifecycles"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.LifecyclesList(app, lifecycleFlags.ptype, version)
		},
	}

	cmd.PersistentFlags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name of the application"))
	cmd.PersistentFlags().StringVarP(&lifecycleFlags.ptype, "ptype", "p", "", i18n.T("The ptype for which the lifecycle needs to be listed"))
	cmd.Flags().IntVarP(&version, "version", "v", 0, i18n.T("The version for which the lifecycle needs to be listed"))

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)
	ptypeCompletion := completion.PtsCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("ptype", ptypeCompletion.CompletionFunc)
	releaseCompletion := completion.ReleaseCompletion{ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("version", releaseCompletion.CompletionFunc)

	cmd.AddCommand(lifecyclesList(cmdr))
	cmd.AddCommand(lifecyclesSet(cmdr))
	cmd.AddCommand(lifecyclesUnset(cmdr))
	return cmd
}

func lifecyclesList(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: i18n.T("List lifecycles for an application"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.LifecyclesList(app, lifecycleFlags.ptype, version)
		},
	}

	cmd.Flags().IntVarP(&version, "version", "v", 0, i18n.T("The version for which the lifecycle needs to be listed"))
	releaseCompletion := completion.ReleaseCompletion{ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("version", releaseCompletion.CompletionFunc)
	return cmd
}

func lifecyclesSet(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		ptype      string
		path       string
		port       int
		headers    string
		stopSignal string
	}

	lifecycleCompletion := completion.LifecycleCompletion{ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "set <handler> <action> [flags] [--] <args>...",
		Args: cobra.MinimumNArgs(2),
		Example: template.CustomExample(
			"drycc lifecycles set postStart httpGet --path=/health -- 8000",
			map[string]string{
				"<handler>": i18n.T("The lifecycle handler, such as 'postStart' 'preStop'"),
				"<action>":  i18n.T("the lifecycle action type, such as 'httpGet', 'exec', 'sleep' or 'tcpSocket'"),
				"<args>": i18n.T(`The arguments required for the lifecycle action. 'exec' accepts a list of arguments;
                  'sleep' accepts duration in seconds; 'httpGet' and 'tcpSocket' accept a port number.`),
			},
		),
		Short: i18n.T("Set lifecycles for an application"),
		Long: i18n.T(`Sets lifecycle handlers for an application.

Lifecycle handlers allow you to run actions at specific points in a container's lifecycle.
Two types of lifecycle handlers are supported: 'postStart' and 'preStop'.

A 'postStart' handler runs immediately after a container is created. If the handler fails,
the container is terminated and restarted according to its restart policy.

A 'preStop' handler runs immediately before a container is terminated. This handler must
complete before the container termination signal is sent. If the handler fails, the
container enters a termination grace period and is then forcefully terminated.

Each lifecycle handler can be configured with one of the following action types:

'httpGet': Performs an HTTP GET request to the container. The action is considered
successful if the response status code is in the 200-399 range. 'httpGet' actions require
a port number to specify where to send the HTTP request.

'exec': Executes a command inside the container. The action is considered successful if
the command exits with status code 0. 'exec' actions accept a list of command arguments
to run inside the container.

'tcpSocket': Attempts to open a TCP socket connection to the container. The action is
considered successful if a connection can be established. 'tcpSocket' actions require a
port number to specify where to attempt the connection.

'sleep': Waits for a specified duration before proceeding. 'sleep' actions accept a
duration in seconds to wait.
`),
		ValidArgsFunction: lifecycleCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			action := args[1]
			lifecycleHandler := &api.LifecycleHandler{}

			switch action {
			case "httpGet":
				headers := []string{}
				if flags.headers != "" {
					headers = strings.Split(flags.headers, ",")
				}
				parsedHeaders, err := parseHeaders(headers)
				if err != nil {
					return fmt.Errorf("could not parse headers: %s", err)
				}

				port, err := strconv.Atoi(args[2])
				if err != nil {
					return fmt.Errorf("could not parse port: %s", err)
				}
				lifecycleHandler.HTTPGet = &api.HTTPGetAction{
					Path:        flags.path,
					Port:        port,
					HTTPHeaders: parsedHeaders,
				}
			case "exec":
				commandArgs := args[2:]
				lifecycleHandler.Exec = &api.ExecAction{
					Command: commandArgs,
				}
			case "sleep":
				seconds, err := strconv.Atoi(args[2])
				if err != nil {
					return fmt.Errorf("could not parse sleep duration: %s", err)
				}
				lifecycleHandler.Sleep = &api.SleepAction{
					Seconds: seconds,
				}
			case "tcpSocket":
				port, err := strconv.Atoi(args[2])
				if err != nil {
					return fmt.Errorf("could not parse port: %s", err)
				}

				lifecycleHandler.TCPSocket = &api.TCPSocketAction{
					Port: port,
				}
			default:
				return fmt.Errorf("invalid action type %s, Must be one of: \"httpGet\", \"sleep\", \"exec\", \"tcpSocket\"", action)
			}
			handler := args[0]
			lifecycle := &api.Lifecycle{StopSignal: flags.stopSignal}
			switch handler {
			case "postStart":
				lifecycle.PostStart = &lifecycleHandler
			case "preStop":
				lifecycle.PreStop = &lifecycleHandler
			default:
				return fmt.Errorf("lifecycle handler %s is invalid. Must be one of 'postStart' or 'preStop'", handler)
			}

			return cmdr.LifecyclesSet(app, flags.ptype, lifecycle)
		},
	}

	cmd.Flags().StringVarP(&flags.ptype, "ptype", "p", "", i18n.T("The ptype for which the lifecycle handler needs to be applied"))
	cmd.MarkFlagRequired("ptype")
	ptypeCompletion := completion.PtsCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("ptype", ptypeCompletion.CompletionFunc)

	cmd.Flags().StringVar(&flags.path, "path", "/", i18n.T("The relative URL path for 'httpGet' actions"))
	cmd.Flags().StringVar(&flags.headers, "headers", "", i18n.T("The HTTP headers to send for 'httpGet' actions, separated by commas"))
	cmd.Flags().StringVar(&flags.stopSignal, "stop-signal", "SIGTERM", i18n.T("The stop signal to send to the container"))
	cmd.Flags().SortFlags = false

	return cmd
}

func lifecyclesUnset(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		app        string
		ptype      string
		lifecycles []string
	}

	lifecycleHandlerCompletion := completion.LifecycleHandlerCompletion{ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "unset <handler>...",
		Args: cobra.MinimumNArgs(1),
		Example: template.CustomExample(
			"drycc lifecycles unset postStart",
			map[string]string{
				"<handler>": i18n.T("the lifecycle handler type, such as 'postStart' or 'preStop'"),
			},
		),
		Short:             i18n.T("Unset lifecycle handler for an application"),
		ValidArgsFunction: lifecycleHandlerCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			flags.lifecycles = args
			for lifecycle := range flags.lifecycles {
				if err := checkLifecycleHandlerType(flags.lifecycles[lifecycle]); err != nil {
					return err
				}
			}
			return cmdr.LifecyclesUnset(app, flags.ptype, flags.lifecycles)
		},
	}
	cmd.Flags().StringVarP(&flags.ptype, "ptype", "p", "", i18n.T("The ptype for which the lifecycle handler needs to be unset"))
	ptypeCompletion := completion.PtsCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("ptype", ptypeCompletion.CompletionFunc)

	return cmd
}

func checkLifecycleHandlerType(lifecycleHandlerType string) error {
	switch lifecycleHandlerType {
	case "postStart", "preStop":
		return nil
	default:
		return fmt.Errorf("lifecycle handler type %s is invalid. Must be one of 'postStart' or 'preStop'", lifecycleHandlerType)
	}
}
