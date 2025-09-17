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

func TestRegistryList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/apps/enterprise/config/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"owner": "jkirk",
			"app": "enterprise",
			"values": [],
			"memory": {},
			"cpu": {},
			"tags": {},
			"registry": {
			  "web": {
				"username": "jkirk",
				"password": "ncc1701"
			}},
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})

	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.RegistryList("enterprise", "web", -1)
	assert.NoError(t, err)
	testutil.AssertOutput(t, b.String(), `PTYPE    USERNAME    PASSWORD
web      jkirk       ncc1701
`)
}

func TestRegistrySet(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/apps/foo/config/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		if r.Method == "POST" {
			testutil.AssertBody(t, api.Config{
				Registry: map[string]map[string]any{
					"web": {
						"username": "jkirk",
						"password": "ncc1701",
					},
				},
			}, r)
		}

		fmt.Fprintf(w, `{
			"owner": "jkirk",
			"app": "foo",
			"values": [],
			"memory": {},
			"cpu": {},
			"registry": {
			  "web": {
			  	"username": "jkirk",
				"password": "ncc1701"
			  }
			},
			"registry": {},
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})

	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.RegistrySet("foo", "web", "jkirk", "ncc1701")
	assert.NoError(t, err)

	testutil.AssertOutput(t, testutil.StripProgress(b.String()), `Applying registry information... done

PTYPE    USERNAME    PASSWORD
web      jkirk       ncc1701
`)
}

func TestRegistryUnset(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/apps/foo/config/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		if r.Method == "POST" {
			testutil.AssertBody(t, api.Config{
				Registry: map[string]map[string]any{
					"web": {
						"username": nil,
						"password": nil,
					},
				},
			}, r)
		}

		fmt.Fprintf(w, `{
			"owner": "jkirk",
			"app": "foo",
			"values": [],
			"memory": {},
			"cpu": {},
			"tags": {},
			"registry": {},
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})

	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.RegistryUnset("foo", "web")
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Applying registry information... done

`, "output")
}
