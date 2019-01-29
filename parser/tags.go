package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// Tags routes tags commands to their specific function
func Tags(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for tags:

tags:list        list tags for an app
tags:set         set tags for an app
tags:unset       unset tags for an app

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "tags:list":
		return tagsList(argv, cmdr)
	case "tags:set":
		return tagsSet(argv, cmdr)
	case "tags:unset":
		return tagsUnset(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "tags" {
			argv[0] = "tags:list"
			return tagsList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func tagsList(argv []string, cmdr cmd.Commander) error {
	usage := `
Lists tags for an application.

Usage: drycc tags:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable name of the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmdr.TagsList(safeGetValue(args, "--app"))
}

func tagsSet(argv []string, cmdr cmd.Commander) error {
	usage := `
Sets tags for an application.

A tag is a key/value pair used to tag an application's containers and is passed to the
scheduler. This is often used to restrict workloads to specific hosts matching the
scheduler-configured metadata.

Usage: drycc tags:set [options] <key>=<value>...

Arguments:
  <key> the tag key, for example: "environ" or "rack"
  <value> the tag value, for example: "prod" or "1"

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)
	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	tags := args["<key>=<value>"].([]string)

	return cmdr.TagsSet(app, tags)
}

func tagsUnset(argv []string, cmdr cmd.Commander) error {
	usage := `
Unsets tags for an application.

Usage: drycc tags:unset [options] <key>...

Arguments:
  <key> the tag key to unset, for example: "environ" or "rack"

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)
	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	tags := args["<key>"].([]string)

	return cmdr.TagsUnset(app, tags)
}
