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

func (d FakeDryccCmd) VolumesList(string, int) error {
	return errors.New("volumes:list")
}

func (d FakeDryccCmd) VolumesInfo(string, string) error {
	return errors.New("volumes:info")
}

func (d FakeDryccCmd) VolumesCreate(string, string, string, string, map[string]interface{}) error {
	return errors.New("volumes:add")
}

func (d FakeDryccCmd) VolumesExpand(string, string, string) error {
	return errors.New("volumes:expand")
}

func (d FakeDryccCmd) VolumesDelete(string, string) error {
	return errors.New("volumes:remove")
}

func (d FakeDryccCmd) VolumesClient(string, string, ...string) error {
	return errors.New("volumes:client")
}

func (d FakeDryccCmd) VolumesMount(string, string, []string) error {
	return errors.New("volumes:mount")
}

func (d FakeDryccCmd) VolumesUnmount(string, string, []string) error {
	return errors.New("volumes:unmount")
}

func TestVolumes(t *testing.T) {
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
			args:     []string{"volumes:add", "myvolume", "500G"},
			expected: "",
		},
		{
			args:     []string{"volumes:expand", "myvolume", "500G"},
			expected: "",
		},
		{
			args:     []string{"volumes:list"},
			expected: "",
		},
		{
			args:     []string{"volumes:remove", "myvolume"},
			expected: "",
		},
		{
			args:     []string{"volumes:client", "ls", "--", "vol://myvolume"},
			expected: "",
		},
		{
			args:     []string{"volumes:mount", "myvolume", "cmd=data/cmd1"},
			expected: "",
		},
		{
			args:     []string{"volumes:unmount", "myvolume", "cmd"},
			expected: "",
		},
		{
			args:     []string{"volumes"},
			expected: "volumes:list",
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
		err = Volumes(c.args, cmdr)
		assert.Error(t, errors.New(expected), err)
	}
}
