package parser

import (
	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/internal/template"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

// NewAutoscaleCommand creates the autoscale command
func NewAutoscaleCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "autoscale",
		Short: i18n.T("Manage autoscale for applications"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmdr.AutoscaleList(app)
		},
	}

	cmd.PersistentFlags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))
	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)

	cmd.AddCommand(autoscaleListCommand(cmdr))
	cmd.AddCommand(autoscaleSetCommand(cmdr))
	cmd.AddCommand(autoscaleUnsetCommand(cmdr))
	return cmd
}

func autoscaleListCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: i18n.T("List autoscale options for an application"),
		Long:  i18n.T("Prints a list of autoscale options for the application"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmdr.AutoscaleList(app)
		},
	}
	return cmd
}

func autoscaleSetCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		min        int
		max        int
		cpuPercent int
	}

	ptsArgsCompletion := completion.PtsArgsCompletion{
		PtsCompletion: &completion.PtsCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile},
	}

	cmd := &cobra.Command{
		Use: "set <ptype> --min=<min> --max=<max> --cpu-percent=<percent>",
		Example: template.CustomExample(
			"drycc autoscale set web --min=1 --max=10 --cpu-percent=50",
			map[string]string{
				"<ptype>": i18n.T("The process type to add to the application's autoscale settings"),
			},
		),
		Args:              cobra.ExactArgs(1),
		Short:             i18n.T("Turn on autoscale for an app"),
		Long:              i18n.T("Set autoscale option per process type for an app"),
		ValidArgsFunction: ptsArgsCompletion.CompletionFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			ptype := args[0]
			return cmdr.AutoscaleSet(app, ptype, flags.min, flags.max, flags.cpuPercent)
		},
	}

	cmd.Flags().IntVarP(&flags.min, "min", "", 0, i18n.T("minimum replicas to keep around"))
	cmd.Flags().IntVarP(&flags.max, "max", "", 0, i18n.T("Max replicas to scale up to"))
	cmd.Flags().IntVarP(&flags.cpuPercent, "cpu-percent", "", 0, i18n.T("Target CPU utilization percentage"))
	cmd.Flags().SortFlags = false

	mustFlags := []string{"min", "max", "cpu-percent"}
	for _, must_flag := range mustFlags {
		cmd.MarkFlagRequired(must_flag)
	}

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)
	return cmd
}

func autoscaleUnsetCommand(cmdr *commands.DryccCmd) *cobra.Command {
	ptsArgsCompletion := completion.PtsArgsCompletion{
		PtsCompletion: &completion.PtsCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile},
	}
	cmd := &cobra.Command{
		Use: "unset <ptype>",
		Example: template.CustomExample(
			"drycc autoscale unset web",
			map[string]string{
				"<ptype>": i18n.T("The process type to remove from the application's autoscale settings"),
			},
		),
		Args:              cobra.ExactArgs(1),
		Short:             i18n.T("Turn off autoscale for an app"),
		Long:              i18n.T("Unset autoscale per process type for an app"),
		ValidArgsFunction: ptsArgsCompletion.CompletionFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			ptype := args[0]
			return cmdr.AutoscaleUnset(app, ptype)
		},
	}
	return cmd
}
