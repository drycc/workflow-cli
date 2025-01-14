package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// Version displays the client version
func Update(argv []string, cmdr cmd.Commander) error {
	usage := `
Update workflow cli to latest release.

Usage: drycc update [options]

Options:
  -d --dry-run
    simulate an update, only print the version info.
`
	args, err := docopt.ParseArgs(usage, argv, "")
	if err != nil {
		return err
	}
	return cmdr.Update(safeGetBool(args, "--dry-run"))
}
