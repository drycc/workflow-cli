package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// Autoscale displays all relevant commands for `drycc autoscale`.
func Autoscale(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for autoscale:

autoscale:list   list autoscale options of an application
autoscale:set    turn on autoscale for an app
autoscale:unset  turn off autoscale for an app

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "autoscale:list":
		return autoscaleList(argv, cmdr)
	case "autoscale:set":
		return autoscaleSet(argv, cmdr)
	case "autoscale:unset":
		return autoscaleUnset(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "autoscale" {
			argv[0] = "autoscale:list"
			return autoscaleList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func autoscaleList(argv []string, cmdr cmd.Commander) error {
	usage := `
Prints a list of autoscale options for the application.

Usage: drycc autoscale:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	return cmdr.AutoscaleList(safeGetString(args, "--app"))
}

func autoscaleSet(argv []string, cmdr cmd.Commander) error {
	usage := `
Set autoscale option per process type for an app.

Usage: drycc autoscale:set <ptype> --min=<min> --max=<max> --cpu-percent=<percent> [options]

Arguments:
  <ptype>
    the process type to add to the application's autoscale settings.
  --min=<min>
	minimum replicas to keep around
  --max=<max>
	max replicas to scale up to
  --cpu-percent=<cpu-percent>
	target CPU utilization

Options:
  -a --app=<app>
    the uniquely identifiable name of the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	ptype := args["<ptype>"].(string)
	app := safeGetString(args, "--app")
	min := safeGetInt(args, "--min")
	max := safeGetInt(args, "--max")
	CPUPercent := safeGetInt(args, "--cpu-percent")

	return cmdr.AutoscaleSet(app, ptype, min, max, CPUPercent)
}

func autoscaleUnset(argv []string, cmdr cmd.Commander) error {
	usage := `
Unset autoscale per process type for an app.

Usage: drycc autoscale:unset <ptype> [options]

Arguments:
  <ptype>
    the process type to remove from the application's autoscale settings.

Options:
  -a --app=<app>
    the uniquely identifiable name of the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	ptype := args["<ptype>"].(string)
	app := safeGetString(args, "--app")

	return cmdr.AutoscaleUnset(app, ptype)
}
