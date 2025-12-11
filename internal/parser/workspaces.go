package parser

import (
	"fmt"

	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/internal/template"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

// NewWorkspacesCommand creates the workspaces command.
func NewWorkspacesCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workspaces",
		Short: i18n.T("Manage workspaces"),
		RunE: func(_ *cobra.Command, _ []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.WorkspacesList(results)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))

	cmd.AddCommand(workspacesList(cmdr))
	cmd.AddCommand(workspacesCreate(cmdr))
	cmd.AddCommand(workspacesInfo(cmdr))
	cmd.AddCommand(workspacesDelete(cmdr))
	cmd.AddCommand(workspacesInvite(cmdr))
	cmd.AddCommand(workspacesRemove(cmdr))
	cmd.AddCommand(workspacesUpdate(cmdr))
	return cmd
}

func workspacesList(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: i18n.T("Lists workspaces for current user"),
		RunE: func(_ *cobra.Command, _ []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.WorkspacesList(results)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))
	return cmd
}

func workspacesCreate(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		email string
	}

	cmd := &cobra.Command{
		Use:  "create <name>",
		Args: cobra.ExactArgs(1),
		Example: template.CustomExample(
			"drycc workspaces create my-workspace --email team@example.com",
			map[string]string{
				"<name>": i18n.T("The name of the workspace"),
			},
		),
		Short: i18n.T("Create a workspace"),
		RunE: func(_ *cobra.Command, args []string) error {
			return cmdr.WorkspacesCreate(args[0], flags.email)
		},
	}

	cmd.Flags().StringVar(&flags.email, "email", "", i18n.T("The contact email for the workspace"))
	return cmd
}

func workspacesInfo(cmdr *commands.DryccCmd) *cobra.Command {
	workspaceCompletion := completion.WorkspaceCompletion{ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "info <name>",
		Args: cobra.ExactArgs(1),
		Example: template.CustomExample(
			"drycc workspaces info my-workspace",
			map[string]string{
				"<name>": i18n.T("The name of the workspace"),
			},
		),
		Short:             i18n.T("Print information about a workspace"),
		ValidArgsFunction: workspaceCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.WorkspacesInfo(args[0], results)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))
	return cmd
}

func workspacesDelete(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		confirm string
	}

	workspaceCompletion := completion.WorkspaceCompletion{ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "delete <name>",
		Args: cobra.ExactArgs(1),
		Example: template.CustomExample(
			"drycc workspaces delete my-workspace",
			map[string]string{
				"<name>": i18n.T("The name of the workspace"),
			},
		),
		Short:             i18n.T("Delete a workspace"),
		ValidArgsFunction: workspaceCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			return cmdr.WorkspacesDelete(args[0], flags.confirm)
		},
	}

	cmd.Flags().StringVar(&flags.confirm, "confirm", "", i18n.T("Skips the prompt for the workspace name"))
	return cmd
}

func workspacesInvite(cmdr *commands.DryccCmd) *cobra.Command {
	workspaceCompletion := completion.WorkspaceCompletion{ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "invite <workspace> <email>",
		Args: cobra.ExactArgs(2),
		Example: template.CustomExample(
			"drycc workspaces invite my-workspace user@example.com",
			map[string]string{
				"<workspace>": i18n.T("The name of the workspace"),
				"<email>":     i18n.T("The email address of the user to invite"),
			},
		),
		Short:             i18n.T("Invite a user to a workspace"),
		ValidArgsFunction: workspaceCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			return cmdr.WorkspacesInvite(args[0], args[1])
		},
	}

	return cmd
}

func workspacesRemove(cmdr *commands.DryccCmd) *cobra.Command {
	removeCompletion := completion.WorkspaceRemoveCompletion{ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "remove <workspace> <username>",
		Args: cobra.ExactArgs(2),
		Example: template.CustomExample(
			"drycc workspaces remove my-workspace john",
			map[string]string{
				"<workspace>": i18n.T("The name of the workspace"),
				"<username>":  i18n.T("The username to remove from the workspace"),
			},
		),
		Short:             i18n.T("Remove a user from a workspace"),
		ValidArgsFunction: removeCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			return cmdr.WorkspacesRemove(args[0], args[1])
		},
	}

	return cmd
}

func workspacesUpdate(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		role   string
		alerts string
	}

	updateCompletion := completion.WorkspaceUpdateCompletion{ConfigFile: &cmdr.ConfigFile}
	roleCompletion := completion.WorkspaceRoleCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "update <workspace> <username>",
		Args: cobra.ExactArgs(2),
		Example: template.CustomExample(
			"drycc workspaces update my-workspace john --role=admin --alerts=true",
			map[string]string{
				"<workspace>": i18n.T("The name of the workspace"),
				"<username>":  i18n.T("The username of the member to update"),
			},
		),
		Short:             i18n.T("Update a member's role and/or alerts in a workspace"),
		ValidArgsFunction: updateCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			var alerts *bool
			if flags.alerts != "" {
				v := flags.alerts == "true"
				alerts = &v
			}
			if flags.role == "" && alerts == nil {
				return fmt.Errorf("at least one of --role or --alerts must be specified")
			}
			return cmdr.WorkspacesUpdate(args[0], args[1], flags.role, alerts)
		},
	}

	cmd.Flags().StringVar(&flags.role, "role", "", i18n.T("The role to assign: admin, member, or viewer"))
	cmd.Flags().StringVar(&flags.alerts, "alerts", "", i18n.T("Enable or disable alerts: true or false"))
	cmd.RegisterFlagCompletionFunc("role", roleCompletion.CompletionFunc)
	return cmd
}
