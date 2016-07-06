package cmd

import (
	"fmt"
	"sort"

	"github.com/deis/workflow-cli/cli"
)

// Shortcuts displays all relevant shortcuts for the CLI.
func ShortcutsList() error {
	var (
		strBuilder string = ""
		keys       []string
	)

	// NOTE(bacongobbler): go does not guarantee an iteration order when iterating over a map,
	// so to work around this we can sort the keys and iterate using the key array
	for k := range cli.Shortcuts {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		strBuilder += fmt.Sprintf("%s -> %s\n", k, cli.Shortcuts[k])
	}

	fmt.Println(strBuilder)

	return nil
}
