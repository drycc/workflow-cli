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

func (d FakeDryccCmd) CanaryInfo(string) error {
	return errors.New("canary:info")
}

func (d FakeDryccCmd) CanaryCreate(string, []string) error {
	return errors.New("canary:create")
}

func (d FakeDryccCmd) CanaryRemove(string, []string) error {
	return errors.New("canary:remove")
}

func (d FakeDryccCmd) CanaryRelease(string) error {
	return errors.New("canary:release")
}

func (d FakeDryccCmd) CanaryRollback(string) error {
	return errors.New("canary:rollback")
}

func TestCanary(t *testing.T) {
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
			args:     []string{"canary:info"},
			expected: "",
		},
		{
			args:     []string{"canary:create", "web cmd"},
			expected: "",
		},
		{
			args:     []string{"canary:remove", "web cmd"},
			expected: "",
		},
		{
			args:     []string{"canary:release"},
			expected: "",
		},
		{
			args:     []string{"canary:rollback"},
			expected: "",
		},
		{
			args:     []string{"canary"},
			expected: "canary:info",
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
		err = Canary(c.args, cmdr)
		assert.Error(t, errors.New(expected), err)
	}
}
