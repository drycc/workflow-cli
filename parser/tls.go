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
tls:auto:issuer       set automatic certificate management environment issuer
tls:auto:enable       enable automatic certificate management environment
tls:auto:disable      disable automatic certificate management environment
tls:force:enable      enable https redirects all your visitor requests from http to https
tls:force:disable     disable https redirects all your visitor requests from http to https

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "tls:info":
		return tlsInfo(argv, cmdr)
	case "tls:auto:issuer":
		return tlsAutoIssuer(argv, cmdr)
	case "tls:auto:enable":
		return tlsAutoEnable(argv, cmdr)
	case "tls:auto:disable":
		return tlsAutoDisable(argv, cmdr)
	case "tls:force:enable":
		return tlsForceEnable(argv, cmdr)
	case "tls:force:disable":
		return tlsForceDisable(argv, cmdr)
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

	return cmdr.TLSInfo(safeGetString(args, "--app"))
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

	return cmdr.TLSForceEnable(safeGetString(args, "--app"))
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

	return cmdr.TLSForceDisable(safeGetString(args, "--app"))
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

	return cmdr.TLSAutoEnable(safeGetString(args, "--app"))
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

	return cmdr.TLSAutoDisable(safeGetString(args, "--app"))
}

func tlsAutoIssuer(argv []string, cmdr cmd.Commander) error {
	usage := `
Disable certs-auto requests to current application.

Usage: drycc tls:auto:issuer --email=<email> --server=<server> --key-id=<key-id> --key-secret=<key-secret> [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  --email=<email>
    the email address to be associated with the ACME account.
  --server=<server>
    Server is the URL used to access the ACME server's 'directory' endpoint.
  --key-id=<key-id>
    keyID is the ID of the CA key that the External Account is bound to.
  --key-secret=<key-secret>
    keySecret holds the symmetric MAC key of the External Account Binding.

`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}
	app := safeGetString(args, "--app")
	email := safeGetString(args, "--email")
	server := safeGetString(args, "--server")
	keyID := safeGetString(args, "--key-id")
	keySecret := safeGetString(args, "--key-secret")

	return cmdr.TLSAutoIssuer(app, email, server, keyID, keySecret)
}
