package commands

import (
	"os"
	"testing"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/snap"
	"github.com/canonical/inference-snaps-cli/pkg/storage"
)

func TestUnsetValueRemovesUserConfigWithoutRestart(t *testing.T) {
	config := storage.NewMockConfig()
	config.Set("api.endpoint", "https://example.com", storage.UserConfig)
	cmd := unsetCommand{
		noRestart: false,
		assumeYes: true,
		Context: &common.Context{
			Config: config,
			Snap:   snap.Mock(),
		},
	}

	if err := cmd.unsetValue("api.endpoint"); err != nil {
		t.Fatalf("unsetValue returned an unexpected error: %v", err)
	}

	values, err := config.Get("api.endpoint")
	if err != nil {
		t.Fatalf("Get returned an unexpected error: %v", err)
	}
	if len(values) != 0 {
		t.Fatalf("expected api.endpoint to be removed, got %#v", values)
	}
}

func TestUnsetKeyToDefaultValue(t *testing.T) {
	config := storage.NewMockConfig()
	config.Set("test-key", "user-value", storage.UserConfig)
	config.Set("test-key", "engine-value", storage.EngineConfig)
	config.Set("test-key", "package-value", storage.PackageConfig)
	cmd := unsetCommand{
		noRestart: true,
		assumeYes: true,
		Context: &common.Context{
			Config: config,
			Snap:   snap.Mock(),
		},
	}
	values, err := config.GetAll()
	if err != nil {
		t.Fatalf("GetAll returned an unexpected error: %v", err)
	}
	if err := cmd.unsetValue("test-key"); err != nil {
		t.Fatalf("unsetValue returned an unexpected error: %v", err)
	}

	values, err = config.Get("test-key")
	if err != nil {
		t.Fatalf("Get returned an unexpected error: %v", err)
	}
	if values["test-key"] != "engine-value" {
		t.Fatalf("expected test-key to be overridden by engine config, got %#v", values["test-key"])
	}
}

func TestUnsetNonexistentKey(t *testing.T) {
	config := storage.NewMockConfig()
	cmd := unsetCommand{
		noRestart: true,
		assumeYes: true,
		Context: &common.Context{
			Config: config,
			Snap:   snap.Mock(),
		},
	}
	if err := os.Setenv("SNAP_INSTANCE_NAME", "mock-snap"); err != nil {
		panic(err)
	}
	defer func() {
		_ = os.Unsetenv("SNAP_INSTANCE_NAME")
	}()
	err := cmd.unsetValue("nonexistent-key")
	if err == nil {
		t.Fatal("expected an error when unsetting a non-existent key, got nil")
	}
	expectedErrMsg := "key \"nonexistent-key\" is not found\n\nUse \"mock-snap get\" to view available keys"
	if err.Error() != expectedErrMsg {
		t.Fatalf("expected error message %q, got %q", expectedErrMsg, err.Error())
	}
}
