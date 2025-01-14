package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// Config routes config commands to their specific function.
func Config(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for config:

config:info        an app config info
config:set         set environment variables for an app
config:unset       unset environment variables for an app
config:pull        pull environment variables to the path
config:push        push environment variables from the path
config:attach      selects a environment groups to attach an app ptype
config:detach      selects a environment groups to detach an app ptype

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "config:info":
		return configInfo(argv, cmdr)
	case "config:set":
		return configSet(argv, cmdr)
	case "config:unset":
		return configUnset(argv, cmdr)
	case "config:pull":
		return configPull(argv, cmdr)
	case "config:push":
		return configPush(argv, cmdr)
	case "config:attach":
		return configAttach(argv, cmdr)
	case "config:detach":
		return configDetach(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "config" {
			argv[0] = "config:info"
			return configInfo(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func configInfo(argv []string, cmdr cmd.Commander) error {
	usage := `
An app config info.

Usage: drycc config:info [options]

Options:
  -a --app=<app>
    the application that you wish to listed.
  -p --ptype=<ptype>
    the ptype for which the config needs to be listed.
  -g --group=<group>
    the group for which the config needs to be listed.
  -v --version=<version>
    the version for which the config needs to be listed.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}
	app := safeGetString(args, "--app")
	ptype := safeGetString(args, "--ptype")
	group := safeGetString(args, "--group")
	var version int
	if safeGetString(args, "--version") != "" {
		if version, err = versionFromString(safeGetString(args, "--version")); err != nil {
			return err
		}
	}

	return cmdr.ConfigInfo(app, ptype, group, version)
}

func configSet(argv []string, cmdr cmd.Commander) error {
	usage := `
Sets environment variables for an application or config group.

Usage: drycc config:set <var>=<value> [<var>=<value>...] [options]

Arguments:
  <var>
    the uniquely identifiable name for the environment variable.
  <value>
    the value of said environment variable.

Options:
  -a --app=<app>
    the application that you wish to set.
  -p --ptype=<ptype>
    the ptype for which the config needs to be set.
  -g --group=<group>
    the group for which the config needs to be set.
  --confirm=yes
    To proceed, type "yes".
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	ptype := safeGetString(args, "--ptype")
	group := safeGetString(args, "--group")
	confirm := safeGetString(args, "--confirm")
	return cmdr.ConfigSet(app, ptype, group, args["<var>=<value>"].([]string), confirm)
}

func configUnset(argv []string, cmdr cmd.Commander) error {
	usage := `
Unsets an environment variable for an application or config group.

Usage: drycc config:unset <key>... [options]

Arguments:
  <key>
    the variable to remove from the application's environment.

Options:
  -a --app=<app>
    the application that you wish to unset.
  -p --ptype=<ptype>
    the ptype for which the config needs to be unset.
  -g --group=<group>
    the group for which the config needs to be unset.
  --confirm=yes
    To proceed, type "yes".
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}
	app := safeGetString(args, "--app")
	ptype := safeGetString(args, "--ptype")
	group := safeGetString(args, "--group")
	confirm := safeGetString(args, "--confirm")
	return cmdr.ConfigUnset(app, ptype, group, args["<key>"].([]string), confirm)
}

func configPull(argv []string, cmdr cmd.Commander) error {
	usage := `
Extract all environment variables from an application or config group. for local use.

The environmental variables can be piped into a file, 'drycc config:pull > file',
or stored locally in a file named .env. This file can be
read by foreman to load the local environment for your app.

Usage: drycc config:pull [options]

Options:
  -a --app=<app>
    the application that you wish to pull.
  -p --ptype=<ptype>
    the ptype for which the config needs to be pull.
  -g --group=<group>
    the group for which the config needs to be push.
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
	ptype := safeGetString(args, "--ptype")
	group := safeGetString(args, "--group")
	path := safeGetValue(args, "--path", ".env")
	interactive := args["--interactive"].(bool)
	overwrite := args["--overwrite"].(bool)

	return cmdr.ConfigPull(app, ptype, group, path, interactive, overwrite)
}

func configPush(argv []string, cmdr cmd.Commander) error {
	usage := `
Sets environment variables for an application or config group.

This file can be read by foreman
to load the local environment for your app. The file should be piped via
stdin, 'drycc config:push < .env', or using the --path option.

Usage: drycc config:push [options]

Options:
  -a --app=<app>
    the application that you wish to push.
  -p --ptype=<ptype>
    the ptype for which the config needs to be push.
  -g --group=<group>
    the group for which the config needs to be push.
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
	ptype := safeGetString(args, "--ptype")
	group := safeGetString(args, "--group")
	path := safeGetValue(args, "--path", ".env")
	confirm := safeGetString(args, "--confirm")
	return cmdr.ConfigPush(app, ptype, group, path, confirm)
}

func configAttach(argv []string, cmdr cmd.Commander) error {
	usage := `
Selects a environment groups to attach an app ptype.

Usage: drycc config:attach <ptype> <groups> [options]

Arguments:
  <ptype>
    the ptype that requires attach configurations.
  <groups>
    comma separated list of config groups.

Options:
  -a --app=<app>
    the application that you wish to attach.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	ptype := safeGetString(args, "<ptype>")
	groups := safeGetString(args, "<groups>")
	return cmdr.ConfigAttach(app, ptype, groups)
}

func configDetach(argv []string, cmdr cmd.Commander) error {
	usage := `
Selects a environment groups to detach an app ptype.

Usage: drycc config:detach <ptype> <groups> [options]

Arguments:
  <ptype>
    the ptype that requires detach configurations.
  <groups>
    comma separated list of config groups.

Options:
  -a --app=<app>
    the application that you wish to detach.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	ptype := safeGetString(args, "<ptype>")
	groups := safeGetString(args, "<groups>")
	return cmdr.ConfigDetach(app, ptype, groups)
}
