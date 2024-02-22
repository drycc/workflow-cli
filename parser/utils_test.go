package parser

import "testing"

func TestSafeGet(t *testing.T) {
	t.Parallel()

	expected := "foo"

	test := make(map[string]interface{}, 1)
	test["test"] = "foo"

	actual := safeGetString(test, "test")

	if expected != actual {
		t.Errorf("Expected %s, Got %s", expected, actual)
	}
}

func TestSafeGetNil(t *testing.T) {
	t.Parallel()

	expected := ""

	test := make(map[string]interface{}, 1)
	test["test"] = nil

	actual := safeGetString(test, "test")

	if expected != actual {
		t.Errorf("Expected %s, Got %s", expected, actual)
	}
}

func TestSafeGetInt(t *testing.T) {
	t.Parallel()

	expected := 1

	test := make(map[string]interface{}, 1)
	test["test"] = "1"

	actual := safeGetInt(test, "test")

	if expected != actual {
		t.Errorf("Expected %d, Got %d", expected, actual)
	}

	if actual = safeGetInt(test, "foo"); actual != 0 {
		t.Errorf("Expected 0, Got %d", actual)
	}
}

func TestPrintHelp(t *testing.T) {
	t.Parallel()

	usage := ""

	if !printHelp([]string{"ps", "--help"}, usage) {
		t.Error("Expected true")
	}

	if !printHelp([]string{"ps", "-h"}, usage) {
		t.Error("Expected true")
	}

	if printHelp([]string{"ps"}, usage) {
		t.Error("Expected false")
	}

	if printHelp([]string{"ps", "--foo"}, usage) {
		t.Error("Expected false")
	}
}
