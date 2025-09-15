package engines

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func templateManifest() Manifest {
	memDisk := "1"
	manifest := Manifest{
		Name:        "test",
		Description: "test",
		Vendor:      "test",
		Grade:       "stable",
		Devices:     Devices{},
		Memory:      &memDisk,
		DiskSpace:   &memDisk,
		Components:  nil,
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
			manifestPath := filepath.Join(enginesDir, engine, "engine.yaml")
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

func TestGradeRequired(t *testing.T) {
	manifest := templateManifest()
	manifest.Grade = ""

	err := manifest.validate("test")
	if err == nil {
		t.Fatal("grade is required")
	}
	t.Log(err)

}

func TestGradeValid(t *testing.T) {
	manifest := templateManifest()

	t.Run("grade stable", func(t *testing.T) {
		manifest.Grade = "stable"

		err := manifest.validate("test")
		if err != nil {
			t.Fatalf("grade stable should be valid: %v", err)
		}
	})
	t.Run("grade devel", func(t *testing.T) {
		manifest.Grade = "devel"

		err := manifest.validate("test")
		if err != nil {
			t.Fatalf("grade devel should be valid: %v", err)
		}
	})
	t.Run("grade invalid", func(t *testing.T) {
		manifest.Grade = "invalid-grade"

		err := manifest.validate("test")
		if err == nil {
			t.Fatal("grade invalid")
		}
		t.Log(err)
	})

}

func TestMemoryValues(t *testing.T) {
	manifest := templateManifest()

	t.Run("valid GB", func(t *testing.T) {
		value := "1G"
		manifest.Memory = &value

		err := manifest.validate("test")
		if err != nil {
			t.Logf("memory should be valid: %v", err)
		}
	})

	t.Run("valid MB", func(t *testing.T) {
		value := "512M"
		manifest.Memory = &value

		err := manifest.validate("test")
		if err != nil {
			t.Logf("memory should be valid: %v", err)
		}
	})

	// Empty memory string in yaml is parsed as nil, which we interpret as unset, which is valid

	t.Run("not numeric", func(t *testing.T) {
		value := "abc"
		manifest.Memory = &value

		err := manifest.validate("test")
		if err == nil {
			t.Fatal("non-numeric memory should be invalid")
		}
		t.Log(err)
	})

}

func TestDiskValues(t *testing.T) {
	manifest := templateManifest()

	t.Run("valid GB", func(t *testing.T) {
		value := "1G"
		manifest.DiskSpace = &value

		err := manifest.validate("test")
		if err != nil {
			t.Logf("disk should be valid: %v", err)
		}
	})

	t.Run("valid MB", func(t *testing.T) {
		value := "512M"
		manifest.DiskSpace = &value

		err := manifest.validate("test")
		if err != nil {
			t.Logf("disk should be valid: %v", err)
		}
	})

	// Empty string in yaml is parsed as nil, which we interpret as unset, which is valid

	t.Run("not numeric", func(t *testing.T) {
		value := "abc"
		manifest.DiskSpace = &value

		err := manifest.validate("test")
		if err == nil {
			t.Fatal("non-numeric disk should be invalid")
		}
		t.Log(err)
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
