package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// Tokens commands to the specific function.
func Tokens(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for tokens:

tokens:list          lists tokens visible to the current controller.
tokens:add           add a token for controller authentication.
tokens:remove        remove a token for controller authentication.

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "tokens:list":
		return tokensList(argv, cmdr)
	case "tokens:add":
		return tokensAdd(argv, cmdr)
	case "tokens:remove":
		return tokensRemove(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}
		if argv[0] == "tokens" {
			argv[0] = "tokens:list"
			return tokensList(argv, cmdr)
		}
		if printHelp(argv, usage) {
			return nil
		}

		PrintUsage(cmdr)
		return nil
	}
}

func tokensList(argv []string, cmdr cmd.Commander) error {
	usage := `
Lists tokens visible to the current controller.

Usage: drycc tokens:list [options]

Options:
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
	return cmdr.TokensList(results)
}

func tokensAdd(argv []string, cmdr cmd.Commander) error {
	usage := `
Add a token for controller authentication.

Usage: drycc tokens:add <alias> [options]

Arguments:
  <alias>
  provide a alias for controller authentication token.

Options:
  -u --username=<username>
    provide a username for the account.
  -p --password=<password>
    provide a password for the account.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}
	alias := safeGetString(args, "<alias>")
	username := safeGetString(args, "--username")
	password := safeGetString(args, "--password")
	_, err = cmdr.TokensAdd(nil, username, password, alias, "", true)
	return err
}

func tokensRemove(argv []string, cmdr cmd.Commander) error {
	usage := `
Remove a token for controller authentication.

Usage: drycc tokens:remove <id>

Arguments:
  <id>
  the id of the token for controller authentication.

`

	args, err := docopt.ParseArgs(usage, argv, "")
	if err != nil {
		return err
	}
	return cmdr.TokensRemove(safeGetString(args, "<id>"), "")
}
