package parser

import (
	"os"

	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
	"golang.org/x/exp/slices"
)

// Ps routes ps commands to their specific function.
func Ps(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for processes:

ps:list        list application processes
ps:logs        print the logs for a container
ps:exec        execute a command in a container
ps:restart     restart an application or process type
ps:scale       scale processes (e.g. web=4 worker=2)
ps:describe    print a detailed description of the selected process

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "ps:list":
		return psList(argv, cmdr)
	case "ps:logs":
		return psLogs(argv, cmdr)
	case "ps:exec":
		return psExec(argv, cmdr)
	case "ps:restart":
		return psRestart(argv, cmdr)
	case "ps:scale":
		return psScale(argv, cmdr)
	case "ps:describe":
		return psDescribe(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "ps" {
			argv[0] = "ps:list"
			return psList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func psList(argv []string, cmdr cmd.Commander) error {
	usage := `
Lists processes servicing an application.

Usage: drycc ps:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")
	if err != nil {
		return err
	}

	// The 1000 is fake for now until API understands limits
	return cmdr.PsList(safeGetString(args, "--app"), 1000)
}

func psLogs(argv []string, cmdr cmd.Commander) error {
	usage := `
Print the logs for a container in a pod or specified resource.

Usage: drycc ps:logs <pod> [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  -n --lines=<lines>
    the number of lines to display.
  -f --follow
    specify if the logs should be streamed.
  -t --container=<container>
    print the logs of this container.
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
	podID := safeGetString(args, "<pod>")
	container := safeGetString(args, "--container")

	return cmdr.PsLogs(app, podID, lines, follow, container)
}

func psExec(argv []string, cmdr cmd.Commander) error {
	usage := `
Execute a command in a container.

Usage: drycc ps:exec <pod> [options] -- <command>...

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  -t --tty
    stdin is a TTY.
  -i --stdin
    pass stdin to the container.
`

	args, err := docopt.ParseArgs(usage, argv, "")
	if err != nil {
		return err
	}
	app := safeGetString(args, "--app")
	pod := safeGetString(args, "<pod>")
	tty := args["--tty"].(bool)
	stdin := args["--stdin"].(bool)
	index := slices.Index(os.Args, "--")
	command := os.Args[index+1:]
	return cmdr.PsExec(app, pod, tty, stdin, command)
}

func psRestart(argv []string, cmdr cmd.Commander) error {
	usage := `
Restart an application or process type.

Usage: drycc ps:restart [<type>...] [options]

Arguments:
  <type>
    the process name as defined in your Procfile, such as 'web' or 'web worker'.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  --confirm=yes
    To proceed, type "yes".
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	apps := safeGetString(args, "--app")
	confirm := safeGetString(args, "--confirm")
	return cmdr.PsRestart(apps, args["<type>"].([]string), confirm)
}

func psScale(argv []string, cmdr cmd.Commander) error {
	usage := `
Scales an application's processes by type.

Usage: drycc ps:scale <type>=<num>... [options]

Arguments:
  <type>
    the process name as defined in your Procfile, such as 'web' or 'worker'.
    Note that Dockerfile apps have a default 'cmd' process type.
  <num>
    the number of processes.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	apps := safeGetString(args, "--app")
	return cmdr.PsScale(apps, args["<type>=<num>"].([]string))
}

func psDescribe(argv []string, cmdr cmd.Commander) error {
	usage := `
Print a detailed description of the selected process.

Usage: drycc ps:describe <pod> [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	pod := safeGetString(args, "<pod>")
	return cmdr.PsDescribe(app, pod)
}
