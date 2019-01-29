package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// Users routes user commands to the specific function.
func Users(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for users:

users:list        list all registered users

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "users:list":
		return usersList(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "users" {
			argv[0] = "users:list"
			return usersList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func usersList(argv []string, cmdr cmd.Commander) error {
	usage := `
Lists all registered users. Workflow administrators will be marked with a *.
Requires admin privileges.

Usage: drycc users:list [options]

Options:
  -l --limit=<num>
    the maximum number of results to display, defaults to config setting
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)
	if err != nil {
		return err
	}

	results, err := responseLimit(safeGetValue(args, "--limit"))

	if err != nil {
		return err
	}

	return cmdr.UsersList(results)
}
