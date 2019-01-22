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

func TestTLSInfo(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/numenor/tls/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
	"uuid": "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
	"app": "numenor",
	"owner": "nazgul",
	"created": "2016-08-22T17:40:16Z",
	"updated": "2016-08-22T17:40:16Z",
	"https_enforced": true
}`)
	})

	err = cmdr.TLSInfo("numenor")
	assert.NoErr(t, err)
	assert.Equal(t, b.String(), `=== numenor TLS
HTTPS Enforced: true
`, "output")
}

func TestTLSEnable(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/numenor/tls/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		b := true
		a := api.NewTLS()
		a.HTTPSEnforced = &b
		testutil.AssertBody(t, a, r)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{
	"uuid": "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
	"app": "numenor",
	"owner": "nazgul",
	"created": "2016-08-22T17:40:16Z",
	"updated": "2016-08-22T17:40:16Z",
	"https_enforced": true
}`)
	})

	err = cmdr.TLSEnable("numenor")
	assert.NoErr(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), "Enabling https-only requests for numenor... done\n", "output")
}

func TestTLSDisable(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/numenor/tls/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		testutil.AssertBody(t, api.NewTLS(), r)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{
	"uuid": "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
	"app": "numenor",
	"owner": "nazgul",
	"created": "2016-08-22T17:40:16Z",
	"updated": "2016-08-22T17:40:16Z",
	"https_enforced": false
}`)
	})

	err = cmdr.TLSDisable("numenor")
	assert.NoErr(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), "Disabling https-only requests for numenor... done\n", "output")
}
