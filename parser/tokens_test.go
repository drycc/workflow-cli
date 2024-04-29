package parser

import (
	"bytes"
	"errors"
	"testing"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/workflow-cli/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

// Create fake implementations of each method that return the argument
// we expect to have called the function (as an error to satisfy the interface).

func (d FakeDryccCmd) TokensList(int) error {
	return errors.New("tokens:list")
}

func (d FakeDryccCmd) TokensAdd(*drycc.Client, string, string, string, string, bool) (*api.AuthTokenResponse, error) {
	return nil, errors.New("tokens:add")
}

func (d FakeDryccCmd) TokensRemove(string, string) error {
	return errors.New("tokens:remove")
}

func TestTokens(t *testing.T) {
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
			args:     []string{"tokens:list"},
			expected: "",
		},
		{
			args:     []string{"gateways:add", "alias"},
			expected: "",
		},
		{
			args:     []string{"gateways:remove", "87587952-0b78-4ecb-ac74-1d9c46197efe"},
			expected: "",
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
		err = Tokens(c.args, cmdr)
		assert.Error(t, errors.New(expected), err)
	}
}
