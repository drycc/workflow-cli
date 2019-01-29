package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// TLS routes tls commands to their specific function.
func TLS(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for tls:

tls:info              view info about an application's TLS settings
tls:enable            enables the router to enforce https-only requests to an application
tls:disable           disables the router to enforce https-only requests to an application

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "tls:info":
		return tlsInfo(argv, cmdr)
	case "tls:enable":
		return tlsEnable(argv, cmdr)
	case "tls:disable":
		return tlsDisable(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "tls" {
			argv[0] = "tls:info"
			return tlsInfo(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func tlsInfo(argv []string, cmdr cmd.Commander) error {
	usage := `
Prints info about the current application's TLS settings.

Usage: drycc tls:info [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmdr.TLSInfo(safeGetValue(args, "--app"))
}

func tlsEnable(argv []string, cmdr cmd.Commander) error {
	usage := `
Enable the router to enforce https-only requests to the current application.

Usage: drycc tls:enable [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmdr.TLSEnable(safeGetValue(args, "--app"))
}

func tlsDisable(argv []string, cmdr cmd.Commander) error {
	usage := `
Disable the router from enforcing https-only requests to the current application.

Usage: drycc tls:disable [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmdr.TLSDisable(safeGetValue(args, "--app"))
}
