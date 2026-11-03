package runtimes

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const runtimeManifestWithSlashInServerName = `
name: test-runtime
servers:
  whisper/server:
    protocol: http
    base-path: /v1
`

func writeRuntimeManifest(t *testing.T, dir, runtimeName, contents string) {
	t.Helper()

	runtimeDir := filepath.Join(dir, runtimeName)
	if err := os.MkdirAll(runtimeDir, 0o755); err != nil {
		t.Fatalf("creating runtime dir: %v", err)
	}
	manifestPath := filepath.Join(runtimeDir, ManifestFilename)
	if err := os.WriteFile(manifestPath, []byte(contents), 0o644); err != nil {
		t.Fatalf("writing manifest: %v", err)
	}
}

func TestLoadManifestRejectsSlashInServerName(t *testing.T) {
	dir := t.TempDir()
	writeRuntimeManifest(t, dir, "test-runtime", runtimeManifestWithSlashInServerName)

	_, err := LoadManifest(dir, "test-runtime")
	if err == nil {
		t.Fatal("expected error for server name containing '/'")
	}
	if !strings.Contains(err.Error(), "invalid server name") {
		t.Fatalf("expected error about invalid server name, got: %v", err)
	}
}

func TestLoadManifestsRejectsSlashInServerName(t *testing.T) {
	dir := t.TempDir()
	writeRuntimeManifest(t, dir, "test-runtime", runtimeManifestWithSlashInServerName)

	_, err := LoadManifests(dir)
	if err == nil {
		t.Fatal("expected error for server name containing '/'")
	}
	if !strings.Contains(err.Error(), "invalid server name") {
		t.Fatalf("expected error about invalid server name, got: %v", err)
	}
}
