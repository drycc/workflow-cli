package parser

import (
	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/internal/template"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

// NewTimeoutsCommand creates the timeouts command
func NewTimeoutsCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "timeouts",
		Short: i18n.T("Manage pods termination grace period"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.TimeoutsList(app, version)
		},
	}

	cmd.PersistentFlags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))
	cmd.Flags().IntVarP(&version, "version", "v", 0, i18n.T("The version for which the timeout needs to be listed"))

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)
	releaseCompletion := completion.ReleaseCompletion{ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("version", releaseCompletion.CompletionFunc)

	cmd.AddCommand(timeoutListCommand(cmdr))
	cmd.AddCommand(timeoutSetCommand(cmdr))
	cmd.AddCommand(timeoutUnsetCommand(cmdr))
	return cmd
}

func timeoutListCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: i18n.T("List resource timeouts for an app"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.TimeoutsList(app, version)
		},
	}

	cmd.Flags().IntVarP(&version, "version", "v", 0, i18n.T("The version for which the timeout needs to be listed"))

	releaseCompletion := completion.ReleaseCompletion{ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("version", releaseCompletion.CompletionFunc)
	return cmd
}

func timeoutSetCommand(cmdr *commands.DryccCmd) *cobra.Command {
	ptsSetArgsCompletion := completion.PtsSetArgsCompletion{
		PtsCompletion: &completion.PtsCompletion{AppID: &app, ArgsLen: -1, ConfigFile: &cmdr.ConfigFile},
	}
	cmd := &cobra.Command{
		Use: "set <ptype>=<value>...",
		Example: template.CustomExample(
			"drycc timeouts set web=30 worker=60",
			map[string]string{
				"<ptype>": i18n.T("The process type as defined in your Procfile"),
				"<value>": i18n.T("The value to apply to the process type in seconds"),
			},
		),
		Short:             i18n.T("Set resource timeouts for an app"),
		Long:              i18n.T("Sets termination grace period for an application"),
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: ptsSetArgsCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			timeouts := args
			return cmdr.TimeoutsSet(app, timeouts)
		},
	}

	return cmd
}

func timeoutUnsetCommand(cmdr *commands.DryccCmd) *cobra.Command {
	ptsArgsCompletion := completion.PtsArgsCompletion{
		PtsCompletion: &completion.PtsCompletion{AppID: &app, ArgsLen: -1, ConfigFile: &cmdr.ConfigFile},
	}
	cmd := &cobra.Command{
		Use: "unset <ptype>...",
		Example: template.CustomExample(
			"drycc timeouts unset web worker",
			map[string]string{
				"<ptype>": i18n.T("The process type as defined in your Procfile"),
			},
		),
		Short:             i18n.T("Unset resource timeouts for an app"),
		Long:              i18n.T("Unsets timeouts for an application. Default value (30s) or set by drycc controller"),
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: ptsArgsCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			ptypes := args
			return cmdr.TimeoutsUnset(app, ptypes)
		},
	}

	return cmd
}
