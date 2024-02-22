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
  "values": {},
  "memory": {},
  "cpu": {},
  "tags": {},
  "registry": {},
  "healthcheck": {
    "web/cmd": {
      "livenessProbe": {
        "initialDelaySeconds": 50,
        "timeoutSeconds": 50,
        "periodSeconds": 10,
        "failureThreshold": 3,
        "httpGet": {
          "port": 80,
          "path": "/"
        },
        "successThreshold": 1
      }
    }
  },
  "created": "2016-09-12T22:20:14Z",
  "updated": "2016-09-12T22:20:14Z"
}`)
	})

	err = cmdr.HealthchecksList("foo", "web/cmd")
	assert.NoError(t, err)

	assert.Equal(t, b.String(), `App:             foo                                                                                                           
UUID:            c039a380-6068-4511-b35a-535a73b86ef5                                                                          
Owner:           bar                                                                                                           
Created:         2016-09-12T22:20:14Z                                                                                          
Updated:         2016-09-12T22:20:14Z                                                                                          
Healthchecks:    
                 liveness web/cmd http-get headers=[] path=/ port=80 delay=50s timeout=50s period=10s #success=1 #failure=3    
`, "output")
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
  "values": {},
  "memory": {},
  "cpu": {},
  "tags": {},
  "registry": {},
  "healthcheck": {},
  "created": "2016-09-12T22:20:14Z",
  "updated": "2016-09-12T22:20:14Z"
}`)
	})

	err = cmdr.HealthchecksList("foo", "")
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
  "values": {},
  "memory": {},
  "cpu": {},
  "tags": {},
  "registry": {},
  "healthcheck": {
    "web/cmd": {
      "livenessProbe": {
        "initialDelaySeconds": 50,
        "timeoutSeconds": 50,
        "periodSeconds": 10,
        "failureThreshold": 3,
        "httpGet": {
          "port": 80,
          "path": "/"
        },
        "successThreshold": 1
      }
    },
		"web": {
      "livenessProbe": {
        "initialDelaySeconds": 50,
        "timeoutSeconds": 50,
        "periodSeconds": 10,
        "failureThreshold": 3,
        "httpGet": {
          "port": 80,
          "path": "/"
        },
        "successThreshold": 1
      }
    }
  },
  "created": "2016-09-12T22:20:14Z",
  "updated": "2016-09-12T22:20:14Z"
}`)
	})

	err = cmdr.HealthchecksList("foo", "")
	assert.NoError(t, err)

	assert.Equal(t, b.String(), `App:             foo                                                                                                           
UUID:            c039a380-6068-4511-b35a-535a73b86ef5                                                                          
Owner:           bar                                                                                                           
Created:         2016-09-12T22:20:14Z                                                                                          
Updated:         2016-09-12T22:20:14Z                                                                                          
Healthchecks:    
                 liveness web http-get headers=[] path=/ port=80 delay=50s timeout=50s period=10s #success=1 #failure=3        
                 liveness web/cmd http-get headers=[] path=/ port=80 delay=50s timeout=50s period=10s #success=1 #failure=3    
`, "output")
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
  "values": {},
  "memory": {},
  "cpu": {},
  "tags": {},
  "registry": {},
  "healthcheck": {
    "web/cmd": {
      "livenessProbe": {
        "initialDelaySeconds": 50,
        "timeoutSeconds": 50,
        "periodSeconds": 10,
        "failureThreshold": 3,
        "httpGet": {
          "port": 80,
          "path": "/"
        },
        "successThreshold": 1
      }
    }
  },
  "created": "2016-09-12T22:20:14Z",
  "updated": "2016-09-12T22:20:14Z"
}`)
	})

	err = cmdr.HealthchecksSet("foo", "liveness", "web/cmd", &api.Healthcheck{})
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), `Applying liveness healthcheck... done

App:             foo                                                                                                           
UUID:            c039a380-6068-4511-b35a-535a73b86ef5                                                                          
Owner:           bar                                                                                                           
Created:         2016-09-12T22:20:14Z                                                                                          
Updated:         2016-09-12T22:20:14Z                                                                                          
Healthchecks:    
                 liveness web/cmd http-get headers=[] path=/ port=80 delay=50s timeout=50s period=10s #success=1 #failure=3    
`, "output")
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
  "values": {},
  "memory": {},
  "cpu": {},
  "tags": {},
  "registry": {},
  "healthcheck": {},
  "created": "2016-09-12T22:20:14Z",
  "updated": "2016-09-12T22:20:14Z"
}`)
	})

	err = cmdr.HealthchecksUnset("foo", "web/cmd", []string{"liveness"})
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), `Removing healthchecks... done

No health checks configured.
`, "output")
}
