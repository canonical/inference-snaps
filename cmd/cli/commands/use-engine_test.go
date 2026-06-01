package commands

import (
	"errors"
	"testing"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/hardware_info"
	"github.com/canonical/inference-snaps-cli/pkg/selector"
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

func ExampleUseEngine_autoSelectEngine() {
	cache := storage.NewMockCache()
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
	cmd.Cache.SetActiveEngine("")
	cmd.Verbose = true
	var allEngines []engines.Manifest
	for _, name := range []string{"not-compatible-engine", "cpu-exptl", "cpu"} {
		e, err := engines.LoadManifest(cmd.Context.EnginesDir, name)
		if err != nil {
			panic(err)
		}
		allEngines = append(allEngines, *e)
	}
	machineInfo, err := hardware_info.GetFromRawData("mustang", true, "../../../test_data")
	if err != nil {
		panic(err)
	}

	scoredEngines, err := selector.ScoreEngines(machineInfo, allEngines)
	if err != nil {
		panic(err)
	}
	if err := cmd.autoSelectScoredEngine(scoredEngines); err != nil {
		panic(err)
	}

	// Output:
	// Evaluating engines for optimal hardware compatibility:
	// ✘ not-compatible-engine: not compatible
	//   - required device not found
	// • cpu-exptl: experimental, score=12
	// ✔ cpu: compatible, score=12
	// Selected engine: cpu
	// Engine changed to "cpu".
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
