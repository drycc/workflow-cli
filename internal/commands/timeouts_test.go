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

type parseTimeoutCase struct {
	Input         string
	Key           string
	Value         string
	ExpectedError bool
	ExpectedMsg   string
}

func TestParseTimeout(t *testing.T) {
	t.Parallel()

	errorHint := ` doesn't fit format type=#
Examples: web=30 worker=300`

	cases := []parseTimeoutCase{
		{"web=20", "web", "20", false, ""},
		{"=1", "", "", true, "=1" + errorHint},
		{"web=", "", "", true, "web=" + errorHint},
		{"1=", "", "", true, "1=" + errorHint},
		{"web=ABCD", "", "", true, "web=ABCD" + errorHint},
	}

	for _, check := range cases {
		key, value, err := parseTimeout(check.Input)
		if check.ExpectedError {
			assert.Equal(t, err.Error(), check.ExpectedMsg, "error")
		} else {
			assert.NoError(t, err)
			assert.Equal(t, key, check.Key, "key")
			assert.Equal(t, value, check.Value, "value")
		}
	}
}

type parseTimeoutsCase struct {
	Input         []string
	ExpectedMap   map[string]any
	ExpectedError bool
	ExpectedMsg   string
}

func TestTimeoutTags(t *testing.T) {
	t.Parallel()

	cases := []parseTimeoutsCase{
		{[]string{"web=10", "worker=20"}, map[string]any{"web": "10", "worker": "20"}, false, ""},
		{[]string{"foo=", "web=10"}, nil, true, `foo= doesn't fit format type=#
Examples: web=30 worker=300`},
	}

	for _, check := range cases {
		actual, err := parseTimeouts(check.Input)
		if check.ExpectedError {
			assert.Equal(t, err.Error(), check.ExpectedMsg, "error")
		} else {
			assert.NoError(t, err)
			assert.Equal(t, actual, check.ExpectedMap, "map")
		}
	}
}

func TestTimeoutsList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/apps/enterprise/config/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
      "owner": "jkirk",
      "app": "enterprise",
      "values": [],
      "memory": {},
      "cpu": {},
      "tags": {},
      "registry": {},
      "created": "2014-01-01T00:00:00UTC",
      "updated": "2014-01-01T00:00:00UTC",
      "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
      "termination_grace_period": {
        "web" : 10,
        "worker" : 20
      }
    }`)
	})

	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.TimeoutsList("enterprise", -1)
	assert.NoError(t, err)
	testutil.AssertOutput(t, b.String(), `PTYPE     TIMEOUT
web       10
worker    20
`)

	server.Mux.HandleFunc("/v2/apps/franklin/config/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
      "owner": "bedison",
      "app": "franklin",
      "values": [],
      "memory": {},
      "cpu": {},
      "tags": {},
      "registry": {},
      "created": "2014-01-01T00:00:00UTC",
      "updated": "2014-01-01T00:00:00UTC",
      "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
      }`)
	})
	b.Reset()

	err = cmdr.TimeoutsList("franklin", -1)
	assert.NoError(t, err)
	assert.Equal(t, b.String(), `Default (30 sec) or controlled by drycc controller.
`, "output")
}

func TestTimeoutsSet(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/apps/foo/config/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		if r.Method == "POST" {
			testutil.AssertBody(t, api.Config{
				Timeout: map[string]any{
					"web": "10",
				},
			}, r)
		}

		fmt.Fprintf(w, `{
      "owner": "jkirk",
      "app": "foo",
      "values": [],
      "memory": {},
      "cpu": {},
      "termination_grace_period": {
        "web": "10"
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

	err = cmdr.TimeoutsSet("foo", []string{"web=10"})
	assert.NoError(t, err)

	testutil.AssertOutput(t, testutil.StripProgress(b.String()), `Applying timeouts... done

PTYPE    TIMEOUT
web      10
`)

	server.Mux.HandleFunc("/v2/apps/franklin/config/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		if r.Method == "POST" {
			testutil.AssertBody(t, api.Config{
				Timeout: map[string]any{
					"web": "10",
				},
			}, r)
		}

		fmt.Fprintf(w, `{
      "owner": "bedison",
      "app": "franklin",
      "values": [],
      "memory": {},
      "cpu": {},
      "termination_grace_period": {
        "web": "10"
      },
      "tags": {},
      "registry": {},
      "created": "2014-01-01T00:00:00UTC",
      "updated": "2014-01-01T00:00:00UTC",
      "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
    }`)
	})
	b.Reset()

	err = cmdr.TimeoutsSet("franklin", []string{"web=10"})
	assert.NoError(t, err)

	testutil.AssertOutput(t, testutil.StripProgress(b.String()), `Applying timeouts... done

PTYPE    TIMEOUT
web      10
`)

	// with requests/timeout parameter
	server.Mux.HandleFunc("/v2/apps/jim/config/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		if r.Method == "POST" {
			testutil.AssertBody(t, api.Config{
				Timeout: map[string]any{
					"web":    "10",
					"worker": "100",
					"db":     "300",
				},
			}, r)
		}

		fmt.Fprintf(w, `{
      "owner": "foo",
      "app": "jim",
      "values": [],
      "memory": {},
      "cpu": {},
      "termination_grace_period": {
        "web": "10",
        "worker": "100",
        "db": "300"
      },
      "tags": {},
      "registry": {},
      "created": "2014-01-01T00:00:00UTC",
      "updated": "2014-01-01T00:00:00UTC",
      "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
    }`)
	})
	b.Reset()

	err = cmdr.TimeoutsSet("jim", []string{"web=10", "worker=100", "db=300"})
	assert.NoError(t, err)

	testutil.AssertOutput(t, testutil.StripProgress(b.String()), `Applying timeouts... done

PTYPE     TIMEOUT
db        300
web       10
worker    100
`)
}

func TestTimeoutsUnset(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/apps/foo/config/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		if r.Method == "POST" {
			testutil.AssertBody(t, api.Config{
				Timeout: map[string]any{
					"web": nil,
				},
			}, r)
		}

		fmt.Fprintf(w, `{
      "owner": "jkirk",
      "app": "foo",
      "values": [],
      "memory": {},
      "cpu": {},
      "termination_grace_period": {
        "web": 10
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

	err = cmdr.TimeoutsUnset("foo", []string{"web"})
	assert.NoError(t, err)

	testutil.AssertOutput(t, testutil.StripProgress(b.String()), `Applying timeouts... done

PTYPE    TIMEOUT
web      10
`)

	server.Mux.HandleFunc("/v2/apps/franklin/config/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		if r.Method == "POST" {
			testutil.AssertBody(t, api.Config{
				Timeout: map[string]any{
					"web": nil,
				},
			}, r)
		}

		fmt.Fprintf(w, `{
      "owner": "bedison",
      "app": "franklin",
      "values": [],
      "memory": {},
      "cpu": {},
      "termination_grace_period": {
        "web": 10
      },
      "tags": {},
      "registry": {},
      "created": "2014-01-01T00:00:00UTC",
      "updated": "2014-01-01T00:00:00UTC",
      "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
    }`)
	})
	b.Reset()

	err = cmdr.TimeoutsUnset("franklin", []string{"web"})
	assert.NoError(t, err)

	testutil.AssertOutput(t, testutil.StripProgress(b.String()), `Applying timeouts... done

PTYPE    TIMEOUT
web      10
`)
}
