package parser

import (
	"strings"

	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/internal/template"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

func NewPermsCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "perms",
		Short: i18n.T("Manage permissions for applications"),
		RunE: func(_ *cobra.Command, _ []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.PermList(app, results)
		},
	}

	cmd.PersistentFlags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))
	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)

	cmd.AddCommand(permsListCommand(cmdr))
	cmd.AddCommand(permAddCommand(cmdr))
	cmd.AddCommand(permUpdateCommand(cmdr))
	cmd.AddCommand(permRemoveCommand(cmdr))
	return cmd
}

func permsListCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: i18n.T("List all user permissions"),
		RunE: func(_ *cobra.Command, _ []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.PermList(app, results)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))

	return cmd
}

func permAddCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		app         string
		username    string
		permissions string
	}

	userPermsArgsCompletion := completion.UserPermsArgsCompletion{
		UserPermsCompletion: &completion.UserPermsCompletion{ConfigFile: &cmdr.ConfigFile},
	}
	cmd := &cobra.Command{
		Use:  "add <username> <permission>...",
		Args: cobra.MinimumNArgs(2),
		Example: template.CustomExample(
			"drycc add username view change delete",
			map[string]string{
				"<username>":   i18n.T("The name of the user"),
				"<permission>": i18n.T("The user permissions (view,change,delete)"),
			},
		),
		Short:             i18n.T("Add a user permissions"),
		ValidArgsFunction: userPermsArgsCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			flags.username = args[0]
			flags.permissions = strings.Join(args[1:], ",")
			return cmdr.PermCreate(app, flags.username, flags.permissions)
		},
	}

	return cmd
}

func permUpdateCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		app         string
		username    string
		permissions string
	}

	permUpdateCompletion := completion.PermUpdateCompletion{AppID: &app, ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "update <username> <permission>... ",
		Args: cobra.MinimumNArgs(2),
		Example: template.CustomExample(
			"drycc update username view,change",
			map[string]string{
				"<username>":   i18n.T("The name of the user"),
				"<permission>": i18n.T("The user's permissions (view,change,delete)"),
			},
		),
		Short:             i18n.T("Update a user permission"),
		ValidArgsFunction: permUpdateCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			flags.username = args[0]
			flags.permissions = strings.Join(args[1:], ",")
			return cmdr.PermUpdate(app, flags.username, flags.permissions)
		},
	}

	return cmd
}

func permRemoveCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		app      string
		username string
	}
	permUsernameCompletion := completion.PermUsernameCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "remove <username>",
		Args: cobra.MinimumNArgs(1),
		Example: template.CustomExample(
			"drycc perms remove username",
			map[string]string{
				"<username>": i18n.T("The name of the user"),
			},
		),
		Short:             i18n.T("Remove a user permissions"),
		ValidArgsFunction: permUsernameCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {

			flags.username = args[0]
			return cmdr.PermDelete(app, flags.username)
		},
	}

	return cmd
}
