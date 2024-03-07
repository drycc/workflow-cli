package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/workflow-cli/pkg/testutil"
	"github.com/stretchr/testify/assert"
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

	server.Mux.HandleFunc("/v2/apps/numenor/tls/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
	"uuid": "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
	"app": "numenor",
	"owner": "nazgul",
	"created": "2016-08-22T17:40:16Z",
	"updated": "2016-08-22T17:40:16Z",
	"https_enforced": true,
	"certs_auto_enabled": true,
	"issuer": {
		"server": "https://acme.zerossl.com/v2/DV90",
		"email": "drycc@drycc.cc",
		"key_id": "AA",
		"key_secret": "BB"
	}
}`)
	})

	err = cmdr.TLSInfo("numenor")
	assert.NoError(t, err)
	assert.Equal(t, b.String(), `UUID                                    OWNER     CERTS-AUTO    HTTPS-ENFORCED    EMAIL             SERVER                           
c4aed81c-d1ca-4ff1-ab89-d2151264e1a3    nazgul    true          true              drycc@drycc.cc    https://acme.zerossl.com/v2/DV90    
`, "output")
}

func TestTLSForceEnable(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/numenor/tls/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		b := true
		a := api.NewTLS()
		a.HTTPSEnforced = &b
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

	err = cmdr.TLSForceEnable("numenor")
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), "Enabling https-only requests for numenor... done\n", "output")
}

func TestTLSForceDisable(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/numenor/tls/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		b := false
		a := api.NewTLS()
		a.HTTPSEnforced = &b
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

	err = cmdr.TLSForceDisable("numenor")
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), "Disabling https-only requests for numenor... done\n", "output")
}

func TestTLSAutoIssuer(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/numenor/tls/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		issuer := api.Issuer{
			Email:     "drycc@drycc.cc",
			Server:    "https://acme.zerossl.com/v2/DV90",
			KeyID:     "keyID",
			KeySecret: "keySecret",
		}
		a := api.NewTLS()
		a.Issuer = &issuer
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{
	"uuid": "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
	"app": "numenor",
	"owner": "nazgul",
	"created": "2016-08-22T17:40:16Z",
	"updated": "2016-08-22T17:40:16Z",
	"https_enforced": false,
	"certs_auto_enabled": false,
	"issuer": {
		"email":"anonymous@cert-manager.io",
		"server":"https://acme-v02.api.letsencrypt.org/directory",
		"key_id":"keyID",
		"key_secret":"keySecret"
	}
}`)
	})

	err = cmdr.TLSAutoIssuer("numenor", "drycc@drycc.cc", "https://acme.zerossl.com/v2/DV90", "keyID", "keySecret")
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), "Adding issuer requests for numenor... done\n", "output")
}
