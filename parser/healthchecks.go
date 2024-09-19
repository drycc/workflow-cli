package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/drycc/workflow-cli/cmd"

	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/controller-sdk-go/api"
)

// TODO: This is for supporting backward compatibility and should be removed
// in future when next major version will be released.
const (
	defaultPtype string = "web"
)

// Healthchecks routes ealthcheck commands to their specific function
func Healthchecks(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for healthchecks:

healthchecks:list        list healthchecks for an app
healthchecks:set         set healthchecks for an app
healthchecks:unset       unset healthchecks for an app

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "healthchecks:list":
		return healthchecksList(argv, cmdr)
	case "healthchecks:set":
		return healthchecksSet(argv, cmdr)
	case "healthchecks:unset":
		return healthchecksUnset(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "healthchecks" {
			argv[0] = "healthchecks:list"
			return healthchecksList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func healthchecksList(argv []string, cmdr cmd.Commander) error {
	usage := `
Lists healthchecks for an application.

Usage: drycc healthchecks:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable name of the application.
  --ptype=<ptype>
    the ptype for which the health check needs to be listed.
  --version=<version>
    the version for which the health check needs to be listed.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	ptype := safeGetString(args, "--ptype")
	var version int
	if safeGetString(args, "--version") != "" {
		if version, err = versionFromString(safeGetString(args, "--version")); err != nil {
			return err
		}
	}
	return cmdr.HealthchecksList(app, ptype, version)
}

func healthchecksSet(argv []string, cmdr cmd.Commander) error {
	usage := `
Sets healthchecks for an application.

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

Usage: drycc healthchecks:set <health-type> <probe-type> [options] [--] <args>...

Arguments:
  <health-type>
    the healthcheck type, such as 'startupProbe' 'livenessProbe' or 'readinessProbe'.
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
  --ptype=<ptype>
    the ptype for which the health check needs to be applied.
  --headers=<headers>...
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

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	path := safeGetString(args, "--path")
	ptype := safeGetString(args, "--ptype")
	initialDelayTimeout := safeGetInt(args, "--initial-delay-timeout")
	timeoutSeconds := safeGetInt(args, "--timeout-seconds")
	periodSeconds := safeGetInt(args, "--period-seconds")
	successThreshold := safeGetInt(args, "--success-threshold")
	failureThreshold := safeGetInt(args, "--failure-threshold")
	headers := []string{}
	if args["--headers"] != nil {
		headers = strings.Split(args["--headers"].(string), ",")
	}
	if ptype == "" {
		ptype = defaultPtype
	}

	healthcheckType := args["<health-type>"].(string)
	probeType := args["<probe-type>"].(string)
	probeArgs := args["<args>"].([]string)

	if err := checkProbeType(healthcheckType); err != nil {
		return err
	}

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
		return fmt.Errorf("invalid probe type. Must be one of: \"httpGet\", \"exec\"")
	}

	return cmdr.HealthchecksSet(app, healthcheckType, ptype, probe)
}

func healthchecksUnset(argv []string, cmdr cmd.Commander) error {
	usage := `
Unsets healthchecks for an application.

Usage: drycc healthchecks:unset <health-type>... [options]

Arguments:
  <health-type>
    the healthcheck type, such as 'startupProbe' 'livenessProbe' or 'readinessProbe'.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  --ptype=<ptype>
    the ptype for which the health check needs to be removed.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	healthchecks := args["<health-type>"].([]string)
	ptype := safeGetString(args, "--ptype")
	if ptype == "" {
		ptype = defaultPtype
	}

	for healthcheck := range healthchecks {
		if err := checkProbeType(healthchecks[healthcheck]); err != nil {
			return err
		}
	}

	return cmdr.HealthchecksUnset(app, ptype, healthchecks)
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
