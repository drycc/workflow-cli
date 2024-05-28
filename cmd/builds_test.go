package cmd

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

func TestParseProcfile(t *testing.T) {
	t.Parallel()

	procMap, err := parseProcfile([]byte(`web: ./test
foo: test --test
`))
	assert.NoError(t, err)
	assert.Equal(t, procMap, map[string]string{"web": "./test", "foo": "test --test"}, "map")

	_, err = parseProcfile([]byte(`web: ./test
foo
`))
	assert.NotEqual(t, err, nil, "yaml")
}

func TestBuildsList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/builds/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"count": 2,
			"next": null,
			"previous": null,
			"results": [
				{
					"app": "",
					"created": "2014-01-01T00:00:00UTC",
					"dockerfile": "",
					"image": "",
					"owner": "",
					"procfile": {},
					"sha": "",
					"updated": "",
					"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
				},
				{
					"app": "",
					"created": "2014-01-05T00:00:00UTC",
					"dockerfile": "",
					"image": "",
					"owner": "",
					"procfile": {},
					"sha": "",
					"updated": "",
					"uuid": "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3"
				}
			]
		}`)
	})

	err = cmdr.BuildsList("foo", -1)
	assert.NoError(t, err)
	assert.Equal(t, b.String(), `OWNER     SHA       CREATED                
<none>    <none>    2014-01-01T00:00:00UTC    
<none>    <none>    2014-01-05T00:00:00UTC    
`, "output")
}

func TestBuildsListLimit(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/builds/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
            "count": 2,
            "next": null,
            "previous": null,
            "results": [
                {
                    "app": "foo",
                    "created": "2014-01-01T00:00:00UTC",
                    "dockerfile": "",
                    "image": "",
                    "owner": "",
                    "procfile": {},
                    "sha": "",
                    "updated": "",
                    "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
                }
            ]
        }`)
	})

	err = cmdr.BuildsList("foo", 1)
	assert.NoError(t, err)
	assert.Equal(t, b.String(), `OWNER     SHA       CREATED                
<none>    <none>    2014-01-01T00:00:00UTC    
`, "output")
}

func TestBuildsCreate(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	// Create a new temporary directory and change to it.
	name, err := os.MkdirTemp("", "client")
	assert.NoError(t, err)
	err = os.Chdir(name)
	assert.NoError(t, err)

	server.Mux.HandleFunc("/v2/apps/enterprise/builds/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		testutil.AssertBody(t, api.CreateBuildRequest{
			Image: "ncc/1701:A",
			Stack: "container",
		}, r)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "{}")
	})

	err = cmdr.BuildsCreate("enterprise", "ncc/1701:A", "container", "", "")
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), "Creating build... done\n", "output")

	server.Mux.HandleFunc("/v2/apps/bradbury/builds/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		testutil.AssertBody(t, api.CreateBuildRequest{
			Image: "nx/72307:latest",
			Stack: "container",
			Procfile: map[string]string{
				"web":  "./drive",
				"warp": "./warp 8",
			},
		}, r)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "{}")
	})
	b.Reset()

	tmpDir, err := os.MkdirTemp("", "tmpdir")
	if err != nil {
		t.Fatalf("error creating temp directory (%s)", err)
	}
	data := `web: ./drive
warp: ./warp 8`
	if err := os.WriteFile(tmpDir+"/Procfile", []byte(data), 0644); err != nil {
		t.Fatalf("error creating %s/Procfile (%s)", tmpDir, err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Fatalf("failed to remove Procfile from %s (%s)", tmpDir, err)
		}
	}()

	err = cmdr.BuildsCreate("bradbury", "nx/72307:latest", "container", tmpDir+"/Procfile", "")
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), "Creating build... done\n", "output")

	server.Mux.HandleFunc("/v2/apps/franklin/builds/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		testutil.AssertBody(t, api.CreateBuildRequest{
			Image: "nx/326:latest",
			Stack: "container",
			Procfile: map[string]string{
				"web":  "./drive",
				"warp": "./warp 8",
			},
			Dryccfile: map[string]interface{}{
				"deploy": map[string]interface{}{
					"web": map[string]interface{}{
						"command": []string{"bash", "-c"},
						"args":    []string{"bundle exec puma -C config/puma.rb"},
					},
				},
			},
		}, r)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "{}")
	})
	b.Reset()

	err = os.WriteFile("Procfile", []byte(`web: ./drive
warp: ./warp 8
`), os.ModePerm)
	assert.NoError(t, err)

	err = os.WriteFile("drycc.yaml", []byte(`
deploy:
  web:
    command:
    - bash
    - -c
    args:
    - bundle exec puma -C config/puma.rb
`), os.ModePerm)
	assert.NoError(t, err)

	err = cmdr.BuildsCreate("franklin", "nx/326:latest", "container", "Procfile", "drycc.yaml")
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), "Creating build... done\n", "output")

}
