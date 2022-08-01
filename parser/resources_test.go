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

func (d FakeDryccCmd) ResourcesServices(int) error {
	return errors.New("resources:services")
}

func (d FakeDryccCmd) ResourcesPlans(string, int) error {
	return errors.New("resources:plans")
}

func (d FakeDryccCmd) ResourcesCreate(string, string, string, []string) error {
	return errors.New("resources:create")
}

func (d FakeDryccCmd) ResourcesList(string, int) error {
	return errors.New("resources:list")
}

func (d FakeDryccCmd) ResourceGet(string, string) error {
	return errors.New("resources:describe")
}

func (d FakeDryccCmd) ResourcePut(string, string, string, []string) error {
	return errors.New("resources:update")
}

func (d FakeDryccCmd) ResourceDelete(string, string) error {
	return errors.New("resources:destroy")
}

func (d FakeDryccCmd) ResourceBind(string, string) error {
	return errors.New("resources:bind")
}

func (d FakeDryccCmd) ResourceUnbind(string, string) error {
	return errors.New("resources:unbind")
}

func TestResources(t *testing.T) {
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
			args:     []string{"resources:services"},
			expected: "",
		},
		{
			args:     []string{"resources:plans", "mysql"},
			expected: "",
		},
		{
			args:     []string{"resources:create", "mysql:5.6", "mysql", "key1=value1"},
			expected: "",
		},
		{
			args:     []string{"resources:list"},
			expected: "",
		},
		{
			args:     []string{"resources:describe", "mysql"},
			expected: "",
		},
		{
			args:     []string{"resources:update", "mysql:5.7", "mysql", "key1=value2"},
			expected: "",
		},
		{
			args:     []string{"resources:destroy", "mysql"},
			expected: "",
		},
		{
			args:     []string{"resources:bind", "mysql"},
			expected: "",
		},
		{
			args:     []string{"resources:unbind", "mysql"},
			expected: "",
		},
		{
			args:     []string{"resources"},
			expected: "resources:list",
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
		err = Resources(c.args, cmdr)
		assert.Error(t, errors.New(expected), err)
	}
}
