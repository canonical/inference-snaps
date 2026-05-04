package commands

import (
	"errors"
	"testing"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/snap"
	"github.com/canonical/inference-snaps-cli/pkg/storage"
)

func ExampleUseEngine_noRestartWhenEngineUnchanged() {
	cache := storage.NewMockCache()
	cache.SetActiveEngine("intel-gpu")
	config := storage.NewMockConfig()
	cmd := useEngineCommand{
		assumeYes: true,
		Context: &common.Context{
			EnginesDir: "../../../test_data/engines",
			Cache:      cache,
			Config:     config,
			Snap:       snap.Mock(),
		},
	}

	if err := cmd.switchEngine("intel-gpu"); err != nil {
		panic(err)
	}

	// Output:
}

func ExampleUseEngine_restartWhenEngineChanged() {
	cache := storage.NewMockCache()
	cache.SetActiveEngine("intel-gpu")
	config := storage.NewMockConfig()
	cmd := useEngineCommand{
		assumeYes: true,
		Context: &common.Context{
			EnginesDir: "../../../test_data/engines",
			Cache:      cache,
			Config:     config,
			Snap:       snap.Mock(),
		},
	}

	if err := cmd.switchEngine("cpu-avx1"); err != nil {
		panic(err)
	}

	// Output:
	// Engine changed to "cpu-avx1".
	// [mock] Restarting all services
}

func TestFixActiveEngine_noActiveEngine(t *testing.T) {
	cache := storage.NewMockCache()
	cmd := useEngineCommand{
		Context: &common.Context{
			EnginesDir: "../../../test_data/engines",
			Cache:      cache,
			Snap:       snap.Mock(),
		},
	}

	err := cmd.fixActiveEngine()
	if !errors.Is(err, common.ErrNoActiveEngine) {
		t.Errorf("expected no active engine error, got %v", err)
	}
}
