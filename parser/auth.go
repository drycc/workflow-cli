package parser

import (
	"fmt"

	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// Auth routes auth commands to the specific function.
func Auth(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for auth:

auth:login             authenticate against a controller
auth:logout            clear the current user session
auth:whoami            display the current user

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "auth:login":
		return authLogin(argv, cmdr)
	case "auth:logout":
		return authLogout(argv, cmdr)
	case "auth:whoami":
		return authWhoami(argv, cmdr)
	case "auth":
		fmt.Print(usage)
		return nil
	default:
		PrintUsage(cmdr)
		return nil
	}
}

func authLogin(argv []string, cmdr cmd.Commander) error {
	usage := `
Logs in by authenticating against a controller.

Usage: drycc auth:login <controller> [options]

Arguments:
  <controller>
    a fully-qualified controller URI, e.g. "http://drycc.local3.dryccapp.com/".

Options:
  --ssl-verify=true
    enables/disables SSL certificate verification for API requests
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	controller := safeGetValue(args, "<controller>")
	sslVerify := true

	if args["--ssl-verify"] != nil && args["--ssl-verify"].(string) == "false" {
		sslVerify = false
	}

	return cmdr.Login(controller, sslVerify)
}

func authLogout(argv []string, cmdr cmd.Commander) error {
	usage := `
Logs out from a controller and clears the user session.

Usage: drycc auth:logout

Options:
`

	if _, err := docopt.ParseArgs(usage, argv, ""); err != nil {
		return err
	}

	return cmdr.Logout()
}

func authWhoami(argv []string, cmdr cmd.Commander) error {
	usage := `
Displays the currently logged in user.

Usage: drycc auth:whoami [options]

Options:
  --all
    fetch a more detailed description about the user.
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	return cmdr.Whoami(args["--all"].(bool))
}
