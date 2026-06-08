package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/storage"
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
