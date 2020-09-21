package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/arschles/assert"
	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/workflow-cli/pkg/testutil"
)

func TestVolumesList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/example-go/volumes/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
		{
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
			"owner": "test",
			"app": "example-go",
			"name": "myvolume",
			"size": "500M",
			"path": {"cmd": "/data/cmd1", "cmd123": "/data/cmd123"},
			"created": "2020-08-26T00:00:00UTC",
			"updated": "2020-08-26T00:00:00UTC"
		}
    ]
}`)
	})

	err = cmdr.VolumesList("example-go", -1)
	assert.NoErr(t, err)

	assert.Equal(t, b.String(), `=== example-go volumes
--- myvolume     500M
cmd              /data/cmd1
cmd123           /data/cmd123
`, "output")
}

func TestVolumesCreate(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/example-go/volumes/", func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertBody(t, api.Volume{Name: "myvolume", Size: "500M"}, r)
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusCreated)
		// Body isn't used by CLI, so it isn't set.
		w.Write([]byte("{}"))
	})

	err = cmdr.VolumesCreate("example-go", "myvolume", "500M")
	assert.NoErr(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Creating myvolume to example-go... done\n", "output")
}

func TestVolumesDelete(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/example-go/volumes/myvolume/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusNoContent)
	})

	err = cmdr.VolumesDelete("example-go", "myvolume")
	assert.NoErr(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Deleting myvolume from example-go... done\n", "output")
}

func TestVolumesMount(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/apps/example-go/volumes/myvolume/path/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		if r.Method == "PATCH" {
			testutil.AssertBody(t, api.Volume{
				Path: map[string]interface{}{
					"cmd": "/data/cmd1",
				},
			}, r)
		}
		fmt.Fprintf(w, `{
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
			"owner": "test",
			"app": "example-go",
			"name": "myvolume",
			"size": "500M",
			"path": {"cmd": "/data/cmd1"},
			"created": "2020-08-26T00:00:00UTC",
			"updated": "2020-08-26T00:00:00UTC"
}`)
	})

	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.VolumesMount("example-go", "myvolume", []string{"cmd=/data/cmd1"})
	assert.NoErr(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Mounting volume... done
`, "output")
}

func TestVolumesUnmount(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/apps/example-go/volumes/myvolume/path/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		if r.Method == "PATCH" {
			testutil.AssertBody(t, api.Volume{
				Path: map[string]interface{}{
					"cmd": nil,
				},
			}, r)
		}
		fmt.Fprintf(w, `{
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
			"owner": "test",
			"app": "example-go",
			"name": "myvolume",
			"size": "500M",
			"path": {},
			"created": "2020-08-26T00:00:00UTC",
			"updated": "2020-08-26T00:00:00UTC"
}`)
	})

	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.VolumesUnmount("example-go", "myvolume", []string{"cmd"})
	assert.NoErr(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Unmounting volume... done
`, "output")
}
