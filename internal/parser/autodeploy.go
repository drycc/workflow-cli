package parser

import (
	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

// NewAutodeployCommand creates the autodeploy command
func NewAutodeployCommand(cmdr *commands.DryccCmd) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "autodeploy",
		Short: i18n.T("Manage autodeploy if or not for applications"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmdr.AutodeployInfo(app)
		},
	}

	cmd.PersistentFlags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))
	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)

	cmd.AddCommand(autodeployInfo(cmdr))
	cmd.AddCommand(autodeployEnable(cmdr))
	cmd.AddCommand(autodeployDisable(cmdr))
	return cmd
}

// AutodeployInfo Information about the autodeploy
func autodeployInfo(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: i18n.T("Prints info about the current application's autodeploy if or not"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmdr.AutodeployInfo(app)
		},
	}
	return cmd
}

// AutodeployEnable creates the autodeploy enable command
func autodeployEnable(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable",
		Short: i18n.T("Enables autodeploy for an app"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmdr.AutodeployEnable(app)
		},
	}
	return cmd
}

// AutodeployDisable creates the autodeploy disable command
func autodeployDisable(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable",
		Short: i18n.T("Disables autodeploy for an app"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmdr.AutodeployDisable(app)
		},
	}
	return cmd
}
