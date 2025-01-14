package cmd

import (
	"testing"
)

func TestShortcutsList(t *testing.T) {
	t.Parallel()

	expected := `create -> apps:create
destroy -> apps:destroy
info -> apps:info
login -> auth:login
logout -> auth:logout
logs -> apps:logs
open -> apps:open
pull -> builds:create
rollback -> releases:rollback
run -> apps:run
scale -> pts:scale
whoami -> auth:whoami
`
	actual := sortShortcuts()
	if actual != expected {
		t.Errorf("Expected %s, Got %s", expected, actual)
	}

}
