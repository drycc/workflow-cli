package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// Certs routes certs commands to their specific function.
func Certs(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for certs:

certs:list            list SSL certificates for an app
certs:add             add an SSL certificate to an app
certs:remove          remove an SSL certificate from an app
certs:info            get detailed informaton about the certificate
certs:attach          attach an SSL certificate to a domain
certs:detach          detach an SSL certificate from a domain

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "certs:list":
		return certsList(argv, cmdr)
	case "certs:add":
		return certAdd(argv, cmdr)
	case "certs:remove":
		return certRemove(argv, cmdr)
	case "certs:info":
		return certInfo(argv, cmdr)
	case "certs:attach":
		return certAttach(argv, cmdr)
	case "certs:detach":
		return certDetach(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "certs" {
			argv[0] = "certs:list"
			return certsList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func certsList(argv []string, cmdr cmd.Commander) error {
	usage := `
Show certificate information for an SSL application.

Usage: drycc certs:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  -l --limit=<num>
    the maximum number of results to display, defaults to config setting
`

	args, err := docopt.ParseArgs(usage, argv, "")
	if err != nil {
		return err
	}

	results, err := responseLimit(safeGetString(args, "--limit"))
	if err != nil {
		return err
	}
	app := safeGetString(args, "--app")
	return cmdr.CertsList(app, results)
}

func certAdd(argv []string, cmdr cmd.Commander) error {
	usage := `
Binds a certificate/key pair to an application.

Usage: drycc certs:add <name> <cert> <key> [options]

Arguments:
  <name>
    Name of the certificate to reference it by.
  <cert>
    The public key of the SSL certificate.
  <key>
    The private key of the SSL certificate.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")
	if err != nil {
		return err
	}

	name := args["<name>"].(string)
	cert := args["<cert>"].(string)
	key := args["<key>"].(string)
	app := safeGetString(args, "--app")
	return cmdr.CertAdd(app, cert, key, name)
}

func certRemove(argv []string, cmdr cmd.Commander) error {
	usage := `
Remove a certificate/key pair from the application.

Usage: drycc certs:remove <name> [options]

Arguments:
  <name>
    the name of the cert to remove from the app.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")
	if err != nil {
		return err
	}
	app := safeGetString(args, "--app")
	name := safeGetString(args, "<name>")
	return cmdr.CertRemove(app, name)
}

func certInfo(argv []string, cmdr cmd.Commander) error {
	usage := `
Fetch more detailed information about a certificate.

Usage: drycc certs:info <name> [options]

Arguments:
  <name>
    the name of the cert to get information from.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")
	if err != nil {
		return err
	}
	app := safeGetString(args, "--app")
	name := safeGetString(args, "<name>")
	return cmdr.CertInfo(app, name)
}

func certAttach(argv []string, cmdr cmd.Commander) error {
	usage := `
Attach a certificate to a domain.

Usage: drycc certs:attach <name> <domain> [options]

Arguments:
  <name>
    name of the certificate to attach domain to.
  <domain>
    common name of the domain to attach to (needs to already be in the system).

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")
	if err != nil {
		return err
	}
	app := safeGetString(args, "--app")
	name := safeGetString(args, "<name>")
	domain := safeGetString(args, "<domain>")
	return cmdr.CertAttach(app, name, domain)
}

func certDetach(argv []string, cmdr cmd.Commander) error {
	usage := `
Detach a certificate from a domain.

Usage: drycc certs:detach <name> <domain> [options]

Arguments:
  <name>
    name of the certificate to deatch from a domain
  <domain>
    common name of the domain to detach from

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")
	if err != nil {
		return err
	}
	app := safeGetString(args, "--app")
	name := safeGetString(args, "<name>")
	domain := safeGetString(args, "<domain>")
	return cmdr.CertDetach(app, name, domain)
}
