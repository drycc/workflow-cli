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
		testutil.AssertBody(t, api.RouteCreateRequest{Name: "example-go", Type: "web", Port: 443, Kind: "HTTPRoute"}, r)
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusCreated)
		// Body isn't used by CLI, so it isn't set.
		w.Write([]byte("{}"))
	})

	err = cmdr.RoutesCreate("foo", "example-go", "web", "HTTPRoute", 443)
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

	server.Mux.HandleFunc("/v2/apps/foo/routes/", func(w http.ResponseWriter, r *http.Request) {
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
			"procfile_type": "web",
			"kind": "HTTPRoute",
			"parent_refs": [
                {
                    "name": "example-go",
                    "sectionName": "example-go-80-http"
                }
            ]
        }
    ]
}`)
	})

	err = cmdr.RoutesList("foo", -1)
	assert.NoError(t, err)

	assert.Equal(t, b.String(), `=== foo Routes
+------------+------+-----------+------------+
|    NAME    | TYPE |   KIND    |  GATEWAY   |
+------------+------+-----------+------------+
| example-go | web  | HTTPRoute | example-go |
+------------+------+-----------+------------+
`, "output")
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

	server.Mux.HandleFunc("/v2/apps/foo/routes/example-go/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `[
  {
    "backendRefs": [
      {
        "kind": "Service",
        "name": "py3django3",
        "port": 80
      }
    ]
  }
]`)
	})

	err = cmdr.RoutesGet("foo", "example-go")
	assert.NoError(t, err)

	assert.Equal(t, b.String(), `[
  {
    "backendRefs": [
      {
        "kind": "Service",
        "name": "py3django3",
        "port": 80
      }
    ]
  }
]
`, "output")
}

func TestRouteSet(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/routes/example-go/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusNoContent)
	})
	rules := `"[{\"backendRefs\": [{\"kind\": \"Service\",\"name\": \"py3django3\",\"port\": 80}]}]"`

	err = cmdr.RoutesSet("foo", "example-go", rules)
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Applying rules... done\n", "output")
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
		testutil.AssertBody(t, api.RouteAttackRequest{Port: 4443, Gateway: "example-go"}, r)
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
		testutil.AssertBody(t, api.RouteDetackRequest{Port: 4443, Gateway: "example-go"}, r)
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusCreated)
		// Body isn't used by CLI, so it isn't set.
		w.Write([]byte("{}"))
	})

	err = cmdr.RoutesDetach("foo", "example-go", 4443, "example-go")
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Detaching route example-go to gateway example-go... done\n", "output")
}

func TestRoutesGet(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/routes/example-go/rules/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusOK)
		// TODO  real rule
		w.Write([]byte(""))
	})

	err = cmdr.RoutesGet("foo", "example-go")
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "\n", "output")
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

	server.Mux.HandleFunc("/v2/apps/foo/routes/example-go/rules/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte(""))
	})

	err = cmdr.RoutesSet("foo", "example-go", "")
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

	server.Mux.HandleFunc("/v2/apps/foo/routes/example-go/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusNoContent)
	})

	err = cmdr.RoutesRemove("foo", "example-go")
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Removing route example-go to foo... done\n", "output")
}
