package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/arschles/assert"
	"github.com/teamhephy/controller-sdk-go/api"
	"github.com/teamhephy/workflow-cli/pkg/testutil"
)

func TestWhitelistList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/whitelist/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
    "addresses": ["1.2.3.4", "0.0.0.0/0"]
}`)
	})

	err = cmdr.WhitelistList("foo")
	assert.NoErr(t, err)

	assert.Equal(t, b.String(), "=== foo Whitelisted Addresses\n1.2.3.4\n0.0.0.0/0\n", "output")
}

func TestWhitelistAdd(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/whitelist/", func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertBody(t, api.Whitelist{Addresses: []string{"1.2.3.4", "0.0.0.0/0"}}, r)
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusCreated)
		// Body isn't used by CLI, so it isn't set.
		w.Write([]byte("{}"))
	})

	err = cmdr.WhitelistAdd("foo", "1.2.3.4,0.0.0.0/0")
	assert.NoErr(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Adding 1.2.3.4,0.0.0.0/0 to foo whitelist...\ndone\n", "output")
}

func TestWhitelistRemove(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/whitelist/", func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertBody(t, api.Whitelist{Addresses: []string{"1.2.3.4"}}, r)
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusCreated)
		// Body isn't used by CLI, so it isn't set.
		w.Write([]byte("{}"))
	})

	err = cmdr.WhitelistRemove("foo", "1.2.3.4")
	assert.NoErr(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Removing 1.2.3.4 from foo whitelist...\ndone\n", "output")
}
