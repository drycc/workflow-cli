package parser

import (
	docopt "github.com/docopt/docopt-go"

	"github.com/drycc/workflow-cli/cmd"
)

// Gateways gateways commands to their specific function.
func Gateways(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for gateways:

gateways:add           create gateways for an application
gateways:list          list application gateways
gateways:remove        remove gateways from an application

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "gateways:add":
		return gatewaysAdd(argv, cmdr)
	case "gateways:list":
		return gatewaysList(argv, cmdr)
	case "gateways:remove":
		return gatewaysRemove(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "gateways" {
			argv[0] = "gateways:list"
			return gatewaysList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func gatewaysAdd(argv []string, cmdr cmd.Commander) error {
	usage := `
Creates gateways for an application and binds it to allow listener of the main app domain

Usage: drycc gateways:add <name> --port=<port> --protocol=<protocol> [options]

Arguments:
  <name>
    the gateway name.
  <port> 
    port is the network port, the listener expects to receive.
  <protocol>
    protocol specifies the network protocol this listener expects to receive. Supports TCP, UDP, TLS, HTTP, and HTTPS.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	name := safeGetString(args, "<name>")
	port := safeGetInt(args, "--port")
	protocol := safeGetString(args, "--protocol")

	return cmdr.GatewaysAdd(app, name, port, protocol)
}

func gatewaysList(argv []string, cmdr cmd.Commander) error {
	usage := `
Lists gateways for an application

Usage: drycc gateways:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  -l --limit=<num>
    the maximum number of results to display, defaults to config setting
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}
	results, err := responseLimit(safeGetString(args, "--limit"))
	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")

	return cmdr.GatewaysList(app, results)
}

func gatewaysRemove(argv []string, cmdr cmd.Commander) error {
	usage := `
Deletes specific gateway for application

Usage: drycc gateways:remove <name> --port=<port> --protocol=<protocol> [options]

Arguments:
  <name>
    the gateway name.
  <port> 
    port is the network port, the listener expects to receive.
  <protocol>
    protocol specifies the network protocol this listener expects to receive.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}
	app := safeGetString(args, "--app")
	name := safeGetString(args, "<name>")
	port := safeGetInt(args, "--port")
	protocol := safeGetString(args, "--protocol")

	return cmdr.GatewaysRemove(app, name, port, protocol)
}
