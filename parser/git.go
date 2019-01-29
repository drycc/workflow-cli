package parser

import (
	"fmt"

	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// Git routes git commands to their specific function.
func Git(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for git:

git:remote          Adds git remote of application to repository
git:remove          Removes git remote of application from repository

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "git:remote":
		return gitRemote(argv, cmdr)
	case "git:remove":
		return gitRemove(argv, cmdr)
	case "git":
		fmt.Print(usage)
		return nil
	default:
		PrintUsage(cmdr)
		return nil
	}
}

func gitRemote(argv []string, cmdr cmd.Commander) error {
	usage := `
Adds git remote of application to repository

Usage: drycc git:remote [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  -r --remote=REMOTE
    name of remote to create. [default: drycc]
  -f --force
    overwrite remote of the given name if it already exists.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	remote := safeGetValue(args, "--remote")
	force := args["--force"].(bool)

	return cmdr.GitRemote(app, remote, force)
}

func gitRemove(argv []string, cmdr cmd.Commander) error {
	usage := `
Removes git remotes of application from repository.

Usage: drycc git:remove [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmdr.GitRemove(safeGetValue(args, "--app"))
}
