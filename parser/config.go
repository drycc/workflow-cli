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
config:pull        extract environment variables to .env
config:push        set environment variables from .env

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
  --oneline
    print output on one line.
  --diff
    print output on multiple lines for comparison against .env files.
  -a --app=<app>
    the uniquely identifiable name of the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}
	app := safeGetValue(args, "--app")
	oneline := args["--oneline"].(bool)
	diff := args["--diff"].(bool)

	format := ""
	if oneline {
		format = "oneline"
	} else if diff {
		format = "diff"
	}

	return cmdr.ConfigList(app, format)
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
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")

	return cmdr.ConfigSet(app, args["<var>=<value>"].([]string))
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
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}
	app := safeGetValue(args, "--app")

	return cmdr.ConfigUnset(app, args["<key>"].([]string))
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
    The application that you wish to pull from
  -i --interactive
    Prompts for each value to be overwritten
  -o --overwrite
    Allows you to have the pull overwrite keys in .env
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	interactive := args["--interactive"].(bool)
	overwrite := args["--overwrite"].(bool)

	return cmdr.ConfigPull(app, interactive, overwrite)
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
    the uniquely identifiable name for the application.
  -p <path>, --path=<path>
    a path leading to an environment file [default: .env]
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	path := safeGetValue(args, "--path")

	return cmdr.ConfigPush(app, path)
}
