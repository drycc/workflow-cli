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

func TestGatewaysList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/gateways/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
        {
            "app": "foo",
            "name": "foo",
            "created": "2023-04-19T00:00:00UTC",
            "owner": "test",
            "updated": "2023-04-19T00:00:00UTC",
            "listeners": [
                {
                    "name": "foo-80-http",
                    "port": 80,
                    "protocol": "HTTP",
                    "allowedRoutes": {"namespaces": {"from": "All"}}
                },
                {
                    "name": "foo-443-https",
                    "port": 443,
                    "protocol": "HTTPS",
                    "allowedRoutes": {"namespaces": {"from": "All"}}
                }
            ]
        }
    ]
}`)
	})

	err = cmdr.GatewaysList("foo", -1)
	assert.NoError(t, err)

	assert.Equal(t, b.String(), `=== foo Gateways
+------+---------------+------+----------+
| NAME |   LISENTER    | PORT | PROTOCOL |
+------+---------------+------+----------+
| foo  | foo-80-http   |   80 | HTTP     |
+      +---------------+------+----------+
|      | foo-443-https |  443 | HTTPS    |
+------+---------------+------+----------+
`, "output")
}

func TestGatewaysAdd(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/gateways/", func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertBody(t, api.GatewayCreateRequest{Name: "example-go", Port: 443, Protocol: "HTTPS"}, r)
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusCreated)
		// Body isn't used by CLI, so it isn't set.
		w.Write([]byte("{}"))
	})

	err = cmdr.GatewaysAdd("foo", "example-go", 443, "HTTPS")
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Adding gateway example-go to foo... done\n", "output")
}

func TestGatewaysDelete(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/gateways/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusNoContent)
	})

	err = cmdr.GatewaysRemove("foo", "example-go", 443, "HTTPS")
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Removing gateway example-go to foo... done\n", "output")
}
