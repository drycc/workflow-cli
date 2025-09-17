package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/internal/template"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

// NewReleasesCommand creates the releases command
func NewReleasesCommand(cmdr *commands.DryccCmd) *cobra.Command {
	ptsArgsCompletion := completion.PtsArgsCompletion{
		PtsCompletion: &completion.PtsCompletion{AppID: &app, ArgsLen: -1, ConfigFile: &cmdr.ConfigFile},
	}

	cmd := &cobra.Command{
		Use:               "releases [<ptype>...]",
		Example:           "drycc releases ptype1 ptype2",
		Short:             i18n.T("Manage releases of an application"),
		ValidArgsFunction: ptsArgsCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.ReleasesList(app, args, results)
		},
	}

	cmd.PersistentFlags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))
	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)

	cmd.AddCommand(releasesListCommand(cmdr))
	cmd.AddCommand(releasesInfoCommand(cmdr))
	cmd.AddCommand(releasesDeployCommand(cmdr))
	cmd.AddCommand(releasesRollbackCommand(cmdr))
	return cmd
}

func releasesListCommand(cmdr *commands.DryccCmd) *cobra.Command {
	ptsArgsCompletion := completion.PtsArgsCompletion{
		PtsCompletion: &completion.PtsCompletion{AppID: &app, ArgsLen: -1, ConfigFile: &cmdr.ConfigFile},
	}

	cmd := &cobra.Command{
		Use:               "list [<ptype>...]",
		Example:           "drycc releases list ptype1 ptype2",
		Short:             i18n.T("List an application's release history"),
		Long:              i18n.T("Lists release history for an application"),
		ValidArgsFunction: ptsArgsCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.ReleasesList(app, args, results)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))
	cmd.Flags().SortFlags = false

	return cmd
}

func releasesInfoCommand(cmdr *commands.DryccCmd) *cobra.Command {
	releaseCompletion := completion.ReleaseCompletion{AppID: &app, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "info <version>",
		Args: cobra.ExactArgs(1),
		Example: template.CustomExample(
			"drycc releases info 1",
			map[string]string{
				"<version>": i18n.T(`The release of the application, such as '1'`),
			},
		),
		Short:             i18n.T("Print information about a specific release"),
		Long:              i18n.T("Prints info about a particular release"),
		ValidArgsFunction: releaseCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			versionStr := args[0]
			version, err := strconv.Atoi(strings.TrimPrefix(versionStr, "v"))
			if err != nil {
				return fmt.Errorf("invalid version format: %v", err)
			}
			return cmdr.ReleasesInfo(app, version)
		},
	}

	return cmd
}

func releasesDeployCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		app     string
		force   bool
		confirm string
	}

	ptsArgsCompletion := completion.PtsArgsCompletion{
		PtsCompletion: &completion.PtsCompletion{AppID: &app, ArgsLen: -1, ConfigFile: &cmdr.ConfigFile},
	}

	cmd := &cobra.Command{
		Use:               "deploy [<ptype>...]",
		Short:             i18n.T("Deploy the latest release by process types"),
		ValidArgsFunction: ptsArgsCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			ptypes := args
			return cmdr.ReleasesDeploy(app, ptypes, flags.force, flags.confirm)
		},
	}

	cmd.Flags().BoolVarP(&flags.force, "force", "f", false, i18n.T("Force deploy"))
	cmd.Flags().StringVar(&flags.confirm, "confirm", "", i18n.T(`To proceed, type "yes"`))
	cmd.Flags().SortFlags = false

	return cmd
}

func releasesRollbackCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use: "rollback [<ptype>...] [version]",
		Example: template.CustomExample(
			"drycc releases rollback web 2",
			map[string]string{
				"<ptype>":   i18n.T(`The process name as defined in your Procfile, such as 'web'`),
				"<version>": i18n.T(`The release of the application, such as '1'`),
			},
		),
		Short: i18n.T("Return to a previous release"),
		Long:  i18n.T("Rolls back to a previous application release"),
		Args:  cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			// Handle arguments safely
			var ptypes []string
			var versionStr string

			if len(args) == 1 {
				versionStr = args[1]
			} else {
				ptypes = args[:len(args)-1]
				versionStr = args[len(args)-1]
			}
			version, err := versionFromString(versionStr)
			if err != nil {
				return err
			}
			return cmdr.ReleasesRollback(app, ptypes, version)
		},
	}

	return cmd
}

func versionFromString(version string) (int, error) {
	if version[:1] == "v" {
		if len(version) < 2 {
			return -1, fmt.Errorf("%s is not in the form 'v#'", version)
		}

		return strconv.Atoi(version[1:])
	}

	return strconv.Atoi(version)
}
