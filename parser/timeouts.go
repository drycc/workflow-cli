package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// Timeouts routes timeouts commands to their specific function
func Timeouts(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for timeouts:

timeouts:list        list resource timeouts for an app
timeouts:set         set resource timeouts for an app
timeouts:unset       unset resource timeouts for an app

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "timeouts:list":
		return timeoutList(argv, cmdr)
	case "timeouts:set":
		return timeoutSet(argv, cmdr)
	case "timeouts:unset":
		return timeoutUnset(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "timeouts" {
			argv[0] = "timeouts:list"
			return timeoutList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func timeoutList(argv []string, cmdr cmd.Commander) error {
	usage := `
Lists resource timeouts for an application.

Usage: drycc timeouts:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable name of the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	return cmdr.TimeoutsList(safeGetString(args, "--app"))
}

func timeoutSet(argv []string, cmdr cmd.Commander) error {
	usage := `
Sets termination grace period for an application.

Usage: drycc timeouts:set <ptype>=<value>... [options]

Arguments:
  <ptype>
    the process type as defined in your Procfile, such as 'web' or 'worker'.
  <value>
    The value to apply to the process type in seconds.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	timeouts := args["<ptype>=<value>"].([]string)

	return cmdr.TimeoutsSet(app, timeouts)
}

func timeoutUnset(argv []string, cmdr cmd.Commander) error {
	usage := `
Unsets timeouts for an application. Default value (30s) or set by drycc controller

Usage: drycc timeouts:unset <ptype>... [options]

Arguments:
  <ptype>
    the process type as defined in your Procfile, such as 'web' or 'worker'.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	timeouts := args["<ptype>"].([]string)

	return cmdr.TimeoutsUnset(app, timeouts)
}
