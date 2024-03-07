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
users:enable      enable a user
users:disable     disable a user

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "users:list":
		return usersList(argv, cmdr)
	case "users:enable":
		return usersEnable(argv, cmdr)
	case "users:disable":
		return usersDisable(argv, cmdr)
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

	args, err := docopt.ParseArgs(usage, argv, "")
	if err != nil {
		return err
	}

	results, err := responseLimit(safeGetString(args, "--limit"))

	if err != nil {
		return err
	}

	return cmdr.UsersList(results)
}

func usersEnable(argv []string, cmdr cmd.Commander) error {
	usage := `
Enable a user when his status is disabled.
Requires admin privileges.

Usage: drycc users:enable <username>

Arguments:
  <username>
  the username you want to enable.
`

	args, err := docopt.ParseArgs(usage, argv, "")
	if err != nil {
		return err
	}
	username := safeGetString(args, "<username>")

	return cmdr.UsersEnable(username)
}

func usersDisable(argv []string, cmdr cmd.Commander) error {
	usage := `
Disable a user when his status is enabled.
Requires admin privileges.

Usage: drycc users:disable <username>

Arguments:
  <username>
  the username you want to disable.
`

	args, err := docopt.ParseArgs(usage, argv, "")
	if err != nil {
		return err
	}
	username := safeGetString(args, "<username>")

	return cmdr.UsersDisable(username)
}
