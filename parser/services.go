package parser

import (
	docopt "github.com/docopt/docopt-go"

	"github.com/drycc/workflow-cli/cmd"
)

// Services routes service commands to their specific function.
func Services(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for services:

services:add           create service for an application
services:list          list application services
services:remove        remove service from an application

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "services:add":
		return servicesAdd(argv, cmdr)
	case "services:list":
		return servicesList(argv, cmdr)
	case "services:remove":
		return servicesRemove(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "services" {
			argv[0] = "services:list"
			return servicesList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func servicesAdd(argv []string, cmdr cmd.Commander) error {
	usage := `
Creates extra service for an application and binds it to specific route of the main app domain.

Usage: drycc services:add <ptype> <port>:<target> [options]

Arguments:
  <ptype>
    procfile type which should handle the request, e.g. webhooks (should be bind to the port PORT)
    only single extra service per Porcfile type could be created.
  <port>
    the port that will be exposed by this service.
  <target>
    number or name of the port to access on the pods targeted by the service.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  --protocol=<protocol>
    the IP protocol for this port. Supports TCP, UDP, and SCTP. Default is TCP.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	ptype := safeGetString(args, "<ptype>")
	protocol := safeGetString(args, "--protocol")
	ports := safeGetString(args, "<port>:<target>")
	return cmdr.ServicesAdd(app, ptype, ports, protocol)
}

func servicesList(argv []string, cmdr cmd.Commander) error {
	usage := `
Lists extra services for an application.

Usage: drycc services:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")

	return cmdr.ServicesList(app)
}

func servicesRemove(argv []string, cmdr cmd.Commander) error {
	usage := `
Deletes specific extra service for application.

Usage: drycc services:remove <ptype> <port> [options]

Arguments:
  <ptype>
    procfile type which should handle the request, e.g. webhooks (should be bind to the port PORT).
    Only single extra service per Porcfile type could be created.
  <port>
    the port exposed by this service.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  --protocol=<protocol>
    the IP protocol for this port. Supports TCP, UDP, and SCTP. Default is TCP.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	ptype := safeGetString(args, "<ptype>")
	protocol := safeGetString(args, "--protocol")
	port := safeGetInt(args, "<port>")

	return cmdr.ServicesRemove(app, ptype, protocol, port)
}
