package models

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"gopkg.in/yaml.v3"
)

func templateManifest() Manifest {
	manifest := Manifest{
		Name:         "test",
		Alias:        "test-alias",
		Description:  "test",
		ModelCardUrl: "https://example.com/model-card",
		Quantization: "Q4_K_M",
		Capabilities: []string{"text"},
		DiskSize:     "6G",
		Components:   []string{"test-component"},
		Environment:  []string{"MODEL_FILE=/tmp/model.gguf"},
	}
	return manifest
}

func TestManifestFiles(t *testing.T) {
	modelsDir := "../../test_data/models"

	_, err := os.ReadDir(modelsDir)
	if err != nil {
		t.Fatalf("Failed reading models directory: %v", err)
	}
	model := "26b-q4-k-m-gguf"
	manifestPath := filepath.Join(modelsDir, model, ManifestFilename)
	t.Run(model, func(t *testing.T) {
		err = Validate(manifestPath)
		if err != nil {
			t.Fatalf("%s: %v", model, err)
		}
	})
}

func TestManifestYamlUnsupportedCapability(t *testing.T) {
	modelsDir := "../../test_data/models"

	_, err := os.ReadDir(modelsDir)
	if err != nil {
		t.Fatalf("Failed reading models directory: %v", err)
	}
	model := "30b-a3b-q4-k-m-gguf"
	manifestPath := filepath.Join(modelsDir, model, ManifestFilename)
	t.Run(model, func(t *testing.T) {
		err = Validate(manifestPath)
		if err == nil {
			t.Fatalf("%s: expected an error for unsupported capability, got nil", model)
		}
	})
}
func TestManifestEmpty(t *testing.T) {
	data := ""
	_, err := parseManifest([]byte(data))
	if err == nil {
		t.Fatal("Empty yaml should fail")
	}
}

func TestUnknownField(t *testing.T) {
	data, _ := yaml.Marshal(templateManifest())
	data = append(data, []byte("unknown-field: test\n")...)

	_, err := parseManifest(data)
	if err == nil {
		t.Fatal("Unknown field should fail")
	}
}

func TestNameRequired(t *testing.T) {
	manifest := templateManifest()
	manifest.Name = ""

	err := manifest.validate("test")
	if err == nil {
		t.Fatal("name field is required")
	}
}

func TestDescriptionRequired(t *testing.T) {
	manifest := templateManifest()
	manifest.Description = ""

	err := manifest.validate("test")
	if err == nil {
		t.Fatal("description field is required")
	}
}

func TestModelCardUrlOptional(t *testing.T) {
	manifest := templateManifest()
	manifest.ModelCardUrl = ""

	err := manifest.validate("test")
	if err != nil {
		t.Fatalf("model-card-url field is optional, got error: %v", err)
	}
}

func TestModelCardUrlInvalid(t *testing.T) {
	manifest := templateManifest()
	manifest.ModelCardUrl = "not-a-url"

	err := manifest.validate("test")
	if err == nil {
		t.Fatal("expected an error for invalid model-card-url")
	}
}

func TestQuantizationOptional(t *testing.T) {
	manifest := templateManifest()
	manifest.Quantization = ""

	err := manifest.validate("test")
	if err != nil {
		t.Fatalf("quantization field is optional, got error: %v", err)
	}
}

func TestCapabilitiesOptional(t *testing.T) {
	manifest := templateManifest()
	manifest.Capabilities = nil

	err := manifest.validate("test")
	if err != nil {
		t.Fatalf("capabilities field is optional, got error: %v", err)
	}
}

func TestCapabilityTextEmbedding(t *testing.T) {
	manifest := templateManifest()
	manifest.Capabilities = []string{"text-embedding"}

	err := manifest.validate("test")
	if err != nil {
		t.Fatalf("text-embedding is a valid capability, got error: %v", err)
	}
}

func TestDiskSizeRequired(t *testing.T) {
	manifest := templateManifest()
	manifest.DiskSize = ""

	err := manifest.validate("test")
	if err == nil {
		t.Fatal("disk-size field is required")
	}
}

func TestComponentsOptional(t *testing.T) {
	manifest := templateManifest()
	manifest.Components = nil

	err := manifest.validate("test")
	if err != nil {
		t.Fatalf("components field is optional, got error: %v", err)
	}
}

func TestComponentEmptyName(t *testing.T) {
	manifest := templateManifest()
	manifest.Components = []string{""}

	err := manifest.validate("test")
	if err == nil {
		t.Fatal("empty component name should fail")
	}
}

func TestEnvironmentOptional(t *testing.T) {
	manifest := templateManifest()
	manifest.Environment = nil

	err := manifest.validate("test")
	if err != nil {
		t.Fatalf("environment field is optional, got error: %v", err)
	}
}

func TestEnvironmentInvalidSyntax(t *testing.T) {
	manifest := templateManifest()
	for _, env := range []string{"MODEL_FILE", "=value", "1BAD=value", "BAD NAME=value"} {
		manifest.Environment = []string{env}
		if err := manifest.validate("test"); err == nil {
			t.Fatalf("environment entry %q should be invalid", env)
		}
	}
}

func TestModelNameMatch(t *testing.T) {
	manifest := templateManifest()
	manifest.Name = "different-name"

	err := manifest.validate("test")
	if err == nil {
		t.Fatal("model directory name should match name in manifest")
	}
}

func TestValidateAliasCollisionSameEngine(t *testing.T) {
	modelA := templateManifest()
	modelA.Name = "model-a"
	modelA.Alias = "shared"

	modelB := templateManifest()
	modelB.Name = "model-b"
	modelB.Alias = "shared"

	engineManifests := []engines.Manifest{
		{Name: "engine-1", Model: engines.Model{Options: []string{"model-a", "model-b"}}},
	}
	if err := validateAliasesPerEngine(engineManifests, []Manifest{modelA, modelB}); err == nil {
		t.Fatalf("expected alias collision within engine, got no error")
	}
}

func TestValidateAliasSharedAcrossEngines(t *testing.T) {
	// The same alias is allowed when the two models never share an engine.
	modelA := templateManifest()
	modelA.Name = "model-a"
	modelA.Alias = "shared"

	modelB := templateManifest()
	modelB.Name = "model-b"
	modelB.Alias = "shared"

	engineManifests := []engines.Manifest{
		{Name: "engine-1", Model: engines.Model{Options: []string{"model-a"}}},
		{Name: "engine-2", Model: engines.Model{Options: []string{"model-b"}}},
	}
	if err := validateAliasesPerEngine(engineManifests, []Manifest{modelA, modelB}); err != nil {
		t.Fatalf("expected no collision across engines, got: %v", err)
	}
}

func TestValidateAliasUniqueWithinEngine(t *testing.T) {
	modelA := templateManifest()
	modelA.Name = "model-a"
	modelA.Alias = "alias-a"

	modelB := templateManifest()
	modelB.Name = "model-b"
	modelB.Alias = "alias-b"

	engineManifests := []engines.Manifest{
		{Name: "engine-1", Model: engines.Model{Options: []string{"model-a", "model-b"}}},
	}
	if err := validateAliasesPerEngine(engineManifests, []Manifest{modelA, modelB}); err != nil {
		t.Fatalf("expected no collision, got: %v", err)
	}
}
