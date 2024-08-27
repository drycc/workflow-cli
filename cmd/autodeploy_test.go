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

func TestAutodeployInfo(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/rivendell/settings/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"owner": "elrond",
			"app": "rivendell",
			"routable": true,
			"autoerollback": true,
			"created": "2024-01-01T00:00:00UTC",
			"updated": "2024-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})

	err = cmdr.AutodeployInfo("rivendell")
	assert.NoError(t, err)
	assert.Equal(t, b.String(), "Autodeploy is enabled.\n", "output")

	server.Mux.HandleFunc("/v2/apps/mordor/settings/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"owner": "sauron",
			"app": "mordor",
			"routable": false,
			"autodeploy": false,
			"created": "2024-01-01T00:00:00UTC",
			"updated": "2024-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})
	b.Reset()

	err = cmdr.AutodeployInfo("mordor")
	assert.NoError(t, err)
	assert.Equal(t, b.String(), "Autodeploy is disabled.\n", "output")

	// test that no autodeploy field doesn't trigger a panic
	server.Mux.HandleFunc("/v2/apps/gondor/settings/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"owner": "aragorn",
			"app": "gondor",
			"created": "2024-01-01T00:00:00UTC",
			"updated": "2024-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})
	b.Reset()

	err = cmdr.AutodeployInfo("gondor")
	assert.NoError(t, err)
	assert.Equal(t, b.String(), "Autodeploy is enabled.\n", "output")
}

func TestAutodeployEnable(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/lothlorien/settings/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		testutil.AssertBody(t, api.AppSettings{Autodeploy: api.NewAutodeploy()}, r)
		fmt.Fprintf(w, `{}`)
	})

	err = cmdr.AutodeployEnable("lothlorien")
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), "Enabling autodeploy for lothlorien... done\n", "output")
}

func TestAutodeployDisable(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/bree/settings/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		autodeploy := false
		testutil.AssertBody(t, api.AppSettings{Autodeploy: &autodeploy}, r)
		fmt.Fprintf(w, `{}`)
	})

	err = cmdr.AutodeployDisable("bree")
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), "Disabling autodeploy for bree... done\n", "output")
}
