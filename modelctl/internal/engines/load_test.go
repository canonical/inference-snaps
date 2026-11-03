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
	if manifest.Experimental != nil {
		t.Fatalf("expected experimental to be unset (nil) for %q", engineName)
	}

	experimentalManifest, err := LoadManifest(enginesDir, "cpu-exptl")
	if err != nil {
		t.Fatalf("expected no error for cpu-exptl, got %v", err)
	}
	if !experimentalManifest.IsExperimental() {
		t.Fatalf("expected cpu-exptl experimental to be true")
	}

	_, err = LoadManifest(enginesDir, "nonexistent")
	if err == nil {
		t.Fatalf("expected error for nonexistent engine, got nil")
	}
	if !errors.Is(err, ErrManifestNotFound) {
		t.Fatalf("unexpected error for nonexistent engine: %s", err)
	}
}

func TestLoadManifestsExperimentalFalse(t *testing.T) {
	enginesDir := "../../test_data/engines"

	manifests, err := LoadManifests(enginesDir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(manifests) == 0 {
		t.Fatal("expected at least one manifest")
	}

	foundCPU := false
	foundCPUExptl := false

	for _, manifest := range manifests {
		switch manifest.Name {
		case "cpu":
			foundCPU = true
			if manifest.Experimental != nil {
				t.Fatalf("expected cpu experimental to be unset (nil)")
			}
		case "cpu-exptl":
			foundCPUExptl = true
			if !manifest.IsExperimental() {
				t.Fatalf("expected cpu-exptl experimental to be true")
			}
		}
	}

	if !foundCPU {
		t.Fatal("expected cpu manifest to be present")
	}
	if !foundCPUExptl {
		t.Fatal("expected cpu-exptl manifest to be present")
	}
}
