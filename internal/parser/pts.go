package parser

import (
	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/internal/template"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

// NewPtsCommand creates a command for managing process types.
func NewPtsCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pts",
		Short: i18n.T("Manage process types inside an app"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.PtsList(app, 1000)
		},
	}

	cmd.PersistentFlags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)

	cmd.AddCommand(ptsListCommand(cmdr))
	cmd.AddCommand(ptsDescribeCommand(cmdr))
	cmd.AddCommand(ptsRestartCommand(cmdr))
	cmd.AddCommand(ptsScaleCommand(cmdr))
	cmd.AddCommand(ptsCleanCommand(cmdr))
	return cmd
}

// ptsListCommand
func ptsListCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: i18n.T("List application process types"),
		Long:  i18n.T("Lists process types servicing an application"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.PtsList(app, 1000)
		},
	}

	return cmd
}

// ptsDescribeCommand
func ptsDescribeCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		ptype string
	}
	ptsArgsCompletion := completion.PtsArgsCompletion{
		PtsCompletion: &completion.PtsCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile},
	}
	cmd := &cobra.Command{
		Use:  "describe <ptype>",
		Args: cobra.ExactArgs(1),
		Example: template.CustomExample(
			"drycc pts describe web",
			map[string]string{
				"<ptype>": i18n.T("The process name as defined in your Procfile"),
			},
		),
		Short:             i18n.T("Describe a process type"),
		Long:              i18n.T("Print a detailed description of the selected process type"),
		ValidArgsFunction: ptsArgsCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			flags.ptype = args[0]
			return cmdr.PtsDescribe(app, flags.ptype)
		},
	}

	return cmd
}

// ptsRestartCommand
func ptsRestartCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		ptypes  []string
		confirm string
	}
	ptsArgsCompletion := completion.PtsArgsCompletion{
		PtsCompletion: &completion.PtsCompletion{AppID: &app, ArgsLen: -1, ConfigFile: &cmdr.ConfigFile},
	}
	cmd := &cobra.Command{
		Use:               "restart [<ptype>...]",
		Args:              cobra.ArbitraryArgs,
		Example:           "drycc pts restart web",
		Short:             i18n.T("Restart application process types"),
		Long:              i18n.T("Restart an application or process types"),
		ValidArgsFunction: ptsArgsCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			flags.ptypes = args
			return cmdr.PtsRestart(app, flags.ptypes, flags.confirm)
		},
	}

	cmd.Flags().StringVar(&flags.confirm, "confirm", "", i18n.T(`To proceed, type "yes"`))

	return cmd
}

// ptsScaleCommand
func ptsScaleCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		app   string
		scale []string // format like "web=5"
	}

	ptsSetArgsCompletion := completion.PtsSetArgsCompletion{
		PtsCompletion: &completion.PtsCompletion{AppID: &app, ArgsLen: -1, ConfigFile: &cmdr.ConfigFile},
	}
	cmd := &cobra.Command{
		Use:  "scale <ptype>=<num>...",
		Args: cobra.MinimumNArgs(1),
		Example: template.CustomExample(
			"drycc pts scale web=3",
			map[string]string{
				"<ptype>": i18n.T("The process name as defined in your Procfile"),
				"<num>":   i18n.T("The number of processes"),
			},
		),
		Short:             i18n.T("Scale process types of replicas"),
		Long:              i18n.T("Scales an application's processes by type"),
		ValidArgsFunction: ptsSetArgsCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			flags.scale = args
			return cmdr.PtsScale(app, flags.scale)
		},
	}

	// shortcuts scale has not app flag
	cmd.Flags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)

	return cmd
}

// ptsCleanCommand
func ptsCleanCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		ptypes []string
	}

	cmd := &cobra.Command{
		Use:  "clean <ptype>...",
		Args: cobra.MinimumNArgs(1),
		Example: template.CustomExample(
			"drycc pts clean web",
			map[string]string{
				"<ptype>": i18n.T("The process name as defined in your Procfile"),
			},
		),
		Short: i18n.T("Clean process types of not used"),
		RunE: func(_ *cobra.Command, args []string) error {
			flags.ptypes = args
			return cmdr.PtsClean(app, flags.ptypes)
		},
	}

	return cmd
}
