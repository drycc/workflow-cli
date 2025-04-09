package parser

import (
	"fmt"

	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/internal/template"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

// NewTLSCommand creates the TLS command
func NewTLSCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tls",
		Short: i18n.T("Manage TLS/SSL settings for applications"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.TLSInfo(app)
		},
	}

	cmd.PersistentFlags().StringVarP(&app, "app", "a", "", i18n.T("Application name"))

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)

	cmd.AddCommand(tlsInfoCommand(cmdr))
	cmd.AddCommand(tlsIssuerCommand(cmdr))
	cmd.AddCommand(tlsAutoCommand(cmdr))
	cmd.AddCommand(tlsForceCommand(cmdr))
	return cmd
}

func tlsInfoCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "info",
		Example: "drycc tls info",
		Short:   i18n.T("View TLS settings information"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.TLSInfo(app)
		},
	}

	return cmd
}

func tlsForceCommand(cmdr *commands.DryccCmd) *cobra.Command {
	TLSActionCompletion := completion.TLSActionCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use: "force <action>",
		Example: template.CustomExample(
			"drycc tls force enable",
			map[string]string{
				"<action>": i18n.T("Action to perform, 'enable' or 'disable'"),
			},
		),
		Args:              cobra.ExactArgs(1),
		ValidArgs:         []string{"enable", "disable"},
		Short:             i18n.T("Force TLS settings to HTTPS-only redirection"),
		ValidArgsFunction: TLSActionCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			action := args[0]
			if action == "enable" {
				return cmdr.TLSForceEnable(app)
			} else if action == "disable" {
				return cmdr.TLSForceDisable(app)
			}
			return fmt.Errorf("invalid action: %s, please use 'enable' or 'disable'", action)
		},
	}

	return cmd
}

func tlsAutoCommand(cmdr *commands.DryccCmd) *cobra.Command {
	TLSActionCompletion := completion.TLSActionCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use: "auto <action>",
		Example: template.CustomExample(
			"drycc tls auto enable",
			map[string]string{
				"<action>": i18n.T("Action to perform, 'enable' or 'disable'"),
			},
		),
		Args:              cobra.ExactArgs(1),
		ValidArgs:         []string{"enable", "disable"},
		Short:             i18n.T("Automatic certificate management"),
		ValidArgsFunction: TLSActionCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			action := args[0]
			if action == "enable" {
				return cmdr.TLSAutoEnable(app)
			} else if action == "disable" {
				return cmdr.TLSAutoDisable(app)
			}
			return fmt.Errorf("invalid action: %s, please use 'enable' or 'disable'", action)
		},
	}

	return cmd
}

func tlsIssuerCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		app       string
		email     string
		server    string
		keyID     string
		keySecret string
	}

	cmd := &cobra.Command{
		Use:   "issuer",
		Short: i18n.T("Configure automatic certificate management issuer details"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.TLSAutoIssuer(app, flags.email, flags.server, flags.keyID, flags.keySecret)
		},
	}

	cmd.Flags().StringVar(&flags.email, "email", "", i18n.T("ACME account email"))
	cmd.Flags().StringVar(&flags.server, "server", "", i18n.T("ACME server URL"))
	cmd.Flags().StringVar(&flags.keyID, "key-id", "", i18n.T("CA key ID"))
	cmd.Flags().StringVar(&flags.keySecret, "key-secret", "", i18n.T("CA key secret"))

	mustFlags := []string{"email", "server", "key-id", "key-secret"}
	for _, mustFlag := range mustFlags {
		cmd.MarkFlagRequired(mustFlag)
	}

	return cmd
}
