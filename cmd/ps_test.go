package cmd

import "testing"

func TestParseType(t *testing.T) {
	t.Parallel()

	// test RC pod name
	appID := "earthy-underdog"
	rcPod := "earthy-underdog-v2-cmd-8yngj"
	psType, psName := parseType(rcPod, appID)
	if psType != "cmd" || psName != rcPod {
		t.Errorf("type was not cmd (got %s) or psName was not %s (got %s)", psType, rcPod, psName)
	}

	// test Deployment pod name - they are longer due to hash
	appID = "nonfat-yearbook"
	deployPod := "nonfat-yearbook-cmd-2180299075-7na91"
	psType, psName = parseType(deployPod, appID)
	if psType != "cmd" || psName != deployPod {
		t.Errorf("type was not cmd (got %s) or psName was not %s (got %s)", psType, deployPod, psName)
	}


	// test type by itself
	psType, psName = parseType("cmd", "fake")
	if psType != "cmd" || psName != "" {
		t.Error("type was not cmd")
	}

}
