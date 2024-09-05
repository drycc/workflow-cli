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

func TestCertsList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/certs/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"count": 4,
			"next": null,
			"previous": null,
			"results": [
				{	
					"app": "foo",
					"name": "test-example-com",
					"common_name": "test.example.com",
					"san": [
						"test.com",
						"example.com"
					],
					"domains": [
						"test.com",
						"example.com"
					],
					"created": "2016-06-09T00:00:00UTC",
					"updated": "2016-06-09T00:00:00UTC",
					"expires": "2014-11-10T00:00:00UTC",
					"fingerprint": "12:34:56:78:90"
				},
				{
					"app": "foo",
					"name": "test-drycc-com",
					"common_name": "test.drycc.com",
					"created": "2016-06-09T00:00:00UTC",
					"updated": "2016-06-09T00:00:00UTC",
					"expires": "2016-08-01T00:00:00UTC",
					"fingerprint": "ab:12:ab:12:ab"
				},
				{
					"app": "foo",
					"name": "test1",
					"common_name": "1.test.drycc.com",
					"expires": "2016-06-11T00:00:00UTC"
				},
				{
					"app": "foo",
					"name": "test2",
					"common_name": "2.test.drycc.com",
					"expires": "2018-01-01T00:00:00UTC"
				}
			]
		}`)
	})

	err = cmdr.CertsList("foo", -1)
	assert.NoError(t, err)

	assert.Equal(t, b.String(), `NAME                COMMON-NAME         EXPIRES        SAN                     DOMAINS              
test-example-com    test.example.com    10 Nov 2014    test.com,example.com    test.com,example.com    
test-drycc-com      test.drycc.com      1 Aug 2016     <none>                  <none>                  
test1               1.test.drycc.com    11 Jun 2016    <none>                  <none>                  
test2               2.test.drycc.com    1 Jan 2018     <none>                  <none>                  
`, "output")

	cf, server, err = testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	cmdr.ConfigFile = cf
	b.Reset()

	server.Mux.HandleFunc("/v2/apps/foo/certs/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"count": 0,
			"next": null,
			"previous": null,
			"results": []
		}`)
	})

	err = cmdr.CertsList("foo", -1)
	assert.NoError(t, err)

	assert.Equal(t, b.String(), "No certs\n", "output")
}

func TestCertsListLimit(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/certs/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"count": 4,
			"next": null,
			"previous": null,
			"results": [
				{
					"name": "test-example-com",
					"common_name": "test.example.com",
					"san": [
						"test.com",
						"drycc.com",
						"example.com"
					],
					"domains": [
						"test.com",
						"drycc.com",
						"example.com"
					],
					"created": "2016-06-09T00:00:00UTC",
					"updated": "2016-06-09T00:00:00UTC",
					"expires": "2014-11-10T00:00:00UTC",
					"fingerprint": "12:34:56:78:90"
				}
			]
		}`)
	})

	err = cmdr.CertsList("foo", 1)
	assert.NoError(t, err)

	assert.Equal(t, b.String(), `NAME                COMMON-NAME         EXPIRES        SAN                               DOMAINS                        
test-example-com    test.example.com    10 Nov 2014    test.com,drycc.com,example.com    test.com,drycc.com,example.com    
`, "output")

}

func TestCertsInfo(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/certs/test-example-com", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"app": "foo",
			"name": "test-example-com",
			"owner": "admin",
			"issuer": "testca",
			"subject": "testing",
			"common_name": "test.drycc.com",
			"created": "2016-06-09T00:00:00Z",
			"updated": "2016-06-09T00:00:00Z",
			"expires": "2016-06-09T00:00:00UTC",
			"starts": "2016-06-09T00:00:00UTC",
			"fingerprint": "ab:12:ab:12:ab",
			"san": [
				"test.com",
				"example.com"
			],
			"domains": [
				"test.com",
				"example.com"
			]
		}`)
	})

	err = cmdr.CertInfo("foo", "test-example-com")
	assert.NoError(t, err)
	assert.Equal(t, b.String(), `Name:                 test-example-com        
Common Name(s):       test.drycc.com          
Expires At:           9 Jun 2016              
Starts At:            9 Jun 2016              
Fingerprint:          ab:12:ab:12:ab          
Subject Alt Name:     test.com,example.com    
Issuer:               testca                  
Subject:              testing                 
                      
Connected Domains:    test.com,example.com    
Owner:                admin                   
Created:              2016-06-09T00:00:00Z    
Updated:              2016-06-09T00:00:00Z    
`, "output")

	server.Mux.HandleFunc("/v2/apps/foo/certs/test-drycc-com", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"name": "test-drycc-com"
		}`)
	})
	b.Reset()

	err = cmdr.CertInfo("foo", "test-drycc-com")
	assert.NoError(t, err)
	assert.Equal(t, b.String(), `Name:                 test-drycc-com    
Common Name(s):       <none>            
Expires At:           <none>            
Starts At:            <none>            
Fingerprint:          <none>            
Subject Alt Name:     <none>            
Issuer:               <none>            
Subject:              <none>            
                      
Connected Domains:    <none>            
Owner:                <none>            
Created:                                
Updated:                                
`, "output")
}

func TestCertsRemove(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/certs/test-example-com", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusNoContent)
	})

	err = cmdr.CertRemove("foo", "test-example-com")
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Removing test-example-com... done\n", "output")
}

func TestCertsAttach(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/certs/test-example-com/domain/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		testutil.AssertBody(t, api.CertAttachRequest{Domain: "drycc.com"}, r)
		w.WriteHeader(http.StatusCreated)
	})

	err = cmdr.CertAttach("foo", "test-example-com", "drycc.com")
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Attaching certificate test-example-com to domain drycc.com... done\n", "output")
}

func TestCertsDetach(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/certs/test-example-com/domain/drycc.com", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusNoContent)
	})

	err = cmdr.CertDetach("foo", "test-example-com", "drycc.com")
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Detaching certificate test-example-com from domain drycc.com... done\n", "output")
}

func TestCertsAdd(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/certs/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		testutil.AssertBody(t, api.CertCreateRequest{Certificate: "cert", Key: "key", Name: "testcert"}, r)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "{}")
	})

	keyFile, err := os.CreateTemp("", "drycc-cli-unit-test-key")
	assert.NoError(t, err)
	_, err = keyFile.Write([]byte("key"))
	assert.NoError(t, err)
	keyFile.Close()

	certFile, err := os.CreateTemp("", "drycc-cli-unit-test-cert")
	assert.NoError(t, err)
	_, err = certFile.Write([]byte("cert"))
	assert.NoError(t, err)
	certFile.Close()

	err = cmdr.CertAdd("foo", certFile.Name(), keyFile.Name(), "testcert")
	assert.NoError(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Adding SSL endpoint... done\n", "output")
}
