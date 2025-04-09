package parser

import (
	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/internal/template"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

// NewTagsCommand creates the tags command
func NewTagsCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		ptype string
	}

	cmd := &cobra.Command{
		Use:   "tags",
		Short: i18n.T("Manage tags for application containers"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.TagsList(app, flags.ptype, version)
		},
	}

	cmd.PersistentFlags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))
	cmd.Flags().StringVarP(&flags.ptype, "ptype", "p", "", i18n.T("The process name as defined in your Procfile"))
	cmd.Flags().IntVarP(&version, "version", "v", 0, i18n.T("The version for which the tag needs to be listed"))

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)
	ptypeCompletion := completion.PtsCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("ptype", ptypeCompletion.CompletionFunc)
	releaseCompletion := completion.ReleaseCompletion{ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("version", releaseCompletion.CompletionFunc)

	cmd.AddCommand(tagsListCommand(cmdr))
	cmd.AddCommand(tagsSetCommand(cmdr))
	cmd.AddCommand(tagsUnsetCommand(cmdr))
	return cmd
}

func tagsListCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		ptype string
	}

	ptsArgsCompletion := completion.PtsArgsCompletion{
		PtsCompletion: &completion.PtsCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile},
	}
	cmd := &cobra.Command{
		Use:               "list",
		Short:             i18n.T("List tags for an application"),
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: ptsArgsCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.TagsList(app, flags.ptype, version)
		},
	}

	cmd.Flags().StringVarP(&flags.ptype, "ptype", "p", "web", i18n.T("The process name as defined in your Procfile"))
	cmd.Flags().IntVarP(&version, "version", "v", 0, i18n.T("The version for which the tag needs to be listed"))

	ptypeCompletion := completion.PtsCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("ptype", ptypeCompletion.CompletionFunc)
	releaseCompletion := completion.ReleaseCompletion{ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("version", releaseCompletion.CompletionFunc)
	return cmd
}

func tagsSetCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		ptype string
	}
	ptsArgsCompletion := completion.PtsArgsCompletion{
		PtsCompletion: &completion.PtsCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile},
	}
	cmd := &cobra.Command{
		Use: "set <key>=<value>...",
		Example: template.CustomExample(
			"drycc tags set environ=prod rack=1",
			map[string]string{
				"<key>":   i18n.T("The tag key"),
				"<value>": i18n.T("The tag value"),
			},
		),
		Short: i18n.T("Set tags for an application"),
		Long: i18n.T(`Sets tags for an application.

A tag is a key/value pair used to tag an application's containers and is passed to the
scheduler. This is often used to restrict workloads to specific hosts matching the
scheduler-configured metadata.`),
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: ptsArgsCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			tags := args[0:]
			return cmdr.TagsSet(app, flags.ptype, tags)
		},
	}

	cmd.Flags().StringVarP(&flags.ptype, "ptype", "p", "web", i18n.T("The process name as defined in your Procfile"))

	ptypeCompletion := completion.PtsCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("ptype", ptypeCompletion.CompletionFunc)

	return cmd
}

func tagsUnsetCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		ptype string
	}
	tagCompletion := completion.TagCompletion{AppID: &app, Ptype: &flags.ptype, ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use: "unset <key>...",
		Example: template.CustomExample(
			"drycc tags unset environ rack",
			map[string]string{
				"<key>": i18n.T("The tag key to unset"),
			},
		),
		Short:             i18n.T("Unset tags for an application"),
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: tagCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			keys := args[0:]
			return cmdr.TagsUnset(app, flags.ptype, keys)
		},
	}

	cmd.Flags().StringVarP(&flags.ptype, "ptype", "p", "web", i18n.T("The process name as defined in your Procfile"))

	ptypeCompletion := completion.PtsCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("ptype", ptypeCompletion.CompletionFunc)

	return cmd
}
