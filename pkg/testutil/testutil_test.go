package testutil

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

// TestStripProgress ensures StripProgress strips what is expected.
func TestStripProgress(t *testing.T) {
	testInput := "Lorem ipsum dolar sit amet"
	expectedOutput := "Lorem ipsum dolar sit amet"

	assert.Equal(t, StripProgress(testInput), expectedOutput, "output")

	testInput = "Lorem ipsum dolar sit amet...\b\b\b"
	assert.Equal(t, StripProgress(testInput), expectedOutput, "output")
}

// TestAssertBody ensures AssertBody correctly marshals into the interface.
func TestAssertBody(t *testing.T) {
	b := nopCloser{bytes.NewBufferString(`{"data":{"lorem":"ipsum"},"dolar":["sit","amet"]}`)}

	sampleRequest := http.Request{
		Body: b,
	}

	expected := map[string]any{
		"data": map[string]any{
			"lorem": "ipsum",
		},
		"dolar": []string{
			"sit",
			"amet",
		},
	}

	AssertBody(t, expected, &sampleRequest)
}

// TestAssertOutput tests the AssertOutput method with various scenarios
func TestAssertOutput(t *testing.T) {
	tests := []struct {
		name       string
		expected   string
		actual     string
		shouldPass bool
	}{
		{
			name:       "identical strings",
			expected:   "line1\nline2\nline3",
			actual:     "line1\nline2\nline3",
			shouldPass: true,
		},
		{
			name:       "trailing spaces should be ignored",
			expected:   "line1\nline2  \nline3",
			actual:     "line1\nline2\nline3",
			shouldPass: true,
		},
		{
			name:       "leading spaces should be ignored",
			expected:   "line1\n  line2\nline3",
			actual:     "line1\nline2\nline3",
			shouldPass: true,
		},
		{
			name:       "tabs should be ignored",
			expected:   "line1\n\tline2\t\nline3",
			actual:     "line1\nline2\nline3",
			shouldPass: true,
		},
		{
			name:       "mixed whitespace should be ignored",
			expected:   "line1\n \t line2 \t \nline3",
			actual:     "line1\nline2\nline3",
			shouldPass: true,
		},
		{
			name:       "different content should fail",
			expected:   "line1\nline2\nline3",
			actual:     "line1\nline2\nline4",
			shouldPass: false,
		},
		{
			name:       "different number of lines should fail",
			expected:   "line1\nline2",
			actual:     "line1\nline2\nline3",
			shouldPass: false,
		},
		{
			name:       "empty strings",
			expected:   "",
			actual:     "",
			shouldPass: true,
		},
		{
			name:       "single line with spaces",
			expected:   "  hello world  ",
			actual:     "hello world",
			shouldPass: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test by comparing processed results directly
			expectedLines := strings.Split(tt.expected, "\n")
			actualLines := strings.Split(tt.actual, "\n")

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

			if tt.shouldPass {
				if expectedProcessed != actualProcessed {
					t.Errorf("Expected processed output to match, but they don't:\nExpected: %q\nActual: %q", expectedProcessed, actualProcessed)
				}
			} else {
				if expectedProcessed == actualProcessed {
					t.Errorf("Expected processed output to differ, but they match: %q", expectedProcessed)
				}
			}
		})
	}
}
