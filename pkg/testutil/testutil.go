// Package testutil provides utilities for testing the workflow CLI.
package testutil

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/workflow-cli/pkg/settings"
)

// TestServer represents a test HTTP server along with a path to a config file
type TestServer struct {
	Server *httptest.Server
	Mux    *http.ServeMux
}

// NewTestServer sets up a test HTTP Server without a configuration file.
func NewTestServer() *TestServer {
	// test server
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	return &TestServer{
		Server: server,
		Mux:    mux,
	}
}

// Close closes the test HTTP server.
func (t *TestServer) Close() {
	t.Server.Close()
}

// NewTestServerAndClient sets up a test HTTP Server along with a configuration file to talk to it
func NewTestServerAndClient() (string, *TestServer, error) {
	server := NewTestServer()

	name, err := os.MkdirTemp("", "client")
	if err != nil {
		server.Close()
		return "", nil, err
	}

	filename := filepath.Join(name, "test.json")

	client, err := drycc.New(false, server.Server.URL, "")
	if err != nil {
		server.Close()
		return "", nil, err
	}

	config := settings.Settings{
		Username: "test",
		Client:   client,
	}

	filename, err = config.Save(filename)
	if err != nil {
		server.Close()
		return "", nil, err
	}
	return filename, server, nil
}

// AssertBody asserts the value of the body of a request.
func AssertBody(t *testing.T, expected any, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatal(err)
	}

	value, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != string(value) {
		t.Errorf("Expected body '%s' actually got '%s'\n", value, body)
	}
}

// StripProgress strips the output from the progress method
func StripProgress(input string) string {
	first := strings.Index(input, "\b")
	// If \b charecter not part of string
	if first == -1 {
		return input
	}
	last := strings.LastIndex(input, "\b")

	// return string without \b or the characters it deletes.
	return input[:first-(last-first+1)] + input[last+1:]
}

// SetHeaders sets standard headers for requests
func SetHeaders(w http.ResponseWriter) {
	w.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)
}

// AssertOutput compares multi-line text after trimming spaces and tabs from each line
func AssertOutput(t *testing.T, expected string, actual string) {
	t.Helper()
	expectedLines := strings.Split(strings.Trim(expected, "\n"), "\n")
	actualLines := strings.Split(strings.Trim(actual, "\n"), "\n")

	// Trim spaces and tabs from each line
	for i, line := range expectedLines {
		expectedLines[i] = strings.Trim(line, " \t")
	}
	for i, line := range actualLines {
		actualLines[i] = strings.Trim(line, " \t")
	}

	// Join lines back with newlines
	expectedProcessed := strings.Join(expectedLines, "\n")
	actualProcessed := strings.Join(actualLines, "\n")

	if expectedProcessed != actualProcessed {
		t.Errorf("Expected output:\n%s\n\nActual output:\n%s", expectedProcessed, actualProcessed)
	}
}
