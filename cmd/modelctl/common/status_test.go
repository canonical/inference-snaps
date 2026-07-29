package common

import (
	"strings"
	"testing"

	"github.com/canonical/inference-snaps-cli/v2/pkg/snap"
	"github.com/canonical/inference-snaps-cli/v2/pkg/storage"
)

func TestStatusStruct(t *testing.T) {
	enginesDir := t.TempDir()
	runtimesDir := t.TempDir()
	modelsDir := t.TempDir()

	writeEngineYAML(t, enginesDir, "my-engine", `name: my-engine
runtime: my-runtime
`)
	writeRuntimeYAML(t, runtimesDir, "my-runtime", `servers:
  openai:
    protocol: http
    base-path: /v1
`)
	writeModelYAML(t, modelsDir, "my-model", `environment:
  - MODEL_NAME=my-model
`)

	cache := storage.NewMockCache()
	if err := cache.SetActiveEngine("my-engine"); err != nil {
		t.Fatalf("SetActiveEngine: %v", err)
	}
	if err := cache.SetActiveModel("my-model"); err != nil {
		t.Fatalf("SetActiveModel: %v", err)
	}

	config := storage.NewMockConfig()
	config.Set("http.port", "8080", storage.UserConfig)
	config.Set("http.host", "127.0.0.1", storage.UserConfig)

	ctx := &Context{
		EnginesDir:  enginesDir,
		RuntimesDir: runtimesDir,
		ModelsDir:   modelsDir,
		Cache:       cache,
		Config:      config,
		Snap:        snap.MockWithServiceStatuses(map[string]string{"llama-server": "active"}),
	}

	status, err := SnapStatus(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if status.Engine != "my-engine" {
		t.Errorf("engine: got %q, want %q", status.Engine, "my-engine")
	}
	if status.Services["llama-server"] != "active" {
		t.Errorf("service status: got %q, want %q", status.Services["llama-server"], "active")
	}
	if status.Entrypoints["openai"].Url != "http://127.0.0.1:8080/v1" {
		t.Errorf("openai entrypoint URL: got %q, want %q", status.Entrypoints["openai"].Url, "http://127.0.0.1:8080/v1")
	}
	if status.Model["name"] != "my-model" {
		t.Errorf("model name: got %q, want %q", status.Model["name"], "my-model")
	}
}

func TestStatusStruct_NoActiveEngine(t *testing.T) {
	cache := storage.NewMockCache()
	ctx := &Context{
		Cache:  cache,
		Config: storage.NewMockConfig(),
		Snap:   snap.Mock(),
	}

	_, err := SnapStatus(ctx)
	if err == nil {
		t.Fatal("expected error for no active engine, got nil")
	}
	if !strings.Contains(err.Error(), "no active engine") {
		t.Errorf("unexpected error: %v", err)
	}
}
