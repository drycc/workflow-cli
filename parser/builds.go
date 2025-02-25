package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// Builds routes build commands to their specific function.
func Builds(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for builds:

builds:info        print information about a specific build
builds:create      imports an image and deploys as a new release
builds:fetch       fetch the Procfile and dryccfile to the local

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "builds:info":
		return buildsInfo(argv, cmdr)
	case "builds:create":
		return buildsCreate(argv, cmdr)
	case "builds:fetch":
		return buildsFetch(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "builds" {
			argv[0] = "builds:info"
			return buildsInfo(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func buildsInfo(argv []string, cmdr cmd.Commander) error {
	usage := `
Print information about a specific build.

Usage: drycc builds:info [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  -v --version=<version>
    the version for which the build info needs to be displayed.
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

	return cmdr.BuildsInfo(safeGetString(args, "--app"), version)
}

func buildsCreate(argv []string, cmdr cmd.Commander) error {
	usage := `
Creates a new build of an application. Imports an <image> and deploys it to Drycc
as a new release. If a Procfile or drycc.yaml is present in the current directory,
it will be used as the default for this application.

Usage: drycc builds:create <image> [options]

Arguments:
  <image>
    A default fully-qualified container image, either from Drycc Registry (e.g. registry.drycc.cc/drycc/example-go:latest)
    or from an in-house registry (e.g. myregistry.example.com:5000/example-go:latest).
    This image must include the tag.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  -s --stack=<stack>
    the stack name for the application, defaults to container.
  -p --procfile=<procfile>
    a YAML file used to supply a Procfile to the application.
  -d --dryccpath=<dryccpath>
    drycc config path to the application, default is '.drycc'.
  --confirm=yes
    to proceed, type "yes".
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	image := safeGetString(args, "<image>")
	confirm := safeGetString(args, "--confirm")

	stack := safeGetValue(args, "--stack", "container")
	procfile := safeGetValue(args, "--procfile", "Procfile")
	dryccpath := safeGetValue(args, "--dryccpath", ".drycc")

	return cmdr.BuildsCreate(app, image, stack, procfile, dryccpath, confirm)
}

func buildsFetch(argv []string, cmdr cmd.Commander) error {
	usage := `
Print process info about a specific build.

Usage: drycc builds:fetch [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  -v --version=<version>
    the version for which the build info needs to be fetched.
  -p --procfile=<procfile>
    the filename of the procfile saved locally, default is 'Procfile'.
  -d --dryccpath=<dryccpath>
    the folder name of the dryccfile saved locally, default is '.drycc'.
  --confirm=yes
    to proceed, type "yes".
  --save
    save process info to the local.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	var version int
	if safeGetString(args, "--version") != "" {
		if version, err = versionFromString(safeGetString(args, "--version")); err != nil {
			return err
		}
	}

	procfile := safeGetValue(args, "--procfile", "Procfile")
	dryccpath := safeGetValue(args, "--dryccpath", ".drycc")
	save := args["--save"].(bool)
	confirm := safeGetString(args, "--confirm")

	return cmdr.BuildsFetch(app, version, procfile, dryccpath, confirm, save)
}
