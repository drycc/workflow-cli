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

func (d FakeDryccCmd) GatewaysList(string, int) error {
	return errors.New("gateways:list")
}

func (d FakeDryccCmd) GatewaysAdd(string, string, int, string) error {
	return errors.New("gateways:add")
}

func (d FakeDryccCmd) GatewaysRemove(string, string, int, string) error {
	return errors.New("gateways:remove")
}

func TestGateways(t *testing.T) {
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
			args:     []string{"gateways:add", "example", "--port=80", "--protocol=http"},
			expected: "",
		},
		{
			args:     []string{"gateways:list"},
			expected: "",
		},
		{
			args:     []string{"gateways:remove", "example", "--port=80", "--protocol=http"},
			expected: "",
		},
		{
			args:     []string{"gateways"},
			expected: "gateways:list",
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
		err = Gateways(c.args, cmdr)
		assert.Error(t, errors.New(expected), err)
	}
}
