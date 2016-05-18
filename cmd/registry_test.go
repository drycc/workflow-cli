package cmd

import "testing"

func TestParseInfo(t *testing.T) {
	t.Parallel()

	// keys can only be username or password
	goodKeys := []string{
		"username=bob",
		"password=isyouruncle",
		// regression test for passwords with equals signs, such as a gcr.io token
		"password=ihaveanequalssign=",
	}

	for _, key := range goodKeys {
		if _, _, err := parseInfo(key); err != nil {
			t.Errorf("failed parsing valid keys, got (%s)", err)
		}
	}

	badKey := "usrname=bob"
	if _, _, err := parseInfo(badKey); err == nil {
		t.Errorf("failed erroring on bad key '%s'", badKey)
	}
}
