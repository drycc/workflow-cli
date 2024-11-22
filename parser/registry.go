package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// Registry routes registry commands to their specific function
func Registry(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for registry:

registry:list        list registry info for an app
registry:set         set registry info for an app
registry:unset       unset registry info for an app

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "registry:list":
		return registryList(argv, cmdr)
	case "registry:set":
		return registrySet(argv, cmdr)
	case "registry:unset":
		return registryUnset(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "registry" {
			argv[0] = "registry:list"
			return registryList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func registryList(argv []string, cmdr cmd.Commander) error {
	usage := `
Lists registry information for an application.

Usage: drycc registry:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable name of the application.
  -p --ptype=<ptype>
    the ptype for registry.
  -v --version=<version>
    the version for which the registry needs to be listed.
`

	args, err := docopt.ParseArgs(usage, argv, "")
	if err != nil {
		return err
	}
	var version int
	if safeGetString(args, "--version") != "" {
		if version, err = versionFromString(safeGetString(args, "--version")); err != nil {
			return err
		}
	}
	ptype := safeGetValue(args, "--ptype", "")
	return cmdr.RegistryList(safeGetString(args, "--app"), ptype, version)
}

func registrySet(argv []string, cmdr cmd.Commander) error {
	usage := `
Sets registry information for an application. These credentials are the same as those used for
'podmain login' to the private registry.

Usage: drycc registry:set <username> <password> [options]

Arguments:
  <username>
    the username of the registry.
  <password>
    the password of the registry.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  -p --ptype=<ptype>
    the ptype for registry, default[web].
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	u := safeGetString(args, "<username>")
	p := safeGetString(args, "<password>")
	ptype := safeGetValue(args, "--ptype", "web")
	return cmdr.RegistrySet(app, ptype, u, p)
}

func registryUnset(argv []string, cmdr cmd.Commander) error {
	usage := `
Unsets registry information for an application.

Usage: drycc registry:unset [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  -p --ptype=<ptype>
    the ptype for registry, default[web].
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	ptype := safeGetValue(args, "--ptype", "web")
	return cmdr.RegistryUnset(app, ptype)
}
