package commands

import (
	"bytes"
	"fmt"
	"net/http"
	"syscall"
	"testing"
	"time"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/workflow-cli/pkg/testutil"
	"github.com/stretchr/testify/assert"
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

	server.Mux.HandleFunc("/v2/apps/example-go/volumes/", func(w http.ResponseWriter, _ *http.Request) {
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
			"size": "500G",
			"path": {"cmd": "/data/cmd1", "cmd123": "/data/cmd123"},
			"type": "csi",
			"parameters": {},
			"created": "2020-08-26T00:00:00UTC",
			"updated": "2020-08-26T00:00:00UTC"
		}
    ]
}`)
	})

	err = cmdr.VolumesList("example-go", -1)
	assert.NoError(t, err)
	testutil.AssertOutput(t, b.String(), `NAME        OWNER    TYPE    PTYPE     PATH            SIZE
myvolume    test     csi     cmd       /data/cmd1      500G
myvolume    test     csi     cmd123    /data/cmd123    500G
`)
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
		testutil.AssertBody(t, api.Volume{Name: "myvolume", Size: "500G", Type: "csi"}, r)
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusCreated)
		// Body isn't used by CLI, so it isn't set.
		w.Write([]byte("{}"))
	})

	err = cmdr.VolumesCreate("example-go", "myvolume", "csi", "500G", map[string]any{})
	assert.NoError(t, err)
	err = cmdr.VolumesCreate("example-go", "myvolume", "csi", "500K", map[string]any{})
	expected := `500K doesn't fit format #unit
Examples: 2G 2g`
	assert.Equal(t, err.Error(), expected, "output")

	assert.Equal(t, testutil.StripProgress(b.String()), "Creating myvolume to example-go... done\n", "output")
}

func TestVolumesInfo(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/example-go/volumes/myvolume/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
	"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	"owner": "test",
	"app": "example-go",
	"name": "myvolume",
	"size": "0G",
	"path": {"cmd": "/data/cmd1", "cmd123": "/data/cmd123"},
	"type": "nfs",
	"parameters": {
		"nfs": {
			"server": "nfs.drycc.cc",
			"path": "/mnt",
			"readOnly": true
		}
	},
	"created": "2020-08-26T00:00:00UTC",
	"updated": "2020-08-26T00:00:00UTC"
}`)
	})

	err = cmdr.VolumesInfo("example-go", "myvolume")
	assert.NoError(t, err)
	testutil.AssertOutput(t, b.String(), `UUID:          de1bf5b5-4a72-4f94-a10c-d2a3741cdf75
Name:          myvolume
Owner:         test
Type:          nfs
Path:
               cmd: /data/cmd1
               cmd123: /data/cmd123

Parameters:
               nfs:
                 path: /mnt
                 readOnly: true
                 server: nfs.drycc.cc

Created:       2020-08-26T00:00:00UTC
Updated:       2020-08-26T00:00:00UTC
`)
}

func TestVolumesServe(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}
	server.Mux.HandleFunc("/v2/apps/example-go/volumes/myvolume/filer/_/ping", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `pong`)
	})
	server.Mux.HandleFunc("/v2/apps/example-go/volumes/myvolume/filer/_/bind", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{"endpoint": "/v2/apps/example-go/volumes/myvolume/filer/webdav", "username": "user", "password": "pass"}`)
	})

	// Use a channel to signal when the method has been called
	done := make(chan error, 1)
	go func() {
		done <- cmdr.VolumesServe("example-go", "myvolume")
	}()

	// Give the goroutine time to start and produce output
	time.Sleep(500 * time.Millisecond)

	// Send interrupt signal to stop the blocking method
	go func() {
		time.Sleep(100 * time.Millisecond)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()

	// Wait for completion with timeout
	select {
	case err := <-done:
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("VolumesServe did not complete within timeout")
	}

	output := b.String()
	assert.Contains(t, output, "WebDAV service for volume myvolume is running.")
	assert.Contains(t, output, "Endpoint:")
	assert.Contains(t, output, "Username:")
	assert.Contains(t, output, "Password:")
}

func TestVolumesExpand(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/example-go/volumes/myvolume/", func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertBody(t, api.Volume{Name: "myvolume", Size: "500G"}, r)
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusCreated)
		// Body isn't used by CLI, so it isn't set.
		w.Write([]byte("{}"))
	})

	err = cmdr.VolumesExpand("example-go", "myvolume", "500G")
	assert.NoError(t, err)
	err = cmdr.VolumesExpand("example-go", "myvolume", "500K")
	expected := `500K doesn't fit format #unit
Examples: 2G 2g`
	assert.Equal(t, err.Error(), expected, "output")

	assert.Equal(t, testutil.StripProgress(b.String()), "Expand myvolume to example-go... done\n", "output")
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

	server.Mux.HandleFunc("/v2/apps/example-go/volumes/myvolume/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusNoContent)
	})

	err = cmdr.VolumesDelete("example-go", "myvolume")
	assert.NoError(t, err)

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
				Path: map[string]any{
					"cmd": "/data/cmd1",
				},
			}, r)
		}
		fmt.Fprintf(w, `{
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
			"owner": "test",
			"app": "example-go",
			"name": "myvolume",
			"size": "500G",
			"path": {"cmd": "/data/cmd1"},
			"type": "csi",
			"parameters": {},
			"created": "2020-08-26T00:00:00UTC",
			"updated": "2020-08-26T00:00:00UTC"
}`)
	})

	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.VolumesMount("example-go", "myvolume", []string{"cmd=/data/cmd1"})
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Mounting volume... done
The pods should be restart, please check the pods up or not.
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
				Path: map[string]any{
					"cmd": nil,
				},
			}, r)
		}
		fmt.Fprintf(w, `{
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
			"owner": "test",
			"app": "example-go",
			"name": "myvolume",
			"size": "500G",
			"path": {},
			"type": "csi",
			"parameters": {},
			"created": "2020-08-26T00:00:00UTC",
			"updated": "2020-08-26T00:00:00UTC"
}`)
	})

	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.VolumesUnmount("example-go", "myvolume", []string{"cmd"})
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Unmounting volume... done
The pods should be restart, please check the pods up or not.
`, "output")
}
