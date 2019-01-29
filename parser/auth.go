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

auth:register          register a new user
auth:login             authenticate against a controller
auth:logout            clear the current user session
auth:passwd            change the password for the current user
auth:whoami            display the current user
auth:cancel            remove the current user account
auth:regenerate        regenerate user tokens

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "auth:register":
		return authRegister(argv, cmdr)
	case "auth:login":
		return authLogin(argv, cmdr)
	case "auth:logout":
		return authLogout(argv, cmdr)
	case "auth:passwd":
		return authPasswd(argv, cmdr)
	case "auth:whoami":
		return authWhoami(argv, cmdr)
	case "auth:cancel":
		return authCancel(argv, cmdr)
	case "auth:regenerate":
		return authRegenerate(argv, cmdr)
	case "auth":
		fmt.Print(usage)
		return nil
	default:
		PrintUsage(cmdr)
		return nil
	}
}

func authRegister(argv []string, cmdr cmd.Commander) error {
	usage := `
Registers a new user with a Drycc controller.

Usage: drycc auth:register <controller> [options]

Arguments:
  <controller>
    fully-qualified controller URI, e.g. 'http://drycc.local3.dryccapp.com/'

Options:
  --username=<username>
    provide a username for the new account.
  --password=<password>
    provide a password for the new account.
  --email=<email>
    provide an email address.
  --login=true
    logs into the new account after registering.
  --ssl-verify=true
    enables/disables SSL certificate verification for API requests
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	controller := safeGetValue(args, "<controller>")
	username := safeGetValue(args, "--username")
	password := safeGetValue(args, "--password")
	email := safeGetValue(args, "--email")
	sslVerify := true
	login := true

	if args["--ssl-verify"] != nil && args["--ssl-verify"].(string) == "false" {
		sslVerify = false
	}

	// NOTE(bacongobbler): two use cases to check here:
	//
	// 1) Legacy; calling `drycc auth:register` without --login
	// 2) calling `drycc auth:register --login false`
	if args["--login"] != nil && args["--login"].(string) == "false" {
		login = false
	}

	return cmdr.Register(controller, username, password, email, sslVerify, login)
}

func authLogin(argv []string, cmdr cmd.Commander) error {
	usage := `
Logs in by authenticating against a controller.

Usage: drycc auth:login <controller> [options]

Arguments:
  <controller>
    a fully-qualified controller URI, e.g. "http://drycc.local3.dryccapp.com/".

Options:
  --username=<username>
    provide a username for the account.
  --password=<password>
    provide a password for the account.
  --ssl-verify=true
    enables/disables SSL certificate verification for API requests
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	controller := safeGetValue(args, "<controller>")
	username := safeGetValue(args, "--username")
	password := safeGetValue(args, "--password")
	sslVerify := true

	if args["--ssl-verify"] != nil && args["--ssl-verify"].(string) == "false" {
		sslVerify = false
	}

	return cmdr.Login(controller, username, password, sslVerify)
}

func authLogout(argv []string, cmdr cmd.Commander) error {
	usage := `
Logs out from a controller and clears the user session.

Usage: drycc auth:logout

Options:
`

	if _, err := docopt.Parse(usage, argv, true, "", false, true); err != nil {
		return err
	}

	return cmdr.Logout()
}

func authPasswd(argv []string, cmdr cmd.Commander) error {
	usage := `
Changes the password for the current user.

Usage: drycc auth:passwd [options]

Options:
  --password=<password>
    the current password for the account.
  --new-password=<new-password>
    the new password for the account.
  --username=<username>
    the account's username.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	username := safeGetValue(args, "--username")
	password := safeGetValue(args, "--password")
	newPassword := safeGetValue(args, "--new-password")

	return cmdr.Passwd(username, password, newPassword)
}

func authWhoami(argv []string, cmdr cmd.Commander) error {
	usage := `
Displays the currently logged in user.

Usage: drycc auth:whoami [options]

Options:
  --all
    fetch a more detailed description about the user.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmdr.Whoami(args["--all"].(bool))
}

func authCancel(argv []string, cmdr cmd.Commander) error {
	usage := `
Cancels and removes the current account.

Usage: drycc auth:cancel [options]

Options:
  --username=<username>
    provide a username for the account.
  --password=<password>
    provide a password for the account.
  --yes
    force "yes" when prompted.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	username := safeGetValue(args, "--username")
	password := safeGetValue(args, "--password")
	yes := args["--yes"].(bool)

	return cmdr.Cancel(username, password, yes)
}

func authRegenerate(argv []string, cmdr cmd.Commander) error {
	usage := `
Regenerates auth token, defaults to regenerating token for the current user.

Usage: drycc auth:regenerate [options]

Options:
  -u --username=<username>
    specify user to regenerate. Requires admin privileges.
  --all
    regenerate token for every user. Requires admin privileges.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	username := safeGetValue(args, "--username")
	all := args["--all"].(bool)

	return cmdr.Regenerate(username, all)
}
