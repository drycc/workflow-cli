package parser

import (
  "github.com/notmaxx/workflow-cli/cmd"
  docopt "github.com/docopt/docopt-go"
)

// Timeouts routes timeouts commands to their specific function
func Timeouts(argv []string, cmdr cmd.Commander) error {
  usage := `
Valid commands for timeouts:

timeouts:list        list resource timeouts for an app
timeouts:set         set resource timeouts for an app
timeouts:unset       unset resource timeouts for an app

Use 'deis help [command]' to learn more.
`

  switch argv[0] {
  case "timeouts:list":
    return timeoutList(argv, cmdr)
  case "timeouts:set":
    return timeoutSet(argv, cmdr)
  case "timeouts:unset":
    return timeoutUnset(argv, cmdr)
  default:
    if printHelp(argv, usage) {
      return nil
    }

    if argv[0] == "timeouts" {
      argv[0] = "timeouts:list"
      return timeoutList(argv, cmdr)
    }

    PrintUsage(cmdr)
    return nil
  }
}

func timeoutList(argv []string, cmdr cmd.Commander) error {
  usage := `
Lists resource timeouts for an application.

Usage: deis timeouts:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable name of the application.
`

  args, err := docopt.Parse(usage, argv, true, "", false, true)

  if err != nil {
    return err
  }

  return cmdr.TimeoutsList(safeGetValue(args, "--app"))
}

func timeoutSet(argv []string, cmdr cmd.Commander) error {
  usage := `
Sets termination grace period for an application.

Usage: deis timeouts:set [options] <type>=<value>...

Arguments:
  <type>
    the process type as defined in your Procfile, such as 'web' or 'worker'.
    Note that Dockerfile apps have a default 'cmd' process type.
  <value>
    The value to apply to the process type in seconds.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

  args, err := docopt.Parse(usage, argv, true, "", false, true)

  if err != nil {
    return err
  }

  app := safeGetValue(args, "--app")
  timeouts := args["<type>=<value>"].([]string)

  return cmdr.TimeoutsSet(app, timeouts)
}

func timeoutUnset(argv []string, cmdr cmd.Commander) error {
  usage := `
Unsets timeouts for an application. Default value (30s)
or KUBERNETES_POD_TERMINATION_GRACE_PERIOD_SECONDS is used

Usage: deis timeouts:unset [options] <type>...

Arguments:
  <type>
    the process type as defined in your Procfile, such as 'web' or 'worker'.
    Note that Dockerfile apps have a default 'cmd' process type.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

  args, err := docopt.Parse(usage, argv, true, "", false, true)

  if err != nil {
    return err
  }

  app := safeGetValue(args, "--app")
  timeouts := args["<type>"].([]string)

  return cmdr.TimeoutsUnset(app, timeouts)
}
