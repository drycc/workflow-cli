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
tls:force:enable      enables the router to enforce https-only requests to an application
tls:force:disable     disables the router to enforce https-only requests to an application
tls:auto:enable       enables the router to automatic generation of certificates to an application
tls:auto:disable      disables the router to automatic generation of certificates to an application

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "tls:info":
		return tlsInfo(argv, cmdr)
	case "tls:force:enable":
		return tlsEnable(argv, cmdr)
	case "tls:force:disable":
		return tlsDisable(argv, cmdr)
	case "tls:auto:enable":
		return tlsAutoEnable(argv, cmdr)
	case "tls:auto:disable":
		return tlsAutoDisable(argv, cmdr)
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

func tlsAutoEnable(argv []string, cmdr cmd.Commander) error {
	usage := `
Enable certs-auto requests to current application.

Usage: drycc tls:auto:enable [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmdr.TLSAutoEnable(safeGetValue(args, "--app"))
}

func tlsAutoDisable(argv []string, cmdr cmd.Commander) error {
	usage := `
Disable certs-auto requests to current application.

Usage: drycc tls:auto:disable [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmdr.TLSAutoDisable(safeGetValue(args, "--app"))
}
