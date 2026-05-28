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

func TestGatewaysList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/gateways/", func(w http.ResponseWriter, _ *http.Request) {
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
            "updated": "2023-04-19T00:00:00UTC",
            "ports": [
                {
                    "port": 80,
                    "protocol": "HTTP"
                },
                {
                    "port": 443,
                    "protocol": "HTTPS"
                }
            ],
            "addresses": [
                {
                    "type": "IPAddress",
                    "value": "192.168.11.1"
                }
            ]
        }
    ]
}`)
	})

	err = cmdr.GatewaysList("foo", -1)
	assert.NoError(t, err)

	assert.Equal(t, b.String(), `NAME    PORT    PROTOCOL    ADDRESSES    
foo     80      HTTP        192.168.11.1    
foo     443     HTTPS       192.168.11.1    
`, "output")
}

func TestGatewaysApply(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/gateways/example-go/", func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertBody(t, api.GatewayUpdateRequest{App: "foo", Name: "example-go", Ports: []api.GatewayPort{{Port: 443, Protocol: "HTTPS"}}}, r)
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	})

	gatewayFile, err := os.CreateTemp("", "gateway.yaml")
	assert.NoError(t, err)
	defer os.Remove(gatewayFile.Name())
	_, err = gatewayFile.Write([]byte("apiVersion: gateway.networking.k8s.io/v1\nkind: Gateway\nmetadata:\n  name: example-go\nspec:\n  ports:\n  - port: 443\n    protocol: HTTPS\n"))
	assert.NoError(t, err)
	gatewayFile.Close()

	err = cmdr.GatewaysApply("foo", gatewayFile.Name())
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Applying gateway example-go to foo... done\n", "output")
}

func TestGatewaysRemove(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/gateways/example-go/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusNoContent)
	})

	err = cmdr.GatewaysRemove("foo", "example-go")
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Removing gateway example-go from foo... done\n", "output")
}
