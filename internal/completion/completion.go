package completion

import (
	"fmt"
	"strconv"
	"strings"

	"slices"

	"github.com/drycc/controller-sdk-go/apps"
	"github.com/drycc/controller-sdk-go/appsettings"
	"github.com/drycc/controller-sdk-go/certs"
	"github.com/drycc/controller-sdk-go/config"
	"github.com/drycc/controller-sdk-go/domains"
	"github.com/drycc/controller-sdk-go/gateways"
	"github.com/drycc/controller-sdk-go/keys"
	"github.com/drycc/controller-sdk-go/limits"
	"github.com/drycc/controller-sdk-go/perms"
	"github.com/drycc/controller-sdk-go/ps"
	"github.com/drycc/controller-sdk-go/pts"
	"github.com/drycc/controller-sdk-go/releases"
	"github.com/drycc/controller-sdk-go/resources"
	"github.com/drycc/controller-sdk-go/routes"
	"github.com/drycc/controller-sdk-go/services"
	"github.com/drycc/controller-sdk-go/tokens"
	"github.com/drycc/controller-sdk-go/volumes"
	"github.com/drycc/workflow-cli/internal/utils"
	"github.com/drycc/workflow-cli/pkg/settings"
	"github.com/spf13/cobra"
)

type Completion interface {
	CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective)
}

type AppCompletion struct {
	ArgsLen    int
	ConfigFile *string
}

func (c *AppCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if s, err := settings.Load(*c.ConfigFile); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
		if apps, _, err := apps.List(s.Client, -1); err == nil {
			var results []string
			for _, app := range apps {
				if strings.HasPrefix(app.ID, toComplete) {
					results = append(results, app.ID)
				}
			}
			return results, cobra.ShellCompDirectiveNoFileComp
		}
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

type CertCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

func (c *CertCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if appID, s, err := utils.LoadAppSettings(*c.ConfigFile, *c.AppID); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
		if certs, _, err := certs.List(s.Client, appID, -1); err == nil {
			var results []string
			for _, cert := range certs {
				if strings.HasPrefix(cert.Name, toComplete) {
					results = append(results, cert.Name)
				}
			}
			return results, cobra.ShellCompDirectiveNoFileComp
		}
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

type CertDomainTachCompletion struct {
	AppID      *string
	ConfigFile *string
}

func (c *CertDomainTachCompletion) CompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		certCompletion := CertCompletion{AppID: c.AppID, ConfigFile: c.ConfigFile}
		return certCompletion.CompletionFunc(cmd, args, toComplete)
	} else {
		domainCompletion := DomainCompletion{AppID: c.AppID, ArgsLen: 1, ConfigFile: c.ConfigFile}
		return domainCompletion.CompletionFunc(cmd, args, toComplete)
	}
}

type ConfigPtsGroupArgsCompletion struct {
	AppID      *string
	ConfigFile *string
}

func (c *ConfigPtsGroupArgsCompletion) CompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		ptsCompletion := PtsCompletion{AppID: c.AppID, ConfigFile: c.ConfigFile}
		return ptsCompletion.CompletionFunc(cmd, args, toComplete)
	} else {
		groups := args[1:]
		if appID, s, err := utils.LoadAppSettings(*c.ConfigFile, *c.AppID); err == nil {
			if config, err := config.List(s.Client, appID, -1); err == nil {
				var results []string
				for _, value := range config.Values {
					if value.Group != "" && strings.HasPrefix(value.Group, toComplete) {
						if !slices.Contains(groups, value.Group) && !slices.Contains(results, value.Group) {
							results = append(results, value.Group)
						}
					}
				}
				return results, cobra.ShellCompDirectiveNoFileComp
			}
		}
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

type ConfigGroupCompletion struct {
	AppID      *string
	ConfigFile *string
	ArgsLen    int
}

func (c *ConfigGroupCompletion) CompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

	if appID, s, err := utils.LoadAppSettings(*c.ConfigFile, *c.AppID); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
		if config, err := config.List(s.Client, appID, -1); err == nil {
			var results []string
			for _, value := range config.Values {
				if value.Group != "" && strings.HasPrefix(value.Group, toComplete) {
					results = append(results, value.Group)
				}
			}
			return results, cobra.ShellCompDirectiveNoFileComp
		}
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

type DomainCompletion struct {
	AppID      *string
	ConfigFile *string
	ArgsLen    int
}

func (c *DomainCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if appID, s, err := utils.LoadAppSettings(*c.ConfigFile, *c.AppID); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
		if domains, _, err := domains.List(s.Client, appID, -1); err == nil {
			var results []string
			for _, domain := range domains {
				if strings.HasPrefix(domain.Domain, toComplete) {
					results = append(results, domain.Domain)
				}
			}
			return results, cobra.ShellCompDirectiveNoFileComp
		}
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

type GatewayProtocolCompletion struct {
	ArgsLen    int
	ConfigFile *string
}

func (c *GatewayProtocolCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	protocols := []string{"TCP", "UDP", "TLS", "HTTP", "HTTPS"}
	if c.ArgsLen < 0 || len(args) == c.ArgsLen {

		var results []string
		for _, protocol := range protocols {
			if strings.HasPrefix(protocol, toComplete) {
				results = append(results, protocol)
			}
		}
		return results, cobra.ShellCompDirectiveNoFileComp

	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

type GatewayNameCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

func (c *GatewayNameCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var results []string
	if appID, s, err := utils.LoadAppSettings(*c.ConfigFile, *c.AppID); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
		if gateways, _, err := gateways.List(s.Client, appID, -1); err == nil {
			for _, gateway := range gateways {
				if strings.HasPrefix(gateway.Name, toComplete) {
					results = append(results, gateway.Name)
				}
			}
		}
	}
	return results, cobra.ShellCompDirectiveNoFileComp
}

type HealthTypeCompletion struct {
	ArgsLen    int
	ConfigFile *string
}

func (c *HealthTypeCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	healthTypes := []string{"startupProbe", "livenessProbe", "readinessProbe"}
	if c.ArgsLen < 0 || len(args) == c.ArgsLen {

		var results []string
		for _, healthType := range healthTypes {
			if strings.HasPrefix(healthType, toComplete) {
				results = append(results, healthType)
			}
		}
		return results, cobra.ShellCompDirectiveNoFileComp

	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

type ProbeTypeCompletion struct {
	ArgsLen    int
	ConfigFile *string
}

func (c *ProbeTypeCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	probeTypes := []string{"httpGet", "exec", "tcpSocket"}
	if c.ArgsLen < 0 || len(args) == c.ArgsLen {

		var results []string
		for _, probeType := range probeTypes {
			if strings.HasPrefix(probeType, toComplete) {
				results = append(results, probeType)
			}
		}
		return results, cobra.ShellCompDirectiveNoFileComp

	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

type HealthChecksCompletion struct {
	ArgsLen    int
	ConfigFile *string
}

func (c *HealthChecksCompletion) CompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		healthTypeCompletion := HealthTypeCompletion{ArgsLen: 0, ConfigFile: c.ConfigFile}
		return healthTypeCompletion.CompletionFunc(cmd, args, toComplete)
	} else {
		probeTypeCompletion := ProbeTypeCompletion{ArgsLen: 1, ConfigFile: c.ConfigFile}
		return probeTypeCompletion.CompletionFunc(cmd, args, toComplete)
	}
}

type ServiceProtocolCompletion struct {
	ArgsLen    int
	ConfigFile *string
}

func (c *ServiceProtocolCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	protocols := []string{"TCP", "UDP", "TLS", "SCTP"}
	if c.ArgsLen < 0 || len(args) == c.ArgsLen {

		var results []string
		for _, protocol := range protocols {
			if strings.HasPrefix(protocol, toComplete) {
				results = append(results, protocol)
			}
		}
		return results, cobra.ShellCompDirectiveNoFileComp

	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

type KeyCompletion struct {
	ArgsLen    int
	ConfigFile *string
}

func (c *KeyCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if s, err := settings.Load(*c.ConfigFile); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
		if keys, _, err := keys.List(s.Client, -1); err == nil {
			var results []string
			for _, key := range keys {
				if strings.HasPrefix(key.ID, toComplete) {
					results = append(results, key.ID)
				}
			}
			return results, cobra.ShellCompDirectiveNoFileComp
		}
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

type LabelCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

func (c *LabelCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if appID, s, err := utils.LoadAppSettings(*c.ConfigFile, *c.AppID); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
		if appsettings, err := appsettings.List(s.Client, appID); err == nil {
			var results []string
			for key := range appsettings.Label {
				if strings.HasPrefix(key, toComplete) {
					if !slices.Contains(args, key) {
						results = append(results, key)
					}
				}
			}
			return results, cobra.ShellCompDirectiveNoFileComp
		}
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

type LimitSpecCompletion struct {
	ArgsLen    int
	ConfigFile *string
}

func (c *LimitSpecCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if s, err := settings.Load(*c.ConfigFile); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
		if specs, _, err := limits.Specs(s.Client, toComplete, -1); err == nil {
			var results []string
			for _, sepc := range specs {
				results = append(results, sepc.ID)
			}
			return results, cobra.ShellCompDirectiveNoFileComp
		}
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

type LimitSetPlanCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

func (c *LimitSetPlanCompletion) CompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if strings.Contains(toComplete, "=") {
		var results []string
		parts := strings.Split(toComplete, "=")
		if s, err := settings.Load(*c.ConfigFile); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
			if plans, _, err := limits.Plans(s.Client, "", 0, 0, -1); err == nil {
				for _, plan := range plans {
					if strings.HasPrefix(plan.ID, parts[1]) {
						results = append(results, fmt.Sprintf("%s=%s", parts[0], plan.ID))
					}
				}
			}
		}
		return results, cobra.ShellCompDirectiveNoFileComp
	} else {
		ptsSetArgsCompletion := PtsSetArgsCompletion{
			PtsCompletion: &PtsCompletion{AppID: c.AppID, ArgsLen: -1, ConfigFile: c.ConfigFile},
		}
		return ptsSetArgsCompletion.CompletionFunc(cmd, args, toComplete)
	}
}

type UserPermsCompletion struct {
	ArgsLen    int
	ConfigFile *string
}

func (c *UserPermsCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	perms := []string{"view", "change", "delete"}
	if len(args) > 0 && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {

		var results []string
		for _, perm := range perms {
			if strings.HasPrefix(perm, toComplete) {
				results = append(results, perm)
			}
		}
		return results, cobra.ShellCompDirectiveNoFileComp

	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

type UserPermsArgsCompletion struct {
	*UserPermsCompletion
}

func (c *UserPermsArgsCompletion) CompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	c.ArgsLen = -1
	var results []string
	completes, _ := c.UserPermsCompletion.CompletionFunc(cmd, args, toComplete)
	for _, complete := range completes {
		if !slices.Contains(args, complete) {
			results = append(results, complete)
		}
	}
	return results, cobra.ShellCompDirectiveNoFileComp
}

type PermUsernameCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

func (c *PermUsernameCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		if appID, s, err := utils.LoadAppSettings(*c.ConfigFile, *c.AppID); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
			if perms, _, err := perms.List(s.Client, appID, -1); err == nil {
				var results []string
				for _, perm := range perms {
					if strings.HasPrefix(perm.Username, toComplete) {
						results = append(results, perm.Username)
					}
				}
				return results, cobra.ShellCompDirectiveNoFileComp
			}
		}
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

type PermUpdateCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

func (c *PermUpdateCompletion) CompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		permUsernameCompletion := PermUsernameCompletion{AppID: c.AppID, ArgsLen: 0, ConfigFile: c.ConfigFile}
		return permUsernameCompletion.CompletionFunc(cmd, args, toComplete)
	} else {
		userPermsArgsCompletion := UserPermsArgsCompletion{
			UserPermsCompletion: &UserPermsCompletion{ConfigFile: c.ConfigFile},
		}
		return userPermsArgsCompletion.CompletionFunc(cmd, args, toComplete)
	}
}

type PsCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

func (c *PsCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if appID, s, err := utils.LoadAppSettings(*c.ConfigFile, *c.AppID); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
		if pods, _, err := ps.List(s.Client, appID, -1); err == nil {
			var results []string
			for _, pod := range pods {
				if strings.HasPrefix(pod.Name, toComplete) {
					results = append(results, pod.Name)
				}
			}
			return results, cobra.ShellCompDirectiveNoFileComp
		}
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

type PtsCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

func (c *PtsCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var results []string
	if appID, s, err := utils.LoadAppSettings(*c.ConfigFile, *c.AppID); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
		if ptypes, _, err := pts.List(s.Client, appID, -1); err == nil {
			for _, ptype := range ptypes {
				if strings.HasPrefix(ptype.Name, toComplete) {
					if !slices.Contains(args, ptype.Name) {
						results = append(results, ptype.Name)
					}
				}
			}
		}
	}
	return results, cobra.ShellCompDirectiveNoFileComp
}

type PtsArgsCompletion struct {
	*PtsCompletion
}

func (c *PtsArgsCompletion) CompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var results []string
	completes, _ := c.PtsCompletion.CompletionFunc(cmd, args, toComplete)
	for _, complete := range completes {
		if !slices.Contains(args, complete) {
			results = append(results, complete)
		}
	}
	return results, cobra.ShellCompDirectiveNoFileComp
}

type PtsSetArgsCompletion struct {
	*PtsCompletion
}

func (c *PtsSetArgsCompletion) CompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var results []string
	completion := PtsCompletion{AppID: c.AppID, ArgsLen: -1, ConfigFile: c.ConfigFile}
	completes, _ := completion.CompletionFunc(cmd, args, toComplete)
	for _, complete := range completes {
		hasArgs := false
		result := fmt.Sprintf("%s=", complete)
		for _, arg := range args {
			if strings.HasPrefix(arg, result) {
				hasArgs = true
				break
			}
		}
		if !hasArgs {
			results = append(results, result)
		}
	}
	return results, cobra.ShellCompDirectiveNoSpace
}

type ResourceServiceCompletion struct {
	ArgsLen    int
	ConfigFile *string
}

func (c *ResourceServiceCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if s, err := settings.Load(*c.ConfigFile); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
		if services, _, err := resources.Services(s.Client, -1); err == nil {
			var results []string
			for _, service := range services {
				if strings.HasPrefix(service.Name, toComplete) {
					results = append(results, service.Name)
				}
			}
			return results, cobra.ShellCompDirectiveNoFileComp
		}
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

type ResourcePlanCompletion struct {
	Service    string
	ArgsLen    int
	ConfigFile *string
}

func (c *ResourcePlanCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if s, err := settings.Load(*c.ConfigFile); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
		var results []string
		if services, _, err := resources.Plans(s.Client, c.Service, -1); err == nil {
			for _, service := range services {
				results = append(results, service.Name)
			}
		}
		return results, cobra.ShellCompDirectiveNoFileComp
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

type ResourceCreateCompletion struct {
	ArgsLen    int
	ConfigFile *string
}

func (c *ResourceCreateCompletion) CompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 1 {
		resourceServiceCompletion := ResourceServiceCompletion{ArgsLen: -1, ConfigFile: c.ConfigFile}
		return resourceServiceCompletion.CompletionFunc(cmd, args, toComplete)
	} else if len(args) == 2 {
		resourcePlanCompletion := ResourcePlanCompletion{Service: args[1], ArgsLen: -1, ConfigFile: c.ConfigFile}
		return resourcePlanCompletion.CompletionFunc(cmd, args, toComplete)
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

type ResourceCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

func (c *ResourceCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if appID, s, err := utils.LoadAppSettings(*c.ConfigFile, *c.AppID); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
		if resources, _, err := resources.List(s.Client, appID, -1); err == nil {
			var results []string
			for _, resource := range resources {
				if strings.HasPrefix(resource.Name, toComplete) {
					results = append(results, resource.Name)
				}
			}
			return results, cobra.ShellCompDirectiveNoFileComp
		}
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

type ResourceUpdateCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

func (c *ResourceUpdateCompletion) CompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		ResourceCompletion := ResourceCompletion{AppID: c.AppID, ArgsLen: -1, ConfigFile: c.ConfigFile}
		return ResourceCompletion.CompletionFunc(cmd, args, toComplete)
	} else if len(args) == 1 {
		resourceServiceCompletion := ResourceServiceCompletion{ArgsLen: -1, ConfigFile: c.ConfigFile}
		return resourceServiceCompletion.CompletionFunc(cmd, args, toComplete)
	} else if len(args) == 2 {
		resourcePlanCompletion := ResourcePlanCompletion{Service: args[1], ArgsLen: -1, ConfigFile: c.ConfigFile}
		return resourcePlanCompletion.CompletionFunc(cmd, args, toComplete)
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

type ReleaseCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

func (c *ReleaseCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if appID, s, err := utils.LoadAppSettings(*c.ConfigFile, *c.AppID); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
		if releases, _, err := releases.List(s.Client, appID, "", -1); err == nil {
			var results []string
			toComplete = strings.TrimPrefix(toComplete, "v")
			for _, release := range releases {
				version := strconv.Itoa(release.Version)
				if strings.HasPrefix(version, toComplete) {
					results = append(results, version)
				}
			}
			return results, cobra.ShellCompDirectiveNoFileComp
		}
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

type RouteKindCompletion struct {
	AppID   *string
	ArgsLen int
}

func (c *RouteKindCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	kinds := []string{"HTTPRoute", "TCPRoute", "UDPRoute", "GRPCRoute", "TLSRoute"}
	if c.ArgsLen < 0 || len(args) == c.ArgsLen {

		var results []string
		for _, kind := range kinds {
			if strings.HasPrefix(kind, toComplete) {
				results = append(results, kind)
			}
		}
		return results, cobra.ShellCompDirectiveNoFileComp

	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

type RouteCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

func (c *RouteCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if appID, s, err := utils.LoadAppSettings(*c.ConfigFile, *c.AppID); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
		if routes, _, err := routes.List(s.Client, appID, -1); err == nil {
			var results []string
			for _, route := range routes {
				if strings.HasPrefix(route.Name, toComplete) {
					results = append(results, route.Name)
				}
			}
			return results, cobra.ShellCompDirectiveNoFileComp
		}
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

type ServiceCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

func (c *ServiceCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if appID, s, err := utils.LoadAppSettings(*c.ConfigFile, *c.AppID); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
		if services, err := services.List(s.Client, appID); err == nil {
			var results []string
			for _, service := range services {
				if strings.HasPrefix(service.Ptype, toComplete) {
					results = append(results, service.Ptype)
				}
			}
			return results, cobra.ShellCompDirectiveNoFileComp
		}
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

type TagCompletion struct {
	AppID      *string
	Ptype      *string
	ArgsLen    int
	ConfigFile *string
}

func (c *TagCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if appID, s, err := utils.LoadAppSettings(*c.ConfigFile, *c.AppID); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
		if config, err := config.List(s.Client, appID, -1); err == nil {
			var results []string
			for tag := range config.Tags[*c.Ptype] {
				if strings.HasPrefix(tag, toComplete) {
					if !slices.Contains(args, tag) {
						results = append(results, tag)
					}
				}
			}
			return results, cobra.ShellCompDirectiveNoFileComp
		}
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

type TlsActionCompletion struct {
	ArgsLen    int
	ConfigFile *string
}

func (c *TlsActionCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	actions := []string{"enable", "disable"}
	if c.ArgsLen < 0 || len(args) == c.ArgsLen {

		var results []string
		for _, action := range actions {
			if strings.HasPrefix(action, toComplete) {
				results = append(results, action)
			}
		}
		return results, cobra.ShellCompDirectiveNoFileComp

	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

type TokenCompletion struct {
	ArgsLen    int
	ConfigFile *string
}

func (c *TokenCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if s, err := settings.Load(*c.ConfigFile); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
		if tokens, _, err := tokens.List(s.Client, -1); err == nil {
			var results []string
			for _, token := range tokens {
				if strings.HasPrefix(token.UUID, toComplete) {
					results = append(results, token.UUID)
				}
			}
			return results, cobra.ShellCompDirectiveNoFileComp
		}
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

type VolumeCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

func (c *VolumeCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if appID, s, err := utils.LoadAppSettings(*c.ConfigFile, *c.AppID); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
		if volumes, _, err := volumes.List(s.Client, appID, -1); err == nil {
			var results []string
			for _, volume := range volumes {
				if strings.HasPrefix(volume.Name, toComplete) {
					results = append(results, volume.Name)
				}
			}
			return results, cobra.ShellCompDirectiveNoFileComp
		}
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

type VolumeTypeCompletion struct {
	ArgsLen    int
	ConfigFile *string
}

func (c *VolumeTypeCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	vts := []string{"csi", "nfs", "oss"}
	if c.ArgsLen < 0 || len(args) == c.ArgsLen {

		var results []string
		for _, vt := range vts {
			if strings.HasPrefix(vt, toComplete) {
				results = append(results, vt)
			}
		}
		return results, cobra.ShellCompDirectiveNoFileComp

	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

type VolumesMountCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

func (c *VolumesMountCompletion) CompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		volumeCompletion := VolumeCompletion{AppID: c.AppID, ArgsLen: 0, ConfigFile: c.ConfigFile}
		return volumeCompletion.CompletionFunc(cmd, args, toComplete)
	} else {
		ptsSetArgsCompletion := PtsSetArgsCompletion{
			PtsCompletion: &PtsCompletion{AppID: c.AppID, ArgsLen: -1, ConfigFile: c.ConfigFile},
		}
		return ptsSetArgsCompletion.CompletionFunc(cmd, args, toComplete)
	}
}

type VolumesUnmountCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

func (c *VolumesUnmountCompletion) CompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		volumeCompletion := VolumeCompletion{AppID: c.AppID, ArgsLen: 0, ConfigFile: c.ConfigFile}
		return volumeCompletion.CompletionFunc(cmd, args, toComplete)
	} else {
		ptsArgCompletion := PtsArgsCompletion{
			PtsCompletion: &PtsCompletion{AppID: c.AppID, ArgsLen: -1, ConfigFile: c.ConfigFile},
		}
		return ptsArgCompletion.CompletionFunc(cmd, args, toComplete)
	}
}

type VolumesCmdCompletion struct {
	ArgsLen    int
	ConfigFile *string
}

func (c *VolumesCmdCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	cmds := []string{"ls", "cp", "rm"}
	if c.ArgsLen < 0 || len(args) == c.ArgsLen {

		var results []string
		for _, cmd := range cmds {
			if strings.HasPrefix(cmd, toComplete) {
				results = append(results, cmd)
			}
		}
		return results, cobra.ShellCompDirectiveNoFileComp

	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}
