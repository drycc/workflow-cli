package parser

import (
	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

func NewGitCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "git",
		Short: i18n.T("Manage git for applications"),
		RunE:  nil,
	}

	cmd.PersistentFlags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))
	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)

	cmd.AddCommand(gitRemote(cmdr))
	cmd.AddCommand(gitRemove(cmdr))
	return cmd
}

func gitRemote(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		app    string
		remote string
		force  bool
	}

	cmd := &cobra.Command{
		Use:   "remote",
		Short: i18n.T("Adds git remote of application to repository"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.GitRemote(app, flags.remote, flags.force)
		},
	}

	cmd.Flags().StringVarP(&flags.remote, "remote", "r", "drycc", i18n.T("Name of remote to create"))
	cmd.Flags().BoolVarP(&flags.force, "force", "f", false, i18n.T("Overwrite remote of the given name if it already exists"))
	cmd.Flags().SortFlags = false

	return cmd
}

func gitRemove(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove",
		Short: i18n.T("Removes git remote of application from repository"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.GitRemove(app)
		},
	}

	return cmd
}
