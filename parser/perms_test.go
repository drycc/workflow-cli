package parser

import (
	"bytes"
	"errors"
	"testing"

	"github.com/drycc/workflow-cli/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

// Create fake implementations of each method that return the argument
// we expect to have called the function (as an error to satisfy the interface).

func (d FakeDryccCmd) PermCodes(int) error {
	return errors.New("perms:codes")
}

func (d FakeDryccCmd) PermList(string, int) error {
	return errors.New("perms:list")
}

func (d FakeDryccCmd) PermCreate(string, string, string) error {
	return errors.New("perms:create")
}

func (d FakeDryccCmd) PermDelete(uint64) error {
	return errors.New("perms:delete")
}

func TestPerms(t *testing.T) {
	t.Parallel()

	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := FakeDryccCmd{WOut: &b, ConfigFile: cf}

	// cases defines the arguments and expected return of the call.
	// if expected is "", it defaults to args[0].
	cases := []struct {
		args     []string
		expected string
	}{
		{
			args:     []string{"perms:codes"},
			expected: "",
		},
		{
			args:     []string{"perms:list"},
			expected: "",
		},
		{
			args:     []string{"perms:create", "test-user", "use_app", "autotest-app"},
			expected: "",
		},
		{
			args:     []string{"perms:delete", "1"},
			expected: "",
		},
		{
			args:     []string{"perms"},
			expected: "perms:list",
		},
	}

	// For each case, check that calling the route with the arguments
	// returns the expected error, which is args[0] if not provided.
	for _, c := range cases {
		var expected string
		if c.expected == "" {
			expected = c.args[0]
		} else {
			expected = c.expected
		}
		err = Perms(c.args, cmdr)
		assert.Error(t, errors.New(expected), err)
	}
}
