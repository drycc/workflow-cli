package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"

	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cli"
	"github.com/drycc/workflow-cli/cmd"
	"github.com/drycc/workflow-cli/parser"
)

const extensionPrefix = "drycc-"

// main exits with the return value of Command(os.Args[1:]), deferring all logic to
// a func we can test.
func main() {
	os.Exit(Command(os.Args[1:], os.Stdout, os.Stderr, os.Stdin))
}

// Command routes drycc commands to their proper parser.
func Command(argv []string, wOut io.Writer, wErr io.Writer, wIn io.Reader) int {
	usage := `
The Drycc command-line client issues API calls to a Drycc controller.

Usage: drycc <command> [<args>...]

Options:
  -h --help
    display help information
  -v --version
    display client version
  -c --config=<config>
    path to configuration file. Equivalent to
    setting $DRYCC_PROFILE. Defaults to ~/.drycc/config.json.
    If value is not a filepath, will assume location ~/.drycc/client.json

Auth commands, use 'drycc help auth' to learn more:

  register      register a new user with a controller
  login         login to a controller
  logout        logout from the current controller

Subcommands, use 'drycc help [subcommand]' to learn more:

  apps          manage applications used to provide services
  autoscale     manage autoscale for applications
  builds        manage builds created using 'git push'
  certs         manage SSL endpoints for an app
  config        manage environment variables that define app config
  domains       manage and assign domain names to your applications
  git           manage git for applications
  healthchecks  manage healthchecks for applications
  keys          manage ssh keys used for 'git push' deployments
  labels        manage labels of application
  limits        manage resource limits for your application
  perms         manage permissions for applications
  ps            manage processes inside an app container
  registry      manage private registry information for your application
  releases      manage releases of an application
  routing       manage routability of an application
  tags          manage tags for application containers
  tls           manage TLS settings for applications
  users         manage users
  version       display client version
  services      manage services for your applications
  routes        manage routes for your applications
  gateways      manage gateways for your applications
  timeouts      manage pods termination grace period
  volumes       manage volumes for your applications
  resources     manage resources for your applications

Shortcut commands, use 'drycc shortcuts' to see all:

  create        create a new application
  destroy       destroy an application
  info          view information about the current app
  logs          view aggregated log info for the app
  open          open a URL to the app in a browser
  pull          imports an image and deploys as a new release
  run           run a command in an ephemeral app container
  scale         scale processes by type (web=2, worker=1)

Use 'git push drycc main' to deploy to an application.
`
	// Reorganize some command line flags and commands.
	command, argv := parseArgs(argv)
	// Give docopt an optional final false arg so it doesn't call os.Exit().
	_, err := docopt.Parse(usage, []string{command}, false, "", true, false)

	if err != nil {
		fmt.Fprintln(wErr, err)
		return 1
	}

	if len(argv) == 0 {
		fmt.Fprintln(wErr, "Usage: drycc <command> [<args>...]")
		return 1
	}

	configFlag := getConfigFlag(argv)
	// Don't pass down config flag to parser because it isn't defined there.
	argv = removeConfigFlag(argv)
	cmdr := cmd.DryccCmd{ConfigFile: configFlag, WOut: wOut, WErr: wErr, WIn: wIn}

	// Dispatch the command, passing the argv through so subcommands can
	// re-parse it according to their usage strings.
	switch command {
	case "apps":
		err = parser.Apps(argv, &cmdr)
	case "auth":
		err = parser.Auth(argv, &cmdr)
	case "autoscale":
		err = parser.Autoscale(argv, &cmdr)
	case "builds":
		err = parser.Builds(argv, &cmdr)
	case "certs":
		err = parser.Certs(argv, &cmdr)
	case "config":
		err = parser.Config(argv, &cmdr)
	case "domains":
		err = parser.Domains(argv, &cmdr)
	case "services":
		err = parser.Services(argv, &cmdr)
	case "gateways":
		err = parser.Gateways(argv, &cmdr)
	case "routes":
		err = parser.Routes(argv, &cmdr)
	case "git":
		err = parser.Git(argv, &cmdr)
	case "healthchecks":
		err = parser.Healthchecks(argv, &cmdr)
	case "help":
		fmt.Fprint(os.Stdout, usage)
		return 0
	case "keys":
		err = parser.Keys(argv, &cmdr)
	case "labels":
		err = parser.Labels(argv, &cmdr)
	case "limits":
		err = parser.Limits(argv, &cmdr)
	case "timeouts":
		err = parser.Timeouts(argv, &cmdr)
	case "perms":
		err = parser.Perms(argv, &cmdr)
	case "ps":
		err = parser.Ps(argv, &cmdr)
	case "registry":
		err = parser.Registry(argv, &cmdr)
	case "releases":
		err = parser.Releases(argv, &cmdr)
	case "routing":
		err = parser.Routing(argv, &cmdr)
	case "shortcuts":
		err = parser.Shortcuts(argv, &cmdr)
	case "tags":
		err = parser.Tags(argv, &cmdr)
	case "tls":
		err = parser.TLS(argv, &cmdr)
	case "users":
		err = parser.Users(argv, &cmdr)
	case "version":
		err = parser.Version(argv, &cmdr)
	case "volumes":
		err = parser.Volumes(argv, &cmdr)
	case "resources":
		err = parser.Resources(argv, &cmdr)
	default:
		env := os.Environ()

		binary, err := exec.LookPath(extensionPrefix + command)
		if err != nil {
			parser.PrintUsage(&cmdr)
			return 1
		}

		cmdArgv := prepareCmdArgs(command, argv)

		err = syscall.Exec(binary, cmdArgv, env)
		if err != nil {
			parser.PrintUsage(&cmdr)
			return 1
		}
	}
	if err != nil {
		fmt.Fprintf(wErr, "Error: %v\n", err)
		return 1
	}
	return 0
}

func removeConfigFlag(argv []string) []string {
	var kept []string
	for i, arg := range argv {
		// -- /bin/sh -c --config  condition
		if arg == "--" {
			kept = append(kept, argv[i:]...)
			break
		}

		if arg == "-c" || strings.HasPrefix(arg, "--config=") {
			continue
			// If the previous option is -c, remove the argument as well
		} else if i != 0 && argv[i-1] == "-c" {
			continue
		}

		kept = append(kept, arg)
	}

	return kept
}

func getConfigFlag(argv []string) string {
	for i, arg := range argv {
		// -- /bin/sh -c/ --config=  condition
		if arg == "--" {
			return ""
		}
		if strings.HasPrefix(arg, "--config=") {
			return strings.TrimPrefix(arg, "--config=")
		} else if i != 0 && argv[i-1] == "-c" {
			return arg
		}
	}

	return ""
}

// parseArgs returns the provided args with "--help" as the last arg if need be,
// expands shortcuts and formats commands to be properly routed.
func parseArgs(argv []string) (string, []string) {
	if len(argv) == 1 {
		if argv[0] == "--help" || argv[0] == "-h" {
			// rearrange "drycc --help" as "drycc help"
			argv[0] = "help"
		} else if argv[0] == "--version" || argv[0] == "-v" {
			// rearrange "drycc --version" as "drycc version"
			argv[0] = "version"
		}
	}

	if len(argv) > 1 {
		// Rearrange "drycc help <command>" to "drycc <command> --help".
		if argv[0] == "help" || argv[0] == "--help" || argv[0] == "-h" {
			argv = append(argv[1:], "--help")
		}
	}

	if len(argv) > 0 {
		argv[0] = replaceShortcut(argv[0])

		index := strings.Index(argv[0], ":")

		if index != -1 {
			command := argv[0]
			return command[:index], argv
		}

		return argv[0], argv
	}

	return "", argv
}

// split original command and pass its first element in arguments
func prepareCmdArgs(command string, argv []string) []string {
	cmdArgv := []string{extensionPrefix + command}
	cmdSplit := strings.Split(argv[0], command+":")

	if len(cmdSplit) > 1 {
		cmdArgv = append(cmdArgv, cmdSplit[1])
	}

	return append(cmdArgv, argv[1:]...)
}

func replaceShortcut(command string) string {
	expandedCommand := cli.Shortcuts[command]
	if expandedCommand == "" {
		return command
	}

	return expandedCommand
}
