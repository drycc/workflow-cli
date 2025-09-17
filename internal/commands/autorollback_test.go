package commands

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/workflow-cli/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestAutorollbackInfo(t *testing.T) {
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

	err = cmdr.AutorollbackInfo("rivendell")
	assert.NoError(t, err)
	assert.Equal(t, b.String(), "Autorollback is enabled.\n", "output")

	server.Mux.HandleFunc("/v2/apps/mordor/settings/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"owner": "sauron",
			"app": "mordor",
			"routable": false,
			"autorollback": false,
			"created": "2024-01-01T00:00:00UTC",
			"updated": "2024-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})
	b.Reset()

	err = cmdr.AutorollbackInfo("mordor")
	assert.NoError(t, err)
	assert.Equal(t, b.String(), "Autorollback is disabled.\n", "output")

	// test that no autorollback field doesn't trigger a panic
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

	err = cmdr.AutorollbackInfo("gondor")
	assert.NoError(t, err)
	assert.Equal(t, b.String(), "Autorollback is enabled.\n", "output")
}

func TestAutorollbackEnable(t *testing.T) {
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
		testutil.AssertBody(t, api.AppSettings{Autorollback: api.NewAutorollback()}, r)
		fmt.Fprintf(w, `{}`)
	})

	err = cmdr.AutorollbackEnable("lothlorien")
	assert.NoError(t, err)
	testutil.AssertOutput(t, testutil.StripProgress(b.String()), "Enabling autorollback for lothlorien... done\n")
}

func TestAutorollbackDisable(t *testing.T) {
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
		autorollback := false
		testutil.AssertBody(t, api.AppSettings{Autorollback: &autorollback}, r)
		fmt.Fprintf(w, `{}`)
	})

	err = cmdr.AutorollbackDisable("bree")
	assert.NoError(t, err)
	testutil.AssertOutput(t, testutil.StripProgress(b.String()), "Disabling autorollback for bree... done\n")
}
