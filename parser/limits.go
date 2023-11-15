package parser

import (
	"regexp"

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

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "limits:list":
		return limitsList(argv, cmdr)
	case "limits:set":
		return limitSet(argv, cmdr)
	case "limits:unset":
		return limitUnset(argv, cmdr)
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

	return cmdr.LimitsList(safeGetValue(args, "--app"))
}

func limitSet(argv []string, cmdr cmd.Commander) error {
	usage := `
Sets resource limits for an application.

A resource limit is a finite resource within a pod which we can apply
restrictions through Kubernetes.The limit is applied to each individual pod,
so setting a memory limit of 1G for an application means that each pod gets 1G of memory.

Usage: drycc limits:set [options] <type>=<value>...

Arguments:
  <type>
    the process type as defined in your Procfile, such as 'web' or 'worker'.
    Note that Dockerfile apps have a default 'cmd' process type.
  <value>
    The value to apply to the process type. By default, this is set to --memory.
    Can be in <limit> format eg. web=2G db=1G
    You can only set one type of limit per call.

    With --memory, units are represented in Megabytes(M), or Gigabytes (G).
	For example, 'drycc limit:set cmd=1G' will restrict all
    "cmd" processes to a maximum of 1 Gigabyte of memory each.

    With --cpu, units are represented in the number of CPUs. For example,
    'drycc limit:set --cpu cmd=1' will restrict all "cmd" processes to a
    maximum of 1 CPU. Alternatively, you can also use milli units to specify the
    number of CPU shares the pod can use. For example, 'drycc limits:set --cpu cmd=500m'
    will restrict all "cmd" processes to half of a CPU.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  --cpu
    value apply to CPU.
  -m --memory
    value apply to memory.

Use 'drycc help [command]' to learn more.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	cpuLimits := []string{}
	memoryLimits := []string{}
	for _, value := range args["<type>=<value>"].([]string) {
		if args["--cpu"].(bool) {
			isCPU, _ := regexp.MatchString("\\d+m?$", value)
			if isCPU {
				cpuLimits = append(cpuLimits, value)
			}
		}
		isMemory, _ := regexp.MatchString("\\d+[M|G]$", value)
		if isMemory {
			memoryLimits = append(memoryLimits, value)
		}
	}

	return cmdr.LimitsSet(app, cpuLimits, memoryLimits)
}

func limitUnset(argv []string, cmdr cmd.Commander) error {
	usage := `
Unsets resource limits for an application.

Usage: drycc limits:unset [options] [--memory | --cpu] <type>...

Arguments:
  <type>
    the process type as defined in your Procfile, such as 'web' or 'worker'.
    Note that Dockerfile apps have a default 'cmd' process type.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  --cpu
    limits cpu shares.
  -m --memory
    limits memory. [default: true]
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	cpuLimits := []string{}
	memoryLimits := []string{}

	if args["--cpu"].(bool) {
		cpuLimits = args["<type>"].([]string)
	}

	if args["--memory"].(bool) {
		memoryLimits = args["<type>"].([]string)
	}

	return cmdr.LimitsUnset(app, cpuLimits, memoryLimits)
}
