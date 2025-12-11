package commands

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/drycc/workflow-cli/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestWorkspacesSwitch(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/workspaces/my-workspace/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"name": "my-workspace",
			"email": "test@example.com",
			"created": "2016-08-22T17:40:16Z",
			"updated": "2016-08-22T17:40:16Z"
		}`)
	})

	err = cmdr.WorkspacesSwitch("my-workspace")
	assert.NoError(t, err)
	assert.Contains(t, b.String(), "Switched to workspace my-workspace")
}

func TestWorkspacesSwitchNotFound(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/workspaces/nonexistent/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	err = cmdr.WorkspacesSwitch("nonexistent")
	assert.Error(t, err)
}
