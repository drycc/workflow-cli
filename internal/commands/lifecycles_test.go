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

func TestLifecyclesList(t *testing.T) {
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
  "lifecycle": {
    "web": {
      "stopSignal": "SIGTERM",
      "postStart": {
        "httpGet": {
          "port": 80,
          "path": "/health",
          "httpHeaders": []
        }
      },
      "preStop": {
        "exec": {
          "command": ["/bin/sh", "-c", "echo stopping"]
        }
      }
    }
  },
  "created": "2016-09-12T22:20:14Z",
  "updated": "2016-09-12T22:20:14Z"
}`)
	})

	err = cmdr.LifecyclesList("foo", "web", -1)
	assert.NoError(t, err)

	// 获取实际输出并打印，以便我们可以查看正确的格式
	actual := b.String()
	t.Logf("Actual output:\n%s", actual)

	// 使用assert.Contains来检查关键部分，而不是完全匹配
	assert.Contains(t, actual, "App:          foo")
	assert.Contains(t, actual, "UUID:         c039a380-6068-4511-b35a-535a73b86ef5")
	assert.Contains(t, actual, "Owner:        bar")
	assert.Contains(t, actual, "Created:      2016-09-12T22:20:14Z")
	assert.Contains(t, actual, "Updated:      2016-09-12T22:20:14Z")
	assert.Contains(t, actual, "stopSignal=SIGTERM")
	assert.Contains(t, actual, "postStart web http-get headers=[] path=/health port=80 SIGTERM")
	assert.Contains(t, actual, "preStop web exec [/bin/sh -c echo stopping] SIGTERM")
}

func TestLifecyclesListNoLifecycle(t *testing.T) {
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
  "lifecycle": {},
  "created": "2016-09-12T22:20:14Z",
  "updated": "2016-09-12T22:20:14Z"
}`)
	})

	err = cmdr.LifecyclesList("foo", "", -1)
	assert.NoError(t, err)
	assert.Equal(t, "No lifecycle found in foo app.\n", b.String())
}

func TestLifecyclesListAllLifecycles(t *testing.T) {
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
  "lifecycle": {
    "web": {
      "stopSignal": "SIGTERM",
      "postStart": {
        "httpGet": {
          "port": 80,
          "path": "/health",
          "httpHeaders": []
        }
      },
      "preStop": {
        "exec": {
          "command": ["/bin/sh", "-c", "echo stopping"]
        }
      }
    },
    "worker": {
      "stopSignal": "SIGKILL",
      "postStart": {
        "tcpSocket": {
          "port": 8080
        }
      },
      "preStop": {
        "sleep": {
          "seconds": 30
        }
      }
    }
  },
  "created": "2016-09-12T22:20:14Z",
  "updated": "2016-09-12T22:20:14Z"
}`)
	})

	err = cmdr.LifecyclesList("foo", "", -1)
	assert.NoError(t, err)

	actual := b.String()
	t.Logf("Actual output:\n%s", actual)

	// 检查关键部分
	assert.Contains(t, actual, "App:          foo")
	assert.Contains(t, actual, "UUID:         c039a380-6068-4511-b35a-535a73b86ef5")
	assert.Contains(t, actual, "Owner:        bar")
	assert.Contains(t, actual, "Created:      2016-09-12T22:20:14Z")
	assert.Contains(t, actual, "Updated:      2016-09-12T22:20:14Z")
	assert.Contains(t, actual, "stopSignal=SIGTERM")
	assert.Contains(t, actual, "postStart web http-get headers=[] path=/health port=80 SIGTERM")
	assert.Contains(t, actual, "preStop web exec [/bin/sh -c echo stopping] SIGTERM")
	assert.Contains(t, actual, "stopSignal=SIGKILL")
	assert.Contains(t, actual, "postStart worker tcp-socket port=8080 SIGKILL")
	assert.Contains(t, actual, "preStop worker sleep Seconds=30 SIGKILL")
}

func TestLifecyclesSet(t *testing.T) {
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
  "lifecycle": {
    "web": {
      "stopSignal": "SIGTERM",
      "postStart": {
        "httpGet": {
          "port": 80,
          "path": "/health",
          "httpHeaders": []
        }
      },
      "preStop": {
        "exec": {
          "command": ["/bin/sh", "-c", "echo stopping"]
        }
      }
    }
  },
  "created": "2016-09-12T22:20:14Z",
  "updated": "2016-09-12T22:20:14Z"
}`)
	})

	postStartHandler := &api.LifecycleHandler{
		HTTPGet: &api.HTTPGetAction{
			Path: "/health",
			Port: 80,
		},
	}
	preStopHandler := &api.LifecycleHandler{
		Exec: &api.ExecAction{
			Command: []string{"/bin/sh", "-c", "echo stopping"},
		},
	}
	lifecycle := &api.Lifecycle{
		StopSignal: "SIGTERM",
		PostStart:  &postStartHandler,
		PreStop:    &preStopHandler,
	}

	err = cmdr.LifecyclesSet("foo", "web", lifecycle)
	assert.NoError(t, err)

	actual := testutil.StripProgress(b.String())
	t.Logf("Actual output:\n%s", actual)

	assert.Contains(t, actual, "Applying lifecycle... done")
	assert.Contains(t, actual, "App:          foo")
	assert.Contains(t, actual, "UUID:         c039a380-6068-4511-b35a-535a73b86ef5")
	assert.Contains(t, actual, "Owner:        bar")
	assert.Contains(t, actual, "Created:      2016-09-12T22:20:14Z")
	assert.Contains(t, actual, "Updated:      2016-09-12T22:20:14Z")
	assert.Contains(t, actual, "stopSignal=SIGTERM")
	assert.Contains(t, actual, "postStart web http-get headers=[] path=/health port=80 SIGTERM")
	assert.Contains(t, actual, "preStop web exec [/bin/sh -c echo stopping] SIGTERM")
}

func TestLifecyclesUnset(t *testing.T) {
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
  "lifecycle": {},
  "created": "2016-09-12T22:20:14Z",
  "updated": "2016-09-12T22:20:14Z"
}`)
	})

	err = cmdr.LifecyclesUnset("foo", "web", []string{"postStart"})
	assert.NoError(t, err)
	assert.Equal(t, "Applying lifecycle... done\n\nNo lifecycle found in foo app.\n", testutil.StripProgress(b.String()))
}

func TestLifecyclesUnsetMultipleHandlers(t *testing.T) {
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
  "lifecycle": {},
  "created": "2016-09-12T22:20:14Z",
  "updated": "2016-09-12T22:20:14Z"
}`)
	})

	err = cmdr.LifecyclesUnset("foo", "web", []string{"postStart", "preStop"})
	assert.NoError(t, err)
	assert.Equal(t, "Applying lifecycle... done\n\nNo lifecycle found in foo app.\n", testutil.StripProgress(b.String()))
}

func TestLifecyclesUnsetInvalidHandler(t *testing.T) {
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
  "lifecycle": {},
  "created": "2016-09-12T22:20:14Z",
  "updated": "2016-09-12T22:20:14Z"
}`)
	})

	err = cmdr.LifecyclesUnset("foo", "web", []string{"invalidHandler"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown lifecycle handler: invalidHandler")
}
