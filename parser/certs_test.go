package parser

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/drycc/workflow-cli/pkg/testutil"
)

// Create fake implementations of each method that return the argument
// we expect to have called the function (as an error to satisfy the interface).

func (d FakeDryccCmd) CertsList(int, time.Time) error {
	return errors.New("certs:list")
}

func (d FakeDryccCmd) CertAdd(string, string, string) error {
	return errors.New("certs:add")
}

func (d FakeDryccCmd) CertRemove(string) error {
	return errors.New("certs:remove")
}

func (d FakeDryccCmd) CertInfo(string) error {
	return errors.New("certs:info")
}

func (d FakeDryccCmd) CertAttach(string, string) error {
	return errors.New("certs:attach")
}

func (d FakeDryccCmd) CertDetach(string, string) error {
	return errors.New("certs:detach")
}

func TestCerts(t *testing.T) {
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
			args:     []string{"certs:list"},
			expected: "",
		},
		{
			args:     []string{"certs:add", "name", "cert", "key"},
			expected: "",
		},
		{
			args:     []string{"certs:remove", "name"},
			expected: "",
		},
		{
			args:     []string{"certs:info", "name"},
			expected: "",
		},
		{
			args:     []string{"certs:attach", "name", "example.com"},
			expected: "",
		},
		{
			args:     []string{"certs:detach", "name", "example.com"},
			expected: "",
		},
		{
			args:     []string{"certs"},
			expected: "certs:list",
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
		err = Certs(c.args, cmdr)
		assert.Err(t, errors.New(expected), err)
	}
}
