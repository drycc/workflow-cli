package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/controller-sdk-go/api"
	"github.com/deis/workflow-cli/pkg/testutil"
	"strings"
)

func TestLabelsList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/rivendell/settings/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"owner": "jim",
			"app": "rivendell",
		    "label": {"team" : "deis", "git_repo": "https://github.com/deis/controller-sdk-go"},
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})

	err = cmdr.LabelsList("rivendell")
	assert.NoErr(t, err)
	assert.Equal(t, strings.TrimSpace(b.String()), `=== rivendell Label
git_repo:      https://github.com/deis/controller-sdk-go
team:          deis`, "output")

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
	assert.NoErr(t, err)
	assert.Equal(t, b.String(), "=== mordor Label\nNo labels found.\n", "output")
}

func TestListsSet(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/lothlorien/settings/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		data := map[string]interface{}{
			"git_repo": "https://github.com/deis/controller-sdk-go",
			"team":     "deis",
		}
		testutil.AssertBody(t, api.AppSettings{Label: data}, r)
		fmt.Fprintf(w, "{}")
	})

	err = cmdr.LabelsSet("lothlorien", []string{
		"team=deis",
		"git_repo=https://github.com/deis/controller-sdk-go",
	})
	assert.NoErr(t, err)
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
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

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
	assert.NoErr(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), "Removing labels on bree... done\n", "output")
}
