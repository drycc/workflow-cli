// Package completion provides types and functions for command completion in the CLI
package completion

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

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
	"github.com/drycc/workflow-cli/internal/loader"
	"github.com/drycc/workflow-cli/pkg/settings"
	"github.com/spf13/cobra"
)

// Completion is an interface for command completion functions
type Completion interface {
	CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective)
}

// AppCompletion provides completion for application names
type AppCompletion struct {
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of application names for completion
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

// CertCompletion provides completion for certificate names
type CertCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of certificate names for completion
func (c *CertCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if appID, s, err := loader.LoadAppSettings(*c.ConfigFile, *c.AppID); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
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

// CertDomainTachCompletion provides completion for certificate domain names
type CertDomainTachCompletion struct {
	AppID      *string
	ConfigFile *string
}

// CompletionFunc returns a list of certificate domain names for completion
func (c *CertDomainTachCompletion) CompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		certCompletion := CertCompletion{AppID: c.AppID, ConfigFile: c.ConfigFile}
		return certCompletion.CompletionFunc(cmd, args, toComplete)
	}
	domainCompletion := DomainCompletion{AppID: c.AppID, ArgsLen: 1, ConfigFile: c.ConfigFile}
	return domainCompletion.CompletionFunc(cmd, args, toComplete)
}

// ConfigPtsGroupArgsCompletion provides completion for config pts group arguments
type ConfigPtsGroupArgsCompletion struct {
	AppID      *string
	ConfigFile *string
}

// CompletionFunc returns a list of config pts group arguments for completion
func (c *ConfigPtsGroupArgsCompletion) CompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		ptsCompletion := PtsCompletion{AppID: c.AppID, ConfigFile: c.ConfigFile}
		return ptsCompletion.CompletionFunc(cmd, args, toComplete)
	}
	groups := args[1:]
	if appID, s, err := loader.LoadAppSettings(*c.ConfigFile, *c.AppID); err == nil {
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
	return nil, cobra.ShellCompDirectiveNoFileComp
}

// ConfigGroupCompletion provides completion for config group names
type ConfigGroupCompletion struct {
	AppID      *string
	ConfigFile *string
	ArgsLen    int
}

// CompletionFunc returns a list of config group names for completion
func (c *ConfigGroupCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if appID, s, err := loader.LoadAppSettings(*c.ConfigFile, *c.AppID); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
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

// DomainCompletion provides completion for domain names
type DomainCompletion struct {
	AppID      *string
	ConfigFile *string
	ArgsLen    int
}

// CompletionFunc returns a list of domain names for completion
func (c *DomainCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if appID, s, err := loader.LoadAppSettings(*c.ConfigFile, *c.AppID); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
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

// GatewayProtocolCompletion provides completion for gateway protocols
type GatewayProtocolCompletion struct {
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of gateway protocols for completion
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

// GatewayNameCompletion provides completion for gateway names
type GatewayNameCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of gateway names for completion
func (c *GatewayNameCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var results []string
	if appID, s, err := loader.LoadAppSettings(*c.ConfigFile, *c.AppID); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
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

// HealthTypeCompletion provides completion for health check types
type HealthTypeCompletion struct {
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of health check types for completion
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

// ProbeTypeCompletion provides completion for probe types
type ProbeTypeCompletion struct {
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of probe types for completion
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

// HealthChecksCompletion provides completion for health checks
type HealthChecksCompletion struct {
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of health checks for completion
func (c *HealthChecksCompletion) CompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		healthTypeCompletion := HealthTypeCompletion{ArgsLen: 0, ConfigFile: c.ConfigFile}
		return healthTypeCompletion.CompletionFunc(cmd, args, toComplete)
	}
	probeTypeCompletion := ProbeTypeCompletion{ArgsLen: 1, ConfigFile: c.ConfigFile}
	return probeTypeCompletion.CompletionFunc(cmd, args, toComplete)
}

// ServiceProtocolCompletion provides completion for service protocols
type ServiceProtocolCompletion struct {
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of service protocols for completion
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

// KeyCompletion provides completion for SSH key names
type KeyCompletion struct {
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of SSH key names for completion
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

// TokenCompletion provides completion for tokens
type TokenCompletion struct {
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of tokens for completion
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

// LabelCompletion provides completion for label names
type LabelCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of label names for completion
func (c *LabelCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if appID, s, err := loader.LoadAppSettings(*c.ConfigFile, *c.AppID); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
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

// LimitSpecCompletion provides completion for limit specifications
type LimitSpecCompletion struct {
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of limit specifications for completion
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

// LimitSetPlanCompletion provides completion for limit set plans
type LimitSetPlanCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of limit set plans for completion
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
	}
	ptsSetArgsCompletion := PtsSetArgsCompletion{
		PtsCompletion: &PtsCompletion{AppID: c.AppID, ArgsLen: -1, ConfigFile: c.ConfigFile},
	}
	return ptsSetArgsCompletion.CompletionFunc(cmd, args, toComplete)
}

// UserPermsCompletion provides completion for user permissions
type UserPermsCompletion struct {
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of user permissions for completion
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

// UserPermsArgsCompletion provides completion for user permission arguments.
type UserPermsArgsCompletion struct {
	*UserPermsCompletion
}

// CompletionFunc returns a list of user permission arguments for completion
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

// PermUsernameCompletion provides completion for permission usernames
type PermUsernameCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of permission usernames for completion
func (c *PermUsernameCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		if appID, s, err := loader.LoadAppSettings(*c.ConfigFile, *c.AppID); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
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

// PermUpdateCompletion provides completion for permission updates
type PermUpdateCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of permission updates for completion
func (c *PermUpdateCompletion) CompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		permUsernameCompletion := PermUsernameCompletion{AppID: c.AppID, ArgsLen: 0, ConfigFile: c.ConfigFile}
		return permUsernameCompletion.CompletionFunc(cmd, args, toComplete)
	}
	userPermsArgsCompletion := UserPermsArgsCompletion{
		UserPermsCompletion: &UserPermsCompletion{ConfigFile: c.ConfigFile},
	}
	return userPermsArgsCompletion.CompletionFunc(cmd, args, toComplete)
}

// PsCompletion provides completion for process names
type PsCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of process names for completion
func (c *PsCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if appID, s, err := loader.LoadAppSettings(*c.ConfigFile, *c.AppID); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
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

// PtsCompletion provides completion for pts names
type PtsCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of pts names for completion
func (c *PtsCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var results []string
	if appID, s, err := loader.LoadAppSettings(*c.ConfigFile, *c.AppID); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
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

// PtsArgsCompletion provides completion for process type arguments.
type PtsArgsCompletion struct {
	*PtsCompletion
}

// CompletionFunc returns a list of pts arguments for completion
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

// PtsSetArgsCompletion provides completion for process type scale arguments.
type PtsSetArgsCompletion struct {
	*PtsCompletion
}

// CompletionFunc returns a list of process type scale arguments for completion.
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

// ResourceServiceCompletion provides completion for resource services
type ResourceServiceCompletion struct {
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of resource services for completion
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

// ResourcePlanCompletion provides completion for resource plans
type ResourcePlanCompletion struct {
	Service    string
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of resource plans for completion
func (c *ResourcePlanCompletion) CompletionFunc(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
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

// ResourceCreateCompletion provides completion for resource creation
type ResourceCreateCompletion struct {
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of resource creation options for completion
func (c *ResourceCreateCompletion) CompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		resourceServiceCompletion := ResourceServiceCompletion{ArgsLen: 0, ConfigFile: c.ConfigFile}
		return resourceServiceCompletion.CompletionFunc(cmd, args, toComplete)
	}
	resourcePlanCompletion := ResourcePlanCompletion{Service: args[0], ArgsLen: 1, ConfigFile: c.ConfigFile}
	return resourcePlanCompletion.CompletionFunc(cmd, args, toComplete)
}

// ResourceCompletion provides completion for resources
type ResourceCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of resources for completion
func (c *ResourceCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if appID, s, err := loader.LoadAppSettings(*c.ConfigFile, *c.AppID); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
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

// ResourceUpdateCompletion provides completion for resource update arguments.
type ResourceUpdateCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of resource update arguments for completion.
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

// ReleaseCompletion provides completion for release versions.
type ReleaseCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of release versions for completion.
func (c *ReleaseCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if appID, s, err := loader.LoadAppSettings(*c.ConfigFile, *c.AppID); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
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

// RouteKindCompletion provides completion for route kinds.
type RouteKindCompletion struct {
	AppID   *string
	ArgsLen int
}

// CompletionFunc returns a list of route kinds for completion.
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

// RouteCompletion provides completion for routes
type RouteCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of routes for completion
func (c *RouteCompletion) CompletionFunc(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	if appID, s, err := loader.LoadAppSettings(*c.ConfigFile, *c.AppID); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
		if routes, _, err := routes.List(s.Client, appID, -1); err == nil {
			var results []string
			for _, route := range routes {
				results = append(results, route.Name)
			}
			return results, cobra.ShellCompDirectiveNoFileComp
		}
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

// ServiceCompletion provides completion for services
type ServiceCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of services for completion
func (c *ServiceCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if appID, s, err := loader.LoadAppSettings(*c.ConfigFile, *c.AppID); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
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

// TagCompletion provides completion for tags
type TagCompletion struct {
	AppID      *string
	Ptype      *string
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of tags for completion
func (c *TagCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if appID, s, err := loader.LoadAppSettings(*c.ConfigFile, *c.AppID); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
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

// TLSActionCompletion provides completion for TLS actions
type TLSActionCompletion struct {
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of TLS actions for completion
func (c *TLSActionCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	actions := []string{"redirect", "passthrough"}
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

// VolumeCompletion provides completion for volumes
type VolumeCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of volumes for completion
func (c *VolumeCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if appID, s, err := loader.LoadAppSettings(*c.ConfigFile, *c.AppID); err == nil && (c.ArgsLen < 0 || len(args) == c.ArgsLen) {
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

// VolumeTypeCompletion provides completion for volume types
type VolumeTypeCompletion struct {
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of volume types for completion
func (c *VolumeTypeCompletion) CompletionFunc(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	types := []string{"static", "dynamic"}
	if c.ArgsLen < 0 || len(args) == c.ArgsLen {
		var results []string
		for _, t := range types {
			if strings.HasPrefix(t, toComplete) {
				results = append(results, t)
			}
		}
		return results, cobra.ShellCompDirectiveNoFileComp
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

// VolumesMountCompletion provides completion for volume mounts
type VolumesMountCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of volume mount options for completion
func (c *VolumesMountCompletion) CompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		volumeCompletion := VolumeCompletion{AppID: c.AppID, ArgsLen: 0, ConfigFile: c.ConfigFile}
		return volumeCompletion.CompletionFunc(cmd, args, toComplete)
	}
	ptsArgCompletion := PtsArgsCompletion{
		PtsCompletion: &PtsCompletion{AppID: c.AppID, ArgsLen: -1, ConfigFile: c.ConfigFile},
	}
	return ptsArgCompletion.CompletionFunc(cmd, args, toComplete)
}

// VolumesUnmountCompletion provides completion for volume unmounts
type VolumesUnmountCompletion struct {
	AppID      *string
	ArgsLen    int
	ConfigFile *string
}

// CompletionFunc returns a list of volume unmount options for completion
func (c *VolumesUnmountCompletion) CompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		volumeCompletion := VolumeCompletion{AppID: c.AppID, ArgsLen: 0, ConfigFile: c.ConfigFile}
		return volumeCompletion.CompletionFunc(cmd, args, toComplete)
	}
	ptsArgCompletion := PtsArgsCompletion{
		PtsCompletion: &PtsCompletion{AppID: c.AppID, ArgsLen: -1, ConfigFile: c.ConfigFile},
	}
	return ptsArgCompletion.CompletionFunc(cmd, args, toComplete)
}
