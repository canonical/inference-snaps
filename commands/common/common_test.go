package common

import (
	"os"
	"path/filepath"
	"testing"
)

// writeRuntimeYAML replaces the runtime.yaml for the named runtime inside the given runtimesDir.
func writeRuntimeYAML(t *testing.T, runtimesDir, name, content string) {
	t.Helper()
	dir := filepath.Join(runtimesDir, name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("MkdirAll %s: %v", dir, err)
	}
	if err := os.WriteFile(filepath.Join(dir, "runtime.yaml"), []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile runtime.yaml: %v", err)
	}
}

// writeEngineYAML writes an engine.yaml for the named engine inside the given enginesDir.
func writeEngineYAML(t *testing.T, enginesDir, name, content string) {
	t.Helper()
	dir := filepath.Join(enginesDir, name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("MkdirAll %s: %v", dir, err)
	}
	if err := os.WriteFile(filepath.Join(dir, "engine.yaml"), []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile engine.yaml: %v", err)
	}
}
