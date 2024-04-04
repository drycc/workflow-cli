package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// Builds routes build commands to their specific function.
func Builds(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for builds:

builds:list        list build history for an application
builds:create      imports an image and deploys as a new release

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "builds:list":
		return buildsList(argv, cmdr)
	case "builds:create":
		return buildsCreate(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "builds" {
			argv[0] = "builds:list"
			return buildsList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func buildsList(argv []string, cmdr cmd.Commander) error {
	usage := `
Lists build history for an application.

Usage: drycc builds:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
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

	return cmdr.BuildsList(safeGetString(args, "--app"), results)
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
    The uniquely identifiable name for the application.
  -s --stack=<stack>
    The stack name for the application, defaults to container.
  -p --procfile=<procfile>
    A YAML file used to supply a Procfile to the application.
  -d --dryccfile=<dryccfile>
    A YAML file used to supply a drycc.yaml to the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	image := safeGetString(args, "<image>")
	stack := safeGetString(args, "--stack")
	if stack == "" {
		stack = "container"
	}
	procfile := safeGetValue(args, "--procfile", "Procfile")
	dryccfile := safeGetValue(args, "--dryccfile", "drycc.yaml")

	return cmdr.BuildsCreate(app, image, stack, procfile, dryccfile)
}
