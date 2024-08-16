package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/drycc/workflow-cli/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestPermCodes(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/perms/codes/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
	"results": [
		{"codename": "use_app", "description": "Can use app"},
		{"codename": "use_cert", "description": "Can use cert"}
	],
	"count": 2
}`)
	})

	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.PermCodes(-1)
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `CODENAME    DESCRIPTION  
use_app     Can use app     
use_cert    Can use cert    
`, "output")
}

func TestListUserPerm(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/perms/rules/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
	"results": [
		{"id": 1, "codename": "use_app", "uniqueid": "autotest-app", "username": "foo"},
		{"id": 2, "codename": "use_cert", "uniqueid": "autotest-cert-1", "username": "foo"}
	],
	"count": 2
}`)
	})

	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.PermList("", -1)
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `ID    CODENAME    UNIQUEID           USERNAME 
1     use_app     autotest-app       foo         
2     use_cert    autotest-cert-1    foo         
`, "output")
}

func TestListUserPermLimit(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/perms/rules/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
	"results": [
		{"id": 1, "codename": "use_app", "uniqueid": "autotest-app", "username": "foo"}
	],
	"count": 1
}`)
	})

	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.PermList("use_app", 1)
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `ID    CODENAME    UNIQUEID        USERNAME 
1     use_app     autotest-app    foo         
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

	server.Mux.HandleFunc("/v2/perms/rules/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
	})

	err = cmdr.PermCreate("use_app", "lorem-ipsum", "test-user")
	assert.NoError(t, err)
	assert.Equal(t,
		testutil.StripProgress(b.String()),
		"Adding user test-user as a collaborator for use_app lorem-ipsum... done\n", "output")
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

	server.Mux.HandleFunc("/v2/perms/rules/1/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
	})

	err = cmdr.PermDelete(1)
	assert.NoError(t, err)
	assert.Equal(t,
		testutil.StripProgress(b.String()),
		"Removing user permission... done\n", "output")
}
