package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// Allowlist displays all relevant commands for `drycc allowlist`.
func Allowlist(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for allowlist:

allowlist:add           adds addresses to the application's allowlist
allowlist:list          list addresses in the application's allowlist
allowlist:remove        remove addresses from the application's allowlist

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "allowlist:add":
		return allowlistAdd(argv, cmdr)
	case "allowlist:list":
		return allowlistList(argv, cmdr)
	case "allowlist:remove":
		return allowlistRemove(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "allowlist" {
			argv[0] = "allowlist:list"
			return allowlistList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func allowlistAdd(argv []string, cmdr cmd.Commander) error {
	usage := `
Adds addresses to an application allowlist.

Usage: drycc allowlist:add <addresses> [options]

Arguments:
  <addresses>
    comma-delimited list of addresses(using IP or CIDR notation) to be allowlisted for the application, such as '1.2.3.4' or '1.2.3.4,0.0.0.0/0'.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	addresses := safeGetValue(args, "<addresses>")

	return cmdr.AllowlistAdd(app, addresses)
}

func allowlistList(argv []string, cmdr cmd.Commander) error {
	usage := `
Lists allowlisted addresses for an application.

Usage: drycc allowlist:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")

	return cmdr.AllowlistList(app)
}

func allowlistRemove(argv []string, cmdr cmd.Commander) error {
	usage := `
Removes addresses from an application allowlist.

Usage: drycc allowlist:remove <addresses> [options]

Arguments:
  <addresses>
    comma-delimited list of addresses(using IP or CIDR notation) to be allowlisted for the application, such as '1.2.3.4' or "1.2.3.4,0.0.0.0/0".

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	addresses := safeGetValue(args, "<addresses>")

	return cmdr.AllowlistRemove(app, addresses)
}
