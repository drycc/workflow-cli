package parser

import (
	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/internal/template"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

// Shortcuts defines the interface for creating shortcut commands.
type Shortcuts interface {
	Create(cmdr *commands.DryccCmd) []*cobra.Command
}

// SupportedShortcuts is a list of all supported shortcut types.
var SupportedShortcuts = []Shortcuts{
	&AuthShortcuts{}, &AppsShortcuts{}, &PsShortcuts{}, &BuildsShortcuts{}, &PtsShortcuts{},
}

// AuthShortcuts provides authentication-related shortcut commands.
type AuthShortcuts struct{}

// AppsShortcuts provides application-related shortcut commands.
type AppsShortcuts struct{}

// PsShortcuts provides process-related shortcut commands.
type PsShortcuts struct{}

// Create returns a slice of authentication-related shortcut commands.
func (a *AuthShortcuts) Create(cmdr *commands.DryccCmd) []*cobra.Command {
	login := authLogin(cmdr)
	login.Example = "drycc login http://drycc.local3.dryccapp.com/"

	logout := authLogout(cmdr)
	logout.Example = "drycc auth logout"

	whoami := authWhoami(cmdr)
	whoami.Example = "drycc whoami"

	return []*cobra.Command{login, logout, whoami}
}

// Create returns a slice of application-related shortcut commands.
func (a *AppsShortcuts) Create(cmdr *commands.DryccCmd) []*cobra.Command {
	destroy := appsDestroy(cmdr)
	destroy.Example = "drycc destroy -a <app> --confirm <app>"

	run := appsRun(cmdr)
	run.Example = template.CustomExample(
		"drycc run --mount=myvolume:/data -- 'echo hello'",
		map[string]string{
			"<volume>":  i18n.T("The volume name"),
			"<path>":    i18n.T("The filesystem path"),
			"<command>": i18n.T("The shell command to run inside the container"),
		},
	)
	return []*cobra.Command{appsCreate(cmdr), destroy, appsInfo(cmdr), appsOpen(cmdr), run}
}

// Create returns a slice of process-related shortcut commands.
func (p *PsShortcuts) Create(cmdr *commands.DryccCmd) []*cobra.Command {
	exec := psExecCommand(cmdr)
	exec.Example = template.CustomExample(
		"drycc exec my-pod -it -- bash",
		map[string]string{
			"<pod>": i18n.T("The pod name for the application"),
		},
	)

	logs := psLogsCommand(cmdr)
	logs.Example = "drycc logs my-pod"
	return []*cobra.Command{exec, logs}
}

// BuildsShortcuts provides build-related shortcut commands.
type BuildsShortcuts struct{}

// Create returns a slice of build-related shortcut commands.
func (b *BuildsShortcuts) Create(cmdr *commands.DryccCmd) []*cobra.Command {
	create := buildsCreate(cmdr)
	create.Use = `pull <image>`
	return []*cobra.Command{create}
}

// PtsShortcuts provides process type-related shortcut commands.
type PtsShortcuts struct{}

// Create returns a slice of process type-related shortcut commands.
func (p *PtsShortcuts) Create(cmdr *commands.DryccCmd) []*cobra.Command {
	scale := ptsScaleCommand(cmdr)
	scale.Example = template.CustomExample(
		"drycc scale web=3",
		map[string]string{
			"<ptype>": i18n.T("The process name as defined in your Procfile"),
			"<num>":   i18n.T("The number of processes"),
		},
	)
	return []*cobra.Command{scale}
}
