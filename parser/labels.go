package parser

import (
	"github.com/deis/workflow-cli/cmd"
	docopt "github.com/docopt/docopt-go"
)

// Labels displays all relevant commands for `deis label`.
func Labels(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for labels:

labels:list   list application's labels
labels:set    add new application's label
labels:unset  remove application's label

Use 'deis help [command]' to learn more.
`

	switch argv[0] {
	case "labels:list":
		return labelsList(argv, cmdr)
	case "labels:set":
		return labelsSet(argv, cmdr)
	case "labels:unset":
		return labelsUnset(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "labels" {
			argv[0] = "labels:list"
			return labelsList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func labelsList(argv []string, cmdr cmd.Commander) error {
	usage := `
Prints a list of labels of the application.

Usage: deis labels:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmdr.LabelsList(safeGetValue(args, "--app"))
}

func labelsSet(argv []string, cmdr cmd.Commander) error {
	usage := `
Sets labels for an application.

A label is a key/value pair used to label an application. This label is a general information for deis user.
Mostly used for administration/maintenance information, note for application. This information isn't send to scheduler.

Usage: deis labels:set [options] <key>=<value>...

Arguments:
  <key> the label key, for example: "git_repo" or "team"
  <value> the label value, for example: "https://github.com/teamhephy/workflow" or "frontend"

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

	return cmdr.LabelsSet(app, tags)
}

func labelsUnset(argv []string, cmdr cmd.Commander) error {
	usage := `
Unsets labels for an application.

Usage: deis labels:unset [options] <key>...

Arguments:
  <key> the label key to unset, for example: "git_repo" or "team"

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

	return cmdr.LabelsUnset(app, tags)
}
