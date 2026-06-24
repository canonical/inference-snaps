package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/canonical/inference-snaps-cli/v2/pkg/storage"
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

func TestModelStatus_ModelNamePresent(t *testing.T) {
	modelsDir := t.TempDir()
	writeModelYAML(t, modelsDir, "my-model", `environment:
  - MODEL_NAME=my-model
  - OTHER_VAR=value
`)
	ctx := makeModelStatusCtx(t, modelsDir, "my-model")

	status, err := ModelStatus(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status["name"] != "my-model" {
		t.Errorf("expected name %q, got %q", "my-model", status["name"])
	}
}

func TestModelStatus_NoModelName(t *testing.T) {
	modelsDir := t.TempDir()
	writeModelYAML(t, modelsDir, "my-model", `environment:
  - OTHER_VAR=value
`)
	ctx := makeModelStatusCtx(t, modelsDir, "my-model")

	status, err := ModelStatus(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := status["name"]; ok {
		t.Errorf("expected no 'name' key, but got: %q", status["name"])
	}
}

func TestModelStatus_ModelNameWithEqualsInValue(t *testing.T) {
	modelsDir := t.TempDir()
	writeModelYAML(t, modelsDir, "my-model", `environment:
  - MODEL_NAME=my=model
`)
	ctx := makeModelStatusCtx(t, modelsDir, "my-model")

	status, err := ModelStatus(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status["name"] != "my=model" {
		t.Errorf("expected name %q, got %q", "my=model", status["name"])
	}
}

func TestModelStatus_InvalidEnvVar(t *testing.T) {
	modelsDir := t.TempDir()
	writeModelYAML(t, modelsDir, "my-model", `environment:
  - INVALID_NO_EQUALS
`)
	ctx := makeModelStatusCtx(t, modelsDir, "my-model")

	_, err := ModelStatus(ctx)
	if err == nil {
		t.Fatal("expected error for invalid env var, got nil")
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


func TestGetModelByNameOrId(t *testing.T) {
	tests := []struct {
		name         string
		activeEngine string
		modelYAML    string // empty means don't write a model manifest
		engineYAML   string // empty means don't write an engine manifest
		query        string
		wantID       string // non-empty: expect this ID in the returned manifest
		wantName     string // non-empty: expect this Name in the returned manifest
		wantErr      bool   // true: expect any non-nil error
	}{
		{
			name:         "found by name",
			activeEngine: "my-engine",
			modelYAML:    "id: my-model-id\nname: my-model\ndisk-size: 1G\n",
			engineYAML:   "name: my-engine\nmodel:\n  options:\n    - my-model-id\n",
			query:        "my-model",
			wantID:       "my-model-id",
		},
		{
			name:         "found by id",
			activeEngine: "my-engine",
			modelYAML:    "id: my-model-id\nname: my-model\ndisk-size: 1G\n",
			engineYAML:   "name: my-engine\nmodel:\n  options:\n    - my-model-id\n",
			query:        "my-model-id",
			wantName:     "my-model",
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

			manifest, err := GetModelByNameOrId(ctx, tc.query)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.wantID != "" && manifest.ID != tc.wantID {
				t.Errorf("expected ID %q, got %q", tc.wantID, manifest.ID)
			}
			if tc.wantName != "" && manifest.Name != tc.wantName {
				t.Errorf("expected Name %q, got %q", tc.wantName, manifest.Name)
			}
		})
	}
}

