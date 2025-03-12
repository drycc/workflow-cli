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
	assert.Equal(t, b.String(), `NAME        OWNER    TYPE    PTYPE     PATH            SIZE 
myvolume    test     csi     cmd       /data/cmd1      500G    
myvolume    test     csi     cmd123    /data/cmd123    500G    
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
		testutil.AssertBody(t, api.Volume{Name: "myvolume", Size: "500G", Type: "csi"}, r)
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusCreated)
		// Body isn't used by CLI, so it isn't set.
		w.Write([]byte("{}"))
	})

	err = cmdr.VolumesCreate("example-go", "myvolume", "csi", "500G", map[string]interface{}{})
	assert.NoError(t, err)
	err = cmdr.VolumesCreate("example-go", "myvolume", "csi", "500K", map[string]interface{}{})
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
	assert.Equal(t, b.String(), `UUID:          de1bf5b5-4a72-4f94-a10c-d2a3741cdf75    
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
`, "output")
}

func TestVolumesClientLs(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}
	server.Mux.HandleFunc("/v2/apps/example-go/volumes/myvolume/client/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{"results": [
  {"name":"handler.go","size":"4159","timestamp":"2024-06-25T22:55:16+08:00","type":"file","path":"/handler.go"},
  {"name":"handler_test.go","size":"2310","timestamp":"2024-06-04T15:29:45+08:00","type":"file","path":"/handler_test.go"}
], "count": 2}`)
	})

	err = cmdr.VolumesClient("example-go", "ls", "vol://myvolume")
	assert.NoError(t, err)
	assert.Contains(t, b.String(), "handler_test.go")
}

func TestVolumesClientCp(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	server.Mux.HandleFunc("/v2/apps/example-go/volumes/myvolume/client/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		if r.URL.RawQuery == "path=etc" {
			fmt.Fprintf(w, `{"results":[],"count":0}`)
		} else if r.Method == http.MethodGet {
			fmt.Fprintf(w, `{"results":[{"name":"hello.txt","size":"4159","timestamp":"2024-06-25T22:55:16+08:00","type":"file","path":"/hello.txt"}], "count": 1}`)
		}
	})
	server.Mux.HandleFunc("/v2/apps/example-go/volumes/myvolume/client/hello.txt", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `hello word`)
	})
	defer server.Close()

	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}
	// test download file
	err = cmdr.VolumesClient("example-go", "cp", "vol://myvolume/hello.txt", "/tmp")
	assert.NoError(t, err)
	result, err := os.ReadFile("/tmp/hello.txt")
	assert.NoError(t, err)
	assert.Equal(t, string(result), `hello word`, "output")
	// test upload file
	err = cmdr.VolumesClient("example-go", "cp", "/tmp/hello.txt", "vol://myvolume/etc")
	assert.NoError(t, err)
}

func TestVolumesClientRm(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}
	// test rm file
	server.Mux.HandleFunc("/v2/apps/example-go/volumes/myvolume/client/etc/hello.txt", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusOK)
	})
	err = cmdr.VolumesClient("example-go", "rm", "vol://myvolume/etc/hello.txt")
	assert.NoError(t, err)
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
