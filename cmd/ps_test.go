package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/workflow-cli/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/websocket"
)

func TestPrintProcesses(t *testing.T) {
	var b bytes.Buffer
	pods := []api.Pods{
		{
			Release:  "v3",
			Name:     "benign-quilting-web-4084101150-c871y",
			Type:     "web",
			State:    "up",
			Ready:    "1/1",
			Restarts: 0,
			Started:  "2023-11-15T11:55:16CST",
		},
		{
			Release:  "v3",
			Name:     "benign-quilting-worker-4084101150-c871y",
			Type:     "worker",
			State:    "up",
			Ready:    "1/1",
			Restarts: 0,
			Started:  "2023-11-15T11:55:16CST",
		},
	}
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	printProcesses(&DryccCmd{WOut: &b, ConfigFile: cf}, "appname", pods)

	assert.Equal(t, b.String(), `NAME                                       RELEASE    STATE    PTYPE     READY    RESTARTS    STARTED                
benign-quilting-web-4084101150-c871y       v3         up       web       1/1      0           2023-11-15T11:55:16CST    
benign-quilting-worker-4084101150-c871y    v3         up       worker    1/1      0           2023-11-15T11:55:16CST    
`, "output")
}

func TestPsList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/pods/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"count": 1,
			"next": null,
			"previous": null,
			"results": [
				{
					"release": "v2",
					"type": "web",
					"name": "foo-web-4084101150-c871y",
					"state": "up",
					"ready": "1/1",
					"restarts": 0,
					"started": "2016-02-13T00:47:52"
				}
			]
		}`)
	})

	err = cmdr.PsList("foo", -1)
	assert.NoError(t, err)

	assert.Equal(t, b.String(), `NAME                        RELEASE    STATE    PTYPE    READY    RESTARTS    STARTED             
foo-web-4084101150-c871y    v2         up       web      1/1      0           2016-02-13T00:47:52    
`, "output")
}

func TestPsExec(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}
	server.Mux.Handle(
		"/v2/apps/foo/pods/foo-web-111/exec/",
		websocket.Handler(func(conn *websocket.Conn) {
			io.Copy(conn, conn)
		}),
	)
	err = cmdr.PsExec("foo", "foo-web-111", true, false, []string{"/bin/sh"})
	assert.NoError(t, err)
}

func TestPsLogs(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}
	server.Mux.Handle(
		"/v2/apps/foo/pods/foo-web-111/logs/",
		websocket.Handler(func(conn *websocket.Conn) {
			conn.WriteClose(100)
		}),
	)
	err = cmdr.PsLogs("foo", "foo-web-111", 300, true, "runner")
	assert.NoError(t, err)
}

type psTargetCases struct {
	Targets       []string
	ExpectedError bool
	ExpectedMap   map[string]int
	ExpectedMsg   string
}

func TestParsePsTargets(t *testing.T) {
	t.Parallel()

	cases := []psTargetCases{
		{[]string{"test"}, true, nil, "'test' does not match the pattern 'ptype=num', ex: web=2"},
		{[]string{"test=a"}, true, nil, "'test=a' does not match the pattern 'ptype=num', ex: web=2"},
		{[]string{"test="}, true, nil, "'test=' does not match the pattern 'ptype=num', ex: web=2"},
		{[]string{"Test=2"}, true, nil, "'Test=2' does not match the pattern 'ptype=num', ex: web=2"},
		{[]string{"test=2"}, false, map[string]int{"test": 2}, ""},
		{[]string{"test-proc=2"}, false, map[string]int{"test-proc": 2}, ""},
		{[]string{"test1=2"}, false, map[string]int{"test1": 2}, ""},
	}

	for _, check := range cases {
		actual, err := parsePsTargets(check.Targets)
		if check.ExpectedError {
			assert.Equal(t, err.Error(), check.ExpectedMsg, "error")
		} else {
			assert.NoError(t, err)
			assert.Equal(t, actual, check.ExpectedMap, "error")
		}
	}
}

func TestPsDescribe(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/pods/foo-web-111/describe/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"count": 1,
			"next": null,
			"previous": null,
			"results": [{
				"container": "web",
				"image": "registry.drycc.cc/base/base",
				"command": ["bash", "-c"],
				"args": ["sleep 3600s"],
				"state": {
					"running": {
					  "startedAt": "2024-05-21T02:27:03+00:00"
					},
					"waiting": {
					  "message": "container create failed: executable file './start.sh' not found in $PATH: No such file or directory\n",
					  "reason": "CreateContainerError"
					}
				},
				"lastState": {
					"terminated": {
					  "containerID": "cri-o://ccfc73b0b4d966af4f93ca871a04fa97460620cd8005c1c36f7734a08ba49ed0",
					  "exitCode": 1,
					  "finishedAt": "2024-05-21T02:27:03+00:00",
					  "reason": "Error",
					  "startedAt": "2024-05-21T02:26:33+00:00"
					}
				},
				"ready": true,
				"restartCount": 1
			}]
		}`)
	})
	server.Mux.HandleFunc("/v2/apps/foo/events/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"count": 1,
			"next": null,
			"previous": null,
			"results": [{
				"reason": "Scheduled",
				"message": "Successfully assigned example-go/example-go-web-6b44dbd6c8-h89cg to node1",
				"created": "2024-07-03T16:28:00"
			}]
		}`)
	})
	err = cmdr.PsDescribe("foo", "foo-web-111")
	assert.NoError(t, err)
}

func TestPsDelete(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}
	server.Mux.HandleFunc("/v2/apps/foo/pods/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusNoContent)
	})
	err = cmdr.PsDelete("foo", []string{"foo-web-111"})
	assert.NoError(t, err)
}
