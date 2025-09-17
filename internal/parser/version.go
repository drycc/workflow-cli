package parser

import (
	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

// NewVersionCommand creates a command for displaying the workflow CLI version.
func NewVersionCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		all bool
	}
	cmd := &cobra.Command{
		Use:   "version",
		Short: i18n.T("Display client version"),
		Long:  i18n.T("Displays the client version"),
		Run: func(_ *cobra.Command, _ []string) {
			cmdr.Version(flags.all)
		},
	}
	cmd.Flags().BoolVarP(&flags.all, "all", "a", false, i18n.T("list api and controller versions."))
	return cmd
}
