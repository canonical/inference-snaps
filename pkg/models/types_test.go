package models

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestParseManifest(t *testing.T) {
	testDataDir := filepath.Join("..", "..", "test_data", "models")

	entries, err := os.ReadDir(testDataDir)
	if err != nil {
		t.Fatalf("failed to read test_data/models directory: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		modelFile := filepath.Join(testDataDir, entry.Name(), "model.yaml")
		t.Run(entry.Name(), func(t *testing.T) {
			data, err := os.ReadFile(modelFile)
			if err != nil {
				t.Fatalf("failed to read %s: %v", modelFile, err)
			}

			var manifest Manifest
			dec := yaml.NewDecoder(bytes.NewReader(data))
			dec.KnownFields(true)
			if err := dec.Decode(&manifest); err != nil {
				t.Fatalf("failed to unmarshal %s: %v", modelFile, err)
			}

			t.Logf("parsed model manifest: id=%q name=%q", manifest.ID, manifest.Name)
		})
	}
}
