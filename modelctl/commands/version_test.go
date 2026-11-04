package commands

import (
	"testing"

	"github.com/canonical/inference-snaps/modelctl/v2/commands/common"
	"github.com/canonical/inference-snaps/modelctl/v2/internal/snap"
)

func TestGetVersionData(t *testing.T) {
	tests := map[string]struct {
		version     string
		expectedCli string
	}{
		"version set": {
			version:     "v2.0.0",
			expectedCli: "v2.0.0",
		},
		"version unset": {
			version:     "",
			expectedCli: "unset",
		},
	}

	for testName, testData := range tests {
		t.Run(testName, func(t *testing.T) {
			cmd := &versionCommand{
				Context: &common.Context{
					Snap:    snap.Mock(),
					Version: testData.version,
				},
			}

			actual := cmd.getVersionData()

			if actual.Cli != testData.expectedCli {
				t.Errorf("Cli = %q, expected %q", actual.Cli, testData.expectedCli)
			}
			if actual.Snap != "mock-version" {
				t.Errorf("Snap = %q, expected %q", actual.Snap, "mock-version")
			}
		})
	}
}

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
