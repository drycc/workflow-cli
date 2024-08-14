package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// Autorollback displays all relevant commands for `drycc autorollback`.
func Autorollback(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for autorollback:

autorollback:info       view autorollback info of an application
autorollback:enable     enable autorollback for an app
autorollback:disable    disable autorollback for an app

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "autorollback:info":
		return autorollbackInfo(argv, cmdr)
	case "autorollback:enable":
		return autorollbackEnable(argv, cmdr)
	case "autorollback:disable":
		return autorollbackDisable(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "autorollback" {
			argv[0] = "autorollback:info"
			return autorollbackInfo(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func autorollbackInfo(argv []string, cmdr cmd.Commander) error {
	usage := `
Prints info about the current application's autorollback if or not.

Usage: drycc autorollback:info [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")
	if err != nil {
		return err
	}

	return cmdr.AutorollbackInfo(safeGetString(args, "--app"))
}

func autorollbackEnable(argv []string, cmdr cmd.Commander) error {
	usage := `
Enables autorollback for an app.

Usage: drycc autorollback:enable [options]

Options:
  -a --app=<app>
    the uniquely identifiable name of the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")
	if err != nil {
		return err
	}

	return cmdr.AutorollbackEnable(safeGetString(args, "--app"))
}

func autorollbackDisable(argv []string, cmdr cmd.Commander) error {
	usage := `
Disables autorollback for an app.

Usage: drycc autorollback:disable [options]

Options:
  -a --app=<app>
    the uniquely identifiable name of the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	return cmdr.AutorollbackDisable(safeGetString(args, "--app"))
}
