package commands

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/drycc/workflow-cli/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestListUserPerm(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/apps/lorem-ipsum/perms/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
	"results": [
		{"app": "lorem-ipsum", "username": "foo", "permissions": ["view"]},
		{"app": "lorem-ipsum", "username": "foo", "permissions": ["view"]}
	],
	"count": 2
}`)
	})

	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.PermList("lorem-ipsum", -1)
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `USERNAME    PERMISSIONS 
foo         view           
foo         view           
`, "output")
}

func TestListUserPermLimit(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/apps/lorem-ipsum/perms/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
	"results": [
		{"app": "lorem-ipsum", "username": "foo", "permissions": ["view"]}
	],
	"count": 1
}`)
	})

	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.PermList("lorem-ipsum", 1)
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `USERNAME    PERMISSIONS 
foo         view           
`, "output")
}

func TestCreateUserPerm(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/lorem-ipsum/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
  "uuid": "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
  "id": "lorem-ipsum",
  "owner": "dolar-sit-amet",
  "structure": {
    "cmd": 1
  },
  "created": "2016-08-22T17:40:16Z",
  "updated": "2016-08-22T17:40:16Z"
}`)
	})

	server.Mux.HandleFunc("/v2/apps/lorem-ipsum/perms/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
	})

	err = cmdr.PermCreate("lorem-ipsum", "test-user", "view")
	assert.NoError(t, err)
	assert.Equal(t,
		testutil.StripProgress(b.String()),
		"Adding user test-user as a collaborator for view... done\n", "output")
}

func TestUpdateUserPerm(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/lorem-ipsum/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
  "uuid": "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
  "id": "lorem-ipsum",
  "owner": "dolar-sit-amet",
  "structure": {
    "cmd": 1
  },
  "created": "2016-08-22T17:40:16Z",
  "updated": "2016-08-22T17:40:16Z"
}`)
	})

	server.Mux.HandleFunc("/v2/apps/lorem-ipsum/perms/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
	})

	err = cmdr.PermUpdate("lorem-ipsum", "test-user", "view")
	assert.NoError(t, err)
	assert.Equal(t,
		testutil.StripProgress(b.String()),
		"Updating user test-user as a collaborator for view... done\n", "output")
}

func TestDeleteUserPerm(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/lorem-ipsum/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
  "uuid": "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
  "id": "lorem-ipsum",
  "owner": "dolar-sit-amet",
  "structure": {
    "cmd": 1
  },
  "created": "2016-08-22T17:40:16Z",
  "updated": "2016-08-22T17:40:16Z"
}`)
	})

	server.Mux.HandleFunc("/v2/apps/lorem-ipsum/perms/test-user/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
	})

	err = cmdr.PermDelete("lorem-ipsum", "test-user")
	assert.NoError(t, err)
	assert.Equal(t,
		testutil.StripProgress(b.String()),
		"Removing user permission... done\n", "output")
}
