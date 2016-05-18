package cmd

import "testing"

func TestParseInfo(t *testing.T) {
	t.Parallel()

	// Regression test for passwords with equals signs, such as a gcr.io token
	key := `password=ihaveanequalssign=`
	if _, _, err := parseInfo(key); err != nil {
		t.Errorf("failed to parse valid token with equals sign: got (%s)", err)
	}
}
