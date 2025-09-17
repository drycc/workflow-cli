package parser

import (
	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/internal/template"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

// NewLabelsCommand creates a command for managing application labels.
func NewLabelsCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "labels",
		Short: i18n.T("Manage labels of application"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.LabelsList(app)
		},
	}

	cmd.PersistentFlags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))
	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)

	cmd.AddCommand(labelsListCommand(cmdr))
	cmd.AddCommand(labelsSetCommand(cmdr))
	cmd.AddCommand(labelsUnsetCommand(cmdr))
	return cmd
}

func labelsListCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: i18n.T("List application's labels"),
		Long:  i18n.T("Prints a list of labels of the application"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.LabelsList(app)
		},
	}

	return cmd
}

func labelsSetCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		app  string
		tags []string
	}

	cmd := &cobra.Command{
		Use:  "set <key>=<value>...",
		Args: cobra.MinimumNArgs(1),
		Example: template.CustomExample(
			"drycc labels set FOO=bar",
			map[string]string{
				"<key>":   i18n.T(`the label key, for example: "git_repo" or "team"`),
				"<value>": i18n.T(`the label value, for example: "https://github.com/drycc/workflow" or "frontend"`),
			},
		),
		Short: i18n.T("Add new application's label"),
		Long: i18n.T(`Sets labels for an application.

A label is a key/value pair used to label an application. This label is a general information for drycc user.
Mostly used for administration/maintenance information, note for application. This information isn't send to scheduler.`),
		RunE: func(_ *cobra.Command, args []string) error {
			flags.tags = args
			return cmdr.LabelsSet(app, flags.tags)
		},
	}

	return cmd
}

func labelsUnsetCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		app  string
		keys []string
	}
	labelCompletion := completion.LabelCompletion{AppID: &app, ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use: "unset <key>...",
		Example: template.CustomExample(
			"drycc labels unset FOO",
			map[string]string{
				"<key>": i18n.T(`the label key to unset, for example: "git_repo" or "team"`),
			},
		),
		Args: cobra.MinimumNArgs(1),

		Short:             i18n.T("Remove application's label"),
		Long:              i18n.T("Unsets labels for an application"),
		ValidArgsFunction: labelCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			flags.keys = args
			return cmdr.LabelsUnset(app, flags.keys)
		},
	}

	return cmd
}
