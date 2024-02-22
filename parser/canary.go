package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// Canary displays all relevant commands for `drycc autoscale`.
func Canary(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for canary:

canary:info     list canary options of an application
canary:create   turn on canary for an app procfile type
canary:remove   turn off canary for an app procfile type
canary:release  release canary deploy for an app

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "canary:info":
		return canaryInfo(argv, cmdr)
	case "canary:create":
		return canaryCreate(argv, cmdr)
	case "canary:remove":
		return canaryRemove(argv, cmdr)
	case "canary:release":
		return canaryRelease(argv, cmdr)
	case "canary:rollback":
		return canaryRollback(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "canary" {
			argv[0] = "canary:info"
			return canaryInfo(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func canaryInfo(argv []string, cmdr cmd.Commander) error {
	usage := `
Prints a list of canary options for the application.

Usage: drycc canary:info [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	return cmdr.CanaryInfo(safeGetString(args, "--app"))
}

func canaryCreate(argv []string, cmdr cmd.Commander) error {
	usage := `
Set canary option type for an app.

Usage: drycc canary:create <process-type>... [options]

Arguments:
  <process-type>
    the process type to add to the application's canary settings.

Options:
  -a --app=<app>
    the uniquely identifiable name of the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	return cmdr.CanaryCreate(app, args["<process-type>"].([]string))
}

func canaryRemove(argv []string, cmdr cmd.Commander) error {
	usage := `
Remove canary option type for an app.

Usage: drycc canary:remove <process-type>... [options]

Arguments:
  <process-type>
    the process type to add to the application's canary settings.

Options:
  -a --app=<app>
    the uniquely identifiable name of the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	return cmdr.CanaryRemove(app, args["<process-type>"].([]string))
}

func canaryRelease(argv []string, cmdr cmd.Commander) error {
	usage := `
Release canary deploy for an app.

Usage: drycc canary:release [options]

Options:
  -a --app=<app>
    the uniquely identifiable name of the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	return cmdr.CanaryRelease(app)
}

func canaryRollback(argv []string, cmdr cmd.Commander) error {
	usage := `
Rollback canary deploy for an app.

Usage: drycc canary:rollback [options]

Options:
  -a --app=<app>
    the uniquely identifiable name of the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	return cmdr.CanaryRollback(app)
}
