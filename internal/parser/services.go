package parser

import (
	"strconv"

	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/internal/template"

	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

func NewServicesCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "services",
		Short: i18n.T("Manage services for your applications"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.ServicesList(app)
		},
	}

	cmd.PersistentFlags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)

	cmd.AddCommand(servicesList(cmdr))
	cmd.AddCommand(servicesAdd(cmdr))
	cmd.AddCommand(servicesRemove(cmdr))
	return cmd
}

func servicesAdd(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		port     string
		protocol string
	}

	ptsArgsCompletion := completion.PtsArgsCompletion{
		PtsCompletion: &completion.PtsCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile},
	}
	cmd := &cobra.Command{
		Use:  "add <ptype> <port>:<target>",
		Args: cobra.ExactArgs(2),
		Example: template.CustomExample(
			"drycc services add web-new 80:8080",
			map[string]string{
				"<ptype>": i18n.T(`procfile type which should handle the request, e.g. webhooks (should be bind to the port PORT)
             only single extra service per Porcfile type could be created.`),
				"<port>":   i18n.T("The port that will be exposed by this service"),
				"<target>": i18n.T("Number or name of the port to access on the pods targeted by the service"),
			},
		),
		Short:             i18n.T("Create service for an application"),
		Long:              i18n.T("Creates extra service for an application and binds it to specific route of the main app domain"),
		ValidArgsFunction: ptsArgsCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			ptype := args[0]
			portTarget := args[1]
			return cmdr.ServicesAdd(app, ptype, portTarget, flags.protocol)
		},
	}

	cmd.Flags().StringVar(&flags.protocol, "protocol", "TCP", i18n.T("The IP protocol for this port. Supports TCP, UDP, and SCTP"))

	serviceProtocolCompletion := completion.ServiceProtocolCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("protocol", serviceProtocolCompletion.CompletionFunc)
	return cmd
}

func servicesList(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Example: "drycc services list",
		Short:   i18n.T("List application services"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.ServicesList(app)
		},
	}

	return cmd
}

func servicesRemove(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		app      string
		port     int
		protocol string
	}

	ServiceCompletion := completion.ServiceCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "remove <ptype> <port>",
		Args: cobra.ExactArgs(2),
		Example: template.CustomExample(
			"drycc services remove web-new 80",
			map[string]string{
				"<ptype>": i18n.T("procfile type which should handle the request, e.g. webhooks"),
				"<port>":  i18n.T("The port exposed by this service"),
			},
		),
		Short:             i18n.T("Remove service from an application"),
		ValidArgsFunction: ServiceCompletion.CompletionFunc,
		RunE: func(_ *cobra.Command, args []string) error {
			ptype := args[0]
			port, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			flags.port = port
			return cmdr.ServicesRemove(app, ptype, flags.protocol, port)
		},
	}

	cmd.Flags().StringVar(&flags.protocol, "protocol", "TCP", i18n.T("The IP protocol for this port. Supports TCP, UDP, and SCTP"))

	serviceProtocolCompletion := completion.ServiceProtocolCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("protocol", serviceProtocolCompletion.CompletionFunc)
	return cmd
}
