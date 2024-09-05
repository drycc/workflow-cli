package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// Autodeploy displays all relevant commands for `drycc autodeploy`.
func Autodeploy(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for autodeploy:

autodeploy:info       view autodeploy info of an application
autodeploy:enable     enable autodeploy for an app
autodeploy:disable    disable autodeploy for an app

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "autodeploy:info":
		return autodeployInfo(argv, cmdr)
	case "autodeploy:enable":
		return autodeployEnable(argv, cmdr)
	case "autodeploy:disable":
		return autodeployDisable(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "autodeploy" {
			argv[0] = "autodeploy:info"
			return autodeployInfo(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func autodeployInfo(argv []string, cmdr cmd.Commander) error {
	usage := `
Prints info about the current application's autodeploy if or not.

Usage: drycc autodeploy:info [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")
	if err != nil {
		return err
	}

	return cmdr.AutodeployInfo(safeGetString(args, "--app"))
}

func autodeployEnable(argv []string, cmdr cmd.Commander) error {
	usage := `
Enables autodeploy for an app.

Usage: drycc autodeploy:enable [options]

Options:
  -a --app=<app>
    the uniquely identifiable name of the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")
	if err != nil {
		return err
	}

	return cmdr.AutodeployEnable(safeGetString(args, "--app"))
}

func autodeployDisable(argv []string, cmdr cmd.Commander) error {
	usage := `
Disables autodeploy for an app.

Usage: drycc autodeploy:disable [options]

Options:
  -a --app=<app>
    the uniquely identifiable name of the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	return cmdr.AutodeployDisable(safeGetString(args, "--app"))
}
