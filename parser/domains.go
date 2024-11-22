package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// Domains routes domain commands to their specific function.
func Domains(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for domains:

domains:add           bind a domain to an application
domains:list          list domains bound to an application
domains:remove        unbind a domain from an application

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "domains:add":
		return domainsAdd(argv, cmdr)
	case "domains:list":
		return domainsList(argv, cmdr)
	case "domains:remove":
		return domainsRemove(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "domains" {
			argv[0] = "domains:list"
			return domainsList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func domainsAdd(argv []string, cmdr cmd.Commander) error {
	usage := `
Binds a domain to an application.

Usage: drycc domains:add <domain> [options]

Arguments:
  <domain>
    the domain name to be bound to the application, such as 'domain.dryccapp.com'.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  -p --ptype=<ptype>
    the ptype type for domain, default[web].
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	ptype := safeGetValue(args, "--ptype", "web")
	domain := safeGetString(args, "<domain>")

	return cmdr.DomainsAdd(app, domain, ptype)
}

func domainsList(argv []string, cmdr cmd.Commander) error {
	usage := `
Lists domains bound to an application.

Usage: drycc domains:list [options]

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

	return cmdr.DomainsList(app, results)
}

func domainsRemove(argv []string, cmdr cmd.Commander) error {
	usage := `
Unbinds a domain for an application.

Usage: drycc domains:remove <domain> [options]

Arguments:
  <domain>
    the domain name to be removed from the application.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	domain := safeGetString(args, "<domain>")

	return cmdr.DomainsRemove(app, domain)
}
