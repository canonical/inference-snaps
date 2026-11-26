package commands

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/canonical/inference-snaps/modelctl/v2/commands/common"
	"github.com/canonical/inference-snaps/modelctl/v2/internal/engines"
	"github.com/canonical/inference-snaps/modelctl/v2/internal/models"
	"github.com/canonical/inference-snaps/modelctl/v2/internal/selector"
	"github.com/canonical/inference-snaps/modelctl/v2/internal/snap"
	"github.com/canonical/inference-snaps/modelctl/v2/internal/storage"
)

// setupSnapComponents creates a temporary directory, populates it with a
// subdirectory for each named component, sets SNAP_COMPONENTS to that
// directory for the duration of the test, and registers a cleanup that
// removes the directory.
func setupSnapComponents(t *testing.T, components ...string) {
	t.Helper()
	dir, err := os.MkdirTemp("", "snap-components-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	for _, c := range components {
		if err := os.Mkdir(dir+"/"+c, 0755); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("SNAP_COMPONENTS", dir)
}

// newUseEngineCmd builds a useEngineCommand backed by mock storage and a fresh
// cache, wired to the shared test_data directories.
func newUseEngineCmd() *useEngineCommand {
	return &useEngineCommand{
		assumeYes: true,
		noRestart: true,
		Context: &common.Context{
			EnginesDir:  "../test_data/engines",
			RuntimesDir: "../test_data/runtimes",
			ModelsDir:   "../test_data/models",
			Cache:       storage.NewMockCache(),
			Config:      storage.NewMockConfig(),
			Snap:        snap.Mock(),
		},
	}
}

// loadScoredEngine loads an engine manifest from test_data and wraps it in a
// ScoredManifest with the given score (0 = incompatible).
func loadScoredEngine(t *testing.T, name string, score int) engines.ScoredManifest {
	t.Helper()
	m, err := engines.LoadManifest("../test_data/engines", name)
	if err != nil {
		t.Fatalf("loading engine %q: %v", name, err)
	}
	return engines.ScoredManifest{Manifest: *m, Score: score}
}

func ExampleUseEngine_noRestartWhenEngineAndModelUnchanged() {
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
	cache.SetActiveModel("4b-it-int4-fq-ov")
	config := storage.NewMockConfig()
	cmd := useEngineCommand{
		assumeYes: true,
		Context: &common.Context{
			EnginesDir:  "../test_data/engines",
			RuntimesDir: "../test_data/runtimes",
			ModelsDir:   "../test_data/models",
			Cache:       cache,
			Config:      config,
			Snap:        snap.Mock(),
		},
	}

	machine, err := machineFixture("dummy-machine")
	if err != nil {
		panic(err)
	}
	machine.PCIDevices[0].DeviceClass = 0x0300
	machine.PCIDevices[0].DeviceId = 0x56a0
	machine.PCIDevices[0].AdditionalProperties = map[string]any{"vram": uint64(16 * 1024 * 1024 * 1024)}
	if err := cmd.switchEngineWithMachine("intel-gpu", machine); err != nil {
		panic(err)
	}

	// Output:
}

func ExampleUseEngine_printMemoryIncompatibleModels() {
	cmd := useEngineCommand{Context: &common.Context{}}
	cmd.printScoredModels([]models.ScoredManifest{
		{Manifest: models.Manifest{Name: "fits"}, CompatibilityReport: models.CompatibilityReport{CompatibleDisk: true, RequiredMemory: 4, TotalMemory: 8}},
		{Manifest: models.Manifest{Name: "exceeds-ram"}, CompatibilityReport: models.CompatibilityReport{CompatibleDisk: true, RequiredMemory: 16, TotalMemory: 8}},
		{Manifest: models.Manifest{Name: "exceeds-vram"}, CompatibilityReport: models.CompatibilityReport{CompatibleDisk: true, RequiredVram: 16, TotalVram: 8}},
	}, "fits")

	// Output:
	// Selecting a compatible model:
	// ✔ fits: best match
	// ✘ exceeds-ram
	// ✘ exceeds-vram
}

func ExampleUseEngine_restartWhenEngineChanged() {
	snapComponents, err := os.MkdirTemp("", "snap-components-*")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(snapComponents)
	for _, comp := range []string{"runtime-llama-cpp-cpu", "model-26b-a4b-q4-k-m-gguf", "mmproj-26b-bf16-gguf"} {
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
			EnginesDir:  "../test_data/engines",
			RuntimesDir: "../test_data/runtimes",
			ModelsDir:   "../test_data/models",
			Cache:       cache,
			Config:      config,
			Snap:        snap.Mock(),
		},
	}

	machine, err := machineFixture("dummy-machine")
	if err != nil {
		panic(err)
	}
	if err := cmd.switchEngineWithMachine("cpu-avx1", machine); err != nil {
		panic(err)
	}

	// Output:
	// Engine changed to "cpu-avx1".
	// Model changed to "26b-q4-k-m-gguf".
	// [mock] Restarting all services
}

func ExampleUseEngine_autoSelectEngine() {
	cache := storage.NewMockCache()
	config := storage.NewMockConfig()
	cmd := useEngineCommand{
		assumeYes: true,
		auto:      true,
		Context: &common.Context{
			EnginesDir:  "../test_data/engines",
			RuntimesDir: "../test_data/runtimes",
			ModelsDir:   "../test_data/models",
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
	if err := os.Mkdir(snapComponents+"/model-30m-q4-k-m-gguf", 0755); err != nil {
		panic(err)
	}
	if err := os.Setenv("SNAP_COMPONENTS", snapComponents); err != nil {
		panic(err)
	}
	defer os.Unsetenv("SNAP_COMPONENTS")
	cmd.Cache.SetActiveEngine("")
	var allEngines []engines.Manifest
	for _, name := range []string{"not-compatible-engine", "cpu-exptl", "cpu"} {
		e, err := engines.LoadManifest(cmd.EnginesDir, name)
		if err != nil {
			panic(err)
		}
		allEngines = append(allEngines, *e)
	}
	machine, err := machineFixture("low-disk-available-machine")
	if err != nil {
		panic(err)
	}

	scoredEngines, err := selector.ScoreEngines(machine, allEngines)
	if err != nil {
		panic(err)
	}
	if err := cmd.autoSelectScoredEngine(scoredEngines, machine); err != nil {
		panic(err)
	}

	// Output:
	// Evaluating engines for optimal hardware compatibility:
	// ✘ not-compatible-engine
	// ✔ cpu-exptl: experimental
	// ✔ cpu: best match
	//
	// Selecting a compatible model:
	// ✔ 30m-q4-k-m-gguf: best match
	// ✘ 26b-q4-k-m-gguf
	// ✘ 30b-a3b-q4-k-m-gguf
	//
	// Engine changed to "cpu".
	// Model changed to "30m-q4-k-m-gguf".
	// [mock] Restarting all services
}

func ExampleUseEngine_autoSelectEngineVerbose() {
	cache := storage.NewMockCache()
	config := storage.NewMockConfig()
	cmd := useEngineCommand{
		assumeYes: true,
		auto:      true,
		Context: &common.Context{
			EnginesDir:  "../test_data/engines",
			RuntimesDir: "../test_data/runtimes",
			ModelsDir:   "../test_data/models",
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
		e, err := engines.LoadManifest(cmd.EnginesDir, name)
		if err != nil {
			panic(err)
		}
		allEngines = append(allEngines, *e)
	}
	machine, err := machineFixture("no-disk-available-machine")
	if err != nil {
		panic(err)
	}

	scoredEngines, err := selector.ScoreEngines(machine, allEngines)
	if err != nil {
		panic(err)
	}
	if err := cmd.autoSelectScoredEngine(scoredEngines, machine); err != nil {
		panic(err)
	}

	// Output:
	// Evaluating engines for optimal hardware compatibility:
	// ✘ not-compatible-engine
	//   - required device not found
	// ✔ cpu-exptl: score=10, experimental
	// ✔ cpu: score=10, best match
	//
	// Selecting a compatible model:
	// ✘ 30m-q4-k-m-gguf:
	//   - requires 1024.0KiB memory, has 1024
	//   - requires 1024.0KiB disk space, has 0
	// ✘ 26b-q4-k-m-gguf:
	//   - requires 8.0GiB memory, has 1024
	// ✘ 30b-a3b-q4-k-m-gguf:
	//   - requires 5.9TiB memory, has 1024
	//   - requires 5.9TiB disk space, has 0
	//
	// Engine changed to "cpu".
	// [mock] Restarting all services
}

func ExampleUseEngine_printIncompatibleModels() {
	cache := storage.NewMockCache()
	config := storage.NewMockConfig()
	cmd := useEngineCommand{
		assumeYes: true,
		auto:      true,
		Context: &common.Context{
			EnginesDir:  "../test_data/engines",
			RuntimesDir: "../test_data/runtimes",
			ModelsDir:   "../test_data/models",
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
	cmd.Verbose = true
	machine, err := machineFixture("no-disk-available-machine")
	if err != nil {
		panic(err)
	}

	if err := cmd.switchEngineWithMachine("cpu", machine); err != nil {
		panic(err)
	}

	// Output:
	// Selecting a compatible model:
	// ✘ 30m-q4-k-m-gguf:
	//   - requires 1024.0KiB memory, has 1024
	//   - requires 1024.0KiB disk space, has 0
	// ✘ 26b-q4-k-m-gguf:
	//   - requires 8.0GiB memory, has 1024
	// ✘ 30b-a3b-q4-k-m-gguf:
	//   - requires 5.9TiB memory, has 1024
	//   - requires 5.9TiB disk space, has 0
	//
	// Engine changed to "cpu".
	// [mock] Restarting all services
}

// TestSwitchEngine_withModelID verifies that switchEngineAndModel respects an
// explicitly supplied modelID:
//   - a valid model supported by the engine is persisted as the active model
//   - an unsupported model returns an error and leaves state unchanged
func TestSwitchEngine_withModelID(t *testing.T) {
	tests := []struct {
		name                    string
		engineName              string
		modelID                 string
		installedComponents     []string
		wantErr                 bool
		wantEngine              string
		wantModel               string
		autoSelect              bool
		setPreviousActiveEngine string
		setPreviousActiveModel  string
	}{
		{
			// Switch to the cpu engine and explicitly request the non-default
			// model (30b-a3b-q4-k-m-gguf).  The engine's default is
			// 26b-q4-k-m-gguf, so this exercises the override branch.
			name:       "valid non-default model is used",
			engineName: "cpu",
			modelID:    "30b-a3b-q4-k-m-gguf",
			installedComponents: []string{
				"runtime-llama-cpp-cpu",
				"model-30b-a3b-q4-k-m-gguf-1-of-6",
				"model-30b-a3b-q4-k-m-gguf-2-of-6",
				"model-30b-a3b-q4-k-m-gguf-3-of-6",
				"model-30b-a3b-q4-k-m-gguf-4-of-6",
				"model-30b-a3b-q4-k-m-gguf-5-of-6",
				"model-30b-a3b-q4-k-m-gguf-6-of-6",
			},
			wantErr:    false,
			wantEngine: "cpu",
			wantModel:  "30b-a3b-q4-k-m-gguf",
		},
		{
			// Switch to the cpu engine and explicitly request the default model
			// (26b-q4-k-m-gguf) – this should also succeed and store the
			// explicit choice.
			name:       "valid default model is used",
			engineName: "cpu",
			modelID:    "26b-q4-k-m-gguf",
			installedComponents: []string{
				"runtime-llama-cpp-cpu",
				"model-26b-a4b-q4-k-m-gguf",
				"mmproj-26b-bf16-gguf",
			},
			wantErr:    false,
			wantEngine: "cpu",
			wantModel:  "26b-q4-k-m-gguf",
		},
		{
			// Request a model that the cpu engine does not support.
			// switchEngine must return an error and must not change any state.
			name:       "unsupported model returns error",
			engineName: "cpu",
			modelID:    "4b-it-int4-fq-ov", // only supported by intel-gpu
			installedComponents: []string{
				"runtime-llama-cpp-cpu",
			},
			wantErr: true,
		},
		{
			// Request an engine with no models that fits the system's available resources.
			// switchEngineAndModel should return an error if autoselection is disable
			name:       "no suitable models available",
			engineName: "cpu",
			modelID:    "",
			installedComponents: []string{
				"runtime-llama-cpp-cpu",
			},
			wantErr:    true,
			wantEngine: "cpu",
			autoSelect: false,
			wantModel:  "",
		},
		{
			// Request an engine with no models that fits the system's available resources.
			// switchEngineAndModel should just print an error if autoselection is enabled
			name:       "no suitable models available autoselection",
			engineName: "cpu",
			modelID:    "",
			installedComponents: []string{
				"runtime-llama-cpp-cpu",
			},
			wantErr:    false,
			autoSelect: true,
			wantEngine: "cpu",
			wantModel:  "",
		},
		{
			// Switch to the same engine with an empty model ID when no model is currently active.
			// The engine has model options available but none were selected.
			// switchEngineAndModel should return an error if autoselection is disabled
			name:       "no change and unselected model with autoselect disabled",
			engineName: "cpu",
			modelID:    "",
			installedComponents: []string{
				"runtime-llama-cpp-cpu",
				"model-26b-a4b-q4-k-m-gguf",
				"mmproj-26b-bf16-gguf",
			},
			wantErr:                 true,
			autoSelect:              false,
			wantEngine:              "cpu",
			wantModel:               "",
			setPreviousActiveEngine: "cpu",
		},
		{
			// Switch to the same engine with an empty model ID when no model is currently active.
			// The engine has model options available but none were selected.
			// switchEngineAndModel should just print an error if autoselection is enabled
			name:       "no change and unselected model with autoselect enabled",
			engineName: "cpu",
			modelID:    "",
			installedComponents: []string{
				"runtime-llama-cpp-cpu",
				"model-26b-a4b-q4-k-m-gguf",
				"mmproj-26b-bf16-gguf",
			},
			wantErr:                 false,
			autoSelect:              true,
			wantEngine:              "cpu",
			wantModel:               "",
			setPreviousActiveEngine: "cpu",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupSnapComponents(t, tt.installedComponents...)

			cmd := newUseEngineCmd()
			if tt.autoSelect {
				cmd.auto = true
			} else {
				cmd.auto = false
			}
			err := cmd.Cache.SetActiveModel(tt.setPreviousActiveModel)
			if err != nil {
				t.Fatalf("unexpected error setting active model: %v", err)
			}

			err = cmd.Cache.SetActiveEngine(tt.setPreviousActiveEngine)
			if err != nil {
				t.Fatalf("unexpected error setting active engine: %v", err)
			}

			err = cmd.switchEngineAndModel(tt.engineName, tt.modelID)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			activeEngine, err := cmd.Cache.GetActiveEngine()
			if err != nil {
				t.Fatalf("getting active engine: %v", err)
			}
			if activeEngine != tt.wantEngine {
				t.Errorf("active engine = %q, want %q", activeEngine, tt.wantEngine)
			}

			activeModel, err := cmd.Cache.GetActiveModel()
			if err != nil {
				t.Fatalf("getting active model: %v", err)
			}
			if activeModel != tt.wantModel {
				t.Errorf("active model = %q, want %q", activeModel, tt.wantModel)
			}
		})
	}
}

func TestSwitchEngineWithMachine_lowDisk(t *testing.T) {
	setupSnapComponents(t)
	cmd := newUseEngineCmd()
	machine, err := machineFixture("low-disk-available-machine")
	if err != nil {
		t.Fatalf("unexpected error getting machine info fixture: %v", err)
	}
	err = cmd.switchEngineWithMachine("cpu", machine)
	if err != nil {
		t.Fatalf("unexpected error switching engine: %v", err)
	}
	activeModel, err := cmd.Cache.GetActiveModel()
	if err != nil {
		t.Fatalf("unexpected error getting active model: %v", err)
	}
	if activeModel != "30m-q4-k-m-gguf" {
		t.Errorf("active model = %q, want %q", activeModel, "30m-q4-k-m-gguf")
	}
}

func TestSwitchEngineWithMachineInfo_noDiskNoModel(t *testing.T) {
	setupSnapComponents(t)
	cmd := newUseEngineCmd()
	machine, err := machineFixture("no-disk-available-machine")
	if err != nil {
		t.Fatalf("unexpected error getting machine info fixture: %v", err)
	}
	err = cmd.switchEngineWithMachine("cpu", machine)
	if err == nil {
		t.Fatal("expected error switching engine when no compatible model fits, got nil")
	}
	// The engine change is still applied, but no model could be selected.
	activeEngine, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		t.Fatalf("unexpected error getting active engine: %v", err)
	}
	if activeEngine != "cpu" {
		t.Errorf("active engine = %q, want %q", activeEngine, "cpu")
	}
	activeModel, err := cmd.Cache.GetActiveModel()
	if err != nil {
		t.Fatalf("unexpected error getting active model: %v", err)
	}
	if activeModel != "" {
		t.Errorf("active model = %q, want %q", activeModel, "")
	}
}

func TestFixActiveEngine_noActiveEngine(t *testing.T) {
	cache := storage.NewMockCache()
	cmd := useEngineCommand{
		Context: &common.Context{
			EnginesDir: "../test_data/engines",
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
			EnginesDir: "../test_data/engines",
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
			EnginesDir: "../test_data/engines",
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

func TestFixActiveEngineWithScores(t *testing.T) {
	tests := []struct {
		name         string
		active       string
		fallback     string
		observable   bool
		architecture string
		engines      []string
		wantEngine   string
		wantErr      string
	}{
		{name: "compatible engine retained", active: "intel-gpu", fallback: "cpu", observable: true, wantEngine: "intel-gpu"},
		{name: "incompatible engine uses fallback", active: "amd-gpu", fallback: "cpu", observable: true, wantEngine: "cpu"},
		{name: "removed engine uses fallback", active: "missing-engine", fallback: "cpu", observable: true, wantEngine: "cpu"},
		{name: "incompatible fallback", active: "amd-gpu", fallback: "cpu", observable: true, architecture: "riscv64", wantEngine: "amd-gpu", wantErr: `fallback engine "cpu" is unavailable or incompatible`},
		{name: "missing fallback", active: "amd-gpu", fallback: "cpu", observable: true, engines: []string{"amd-gpu"}, wantEngine: "amd-gpu", wantErr: `fallback engine "cpu" is unavailable or incompatible`},
		{name: "no fallback", active: "missing-engine", observable: true, fallback: "", wantEngine: "missing-engine", wantErr: "no fallback engine was specified"},
		{name: "unobservable hardware retains engine", active: "intel-gpu", fallback: "cpu", wantEngine: "intel-gpu"},
		{name: "unobservable hardware without fallback", active: "missing-engine", fallback: "", wantEngine: "missing-engine", wantErr: "no fallback engine was specified"},
	}
	manifests, err := engines.LoadManifests("../test_data/engines")
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupSnapComponents(t, "runtime-openvino-model-server", "model-4b-it-int4-fq-ov", "runtime-llama-cpp-cpu", "model-26b-a4b-q4-k-m-gguf", "mmproj-26b-bf16-gguf")
			cmd := newUseEngineCmd()
			cmd.fallback = tt.fallback
			if err := cmd.Cache.SetActiveEngine(tt.active); err != nil {
				t.Fatal(err)
			}
			if tt.active == "intel-gpu" {
				if err := cmd.Cache.SetActiveModel("4b-it-int4-fq-ov"); err != nil {
					t.Fatal(err)
				}
			}
			machineInfo, err := machineFixture("i7-1165G7")
			if err != nil {
				t.Fatal(err)
			}
			if tt.architecture != "" {
				machineInfo.CPUs[0].Architecture = tt.architecture
			}
			enginesToScore := manifests
			if tt.engines != nil {
				enginesToScore = nil
				for _, name := range tt.engines {
					manifest, err := engines.LoadManifest(cmd.EnginesDir, name)
					if err != nil {
						t.Fatal(err)
					}
					enginesToScore = append(enginesToScore, *manifest)
				}
			}
			scored, err := selector.ScoreEngines(machineInfo, enginesToScore)
			if err != nil {
				t.Fatal(err)
			}
			machineForFix := machineInfo
			if !tt.observable {
				machineForFix = nil
			}
			err = cmd.fixActiveEngineWithScores(tt.active, tt.observable, scored, machineForFix)
			if tt.wantErr == "" && err != nil || tt.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tt.wantErr)) {
				t.Fatalf("got error %v, want %q", err, tt.wantErr)
			}
			active, err := cmd.Cache.GetActiveEngine()
			if err != nil {
				t.Fatal(err)
			}
			if active != tt.wantEngine {
				t.Errorf("active engine = %q, want %q", active, tt.wantEngine)
			}
			if tt.wantErr == "" {
				model, err := cmd.Cache.GetActiveModel()
				if err != nil {
					t.Fatal(err)
				}
				wantModel := "4b-it-int4-fq-ov"
				configKey, wantConfig := "target-device", any("GPU")
				if tt.wantEngine == "cpu" {
					wantModel = "26b-q4-k-m-gguf"
					configKey, wantConfig = "sleep-idle-seconds", 600
				}
				if model != wantModel {
					t.Errorf("active model = %q, want %q", model, wantModel)
				}
				config, err := cmd.Config.Get(configKey)
				if err != nil {
					t.Fatal(err)
				}
				if config[configKey] != wantConfig {
					t.Errorf("%s = %v, want %v", configKey, config[configKey], wantConfig)
				}
			}
		})
	}
}

func TestFixActiveEngineObservabilityError(t *testing.T) {
	cmd := newUseEngineCmd()
	cmd.Snap = snap.MockWithHardwareObservability(false, errors.New("observe failed"))
	if err := cmd.Cache.SetActiveEngine("cpu"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.fixActiveEngine(); err == nil || !strings.Contains(err.Error(), "observe failed") {
		t.Fatalf("expected observability error, got %v", err)
	}
}

func TestFixActiveEngineUnobservableKeepsEngine(t *testing.T) {
	setupSnapComponents(t, "runtime-openvino-model-server", "model-4b-it-int4-fq-ov")
	cmd := newUseEngineCmd()
	cmd.fallback = "cpu"
	if err := cmd.Cache.SetActiveEngine("intel-gpu"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Cache.SetActiveModel("4b-it-int4-fq-ov"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.fixActiveEngine(); err != nil {
		t.Fatal(err)
	}
	active, err := cmd.Cache.GetActiveEngine()
	if err != nil || active != "intel-gpu" {
		t.Fatalf("active engine = %q, error = %v", active, err)
	}
	config, err := cmd.Config.Get("target-device")
	if err != nil || config["target-device"] != "GPU" {
		t.Fatalf("engine config = %v, error = %v", config, err)
	}
}

func TestFixActiveEngineChecksComponents(t *testing.T) {
	cmd := newUseEngineCmd()
	cmd.RuntimesDir = t.TempDir()
	if err := cmd.Cache.SetActiveEngine("intel-gpu"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Cache.SetActiveModel("4b-it-int4-fq-ov"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.fixActiveEngine(); err == nil || !strings.Contains(err.Error(), "installing missing components") {
		t.Fatalf("expected component check error, got %v", err)
	}
}

func TestFixActiveEngineScoringError(t *testing.T) {
	cmd := newUseEngineCmd()
	cmd.EnginesDir = t.TempDir() + "/missing"
	cmd.Snap = snap.MockWithHardwareObservability(true, nil)
	if err := cmd.Cache.SetActiveEngine("cpu"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.fixActiveEngine(); err == nil || !strings.Contains(err.Error(), "scoring engines") {
		t.Fatalf("expected scoring error, got %v", err)
	}
}

// TestSwitchToPreinstalledEngineAndModel covers the component-aware engine
// selection logic introduced in switchToPreinstalledEngineAndModel.
func TestSwitchToPreinstalledEngineAndModel(t *testing.T) {
	// All required components for cpu + 26b-q4-k-m-gguf (runtime + model shards).
	// Pre-installing them makes InstallMissingComponents a no-op so the test
	// output stays clean.
	cpuAnd26bComponents := []string{
		"runtime-llama-cpp-cpu",
		"model-26b-a4b-q4-k-m-gguf",
		"mmproj-26b-bf16-gguf",
	}
	// All required components for cpu + 30b-a3b-q4-k-m-gguf.
	cpuAnd30bComponents := []string{
		"runtime-llama-cpp-cpu",
		"model-30b-a3b-q4-k-m-gguf-1-of-6",
		"model-30b-a3b-q4-k-m-gguf-2-of-6",
		"model-30b-a3b-q4-k-m-gguf-3-of-6",
		"model-30b-a3b-q4-k-m-gguf-4-of-6",
		"model-30b-a3b-q4-k-m-gguf-5-of-6",
		"model-30b-a3b-q4-k-m-gguf-6-of-6",
	}
	// All required components for intel-gpu + 4b-it-int4-fq-ov.
	intelGpuAnd4bComponents := []string{
		"runtime-openvino-model-server",
		"model-4b-it-int4-fq-ov",
	}

	tests := []struct {
		name                string
		installedComponents []string
		// engineScores lists the engines exposed to the function as already
		// scored.  Use score=0 to mark an engine as hardware-incompatible.
		engineScores []struct {
			name  string
			score int
		}
		wantSwitched bool
		wantEngine   string
		wantModel    string
	}{
		{
			// A single GGUF model component is enough to seed the model, which
			// in turn points to the cpu engine – even though intel-gpu has a
			// higher hardware-compatibility score.
			name:                "gguf model component selects cpu over higher-scored intel-gpu",
			installedComponents: cpuAnd26bComponents,
			engineScores: []struct {
				name  string
				score int
			}{
				{"cpu", 10},
				{"intel-gpu", 20}, // higher score, but its model is not seeded
			},
			wantSwitched: true,
			wantEngine:   "cpu",
			wantModel:    "26b-q4-k-m-gguf",
		},
		{
			// The mmproj component alone is sufficient to identify the 26b
			// model and choose the cpu engine.  The missing model shard
			// (model-26b-a4b-q4-k-m-gguf) would be installed by the mock snap
			// when switchEngine runs InstallMissingComponents.
			name: "mmproj component alone seeds cpu engine via 26b model",
			installedComponents: []string{
				"mmproj-26b-bf16-gguf",  // seed: one component of 26b-q4-k-m-gguf model
				"runtime-llama-cpp-cpu", // required by cpu runtime; avoids runtime-install noise
			},
			engineScores: []struct {
				name  string
				score int
			}{
				{"cpu", 10},
				{"intel-gpu", 20},
			},
			wantSwitched: true,
			wantEngine:   "cpu",
			wantModel:    "26b-q4-k-m-gguf",
		},
		{
			// An OpenVINO model component seeds the intel-gpu engine.
			name:                "openvino model component selects intel-gpu",
			installedComponents: intelGpuAnd4bComponents,
			engineScores: []struct {
				name  string
				score int
			}{
				{"cpu", 10},
				{"intel-gpu", 20},
			},
			wantSwitched: true,
			wantEngine:   "intel-gpu",
			wantModel:    "4b-it-int4-fq-ov",
		},
		{
			// The cpu runtime component seeds the cpu engine.  Because no model
			// component is installed the engine's default model is used.
			name:                "cpu runtime component seeds cpu engine with default model",
			installedComponents: cpuAnd26bComponents,
			engineScores: []struct {
				name  string
				score int
			}{
				{"cpu", 10},
				{"intel-gpu", 20},
			},
			wantSwitched: true,
			wantEngine:   "cpu",
			wantModel:    "26b-q4-k-m-gguf",
		},
		{
			// A single shard of a multi-shard model is enough to seed the whole
			// model (any one component match is sufficient).
			name:                "single shard of multi-shard model seeds cpu engine with 30b model",
			installedComponents: cpuAnd30bComponents, // shard 3 is the seed; all present for clean switchEngine
			engineScores: []struct {
				name  string
				score int
			}{
				{"cpu", 10},
				{"intel-gpu", 0}, // incompatible on this (hypothetical) machine
			},
			wantSwitched: true,
			wantEngine:   "cpu",
			wantModel:    "30b-a3b-q4-k-m-gguf",
		},
		{
			// Installing a CUDA runtime component when the CUDA engine is
			// hardware-incompatible (score=0) must not influence selection.
			// The function should return false so the caller falls back to
			// standard auto-selection.
			name:                "incompatible cuda runtime component falls back to standard selection",
			installedComponents: []string{"runtime-llama-cpp-cuda"},
			engineScores: []struct {
				name  string
				score int
			}{
				{"cuda-generic", 0}, // seeded by installed component, but incompatible
				{"cpu", 10},
				{"intel-gpu", 20},
			},
			wantSwitched: false,
			wantEngine:   "", // no engine set by this function
			wantModel:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupSnapComponents(t, tt.installedComponents...)

			var scored []engines.ScoredManifest
			for _, es := range tt.engineScores {
				scored = append(scored, loadScoredEngine(t, es.name, es.score))
			}

			cmd := newUseEngineCmd()
			machine, err := machineFixture("dummy-machine")
			if err != nil {
				t.Fatalf("loading machine fixture: %v", err)
			}
			switched, err := selectEngineForSeededComponents(cmd, scored, machine)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if switched != tt.wantSwitched {
				t.Errorf("switched = %v, want %v", switched, tt.wantSwitched)
			}

			if tt.wantEngine != "" {
				activeEngine, err := cmd.Cache.GetActiveEngine()
				if err != nil {
					t.Fatalf("getting active engine: %v", err)
				}
				if activeEngine != tt.wantEngine {
					t.Errorf("active engine = %q, want %q", activeEngine, tt.wantEngine)
				}
			}

			if tt.wantModel != "" {
				activeModel, err := cmd.Cache.GetActiveModel()
				if err != nil {
					t.Fatalf("getting active model: %v", err)
				}
				if activeModel != tt.wantModel {
					t.Errorf("active model = %q, want %q", activeModel, tt.wantModel)
				}
			}
		})
	}
}
