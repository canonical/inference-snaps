package runtimes

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestParseManifest(t *testing.T) {
	testDataDir := filepath.Join("..", "..", "test_data", "runtimes")

	entries, err := os.ReadDir(testDataDir)
	if err != nil {
		t.Fatalf("failed to read test_data/runtimes directory: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		runtimeFile := filepath.Join(testDataDir, entry.Name(), "runtime.yaml")
		t.Run(entry.Name(), func(t *testing.T) {
			data, err := os.ReadFile(runtimeFile)
			if err != nil {
				t.Fatalf("failed to read %s: %v", runtimeFile, err)
			}

			var manifest Manifest
			dec := yaml.NewDecoder(bytes.NewReader(data))
			dec.KnownFields(true)
			if err := dec.Decode(&manifest); err != nil {
				t.Fatalf("failed to unmarshal %s: %v", runtimeFile, err)
			}

			t.Logf("parsed runtime manifest: servers=%v components=%v", manifest.Servers, manifest.Components)
		})
	}
}
