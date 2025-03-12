package parser

import (
	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/internal/template"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

// NewTokensCommand creates the tokens command
func NewTokensCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tokens",
		Short: i18n.T("Manage user tokens"),
		RunE: func(cmd *cobra.Command, args []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.TokensList(results)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))

	cmd.AddCommand(tokensListCommand(cmdr))
	cmd.AddCommand(tokensAddCommand(cmdr))
	cmd.AddCommand(tokensRemoveCommand(cmdr))
	return cmd
}

func tokensListCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: i18n.T("Lists tokens visible to the current controller"),
		RunE: func(cmd *cobra.Command, args []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.TokensList(results)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))
	return cmd
}

func tokensAddCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		username string
		password string
	}

	cmd := &cobra.Command{
		Use: "add <alias>",
		Example: template.CustomExample(
			"drycc tokens add mytoken -u uname -p passwd",
			map[string]string{
				"<alias>": i18n.T("Provide a alias for controller authentication token"),
			},
		),
		Short: i18n.T("Add a token for controller authentication"),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			alias := args[0]
			_, err := cmdr.TokensAdd(nil, flags.username, flags.password, alias, "", true)
			return err
		},
	}

	cmd.Flags().StringVarP(&flags.username, "username", "u", "", i18n.T("Provide a username for the account"))
	cmd.Flags().StringVarP(&flags.password, "password", "p", "", i18n.T("Provide a password for the account"))

	must_flags := []string{"username", "password"}
	for _, must_flag := range must_flags {
		cmd.MarkFlagRequired(must_flag)
	}

	return cmd
}

func tokensRemoveCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		id      string
		confirm string
	}

	// coTokenCompletion
	tokenCompletion := completion.TokenCompletion{ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use: "remove <id>",
		Example: template.CustomExample(
			"drycc tokens remove 34103c3b-077e-4e37-a9ad-29ba4324ad8b",
			map[string]string{
				"<id>": i18n.T("The id of the token for controller authentication"),
			},
		),
		Short:             i18n.T("Remove a token for controller authentication"),
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: tokenCompletion.CompletionFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.id = args[0]
			return cmdr.TokensRemove(flags.id, flags.confirm)
		},
	}

	cmd.Flags().StringVarP(&flags.confirm, "confirm", "", "", i18n.T(`To proceed, type "yes"`))

	return cmd
}
