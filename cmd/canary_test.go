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

func TestCanaryInfo(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/rivendell/settings/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"owner": "elrond",
			"app": "rivendell",
			"canaries": ["cmd"],
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})

	err = cmdr.CanaryInfo("rivendell")
	assert.NoError(t, err)
	assert.Equal(t, b.String(), "=== rivendell Canary\n\ncmd\n", "output")

	server.Mux.HandleFunc("/v2/apps/mordor/settings/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"owner": "sauron",
			"app": "mordor",
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})
	b.Reset()

	err = cmdr.CanaryInfo("mordor")
	assert.NoError(t, err)
	assert.Equal(t, b.String(), "=== mordor Canary\n\n", "output")
}

func TestCanaryCreate(t *testing.T) {
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
		data := []string{"cmd"}
		testutil.AssertBody(t, api.AppSettings{Canaries: data}, r)
		fmt.Fprintf(w, `{}`)
	})

	err = cmdr.CanaryCreate("lothlorien", []string{"cmd"})
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), "Applying canary settings for process type cmd on lothlorien... done\n", "output")
}

func TestCanaryRemove(t *testing.T) {
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
		data := []string{"cmd"}
		testutil.AssertBody(t, api.AppSettings{Canaries: data}, r)
		w.WriteHeader(http.StatusNoContent)
	})

	err = cmdr.CanaryRemove("bree", []string{"cmd"})
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), "Removing canary for process type cmd on bree... done\n", "output")
}

func TestCanaryRelease(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/bree/canary/release/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusNoContent)
	})

	err = cmdr.CanaryRelease("bree")
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), "Release canary for bree... done\n", "output")
}

func TestCanaryRollback(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/bree/canary/rollback/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusNoContent)
	})

	err = cmdr.CanaryRollback("bree")
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), "Rollback canary for bree... done\n", "output")
}
