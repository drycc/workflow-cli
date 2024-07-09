package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// Config routes config commands to their specific function.
func Config(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for config:

config:list        list environment variables for an app
config:set         set environment variables for an app
config:unset       unset environment variables for an app
config:pull        pull environment variables to the path
config:push        push environment variables from the path

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "config:list":
		return configList(argv, cmdr)
	case "config:set":
		return configSet(argv, cmdr)
	case "config:unset":
		return configUnset(argv, cmdr)
	case "config:pull":
		return configPull(argv, cmdr)
	case "config:push":
		return configPush(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "config" {
			argv[0] = "config:list"
			return configList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func configList(argv []string, cmdr cmd.Commander) error {
	usage := `
Lists environment variables for an application.

Usage: drycc config:list [options]

Options:
  -a --app=<app>
    the application that you wish to listed.
  --type=<type>
    the procType for which the config needs to be listed.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}
	app := safeGetString(args, "--app")
	procType := safeGetString(args, "--type")
	return cmdr.ConfigList(app, procType)
}

func configSet(argv []string, cmdr cmd.Commander) error {
	usage := `
Sets environment variables for an application.

Usage: drycc config:set <var>=<value> [<var>=<value>...] [options]

Arguments:
  <var>
    the uniquely identifiable name for the environment variable.
  <value>
    the value of said environment variable.

Options:
  -a --app=<app>
    the application that you wish to set.
  --type=<type>
    the procType for which the config needs to be set.
  --confirm=yes
    To proceed, type "yes".
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	procType := safeGetString(args, "--type")
	confirm := safeGetString(args, "--confirm")
	return cmdr.ConfigSet(app, procType, args["<var>=<value>"].([]string), confirm)
}

func configUnset(argv []string, cmdr cmd.Commander) error {
	usage := `
Unsets an environment variable for an application.

Usage: drycc config:unset <key>... [options]

Arguments:
  <key>
    the variable to remove from the application's environment.

Options:
  -a --app=<app>
    the application that you wish to unset.
  --type=<type>
    the procType for which the config needs to be unset.
  --confirm=yes
    To proceed, type "yes".
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}
	app := safeGetString(args, "--app")
	procType := safeGetString(args, "--type")
	confirm := safeGetString(args, "--confirm")
	return cmdr.ConfigUnset(app, procType, args["<key>"].([]string), confirm)
}

func configPull(argv []string, cmdr cmd.Commander) error {
	usage := `
Extract all environment variables from an application for local use.

The environmental variables can be piped into a file, 'drycc config:pull > file',
or stored locally in a file named .env. This file can be
read by foreman to load the local environment for your app.

Usage: drycc config:pull [options]

Options:
  -a --app=<app>
    the application that you wish to pull.
  --type=<type>
    the procType for which the config needs to be pull.
  --path=<path>
    a path leading to an environment file [default: .env]
  -i --interactive
    prompts for each value to be overwritten.
  -o --overwrite
    allows you to have the pull overwrite keys to the path.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	procType := safeGetString(args, "--type")
	path := safeGetValue(args, "--path", ".env")
	interactive := args["--interactive"].(bool)
	overwrite := args["--overwrite"].(bool)

	return cmdr.ConfigPull(app, procType, path, interactive, overwrite)
}

func configPush(argv []string, cmdr cmd.Commander) error {
	usage := `
Sets environment variables for an application.

This file can be read by foreman
to load the local environment for your app. The file should be piped via
stdin, 'drycc config:push < .env', or using the --path option.

Usage: drycc config:push [options]

Options:
  -a --app=<app>
    the application that you wish to push.
  --type=<type>
    the procType for which the config needs to be push.
  --path=<path>
    a path leading to an environment file [default: .env]
  --confirm=yes
    To proceed, type "yes".
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	procType := safeGetString(args, "--type")
	path := safeGetValue(args, "--path", ".env")
	confirm := safeGetString(args, "--confirm")
	return cmdr.ConfigPush(app, procType, path, confirm)
}
