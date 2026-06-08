package engines

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func templateManifest() Manifest {
	manifest := Manifest{
		Name:         "test",
		Description:  "test",
		Vendor:       "test",
		Experimental: nil,
		Model: Model{
			Default: "26b-q4-k-m-gguf",
			Options: []string{"26b-q4-k-m-gguf", "30b-a3b-q4-k-m-gguf "},
		},
		Configurations: map[string]interface{}{
			"engine": "test",
			"model":  "test",
		},
	}
	return manifest
}

func TestManifestFiles(t *testing.T) {
	enginesDir := "../../test_data/engines"

	entries, err := os.ReadDir(enginesDir)
	if err != nil {
		t.Fatalf("Failed reading engines directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			engine := entry.Name()
			manifestPath := filepath.Join(enginesDir, engine, ManifestFilename)
			t.Run(engine, func(t *testing.T) {
				err = Validate(manifestPath)
				if err != nil {
					t.Fatalf("%s: %v", engine, err)
				}
			})
		}
	}
}

func TestManifestEmpty(t *testing.T) {
	data := ""
	err := validateManifestYaml("", []byte(data))
	if err == nil {
		t.Fatal("Empty yaml should fail")
	}
	t.Log(err)
}

func TestUnknownField(t *testing.T) {
	data, _ := yaml.Marshal(templateManifest())
	data = append(data, []byte("unknown-field: test\n")...)

	err := validateManifestYaml("test", data)
	if err == nil {
		t.Fatal("Unknown field should fail")
	}
	t.Log(err)
}

func TestNameRequired(t *testing.T) {
	manifest := templateManifest()
	manifest.Name = ""

	err := manifest.validate("test")
	if err == nil {
		t.Fatal("name field is required")
	}
	t.Log(err)

}

func TestDescriptionRequired(t *testing.T) {
	manifest := templateManifest()
	manifest.Description = ""

	err := manifest.validate("test")
	if err == nil {
		t.Fatal("description is required")
	}
	t.Log(err)

}

func TestVendorRequired(t *testing.T) {
	manifest := templateManifest()
	manifest.Vendor = ""

	err := manifest.validate("test")
	if err == nil {
		t.Fatal("vendor is required")
	}
	t.Log(err)

}

func TestExperimentalValid(t *testing.T) {
	manifest := templateManifest()

	t.Run("experimental false", func(t *testing.T) {
		value := false
		manifest.Experimental = &value

		err := manifest.validate("test")
		if err != nil {
			t.Fatalf("experimental false should be valid: %v", err)
		}
	})
	t.Run("experimental true", func(t *testing.T) {
		value := true
		manifest.Experimental = &value

		err := manifest.validate("test")
		if err != nil {
			t.Fatalf("experimental true should be valid: %v", err)
		}
	})
}

func TestConfig(t *testing.T) {
	manifest := templateManifest()

	t.Run("config is primitive", func(t *testing.T) {
		manifest.Configurations = map[string]interface{}{"model": true}
		err := manifest.validate("test")
		if err != nil {
			t.Fatalf("primitive model field should be valid: %v", err)
		}
	})

	t.Run("config is not primitive", func(t *testing.T) {
		manifest.Configurations = map[string]interface{}{"model": []string{"one", "two"}}
		err := manifest.validate("test")
		if err == nil {
			t.Fatal("non-primitive model field should be invalid")
		}
		t.Log(err)
	})
}
