package parser

import (
	"bytes"
	"errors"
	"testing"

	"github.com/drycc/workflow-cli/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func (d FakeDryccCmd) BuildsList(string, int) error {
	return errors.New("builds:list")
}

func (d FakeDryccCmd) BuildsCreate(string, string, string, string) error {
	return errors.New("builds:create")
}

func TestBuilds(t *testing.T) {
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
			args:     []string{"builds:list"},
			expected: "",
		},
		{
			args:     []string{"builds:create", "drycc/example-go:latest"},
			expected: "",
		},
		{
			args:     []string{"builds"},
			expected: "builds:list",
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
		err = Builds(c.args, cmdr)
		assert.Error(t, errors.New(expected), err)
	}
}
