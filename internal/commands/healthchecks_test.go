package commands

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/workflow-cli/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestHealthchecksList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/config/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
  "uuid": "c039a380-6068-4511-b35a-535a73b86ef5",
  "app": "foo",
  "owner": "bar",
  "values": [],
  "memory": {},
  "cpu": {},
  "tags": {},
  "registry": {},
  "healthcheck": {
    "web": {
      "liveness_probe": {
        "initial_delay_seconds": 50,
        "timeout_seconds": 50,
        "period_seconds": 10,
        "failure_threshold": 3,
        "http_get": {
          "port": 80,
          "path": "/"
        },
        "success_threshold": 1
      }
    }
  },
  "created": "2016-09-12T22:20:14Z",
  "updated": "2016-09-12T22:20:14Z"
}`)
	})

	err = cmdr.HealthchecksList("foo", "web", -1)
	assert.NoError(t, err)

	actual := b.String()
	assert.Contains(t, actual, "App:")
	assert.Contains(t, actual, "foo")
	assert.Contains(t, actual, "c039a380-6068-4511-b35a-535a73b86ef5")
	assert.Contains(t, actual, "Healthchecks:")
	assert.Contains(t, actual, "liveness_probe web http-get headers=[] path=/ port=80 delay=50s timeout=50s period=10s #success=1 #failure=3")
}

func TestHealthchecksListNoHealthCheck(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/config/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
  "uuid": "c039a380-6068-4511-b35a-535a73b86ef5",
  "app": "foo",
  "owner": "bar",
  "values": [],
  "memory": {},
  "cpu": {},
  "tags": {},
  "registry": {},
  "healthcheck": {},
  "created": "2016-09-12T22:20:14Z",
  "updated": "2016-09-12T22:20:14Z"
}`)
	})

	err = cmdr.HealthchecksList("foo", "", -1)
	assert.NoError(t, err)

	assert.Equal(t, b.String(), `No health checks configured.
`, "output")
}

func TestHealthchecksListAllHealthChecks(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/config/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
  "uuid": "c039a380-6068-4511-b35a-535a73b86ef5",
  "app": "foo",
  "owner": "bar",
  "values": [],
  "memory": {},
  "cpu": {},
  "tags": {},
  "registry": {},
  "healthcheck": {
    "web": {
      "liveness_probe": {
        "initial_delay_seconds": 50,
        "timeout_seconds": 50,
        "period_seconds": 10,
        "failure_threshold": 3,
        "http_get": {
          "port": 80,
          "path": "/"
        },
        "success_threshold": 1
      }
    },
		"task": {
      "liveness_probe": {
        "initial_delay_seconds": 50,
        "timeout_seconds": 50,
        "period_seconds": 10,
        "failure_threshold": 3,
        "http_get": {
          "port": 8000,
          "path": "/"
        },
        "success_threshold": 1
      }
    }
  },
  "created": "2016-09-12T22:20:14Z",
  "updated": "2016-09-12T22:20:14Z"
}`)
	})

	err = cmdr.HealthchecksList("foo", "", -1)
	assert.NoError(t, err)

	actual := b.String()
	assert.Contains(t, actual, "App:")
	assert.Contains(t, actual, "foo")
	assert.Contains(t, actual, "c039a380-6068-4511-b35a-535a73b86ef5")
	assert.Contains(t, actual, "Healthchecks:")
	assert.Contains(t, actual, "liveness_probe task http-get headers=[] path=/ port=8000 delay=50s timeout=50s period=10s #success=1 #failure=3")
	assert.Contains(t, actual, "liveness_probe web http-get headers=[] path=/ port=80 delay=50s timeout=50s period=10s #success=1 #failure=3")
}

func TestHealthchecksSet(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/config/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
  "uuid": "c039a380-6068-4511-b35a-535a73b86ef5",
  "app": "foo",
  "owner": "bar",
  "values": [],
  "memory": {},
  "cpu": {},
  "tags": {},
  "registry": {},
  "healthcheck": {
    "web": {
      "liveness_probe": {
        "initial_delay_seconds": 50,
        "timeout_seconds": 50,
        "period_seconds": 10,
        "failure_threshold": 3,
        "http_get": {
          "port": 80,
          "path": "/"
        },
        "success_threshold": 1
      }
    }
  },
  "created": "2016-09-12T22:20:14Z",
  "updated": "2016-09-12T22:20:14Z"
}`)
	})

	err = cmdr.HealthchecksSet("foo", "liveness_probe", "web", &api.ContainerProbe{})
	assert.NoError(t, err)
	actual := testutil.StripProgress(b.String())
	assert.Contains(t, actual, "Applying liveness_probe healthcheck... done")
	assert.Contains(t, actual, "App:")
	assert.Contains(t, actual, "foo")
	assert.Contains(t, actual, "c039a380-6068-4511-b35a-535a73b86ef5")
	assert.Contains(t, actual, "Healthchecks:")
	assert.Contains(t, actual, "liveness_probe web http-get headers=[] path=/ port=80 delay=50s timeout=50s period=10s #success=1 #failure=3")
}

func TestHealthchecksUnset(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/config/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
  "uuid": "c039a380-6068-4511-b35a-535a73b86ef5",
  "app": "foo",
  "owner": "bar",
  "values": [],
  "memory": {},
  "cpu": {},
  "tags": {},
  "registry": {},
  "healthcheck": {},
  "created": "2016-09-12T22:20:14Z",
  "updated": "2016-09-12T22:20:14Z"
}`)
	})

	err = cmdr.HealthchecksUnset("foo", "web", []string{"liveness_probe"})
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), `Removing healthchecks... done

No health checks configured.
`, "output")
}
