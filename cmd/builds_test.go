package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"reflect"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/workflow-cli/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
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

func TestBuildsInfo(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/build/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
				"app": "",
				"created": "2014-01-01T00:00:00UTC",
				"dockerfile": "",
				"image": "",
				"owner": "",
				"procfile": {},
				"sha": "",
				"updated": "",
				"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
			}
		`)
	})

	err = cmdr.BuildsInfo("foo", -1)
	assert.NoError(t, err)
	assert.Equal(t, b.String(), `App:                                                
Sha:                                                
UUID:       de1bf5b5-4a72-4f94-a10c-d2a3741cdf75    
Owner:                                              
Image:                                              
Stack:                                              
Created:    2014-01-01T00:00:00UTC                  
Updated:                                            
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

	server.Mux.HandleFunc("/v2/apps/enterprise/build/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		if r.Method == "POST" {
			testutil.AssertBody(t, api.CreateBuildRequest{
				Image: "ncc/1701:A",
				Stack: "container",
			}, r)
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, "{}")
		}
	})

	err = cmdr.BuildsCreate("enterprise", "ncc/1701:A", "container", "", "", "yes")
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), "Creating build... done\n", "output")

	server.Mux.HandleFunc("/v2/apps/bradbury/build/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		if r.Method == "POST" {
			testutil.AssertBody(t, api.CreateBuildRequest{
				Image: "nx/72307:latest",
				Stack: "container",
				Procfile: map[string]string{
					"web":  "./drive",
					"warp": "./warp 8",
				},
			}, r)
		}

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

	err = cmdr.BuildsCreate("bradbury", "nx/72307:latest", "container", tmpDir+"/Procfile", "", "yes")
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), "Creating build... done\n", "output")

	server.Mux.HandleFunc("/v2/apps/franklin/build/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		if r.Method == "POST" {
			testutil.AssertBody(t, api.CreateBuildRequest{
				Image: "nx/326:latest",
				Stack: "container",
				Procfile: map[string]string{
					"web":  "./drive",
					"warp": "./warp 8",
				},
				Dryccfile: map[string]interface{}{
					"pipeline": map[string]interface{}{
						"web.yaml": map[string]interface{}{
							"kind":  "pipeline",
							"ptype": "web",
							"deploy": map[string]interface{}{
								"command": []string{
									"bash",
									"-c",
								},
								"args": []string{
									"bundle exec puma -C config/puma.rb",
								},
							},
						},
					},
				},
			}, r)
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "{}")
	})
	b.Reset()

	err = os.WriteFile("Procfile", []byte(`web: ./drive
warp: ./warp 8
`), os.ModePerm)
	assert.NoError(t, err)

	dryccpath := ".drycc"
	os.MkdirAll(dryccpath, 0700)
	defer os.RemoveAll(dryccpath)
	err = os.WriteFile(filepath.Join(dryccpath, "web.yaml"), []byte(`
kind: pipeline
ptype: web
deploy:
  command:
  - bash
  - -c
  args:
  - bundle exec puma -C config/puma.rb
`), os.ModePerm)
	assert.NoError(t, err)

	err = cmdr.BuildsCreate("franklin", "nx/326:latest", "container", "Procfile", dryccpath, "yes")
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), "Creating build... done\n", "output")

}

func TestBuildsFetch(t *testing.T) {
	t.Parallel()

	// Mock server setup
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	cmdr := DryccCmd{WOut: new(bytes.Buffer), ConfigFile: cf}

	// Mock build data
	server.Mux.HandleFunc("/v2/apps/testapp/build/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"app": "testapp",
			"procfile": {"web": "node server.js"},
			"dryccfile": {
				"pipeline": {
					"web.yaml": {
						"kind": "pipeline",
						"ptype": "web",
						"deploy": {
							"command": ["bash", "-c"],
							"args": ["echo hello"]
						}
					}
				},
				"config": {
					"env1.env": {"KEY1": "VALUE1"},
					"env2.env": {"KEY2": "VALUE2"}
				}
			}
		}`)
	})

	// Helper function to compare YAML content ignoring field order
	isEqualYAML := func(actual, expected string) bool {
		var actualMap, expectedMap map[string]interface{}
		err := yaml.Unmarshal([]byte(actual), &actualMap)
		if err != nil {
			return false
		}
		err = yaml.Unmarshal([]byte(expected), &expectedMap)
		if err != nil {
			return false
		}
		return reflect.DeepEqual(actualMap, expectedMap)
	}

	// Test case 1: Successful fetch with valid confirm
	t.Run("Successful Fetch", func(t *testing.T) {
		tmpDir := t.TempDir() // Create a unique temporary directory for this test
		procfilePath := filepath.Join(tmpDir, "Procfile")
		dryccpath := filepath.Join(tmpDir, ".drycc")

		// Mock user input for confirmation
		restoreStdin := os.Stdin
		r, w, _ := os.Pipe()
		os.Stdin = r
		go func() {
			defer w.Close() // Ensure the pipe is closed after writing
			w.Write([]byte("yes\n"))
		}()

		err := cmdr.BuildsFetch("testapp", 0, procfilePath, dryccpath, "", true)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify Procfile content
		content, err := os.ReadFile(procfilePath)
		if err != nil {
			t.Fatalf("failed to read Procfile: %v", err)
		}
		expectedProcfile := "web: node server.js\n"
		if string(content) != expectedProcfile {
			t.Errorf("expected Procfile content %q, got %q", expectedProcfile, string(content))
		}

		// Verify .drycc/config/env1.env content
		env1Content, err := os.ReadFile(filepath.Join(dryccpath, "config", "env1.env"))
		if err != nil {
			t.Fatalf("failed to read env1.env: %v", err)
		}
		expectedEnv1 := "KEY1=VALUE1\n"
		if string(env1Content) != expectedEnv1 {
			t.Errorf("expected env1.env content %q, got %q", expectedEnv1, string(env1Content))
		}

		// Verify .drycc/web.yaml content
		webYamlContent, err := os.ReadFile(filepath.Join(dryccpath, "web.yaml"))
		if err != nil {
			t.Fatalf("failed to read web.yaml: %v", err)
		}
		expectedWebYaml := "kind: pipeline\nptype: web\ndeploy:\n  command:\n  - bash\n  - -c\n  args:\n  - echo hello\n"
		if !isEqualYAML(string(webYamlContent), expectedWebYaml) {
			t.Errorf("expected web.yaml content %q, got %q", expectedWebYaml, string(webYamlContent))
		}

		os.Stdin = restoreStdin
	})

	// Test case 2: User cancels the operation
	t.Run("User Cancels Operation", func(t *testing.T) {
		tmpDir := t.TempDir() // Create a unique temporary directory for this test
		procfilePath := filepath.Join(tmpDir, "Procfile")
		dryccpath := filepath.Join(tmpDir, ".drycc")

		// Mock user input for cancellation
		restoreStdin := os.Stdin
		r, w, _ := os.Pipe()
		os.Stdin = r
		go func() {
			defer w.Close() // Ensure the pipe is closed after writing
			w.Write([]byte("no\n"))
		}()

		err := cmdr.BuildsFetch("testapp", 0, procfilePath, dryccpath, "", true)
		if err == nil || err.Error() != "cancel the build fetch action" {
			t.Fatalf("expected cancellation error, got %v", err)
		}

		// Verify files were not created
		if _, err := os.Stat(procfilePath); !os.IsNotExist(err) {
			t.Errorf("expected Procfile to not exist, but it does")
		}

		if _, err := os.Stat(dryccpath); !os.IsNotExist(err) {
			t.Errorf("expected .drycc directory to not exist, but it does")
		}

		os.Stdin = restoreStdin
	})
}
