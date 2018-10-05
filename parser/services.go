package parser

import (
  "github.com/teamhephy/workflow-cli/cmd"
  docopt "github.com/docopt/docopt-go"
)

// Services routes service commands to their specific function.
func Services(argv []string, cmdr cmd.Commander) error {
  usage := `
Valid commands for services:

services:add           create service for an application
services:list          list application services
services:remove        remove service from an application

Use 'deis help [command]' to learn more.
`

  switch argv[0] {
  case "services:add":
    return servicesAdd(argv, cmdr)
  case "services:list":
    return servicesList(argv, cmdr)
  case "services:remove":
    return servicesRemove(argv, cmdr)
  default:
    if printHelp(argv, usage) {
      return nil
    }

    if argv[0] == "services" {
      argv[0] = "services:list"
      return servicesList(argv, cmdr)
    }

    PrintUsage(cmdr)
    return nil
  }
}

func servicesAdd(argv []string, cmdr cmd.Commander) error {
  usage := `
Creates extra service for an application and binds it to specific route of the main app domain

Usage: deis services:add --type <procfile_type> --route <path_pattern> [options]

Arguments:
  <procfile_type>
    Procfile type which should handle the request, e.g. webhooks (should be bind to the port PORT).
    Only single extra service per Porcfile type could be created

  <path_pattern>
    Nginx locations where route requests, one or many via comma,
    e.g. /webhooks/notify
    OR "/webhooks/notify,~ ^/users/[0-9]+/.*/webhooks/notify,/webhooks/rest"

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

  args, err := docopt.Parse(usage, argv, true, "", false, true)

  if err != nil {
    return err
  }

  app := safeGetValue(args, "--app")
  procfileType := safeGetValue(args, "<procfile_type>")
  pathPattern := safeGetValue(args, "<path_pattern>")

  return cmdr.ServicesAdd(app, procfileType, pathPattern)
}

func servicesList(argv []string, cmdr cmd.Commander) error {
  usage := `
Lists extra services for an application

Usage: deis services:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

  args, err := docopt.Parse(usage, argv, true, "", false, true)

  if err != nil {
    return err
  }

  app := safeGetValue(args, "--app")

  return cmdr.ServicesList(app)
}

func servicesRemove(argv []string, cmdr cmd.Commander) error {
  usage := `
Deletes specific extra service for application

Usage: deis services:remove <procfile_type> [options]

Arguments:
  <procfile_type>
    extra service for procfile type that should be removed

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

  args, err := docopt.Parse(usage, argv, true, "", false, true)

  if err != nil {
    return err
  }

  app := safeGetValue(args, "--app")
  procfileType := safeGetValue(args, "<procfile_type>")

  return cmdr.ServicesRemove(app, procfileType)
}
