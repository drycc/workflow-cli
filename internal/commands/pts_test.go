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

func TestPrintProcessTypes(t *testing.T) {
	var b bytes.Buffer
	ptypes := api.Ptypes{
		{
			Name:              "web",
			Release:           "v1",
			Ready:             "1/1",
			UpToDate:          1,
			AvailableReplicas: 1,
			Started:           "2024-07-04T14:33:00CST",
			Garbage:           false,
		},
		{
			Name:              "worker",
			Release:           "v1",
			Ready:             "1/1",
			UpToDate:          1,
			AvailableReplicas: 1,
			Started:           "2024-07-04T14:33:00CST",
			Garbage:           false,
		},
	}
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	printProcessTypes(&DryccCmd{WOut: &b, ConfigFile: cf}, "appname", ptypes)

	assert.Equal(t, b.String(), `NAME      RELEASE    READY    UP-TO-DATE    AVAILABLE    STARTED                   GARBAGE 
web       v1         1/1      1             1            2024-07-04T14:33:00CST    false      
worker    v1         1/1      1             1            2024-07-04T14:33:00CST    false      
`, "output")
}

func TestPtsList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/ptypes/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"count": 1,
			"next": null,
			"previous": null,
			"results": [
				{
					"name": "web",
					"release": "v1",
					"ready": "1/1",
					"up_to_date": 1,
					"available_replicas": 1,
					"started": "2016-02-13T00:47:52",
					"garbage": false
				}
			]
		}`)
	})

	err = cmdr.PtsList("foo", -1)
	assert.NoError(t, err)

	assert.Equal(t, b.String(), `NAME    RELEASE    READY    UP-TO-DATE    AVAILABLE    STARTED                GARBAGE 
web     v1         1/1      1             1            2016-02-13T00:47:52    false      
`, "output")
}

type ptsTargetCases struct {
	Targets       []string
	ExpectedError bool
	ExpectedMap   map[string]int
	ExpectedMsg   string
}

func TestParsePtsTargets(t *testing.T) {
	t.Parallel()

	cases := []ptsTargetCases{
		{[]string{"test"}, true, nil, "'test' does not match the pattern 'ptype=num', ex: web=2"},
		{[]string{"test=a"}, true, nil, "'test=a' does not match the pattern 'ptype=num', ex: web=2"},
		{[]string{"test="}, true, nil, "'test=' does not match the pattern 'ptype=num', ex: web=2"},
		{[]string{"Test=2"}, true, nil, "'Test=2' does not match the pattern 'ptype=num', ex: web=2"},
		{[]string{"test=2"}, false, map[string]int{"test": 2}, ""},
		{[]string{"test-proc=2"}, false, map[string]int{"test-proc": 2}, ""},
		{[]string{"test1=2"}, false, map[string]int{"test1": 2}, ""},
	}

	for _, check := range cases {
		actual, err := parsePtsTargets(check.Targets)
		if check.ExpectedError {
			assert.Equal(t, err.Error(), check.ExpectedMsg, "error")
		} else {
			assert.NoError(t, err)
			assert.Equal(t, actual, check.ExpectedMap, "error")
		}
	}
}

func TestPtsScale(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}
	err = cmdr.PtsScale("foo", []string{"test"})
	assert.Equal(t, err.Error(), "'test' does not match the pattern 'ptype=num', ex: web=2", "error")

	server.Mux.HandleFunc("/v2/apps/foo/ptypes/scale/", func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertBody(t, map[string]int{"web": 1}, r)
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusNoContent)
	})

	b.Reset()
	err = cmdr.PtsScale("foo", []string{"web=1"})
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), `Scaling process types... but first, coffee!
done in 0s
`, "output")
}

func TestPtsRestart(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/ptypes/restart/", func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertBody(t, map[string]string{"ptypes": ""}, r)
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusNoContent)
	})

	b.Reset()
	err = cmdr.PtsRestart("foo", []string{""}, "yes")
	assert.NoError(t, err)

	server.Mux.HandleFunc("/v2/apps/coolapp/ptypes/restart/", func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertBody(t, map[string]string{"ptypes": "web"}, r)
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusNoContent)
	})

	b.Reset()

	err = cmdr.PtsRestart("coolapp", []string{"web"}, "")
	assert.NoError(t, err)

	server.Mux.HandleFunc("/v2/apps/testapp/ptypes/restart/", func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertBody(t, map[string]string{"ptypes": "web,worker"}, r)
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusNoContent)
	})

	b.Reset()

	err = cmdr.PtsRestart("testapp", []string{"web", "worker"}, "")
	assert.NoError(t, err)
}

func TestPtsDescribe(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}
	server.Mux.HandleFunc("/v2/apps/foo/ptypes/foo-web/describe/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
            "count": 1,
            "next": null,
            "previous": null,
            "results": [{
                "container": "web",
                "image": "registry.drycc.cc/base/base",
                "command": ["bash", "-c"],
                "args": ["sleep", "3600s"],
                "readiness_probe": {
                    "exec": {
                        "command": ["ls", "-la"]
                    },
                    "failureThreshold": 3,
                    "initialDelaySeconds": 50,
                    "periodSeconds": 10,
                    "successThreshold": 1,
                    "timeoutSeconds": 50
                },
                "limits": {
                    "cpu": "1",
                    "memory": "2Gi"
                },
                "volume_mounts": [
                    {
                        "mountPath": "/data",
                        "name": "myvolume"
                    }
                ]
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
                "reason": "ScalingReplicaSet",
                "message": "Scaled up replica set example-go-web-6b44dbd6c8 to 2 from 1",
                "created": "2024-07-03T16:28:00"
            }]
        }`)
	})
	err = cmdr.PtsDescribe("foo", "web")
	assert.NoError(t, err)
}

func TestPtsClean(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/testapp/ptypes/clean/", func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertBody(t, map[string]string{"ptypes": "worker"}, r)
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusNoContent)
	})

	b.Reset()

	err = cmdr.PtsClean("testapp", []string{"worker"})
	assert.NoError(t, err)
}
