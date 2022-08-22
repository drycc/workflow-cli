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
ps:exec        execute a command in a container
ps:restart     restart an application or its process types
ps:scale       scale processes (e.g. web=4 worker=2)
ps:stop        stop processes (e.g. web worker)
ps:start       start processes (e.g. web worker)

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "ps:list":
		return psList(argv, cmdr)
	case "ps:exec":
		return psExec(argv, cmdr)
	case "ps:restart":
		return psRestart(argv, cmdr)
	case "ps:scale":
		return psScale(argv, cmdr)
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

	args, err := docopt.Parse(usage, argv, true, "", false, true)
	if err != nil {
		return err
	}

	// The 1000 is fake for now until API understands limits
	return cmdr.PsList(safeGetValue(args, "--app"), 1000)
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

	args, err := docopt.Parse(usage, argv, true, "", false, true)
	if err != nil {
		return err
	}
	app := safeGetValue(args, "--app")
	pod := safeGetValue(args, "<pod>")
	tty := args["--tty"].(bool)
	stdin := args["--stdin"].(bool)
	index := slices.Index(os.Args, "--")
	command := os.Args[index+1:]
	return cmdr.PsExec(app, pod, tty, stdin, command)
}

func psRestart(argv []string, cmdr cmd.Commander) error {
	usage := `
Restart an application, a process type or a specific process.

Usage: drycc ps:restart [<type>] [options]

Arguments:
  <type>
    the process name as defined in your Procfile, such as 'web' or 'worker'.
    To restart a particular process, use 'web-asdfg' or 'app-v2-web-asdfg'.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	apps := safeGetValue(args, "--app")
	tp := safeGetValue(args, "<type>")
	return cmdr.PsRestart(apps, tp)
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

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	apps := safeGetValue(args, "--app")
	return cmdr.PsScale(apps, args["<type>=<num>"].([]string))
}
