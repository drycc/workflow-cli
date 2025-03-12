package parser

import (
	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

// NewAuthCommand creates the auth command
func NewAuthCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: i18n.T("Manage authentication"),
		RunE:  nil,
	}
	cmd.AddCommand(authLogin(cmdr))
	cmd.AddCommand(authLogout(cmdr))
	cmd.AddCommand(authWhoami(cmdr))
	return cmd
}

// AuthLogin creates the auth:login command
func authLogin(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		username  string
		password  string
		sslVerify bool
	}

	cmd := &cobra.Command{
		Use:     "login <controller>",
		Args:    cobra.ExactArgs(1),
		Example: "drycc auth login http://drycc.local3.dryccapp.com/",
		Short:   i18n.T("Authenticate against a controller"),
		Long:    i18n.T("Logs in by authenticating against a controller"),
		RunE: func(cmd *cobra.Command, args []string) error {
			controller := args[0]
			return cmdr.Login(controller, flags.sslVerify, flags.username, flags.password)
		},
	}
	cmd.Flags().StringVarP(&flags.username, "username", "u", "", i18n.T("Provide a username for the account"))
	cmd.Flags().StringVarP(&flags.password, "password", "p", "", i18n.T("Provide a password for the account"))
	cmd.Flags().BoolVar(&flags.sslVerify, "ssl-verify", true, i18n.T("Enables or disables SSL certificate verification for API requests"))
	cmd.Flags().SortFlags = false
	return cmd
}

// AuthLogout creates the auth:logout command
func authLogout(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "logout",
		Example: "drycc auth logout",
		Short:   i18n.T("Clear the current user session"),
		Long:    i18n.T("Logs out from a controller and clears the user session"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmdr.Logout()
		},
	}
	return cmd
}

// AuthWhoami creates the auth:whoami command
func authWhoami(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		all bool
	}

	cmd := &cobra.Command{
		Use:     "whoami",
		Example: "drycc auth whoami",
		Short:   i18n.T("Display the current user"),
		Long:    i18n.T("Displays the currently logged in user"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmdr.Whoami(flags.all)
		},
	}

	cmd.Flags().BoolVarP(&flags.all, "all", "a", false, i18n.T("Fetch detailed user information"))
	return cmd
}
