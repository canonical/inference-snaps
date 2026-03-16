package common

import (
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/engines"
)

func TestFillIncompatibilityIssues(t *testing.T) {
	var compatibilityReport engines.CompatibilityReport = engines.CompatibilityReport{
		CompatibleMemory:   false,
		RequiredMemory:     8 * 1024 * 1024 * 1024, // 8 GiB
		TotalRAM:           4 * 1024 * 1024 * 1024, // 4 GiB
		TotalSwap:          0,
		CompatibleDisk:     false,
		RequiredDiskSpace:  100 * 1024 * 1024 * 1024, // 100 GiB
		AvailableDiskSpace: 50 * 1024 * 1024 * 1024,  // 50 GiB
		CompatibleDevices:  false,
	}

	expectedReasons := []string{
		"insufficient memory",
		"insufficient disk space",
		"required device not found",
	}

	var engineDetails EngineDetails
	engineDetails.fillIncompatibilityIssues(compatibilityReport)
	actualReasons := engineDetails.CompatibilityIssues

	if len(actualReasons) != len(expectedReasons) {
		t.Fatalf("Expected to have %d compatibility issues, got: %d", len(expectedReasons), len(actualReasons))
	}

	for i, reason := range actualReasons {
		if reason != expectedReasons[i] {
			t.Errorf("Expected reason: %s, got: %s", expectedReasons[i], reason)
		}
	}
}
