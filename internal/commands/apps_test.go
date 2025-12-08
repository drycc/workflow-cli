package commands

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/drycc/workflow-cli/pkg/git"
	"github.com/drycc/workflow-cli/pkg/settings"
	"github.com/drycc/workflow-cli/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

type expandURLCases struct {
	Input    string
	Expected string
}

func TestAppsList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"count": 2,
			"next": null,
			"previous": null,
			"results": [
			    {
					"uuid": "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
					"id": "lorem-ipsum",
					"owner": "dolar-sit-amet",
					"created": "2016-08-22T17:40:16Z",
					"updated": "2016-08-22T17:40:16Z",
					"structure": {
						"cmd": 1
					}
				},
				{
					"uuid": "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
					"id": "consectetur",
					"owner": "adipiscing",
					"created": "2016-08-22T17:40:16Z",
					"updated": "2016-08-22T17:40:16Z",
					"structure": {
						"cmd": 1
					}
				}
			]
		}`)
	})

	err = cmdr.AppsList(-1)
	assert.NoError(t, err)
	testutil.AssertOutput(t, b.String(), `ID             OWNER             CREATED                 UPDATED
lorem-ipsum    dolar-sit-amet    2016-08-22T17:40:16Z    2016-08-22T17:40:16Z
consectetur    adipiscing        2016-08-22T17:40:16Z    2016-08-22T17:40:16Z
`)
}

func TestAppsListLimit(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"count": 2,
			"next": null,
			"previous": null,
			"results": [
			    {
					"uuid": "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
					"id": "lorem-ipsum",
					"owner": "dolar-sit-amet",
					"created": "2016-08-22T17:40:16Z",
					"updated": "2016-08-22T17:40:16Z",
					"structure": {
						"cmd": 1
					}
				}
			]
		}`)
	})

	err = cmdr.AppsList(1)
	assert.NoError(t, err)
	testutil.AssertOutput(t, b.String(), `ID             OWNER             CREATED                 UPDATED
lorem-ipsum    dolar-sit-amet    2016-08-22T17:40:16Z    2016-08-22T17:40:16Z
`)
}

func TestAppsInfo(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/lorem-ipsum/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
    "uuid": "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
    "id": "lorem-ipsum",
    "owner": "dolar-sit-amet",
    "structure": {
      "cmd": 1
    },
    "created": "2016-08-22T17:40:16Z",
    "updated": "2016-08-22T17:40:16Z"
}`)
	})

	server.Mux.HandleFunc("/v2/apps/lorem-ipsum/pods/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
    "count": 1,
    "results": [
      {
        "state": "up",
        "started": "2016-08-22T17:42:16Z",
        "name": "lorem-ipsum-cmd-1911796442-48b58",
        "release": "v2",
        "type": "cmd"
      }
    ]
}`)
	})

	server.Mux.HandleFunc("/v2/apps/lorem-ipsum/domains/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
      {
         "owner": "dolar-sit-amet",
         "created": "2016-08-22T17:40:16Z",
         "updated": "2016-08-22T17:40:16Z",
         "app": "lorem-ipsum",
         "domain": "lorem-ipsum"
      }
    ]
}`)
	})

	server.Mux.HandleFunc("/v2/apps/lorem-ipsum/settings/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"owner": "elrond",
			"app": "lorem-ipsum",
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
			"label": {
				"team": "frontend"
			}
		}`)
	})
	s, err := settings.Load(cmdr.ConfigFile)
	if err != nil {
		t.Fatal(err)
	}
	url, err := cmdr.appURL(s, "lorem-ipsum")
	if err != nil {
		t.Fatal(err)
	}
	err = cmdr.AppInfo("lorem-ipsum")
	assert.NoError(t, err)
	testutil.AssertOutput(t, b.String(), `App:          lorem-ipsum
URL:          `+url+`
UUID:         c4aed81c-d1ca-4ff1-ab89-d2151264e1a3
Owner:        dolar-sit-amet
Created:      2016-08-22T17:40:16Z
Updated:      2016-08-22T17:40:16Z
Processes:
              Name:                                   lorem-ipsum-cmd-1911796442-48b58
              Release:                                v2
              State:                                  up
              Ptype:                                  cmd
              Started:                                2016-08-22T17:42:16Z
Domains:
              Domain:                                 lorem-ipsum
              Created:                                2016-08-22T17:40:16Z
              Updated:                                2016-08-22T17:40:16Z
Labels:
              Key:                                    team
              Value:                                  frontend
`)
}

func TestAppDestroy(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/lorem-ipsum/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
    "uuid": "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
    "id": "lorem-ipsum",
    "owner": "dolar-sit-amet",
    "structure": {
      "cmd": 1
    },
    "created": "2016-08-22T17:40:16Z",
    "updated": "2016-08-22T17:40:16Z"
}`)
	})
	err = cmdr.AppDestroy("lorem-ipsum", "bad-confirm-string")
	assert.Equal(t, err.Error(), `app lorem-ipsum does not match confirm bad-confirm-string, aborting`, "output")

	err = cmdr.AppDestroy("lorem-ipsum", "lorem-ipsum")
	assert.NoError(t, err)
	assert.Equal(t, b.String(), `Destroying lorem-ipsum...
done in 0s
`, "output")
}

func TestAppTransfer(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/lorem-ipsum/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
    "uuid": "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
    "id": "lorem-ipsum",
    "owner": "dolar-sit-amet",
    "structure": {
      "cmd": 1
    },
    "created": "2016-08-22T17:40:16Z",
    "updated": "2016-08-22T17:40:16Z"
}`)
	})

	err = cmdr.AppTransfer("lorem-ipsum", "test-user")
	assert.NoError(t, err)
	assert.Equal(t, b.String(), `Transferring lorem-ipsum to test-user... done
`, "output")
}

func TestExpandUrl(t *testing.T) {
	checks := []expandURLCases{
		{
			Input:    "test.com",
			Expected: "test.com",
		},
		{
			Input:    "test",
			Expected: "test.foo.com",
		},
	}

	for _, check := range checks {
		out := expandURL("drycc.foo.com", check.Input)

		if out != check.Expected {
			t.Errorf("Expected %s, Got %s", check.Expected, out)
		}
	}
}

func TestRemoteExists(t *testing.T) {
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/apps/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
    "owner": "jkirk",
    "id": "foo",
    "created": "2014-01-01T00:00:00UTC",
    "updated": "2014-01-01T00:00:00UTC",
    "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
}`)
	})

	// create a remote first before running apps:create
	dir, err := os.MkdirTemp("", "apps")
	assert.NoError(t, err)

	defer os.RemoveAll(dir)

	assert.NoError(t, os.Chdir(dir))

	assert.NoError(t, git.Init(git.DefaultCmd))
	assert.NoError(t, git.CreateRemote(git.DefaultCmd, "localhost", "drycc", "appname"))

	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.AppCreate("foo", "drycc", false)

	// Check that an error occurred and it contains the remote name
	// This works for any language since the remote name "drycc" is always in the error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "drycc", 
		"error message should contain the remote name 'drycc', got: %s", err.Error())
}
