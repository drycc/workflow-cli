package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// Shortcuts displays all relevant shortcuts for the CLI.
func Shortcuts(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for shortcuts:

shortcuts:list       list all relevant shortcuts for the CLI

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "shortcuts:list":
		return shortcutsList(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "shortcuts" {
			argv[0] = "shortcuts:list"
			return shortcutsList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func shortcutsList(argv []string, cmdr cmd.Commander) error {
	usage := `
Lists all relevant shortcuts for the CLI.

Usage: drycc shortcuts:list
`

	_, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	return cmdr.ShortcutsList()
}
