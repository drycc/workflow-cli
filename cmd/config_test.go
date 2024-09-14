package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/workflow-cli/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {
	t.Parallel()

	_, err := parseConfig([]string{"FOO=bar", "CAR star"})
	assert.NotEqual(t, err, nil, "config")

	actual, err := parseConfig([]string{"FOO=bar"})
	assert.NoError(t, err)
	assert.Equal(t, actual, api.ConfigValues{"FOO": "bar"}, "map")

	actual, err = parseConfig([]string{"FOO="})
	assert.NoError(t, err)
	assert.Equal(t, actual, api.ConfigValues{"FOO": ""}, "map")
}

func TestParseSSHKey(t *testing.T) {
	t.Parallel()

	_, err := parseSSHKey("foobar")
	assert.NotEqual(t, err, "bogus key")

	validSSHKey := "-----BEGIN OPENSSH PRIVATE KEY-----"

	actual, err := parseSSHKey(validSSHKey)
	assert.NoError(t, err)
	assert.Equal(t, actual, validSSHKey, "plain key")

	encodedSSHKey := "LS0tLS1CRUdJTiBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0="

	actual, err = parseSSHKey(encodedSSHKey)
	assert.NoError(t, err)
	assert.Equal(t, actual, validSSHKey, "base64 key")

	keyFile, err := os.CreateTemp("", "drycc-cli-unit-test-sshkey")
	assert.NoError(t, err)
	defer os.Remove(keyFile.Name())
	_, err = keyFile.Write([]byte(validSSHKey))
	assert.NoError(t, err)
	keyFile.Close()

	actual, err = parseSSHKey(keyFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, actual, validSSHKey, "key path")
}

func TestFormatConfig(t *testing.T) {
	t.Parallel()

	testMap := map[string]interface{}{
		"TEST":  "testing",
		"NCC":   1701,
		"TRUE":  false,
		"FLOAT": 12.34,
	}

	testOut := formatConfig(testMap)
	assert.Equal(t, testOut, `FLOAT=12.34
NCC=1701
TEST=testing
TRUE=false
`, "output")
}

func TestSortKeys(t *testing.T) {
	test := map[string]interface{}{
		"d": nil,
		"b": nil,
		"c": nil,
		"a": nil,
	}

	assert.Equal(t, *sortKeys(test), []string{"a", "b", "c", "d"}, "map")
}

func TestConfigList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/apps/foo/config/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
    "owner": "jkirk",
    "app": "foo",
    "values": {
        "TEST":  "testing",
        "NCC":   "1701",
        "TRUE":  "false",
        "FLOAT": "12.34"
    },
    "typed_values": {
		"web": {
            "PORT":  "9000"
		}
    },
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

	err = cmdr.ConfigList("foo", "")
	assert.NoError(t, err)

	assert.Equal(t, b.String(), `PTYPE    NAME     VALUE   
N/A      FLOAT    12.34      
N/A      NCC      1701       
N/A      TEST     testing    
N/A      TRUE     false      
web      PORT     9000       
`, "output")

	b.Reset()
	err = cmdr.ConfigList("foo", "web")
	assert.NoError(t, err)
	assert.Equal(t, b.String(), `PTYPE    NAME     VALUE   
N/A      FLOAT    12.34      
N/A      NCC      1701       
N/A      TEST     testing    
N/A      TRUE     false      
web      PORT     9000       
`, "output")

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
				Values: map[string]interface{}{
					"TRUE":    "false",
					"SSH_KEY": "LS0tLS1CRUdJTiBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0=",
				},
			}, r)
		}

		fmt.Fprintf(w, `{
	"owner": "jkirk",
	"app": "foo",
	"values": {
			"TEST":  "testing",
			"NCC":   "1701",
			"TRUE":  "false",
			"SSH_KEY": "LS0tLS1CRUdJTiBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0=",
			"FLOAT": "12.34"
	},
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

	err = cmdr.ConfigSet("foo", "", []string{"TRUE=false", "SSH_KEY=-----BEGIN OPENSSH PRIVATE KEY-----"}, "yes")
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Creating config... done

PTYPE    NAME       VALUE                                            
N/A      FLOAT      12.34                                               
N/A      NCC        1701                                                
N/A      SSH_KEY    LS0tLS1CRUdJTiBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0=    
N/A      TEST       testing                                             
N/A      TRUE       false                                               
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
				Values: map[string]interface{}{
					"FOO": nil,
				},
			}, r)
		}

		fmt.Fprintf(w, `{
	"owner": "jkirk",
	"app": "foo",
	"values": {
			"TEST":  "testing",
			"NCC":   "1701",
			"TRUE":  "false",
			"FLOAT": "12.34"
	},
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

	err = cmdr.ConfigUnset("foo", "", []string{"FOO"}, "yes")
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Removing config... done

PTYPE    NAME     VALUE   
N/A      FLOAT    12.34      
N/A      NCC      1701       
N/A      TEST     testing    
N/A      TRUE     false      
`, "output")
}

func TestConfigUnsetTypedValues(t *testing.T) {
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
				TypedValues: map[string]api.ConfigValues{
					"web": {
						"FOO": nil,
					},
				},
			}, r)
		}

		fmt.Fprintf(w, `{
	"owner": "jkirk",
	"app": "foo",
	"values": {
		"RELEASE_VERSION": "v1"
	},
	"typed_values": {
		"web": {
			"TEST":  "testing",
			"NCC":   "1701",
			"TRUE":  "false",
			"FLOAT": "12.34"
		}
	},
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

	err = cmdr.ConfigUnset("foo", "web", []string{"FOO"}, "")
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Removing config... done

PTYPE    NAME               VALUE   
N/A      RELEASE_VERSION    v1         
web      FLOAT              12.34      
web      NCC                1701       
web      TEST               testing    
web      TRUE               false      
`, "output")
}
