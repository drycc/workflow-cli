package cmd

import (
	"fmt"
	"io"
	"time"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// Commander is interface definition for running commands
type Commander interface {
	AppCreate(string, string, bool) error
	AppsList(int) error
	AppInfo(string) error
	AppOpen(string) error
	AppLogs(string, int, bool, int) error
	AppRun(string, string, []string, uint32, uint32) error
	AppDestroy(string, string) error
	AppTransfer(string, string) error
	AutodeployInfo(string) error
	AutodeployEnable(string) error
	AutodeployDisable(string) error
	AutorollbackInfo(string) error
	AutorollbackEnable(string) error
	AutorollbackDisable(string) error
	AutoscaleList(string) error
	AutoscaleSet(string, string, int, int, int) error
	AutoscaleUnset(string, string) error
	Login(string, bool, string, string) error
	Logout() error
	Whoami(bool) error
	TokensList(int) error
	TokensAdd(*drycc.Client, string, string, string, string, bool) (*api.AuthTokenResponse, error)
	TokensRemove(string, string) error
	BuildsInfo(string, int) error
	BuildsCreate(string, string, string, string, string, string) error
	CertsList(string, int) error
	CertAdd(string, string, string, string) error
	CertRemove(string, string) error
	CertInfo(string, string) error
	CertAttach(string, string, string) error
	CertDetach(string, string, string) error
	ConfigInfo(string, string, string, int) error
	ConfigSet(string, string, string, []string, string) error
	ConfigUnset(string, string, string, []string, string) error
	ConfigPull(string, string, string, string, bool, bool) error
	ConfigPush(string, string, string, string, string) error
	ConfigAttach(string, string, string) error
	ConfigDetach(string, string, string) error
	DomainsList(string, int) error
	DomainsAdd(string, string, string) error
	DomainsRemove(string, string) error
	ServicesList(string) error
	ServicesAdd(string, string, string, string) error
	ServicesRemove(string, string, string, int) error
	GatewaysAdd(string, string, int, string) error
	GatewaysList(string, int) error
	GatewaysRemove(string, string, int, string) error
	RoutesCreate(string, string, string, ...api.BackendRefRequest) error
	RoutesList(string, int) error
	RoutesGet(string, string) error
	RoutesSet(string, string, string) error
	RoutesAttach(string, string, int, string) error
	RoutesDetach(string, string, int, string) error
	RoutesRemove(string, string) error
	GitRemote(string, string, bool) error
	GitRemove(string) error
	HealthchecksList(string, string, int) error
	HealthchecksSet(string, string, string, *api.Healthcheck) error
	HealthchecksUnset(string, string, []string) error
	KeysList(int) error
	KeyRemove(string) error
	KeyAdd(string, string) error
	LabelsList(string) error
	LabelsSet(string, []string) error
	LabelsUnset(string, []string) error
	LimitsList(string, int) error
	LimitsSet(string, []string) error
	LimitsUnset(string, []string) error
	LimitsSpecs(string, int) error
	LimitsPlans(string, int, int, int) error
	TimeoutsList(string, int) error
	TimeoutsSet(string, []string) error
	TimeoutsUnset(string, []string) error
	PermList(string, int) error
	PermCreate(string, string, string) error
	PermUpdate(string, string, string) error
	PermDelete(string, string) error
	PsList(string, int) error
	PsLogs(string, string, int, bool, string) error
	PsExec(string, string, bool, bool, []string) error
	PsDescribe(string, string) error
	PsDelete(string, []string) error
	PtsList(string, int) error
	PtsDescribe(string, string) error
	PtsScale(string, []string) error
	PtsRestart(string, []string, string) error
	PtsClean(string, []string) error
	RegistryList(string, string, int) error
	RegistrySet(string, string, string, string) error
	RegistryUnset(string, string) error
	ReleasesList(string, string, int) error
	ReleasesInfo(string, int) error
	ReleasesDeploy(string, []string, bool, string) error
	ReleasesRollback(string, []string, int) error
	RoutingInfo(string) error
	RoutingEnable(string) error
	RoutingDisable(string) error
	ShortcutsList() error
	TagsList(string, string, int) error
	TagsSet(string, string, []string) error
	TagsUnset(string, string, []string) error
	TLSInfo(string) error
	TLSForceEnable(string) error
	TLSForceDisable(string) error
	TLSAutoEnable(string) error
	TLSAutoDisable(string) error
	TLSAutoIssuer(string, string, string, string, string) error
	Update(bool) error
	UsersList(int) error
	UsersEnable(string) error
	UsersDisable(string) error
	Println(...interface{}) (int, error)
	Print(...interface{}) (int, error)
	Printf(string, ...interface{}) (int, error)
	PrintErrln(...interface{}) (int, error)
	PrintErr(...interface{}) (int, error)
	PrintErrf(string, ...interface{}) (int, error)
	Version(bool) error
	VolumesCreate(string, string, string, string, map[string]interface{}) error
	VolumesExpand(string, string, string) error
	VolumesDelete(string, string) error
	VolumesList(string, int) error
	VolumesInfo(string, string) error
	VolumesClient(string, string, ...string) error
	VolumesMount(string, string, []string) error
	VolumesUnmount(string, string, []string) error
	ResourcesServices(int) error
	ResourcesPlans(string, int) error
	ResourcesCreate(string, string, string, []string, string) error
	ResourcesList(string, int) error
	ResourceDelete(string, string, string) error
	ResourceGet(string, string, bool) error
	ResourcePut(string, string, string, []string, string) error
	ResourceBind(string, string) error
	ResourceUnbind(string, string) error
}

// DryccCmd is an implementation of Commander.
type DryccCmd struct {
	ConfigFile string
	Warned     bool
	WOut       io.Writer
	WErr       io.Writer
	WIn        io.Reader
	Location   *time.Location
}

// Println prints a line to an output writer.
func (d *DryccCmd) Println(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(d.WOut, a...)
}

// Print prints a line to an output writer.
func (d *DryccCmd) Print(a ...interface{}) (n int, err error) {
	return fmt.Fprint(d.WOut, a...)
}

// Printf prints a line to an error writer.
func (d *DryccCmd) Printf(s string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(d.WOut, s, a...)
}

// PrintErrln prints a line to an error writer.
func (d *DryccCmd) PrintErrln(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(d.WErr, a...)
}

// PrintErr prints a line to an error writer.
func (d *DryccCmd) PrintErr(a ...interface{}) (n int, err error) {
	return fmt.Fprint(d.WErr, a...)
}

// PrintErrf prints a line to an error writer.
func (d *DryccCmd) PrintErrf(s string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(d.WErr, s, a...)
}
