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

func (d FakeDryccCmd) ReleasesList(string, string, int) error {
	return errors.New("releases:list")
}

func (d FakeDryccCmd) ReleasesInfo(string, int) error {
	return errors.New("releases:info")
}

func (d FakeDryccCmd) ReleasesDeploy(string, []string, bool, string) error {
	return errors.New("releases:deploy")
}

func (d FakeDryccCmd) ReleasesRollback(string, []string, int) error {
	return errors.New("releases:rollback")
}

func TestReleases(t *testing.T) {
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
			args:     []string{"releases:list"},
			expected: "",
		},
		{
			args:     []string{"releases:info", "v1"},
			expected: "",
		},
		{
			args:     []string{"releases:deploy", "web"},
			expected: "",
		},
		{
			args:     []string{"releases:rollback"},
			expected: "",
		},
		{
			args:     []string{"releases:rollback", "v1"},
			expected: "",
		},
		{
			args:     []string{"releases"},
			expected: "releases:list",
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
		err = Releases(c.args, cmdr)
		assert.Error(t, errors.New(expected), err)
	}
}
