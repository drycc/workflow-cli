package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/workflow-cli/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

const getPlanFixture string = `
{
	"id": "std1.large.c1m1",
	"spec": {
	  "id": "std1",
	  "cpu": {
		"name": "Unknown CPU",
		"cores": 32,
		"clock": "3100MHZ",
		"boost": "3700MHZ",
		"threads": 64
	  },
	  "memory": {
		"size": "64GB",
		"type": "DDR4-ECC"
	  },
	  "features": {
		"gpu": {
		  "name": "Unknown Integrated GPU",
		  "tmus": 1,
		  "rops": 1,
		  "cores": 128,
		  "memory": {
			"size": "shared",
			"type": "shared"
		  }
		},
		"network": "10G"
	  },
	  "keywords": [
		"amd",
		"intel",
		"unknown"
	  ],
	  "disabled": false
	},
	"cpu": 1,
	"memory": 1,
	"features": {
		"gpu": 1,
		"network": 1
	},
	"disabled": false
}
`
const specsFixture string = `
{
	"results": [{
		"id": "std1",
		"cpu": {
			"name": "Unknown CPU",
			"cores": 32,
			"clock": "3100MHZ",
			"boost": "3700MHZ",
			"threads": 64
		},
		"memory": {
			"size": "64GB",
			"type": "DDR4-ECC"
		},
		"features": {
			"gpu": {
				"name": "Unknown Integrated GPU",
				"tmus": 1,
				"rops": 1,
				"cores": 128,
				"memory": {
					"size": "shared",
					"type": "shared"
				}
			},
			"network": "10G"
		},
		"keywords": [
			"amd",
			"intel",
			"unknown"
		],
		"disabled": false
	}],
	"count": 1
}
`

const plansFixture string = `
{
	"results": [{
			"id": "std1.large.c1m1",
			"spec": {
				"id": "std1",
				"cpu": {
					"name": "Unknown CPU",
					"cores": 32,
					"clock": "3100MHZ",
					"boost": "3700MHZ",
					"threads": 64
				},
				"memory": {
					"size": "64GB",
					"type": "DDR4-ECC"
				},
				"features": {
					"gpu": {
						"name": "Unknown Integrated GPU",
						"tmus": 1,
						"rops": 1,
						"cores": 128,
						"memory": {
							"size": "shared",
							"type": "shared"
						}
					},
					"network": "10G"
				},
				"keywords": [
					"amd",
					"intel",
					"unknown"
				],
				"disabled": false
			},
			"cpu": 1,
			"memory": 1,
			"disabled": false
		},
		{
			"id": "std1.large.c1m2",
			"spec": {
				"id": "std1",
				"cpu": {
					"name": "Unknown CPU",
					"cores": 32,
					"clock": "3100MHZ",
					"boost": "3700MHZ",
					"threads": 64
				},
				"memory": {
					"size": "64GB",
					"type": "DDR4-ECC"
				},
				"features": {
					"gpu": {
						"name": "Unknown Integrated GPU",
						"tmus": 1,
						"rops": 1,
						"cores": 128,
						"memory": {
							"size": "shared",
							"type": "shared"
						}
					},
					"network": "10G"
				},
				"keywords": [
					"amd",
					"intel",
					"unknown"
				],
				"disabled": false
			},
			"cpu": 1,
			"memory": 2,
			"disabled": false
		}
	],
	"count": 2
}
`

type parseLimitCase struct {
	Input         string
	Key           string
	Value         string
	ExpectedError bool
	ExpectedMsg   string
}

func newTestServer(t *testing.T) (string, *testutil.TestServer) {
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	server.Mux.HandleFunc("/v2/limits/plans/std1.large.c1m1/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)
		testutil.SetHeaders(w)
		fmt.Fprint(w, getPlanFixture)
	})
	server.Mux.HandleFunc("/v2/limits/specs/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)
		testutil.SetHeaders(w)
		fmt.Fprint(w, specsFixture)
	})

	server.Mux.HandleFunc("/v2/limits/plans/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)
		testutil.SetHeaders(w)
		fmt.Fprint(w, plansFixture)
	})

	return cf, server
}

func TestParseLimit(t *testing.T) {
	t.Parallel()

	var errorHint = ` doesn't fit format type=#unit or type=#
Examples: web=std1.large.c1m1`

	cases := []parseLimitCase{
		{"web=std1.large.c1m1", "web", "std1.large.c1m1", false, ""},
		{"web=std1.large.c2m2", "web", "std1.large.c2m2", false, ""},
		{"task=std1.large.c2m2", "task", "std1.large.c2m2", false, ""},
		{"task=std1.large.c2m4", "task", "std1.large.c2m4", false, ""},
		{"task-big=std1.large.c2m4", "task-big", "std1.large.c2m4", false, ""},
		{"task=[]std1.large.c2m4", "", "", true, "task=[]std1.large.c2m4" + errorHint},
		{"task[]=&std1.large.c2m4", "", "", true, "task[]=&std1.large.c2m4" + errorHint},
		{"task~!=&std1.large.c2m4", "", "", true, "task~!=&std1.large.c2m4" + errorHint},
	}

	for _, check := range cases {
		key, value, err := parseLimit(check.Input)
		if check.ExpectedError {
			assert.Equal(t, err.Error(), check.ExpectedMsg, "error")
		} else {
			assert.NoError(t, err)
			assert.Equal(t, key, check.Key, "key")
			assert.Equal(t, value, check.Value, "value")
		}
	}
}

type parseLimitsCase struct {
	Input         []string
	ExpectedMap   map[string]interface{}
	ExpectedError bool
	ExpectedMsg   string
}

func TestLimitTags(t *testing.T) {
	t.Parallel()

	cases := []parseLimitsCase{
		{[]string{"web=std1.large.c1m1", "worker=std1.large.c1m2"}, map[string]interface{}{"web": "std1.large.c1m1", "worker": "std1.large.c1m2"}, false, ""},
		{[]string{"foo=", "web=std1.large.c1m1"}, nil, true, `foo= doesn't fit format type=#unit or type=#
Examples: web=std1.large.c1m1`},
	}

	for _, check := range cases {
		actual, err := parseLimits(check.Input)
		if check.ExpectedError {
			assert.Equal(t, err.Error(), check.ExpectedMsg, "error")
		} else {
			assert.NoError(t, err)
			assert.Equal(t, actual, check.ExpectedMap, "map")
		}
	}
}

func TestLimitsList(t *testing.T) {
	t.Parallel()
	cf, server := newTestServer(t)
	defer server.Close()

	server.Mux.HandleFunc("/v2/apps/enterprise/config/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"owner": "jkirk",
			"app": "enterprise",
			"values": [],
			"limits": {
				"web": "std1.large.c1m1",
				"worker": "std1.large.c1m1",
				"db": "std1.large.c1m1"
			},
			"tags": {},
			"registry": {},
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})

	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	err := cmdr.LimitsList("enterprise", -1)
	assert.NoError(t, err)
	assert.Equal(t, b.String(), `PTYPE     PLAN               VCPUS    MEMORY    FEATURES                          
db        std1.large.c1m1    1        1 GiB     Unknown Integrated GPU shared * 1    
web       std1.large.c1m1    1        1 GiB     Unknown Integrated GPU shared * 1    
worker    std1.large.c1m1    1        1 GiB     Unknown Integrated GPU shared * 1    
`, "output")

	server.Mux.HandleFunc("/v2/apps/franklin/config/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"owner": "bedison",
			"app": "franklin",
			"limits": {},
			"cpu": {},
			"tags": {},
			"registry": {},
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
			}`)
	})
	b.Reset()

	err = cmdr.LimitsList("franklin", -1)
	assert.NoError(t, err)
	assert.Equal(t, b.String(), `No limits found in franklin app.
`, "output")
}

func TestLimitsSet(t *testing.T) {
	t.Parallel()
	cf, server := newTestServer(t)
	defer server.Close()

	server.Mux.HandleFunc("/v2/apps/foo/config/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		if r.Method == "POST" {
			testutil.AssertBody(t, api.Config{
				Limits: map[string]interface{}{
					"web": "std1.large.c1m1",
				},
			}, r)
		}

		fmt.Fprintf(w, `{
			"owner": "jkirk",
			"app": "foo",
			"values": [],
			"limits": {"web": "std1.large.c1m1"},
			"tags": {},
			"registry": {},
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})

	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	err := cmdr.LimitsSet("foo", []string{"web=std1.large.c1m1"})
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Applying limits... done

PTYPE    PLAN               VCPUS    MEMORY    FEATURES                          
web      std1.large.c1m1    1        1 GiB     Unknown Integrated GPU shared * 1    
`, "output")

	server.Mux.HandleFunc("/v2/apps/franklin/config/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		if r.Method == "POST" {
			testutil.AssertBody(t, api.Config{
				Limits: map[string]interface{}{
					"web": "std1.large.c1m1",
				},
			}, r)
		}

		fmt.Fprintf(w, `{
			"owner": "bedison",
			"app": "franklin",
			"values": [],
			"limits": {
				"web": "std1.large.c1m1"
			},
			"tags": {},
			"registry": {},
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})
	b.Reset()

	err = cmdr.LimitsSet("franklin", []string{"web=std1.large.c1m1"})
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Applying limits... done

PTYPE    PLAN               VCPUS    MEMORY    FEATURES                          
web      std1.large.c1m1    1        1 GiB     Unknown Integrated GPU shared * 1    
`, "output")

	// with requests/limit parameter
	server.Mux.HandleFunc("/v2/apps/jim/config/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		if r.Method == "POST" {
			testutil.AssertBody(t, api.Config{
				Limits: map[string]interface{}{
					"web":    "std1.large.c1m1",
					"worker": "std1.large.c1m1",
					"db":     "std1.large.c1m1",
				},
			}, r)
		}

		fmt.Fprintf(w, `{
			"owner": "foo",
			"app": "jim",
			"values": [],
			"limits": {
				"web": "std1.large.c1m1",
				"worker": "std1.large.c1m1",
				"db": "std1.large.c1m1"
			},
			"tags": {},
			"registry": {},
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})
	b.Reset()

	err = cmdr.LimitsSet("jim", []string{"web=std1.large.c1m1", "worker=std1.large.c1m1", "db=std1.large.c1m1"})
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Applying limits... done

PTYPE     PLAN               VCPUS    MEMORY    FEATURES                          
db        std1.large.c1m1    1        1 GiB     Unknown Integrated GPU shared * 1    
web       std1.large.c1m1    1        1 GiB     Unknown Integrated GPU shared * 1    
worker    std1.large.c1m1    1        1 GiB     Unknown Integrated GPU shared * 1    
`, "output")

	// with requests/limit parameter
	server.Mux.HandleFunc("/v2/apps/phew/config/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		if r.Method == "POST" {
			testutil.AssertBody(t, api.Config{
				Limits: map[string]interface{}{
					"web":    "std1.large.c1m1",
					"worker": "std1.large.c1m1",
					"db":     "std1.large.c1m1",
				},
			}, r)
		}

		fmt.Fprintf(w, `{
			"owner": "foo",
			"app": "jim",
			"values": [],
			"limits": {
				"web": "std1.large.c1m1",
				"worker": "std1.large.c1m1",
				"db": "std1.large.c1m1"
			},
			"tags": {},
			"registry": {},
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})
	b.Reset()

	err = cmdr.LimitsSet("phew", []string{"web=std1.large.c1m1", "worker=std1.large.c1m1", "db=std1.large.c1m1"})
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Applying limits... done

PTYPE     PLAN               VCPUS    MEMORY    FEATURES                          
db        std1.large.c1m1    1        1 GiB     Unknown Integrated GPU shared * 1    
web       std1.large.c1m1    1        1 GiB     Unknown Integrated GPU shared * 1    
worker    std1.large.c1m1    1        1 GiB     Unknown Integrated GPU shared * 1    
`, "output")
}

func TestLimitsUnset(t *testing.T) {
	t.Parallel()
	cf, server := newTestServer(t)
	defer server.Close()

	server.Mux.HandleFunc("/v2/apps/foo/config/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		if r.Method == "POST" {
			testutil.AssertBody(t, api.Config{
				Limits: map[string]interface{}{
					"web": nil,
				},
			}, r)
		}

		fmt.Fprintf(w, `{
			"owner": "jkirk",
			"app": "foo",
			"values": [],
			"limits": {
				"web": "std1.large.c1m1"
			},
			"tags": {},
			"registry": {},
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})

	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	err := cmdr.LimitsUnset("foo", []string{"web"})
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Applying limits... done

PTYPE    PLAN               VCPUS    MEMORY    FEATURES                          
web      std1.large.c1m1    1        1 GiB     Unknown Integrated GPU shared * 1    
`, "output")

	server.Mux.HandleFunc("/v2/apps/franklin/config/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		if r.Method == "POST" {
			testutil.AssertBody(t, api.Config{
				Limits: map[string]interface{}{
					"web": nil,
				},
			}, r)
		}

		fmt.Fprintf(w, `{
			"owner": "bedison",
			"app": "franklin",
			"values": [],
			"limits": {},
			"tags": {},
			"registry": {},
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})
	b.Reset()

	err = cmdr.LimitsUnset("franklin", []string{"web"})
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Applying limits... done

No limits found in franklin app.
`, "output")
}

func TestLimitsSpecs(t *testing.T) {
	t.Parallel()
	cf, server := newTestServer(t)
	defer server.Close()

	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	err := cmdr.LimitsSpecs("", 10)
	assert.NoError(t, err)
	assert.Equal(t, b.String(), `ID      CPU            CLOCK      BOOST      CORES    THREADS    NETWORK    FEATURES                      
std1    Unknown CPU    3100MHZ    3700MHZ    32       64         10G        Unknown Integrated GPU shared    
`, "output")
}

func TestLimitsPlans(t *testing.T) {
	t.Parallel()
	cf, server := newTestServer(t)
	defer server.Close()

	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	err := cmdr.LimitsPlans("", 0, 0, 100)
	assert.NoError(t, err)
	assert.Equal(t, b.String(), `ID                 SPEC    CPU            VCPUS    MEMORY    FEATURES                      
std1.large.c1m1    std1    Unknown CPU    1        1 GiB     Unknown Integrated GPU shared    
std1.large.c1m2    std1    Unknown CPU    1        2 GiB     Unknown Integrated GPU shared    
`, "output")
}
