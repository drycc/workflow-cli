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

func (d FakeDryccCmd) PtsList(string, int) error {
	return errors.New("pts:list")
}

func (d FakeDryccCmd) PtsDescribe(string, string) error {
	return errors.New("pts:describe")
}

func (d FakeDryccCmd) PtsScale(string, []string) error {
	return errors.New("pts:scale")
}

func (d FakeDryccCmd) PtsRestart(string, []string, string) error {
	return errors.New("pts:restart")
}

func TestPts(t *testing.T) {
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
			args:     []string{"pts:list"},
			expected: "",
		},
		{
			args:     []string{"pts:describe", "web"},
			expected: "",
		},
		{
			args:     []string{"pts:restart", "web"},
			expected: "",
		},
		{
			args:     []string{"pts:scale", "web", "5"},
			expected: "",
		},
		{
			args:     []string{"pts"},
			expected: "pts:list",
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
		err = Pts(c.args, cmdr)
		assert.Error(t, errors.New(expected), err)
	}
}
