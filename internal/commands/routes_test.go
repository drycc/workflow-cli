package commands

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/workflow-cli/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestRoutesCreate(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/routes/", func(w http.ResponseWriter, r *http.Request) {
		request := api.RouteCreateRequest{
			Name: "example-go",
			Kind: "HTTPRoute",
			Rules: []api.RequestRouteRule{{
				BackendRefs: []api.BackendRefRequest{{Name: "example-go", Port: 443}},
			}},
		}
		testutil.AssertBody(t, request, r)
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusCreated)
		// Body isn't used by CLI, so it isn't set.
		w.Write([]byte("{}"))
	})

	err = cmdr.RoutesCreate("foo", "example-go", "HTTPRoute", api.BackendRefRequest{Name: "example-go", Port: 443})
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Adding route example-go to foo... done\n", "output")
}

func TestRoutesList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/routes/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
        {
            "app": "example-go",
            "created": "2023-04-19T00:00:00UTC",
            "owner": "test",
            "updated": "2023-04-19T00:00:00UTC",
            "name": "example-go",
            "ptype": "web",
            "kind": "HTTPRoute",
            "port": 80,
            "parent_refs": [
                {
                    "name": "example-go",
                    "port": 80
                },
                {
                    "name": "example-go",
                    "port": 8080
                }
            ],
			"rules": [{
				"backendRefs": [
					{
						"kind": "Service",
						"name": "yygl-nextcloud",
						"port": 80,
						"weight": 100
					},
					{
						"kind": "Service",
						"name": "yygl-nextcloud",
						"port": 8080,
						"weight": 100
					}
				]
			}]
        }
    ]
}`)
	})

	err = cmdr.RoutesList("foo", -1)
	assert.NoError(t, err)

	assert.Equal(t, b.String(), `NAME          OWNER    KIND         GATEWAYS                               SERVICES                                    
example-go    test     HTTPRoute    ["example-go:80","example-go:8080"]    ["yygl-nextcloud:80","yygl-nextcloud:8080"]    
`, "output")
}

func TestRoutesAttach(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/routes/example-go/attach/", func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertBody(t, api.RouteAttachRequest{Port: 4443, Gateway: "example-go"}, r)
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusCreated)
		// Body isn't used by CLI, so it isn't set.
		w.Write([]byte("{}"))
	})

	err = cmdr.RoutesAttach("foo", "example-go", 4443, "example-go")
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Attaching route example-go to gateway example-go... done\n", "output")
}

func TestRoutesDetach(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/routes/example-go/detach/", func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertBody(t, api.RouteDetachRequest{Port: 4443, Gateway: "example-go"}, r)
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusCreated)
		// Body isn't used by CLI, so it isn't set.
		w.Write([]byte("{}"))
	})

	err = cmdr.RoutesDetach("foo", "example-go", 4443, "example-go")
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Detaching route example-go to gateway example-go... done\n", "output")
}

func TestRouteGet(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/routes/example-go/rules/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `
[{
	"backendRefs": [
		{
			"group": "",
			"kind": "Service",
			"name": "example-go",
			"port": 1234,
			"weight": 1
		}
		],
		"matches": [
		{
			"path": {
			"type": "PathPrefix",
			"value": "/get"
			}
		}
		]
	}
]`)
	})

	err = cmdr.RoutesGet("foo", "example-go")
	assert.NoError(t, err)

	assert.Equal(t, b.String(), `- backendRefs:
  - group: ""
    kind: Service
    name: example-go
    port: 1234
    weight: 1
  matches:
  - path:
      type: PathPrefix
      value: /get

`, "output")
}

func TestRoutesSet(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/routes/example-go/rules/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte(""))
	})
	ruleFile, err := os.CreateTemp("", "rules.yaml")
	rules := `
- backendRefs:
  - group: ""
    kind: Service
    name: example-go
    port: 1234
    weight: 1
  matches:
  - path:
      type: PathPrefix
      value: /get`
	assert.NoError(t, err)
	defer os.Remove(ruleFile.Name())
	_, err = ruleFile.Write([]byte(rules))
	assert.NoError(t, err)
	ruleFile.Close()
	err = cmdr.RoutesSet("foo", "example-go", ruleFile.Name())
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Applying rules... done\n", "output")
}

func TestRoutesRemove(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/routes/example-go/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusNoContent)
	})

	err = cmdr.RoutesRemove("foo", "example-go")
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Removing route example-go to foo... done\n", "output")
}
