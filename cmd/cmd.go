package cmd

import (
	"fmt"
	"io"
	"time"

	"github.com/drycc/controller-sdk-go/api"
)

// Commander is interface definition for running commands
type Commander interface {
	AppCreate(string, string, bool) error
	AppsList(int) error
	AppInfo(string) error
	AppOpen(string) error
	AppLogs(string, int, bool, int) error
	AppRun(string, string, []string) error
	AppDestroy(string, string) error
	AppTransfer(string, string) error
	AutoscaleList(string) error
	AutoscaleSet(string, string, int, int, int) error
	AutoscaleUnset(string, string) error
	Login(string, bool) error
	Logout() error
	Whoami(bool) error
	BuildsList(string, int) error
	BuildsCreate(string, string, string, string) error
	CanaryInfo(string) error
	CanaryCreate(string, []string) error
	CanaryRemove(string, []string) error
	CanaryRelease(string) error
	CanaryRollback(string) error
	CertsList(int, time.Time) error
	CertAdd(string, string, string) error
	CertRemove(string) error
	CertInfo(string) error
	CertAttach(string, string) error
	CertDetach(string, string) error
	ConfigList(string, string) error
	ConfigSet(string, []string) error
	ConfigUnset(string, []string) error
	ConfigPull(string, bool, bool) error
	ConfigPush(string, string) error
	DomainsList(string, int) error
	DomainsAdd(string, string) error
	DomainsRemove(string, string) error
	ServicesList(string) error
	ServicesAdd(string, string, string, string) error
	ServicesRemove(string, string, string, int) error
	GatewaysAdd(string, string, int, string) error
	GatewaysList(string, int) error
	GatewaysRemove(string, string, int, string) error
	RoutesCreate(string, string, string, string, int) error
	RoutesList(string, int) error
	RoutesGet(string, string) error
	RoutesSet(string, string, string) error
	RoutesAttach(string, string, int, string) error
	RoutesDetach(string, string, int, string) error
	RoutesRemove(string, string) error
	GitRemote(string, string, bool) error
	GitRemove(string) error
	HealthchecksList(string, string) error
	HealthchecksSet(string, string, string, *api.Healthcheck) error
	HealthchecksUnset(string, string, []string) error
	KeysList(int) error
	KeyRemove(string) error
	KeyAdd(string, string) error
	LabelsList(string) error
	LabelsSet(string, []string) error
	LabelsUnset(string, []string) error
	LimitsList(string) error
	LimitsSet(string, []string, []string) error
	LimitsUnset(string, []string, []string) error
	TimeoutsList(string) error
	TimeoutsSet(string, []string) error
	TimeoutsUnset(string, []string) error
	PermsList(string, bool, int) error
	PermCreate(string, string, bool) error
	PermDelete(string, string, bool) error
	PsList(string, int) error
	PsExec(string, string, bool, bool, []string) error
	PsScale(string, []string) error
	PsRestart(string, string) error
	RegistryList(string) error
	RegistrySet(string, []string) error
	RegistryUnset(string, []string) error
	ReleasesList(string, int) error
	ReleasesInfo(string, int) error
	ReleasesRollback(string, int) error
	RoutingInfo(string) error
	RoutingEnable(string) error
	RoutingDisable(string) error
	ShortcutsList() error
	TagsList(string) error
	TagsSet(string, []string) error
	TagsUnset(string, []string) error
	TLSInfo(string) error
	TLSForceEnable(string) error
	TLSForceDisable(string) error
	TLSAutoEnable(string) error
	TLSAutoDisable(string) error
	TLSAutoIssuer(string, string, string, string, string) error
	UsersList(results int) error
	UsersEnable(string) error
	UsersDisable(string) error
	Println(...interface{}) (int, error)
	Print(...interface{}) (int, error)
	Printf(string, ...interface{}) (int, error)
	PrintErrln(...interface{}) (int, error)
	PrintErr(...interface{}) (int, error)
	PrintErrf(string, ...interface{}) (int, error)
	Version(bool) error
	VolumesCreate(string, string, string) error
	VolumesExpand(string, string, string) error
	VolumesDelete(string, string) error
	VolumesList(string, int) error
	VolumesMount(string, string, []string) error
	VolumesUnmount(string, string, []string) error
	ResourcesServices(int) error
	ResourcesPlans(string, int) error
	ResourcesCreate(string, string, string, []string, string) error
	ResourcesList(string, int) error
	ResourceDelete(string, string) error
	ResourceGet(string, string) error
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
