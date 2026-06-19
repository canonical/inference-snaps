package commands

import (
	"errors"
	"os"
	"testing"

	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/common"
	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/canonical/inference-snaps-cli/v2/pkg/hardware_info"
	"github.com/canonical/inference-snaps-cli/v2/pkg/selector"
	"github.com/canonical/inference-snaps-cli/v2/pkg/snap"
	"github.com/canonical/inference-snaps-cli/v2/pkg/storage"
)

func ExampleUseEngine_noRestartWhenEngineUnchanged() {
	// intel-gpu now requires runtime and model components, so we need SNAP_COMPONENTS
	// to be set with those component directories so InstallMissingComponents is a no-op.
	snapComponents, err := os.MkdirTemp("", "snap-components-*")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(snapComponents)
	for _, comp := range []string{"runtime-openvino-model-server", "model-4b-it-int4-fq-ov"} {
		if err := os.Mkdir(snapComponents+"/"+comp, 0755); err != nil {
			panic(err)
		}
	}
	if err := os.Setenv("SNAP_COMPONENTS", snapComponents); err != nil {
		panic(err)
	}
	defer os.Unsetenv("SNAP_COMPONENTS")

	cache := storage.NewMockCache()
	cache.SetActiveEngine("intel-gpu")
	config := storage.NewMockConfig()
	cmd := useEngineCommand{
		assumeYes: true,
		Context: &common.Context{
			EnginesDir:  "../../../test_data/engines",
			RuntimesDir: "../../../test_data/runtimes",
			ModelsDir:   "../../../test_data/models",
			Cache:       cache,
			Config:      config,
			Snap:        snap.Mock(),
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
			EnginesDir:  "../../../test_data/engines",
			RuntimesDir: "../../../test_data/runtimes",
			ModelsDir:   "../../../test_data/models",
			Cache:       cache,
			Config:      config,
			Snap:        snap.Mock(),
		},
	}
	// Create a temporary SNAP_COMPONENTS directory with stub component directories so that
	// required components appear "installed" and the install flow produces no extra output.
	snapComponents, err := os.MkdirTemp("", "snap-components-*")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(snapComponents)
	if err := os.Mkdir(snapComponents+"/runtime-llama-cpp-cpu", 0755); err != nil {
		panic(err)
	}
	if err := os.Mkdir(snapComponents+"/model-26b-a4b-q4-k-m-gguf", 0755); err != nil {
		panic(err)
	}
	if err := os.Mkdir(snapComponents+"/mmproj-26b-bf16-gguf", 0755); err != nil {
		panic(err)
	}
	if err := os.Setenv("SNAP_COMPONENTS", snapComponents); err != nil {
		panic(err)
	}
	defer os.Unsetenv("SNAP_COMPONENTS")
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
	// • cpu-exptl: experimental, score=10
	// ✔ cpu: compatible, score=10
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

func TestAutoSelectEngine_fallbackToEngine(t *testing.T) {

	cache := storage.NewMockCache()
	cmd := useEngineCommand{
		Context: &common.Context{
			EnginesDir: "../../../test_data/engines",
			Cache:      cache,
			Snap:       snap.Mock(),
		},
		fallback:  "amd-gpu",
		auto:      true,
		assumeYes: true,
	}
	err := cmd.autoSelectEngine()
	if err != nil {
		t.Fatalf("unexpected error auto selecting engine: %v", err)
	}

	activeEngine, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		t.Fatalf("unexpected error getting active engine: %v", err)
	}
	if activeEngine != "amd-gpu" {
		t.Errorf("expected active engine to be 'amd-gpu', got %q", activeEngine)
	}
}

func TestFixActiveEngine_fallbackWhenActiveEngineMissing(t *testing.T) {
	cache := storage.NewMockCache()
	cache.SetActiveEngine("missing-engine")

	cmd := useEngineCommand{
		Context: &common.Context{
			EnginesDir: "../../../test_data/engines",
			Cache:      cache,
			Config:     storage.NewMockConfig(),
			Snap:       snap.Mock(),
		},
		fallback:  "amd-gpu",
		fix:       true,
		assumeYes: true,
	}

	err := cmd.fixActiveEngine()
	if err != nil {
		t.Fatalf("unexpected error fixing active engine: %v", err)
	}

	activeEngine, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		t.Fatalf("unexpected error getting active engine: %v", err)
	}
	if activeEngine != "amd-gpu" {
		t.Errorf("expected active engine to be 'amd-gpu', got %q", activeEngine)
	}
}
