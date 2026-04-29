package commands

import (
	"os"
	"testing"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/storage"
)

func TestGetEnvVarsFromPassthroughConfigs(t *testing.T) {
	cmd := runCommand{}
	passthrough := map[string]any{
		"passthrough.environment.my-key": "value",
		"passthrough.environment.other":  123,
		"passthrough.not-environment":    "ignored",
	}

	got, err := cmd.getEnvVarsFromPassthroughConfigs(passthrough)
	if err != nil {
		t.Fatalf("getEnvVarsFromPassthroughConfigs returned error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("getEnvVarsFromPassthroughConfigs returned %d keys, want 2", len(got))
	}

	if got["MY_KEY"] != "value" {
		t.Fatalf("getEnvVarsFromPassthroughConfigs returned %v for MY_KEY, want value", got["MY_KEY"])
	}

	if got["OTHER"] != "123" {
		t.Fatalf("getEnvVarsFromPassthroughConfigs returned %v for OTHER, want 123", got["OTHER"])
	}
}

func TestProcessPassthroughConfigs(t *testing.T) {
	mockConfig := storage.NewMockConfig()
	mockConfig.Set("passthrough.environment.my-key", "value", storage.UserConfig)
	mockConfig.Set("passthrough.environment.other", "123", storage.UserConfig)
	mockConfig.Set("passthrough.not-environment", "ignored", storage.UserConfig)
	cmd := runCommand{
		Context: &common.Context{
			Config: mockConfig,
		},
	}

	err := cmd.processPassthroughConfigs()
	if err != nil {
		t.Fatalf("processPassthroughConfigs returned error: %v", err)
	}

	if got := os.Getenv("MY_KEY"); got != "value" {
		t.Fatalf("expected MY_KEY to be %q, got %q", "value", got)
	}

	if got := os.Getenv("OTHER"); got != "123" {
		t.Fatalf("expected OTHER to be %q, got %q", "123", got)
	}
}
