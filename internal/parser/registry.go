package parser

import (
	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/internal/template"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

// no default web
var registryFlags struct {
	ptype string
}

// NewRegistryCommand creates a command for managing registry information.
func NewRegistryCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "registry",
		Short: i18n.T("Manage private registry information for your application"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.RegistryList(app, registryFlags.ptype, version)
		},
	}

	cmd.PersistentFlags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))
	cmd.Flags().StringVarP(&registryFlags.ptype, "ptype", "p", "", i18n.T("The ptype for registry"))
	cmd.Flags().IntVarP(&version, "version", "v", 0, i18n.T("The version for which the registry needs to be listed"))

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)
	ptypeCompletion := completion.PtsCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("ptype", ptypeCompletion.CompletionFunc)

	cmd.AddCommand(registryListCommand(cmdr))
	cmd.AddCommand(registrySetCommand(cmdr))
	cmd.AddCommand(registryUnsetCommand(cmdr))
	return cmd
}

func registryListCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: i18n.T("List registry info for an app"),
		Long:  i18n.T("Lists registry information for an application"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.RegistryList(app, registryFlags.ptype, version)
		},
	}

	cmd.Flags().StringVarP(&registryFlags.ptype, "ptype", "p", "", i18n.T("The ptype for registry"))
	cmd.Flags().IntVarP(&version, "version", "v", 0, i18n.T("The version for which the registry needs to be listed"))

	ptypeCompletion := completion.PtsCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("ptype", ptypeCompletion.CompletionFunc)
	releaseCompletion := completion.ReleaseCompletion{ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("version", releaseCompletion.CompletionFunc)
	return cmd
}

func registrySetCommand(cmdr *commands.DryccCmd) *cobra.Command {
	// default web
	var flags struct {
		ptype string
	}
	cmd := &cobra.Command{
		Use:  "set <username> <password>",
		Args: cobra.ExactArgs(2),
		Example: template.CustomExample(
			"drycc registry set username password",
			map[string]string{
				"<username>": i18n.T("The username of the registry"),
				"<password>": i18n.T("The password of the registry"),
			},
		),
		Short: i18n.T("Set registry info for an app"),
		Long: i18n.T(`Sets registry information for an application. These credentials are the same as those used for
'podmain login' to the private registry.`),
		RunE: func(_ *cobra.Command, args []string) error {
			username := args[0]
			password := args[1]
			return cmdr.RegistrySet(app, flags.ptype, username, password)
		},
	}

	cmd.Flags().StringVarP(&flags.ptype, "ptype", "p", "", i18n.T("The ptype for registry"))
	cmd.MarkFlagRequired("ptype")

	ptypeCompletion := completion.PtsCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("ptype", ptypeCompletion.CompletionFunc)
	return cmd
}

func registryUnsetCommand(cmdr *commands.DryccCmd) *cobra.Command {
	// default web
	var flags struct {
		ptype string
	}
	cmd := &cobra.Command{
		Use:   "unset",
		Short: i18n.T("Unset registry info for an app"),
		Long:  i18n.T("Unsets registry information for an application"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.RegistryUnset(app, flags.ptype)
		},
	}

	cmd.Flags().StringVarP(&flags.ptype, "ptype", "p", "", i18n.T("The ptype for registry"))
	cmd.MarkFlagRequired("ptype")

	ptypeCompletion := completion.PtsCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("ptype", ptypeCompletion.CompletionFunc)
	return cmd
}
