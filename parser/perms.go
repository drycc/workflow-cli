package parser

import (
	"fmt"

	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// Perms routes perms commands to their specific function.
func Perms(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for perms:

perms:codes           list all policy codenames
perms:list            list all user permission for objects
perms:create          create a user permission for objects
perms:delete          delete a user permission for objects

Use 'drycc help perms:[command]' to learn more.
`

	switch argv[0] {
	case "perms:codes":
		return permsCodes(argv, cmdr)
	case "perms:list":
		return permsList(argv, cmdr)
	case "perms:create":
		return permCreate(argv, cmdr)
	case "perms:delete":
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

func permsCodes(argv []string, cmdr cmd.Commander) error {
	usage := `
List all object policy codenames.

Usage: drycc perms:codes [options]

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

	return cmdr.PermCodes(results)
}

func permsList(argv []string, cmdr cmd.Commander) error {
	usage := `
List all user permission for objects

Usage: drycc perms:list [options]

Options:
  -c --codename=<codename>
    filter all user permissions by codename
  -l --limit=<num>
    the maximum number of results to display, defaults to config setting
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	codename := safeGetString(args, "--codename")
	results, err := responseLimit(safeGetString(args, "--limit"))

	if err != nil {
		return err
	}

	return cmdr.PermList(codename, results)
}

func permCreate(argv []string, cmdr cmd.Commander) error {
	usage := `
Gives another user permission to use an object.

Usage: drycc perms:create <username> <codename> <uniqueid>

Arguments:
  <username>
    the name of the new user
  <codename>
    the object policy codename
  <uniqueid>
    unique identifier for shared objects
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	username := safeGetString(args, "<username>")
	codename := safeGetString(args, "<codename>")
	uniqueid := safeGetString(args, "<uniqueid>")
	fmt.Println(username, codename, uniqueid)
	return cmdr.PermCreate(codename, uniqueid, username)
}

func permDelete(argv []string, cmdr cmd.Commander) error {
	usage := `
Revokes another user's permission to use an object.

Usage: drycc perms:delete <id>

Arguments:
  <id>
    the id of the user perm.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	id := uint64(safeGetInt(args, "<id>"))
	return cmdr.PermDelete(id)
}
