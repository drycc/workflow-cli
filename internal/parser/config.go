package parser

import (
	"strings"

	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/internal/template"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

var configFlags struct {
	ptype string
	group string
}

// NewConfigCommand creates the config command
func NewConfigCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cobra.EnableCommandSorting = false
	cmd := &cobra.Command{
		Use:   "config",
		Short: i18n.T("Manage environment variables that define app config"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.ConfigInfo(app, configFlags.ptype, configFlags.group, version)
		},
	}

	cmd.PersistentFlags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))
	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)

	cmd.Flags().StringVarP(&configFlags.ptype, "ptype", "p", "", i18n.T("The ptype for which the config needs to be listed"))
	cmd.Flags().StringVarP(&configFlags.group, "group", "g", "", i18n.T("The group for which the config needs to be listed"))
	cmd.Flags().IntVarP(&version, "version", "v", 0, i18n.T("The version for which the config needs to be listed"))

	ptsCompletion := completion.PtsCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("ptype", ptsCompletion.CompletionFunc)
	configGroupCompletion := completion.ConfigGroupCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("group", configGroupCompletion.CompletionFunc)
	releaseCompletion := completion.ReleaseCompletion{ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("version", releaseCompletion.CompletionFunc)

	cmd.AddCommand(configInfoCommand(cmdr))
	cmd.AddCommand(configSetCommand(cmdr))
	cmd.AddCommand(configUnsetCommand(cmdr))
	cmd.AddCommand(configPullCommand(cmdr))
	cmd.AddCommand(configPushCommand(cmdr))
	cmd.AddCommand(configAttachCommand(cmdr))
	cmd.AddCommand(configDetachCommand(cmdr))
	return cmd
}

func configInfoCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: i18n.T("An app config info"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.ConfigInfo(app, configFlags.ptype, configFlags.group, version)
		},
	}

	cmd.Flags().StringVarP(&configFlags.ptype, "ptype", "p", "", i18n.T("The ptype for which the config needs to be listed"))
	cmd.Flags().StringVarP(&configFlags.group, "group", "g", "", i18n.T("The group for which the config needs to be listed"))
	cmd.Flags().IntVarP(&version, "version", "v", 0, i18n.T("The version for which the config needs to be listed"))

	ptsCompletion := completion.PtsCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("ptype", ptsCompletion.CompletionFunc)
	configGroupCompletion := completion.ConfigGroupCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("group", configGroupCompletion.CompletionFunc)
	releaseCompletion := completion.ReleaseCompletion{ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("version", releaseCompletion.CompletionFunc)
	return cmd
}

func configSetCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		group   string
		confirm string
	}

	cmd := &cobra.Command{
		Use:   "set <key>=<value>...",
		Args:  cobra.MinimumNArgs(1),
		Short: i18n.T("Set environment variables for an app"),
		Long:  i18n.T("Sets environment variables for an application or config group"),
		RunE: func(_ *cobra.Command, args []string) error {
			return cmdr.ConfigSet(app, configFlags.ptype, configFlags.group, args, flags.confirm)
		},
	}

	cmd.Flags().StringVarP(&configFlags.ptype, "ptype", "p", "", i18n.T("The ptype for which the config needs to be set"))
	cmd.Flags().StringVarP(&configFlags.group, "group", "g", "", i18n.T("The group for which the config needs to be set"))
	cmd.Flags().StringVarP(&flags.confirm, "confirm", "", "", i18n.T("To proceed, type 'yes'"))
	cmd.Flags().SortFlags = false
	cmd.MarkFlagsMutuallyExclusive("ptype", "group")

	ptsCompletion := completion.PtsCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("ptype", ptsCompletion.CompletionFunc)
	configGroupCompletion := completion.ConfigGroupCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("group", configGroupCompletion.CompletionFunc)

	return cmd
}

func configUnsetCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		confirm string
	}

	cmd := &cobra.Command{
		Use:   "unset <key>...",
		Args:  cobra.MinimumNArgs(1),
		Short: i18n.T("Unset environment variables for an app"),
		Long:  i18n.T("Unsets an environment variable for an application or config group"),
		RunE: func(_ *cobra.Command, args []string) error {
			return cmdr.ConfigUnset(app, configFlags.ptype, configFlags.group, args, flags.confirm)
		},
	}

	cmd.Flags().StringVarP(&configFlags.ptype, "ptype", "p", "", i18n.T("The ptype for which the config needs to be unset"))
	cmd.Flags().StringVarP(&configFlags.group, "group", "g", "", i18n.T("The group for which the config needs to be unset"))
	cmd.Flags().StringVarP(&flags.confirm, "confirm", "", "", i18n.T("To proceed, type 'yes'"))
	cmd.Flags().SortFlags = false
	cmd.MarkFlagsMutuallyExclusive("ptype", "group")

	ptsCompletion := completion.PtsCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("ptype", ptsCompletion.CompletionFunc)
	configGroupCompletion := completion.ConfigGroupCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("group", configGroupCompletion.CompletionFunc)

	return cmd
}

func configPullCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		path        string
		interactive bool
		overwrite   bool
	}

	cmd := &cobra.Command{
		Use:   "pull",
		Short: i18n.T("Pull environment variables to the path"),
		Long: i18n.T(`Extract all environment variables from an application or config group. for local use.

The environmental variables can be piped into a file, 'drycc config pull > file',
or stored locally in a file named .env. This file can be
read by foreman to load the local environment for your app.`),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.ConfigPull(app, configFlags.ptype, configFlags.group, flags.path, flags.interactive, flags.overwrite)
		},
	}

	cmd.Flags().StringVarP(&configFlags.ptype, "ptype", "p", "", i18n.T("The ptype for which the config needs to be pull"))
	cmd.Flags().IntVarP(&version, "version", "v", 0, i18n.T("The version for which the config needs to be pull"))
	cmd.Flags().StringVarP(&configFlags.group, "group", "g", "", i18n.T("The group for which the config needs to be pull"))
	cmd.Flags().StringVar(&flags.path, "path", ".env", i18n.T("A path leading to an environment file"))
	cmd.Flags().BoolVarP(&flags.interactive, "interactive", "i", false, i18n.T("Prompts for each value to be overwritten"))
	cmd.Flags().BoolVarP(&flags.overwrite, "overwrite", "o", false, i18n.T("Allows you to have the pull overwrite keys to the path"))
	cmd.Flags().SortFlags = false
	cmd.MarkFlagsMutuallyExclusive("ptype", "group")

	ptsCompletion := completion.PtsCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("ptype", ptsCompletion.CompletionFunc)
	configGroupCompletion := completion.ConfigGroupCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("group", configGroupCompletion.CompletionFunc)
	releaseCompletion := completion.ReleaseCompletion{ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("version", releaseCompletion.CompletionFunc)

	return cmd
}

func configPushCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		path    string
		confirm string
	}

	cmd := &cobra.Command{
		Use:   "push",
		Short: i18n.T("Push environment variables from the path"),
		Long: i18n.T(`Sets environment variables for an application or config group.

This file can be read by foreman
to load the local environment for your app. The file should be piped via
stdin, 'drycc config push < .env', or using the --path option.`),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.ConfigPush(app, configFlags.ptype, configFlags.group, flags.path, flags.confirm)
		},
	}

	cmd.Flags().StringVarP(&configFlags.ptype, "ptype", "p", "", i18n.T("The ptype for which the config needs to be push"))
	cmd.Flags().StringVarP(&configFlags.group, "group", "g", "", i18n.T("The group for which the config needs to be push"))
	cmd.Flags().StringVar(&flags.path, "path", ".env", i18n.T("A path leading to an environment file"))
	cmd.Flags().StringVar(&flags.confirm, "confirm", "", i18n.T("To proceed, type 'yes'"))
	cmd.Flags().SortFlags = false
	cmd.MarkFlagsMutuallyExclusive("ptype", "group")

	ptypeCompletion := completion.PtsCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("ptype", ptypeCompletion.CompletionFunc)
	configGroupCompletion := completion.ConfigGroupCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("group", configGroupCompletion.CompletionFunc)
	return cmd
}

func configAttachCommand(cmdr *commands.DryccCmd) *cobra.Command {
	configPtsGroupArgsCompletion := completion.ConfigPtsGroupArgsCompletion{
		AppID:      &app,
		ConfigFile: &cmdr.ConfigFile,
	}
	cmd := &cobra.Command{
		Use: "attach <ptype> <group>...",
		Example: template.CustomExample(
			"drycc config attach web group1 group2",
			map[string]string{
				"<ptype>": i18n.T("The ptype that requires attach configurations"),
				"<group>": i18n.T("The group that requires attach configurations"),
			},
		),
		Args:              cobra.MinimumNArgs(2),
		ValidArgsFunction: configPtsGroupArgsCompletion.CompletionFunc,
		Short:             i18n.T("Selects environment groups to attach an app ptype"),
		RunE: func(_ *cobra.Command, args []string) error {
			configFlags.ptype = args[0]
			return cmdr.ConfigAttach(app, configFlags.ptype, strings.Join(args[1:], ","))
		},
	}

	return cmd
}

func configDetachCommand(cmdr *commands.DryccCmd) *cobra.Command {
	configPtsGroupArgsCompletion := completion.ConfigPtsGroupArgsCompletion{
		AppID:      &app,
		ConfigFile: &cmdr.ConfigFile,
	}
	cmd := &cobra.Command{
		Use:  "detach <ptype> <group>...",
		Args: cobra.MinimumNArgs(2),
		Example: template.CustomExample(
			"drycc config detach web group1 group2",
			map[string]string{
				"<ptype>": i18n.T("The ptype that requires detach configurations"),
				"<group>": i18n.T("The group that requires detach configurations"),
			},
		),
		Short:             i18n.T("Selects environment groups to detach an app ptype"),
		ValidArgsFunction: configPtsGroupArgsCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			configFlags.ptype = args[0]
			return cmdr.ConfigDetach(app, configFlags.ptype, strings.Join(args[1:], ","))
		},
	}

	return cmd
}
