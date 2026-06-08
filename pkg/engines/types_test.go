package engines

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestParseManifest(t *testing.T) {
	testDataDir := filepath.Join("..", "..", "test_data", "engines")

	entries, err := os.ReadDir(testDataDir)
	if err != nil {
		t.Fatalf("failed to read test_data/engines directory: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		engineFile := filepath.Join(testDataDir, entry.Name(), "engine.yaml")
		t.Run(entry.Name(), func(t *testing.T) {
			data, err := os.ReadFile(engineFile)
			if err != nil {
				t.Fatalf("failed to read %s: %v", engineFile, err)
			}

			var manifest Manifest
			dec := yaml.NewDecoder(bytes.NewReader(data))
			dec.KnownFields(true)
			if err := dec.Decode(&manifest); err != nil {
				t.Fatalf("failed to unmarshal %s: %v", engineFile, err)
			}

			t.Logf("parsed engine manifest: name=%q vendor=%q experimental=%b", manifest.Name, manifest.Vendor, manifest.Experimental)
		})
	}
}
