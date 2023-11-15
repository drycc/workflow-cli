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

func TestLabelsList(t *testing.T) {
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
			"owner": "jim",
			"app": "rivendell",
		    "label": {"team" : "drycc", "git_repo": "https://github.com/drycc/controller-sdk-go"},
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})

	err = cmdr.LabelsList("rivendell")
	assert.NoError(t, err)
	assert.Equal(t, b.String(), `UUID                                    OWNER    KEY         VALUE                                      
de1bf5b5-4a72-4f94-a10c-d2a3741cdf75    jim      git_repo    https://github.com/drycc/controller-sdk-go    
de1bf5b5-4a72-4f94-a10c-d2a3741cdf75    jim      team        drycc                                         
`, "output")

	server.Mux.HandleFunc("/v2/apps/mordor/settings/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"owner": "priw",
			"app": "mordor",
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})
	b.Reset()

	err = cmdr.LabelsList("mordor")
	assert.NoError(t, err)
	assert.Equal(t, b.String(), "No labels found in mordor app.\n", "output")
}

func TestListsSet(t *testing.T) {
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
		data := map[string]interface{}{
			"git_repo": "https://github.com/drycc/controller-sdk-go",
			"team":     "drycc",
		}
		testutil.AssertBody(t, api.AppSettings{Label: data}, r)
		fmt.Fprintf(w, "{}")
	})

	err = cmdr.LabelsSet("lothlorien", []string{
		"team=drycc",
		"git_repo=https://github.com/drycc/controller-sdk-go",
	})
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), "Applying labels on lothlorien... done\n", "output")
}

func TestListsUnset(t *testing.T) {
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
		testutil.AssertBody(t, api.AppSettings{Label: map[string]interface{}{
			"team":     nil,
			"git_repo": nil,
		}}, r)
		fmt.Fprintf(w, "{}")
	})

	err = cmdr.LabelsUnset("bree", []string{
		"team",
		"git_repo",
	})
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), "Removing labels on bree... done\n", "output")
}
