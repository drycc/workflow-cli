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
	DryccCmd := &cobra.Command{
		Use:   "version",
		Short: i18n.T("Display client version"),
		Long:  i18n.T("Displays the client version"),
		Run: func(DryccCmd *cobra.Command, args []string) {
			cmdr.Version(flags.all)
		},
	}
	DryccCmd.Flags().BoolVarP(&flags.all, "all", "a", false, i18n.T("list api and controller versions."))
	return DryccCmd
}
