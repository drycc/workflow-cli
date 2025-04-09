package parser

import (
	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

// Build an image from source code
func NewVersionCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var flags struct {
		all bool
	}
	cmd := &cobra.Command{
		Use:   "version",
		Short: i18n.T("Display client version"),
		Long:  i18n.T("Displays the client version"),
		Run: func(_ *cobra.Command, args []string) {
			cmdr.Version(flags.all)
		},
	}
	cmd.Flags().BoolVarP(&flags.all, "all", "a", false, i18n.T("list api and controller versions."))
	return cmd
}
