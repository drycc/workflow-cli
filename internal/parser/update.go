package parser

import (
	"github.com/drycc/workflow-cli/internal/commands"
	"github.com/drycc/workflow-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

func NewUpdateCommand(cmdr *commands.DryccCmd) *cobra.Command {
	var dryRun bool
	var cmd = &cobra.Command{
		Use:   "update",
		Args:  cobra.NoArgs,
		Short: i18n.T("Update workflow cli to latest release"),
		Long:  i18n.T("Update workflow cli to latest release"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmdr.Update(dryRun)
		},
	}
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, i18n.T("Simulate an update, only print the version info."))
	return cmd
}
