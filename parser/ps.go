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
ps:describe    print a detailed description of the selected process
ps:delete      delete the selected processes

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "ps:list":
		return psList(argv, cmdr)
	case "ps:logs":
		return psLogs(argv, cmdr)
	case "ps:exec":
		return psExec(argv, cmdr)
	case "ps:describe":
		return psDescribe(argv, cmdr)
	case "ps:delete":
		return psDelete(argv, cmdr)
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

Arguments:
  <pod> the pod name for the application.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  -n --lines=<lines>
    the number of lines to display, default to 300 lines, -1 showing all log lines.
  -f --follow
    specify if the logs should be streamed.
  -c --container=<container>
    print the logs of this container.
  --previous
    print the logs for the previous instance of the container in a pod if it exists.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	lines := safeGetInt(args, "--lines")
	if lines < 0 {
		lines = -1
	} else if lines == 0 {
		lines = 300
	}
	follow := safeGetBool(args, "--follow")
	podID := safeGetString(args, "<pod>")
	container := safeGetString(args, "--container")
	previous := args["--previous"].(bool)

	return cmdr.PsLogs(app, podID, lines, follow, container, previous)
}

func psExec(argv []string, cmdr cmd.Commander) error {
	usage := `
Execute a command in a container.

Usage: drycc ps:exec <pod> [options] -- <command>...

Arguments:
  <pod> the pod name for the application.

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

func psDescribe(argv []string, cmdr cmd.Commander) error {
	usage := `
Print a detailed description of the selected process.

Usage: drycc ps:describe <pod> [options]

Arguments:
  <pod> the pod name for the application.

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

func psDelete(argv []string, cmdr cmd.Commander) error {
	usage := `
Delete the selected processes.

Usage: drycc ps:delete <pod>... [options]

Arguments:
  <pod> the pod name for the application.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")
	if err != nil {
		return err
	}
	app := safeGetString(args, "--app")
	return cmdr.PsDelete(app, args["<pod>"].([]string))
}
