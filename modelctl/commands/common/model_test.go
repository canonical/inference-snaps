package common

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/canonical/inference-snaps-cli/v2/internal/engines"
	"github.com/canonical/inference-snaps-cli/v2/internal/storage"
	"github.com/canonical/inference-snaps-cli/v2/internal/utils"
	"github.com/canonical/lscompute/pkg/machine"
	"github.com/canonical/lscompute/pkg/machine/device/pci"
	"github.com/canonical/lscompute/pkg/machine/disk"
)

// writeModelYAML creates a model manifest at modelsDir/<name>/model.yaml with the given content.
func writeModelYAML(t *testing.T, modelsDir, name, content string) {
	t.Helper()
	dir := filepath.Join(modelsDir, name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("MkdirAll %s: %v", dir, err)
	}
	if err := os.WriteFile(filepath.Join(dir, "model.yaml"), []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile model.yaml: %v", err)
	}
}

// makeModelStatusCtx builds a Context with the given active model and a temp modelsDir.
func makeModelStatusCtx(t *testing.T, modelsDir, modelName string) *Context {
	t.Helper()
	cache := storage.NewMockCache()
	if err := cache.SetActiveModel(modelName); err != nil {
		t.Fatalf("SetActiveModel: %v", err)
	}
	return &Context{
		ModelsDir: modelsDir,
		Cache:     cache,
	}
}

func TestModelStatus_ModelNameWithEqualsInValue(t *testing.T) {
	modelsDir := t.TempDir()
	writeModelYAML(t, modelsDir, "my-model", `name: my=model`)
	ctx := makeModelStatusCtx(t, modelsDir, "my-model")

	status, err := ModelStatus(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status["name"] != "my=model" {
		t.Errorf("expected name %q, got %q", "my=model", status["name"])
	}
}

func TestModelStatus_NonExistentModel(t *testing.T) {
	modelsDir := t.TempDir()
	ctx := makeModelStatusCtx(t, modelsDir, "non-existent")

	_, err := ModelStatus(ctx)
	if err == nil {
		t.Fatal("expected error for non-existent model, got nil")
	}
}

func TestGetModelManifestByNameOrAlias(t *testing.T) {
	tests := []struct {
		name         string
		activeEngine string
		modelYAML    string // empty means don't write a model manifest
		engineYAML   string // empty means don't write an engine manifest
		query        string
		wantName     string // non-empty: expect this Name in the returned manifest
		wantAlias    string // non-empty: expect this Alias in the returned manifest
		wantErr      bool   // true: expect any non-nil error
	}{
		{
			name:         "found by alias",
			activeEngine: "my-engine",
			modelYAML:    "name: my-model-id\nalias: my-model\ndisk-size: 1G\n",
			engineYAML:   "name: my-engine\nmodel:\n  options:\n    - my-model-id\n",
			query:        "my-model",
			wantName:     "my-model-id",
		},
		{
			name:         "found by name",
			activeEngine: "my-engine",
			modelYAML:    "name: my-model-id\nalias: my-model\ndisk-size: 1G\n",
			engineYAML:   "name: my-engine\nmodel:\n  options:\n    - my-model-id\n",
			query:        "my-model-id",
			wantAlias:    "my-model",
		},
		{
			name:         "no active engine",
			activeEngine: "",
			query:        "my-model",
			wantErr:      true,
		},
		{
			name:         "incompatible with active engine",
			activeEngine: "my-engine",
			modelYAML:    "id: other-model-id\nname: other-model\ndisk-size: 1G\n",
			engineYAML:   "name: my-engine\nmodel:\n  options:\n    - some-other-model-id\n",
			query:        "other-model",
			wantErr:      true,
		},
		{
			name:         "model does not exist",
			activeEngine: "my-engine",
			engineYAML:   "name: my-engine\nmodel:\n  options: []\n",
			query:        "nonexistent-model",
			wantErr:      true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			modelsDir := t.TempDir()
			enginesDir := t.TempDir()

			if tc.modelYAML != "" {
				writeModelYAML(t, modelsDir, "my-model", tc.modelYAML)
			}
			if tc.engineYAML != "" {
				writeEngineYAML(t, enginesDir, "my-engine", tc.engineYAML)
			}

			cache := storage.NewMockCache()
			if err := cache.SetActiveEngine(tc.activeEngine); err != nil {
				t.Fatalf("SetActiveEngine: %v", err)
			}
			ctx := &Context{ModelsDir: modelsDir, EnginesDir: enginesDir, Cache: cache}

			manifest, err := GetModelManifestByNameOrAlias(ctx, tc.query)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.wantName != "" && manifest.Name != tc.wantName {
				t.Errorf("expected Name %q, got %q", tc.wantName, manifest.Name)
			}
			if tc.wantAlias != "" && manifest.Alias != tc.wantAlias {
				t.Errorf("expected Alias %q, got %q", tc.wantAlias, manifest.Alias)
			}
		})
	}
}

func TestGetAllModels(t *testing.T) {
	modelsDir := t.TempDir()
	enginesDir := t.TempDir()
	writeModelYAML(t, modelsDir, "model1", "name: model1\nalias: m1\ndisk-size: 1G\n")
	writeModelYAML(t, modelsDir, "model2", "name: model2\nalias: m2\ndisk-size: 2G\n")

	ctx := &Context{ModelsDir: modelsDir, EnginesDir: enginesDir}

	models, err := GetAllModels(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(models) != 2 {
		t.Fatalf("expected 2 models, got %d", len(models))
	}

	expectedNames := map[string]bool{"model1": true, "model2": true}
	for _, model := range models {
		if !expectedNames[model.Name] {
			t.Errorf("unexpected model name: %s", model.Name)
		}
	}
}

func TestGetAllModelsWithEngines(t *testing.T) {
	modelsDir := t.TempDir()
	enginesDir := t.TempDir()

	writeModelYAML(t, modelsDir, "model1", "name: model1\nalias: m1\ndisk-size: 1G\n")
	writeModelYAML(t, modelsDir, "model2", "name: model2\nalias: m2\ndisk-size: 2G\n")
	writeModelYAML(t, modelsDir, "model3", "name: model3\nalias: m3\ndisk-size: 3G\n")
	writeModelYAML(t, modelsDir, "model4", "name: model4\nalias: m4\ndisk-size: 4G\n")

	writeEngineYAML(t, enginesDir, "engine1", "name: engine1\nmodel:\n  options:\n    - model1\n    - model3\n")
	writeEngineYAML(t, enginesDir, "engine2", "name: engine2\nmodel:\n  options:\n    - model1\n    - model2\n")

	ctx := &Context{ModelsDir: modelsDir, EnginesDir: enginesDir}

	modelsWithEngines, err := GetAllModels(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(modelsWithEngines) != 4 {
		t.Fatalf("expected 4 models, got %d", len(modelsWithEngines))
	}

	for _, modelWithEngines := range modelsWithEngines {
		switch modelWithEngines.Name {
		case "model1":
			if len(modelWithEngines.CompatibleEngines) != 2 || modelWithEngines.CompatibleEngines[0] != "engine1" || modelWithEngines.CompatibleEngines[1] != "engine2" {
				t.Errorf("model1 should be compatible with engine1 and engine2, got: %v", modelWithEngines.CompatibleEngines)
			}
		case "model2":
			if len(modelWithEngines.CompatibleEngines) != 1 || modelWithEngines.CompatibleEngines[0] != "engine2" {
				t.Errorf("model2 should be compatible with engine2, got: %v", modelWithEngines.CompatibleEngines)
			}
		case "model3":
			if len(modelWithEngines.CompatibleEngines) != 1 || modelWithEngines.CompatibleEngines[0] != "engine1" {
				t.Errorf("model3 should be compatible with engine1, got: %v", modelWithEngines.CompatibleEngines)
			}
		case "model4":
			if len(modelWithEngines.CompatibleEngines) != 0 {
				t.Errorf("model4 should not be compatible with any engines, got: %v", modelWithEngines.CompatibleEngines)
			}
		default:
			t.Errorf("unexpected model name: %s", modelWithEngines.Name)
		}
	}
}

func TestAvailableDiskSpace(t *testing.T) {
	machine := new(machine.Machine{})
	machine.Disk = []disk.Disk{
		{
			Available:  1024,
			Total:      2048,
			Path:       "/var/lib/snapd/snaps",
			MountPoint: new("/"),
		},
	}

	diskSpace, err := availableDiskSpace(machine)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if diskSpace != 1024 {
		t.Errorf("expected positive available disk space, got: %d", diskSpace)
	}
}

func TestSelectModelMemory(t *testing.T) {
	for _, test := range []struct {
		name          string
		options       []string
		preferred     string
		totalMemory   uint64
		availableDisk uint64
		gpu           bool
		vram          any
		wantModel     string
		wantErr       error
	}{
		{name: "preferred exceeds RAM", options: []string{"small", "large"}, preferred: "large", totalMemory: 8, availableDisk: 100, wantModel: "small"},
		{name: "largest exceeds RAM", options: []string{"small", "large"}, totalMemory: 8, availableDisk: 100, wantModel: "small"},
		{name: "exact memory fit", options: []string{"large"}, preferred: "large", totalMemory: 16, availableDisk: 100, wantModel: "large"},
		{name: "no model fits RAM", options: []string{"small", "large"}, totalMemory: 1, availableDisk: 100, wantErr: ErrInsufficientMemoryForModel},
		{name: "no model fits disk", options: []string{"small", "large"}, totalMemory: 100, availableDisk: 1, wantErr: ErrInsufficientDiskSpaceForModel},
		{name: "memory defaults to disk size", options: []string{"implicit"}, totalMemory: 8, availableDisk: 100, wantErr: ErrInsufficientMemoryForModel},
		{name: "preferred exceeds VRAM", options: []string{"small", "large"}, preferred: "large", totalMemory: 100, availableDisk: 100, gpu: true, vram: uint64(8), wantModel: "small"},
		{name: "no model fits VRAM", options: []string{"small", "large"}, totalMemory: 100, availableDisk: 100, gpu: true, vram: uint64(1), wantErr: ErrInsufficientMemoryForModel},
		{name: "exact VRAM fit", options: []string{"large"}, totalMemory: 1, availableDisk: 100, gpu: true, vram: uint64(16), wantModel: "large"},
		{name: "unified memory uses RAM", options: []string{"small", "large"}, preferred: "large", totalMemory: 8, availableDisk: 100, gpu: true, wantModel: "small"},
	} {
		t.Run(test.name, func(t *testing.T) {
			modelsDir := t.TempDir()
			writeModelYAML(t, modelsDir, "small", "name: small\ndisk-size: 2\nrequired-memory: 4\n")
			writeModelYAML(t, modelsDir, "large", "name: large\ndisk-size: 4\nrequired-memory: 16\n")
			writeModelYAML(t, modelsDir, "implicit", "name: implicit\ndisk-size: 16\n")
			hostMachine := &machine.Machine{}
			hostMachine.Memory.TotalRam = test.totalMemory
			hostMachine.Disk = []disk.Disk{{Path: "/var/lib/snapd/snaps", Available: test.availableDisk}}
			engine := engines.Manifest{Model: engines.Model{Options: test.options}}
			if test.gpu {
				engine.Devices.Allof = []engines.Device{{Type: "gpu", VendorId: new(utils.HexInt(0x8086))}}
				hostMachine.PCIDevices = []pci.Device{{
					DeviceClass:          0x0300,
					VendorId:             0x8086,
					AdditionalProperties: map[string]any{"vram": test.vram},
				}}
			}
			selected, scoredModels, err := SelectModel(&Context{ModelsDir: modelsDir}, engine, test.preferred, hostMachine)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("error = %v, want %v", err, test.wantErr)
			}
			if selected != test.wantModel {
				t.Errorf("selected = %q, want %q", selected, test.wantModel)
			}
			for _, model := range scoredModels {
				if !model.CompatibilityReport.Compatible() && model.Score != 0 {
					t.Errorf("incompatible model %q has score %d", model.Name, model.Score)
				}
			}
		})
	}
}
