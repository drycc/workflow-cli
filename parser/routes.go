package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// Routes routes commands to their specific function.
func Routes(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for routes:

routes:create        create a route for an application
routes:list          list application routes
routes:get           get rule of route
routes:set           set rule of route
routes:attach        attach to gateway
routes:detach        detach to gateway
routes:remove        remove routes from an application

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "routes:create":
		return routesCreate(argv, cmdr)
	case "routes:list":
		return routesList(argv, cmdr)
	case "routes:get":
		return routesGet(argv, cmdr)
	case "routes:set":
		return routesSet(argv, cmdr)
	case "routes:attach":
		return routesAttach(argv, cmdr)
	case "routes:detach":
		return routesDetach(argv, cmdr)
	case "routes:remove":
		return routesRemove(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "routes" {
			argv[0] = "routes:list"
			return routesList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func routesCreate(argv []string, cmdr cmd.Commander) error {
	usage := `
Creates routes for an application, provides a way to route requests

Usage: drycc routes:create <name> --type=<type> --kind=<kind> --port=<port> [options]

Arguments:
  <name>
    the route name.
  <type>
    the process type needs to create route.
  <kind>
    the route kind. Supports HTTPRoute, TCPRoute, UDPRoute, GRPCRoute, and TLSRoute.
  <port>
    the port that should be exposed by this service.

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
	procType := safeGetString(args, "--type")
	kind := safeGetString(args, "--kind")
	port := safeGetInt(args, "--port")

	return cmdr.RoutesCreate(app, name, procType, kind, port)
}

func routesList(argv []string, cmdr cmd.Commander) error {
	usage := `
List routes for an application

Usage: drycc routes:list [options]

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

	return cmdr.RoutesList(app, results)
}

func routesGet(argv []string, cmdr cmd.Commander) error {
	usage := `
Get a route for an application

Usage: drycc routes:get <name> [options]

Arguments:
  <name>
    the route name.

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

	return cmdr.RoutesGet(app, name)
}

func routesSet(argv []string, cmdr cmd.Commander) error {
	usage := `
Set route for an application

Usage: drycc routes:set <name> --rules-file=<rules-file> [options]

Arguments:
  <name>
    the route name.
  <rules-file>
    rules-file is the file name of rule to apply.

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
	rulesFile := safeGetString(args, "--rules-file")

	return cmdr.RoutesSet(app, name, rulesFile)
}

func routesAttach(argv []string, cmdr cmd.Commander) error {
	usage := `
Attach a route for an application

Usage: drycc routes:attach <name> --gateway=<gateway> --port=<port> [options]

Arguments:
  <name>
    the route name.
  <gateway>
    the gateway name.
  <port>
    port is the network port, the gateway listener expects to receive.

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
	gateway := safeGetString(args, "--gateway")

	return cmdr.RoutesAttach(app, name, port, gateway)
}

func routesDetach(argv []string, cmdr cmd.Commander) error {
	usage := `
Detach a route for an application

Usage: drycc routes:detach <name> --gateway=<gateway> --port=<port> [options]

Arguments:
  <name>
    the route name.
  <gateway>
    the gateway name.
  <port>
    port is the network port, the gateway listener expects to receive.

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
	gateway := safeGetString(args, "--gateway")

	return cmdr.RoutesDetach(app, name, port, gateway)
}

func routesRemove(argv []string, cmdr cmd.Commander) error {
	usage := `
Delete a route from the application.

Usage: drycc routes:remove <name> [options]

Arguments:
  <name>
    the route name.

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

	return cmdr.RoutesRemove(app, name)
}
