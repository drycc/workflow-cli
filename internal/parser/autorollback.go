package parser

import (
	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

// NewAutorollbackCommand creates the autorollback command
func NewAutorollbackCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "autorollback",
		Short: i18n.T("Manage autorollback if or not for application"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmdr.AutorollbackInfo(app)
		},
	}
	cmd.PersistentFlags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))
	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)

	cmd.AddCommand(autorollbackInfo(cmdr))
	cmd.AddCommand(autorollbackEnable(cmdr))
	cmd.AddCommand(autorollbackDisable(cmdr))
	return cmd
}

// AutorollbackInfo creates the autorollback info command
func autorollbackInfo(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: i18n.T("View autorollback info for an application"),
		Long:  i18n.T("Prints info about the current application's autorollback if or not"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmdr.AutorollbackInfo(app)
		},
	}
	return cmd
}

// AutorollbackEnable creates the autorollback enable command
func autorollbackEnable(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable",
		Short: i18n.T("Enable autorollback for an application"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmdr.AutorollbackEnable(app)
		},
	}
	return cmd
}

// AutorollbackDisable creates the autorollback disable command
func autorollbackDisable(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable",
		Short: i18n.T("Disable autorollback for an application"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmdr.AutorollbackDisable(app)
		},
	}
	return cmd
}
