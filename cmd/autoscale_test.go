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

func TestAutoscaleList(t *testing.T) {
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
			"autoscale": {"cmd": {"min": 3, "max": 8, "cpu_percent": 40}},
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})

	err = cmdr.AutoscaleList("rivendell")
	assert.NoError(t, err)
	assert.Equal(t, b.String(), `PTYPE    PERCENT    MIN    MAX 
cmd      40         3      8      
`, "output")

	server.Mux.HandleFunc("/v2/apps/mordor/settings/", func(w http.ResponseWriter, _ *http.Request) {
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

	err = cmdr.AutoscaleList("mordor")
	assert.NoError(t, err)
	assert.Equal(t, b.String(), "No autoscale rules found.\n", "output")
}

func TestAutoscaleSet(t *testing.T) {
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
		data := map[string]*api.Autoscale{
			"cmd": {
				Min:        3,
				Max:        8,
				CPUPercent: 40,
			},
		}
		testutil.AssertBody(t, api.AppSettings{Autoscale: data}, r)
		fmt.Fprintf(w, `{}`)
	})

	err = cmdr.AutoscaleSet("lothlorien", "cmd", 3, 8, 40)
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), "Applying autoscale settings for process type cmd on lothlorien... done\n", "output")
}

func TestAutoscaleUnset(t *testing.T) {
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
		testutil.AssertBody(t, api.AppSettings{Autoscale: map[string]*api.Autoscale{"cmd": nil}}, r)
		fmt.Fprintf(w, `{"autoscale":{"cmd":null}}`)
	})

	err = cmdr.AutoscaleUnset("bree", "cmd")
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), "Removing autoscale for process type cmd on bree... done\n", "output")
}
