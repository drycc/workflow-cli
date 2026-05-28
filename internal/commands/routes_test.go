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

func TestRoutesApply(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/routes/example-go/", func(w http.ResponseWriter, r *http.Request) {
		request := api.RouteUpdateRequest{
			App:        "foo",
			Name:       "example-go",
			Kind:       "HTTPRoute",
			ParentRefs: []api.RouteParentRef{},
			Rules: []api.RouteRule{{
				"backendRefs": []map[string]any{{"name": "example-go", "port": 443}},
			}},
		}
		testutil.AssertBody(t, request, r)
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	})
	routeFile, err := os.CreateTemp("", "route.yaml")
	assert.NoError(t, err)
	defer os.Remove(routeFile.Name())
	_, err = routeFile.Write([]byte("apiVersion: gateway.networking.k8s.io/v1\nkind: HTTPRoute\nmetadata:\n  name: example-go\nspec:\n  parents: []\n  rules:\n  - backends:\n    - name: example-go\n      port: 443\n"))
	assert.NoError(t, err)
	routeFile.Close()

	err = cmdr.RoutesApply("foo", routeFile.Name())
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Applying route example-go to foo... done\n", "output")
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

	assert.Equal(t, b.String(), `NAME          KIND         GATEWAYS                               SERVICES                                    
example-go    HTTPRoute    ["example-go:80","example-go:8080"]    ["yygl-nextcloud:80","yygl-nextcloud:8080"]    
`, "output")
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

	assert.Equal(t, testutil.StripProgress(b.String()), "Removing route example-go from foo... done\n", "output")
}
