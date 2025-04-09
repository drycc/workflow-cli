package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/completion"
	"github.com/drycc/workflow-cli/internal/template"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

func NewRoutesCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "routes",
		Short: i18n.T("Manage routes for your applications"),
		RunE: func(cmd *cobra.Command, args []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.RoutesList(app, results)
		},
	}

	cmd.PersistentFlags().StringVarP(&app, "app", "a", "", i18n.T("The uniquely identifiable name for the application"))
	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))

	appCompletion := completion.AppCompletion{ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("app", appCompletion.CompletionFunc)

	cmd.AddCommand(routesListCommand(cmdr))
	cmd.AddCommand(routesAddCommand(cmdr))
	cmd.AddCommand(routesGetCommand(cmdr))
	cmd.AddCommand(routesSetCommand(cmdr))
	cmd.AddCommand(routesAttachCommand(cmdr))
	cmd.AddCommand(routesDetachCommand(cmdr))
	cmd.AddCommand(routesRemoveCommand(cmdr))
	return cmd
}

func routesAddCommand(cmdr *commands.DryccCmd) *cobra.Command {

	routeKindCompletion := completion.RouteKindCompletion{ArgsLen: 1}
	cmd := &cobra.Command{
		Use:  "add <name> <kind> <backend>...",
		Args: cobra.MinimumNArgs(3),
		Example: template.CustomExample(
			`drycc routes add myroute HTTPRoute svc:8080,100`,
			map[string]string{
				"<name>":    i18n.T("The unique name of the route"),
				"<kind>":    i18n.T("The route kind, range: HTTPRoute,TCPRoute,UDPRoute,GRPCRoute,TLSRoute"),
				"<backend>": i18n.T(`The route's backend, pattern: <service>:<port>,<weight>`),
			},
		),
		Short:             i18n.T("Create a route for an application"),
		Long:              i18n.T("Creates routes for an application, provides a way to route requests"),
		ValidArgsFunction: routeKindCompletion.CompletionFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			kind := args[1]
			backends := args[2:]

			var backendRefs []api.BackendRefRequest
			for _, backend := range backends {
				params := strings.Split(backend, ",")
				if len(params) != 2 {
					return fmt.Errorf("invalid backend format %s", backend)
				}
				servicePort := strings.Split(params[0], ":")
				if len(servicePort) != 2 {
					return fmt.Errorf("backend must be in 'service:port,weight' format")
				}
				port, err := strconv.Atoi(servicePort[1])
				if err != nil {
					return err
				}
				weight, err := strconv.Atoi(params[1])
				if err != nil {
					return err
				}
				backendRefs = append(backendRefs, api.BackendRefRequest{
					Kind:   "Service",
					Name:   servicePort[0],
					Port:   int32(port),
					Weight: int32(weight),
				})
			}

			return cmdr.RoutesCreate(app, name, kind, backendRefs...)
		},
	}

	return cmd
}

func routesListCommand(cmdr *commands.DryccCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: i18n.T("List application routes"),
		RunE: func(cmd *cobra.Command, args []string) error {
			results, _ := commands.ResponseLimit(limit)
			return cmdr.RoutesList(app, results)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 0, i18n.T("The maximum number of results to display"))

	return cmd
}

func routesGetCommand(cmdr *commands.DryccCmd) *cobra.Command {
	routeCompletion := completion.RouteCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "get <name>",
		Args: cobra.ExactArgs(1),
		Example: template.CustomExample(
			"drycc routes get myroute",
			map[string]string{
				"<name>": i18n.T("The unique name of the route"),
			},
		),
		Short:             i18n.T("Get route rules"),
		Long:              i18n.T("Get route rules for an application"),
		ValidArgsFunction: routeCompletion.CompletionFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			return cmdr.RoutesGet(app, name)
		},
	}

	return cmd
}

func routesSetCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		rulesFile string
	}
	routeCompletion := completion.RouteCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "set <name>",
		Args: cobra.ExactArgs(1),
		Example: template.CustomExample(
			"drycc routes set myroute --rules-file=rules.yaml",
			map[string]string{
				"<name>": i18n.T("The unique name of the route"),
			},
		),
		Short:             i18n.T("Set route rules"),
		Long:              i18n.T("Set route rules for an application"),
		ValidArgsFunction: routeCompletion.CompletionFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			return cmdr.RoutesSet(app, name, flags.rulesFile)
		},
	}

	cmd.Flags().StringVar(&flags.rulesFile, "rules-file", "", i18n.T("Rules-file is the file name of rule to apply"))
	cmd.MarkFlagRequired("rules-file")

	return cmd
}

func routesAttachCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		port    int
		gateway string
	}
	routeCompletion := completion.RouteCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "attach <name>",
		Args: cobra.ExactArgs(1),
		Example: template.CustomExample(
			"drycc routes attach myroute --gateway=gw1 --port=80",
			map[string]string{
				"<name>": i18n.T("The unique name of the route"),
			},
		),
		Short:             i18n.T("Attach to gateway"),
		ValidArgsFunction: routeCompletion.CompletionFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			return cmdr.RoutesAttach(app, name, flags.port, flags.gateway)
		},
	}

	cmd.Flags().StringVar(&flags.gateway, "gateway", "", i18n.T("The unique name of the gaetway"))
	cmd.Flags().IntVarP(&flags.port, "port", "", 0, i18n.T("Port is the network port, the gateway listener expects to receive"))

	mustFlags := []string{"gateway", "port"}
	for _, mustFlag := range mustFlags {
		cmd.MarkFlagRequired(mustFlag)
	}

	gatewayNameCompletion := completion.GatewayNameCompletion{AppID: &app, ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("gateway", gatewayNameCompletion.CompletionFunc)
	return cmd
}

func routesDetachCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		port    int
		gateway string
	}
	routeCompletion := completion.RouteCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "detach <name>",
		Args: cobra.ExactArgs(1),
		Example: template.CustomExample(
			"drycc routes detach myroute --gateway=gw1 --port=80",
			map[string]string{
				"<name>": i18n.T("The unique name of the route"),
			},
		),
		Short:             i18n.T("Dettach to gateway"),
		ValidArgsFunction: routeCompletion.CompletionFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			return cmdr.RoutesDetach(app, name, flags.port, flags.gateway)
		},
	}

	cmd.Flags().StringVar(&flags.gateway, "gateway", "", i18n.T("The unique name of the gaetway"))
	cmd.Flags().IntVarP(&flags.port, "port", "", 0, i18n.T("Port is the network port, the gateway listener expects to receive"))

	mustFlags := []string{"port", "protocol"}
	for _, mustFlag := range mustFlags {
		cmd.MarkFlagRequired(mustFlag)
	}

	gatewayNameCompletion := completion.GatewayNameCompletion{AppID: &app, ArgsLen: -1, ConfigFile: &cmdr.ConfigFile}
	cmd.RegisterFlagCompletionFunc("gateway", gatewayNameCompletion.CompletionFunc)
	return cmd
}

func routesRemoveCommand(cmdr *commands.DryccCmd) *cobra.Command {

	routeCompletion := completion.RouteCompletion{AppID: &app, ArgsLen: 0, ConfigFile: &cmdr.ConfigFile}
	cmd := &cobra.Command{
		Use:  "remove <name>",
		Args: cobra.ExactArgs(1),
		Example: template.CustomExample(
			"drycc routes remove myroute",
			map[string]string{
				"<name>": i18n.T("The unique name of the route"),
			},
		),
		Short:             i18n.T("Remove a route from an application"),
		ValidArgsFunction: routeCompletion.CompletionFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			return cmdr.RoutesRemove(app, name)
		},
	}

	return cmd
}
