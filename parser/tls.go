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
tls:auto:issuer       add a issuer to an application

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "tls:info":
		return tlsInfo(argv, cmdr)
	case "tls:force:enable":
		return tlsForceEnable(argv, cmdr)
	case "tls:force:disable":
		return tlsForceDisable(argv, cmdr)
	case "tls:auto:enable":
		return tlsAutoEnable(argv, cmdr)
	case "tls:auto:disable":
		return tlsAutoDisable(argv, cmdr)
	case "tls:auto:issuer":
		return tlsAutoIssuer(argv, cmdr)
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

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	return cmdr.TLSInfo(safeGetValue(args, "--app"))
}

func tlsForceEnable(argv []string, cmdr cmd.Commander) error {
	usage := `
Enable the router to enforce https-only requests to the current application.

Usage: drycc tls:force:enable [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	return cmdr.TLSForceEnable(safeGetValue(args, "--app"))
}

func tlsForceDisable(argv []string, cmdr cmd.Commander) error {
	usage := `
Disable the router from enforcing https-only requests to the current application.

Usage: drycc tls:force:disable [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	return cmdr.TLSForceDisable(safeGetValue(args, "--app"))
}

func tlsAutoEnable(argv []string, cmdr cmd.Commander) error {
	usage := `
Enable certs-auto requests to current application.

Usage: drycc tls:auto:enable [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

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

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	return cmdr.TLSAutoDisable(safeGetValue(args, "--app"))
}

func tlsAutoIssuer(argv []string, cmdr cmd.Commander) error {
	usage := `
Disable certs-auto requests to current application.

Usage: drycc tls:auto:issuer --email=<email> --server=<server> --key-id=<key-id> --key-secret=<key-secret> [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  <email>
    the email address to be associated with the ACME account.
  <server>
    Server is the URL used to access the ACME server's 'directory' endpoint.
  <key-id>
    keyID is the ID of the CA key that the External Account is bound to.
  <key-secret>
    keySecret holds the symmetric MAC key of the External Account Binding.

`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}
	app := safeGetValue(args, "--app")
	email := safeGetValue(args, "--email")
	server := safeGetValue(args, "--server")
	keyID := safeGetValue(args, "--key-id")
	keySecret := safeGetValue(args, "--key-secret")

	return cmdr.TLSAutoIssuer(app, email, server, keyID, keySecret)
}
