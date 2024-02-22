package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/drycc/workflow-cli/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestPermsListUsers(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/apps/foo/perms/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)

		fmt.Fprintf(w, `{
  "users": [
    "baz",
    "bar"
  ]
}`)
	})

	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.PermsList("foo", false, -1)
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `USERNAME    ADMIN 
baz         false    
bar         false    
`, "output")
}

func TestPermsListUsersLimit(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/apps/foo/perms/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)

		fmt.Fprintf(w, `{
  "users": [
    "baz"
  ]
}`)
	})

	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.PermsList("foo", false, 1)
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `USERNAME    ADMIN 
baz         false    
`, "output")
}

func TestPermsListAdmins(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/admin/perms/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)

		fmt.Fprintf(w, `{
  "count": 2,
  "next": null,
  "previous": null,
  "results": [
    {
      "username": "fred",
      "is_superuser": true
    },
    {
      "username": "bob",
      "is_superuser": true
    }
]
}`)
	})

	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.PermsList("foo", true, -1)
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `USERNAME    ADMIN 
fred        true     
bob         true     
`, "output")
}

func TestPermsListAdminsLimit(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/admin/perms/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)

		fmt.Fprintf(w, `{
  "count": 2,
  "next": null,
  "previous": null,
  "results": [
    {
      "username": "fred",
      "is_superuser": true
    }
  ]
}`)
	})

	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.PermsList("foo", true, 1)
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `USERNAME    ADMIN 
fred        true     
`, "output")
}

func TestPermsCreateUser(t *testing.T) {
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

	server.Mux.HandleFunc("/v2/apps/lorem-ipsum/perms", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
	})

	err = cmdr.PermCreate("lorem-ipsum", "test-user", false)
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), `Adding test-user to lorem-ipsum collaborators... done
`, "output")
}

func TestPermsCreateAdmin(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/admin/perms/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
	})

	err = cmdr.PermCreate("lorem-ipsum", "test-admin", true)
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), `Adding test-admin to system administrators... done
`, "output")
}

func TestPermsDeleteUser(t *testing.T) {
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

	server.Mux.HandleFunc("/v2/apps/lorem-ipsum/perms", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
	})

	err = cmdr.PermDelete("lorem-ipsum", "test-user", false)
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), `Removing test-user from lorem-ipsum collaborators... done
`, "output")
}

func TestPermsDeleteAdmin(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/admin/perms/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
	})

	err = cmdr.PermDelete("lorem-ipsum", "test-admin", true)
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), `Removing test-admin from system administrators... done
`, "output")
}
