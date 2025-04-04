package parser

import (
	"strings"

	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// Apps routes app commands to their specific function.
func Apps(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for apps:

apps:create        create a new application
apps:list          list accessible applications
apps:info          view info about an application
apps:open          open the application in a browser
apps:logs          view aggregated application logs
apps:run           run a command in an ephemeral app container
apps:destroy       destroy an application
apps:transfer      transfer app ownership to another user

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "apps:create":
		return appCreate(argv, cmdr)
	case "apps:list":
		return appsList(argv, cmdr)
	case "apps:info":
		return appInfo(argv, cmdr)
	case "apps:open":
		return appOpen(argv, cmdr)
	case "apps:logs":
		return appLogs(argv, cmdr)
	case "apps:run":
		return appRun(argv, cmdr)
	case "apps:destroy":
		return appDestroy(argv, cmdr)
	case "apps:transfer":
		return appTransfer(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "apps" {
			argv[0] = "apps:list"
			return appsList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func appCreate(argv []string, cmdr cmd.Commander) error {
	usage := `
Creates a new application.

- if no <id> is provided, one will be generated automatically.

Usage: drycc apps:create [<id>] [options]

Arguments:
  <id>
    a uniquely identifiable name for the application. No other app can already
    exist with this name.

Options:
  --no-remote
    do not create a 'drycc' git remote.
  -r --remote REMOTE
    name of remote to create. [default: drycc].
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	id := safeGetString(args, "<id>")
	remote := safeGetString(args, "--remote")
	noRemote := safeGetBool(args, "--no-remote")

	return cmdr.AppCreate(id, remote, noRemote)
}

func appsList(argv []string, cmdr cmd.Commander) error {
	usage := `
Lists applications visible to the current user.

Usage: drycc apps:list [options]

Options:
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

	return cmdr.AppsList(results)
}

func appInfo(argv []string, cmdr cmd.Commander) error {
	usage := `
Prints info about the current application.

Usage: drycc apps:info [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")

	return cmdr.AppInfo(app)
}

func appOpen(argv []string, cmdr cmd.Commander) error {
	usage := `
Opens a URL to the application in the default browser.

Usage: drycc apps:open [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")

	return cmdr.AppOpen(app)
}

func appLogs(argv []string, cmdr cmd.Commander) error {
	usage := `
Retrieves the most recent log events.

Usage: drycc apps:logs [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  -n --lines=<lines>
    the number of lines to display.
  -f --follow
    specify if the logs should be streamed.
  -t --timeout=<timeout>
    the max seconds of follow the log stream.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	lines := safeGetInt(args, "--lines")
	if lines <= 0 {
		lines = 300
	}
	follow := safeGetBool(args, "--follow")
	timeout := safeGetInt(args, "--timeout")
	if timeout <= 0 {
		timeout = 300
	}

	return cmdr.AppLogs(app, lines, follow, timeout)
}

func appRun(argv []string, cmdr cmd.Commander) error {
	usage := `
Runs a command inside an ephemeral app container.

Usage: drycc apps:run [--mount=<volume>:<path>...] [options] [--] <command>...

Arguments:
  <volume>
    the volume name.
  <path>
    the filesystem path.
  <command>
    the shell command to run inside the container.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  --timeout=<timeout>
    the timeout for command run, default to 3600 seconds.
  --expires=<expires>
    retention time of running records, default to 3600 seconds.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	timeout := uint32(safeGetInt(args, "--timeout"))
	expires := uint32(safeGetInt(args, "--expires"))
	command := strings.Join(args["<command>"].([]string), " ")
	mounts := args["--mount"].([]string)
	return cmdr.AppRun(app, command, mounts, timeout, expires)
}

func appDestroy(argv []string, cmdr cmd.Commander) error {
	usage := `
Destroys an application.

Usage: drycc apps:destroy [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  --confirm=<app>
    skips the prompt for the application name. <app> is the uniquely identifiable
    name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	confirm := safeGetString(args, "--confirm")

	return cmdr.AppDestroy(app, confirm)
}

func appTransfer(argv []string, cmdr cmd.Commander) error {
	usage := `
Transfer app ownership to another user.

Usage: drycc apps:transfer <username> [options]

Arguments:
  <username>
    the user that the app will be transferred to.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	user := safeGetString(args, "<username>")

	return cmdr.AppTransfer(app, user)
}
