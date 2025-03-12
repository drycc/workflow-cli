package parser

import (
	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

// NewDomainsCommand creates the domains command
func NewDomainsCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "domains",
		Short: i18n.T("Manage and assign domain names to your applications"),
		RunE: func(cmd *cobra.Command, args []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.DomainsList(app, results)
		},
	}

	cmd.PersistentFlags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))
	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)

	cmd.AddCommand(domainsListCommand(cmdr))
	cmd.AddCommand(domainsAddCommand(cmdr))
	cmd.AddCommand(domainsRemoveCommand(cmdr))
	return cmd
}

func domainsListCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: i18n.T("List domains bound to an application"),
		RunE: func(cmd *cobra.Command, args []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.DomainsList(app, results)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))

	return cmd
}

func domainsAddCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		ptype string
	}
	cmd := &cobra.Command{
		Use:   "add <domain>",
		Args:  cobra.ExactArgs(1),
		Short: i18n.T("Bind a domain to an application"),
		RunE: func(cmd *cobra.Command, args []string) error {
			domain := args[0]
			return cmdr.DomainsAdd(app, domain, flags.ptype)
		},
	}

	cmd.Flags().StringVarP(&flags.ptype, "ptype", "p", "web", i18n.T("The ptype for domain"))

	ptypeCompletion := completion.PtsCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile, AppID: &app}
	cmd.RegisterFlagCompletionFunc("ptype", ptypeCompletion.CompletionFunc)
	return cmd
}

func domainsRemoveCommand(cmdr *commands.DryccCmd) *cobra.Command {
	domainCompletion := completion.DomainCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:               "remove <domain>",
		Args:              cobra.ExactArgs(1),
		Short:             i18n.T("Unbinds a domain for an application"),
		ValidArgsFunction: domainCompletion.CompletionFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			domain := args[0]
			return cmdr.DomainsRemove(app, domain)
		},
	}

	return cmd
}
