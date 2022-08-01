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

type parseLimitCase struct {
	Input         string
	Key           string
	Value         string
	ExpectedError bool
	ExpectedMsg   string
}

func TestParseLimit(t *testing.T) {
	t.Parallel()

	var errorHint = ` doesn't fit format type=#unit or type=#
Examples: web=2G worker=500M db=1G`

	cases := []parseLimitCase{
		{"web=2G", "web", "2G", false, ""},
		{"web=2", "web", "2", false, ""},
		{"web=100m", "web", "100m", false, ""},
		{"web1=2G", "web1", "2G", false, ""},
		{"web-server=2G", "web-server", "2G", false, ""},
		{"web-server1=2G", "web-server1", "2G", false, ""},
		{"web=0.1", "", "", true, "web=0.1" + errorHint},
		{"web=.123", "", "", true, "web=.123" + errorHint},
		{"=1", "", "", true, "=1" + errorHint},
		{"web=", "", "", true, "web=" + errorHint},
		{"1=", "", "", true, "1=" + errorHint},
		{"web=G", "", "", true, "web=G" + errorHint},
		{"web-=2G", "", "", true, "web-=2G" + errorHint},
		{"-web=2G", "", "", true, "-web=2G" + errorHint},
		{"Web=2G", "", "", true, "Web=2G" + errorHint},
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
		{[]string{"web=1G", "worker=2"}, map[string]interface{}{"web": "1G", "worker": "2"}, false, ""},
		{[]string{"foo=", "web=1G"}, nil, true, `foo= doesn't fit format type=#unit or type=#
Examples: web=2G worker=500M db=1G`},
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
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/apps/enterprise/config/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"owner": "jkirk",
			"app": "enterprise",
			"values": {},
			"memory": {
				"web": "2G",
				"worker": "1G",
				"db": "1000M"
			},
			"cpu": {
				"web": "2",
				"worker": "1",
				"db": "500m"
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

	err = cmdr.LimitsList("enterprise")
	assert.NoError(t, err)
	assert.Equal(t, b.String(), `=== enterprise Limits

--- Memory
db         1000M
web        2G
worker     1G

--- CPU
db         500m
web        2
worker     1
`, "output")

	server.Mux.HandleFunc("/v2/apps/franklin/config/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"owner": "bedison",
			"app": "franklin",
			"values": {},
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

	err = cmdr.LimitsList("franklin")
	assert.NoError(t, err)
	assert.Equal(t, b.String(), `=== franklin Limits

--- Memory

--- CPU
`, "output")
}

func TestLimitsSet(t *testing.T) {
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
				CPU: map[string]interface{}{
					"web": "100m",
				},
			}, r)
		}

		fmt.Fprintf(w, `{
			"owner": "jkirk",
			"app": "foo",
			"values": {},
			"memory": {"web": "128M"},
			"cpu": {
				"web": "100m"
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

	err = cmdr.LimitsSet("foo", []string{"web=100m"}, []string{})
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Applying limits... done

=== foo Limits

--- Memory
web     128M

--- CPU
web     100m
`, "output")

	server.Mux.HandleFunc("/v2/apps/franklin/config/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		if r.Method == "POST" {
			testutil.AssertBody(t, api.Config{
				Memory: map[string]interface{}{
					"web": "1G",
				},
			}, r)
		}

		fmt.Fprintf(w, `{
			"owner": "bedison",
			"app": "franklin",
			"values": {},
			"memory": {
				"web": "1G"
			},
			"cpu": {
				"web": "1"
			},
			"tags": {},
			"registry": {},
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})
	b.Reset()

	err = cmdr.LimitsSet("franklin", []string{}, []string{"web=1G"})
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Applying limits... done

=== franklin Limits

--- Memory
web     1G

--- CPU
web     1
`, "output")

	// with requests/limit parameter
	server.Mux.HandleFunc("/v2/apps/jim/config/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		if r.Method == "POST" {
			testutil.AssertBody(t, api.Config{
				Memory: map[string]interface{}{
					"web":    "2000M",
					"worker": "3G",
					"db":     "5G",
				},
			}, r)
		}

		fmt.Fprintf(w, `{
			"owner": "foo",
			"app": "jim",
			"values": {},
			"memory": {
				"web": "2000M",
				"worker": "3G",
				"db": "5G"
			},
			"cpu": {
				"web": "1",
				"worker": "1",
				"db": "5"
			},
			"tags": {},
			"registry": {},
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})
	b.Reset()

	err = cmdr.LimitsSet("jim", []string{}, []string{"web=2000M", "worker=3G", "db=5G"})
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Applying limits... done

=== jim Limits

--- Memory
db         5G
web        2000M
worker     3G

--- CPU
db         5
web        1
worker     1
`, "output")

	// with requests/limit parameter
	server.Mux.HandleFunc("/v2/apps/phew/config/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		if r.Method == "POST" {
			testutil.AssertBody(t, api.Config{
				CPU: map[string]interface{}{
					"web":    "2",
					"worker": "300m",
					"db":     "5",
				},
			}, r)
		}

		fmt.Fprintf(w, `{
			"owner": "foo",
			"app": "jim",
			"values": {},
			"cpu": {
				"web": "2",
				"worker": "300m",
				"db": "5"
			},
			"memory": {
				"web": "1G",
				"worker": "1G",
				"db": "1G"
			},
			"tags": {},
			"registry": {},
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})
	b.Reset()

	err = cmdr.LimitsSet("phew", []string{"web=2", "worker=300m", "db=5"}, []string{})
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Applying limits... done

=== phew Limits

--- Memory
db         1G
web        1G
worker     1G

--- CPU
db         5
web        2
worker     300m
`, "output")
}

func TestLimitsUnset(t *testing.T) {
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
				Memory: map[string]interface{}{
					"web": nil,
				},
			}, r)
		}

		fmt.Fprintf(w, `{
			"owner": "jkirk",
			"app": "foo",
			"values": {},
			"memory": {
				"web": "128M"
			},
			"cpu": {
				"web": "125m"
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

	err = cmdr.LimitsUnset("foo", []string{}, []string{"web"})
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Applying limits... done

=== foo Limits

--- Memory
web     128M

--- CPU
web     125m
`, "output")

	server.Mux.HandleFunc("/v2/apps/franklin/config/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		if r.Method == "POST" {
			testutil.AssertBody(t, api.Config{
				CPU: map[string]interface{}{
					"web": nil,
				},
				Memory: map[string]interface{}{
					"web": nil,
				},
			}, r)
		}

		fmt.Fprintf(w, `{
			"owner": "bedison",
			"app": "franklin",
			"values": {},
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

	err = cmdr.LimitsUnset("franklin", []string{"web"}, []string{"web"})
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Applying limits... done

=== franklin Limits

--- Memory

--- CPU
`, "output")
}
