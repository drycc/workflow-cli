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

var healthchecksFlags struct {
	ptype string
}

// NewHealthchecksCommand creates a command for managing application healthchecks.
func NewHealthchecksCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "healthchecks",
		Short: i18n.T("Manage application healthchecks"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.HealthchecksList(app, healthchecksFlags.ptype, version)
		},
	}

	cmd.PersistentFlags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name of the application"))
	cmd.PersistentFlags().StringVarP(&healthchecksFlags.ptype, "ptype", "p", "", i18n.T("The ptype for which the health check needs to be listed"))
	cmd.Flags().IntVarP(&version, "version", "v", 0, i18n.T("The version for which the health check needs to be listed"))

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)
	ptypeCompletion := completion.PtsCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("ptype", ptypeCompletion.CompletionFunc)
	releaseCompletion := completion.ReleaseCompletion{ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("version", releaseCompletion.CompletionFunc)

	cmd.AddCommand(healthchecksList(cmdr))
	cmd.AddCommand(healthchecksSet(cmdr))
	cmd.AddCommand(healthchecksUnset(cmdr))
	return cmd
}

func healthchecksList(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: i18n.T("List healthchecks for an application"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.HealthchecksList(app, healthchecksFlags.ptype, version)
		},
	}

	cmd.Flags().IntVarP(&version, "version", "v", 0, i18n.T("The version for which the health check needs to be listed"))

	releaseCompletion := completion.ReleaseCompletion{ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("version", releaseCompletion.CompletionFunc)
	return cmd
}

func healthchecksSet(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		ptype         string
		healthType    string
		probeType     string
		path          string
		port          int
		headers       string
		initialDelay  int
		timeout       int
		period        int
		successThresh int
		failureThresh int
		commandArgs   []string
	}

	healthChecksCompletion := completion.HealthChecksCompletion{ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "set <health-type> <probe-type> [flags] [--] <args>...",
		Args: cobra.MinimumNArgs(2),
		Example: template.CustomExample(
			"drycc healthchecks set readinessProbe httpGet --path=/health -- 8000",
			map[string]string{
				"<health-type>": i18n.T("The healthcheck type, such as 'startupProbe' 'livenessProbe' or 'readinessProbe'"),
				"<probe-type>":  i18n.T("the healthcheck probe type, such as 'httpGet', 'exec' or 'tcpSocket'"),
				"<args>": i18n.T(`The arguments required for the healthcheck probe. 'exec', accepts a list of arguments;
                  'httpGet' and 'tcpSocket' accept a port number.`),
			},
		),
		Short: i18n.T("Set healthchecks for an application"),
		Long: i18n.T(`Sets healthchecks for an application.

By default, Workflow only checks that the application starts in their Container. A health
check may be added by configuring a health check probe for the application. The health
checks are implemented as Kubernetes Container Probes. A 'startupProbe' 'livenessProbe' 
and a 'readinessProbe' can be configured, and each probe can be of type 'httpGet', 'exec' 
or 'tcpSocket' depending on the type of probe the Container requires.

A 'startupProbe' indicates whether the application within the container is started.
All other probes are disabled if a startup probe is provided, until it succeeds.
If the startup probe fails, the container is subjected to its restart policy.

A 'livenessProbe' is useful for applications running for long periods of time, eventually
transitioning to broken states and cannot recover except by restarting them.

Other times, a 'readinessProbe' is useful when the Container is only temporarily unable
to serve, and will recover on its own. In this case, if a Container fails its 'readinessProbe'
, the Container will not be shut down, but rather the Container will stop receiving
incoming requests.

'httpGet' probes are just as it sounds: it performs a HTTP GET operation on the Container.
A response code inside the 200-399 range is considered a pass. 'httpGet' probes accept a
port number to perform the HTTP GET operation on the Container.

'exec' probes run a command inside the Container to determine its health. An exit code of
zero is considered a pass, while a non-zero status code is considered a fail. 'exec'
probes accept a string of arguments to be run inside the Container.

'tcpSocket' probes attempt to open a socket in the Container. The Container is only
considered healthy if the check can establish a connection. 'tcpSocket' probes accept a
port number to perform the socket connection on the Container.
`),
		ValidArgsFunction: healthChecksCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			healthType := args[0]
			probeType := args[1]

			probe := &api.Healthcheck{
				InitialDelaySeconds: flags.initialDelay,
				TimeoutSeconds:      flags.timeout,
				PeriodSeconds:       flags.period,
				SuccessThreshold:    flags.successThresh,
				FailureThreshold:    flags.failureThresh,
			}

			switch probeType {
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
				probe.HTTPGet = &api.HTTPGetProbe{
					Path:        flags.path,
					Port:        port,
					HTTPHeaders: parsedHeaders,
				}
			case "exec":
				probeArgs := args[2:]
				probe.Exec = &api.ExecProbe{
					Command: probeArgs,
				}
			case "tcpSocket":
				port, err := strconv.Atoi(args[2])
				if err != nil {
					return fmt.Errorf("could not parse port: %s", err)
				}

				probe.TCPSocket = &api.TCPSocketProbe{
					Port: port,
				}
			default:
				return fmt.Errorf("invalid probe type %s, Must be one of: \"httpGet\", \"exec\", \"tcpSocket\"", probeType)
			}

			return cmdr.HealthchecksSet(app, healthType, healthchecksFlags.ptype, probe)
		},
	}

	cmd.Flags().StringVarP(&healthchecksFlags.ptype, "ptype", "p", "", i18n.T("The ptype for which the health check needs to be applied"))
	ptypeCompletion := completion.PtsCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("ptype", ptypeCompletion.CompletionFunc)

	cmd.Flags().StringVar(&flags.path, "path", "/", i18n.T("The relative URL path for 'httpGet' probes"))
	cmd.Flags().StringVar(&flags.headers, "headers", "", i18n.T("The HTTP headers to send for 'httpGet' probes, separated by commas"))
	cmd.Flags().IntVar(&flags.initialDelay, "initial-delay-timeout", 50, i18n.T("The initial delay timeout for the probe"))
	cmd.Flags().IntVar(&flags.timeout, "timeout-seconds", 50, i18n.T("The number of seconds after which the probe times out"))
	cmd.Flags().IntVar(&flags.period, "period-seconds", 10, i18n.T("How often (in seconds) to perform the probe"))
	cmd.Flags().IntVar(&flags.successThresh, "success-threshold", 1, i18n.T("Minimum consecutive successes for the probe to be considered successful after having failed"))
	cmd.Flags().IntVar(&flags.failureThresh, "failure-threshold", 3, i18n.T("Minimum consecutive successes for the probe to be considered failed after having succeeded"))
	cmd.Flags().SortFlags = false

	return cmd
}

func healthchecksUnset(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		app     string
		ptype   string
		healths []string
	}

	healthTypeCompletion := completion.HealthTypeCompletion{ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "unset <health-type>...",
		Args: cobra.MinimumNArgs(1),
		Example: template.CustomExample(
			"drycc healthchecks unset startupProbe",
			map[string]string{
				"<health-type>": i18n.T("the healthcheck type, such as 'startupProbe' 'livenessProbe' or 'readinessProbe'"),
			},
		),
		Short:             i18n.T("Unset healthchecks for an application"),
		ValidArgsFunction: healthTypeCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			flags.healths = args
			for healthcheck := range flags.healths {
				if err := checkProbeType(flags.healths[healthcheck]); err != nil {
					return err
				}
			}
			return cmdr.HealthchecksUnset(app, healthchecksFlags.ptype, flags.healths)
		},
	}
	cmd.Flags().StringVarP(&healthchecksFlags.ptype, "ptype", "p", "", i18n.T("The ptype for which the health check needs to be removed"))
	ptypeCompletion := completion.PtsCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("ptype", ptypeCompletion.CompletionFunc)

	return cmd
}

func parseHeaders(headers []string) ([]*api.KVPair, error) {
	var parsedHeaders []*api.KVPair
	for _, header := range headers {
		parsedHeader, err := parseHeader(header)
		if err != nil {
			return nil, err
		}
		parsedHeaders = append(parsedHeaders, parsedHeader)
	}
	return parsedHeaders, nil
}

func parseHeader(header string) (*api.KVPair, error) {
	headerParts := strings.SplitN(header, ":", 2)
	if len(headerParts) != 2 {
		return nil, fmt.Errorf("could not find separator in header (%s)", header)
	}
	return &api.KVPair{
		Name:  strings.TrimSpace(headerParts[0]),
		Value: strings.TrimSpace(headerParts[1]),
	}, nil
}

func checkProbeType(probe string) error {
	var found bool
	probeTypes := []string{
		"startupProbe",
		"livenessProbe",
		"readinessProbe",
	}
	for _, probeType := range probeTypes {
		if probe == probeType {
			found = true
		}
	}
	if !found {
		return fmt.Errorf("probe type %s is invalid. Must be one of %s", probe, probeTypes)
	}
	return nil
}
