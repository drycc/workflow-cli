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

func TestParseConfig(t *testing.T) {
	t.Parallel()

	_, err := parseConfig("", "", []string{"FOO=bar", "CAR star"})
	assert.NotEqual(t, err, nil, "config")

	actual, err := parseConfig("", "", []string{"FOO=bar"})
	assert.NoError(t, err)
	assert.Equal(t, actual, []api.ConfigValue{{ConfigVar: api.ConfigVar{Name: "FOO", Value: "bar"}}}, "map")

	actual, err = parseConfig("", "", []string{"FOO="})
	assert.NoError(t, err)
	assert.Equal(t, actual, []api.ConfigValue{{ConfigVar: api.ConfigVar{Name: "FOO", Value: ""}}}, "map")
}

func TestFormatConfig(t *testing.T) {
	t.Parallel()

	testMap := []api.ConfigValue{
		{
			Ptype: "web",
			Group: "",
			ConfigVar: api.ConfigVar{
				Name:  "TEST",
				Value: "testing",
			},
		}, {
			Ptype: "web",
			Group: "",
			ConfigVar: api.ConfigVar{
				Name:  "NCC",
				Value: "1701",
			},
		}, {
			Ptype: "web",
			Group: "",
			ConfigVar: api.ConfigVar{
				Name:  "TRUE",
				Value: false,
			},
		}, {
			Ptype: "web",
			Group: "",
			ConfigVar: api.ConfigVar{
				Name:  "FLOAT",
				Value: 12.34,
			},
		},
	}

	testOut := formatConfig(testMap)
	assert.Equal(t, testOut, `FLOAT=12.34
NCC=1701
TEST=testing
TRUE=false
`, "output")
}

func TestSortKeys(t *testing.T) {
	test := map[string]any{
		"d": nil,
		"b": nil,
		"c": nil,
		"a": nil,
	}

	assert.Equal(t, *sortKeys(test), []string{"a", "b", "c", "d"}, "map")
}

func TestConfigSet(t *testing.T) {
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
				Values: []api.ConfigValue{
					{
						Ptype: "web",
						ConfigVar: api.ConfigVar{
							Name:  "TRUE",
							Value: "false",
						},
					},
					{
						Ptype: "web",
						ConfigVar: api.ConfigVar{
							Name:  "DEBUG",
							Value: "true",
						},
					},
				},
			}, r)
		}

		fmt.Fprintf(w, `{
	"owner": "jkirk",
	"app": "foo",
	"values": [
	  {
	    "ptype": "web",
	    "name": "TRUE",
		"value": "false"
	  },
	  {
	    "ptype": "web",
	    "name": "DEBUG",
		"value": "true"
	  }
	],
	"memory": {},
	"cpu": {},
	"tags": {},
	"registry": {},
	"created": "2014-01-01T00:00:00UTC",
	"updated": "2014-01-01T00:00:00UTC",
	"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
}`)
	})

	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.ConfigSet("foo", "web", "", []string{"TRUE=false", "DEBUG=true"}, "yes")
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Creating config... done

`, "output")
}

func TestConfigUnset(t *testing.T) {
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
				Values: []api.ConfigValue{
					{
						Ptype: "web",
						ConfigVar: api.ConfigVar{
							Name:  "FOO",
							Value: nil,
						},
					},
				},
			}, r)
		}

		fmt.Fprintf(w, `{
	"owner": "jkirk",
	"app": "foo",
	"values": [
	  {
	    "ptype": "web",
	    "name": "FLOAT",
		"value": "12.34"
	  },
	  {
	    "ptype": "web",
	    "name": "NCC",
		"value": "1701"
	  }
	],
	"memory": {},
	"cpu": {},
	"tags": {},
	"registry": {},
	"created": "2014-01-01T00:00:00UTC",
	"updated": "2014-01-01T00:00:00UTC",
	"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
}`)
	})

	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.ConfigUnset("foo", "web", "", []string{"FOO"}, "yes")
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Removing config... done

`, "output")
}
