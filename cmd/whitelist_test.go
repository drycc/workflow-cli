package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/arschles/assert"
	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/workflow-cli/pkg/testutil"
)

func TestAllowlistList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/allowlist/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
    "addresses": ["1.2.3.4", "0.0.0.0/0"]
}`)
	})

	err = cmdr.AllowlistList("foo")
	assert.NoErr(t, err)

	assert.Equal(t, b.String(), "=== foo Allowlisted Addresses\n1.2.3.4\n0.0.0.0/0\n", "output")
}

func TestAllowlistAdd(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/allowlist/", func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertBody(t, api.Allowlist{Addresses: []string{"1.2.3.4", "0.0.0.0/0"}}, r)
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusCreated)
		// Body isn't used by CLI, so it isn't set.
		w.Write([]byte("{}"))
	})

	err = cmdr.AllowlistAdd("foo", "1.2.3.4,0.0.0.0/0")
	assert.NoErr(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Adding 1.2.3.4,0.0.0.0/0 to foo allowlist...\ndone\n", "output")
}

func TestAllowlistRemove(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/allowlist/", func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertBody(t, api.Allowlist{Addresses: []string{"1.2.3.4"}}, r)
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusCreated)
		// Body isn't used by CLI, so it isn't set.
		w.Write([]byte("{}"))
	})

	err = cmdr.AllowlistRemove("foo", "1.2.3.4")
	assert.NoErr(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Removing 1.2.3.4 from foo allowlist...\ndone\n", "output")
}
