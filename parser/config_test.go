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

func (d FakeDryccCmd) ConfigInfo(string, string, string, int) error {
	return errors.New("config:info")
}

func (d FakeDryccCmd) ConfigSet(string, string, string, []string, string) error {
	return errors.New("config:set")
}

func (d FakeDryccCmd) ConfigUnset(string, string, string, []string, string) error {
	return errors.New("config:unset")
}

func (d FakeDryccCmd) ConfigPull(string, string, string, string, bool, bool) error {
	return errors.New("config:pull")
}

func (d FakeDryccCmd) ConfigPush(string, string, string, string, string) error {
	return errors.New("config:push")
}

func (d FakeDryccCmd) ConfigAttach(string, string, string) error {
	return errors.New("config:attach")
}

func (d FakeDryccCmd) ConfigDetach(string, string, string) error {
	return errors.New("config:detach")
}

func TestConfig(t *testing.T) {
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
			args:     []string{"config:info"},
			expected: "",
		},
		{
			args:     []string{"config:set", "var=value"},
			expected: "",
		},
		{
			args:     []string{"config:unset", "var"},
			expected: "",
		},
		{
			args:     []string{"config:pull"},
			expected: "",
		},
		{
			args:     []string{"config:push"},
			expected: "",
		},
		{
			args:     []string{"config:attach", "web", "g1"},
			expected: "",
		},
		{
			args:     []string{"config:detach", "web", "g1"},
			expected: "",
		},
		{
			args:     []string{"config"},
			expected: "config:info",
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
		err = Config(c.args, cmdr)
		assert.Error(t, errors.New(expected), err)
	}
}
