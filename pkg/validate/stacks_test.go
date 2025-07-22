package validate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/canonical/stack-utils/pkg/types"
	"gopkg.in/yaml.v3"
)

func templateManifest() types.Stack {
	memDisk := "1"
	manifest := types.Stack{
		Name:        "test",
		Description: "test",
		Vendor:      "test",
		Grade:       "stable",
		Devices:     types.StackDevices{},
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
	stacksDir := "../../test_data/stacks"

	entries, err := os.ReadDir(stacksDir)
	if err != nil {
		t.Fatalf("Failed reading stacks dir: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			stack := entry.Name()
			stackPath := filepath.Join(stacksDir, stack, "stack.yaml")
			t.Run(stack, func(t *testing.T) {
				err = Stack(stackPath)
				if err != nil {
					t.Fatalf("%s: %v", stack, err)
				}
			})
		}
	}
}

func TestManifestEmpty(t *testing.T) {
	data := ""
	err := validateStackYaml("", []byte(data))
	if err == nil {
		t.Fatal("Empty yaml should fail")
	}
	t.Log(err)
}

func TestUnknownField(t *testing.T) {
	data, _ := yaml.Marshal(templateManifest())
	data = append(data, []byte("unknown-field: test\n")...)

	err := validateStackYaml("test", data)
	if err == nil {
		t.Fatal("Unknown field should fail")
	}
	t.Log(err)
}

func TestNameRequired(t *testing.T) {
	manifest := templateManifest()
	manifest.Name = ""

	err := validateStackStruct("test", manifest)
	if err == nil {
		t.Fatal("name field is required")
	}
	t.Log(err)

}

func TestDescriptionRequired(t *testing.T) {
	manifest := templateManifest()
	manifest.Description = ""

	err := validateStackStruct("test", manifest)
	if err == nil {
		t.Fatal("description is required")
	}
	t.Log(err)

}

func TestVendorRequired(t *testing.T) {
	manifest := templateManifest()
	manifest.Vendor = ""

	err := validateStackStruct("test", manifest)
	if err == nil {
		t.Fatal("vendor is required")
	}
	t.Log(err)

}

func TestGradeRequired(t *testing.T) {
	manifest := templateManifest()
	manifest.Grade = ""

	err := validateStackStruct("test", manifest)
	if err == nil {
		t.Fatal("grade is required")
	}
	t.Log(err)

}

func TestGradeValid(t *testing.T) {
	manifest := templateManifest()

	t.Run("grade stable", func(t *testing.T) {
		manifest.Grade = "stable"

		err := validateStackStruct("test", manifest)
		if err != nil {
			t.Fatalf("grade stable should be valid: %v", err)
		}
	})
	t.Run("grade devel", func(t *testing.T) {
		manifest.Grade = "devel"

		err := validateStackStruct("test", manifest)
		if err != nil {
			t.Fatalf("grade devel should be valid: %v", err)
		}
	})
	t.Run("grade invalid", func(t *testing.T) {
		manifest.Grade = "invalid-grade"

		err := validateStackStruct("test", manifest)
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

		err := validateStackStruct("test", manifest)
		if err != nil {
			t.Logf("memory should be valid: %v", err)
		}
	})

	t.Run("valid MB", func(t *testing.T) {
		value := "512M"
		manifest.Memory = &value

		err := validateStackStruct("test", manifest)
		if err != nil {
			t.Logf("memory should be valid: %v", err)
		}
	})

	// Empty memory string in yaml is parsed as nil, which we interpret as unset, which is valid

	t.Run("not numeric", func(t *testing.T) {
		value := "abc"
		manifest.Memory = &value

		err := validateStackStruct("test", manifest)
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

		err := validateStackStruct("test", manifest)
		if err != nil {
			t.Logf("disk should be valid: %v", err)
		}
	})

	t.Run("valid MB", func(t *testing.T) {
		value := "512M"
		manifest.DiskSpace = &value

		err := validateStackStruct("test", manifest)
		if err != nil {
			t.Logf("disk should be valid: %v", err)
		}
	})

	// Empty string in yaml is parsed as nil, which we interpret as unset, which is valid

	t.Run("not numeric", func(t *testing.T) {
		value := "abc"
		manifest.DiskSpace = &value

		err := validateStackStruct("test", manifest)
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
		err := validateStackStruct("test", manifest)
		if err != nil {
			t.Fatalf("primitive model field should be valid: %v", err)
		}
	})

	t.Run("config is not primitive", func(t *testing.T) {
		manifest.Configurations = map[string]interface{}{"model": []string{"one", "two"}}
		err := validateStackStruct("test", manifest)
		if err == nil {
			t.Fatal("non-primitive model field should be invalid")
		}
		t.Log(err)
	})
}
