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

type parseInfoCase struct {
	Input         string
	Key           string
	Value         string
	ExpectedError bool
	ExpectedMsg   string
}

func TestParseInfo(t *testing.T) {
	t.Parallel()

	cases := []parseInfoCase{
		{"username=test", "username", "test", false, ""},
		{"password=test=", "password", "test=", false, ""},
		{"test=1", "", "", true, `test is invalid. Valid keys are "username" or "password"`},
		{"test", "", "", true, `test is invalid. Must be in format key=value
Examples: username=bob password=s3cur3pw1`},
		{"test=", "", "", true, `test= is invalid. Must be in format key=value
Examples: username=bob password=s3cur3pw1`},
		{"=test", "", "", true, `=test is invalid. Must be in format key=value
Examples: username=bob password=s3cur3pw1`},
	}

	for _, check := range cases {
		key, value, err := parseInfo(check.Input)
		if check.ExpectedError {
			assert.Equal(t, err.Error(), check.ExpectedMsg, "error")
		} else {
			assert.NoError(t, err)
			assert.Equal(t, key, check.Key, "key")
			assert.Equal(t, value, check.Value, "value")
		}
	}
}

type parseInfosCase struct {
	Input         []string
	ExpectedMap   map[string]interface{}
	ExpectedError bool
	ExpectedMsg   string
}

func TestParseInfos(t *testing.T) {
	t.Parallel()

	cases := []parseInfosCase{
		{[]string{"username=test", "password=abc123"}, map[string]interface{}{"username": "test", "password": "abc123"}, false, ""},
		{[]string{"foo=", "true=false"}, nil, true, `foo= is invalid. Must be in format key=value
Examples: username=bob password=s3cur3pw1`},
	}

	for _, check := range cases {
		actual, err := parseInfos(check.Input)
		if check.ExpectedError {
			assert.Equal(t, err.Error(), check.ExpectedMsg, "error")
		} else {
			assert.NoError(t, err)
			assert.Equal(t, actual, check.ExpectedMap, "map")
		}
	}
}

func TestRegistryList(t *testing.T) {
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
			"values": {},
			"memory": {},
			"cpu": {},
			"tags": {},
			"registry": {
				"username": "jkirk",
				"password": "ncc1701"
			},
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})

	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.RegistryList("enterprise")
	assert.NoError(t, err)
	assert.Equal(t, b.String(), `UUID                                    OWNER    KEY         VALUE   
de1bf5b5-4a72-4f94-a10c-d2a3741cdf75    jkirk    password    ncc1701    
de1bf5b5-4a72-4f94-a10c-d2a3741cdf75    jkirk    username    jkirk      
`, "output")
}

func TestRegistrySet(t *testing.T) {
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
				Registry: map[string]interface{}{
					"username": "jkirk",
					"password": "ncc1701",
				},
			}, r)
		}

		fmt.Fprintf(w, `{
			"owner": "jkirk",
			"app": "foo",
			"values": {},
			"memory": {},
			"cpu": {},
			"registry": {
				"username": "jkirk",
				"password": "ncc1701"
			},
			"registry": {},
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})

	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.RegistrySet("foo", []string{"username=jkirk", "password=ncc1701"})
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Applying registry information... done

UUID                                    OWNER    KEY         VALUE   
de1bf5b5-4a72-4f94-a10c-d2a3741cdf75    jkirk    password    ncc1701    
de1bf5b5-4a72-4f94-a10c-d2a3741cdf75    jkirk    username    jkirk      
`, "output")
}

func TestRegistryUnset(t *testing.T) {
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
				Registry: map[string]interface{}{
					"username": nil,
					"password": nil,
				},
			}, r)
		}

		fmt.Fprintf(w, `{
			"owner": "jkirk",
			"app": "foo",
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

	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.RegistryUnset("foo", []string{"username", "password"})
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Applying registry information... done

No registrys found in foo app.
`, "output")
}
