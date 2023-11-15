package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// Version displays the client version
func Version(argv []string, cmdr cmd.Commander) error {
	usage := `
Displays the client version.

Usage: drycc version [options]

Options:
  -a --all
    list api and controller versions
`
	args, err := docopt.ParseArgs(usage, argv, "")
	if err != nil {
		return err
	}

	return cmdr.Version(args["--all"].(bool))
}
