package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/drycc/controller-sdk-go/api"

	"github.com/arschles/assert"
	"github.com/drycc/workflow-cli/pkg/testutil"
)

func TestResourcesServices(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/resources/services/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
	"results": [
		{
			"id": "332588e0-6c2c-4f56-a6af-a56fd01ec4b4",
			"name": "mysql",
			"updateable": true
		}
	],
	"count": 1
}`)
	})

	err = cmdr.ResourcesServices(100)
	assert.NoErr(t, err)

	assert.Equal(t, b.String(), `+-------+------------+
| NAME  | UPDATEABLE |
+-------+------------+
| mysql | true       |
+-------+------------+
`, "output")
}

func TestResourcesPlans(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/resources/services/mysql/plans/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
	"results": [{
		"id": "4d1dbd33-201b-45bc-9abb-757584ef7ab8",
		"name": "standard-1600",
		"description": "mysql standard-1600 plan which limit persistence size 1600Gi."
	}],
	"count": 1
}`)
	})

	err = cmdr.ResourcesPlans("mysql", 100)
	assert.NoErr(t, err)

	assert.Equal(t, b.String(), `+---------------+--------------------------------+
|     NAME      |          DESCRIPTION           |
+---------------+--------------------------------+
| standard-1600 | mysql standard-1600 plan which |
|               | limit persistence size 1600Gi. |
+---------------+--------------------------------+
`, "output")
}

func TestResourcesCreate(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/example-go/resources/", func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertBody(t, api.Resource{Name: "mysql", Plan: "mysql:5.6"}, r)
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusCreated)
		// Body isn't used by CLI, so it isn't set.
		w.Write([]byte("{}"))
	})

	err = cmdr.ResourcesCreate("example-go", "mysql:5.6", "mysql", nil)
	assert.NoErr(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Creating mysql to example-go... done\n", "output")
}

func TestResourcesList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/example-go/resources/", func(w http.ResponseWriter, r *http.Request) {
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
			"name": "mysql",
			"plan": "mysql:5.6",
			"data": {},
			"options": {},
			"status": null,
			"binding": null,
			"created": "2020-09-08T00:00:00UTC",
			"updated": "2020-09-08T00:00:00UTC"
		}
   ]
}`)
	})

	err = cmdr.ResourcesList("example-go", -1)
	assert.NoErr(t, err)

	assert.Equal(t, b.String(), `=== example-go resources
mysql     mysql:5.6
`, "output")
}

func TestResourceGet(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/example-go/resources/mysql/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
			"owner": "test",
			"app": "example-go",
			"name": "mysql",
			"plan": "mysql:5.6",
			"data": {"data12":"value1","data3":"value1"},
			"options": {"para13451":"value2","para122":"value1"},
			"status": "Ready",
			"binding": "Ready",
			"created": "2020-09-08T00:00:00UTC",
			"updated": "2020-09-08T00:00:00UTC"
}`)
	})

	err = cmdr.ResourceGet("example-go", "mysql")
	assert.NoErr(t, err)
	// todo format data json to yaml
	assert.Equal(t, b.String(), `=== example-go resource mysql
plan:          mysql:5.6
status:        Ready
binding:       Ready

data12:        value1
data3:         value1

para122:       value1
para13451:     value2
`, "output")
}

func TestResourceDelete(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/example-go/resources/mysql/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusNoContent)
	})

	err = cmdr.ResourceDelete("example-go", "mysql")
	assert.NoErr(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Deleting mysql from example-go... done\n", "output")
}

func TestResourcePut(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/example-go/resources/mysql/", func(w http.ResponseWriter, r *http.Request) {
		paras := make(map[string]interface{}, 1)
		paras["para1.para2"] = "v1"
		testutil.AssertBody(t, api.Resource{Plan: "mysql:5.7", Options: paras}, r)
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusCreated)
		// Body isn't used by CLI, so it isn't set.
		w.Write([]byte("{}"))
	})

	err = cmdr.ResourcePut("example-go", "mysql:5.7", "mysql", []string{"para1.para2=v1"})
	assert.NoErr(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Updating mysql to example-go... done\n", "output")
}

func TestResourceBind(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/apps/example-go/resources/mysql/binding/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		if r.Method == "PATCH" {
			testutil.AssertBody(t, api.ResourceBinding{
				BindAction: "bind",
			}, r)
		}
		fmt.Fprintf(w, `{
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
			"owner": "test",
			"app": "example-go",
			"name": "mysql",
			"plan": "mysql:5.7",
			"data": {},
			"options": {},
			"status": null,
			"binding": null,
			"created": "2020-09-08T00:00:00UTC",
			"updated": "2020-09-08T00:00:00UTC"
}`)
	})

	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.ResourceBind("example-go", "mysql")
	assert.NoErr(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Binding resource... done
`, "output")
}

func TestResourceUnbind(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/apps/example-go/resources/mysql/binding/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		if r.Method == "PATCH" {
			testutil.AssertBody(t, api.ResourceBinding{
				BindAction: "unbind",
			}, r)
		}
		fmt.Fprintf(w, `{
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
			"owner": "test",
			"app": "example-go",
			"name": "mysql",
			"plan": "mysql:5.7",
			"data": {},
			"options": {},
			"status": null,
			"binding": null,
			"created": "2020-09-08T00:00:00UTC",
			"updated": "2020-09-08T00:00:00UTC"
}`)
	})

	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.ResourceUnbind("example-go", "mysql")
	assert.NoErr(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Unbinding resource... done
`, "output")
}
