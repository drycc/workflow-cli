package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// Perms routes perms commands to their specific function.
func Perms(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for perms:

perms:list            list all user permission
perms:add             create a user permission
perms:update          update a user permission
perms:remove          delete a user permission

Use 'drycc help perms:[command]' to learn more.
`

	switch argv[0] {
	case "perms:list":
		return permsList(argv, cmdr)
	case "perms:add":
		return permCreate(argv, cmdr)
	case "perms:update":
		return permUpdate(argv, cmdr)
	case "perms:remove":
		return permDelete(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "perms" {
			argv[0] = "perms:list"
			return permsList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func permsList(argv []string, cmdr cmd.Commander) error {
	usage := `
List all user permissions.

Usage: drycc perms:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  -l --limit=<num>
    the maximum number of results to display, defaults to config setting.
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
	return cmdr.PermList(app, results)
}

func permCreate(argv []string, cmdr cmd.Commander) error {
	usage := `
Grant permissions to user.

Usage: drycc perms:add <username> <permissions> [options]

Arguments:
  <username>
    the name of the user.
  <permissions>
	comma-delimited list of permissions (view,change,delete).
Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	username := safeGetString(args, "<username>")
	permissions := safeGetString(args, "<permissions>")
	return cmdr.PermCreate(app, username, permissions)
}

func permUpdate(argv []string, cmdr cmd.Commander) error {
	usage := `
Update permissions to user.

Usage: drycc perms:update <username> <permissions> [options]

Arguments:
  <username>
    the name of the user.
  <permissions>
    comma-delimited list of permissions (view,change,delete).
Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	username := safeGetString(args, "<username>")
	permissions := safeGetString(args, "<permissions>")
	return cmdr.PermUpdate(app, username, permissions)
}

func permDelete(argv []string, cmdr cmd.Commander) error {
	usage := `
Delete a user from the app.

Usage: drycc perms:remove <username> [options]

Arguments:
  <username>
    the name of the user.
Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	username := safeGetString(args, "<username>")
	return cmdr.PermDelete(app, username)
}
