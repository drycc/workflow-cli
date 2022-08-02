package parser

import (
	"bytes"
	"errors"
	"testing"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/workflow-cli/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

// Create fake implementations of each method that return the argument
// we expect to have called the function (as an error to satisfy the interface).

func (d FakeDryccCmd) HealthchecksList(string, string) error {
	return errors.New("healthchecks:list")
}

func (d FakeDryccCmd) HealthchecksSet(string, string, string, *api.Healthcheck) error {
	return errors.New("healthchecks:set")
}

func (d FakeDryccCmd) HealthchecksUnset(string, string, []string) error {
	return errors.New("healthchecks:unset")
}

func TestHealthchecks(t *testing.T) {
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
			args:     []string{"healthchecks:list"},
			expected: "",
		},
		{
			args:     []string{"healthchecks:set", "liveness", "httpGet", "80"},
			expected: "",
		},
		{
			args:     []string{"healthchecks:set", "liveness", "httpGet", "80", "--headers=test-header:test-value"},
			expected: "",
		},
		{
			args:     []string{"healthchecks:set", "liveness", "exec", "ls"},
			expected: "",
		},
		{
			args:     []string{"healthchecks:set", "liveness", "tcpSocket", "80"},
			expected: "",
		},
		{
			args:     []string{"healthchecks:unset", "liveness"},
			expected: "",
		},
		{
			args:     []string{"healthchecks"},
			expected: "healthchecks:list",
		},
		{
			args:     []string{"healthchecks:set", "alien", "httpGet", "80"},
			expected: "probe type alien is invalid. Must be one of [liveness readiness]",
		},
		{
			args:     []string{"healthchecks:unset", "alien", "httpGet", "80"},
			expected: "probe type alien is invalid. Must be one of [liveness readiness]",
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
		err = Healthchecks(c.args, cmdr)
		assert.Error(t, errors.New(expected), err)
	}
}
