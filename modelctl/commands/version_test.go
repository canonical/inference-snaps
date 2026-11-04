package commands

import (
	"testing"
)

func TestCleanVersionString(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected string
	}{
		"empty version": {
			input:    "",
			expected: "unset",
		},
		"normal version": {
			input:    "1.2.3",
			expected: "1.2.3",
		},
		"devel version": {
			input:    "(devel)",
			expected: "(devel)",
		},
		"dirty version": {
			input:    "1.0.0+dirty",
			expected: "1.0.0+dirty",
		},
	}

	for testName, testData := range tests {
		t.Run(testName, func(t *testing.T) {
			actual := cleanVersionString(testData.input)
			if actual != testData.expected {
				t.Errorf("cleanVersionString(%q) returned %q, expected %q", testData.input, actual, testData.expected)
			}
		})
	}
}
