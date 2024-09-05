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

func TestDomainsList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/domains/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
    "count": 2,
    "next": null,
    "previous": null,
    "results": [
        {
            "app": "foo",
            "created": "2014-01-01T00:00:00UTC",
            "domain": "example.example.com",
			"ptype": "web",
            "owner": "test",
            "updated": "2014-01-01T00:00:00UTC"
        },
        {
            "app": "foo",
            "created": "2014-01-01T00:00:00UTC",
            "domain": "foo",
			"ptype": "web",
            "owner": "test",
            "updated": "2014-01-01T00:00:00UTC"
        }
    ]
}`)
	})

	err = cmdr.DomainsList("foo", -1)
	assert.NoError(t, err)

	assert.Equal(t, b.String(), `APP    OWNER    PTYPE    CREATED                   UPDATED                   DOMAIN              
foo    test     web      2014-01-01T00:00:00UTC    2014-01-01T00:00:00UTC    example.example.com    
foo    test     web      2014-01-01T00:00:00UTC    2014-01-01T00:00:00UTC    foo                    
`, "output")
}

func TestDomainsListLimit(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/domains/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
    "count": 2,
    "next": null,
    "previous": null,
    "results": [
        {
            "app": "foo",
            "created": "2014-01-01T00:00:00UTC",
            "domain": "example.example.com",
			"ptype": "web",
            "owner": "test",
            "updated": "2014-01-01T00:00:00UTC"
        }
    ]
}`)
	})

	err = cmdr.DomainsList("foo", 1)
	assert.NoError(t, err)

	assert.Equal(t, b.String(), `APP    OWNER    PTYPE    CREATED                   UPDATED                   DOMAIN              
foo    test     web      2014-01-01T00:00:00UTC    2014-01-01T00:00:00UTC    example.example.com    
`, "output")
}

func TestDomainsAdd(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/domains/", func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertBody(t, api.DomainCreateRequest{Domain: "example.example.com", Ptype: "web"}, r)
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusCreated)
		// Body isn't used by CLI, so it isn't set.
		w.Write([]byte("{}"))
	})

	err = cmdr.DomainsAdd("foo", "example.example.com", "web")
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Adding example.example.com to foo... done\n", "output")
}

func TestDomainsDelete(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/domains/example.example.com", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusNoContent)
	})

	err = cmdr.DomainsRemove("foo", "example.example.com")
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Removing example.example.com from foo... done\n", "output")
}
