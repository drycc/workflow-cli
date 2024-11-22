package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// Pts routes pts commands to their specific function.
func Pts(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for processes:

pts:list        list application process types
pts:describe    print a detailed description of the selected process type
pts:restart     restart an application or process types
pts:scale       scale process types of replicas (e.g. web=4 worker=2)
pts:clean       clean process types of not used

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "pts:list":
		return ptsList(argv, cmdr)
	case "pts:describe":
		return ptsDescribe(argv, cmdr)
	case "pts:restart":
		return ptsRestart(argv, cmdr)
	case "pts:scale":
		return ptsScale(argv, cmdr)
	case "pts:clean":
		return ptsClean(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "pts" {
			argv[0] = "pts:list"
			return ptsList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func ptsList(argv []string, cmdr cmd.Commander) error {
	usage := `
Lists process types servicing an application.

Usage: drycc pts:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")
	if err != nil {
		return err
	}

	// The 1000 is fake for now until API understands limits
	return cmdr.PtsList(safeGetString(args, "--app"), 1000)
}

func ptsDescribe(argv []string, cmdr cmd.Commander) error {
	usage := `
Print a detailed description of the selected process type.

Usage: drycc pts:describe <ptype> [options]

Arguments:
  <ptype> the process name as defined in your Procfile, such as 'web'.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	ptype := safeGetString(args, "<ptype>")
	return cmdr.PtsDescribe(app, ptype)
}

func ptsRestart(argv []string, cmdr cmd.Commander) error {
	usage := `
Restart an application or process types.

Usage: drycc pts:restart [<ptype>...] [options]

Arguments:
  <ptype>
    the process name as defined in your Procfile, such as 'web' or 'web worker'.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  --confirm=yes
    To proceed, type "yes".
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	apps := safeGetString(args, "--app")
	confirm := safeGetString(args, "--confirm")
	return cmdr.PtsRestart(apps, args["<ptype>"].([]string), confirm)
}

func ptsScale(argv []string, cmdr cmd.Commander) error {
	usage := `
Scales an application's processes by type.

Usage: drycc pts:scale <ptype>=<num>... [options]

Arguments:
  <ptype>
    the process name as defined in your Procfile, such as 'web' or 'worker'.
  <num>
    the number of processes.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	apps := safeGetString(args, "--app")
	return cmdr.PtsScale(apps, args["<ptype>=<num>"].([]string))
}

func ptsClean(argv []string, cmdr cmd.Commander) error {
	usage := `
Clean process types of not used.

Usage: drycc pts:clean <ptype>... [options]

Arguments:
  <ptype>
    the process name as defined in your Procfile, such as 'web' or 'worker'.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	apps := safeGetString(args, "--app")
	return cmdr.PtsScale(apps, args["<ptype>"].([]string))
}
