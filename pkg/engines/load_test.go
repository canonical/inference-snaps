package engines

import (
	"errors"
	"testing"
)

func TestLoadManifest(t *testing.T) {
	enginesDir := "../../test_data/engines"

	const engineName = "intel-cpu"
	manifest, err := LoadManifest(enginesDir, engineName)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if manifest.Name != engineName {
		t.Fatalf("expected engine %q, got %q", engineName, manifest.Name)
	}

	_, err = LoadManifest(enginesDir, "nonexistent")
	if err == nil {
		t.Fatalf("expected error for nonexistent engine, got nil")
	}
	if !errors.Is(err, ErrManifestNotFound) {
		t.Fatalf("unexpected error for nonexistent engine: %s", err)
	}
}
