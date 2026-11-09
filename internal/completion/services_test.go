package completion

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/drycc/workflow-cli/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestServiceCompletionUsesNamesForDuplicatePtype(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	assert.NoError(t, err)
	defer server.Close()

	server.Mux.HandleFunc("/v2/apps/foo/services/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprint(w, `{
    "count": 2,
    "next": null,
    "previous": null,
    "results": [
        {"name": "foo-worker-a", "ptype": "worker", "ports": []},
        {"name": "foo-worker-b", "ptype": "worker", "ports": []},
        {"name": "other-web", "ptype": "web", "ports": []}
    ]
}`)
	})

	appID := "foo"
	completion := ServiceCompletion{AppID: &appID, ArgsLen: 0, ConfigFile: &cf}
	results, _ := completion.CompletionFunc(nil, nil, "foo-worker")
	assert.ElementsMatch(t, []string{"foo-worker-a", "foo-worker-b"}, results)
}
