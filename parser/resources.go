package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// Resources commands to their specific function.
func Resources(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for resources:

resources:services         list all available resource services
resources:plans            list all available plans for an resource services
resources:create           create a resource for the application
resources:list             list resources in the application
resources:describe         get a resource detail info in the application
resources:update           update a resource from the application
resources:destroy          delete a resource from the applicationa
resources:bind             bind a resource to servicebroker
resources:unbind           unbind a resource from servicebroker

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "resources:services":
		return resourcesServices(argv, cmdr)
	case "resources:plans":
		return resourcesPlans(argv, cmdr)
	case "resources:create":
		return resourcesCreate(argv, cmdr)
	case "resources:list":
		return resourcesList(argv, cmdr)
	case "resources:describe":
		return resourceGet(argv, cmdr)
	case "resources:update":
		return resourcePut(argv, cmdr)
	case "resources:destroy":
		return resourceDelete(argv, cmdr)
	case "resources:bind":
		return resourceBind(argv, cmdr)
	case "resources:unbind":
		return resourceUnbind(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "resources" {
			argv[0] = "resources:list"
			return resourcesList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func resourcesServices(argv []string, cmdr cmd.Commander) error {
	usage := `
List all available resource services.

Usage: drycc resources:services [options]

Options:
  -l --limit=<num>
    the maximum number of results to display, defaults to config setting
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	results, err := responseLimit(safeGetValue(args, "--limit"))

	if err != nil {
		return err
	}

	return cmdr.ResourcesServices(results)
}

func resourcesPlans(argv []string, cmdr cmd.Commander) error {
	usage := `
List all available plans for an resource services.

Usage: drycc resources:plans <service> [options]

Arguments:
  <service>
    the service name for plans.

Options:
  -l --limit=<num>
    the maximum number of results to display, defaults to config setting
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	service := safeGetValue(args, "<service>")

	results, err := responseLimit(safeGetValue(args, "--limit"))

	if err != nil {
		return err
	}

	return cmdr.ResourcesPlans(service, results)
}

func resourcesCreate(argv []string, cmdr cmd.Commander) error {
	usage := `
Create a resource for the application.

Usage: drycc resources:create <plan> <name> [<param>=<value>...] [options]

Arguments:
  <plan>
    the resource's plan.
  <name>
    this resource instance alias.
  <param>
    the resource instance parameters key.
  <value>
    the resource instance parameters value.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	plan := safeGetValue(args, "<plan>")
	name := safeGetValue(args, "<name>")

	return cmdr.ResourcesCreate(app, plan, name, args["<param>=<value>"].([]string))
}

func resourcesList(argv []string, cmdr cmd.Commander) error {
	usage := `
List resources in the application.

Usage: drycc resources:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  -l --limit=<num>
    the maximum number of results to display, defaults to config setting
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	results, err := responseLimit(safeGetValue(args, "--limit"))

	if err != nil {
		return err
	}
	app := safeGetValue(args, "--app")

	return cmdr.ResourcesList(app, results)
}

func resourceGet(argv []string, cmdr cmd.Commander) error {
	usage := `
get a resource's detail in the application.

Usage: drycc resources:describe <name> [options]

Arguments:
  <name>
    this resource instance alias.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  -l --limit=<num>
    the maximum number of results to display, defaults to config setting
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	name := safeGetValue(args, "<name>")

	return cmdr.ResourceGet(app, name)
}

func resourcePut(argv []string, cmdr cmd.Commander) error {

	usage := `
update a resource from the application

Usage: drycc resources:update <plan> <name> [<param>=<value>...] [options]

Arguments:
  <plan>
    the resource's plan.
  <name>
    this resource instance alias.
  <param>
    the resource instance parameters key.
  <value>
    the resource instance parameters value.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	plan := safeGetValue(args, "<plan>")
	name := safeGetValue(args, "<name>")

	return cmdr.ResourcePut(app, plan, name, args["<param>=<value>"].([]string))
}

func resourceDelete(argv []string, cmdr cmd.Commander) error {

	usage := `
Delete a resource from the application.

Usage: drycc resources:destroy <name> [options]

Arguments:
  <name>
    the resource instance alias name to be removed.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	name := safeGetValue(args, "<name>")

	return cmdr.ResourceDelete(app, name)
}

func resourceBind(argv []string, cmdr cmd.Commander) error {
	usage := `
bind a resource for an application.

Usage: drycc resources:bind <name> [options]

Arguments:
  <name>
    the resource instance alias name.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	name := safeGetValue(args, "<name>")

	return cmdr.ResourceBind(app, name)
}

func resourceUnbind(argv []string, cmdr cmd.Commander) error {
	usage := `
unbind a resources for an application.

Usage: drycc resources:unbind <name> [options]

Arguments:
  <name>
    the resource instance alias name.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	name := safeGetValue(args, "<name>")

	return cmdr.ResourceUnbind(app, name)
}
