package parser

import (
	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/template"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

// NewUsersCommand creates the users command
func NewUsersCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "users",
		Short: i18n.T("Manage registered users"),
		RunE: func(_ *cobra.Command, _ []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.UsersList(results)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))

	cmd.AddCommand(usersListCommand(cmdr))
	cmd.AddCommand(usersEnableCommand(cmdr))
	cmd.AddCommand(usersDisableCommand(cmdr))
	return cmd
}

func usersListCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: i18n.T("List all registered users"),
		RunE: func(_ *cobra.Command, _ []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.UsersList(results)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))
	return cmd
}

func usersEnableCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		username string
	}

	cmd := &cobra.Command{
		Use: "enable <username>",
		Example: template.CustomExample(
			"drycc users enable john_doe",
			map[string]string{
				"<username>": i18n.T("The username you want to enable"),
			},
		),
		Short: i18n.T("Enable a user"),
		Long: i18n.T(`Enable a user when his status is disabled.
Requires admin privileges.`),
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			flags.username = args[0]
			return cmdr.UsersEnable(flags.username)
		},
	}
	return cmd
}

func usersDisableCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		username string
	}

	cmd := &cobra.Command{
		Use: "disable <username>",
		Example: template.CustomExample(
			"drycc users disable john_doe",
			map[string]string{
				"<username>": i18n.T("The username you want to disable"),
			},
		),
		Short: i18n.T("Disable a user"),
		Long: i18n.T(`Disable a user when his status is disabled.
Requires admin privileges.`),
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			flags.username = args[0]
			return cmdr.UsersDisable(flags.username)
		},
	}
	return cmd
}
