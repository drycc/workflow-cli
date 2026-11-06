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

func TestServicesList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	assert.NoError(t, err)
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/services/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprint(w, `{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [{
        "app": "foo",
        "created": "2023-04-19T00:00:00UTC",
        "updated": "2023-04-19T00:00:00UTC",
        "domain": "foo-worker.foo.svc",
        "ptype": "worker",
        "ports": [{"name": "tcp-5000", "port": 5000, "protocol": "TCP", "target_port": 5000}]
    }]
}`)
	})

	err = cmdr.ServicesList("foo", -1)
	assert.NoError(t, err)
	assert.Contains(t, b.String(), "PTYPE")
	assert.Contains(t, b.String(), "TARGET PORT")
	assert.Contains(t, b.String(), "DOMAIN")
	assert.Contains(t, b.String(), "worker")
	assert.Contains(t, b.String(), "5000")
	assert.Contains(t, b.String(), "foo-worker.foo.svc")
	assert.NotContains(t, b.String(), "NAME")
	assert.NotContains(t, b.String(), "TARGET-PORT")
	assert.NotContains(t, b.String(), "tcp-5000")
}

func TestServicesListEmpty(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	assert.NoError(t, err)
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/services/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprint(w, `{"count": 0, "next": null, "previous": null, "results": []}`)
	})

	err = cmdr.ServicesList("foo", -1)
	assert.NoError(t, err)
	assert.Contains(t, b.String(), "No services found in foo app.")
}

func TestServicesInfo(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	assert.NoError(t, err)
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/services/web/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprint(w, `{
    "app": "foo",
    "domain": "foo.foo.svc",
    "ptype": "web",
    "ports": [{"name": "tcp-80", "port": 80, "protocol": "TCP", "target_port": 5000}]
}`)
	})

	err = cmdr.ServicesInfo("foo", "web")
	assert.NoError(t, err)
	assert.Contains(t, b.String(), `domain = "foo.foo.svc"`)
	assert.Contains(t, b.String(), `ptype = "web"`)
	assert.Contains(t, b.String(), "target_port = 5000")
	assert.NotContains(t, b.String(), "app =")
	assert.NotContains(t, b.String(), "name =")
	assert.NotContains(t, b.String(), "tcp-80")
}

func TestServicesApply(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	assert.NoError(t, err)
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/services/web/", func(w http.ResponseWriter, r *http.Request) {
		request := api.ServiceUpdateRequest{
			Ptype: "web",
			Ports: []api.ServicePort{{Port: 80, Protocol: "TCP", TargetPort: 5000}},
		}
		testutil.AssertBody(t, request, r)
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"domain":"foo.foo.svc","ptype":"web","ports":[]}`)
	})
	serviceFile, err := os.CreateTemp("", "service.toml")
	assert.NoError(t, err)
	defer os.Remove(serviceFile.Name())
	_, err = serviceFile.Write([]byte("domain = \"foo.foo.svc\"\nptype = \"web\"\n\n[[ports]]\nport = 80\nprotocol = \"TCP\"\ntarget_port = 5000\n"))
	assert.NoError(t, err)
	assert.NoError(t, serviceFile.Close())

	err = cmdr.ServicesApply("foo", serviceFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, "Applying service web to foo... done\n", testutil.StripProgress(b.String()))
}

func TestServicesApplyMissingPtype(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	assert.NoError(t, err)
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	serviceFile, err := os.CreateTemp("", "service.toml")
	assert.NoError(t, err)
	defer os.Remove(serviceFile.Name())
	_, err = serviceFile.Write([]byte("[[ports]]\nport = 80\nprotocol = \"TCP\"\ntarget_port = 5000\n"))
	assert.NoError(t, err)
	assert.NoError(t, serviceFile.Close())

	err = cmdr.ServicesApply("foo", serviceFile.Name())
	assert.EqualError(t, err, "invalid service configuration: missing ptype")
}

func TestServicesRemove(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	assert.NoError(t, err)
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/services/web/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusNoContent)
	})

	err = cmdr.ServicesRemove("foo", "web")
	assert.NoError(t, err)
	assert.Equal(t, "Removing service web from foo... done\n", testutil.StripProgress(b.String()))
}
