package parser

import (
	"testing"

	"github.com/drycc/workflow-cli/internal/commands"
)

func TestAppsCommandsWithAppFlag(t *testing.T) {
	t.Parallel()

	cmd := NewAppsCommand(&commands.DryccCmd{})
	for _, name := range []string{"archive", "unarchive", "quota"} {
		t.Run(name, func(t *testing.T) {
			subcmd, _, err := cmd.Find([]string{name})
			if err != nil {
				t.Fatal(err)
			}
			if subcmd.Name() != name {
				t.Fatalf("expected %s command, got %s", name, subcmd.Name())
			}
			flag := subcmd.Flags().Lookup("app")
			if flag == nil {
				t.Fatal("expected --app flag")
			}
			if flag.Shorthand != "a" {
				t.Fatalf("expected -a shorthand, got %q", flag.Shorthand)
			}
		})
	}
}
