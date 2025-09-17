package parser

import (
	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/internal/template"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

// NewBuildsCommand creates a command for managing builds.
func NewBuildsCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "builds",
		Short: i18n.T("Manage builds created using 'git push'"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.BuildsInfo(app, version)
		},
	}

	cmd.PersistentFlags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))
	cmd.Flags().IntVarP(&version, "version", "v", 0, i18n.T("The version for which the build info needs to be displayed"))

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)
	releaseCompletion := completion.ReleaseCompletion{ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("version", releaseCompletion.CompletionFunc)

	cmd.AddCommand(buildsInfo(cmdr))
	cmd.AddCommand(buildsCreate(cmdr))
	cmd.AddCommand(buildsFetch(cmdr))

	return cmd
}

func buildsCreate(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		stack     string
		procfile  string
		dryccPath string
		confirm   string
	}
	cmd := &cobra.Command{
		Use:  `create <image>`,
		Args: cobra.ExactArgs(1),
		Example: template.CustomExample(
			`drycc create docker.io/library/nginx:latest`,
			map[string]string{
				"<image>": i18n.T("default container image"),
			},
		),
		Short: i18n.T("imports an image and deploys as a new release"),
		Long: i18n.T(`Creates a new build of an application. Imports an <image> and deploys it to Drycc
as a new release. If a Procfile or drycc.yaml is present in the current directory,
it will be used as the default for this application.`),
		RunE: func(_ *cobra.Command, args []string) error {
			image := args[0]
			err := cmdr.BuildsCreate(app, image, flags.stack, flags.procfile, flags.dryccPath, flags.confirm)
			return err
		},
	}

	// // shortcuts pull has not app flag
	cmd.Flags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))
	cmd.Flags().StringVarP(&flags.stack, "stack", "s", "container", i18n.T("The stack name for the application"))
	cmd.Flags().StringVarP(&flags.procfile, "procfile", "p", "Procfile", i18n.T("A YAML file used to supply a Procfile to the application"))
	cmd.Flags().StringVarP(&flags.dryccPath, "dryccpath", "d", ".drycc", i18n.T("Drycc config path to the application"))
	cmd.Flags().StringVarP(&flags.confirm, "confirm", "", "", i18n.T(`To proceed, type "yes"`))
	cmd.Flags().SortFlags = false

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)

	return cmd
}

// buildsInfo
func buildsInfo(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "info",
		Example: "drycc builds info",
		Short:   i18n.T("Print information about a specific build"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.BuildsInfo(app, version)
		},
	}
	cmd.Flags().IntVarP(&version, "version", "v", 0, i18n.T("The version for which the build info needs to be displayed"))

	releaseCompletion := completion.ReleaseCompletion{ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("version", releaseCompletion.CompletionFunc)
	return cmd
}

// buildsFetch
func buildsFetch(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		stack     string
		procfile  string
		dryccPath string
		confirm   string
		save      bool
	}

	cmd := &cobra.Command{
		Use:     "fetch",
		Example: "drycc builds fetch",
		Short:   i18n.T("Print process info about a specific build"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.BuildsFetch(app, version, flags.procfile, flags.dryccPath, flags.confirm, flags.save)
		},
	}
	cmd.Flags().IntVarP(&version, "version", "v", 0, i18n.T("The version for which the build info needs to be displayed"))
	cmd.Flags().StringVarP(&flags.procfile, "procfile", "p", "Procfile", i18n.T("A YAML file used to supply a Procfile to the application"))
	cmd.Flags().StringVarP(&flags.dryccPath, "dryccpath", "d", ".drycc", i18n.T("Drycc config path to the application"))
	cmd.Flags().StringVarP(&flags.confirm, "confirm", "", "", i18n.T(`To proceed, type "yes"`))
	cmd.Flags().BoolVarP(&flags.save, "save", "", false, i18n.T("Save process info to the local"))
	cmd.Flags().SortFlags = false

	releaseCompletion := completion.ReleaseCompletion{ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("version", releaseCompletion.CompletionFunc)
	return cmd
}
