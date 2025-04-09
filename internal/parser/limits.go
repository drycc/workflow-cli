package parser

import (
	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/internal/template"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

func NewLimitsCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "limits",
		Short: i18n.T("Manage resource limits for your application"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.LimitsList(app, version)
		},
	}

	cmd.PersistentFlags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))
	cmd.Flags().IntVarP(&version, "version", "v", 0, i18n.T("The version for which the limit needs to be listed"))

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)
	releaseCompletion := completion.ReleaseCompletion{ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("version", releaseCompletion.CompletionFunc)

	cmd.AddCommand(limitsListCommand(cmdr))
	cmd.AddCommand(limitSetCommand(cmdr))
	cmd.AddCommand(limitUnsetCommand(cmdr))
	cmd.AddCommand(limitSpecsCommand(cmdr))
	cmd.AddCommand(limitPlansCommand(cmdr))
	return cmd
}

func limitsListCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Example: "drycc limits list",
		Short:   i18n.T("List resource limits for an app"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.LimitsList(app, version)
		},
	}

	cmd.Flags().IntVarP(&version, "version", "v", 0, i18n.T("The version for which the limit needs to be listed"))

	releaseCompletion := completion.ReleaseCompletion{ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("version", releaseCompletion.CompletionFunc)

	return cmd
}

func limitSetCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		limits []string
	}
	limitSetPlanCompletion := completion.LimitSetPlanCompletion{AppID: &app, ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "set <ptype>=<value>...",
		Args: cobra.MinimumNArgs(1),
		Example: template.CustomExample(
			"drycc limits set web=std1.large.c1m2",
			map[string]string{
				"<ptype>": i18n.T("The process type as defined in your Procfile, such as 'web' or 'worker'."),
				"<value>": i18n.T("The limit plan id to apply to the process type"),
			},
		),
		Short: i18n.T("Set resource limits for an app"),
		Long: i18n.T(`Sets resource limits for an application.

A resource limit is a finite resource within a pod which we can apply
restrictions through Kubernetes.`),
		ValidArgsFunction: limitSetPlanCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			flags.limits = args
			return cmdr.LimitsSet(app, flags.limits)
		},
	}

	return cmd
}

func limitUnsetCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		limits []string
	}
	ptsCompletion := completion.PtsCompletion{AppID: &app, ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "unset <ptype>...",
		Args: cobra.MinimumNArgs(1),
		Example: template.CustomExample(
			"drycc limits unset web",
			map[string]string{
				"<ptype>": i18n.T("The process type as defined in your Procfile, such as 'web' or 'worker'"),
			},
		),
		Short:             i18n.T("Unset resource limits for an app"),
		ValidArgsFunction: ptsCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			flags.limits = args
			return cmdr.LimitsUnset(app, flags.limits)
		},
	}

	return cmd
}

func limitSpecsCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		keywords string
		limit    int
	}

	cmd := &cobra.Command{
		Use:   "specs",
		Short: i18n.T("List specification information of the server"),
		Long:  i18n.T("List all available limit specs"),
		RunE: func(_ *cobra.Command, _ []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.LimitsSpecs(flags.keywords, results)
		},
	}

	cmd.Flags().StringVarP(&flags.keywords, "keywords", "k", "", i18n.T("Search keywords separated by commas, matching must satisfy all of the specified."))
	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))
	return cmd
}

func limitPlansCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		specID string
		cpu    int
		memory int
		limit  int
	}

	cmd := &cobra.Command{
		Use:   "plans",
		Short: i18n.T("List all available limit plans"),
		RunE: func(_ *cobra.Command, _ []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.LimitsPlans(flags.specID, flags.cpu, flags.memory, results)
		},
	}

	cmd.Flags().IntVar(&flags.cpu, "cpu", 0, i18n.T("Query plans that meet the specified number of cpu cores"))
	cmd.Flags().IntVarP(&flags.memory, "memory", "m", 0, i18n.T("Query plans that meet the specified memory capacity, unit GiB"))
	cmd.Flags().StringVar(&flags.specID, "spec-id", "", i18n.T("Query plans that meet the specified spec id, see [specs] subcommand."))
	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))
	cmd.Flags().SortFlags = false

	completion := completion.LimitSpecCompletion{ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("spec-id", completion.CompletionFunc)

	return cmd
}
