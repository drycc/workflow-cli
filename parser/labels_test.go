package parser

import (
	"bytes"
	"errors"
	"testing"

	"github.com/arschles/assert"
	"github.com/teamhephy/workflow-cli/pkg/testutil"
)

// Create fake implementations of each method that return the argument
// we expect to have called the function (as an error to satisfy the interface).

func (d FakeDeisCmd) LabelsList(string) error {
	return errors.New("labels:list")
}

func (d FakeDeisCmd) LabelsSet(string, []string) error {
	return errors.New("labels:set")
}

func (d FakeDeisCmd) LabelsUnset(string, []string) error {
	return errors.New("labels:unset")
}

func TestLabels(t *testing.T) {
	t.Parallel()

	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := FakeDeisCmd{WOut: &b, ConfigFile: cf}

	// cases defines the arguments and expected return of the call.
	// if expected is "", it defaults to args[0].
	cases := []struct {
		args     []string
		expected string
	}{
		{
			args:     []string{"labels:list"},
			expected: "",
		},
		{
			args:     []string{"labels:set", "git_repo=https://github.com/teamhephy/workflow", "team=deis"},
			expected: "",
		},
		{
			args:     []string{"labels:unset", "git_repo", "team"},
			expected: "",
		},
		{
			args:     []string{"labels"},
			expected: "labels:list",
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
		err = Labels(c.args, cmdr)
		assert.Err(t, errors.New(expected), err)
	}
}
