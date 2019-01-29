package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// Routing displays all relevant commands for `drycc routing`.
func Routing(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for routing:

routing:info       view routability of an application
routing:enable     enable routing for an app
routing:disable    disable routing for an app

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "routing:info":
		return routingInfo(argv, cmdr)
	case "routing:enable":
		return routingEnable(argv, cmdr)
	case "routing:disable":
		return routingDisable(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "routing" {
			argv[0] = "routing:info"
			return routingInfo(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func routingInfo(argv []string, cmdr cmd.Commander) error {
	usage := `
Prints info about the current application's routability.

Usage: drycc routing:info [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)
	if err != nil {
		return err
	}

	return cmdr.RoutingInfo(safeGetValue(args, "--app"))
}

func routingEnable(argv []string, cmdr cmd.Commander) error {
	usage := `
Enables routability for an app.

Usage: drycc routing:enable [options]

Options:
  -a --app=<app>
    the uniquely identifiable name of the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)
	if err != nil {
		return err
	}

	return cmdr.RoutingEnable(safeGetValue(args, "--app"))
}

func routingDisable(argv []string, cmdr cmd.Commander) error {
	usage := `
Disables routability for an app.

Usage: drycc routing:disable [options]

Options:
  -a --app=<app>
    the uniquely identifiable name of the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmdr.RoutingDisable(safeGetValue(args, "--app"))
}
