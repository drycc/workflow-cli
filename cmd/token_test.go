package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/drycc/workflow-cli/pkg/testutil"
)

func TestTokensList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/tokens/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"count": 2,
			"next": null,
			"previous": null,
			"results": [
				{
					"uuid": "f71e3b18-e702-409e-bd7f-8fb0a66d7b12",
					"owner": "test",
					"alias": "",
					"fuzzy_key": "c8e74fa4cbf...e4954d602ec5ed19ba",
					"created": "2023-04-18T00:00:00UTC",
					"updated": "2023-04-19T00:00:00UTC"
				},
				{
					"uuid": "f71e3b18-e702-499e-bd7f-8fb0a66d7b12",
					"owner": "test",
					"alias": "test",
					"fuzzy_key": "c8e74fa4cbf...e4954d60cec5ed19ba",
					"created": "2023-04-18T10:00:00UTC",
					"updated": "2023-04-19T12:00:00UTC"
				}
			]
		}`)
	})

	err = cmdr.TokensList(-1)
	assert.NoError(t, err)

	assert.Equal(t, b.String(), `UUID                                    OWNER    ALIAS     KEY                                 CREATE                    UPDATED                
f71e3b18-e702-409e-bd7f-8fb0a66d7b12    test     <none>    c8e74fa4cbf...e4954d602ec5ed19ba    2023-04-18T00:00:00UTC    2023-04-19T00:00:00UTC    
f71e3b18-e702-499e-bd7f-8fb0a66d7b12    test     test      c8e74fa4cbf...e4954d60cec5ed19ba    2023-04-18T10:00:00UTC    2023-04-19T12:00:00UTC    
`, "output")
}

func TestTokenDelete(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/tokens/f71e3b18-e702-499e-bd7f-8fb0a66d7b12/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(204)
	})

	err = cmdr.TokensRemove("f71e3b18-e702-499e-bd7f-8fb0a66d7b12", "yes")

	assert.NoError(t, err)
	assert.Equal(t, b.String(), "done\n")
}
