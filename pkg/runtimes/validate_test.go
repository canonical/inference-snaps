package runtimes

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func templateManifest() Manifest {
	manifest := Manifest{
		Name: "test",
		Servers: map[string]Server{
			"openai": {
				Protocol: "http",
				BasePath: "/v1",
			},
		},
		Environment: []string{"PATH=$PATH:$SNAP_COMPONENTS/runtime-test/bin"},
		Components:  []string{"runtime-test"},
	}
	return manifest
}

func TestManifestFiles(t *testing.T) {
	runtimesDir := "../../test_data/runtimes"

	entries, err := os.ReadDir(runtimesDir)
	if err != nil {
		t.Fatalf("Failed reading runtimes directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			runtime := entry.Name()
			manifestPath := filepath.Join(runtimesDir, runtime, ManifestFilename)
			t.Run(runtime, func(t *testing.T) {
				err = Validate(manifestPath)
				if err != nil {
					t.Fatalf("%s: %v", runtime, err)
				}
			})
		}
	}
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

func TestNameMatch(t *testing.T) {
	manifest := templateManifest()
	manifest.Name = "different-name"

	err := manifest.validate("test")
	if err == nil {
		t.Fatal("runtime directory name should match name in manifest")
	}
}

func TestServersRequired(t *testing.T) {
	manifest := templateManifest()
	manifest.Servers = nil

	err := manifest.validate("test")
	if err == nil {
		t.Fatal("servers field is required")
	}
}

func TestServerProtocolRequired(t *testing.T) {
	manifest := templateManifest()
	manifest.Servers = map[string]Server{
		"openai": {
			BasePath: "/v1",
		},
	}

	err := manifest.validate("test")
	if err == nil {
		t.Fatal("server protocol field is required")
	}
}

func TestServerBasePathRequired(t *testing.T) {
	manifest := templateManifest()
	manifest.Servers = map[string]Server{
		"openai": {
			Protocol: "http",
		},
	}

	err := manifest.validate("test")
	if err == nil {
		t.Fatal("server base-path field is required")
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
	for _, env := range []string{"PATH", "=value", "1BAD=value", "BAD NAME=value"} {
		manifest.Environment = []string{env}
		if err := manifest.validate("test"); err == nil {
			t.Fatalf("environment entry %q should be invalid", env)
		}
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
