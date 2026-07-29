package runtimes

import (
	"os"
	"path/filepath"
	"strings"
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

func TestServerProtocolInvalid(t *testing.T) {
	manifest := templateManifest()

	invalidProtocols := []string{"ftp", "ssh", "telnet", "http+invalid"}
	for _, protocol := range invalidProtocols {
		manifest.Servers = map[string]Server{
			"openai": {
				Protocol: protocol,
				BasePath: "/v1",
			},
		}

		err := manifest.validate("test")
		if err == nil {
			t.Fatalf("server protocol %q should be invalid", protocol)
		}
		if !strings.Contains(err.Error(), "invalid protocol") {
			t.Fatalf("expected error about invalid protocol for %q, got: %v", protocol, err)
		}
	}
}

func TestServerProtocolValid(t *testing.T) {
	manifest := templateManifest()

	validProtocols := []string{"http", "https", "http+unix", "https+unix", "ws", "wss", "ws+unix", "wss+unix"}
	for _, protocol := range validProtocols {
		manifest.Servers = map[string]Server{
			"openai": {
				Protocol: protocol,
				BasePath: "/v1",
			},
		}

		err := manifest.validate("test")
		if err != nil {
			t.Fatalf("server protocol %q should be valid, got error: %v", protocol, err)
		}
	}
}

func TestServerNamespaceInvalid(t *testing.T) {
	manifest := templateManifest()

	invalidNamespaces := []string{"-namespace", "1namespace", "Name-Space", "name space", "-", "1"}
	for _, namespace := range invalidNamespaces {
		manifest.Servers = map[string]Server{
			"openai": {
				Protocol:  "http",
				Namespace: namespace,
			},
		}

		err := manifest.validate("test")
		if err == nil {
			t.Fatalf("server namespace %q should be invalid", namespace)
		}
		if !strings.Contains(err.Error(), "invalid namespace") {
			t.Fatalf("expected error about invalid namespace for %q, got: %v", namespace, err)
		}
	}
}

func TestServerNamespaceValid(t *testing.T) {
	manifest := templateManifest()

	validNamespaces := []string{"", "namespace", "myservice", "logger", "kserve", "abc", "my-service", "service123", "my-service-1"}
	for _, namespace := range validNamespaces {
		manifest.Servers = map[string]Server{
			"openai": {
				Protocol:  "http",
				Namespace: namespace,
			},
		}

		err := manifest.validate("test")
		if err != nil {
			t.Fatalf("server namespace %q should be valid, got error: %v", namespace, err)
		}
	}
}

func TestServerBasePathOptional(t *testing.T) {
	manifest := templateManifest()
	manifest.Servers = map[string]Server{
		"openai": {
			Protocol: "http",
		},
	}

	err := manifest.validate("test")
	if err != nil {
		t.Fatalf("server base-path should be optional, got error: %v", err)
	}
}

func TestServerBasePathValidUrlPath(t *testing.T) {
	manifest := templateManifest()

	validPaths := []string{"", "/", "/v1", "/api/v2/stream"}
	for _, basePath := range validPaths {
		manifest.Servers = map[string]Server{
			"openai": {
				Protocol: "http",
				BasePath: basePath,
			},
		}

		if err := manifest.validate("test"); err != nil {
			t.Fatalf("base-path %q should be valid, got error: %v", basePath, err)
		}
	}
}

func TestServerBasePathRejectsNonPathURLForms(t *testing.T) {
	testCases := []struct {
		name            string
		basePath        string
		wantErrContains string
	}{
		{
			name:            "missing leading slash",
			basePath:        "v1",
			wantErrContains: "must start with '/'",
		},
		{
			name:            "full URL",
			basePath:        "http://example.com/v1",
			wantErrContains: "must start with '/'",
		},
		{
			name:            "query string",
			basePath:        "/v1?model=x",
			wantErrContains: "query and fragment are not allowed",
		},
		{
			name:            "fragment",
			basePath:        "/v1#section",
			wantErrContains: "query and fragment are not allowed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			manifest := templateManifest()
			manifest.Servers = map[string]Server{
				"openai": {
					Protocol: "http",
					BasePath: tc.basePath,
				},
			}

			err := manifest.validate("test")
			if err == nil {
				t.Fatalf("base-path %q should be invalid", tc.basePath)
			}
			if !strings.Contains(err.Error(), tc.wantErrContains) {
				t.Fatalf("base-path %q error = %q, want it to contain %q", tc.basePath, err.Error(), tc.wantErrContains)
			}
		})
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
