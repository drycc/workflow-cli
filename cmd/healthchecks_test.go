package cmd

import (
	"bytes"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/controller-sdk-go/api"
)

func TestPrintHealthCheck(t *testing.T) {
	var b bytes.Buffer
	testHealthCheck := api.Healthchecks{}
	printHealthCheck(&b, testHealthCheck)
	assert.Equal(t, b.String(), "--- Liveness\nNo liveness probe configured.\n\n--- Readiness\nNo readiness probe configured.\n", "healthcheck")
	b.Reset()
	testHealthCheck["livenessProbe"] = &api.Healthcheck{}
	testHealthCheck["readinessProbe"] = &api.Healthcheck{}
	printHealthCheck(&b, testHealthCheck)
	assert.Equal(t, b.String(), "--- Liveness\nInitial Delay (seconds): 0\nTimeout (seconds): 0\nPeriod (seconds): 0\nSuccess Threshold: 0\nFailure Threshold: 0\nExec Probe: N/A\nHTTP GET Probe: N/A\nTCP Socket Probe: N/A\n\n--- Readiness\nInitial Delay (seconds): 0\nTimeout (seconds): 0\nPeriod (seconds): 0\nSuccess Threshold: 0\nFailure Threshold: 0\nExec Probe: N/A\nHTTP GET Probe: N/A\nTCP Socket Probe: N/A\n", "healthcheck")
}
