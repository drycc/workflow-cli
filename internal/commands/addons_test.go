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

func TestAddonsList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/addons/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
        {
            "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
            "name": "my-addon",
            "app": "foo",
            "plan": "micro",
            "kind": "Valkey",
            "multiplier": 1,
            "parameters": {},
            "created": "2025-07-20T00:00:00UTC",
            "updated": "2025-07-20T00:00:00UTC"
        }
    ]
}`)
	})

	err = cmdr.AddonsList("foo", -1)
	assert.NoError(t, err)

	assert.Equal(t, "NAME        KIND      PLAN  \nmy-addon    Valkey    micro    \n", b.String())
}

func TestAddonsListEmpty(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/addons/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{"count": 0, "next": null, "previous": null, "results": []}`)
	})

	err = cmdr.AddonsList("foo", -1)
	assert.NoError(t, err)
	assert.Equal(t, b.String(), "No addons found in foo app.\n", "output")
}

func TestAddonsInfo(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/addons/my-addon/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
    "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
    "name": "my-addon",
    "app": "foo",
    "plan": "micro",
    "kind": "Valkey",
    "multiplier": 1,
    "parameters": {},
    "created": "2025-07-20T00:00:00UTC",
    "updated": "2025-07-20T00:00:00UTC"
}`)
	})

	err = cmdr.AddonsInfo("foo", "my-addon")
	assert.NoError(t, err)
	assert.Contains(t, b.String(), "kind: Valkey")
	assert.Contains(t, b.String(), "name: my-addon")
	assert.Contains(t, b.String(), "plan: micro")
}

func TestAddonsConnection(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/addons/my-addon/connection/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{"host":"localhost","port":"6379"}`)
	})

	err = cmdr.AddonsConnection("foo", "my-addon")
	assert.NoError(t, err)
	assert.Contains(t, b.String(), "KEY")
	assert.Contains(t, b.String(), "VALUE")
	assert.Contains(t, b.String(), "host")
	assert.Contains(t, b.String(), "localhost")
	assert.Contains(t, b.String(), "port")
	assert.Contains(t, b.String(), "6379")
}

func TestAddonsApply(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/addons/my-addon/", func(w http.ResponseWriter, r *http.Request) {
		request := api.AddonInstanceUpsertRequest{
			Kind:       "Valkey",
			Plan:       "micro",
			Parameters: map[string]any{},
		}
		testutil.AssertBody(t, request, r)
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	})

	addonFile, err := os.CreateTemp("", "addon.yaml")
	assert.NoError(t, err)
	defer os.Remove(addonFile.Name())
	_, err = addonFile.Write([]byte("apiVersion: controller.drycc.cc/v2\nkind: Valkey\nmetadata:\n  name: my-addon\nspec:\n  plan: micro\n  parameters: {}\n"))
	assert.NoError(t, err)
	addonFile.Close()

	err = cmdr.AddonsApply("foo", addonFile.Name())
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Applying addon my-addon to foo... done\n", "output")
}

func TestAddonsApplyMissingName(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	addonFile, err := os.CreateTemp("", "addon.yaml")
	assert.NoError(t, err)
	defer os.Remove(addonFile.Name())
	_, err = addonFile.Write([]byte("kind: Valkey\nspec:\n  plan: micro\n"))
	assert.NoError(t, err)
	addonFile.Close()

	err = cmdr.AddonsApply("foo", addonFile.Name())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing metadata.name")
}

func TestAddonsRemove(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/addons/my-addon/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusNoContent)
	})

	err = cmdr.AddonsRemove("foo", "my-addon")
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Removing addon my-addon from foo... done\n", "output")
}

func TestAddonsClassesList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/addon-classes/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
    "count": 1,
    "results": [
        {
            "name": "valkey",
            "description": "Valkey in-memory data store",
            "kind": "Valkey",
            "storage_model": "bundle",
            "plans": []
        }
    ]
}`)
	})

	err = cmdr.AddonsClassesList(-1)
	assert.NoError(t, err)

	assert.Equal(t, "NAME      KIND      DESCRIPTION                 \nvalkey    Valkey    Valkey in-memory data store    \n", b.String())
}

func TestAddonsClassesListEmpty(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/addon-classes/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{"count": 0, "results": []}`)
	})

	err = cmdr.AddonsClassesList(-1)
	assert.NoError(t, err)
	assert.Equal(t, b.String(), "No addon classes found.\n", "output")
}

func TestAddonsClassesInfo(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/addon-classes/valkey/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
    "name": "valkey",
    "description": "Valkey in-memory data store",
    "kind": "Valkey",
    "storage_model": "bundle",
    "plans": [
        {
            "name": "micro",
            "description": "1 GB, Standalone"
        }
    ]
}`)
	})

	err = cmdr.AddonsClassesInfo("valkey")
	assert.NoError(t, err)
	assert.Contains(t, b.String(), "kind: Valkey")
	assert.Contains(t, b.String(), "name: valkey")
	assert.Contains(t, b.String(), "storageModel: bundle")
	assert.Contains(t, b.String(), "micro")
}
