package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// Keys routes key commands to the specific function.
func Keys(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for SSH keys:

keys:list        list SSH keys for the logged in user
keys:add         add an SSH key
keys:remove      remove an SSH key

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "keys:list":
		return keysList(argv, cmdr)
	case "keys:add":
		return keyAdd(argv, cmdr)
	case "keys:remove":
		return keyRemove(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "keys" {
			argv[0] = "keys:list"
			return keysList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func keysList(argv []string, cmdr cmd.Commander) error {
	usage := `
Lists SSH keys for the logged in user.

Usage: drycc keys:list [options]

Options:
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

	return cmdr.KeysList(results)
}

func keyAdd(argv []string, cmdr cmd.Commander) error {
	usage := `
Adds SSH keys for the logged in user.

Usage: drycc keys:add [<name>] [<key>]

<name> and <key> can be used in either order and are both optional.

Arguments:
  <name>
    name of the SSH key.
  <key>
    a local file path to an SSH public key used to push application code.
`

	args, err := docopt.ParseArgs(usage, argv, "")
	if err != nil {
		return err
	}

	return cmdr.KeyAdd(safeGetString(args, "<name>"), safeGetString(args, "<key>"))
}

func keyRemove(argv []string, cmdr cmd.Commander) error {
	usage := `
Removes an SSH key for the logged in user.

Usage: drycc keys:remove <key>

Arguments:
  <key>
    the SSH public key to revoke source code push access.
`

	args, err := docopt.ParseArgs(usage, argv, "")
	if err != nil {
		return err
	}

	return cmdr.KeyRemove(safeGetString(args, "<key>"))
}
