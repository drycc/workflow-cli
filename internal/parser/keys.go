package parser

import (
	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/internal/template"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

func NewKeysCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keys",
		Short: i18n.T("Manage ssh keys used for 'git push' deployments"),
		RunE: func(_ *cobra.Command, _ []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.KeysList(results)
		},
	}
	cmd.AddCommand(keysListCommand(cmdr))
	cmd.AddCommand(keyAddCommand(cmdr))
	cmd.AddCommand(keyRemoveCommand(cmdr))
	return cmd
}

func keysListCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: i18n.T("List SSH keys for the logged in user"),

		RunE: func(_ *cobra.Command, _ []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.KeysList(results)
		},
	}

	return cmd
}

func keyAddCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		name string
		key  string
	}

	cmd := &cobra.Command{
		Use: "add [<name>] [<key>]",
		Example: template.CustomExample(
			"drycc keys add mykey /path/to/public.key",
			map[string]string{
				"<name>": i18n.T("Name of the SSH key"),
				"<key>":  i18n.T("A local file path to an SSH public key used to push application code"),
			},
		),
		Args:  cobra.MaximumNArgs(2),
		Short: i18n.T("Add an SSH key"),
		Long:  i18n.T("Adds SSH keys for the logged in user"),
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) == 1 {
				flags.key = args[0]
			} else if len(args) == 2 {
				flags.name = args[0]
				flags.key = args[1]
			}
			return cmdr.KeyAdd(flags.name, flags.key)
		},
	}

	return cmd
}

func keyRemoveCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		key string
	}
	keyCompletion := completion.KeyCompletion{ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:               "remove <key>",
		Args:              cobra.ExactArgs(1),
		Short:             i18n.T("Remove an SSH key"),
		Long:              i18n.T("Removes an SSH key for the logged in user"),
		ValidArgsFunction: keyCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			flags.key = args[0]
			return cmdr.KeyRemove(flags.key)
		},
	}

	return cmd
}
