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
func (d FakeDryccCmd) RoutesCreate(string, string, string, string, int) error {
	return errors.New("routes:add")
}

func (d FakeDryccCmd) RoutesList(string, int) error {
	return errors.New("routes:list")
}

func (d FakeDryccCmd) RoutesGet(string, string) error {
	return errors.New("routes:get")
}

func (d FakeDryccCmd) RoutesSet(string, string, string) error {
	return errors.New("routes:Set")
}

func (d FakeDryccCmd) RoutesAttach(string, string, int, string) error {
	return errors.New("routes:attach")
}

func (d FakeDryccCmd) RoutesDetach(string, string, int, string) error {
	return errors.New("routes:detach")
}

func (d FakeDryccCmd) RoutesRemove(string, string) error {
	return errors.New("routes:remove")
}

func TestRoutes(t *testing.T) {
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
			args:     []string{"routes:add", "example", "--ptype=web", "--kind=http", "--port=80"},
			expected: "",
		},
		{
			args:     []string{"routes:list"},
			expected: "",
		},
		{
			args:     []string{"routes:get", "example"},
			expected: "",
		}, {
			args:     []string{"routes:set", "example", "--rules-file=rules,json"},
			expected: "",
		},
		{
			args:     []string{"routes:attach", "example", "--port=80", "--gateway=example"},
			expected: "",
		},
		{
			args:     []string{"routes:detach", "example", "--port=80", "--gateway=example"},
			expected: "",
		},
		{
			args:     []string{"routes:remove", "example"},
			expected: "",
		},
		{
			args:     []string{"routes"},
			expected: "routes:list",
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
		err = Routes(c.args, cmdr)
		assert.Error(t, errors.New(expected), err)
	}
}
