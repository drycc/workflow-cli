package parser

import (
	"fmt"
	"strconv"

	"github.com/drycc/workflow-cli/cmd"
	docopt "github.com/docopt/docopt-go"
)

// Releases routes releases commands to their specific function.
func Releases(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for releases:

releases:list        list an application's release history
releases:info        print information about a specific release
releases:rollback    return to a previous release

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "releases:list":
		return releasesList(argv, cmdr)
	case "releases:info":
		return releasesInfo(argv, cmdr)
	case "releases:rollback":
		return releasesRollback(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "releases" {
			argv[0] = "releases:list"
			return releasesList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func releasesList(argv []string, cmdr cmd.Commander) error {
	usage := `
Lists release history for an application.

Usage: drycc releases:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  -l --limit=<num>
    the maximum number of results to display, defaults to config setting
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)
	if err != nil {
		return err
	}

	results, err := responseLimit(safeGetValue(args, "--limit"))
	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")

	return cmdr.ReleasesList(app, results)
}

func releasesInfo(argv []string, cmdr cmd.Commander) error {
	usage := `
Prints info about a particular release.

Usage: drycc releases:info <version> [options]

Arguments:
  <version>
    the release of the application, such as 'v1'.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)
	if err != nil {
		return err
	}

	version, err := versionFromString(args["<version>"].(string))

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")

	return cmdr.ReleasesInfo(app, version)
}

func releasesRollback(argv []string, cmdr cmd.Commander) error {
	usage := `
Rolls back to a previous application release.

Usage: drycc releases:rollback [<version>] [options]

Arguments:
  <version>
    the release of the application, such as 'v1'.

Options:
  -a --app=<app>
    the uniquely identifiable name of the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)
	if err != nil {
		return err
	}

	var version int

	if args["<version>"] == nil {
		version = -1
	} else {
		version, err = versionFromString(args["<version>"].(string))

		if err != nil {
			return err
		}
	}

	app := safeGetValue(args, "--app")

	return cmdr.ReleasesRollback(app, version)
}

func versionFromString(version string) (int, error) {
	if version[:1] == "v" {
		if len(version) < 2 {
			return -1, fmt.Errorf("%s is not in the form 'v#'", version)
		}

		return strconv.Atoi(version[1:])
	}

	return strconv.Atoi(version)
}
