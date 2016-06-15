package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/deis/workflow-cli/cmd"

	"github.com/deis/controller-sdk-go/api"
	docopt "github.com/docopt/docopt-go"
)

// Healthchecks routes ealthcheck commands to their specific function
func Healthchecks(argv []string) error {
	usage := `
Valid commands for healthchecks:

healthchecks:list        list healthchecks for an app
healthchecks:set         set healthchecks for an app
healthchecks:unset       unset healthchecks for an app

Use 'deis help [command]' to learn more.
`

	switch argv[0] {
	case "healthchecks:list":
		return healthchecksList(argv)
	case "healthchecks:set":
		return healthchecksSet(argv)
	case "healthchecks:unset":
		return healthchecksUnset(argv)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "healthchecks" {
			argv[0] = "healthchecks:list"
			return healthchecksList(argv)
		}

		PrintUsage()
		return nil
	}
}

func healthchecksList(argv []string) error {
	usage := `
Lists healthchecks for an application.

Usage: deis healthchecks:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable name of the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmd.HealthchecksList(safeGetValue(args, "--app"))
}

func healthchecksSet(argv []string) error {
	usage := `
Sets healthchecks for an application.

By default, Workflow only checks that the application starts in their Container. A health
check may be added by configuring a health check probe for the application. The health
checks are implemented as Kubernetes Container Probes. A 'liveness' and a 'readiness'
probe can be configured, and each probe can be of type 'httpGet', 'exec' or 'tcpSocket'
depending on the type of probe the Container requires.

A 'liveness' probe is useful for applications running for long periods of time, eventually
transitioning to broken states and cannot recover except by restarting them.

Other times, a 'readiness' probe is useful when the Container is only temporarily unable
to serve, and will recover on its own. In this case, if a Container fails its 'readiness'
probe, the Container will not be shut down, but rather the Container will stop receiving
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

Usage: deis healthchecks:set <health-type> <probe-type> [options] [--] <args>...

Arguments:
  <health-type>
    the healthcheck type, such as 'liveness' or 'readiness'.
  <probe-type>
    the healthcheck probe type, such as 'httpGet', 'exec' or 'tcpSocket'.
  <args>
    The arguments required for the healthcheck probe. 'exec', accepts a list of arguments;
    'httpGet' and 'tcpSocket' accept a port number.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  -p --path=<path>
    the relative URL path for 'httpGet' probes. [default: /]
  --header=<header>...
    the HTTP headers to send for 'httpGet' probes, separated by commas.
  --initial-delay-timeout=<initial-delay-timeout>
    the initial delay timeout for the probe [default: 50]
  --timeout-seconds=<timeout-seconds>
    the number of seconds after which the probe times out [default: 50]
  --period-seconds=<period-seconds>
    how often (in seconds) to perform the probe [default: 10]
  --success-threshold=<success-threshold>
    minimum consecutive successes for the probe to be considered successful after having failed [default: 1]
  --failure-threshold=<failure-threshold>
    minimum consecutive successes for the probe to be considered failed after having succeeded [default: 3]
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	path := safeGetValue(args, "--path")
	initialDelayTimeout := safeGetInt(args, "--initial-delay-timeout")
	timeoutSeconds := safeGetInt(args, "--timeout-seconds")
	periodSeconds := safeGetInt(args, "--period-seconds")
	successThreshold := safeGetInt(args, "--success-threshold")
	failureThreshold := safeGetInt(args, "--failure-threshold")
	headers := []string{}
	if args["--headers"] != nil {
		headers = args["--headers"].([]string)
	}

	healthcheckType := args["<health-type>"].(string)
	probeType := args["<probe-type>"].(string)
	probeArgs := args["<args>"].([]string)

	if healthcheckType != "liveness" && healthcheckType != "readiness" {
		return fmt.Errorf("Invalid healthcheck type. Must be one of: \"liveness\", \"readiness\"")
	}

	// NOTE(bacongobbler): k8s healthchecks use the term "livenessProbe" and "readinessProbe", so let's
	// add that to the end of the healthcheck type so the controller sees the right probe type
	healthcheckType += "Probe"

	probe := &api.Healthcheck{
		InitialDelaySeconds: initialDelayTimeout,
		TimeoutSeconds:      timeoutSeconds,
		PeriodSeconds:       periodSeconds,
		SuccessThreshold:    successThreshold,
		FailureThreshold:    failureThreshold,
	}

	switch probeType {
	case "httpGet":
		parsedHeaders, err := parseHeaders(headers)
		if err != nil {
			return fmt.Errorf("could not parse headers: %s", err)
		}
		port, err := strconv.Atoi(probeArgs[0])
		if err != nil {
			return fmt.Errorf("could not parse port: %s", err)
		}
		probe.HTTPGet = &api.HTTPGetProbe{
			Path:        path,
			Port:        port,
			HTTPHeaders: parsedHeaders,
		}
	case "exec":
		probe.Exec = &api.ExecProbe{
			Command: probeArgs,
		}
	case "tcpSocket":
		port, err := strconv.Atoi(probeArgs[0])
		if err != nil {
			return fmt.Errorf("could not parse port: %s", err)
		}
		probe.TCPSocket = &api.TCPSocketProbe{
			Port: port,
		}
	default:
		return fmt.Errorf("Invalid probe type. Must be one of: \"httpGet\", \"exec\"")
	}
	return cmd.HealthchecksSet(app, healthcheckType, probe)
}

func healthchecksUnset(argv []string) error {
	usage := `
Unsets healthchecks for an application.

Usage: deis healthchecks:unset [options] <type>...

Arguments:
  <type>
    the healthcheck type, such as 'liveness' or 'readiness'.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	healthchecks := args["<type>"].([]string)

	// NOTE(bacongobbler): k8s healthchecks use the term "livenessProbe" and "readinessProbe", so let's
	// add that to the end of the healthcheck type so the controller sees the right probe type
	for healthcheck := range healthchecks {
		healthchecks[healthcheck] += "Probe"
	}

	return cmd.HealthchecksUnset(app, healthchecks)
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
		Key:   strings.TrimSpace(headerParts[0]),
		Value: strings.TrimSpace(headerParts[1]),
	}, nil
}
