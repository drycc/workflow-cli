package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// Limits routes limits commands to their specific function
func Limits(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for limits:

limits:list        list resource limits for an app
limits:set         set resource limits for an app
limits:unset       unset resource limits for an app
limits:specs       list specification information of the server
limits:plans       list resource limit plans

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "limits:list":
		return limitsList(argv, cmdr)
	case "limits:set":
		return limitSet(argv, cmdr)
	case "limits:unset":
		return limitUnset(argv, cmdr)
	case "limits:specs":
		return limitSpecs(argv, cmdr)
	case "limits:plans":
		return limitPlans(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "limits" {
			argv[0] = "limits:list"
			return limitsList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func limitsList(argv []string, cmdr cmd.Commander) error {
	usage := `
Lists resource limits for an application.

Usage: drycc limits:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable name of the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	return cmdr.LimitsList(safeGetString(args, "--app"))
}

func limitSet(argv []string, cmdr cmd.Commander) error {
	usage := `
Sets resource limits for an application.

A resource limit is a finite resource within a pod which we can apply
restrictions through Kubernetes.

Usage: drycc limits:set <ptype>=<value>... [options]

Arguments:
  <ptype>
    the process type as defined in your Procfile, such as 'web' or 'worker'.
  <value>
    The limit plan id to apply to the process type.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.

Use 'drycc help [command]' to learn more.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	limits := args["<ptype>=<value>"].([]string)
	return cmdr.LimitsSet(app, limits)
}

func limitUnset(argv []string, cmdr cmd.Commander) error {
	usage := `
Unsets resource limits for an application.

Usage: drycc limits:unset <ptype>... [options]

Arguments:
  <ptype>
    the process type as defined in your Procfile, such as 'web' or 'worker'.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	limits := args["<ptype>"].([]string)

	return cmdr.LimitsUnset(app, limits)
}

func limitSpecs(argv []string, cmdr cmd.Commander) error {
	usage := `
List all available limit specs.

Usage: drycc limits:specs [options]

Options:
  -l --limit=<num>
    the maximum number of results to display, defaults to config setting.
  -k --keywords=<keywords>
    search keywords separated by commas, matching must satisfy all of the specified.
`
	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}
	results, err := responseLimit(safeGetString(args, "--limit"))

	if err != nil {
		return err
	}
	keywords := safeGetString(args, "--keywords")

	return cmdr.LimitsSpecs(keywords, results)
}

func limitPlans(argv []string, cmdr cmd.Commander) error {
	usage := `
List all available limit plans.

Usage: drycc limits:plans [options]

Options:
  --cpu=<cpu>
    query plans that meet the specified number of cpu cores.
  --memory=<memory>
    query plans that meet the specified memory capacity, unit GiB.
  --spec-id=<spec-id>
    query plans that meet the specified spec id, see [specs] subcommand.
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

	specID := safeGetString(args, "--spec-id")
	cpu := safeGetInt(args, "--cpu")
	memory := safeGetInt(args, "--memory")

	return cmdr.LimitsPlans(specID, cpu, memory, results)
}
