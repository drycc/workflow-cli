package parser

import (
	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/internal/template"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

// NewCertsCommand creates the certs command
func NewCertsCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "certs",
		Short: i18n.T("Manage SSL endpoints for an app"),
		RunE: func(_ *cobra.Command, _ []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.CertsList(app, results)
		},
	}

	cmd.PersistentFlags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))
	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)

	cmd.AddCommand(certsListCommand(cmdr))
	cmd.AddCommand(certAddCommand(cmdr))
	cmd.AddCommand(certRemoveCommand(cmdr))
	cmd.AddCommand(certInfoCommand(cmdr))
	cmd.AddCommand(certAttachCommand(cmdr))
	cmd.AddCommand(certDetachCommand(cmdr))
	return cmd
}

func certsListCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: i18n.T("List SSL certificates for an application"),
		Long:  i18n.T("Show certificate information for an SSL application"),
		RunE: func(_ *cobra.Command, _ []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.CertsList(app, results)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))

	return cmd
}

func certAddCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use: "add <name> <cert> <key>",
		Example: template.CustomExample(
			"drycc certs add cert-name default.crt default.key",
			map[string]string{
				"<name>": i18n.T("Name of the certificate to reference it by"),
				"<cert>": i18n.T("The public key of the SSL certificate"),
				"<key>":  i18n.T("The private key of the SSL certificate"),
			},
		),
		Args:  cobra.ExactArgs(3),
		Short: i18n.T("Add an SSL certificate to an app"),
		RunE: func(_ *cobra.Command, args []string) error {
			name := args[0]
			cert := args[1]
			key := args[2]
			return cmdr.CertAdd(app, cert, key, name)
		},
	}

	return cmd
}

func certRemoveCommand(cmdr *commands.DryccCmd) *cobra.Command {

	certCompletion := completion.CertCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use: "remove <name>",
		Example: template.CustomExample(
			"drycc certs remove cert-name",
			map[string]string{
				"<name>": i18n.T("The name of the cert to remove from the app"),
			},
		),
		Args:              cobra.ExactArgs(1),
		Short:             i18n.T("Remove an SSL certificate from an app"),
		ValidArgsFunction: certCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			name := args[0]
			return cmdr.CertRemove(app, name)
		},
	}

	return cmd
}

func certInfoCommand(cmdr *commands.DryccCmd) *cobra.Command {
	certCompletion := completion.CertCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use: "info <name>",
		Example: template.CustomExample(
			"drycc certs info cert-name",
			map[string]string{
				"<name>": i18n.T("The name of the cert to get information from"),
			},
		),
		Args:              cobra.ExactArgs(1),
		Short:             i18n.T("Get detailed informaton about the certificate"),
		Long:              i18n.T("Fetch more detailed information about a certificate"),
		ValidArgsFunction: certCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			name := args[0]
			return cmdr.CertInfo(app, name)
		},
	}

	return cmd
}

func certAttachCommand(cmdr *commands.DryccCmd) *cobra.Command {
	certTachCompletion := completion.CertDomainTachCompletion{AppID: &app, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use: "attach <name> <domain>",
		Example: template.CustomExample(
			"drycc certs attach cert-name id.example.com",
			map[string]string{
				"<name>":   i18n.T("Name of the certificate to attach domain to"),
				"<domain>": i18n.T("Common name of the domain to attach to (needs to already be in the system)"),
			},
		),
		Args:              cobra.ExactArgs(2),
		Short:             i18n.T("Attach an SSL certificate to a domain"),
		Long:              i18n.T("Attach a certificate to a domain"),
		ValidArgsFunction: certTachCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			name := args[0]
			domain := args[1]
			return cmdr.CertAttach(app, name, domain)
		},
	}

	return cmd
}

func certDetachCommand(cmdr *commands.DryccCmd) *cobra.Command {
	certTachCompletion := completion.CertDomainTachCompletion{AppID: &app, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use: "detach <name> <domain>",
		Example: template.CustomExample(
			"drycc certs attach cert-name id.example.com",
			map[string]string{
				"<name>":   i18n.T("Name of the certificate to deatch from a domain"),
				"<domain>": i18n.T("Common name of the domain to detach from"),
			},
		),
		Args:              cobra.ExactArgs(2),
		Short:             i18n.T("Detach an SSL certificate from a domain"),
		Long:              i18n.T("Detach certificate from a domain"),
		ValidArgsFunction: certTachCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			name := args[0]
			domain := args[1]
			return cmdr.CertDetach(app, name, domain)
		},
	}
	return cmd
}
