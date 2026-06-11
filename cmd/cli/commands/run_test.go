package commands

import (
	"os"
	"testing"

	"github.com/canonical/inference-snaps-cli/v2/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/v2/pkg/storage"
)

func TestProcessEnvConfigs(t *testing.T) {
	mockConfig := storage.NewMockConfig()
	mockConfig.Set("env.my-key", "value", storage.UserConfig)
	mockConfig.Set("env.other", "123", storage.UserConfig)
	mockConfig.Set("other.ignored", "ignored", storage.UserConfig)
	cmd := runCommand{
		Context: &common.Context{
			Config: mockConfig,
		},
	}

	err := cmd.processEnvConfigs()
	if err != nil {
		t.Fatalf("processEnvConfigs returned error: %v", err)
	}

	if got := os.Getenv("MY_KEY"); got != "value" {
		t.Fatalf("expected MY_KEY to be %q, got %q", "value", got)
	}

	if got := os.Getenv("OTHER"); got != "123" {
		t.Fatalf("expected OTHER to be %q, got %q", "123", got)
	}
}
