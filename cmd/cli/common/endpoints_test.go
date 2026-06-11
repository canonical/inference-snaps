package common

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/canonical/inference-snaps-cli/v2/pkg/runtimes"
	"github.com/canonical/inference-snaps-cli/v2/pkg/storage"
)

// writeRuntimeYAML replaces the runtime.yaml for the named runtime inside the given runtimesDir.
func writeRuntimeYAML(t *testing.T, runtimesDir, name, content string) {
	t.Helper()
	dir := filepath.Join(runtimesDir, name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("MkdirAll %s: %v", dir, err)
	}
	if err := os.WriteFile(filepath.Join(dir, "runtime.yaml"), []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile runtime.yaml: %v", err)
	}
}

// writeEngineYAML writes an engine.yaml for the named engine inside the given enginesDir.
func writeEngineYAML(t *testing.T, enginesDir, name, content string) {
	t.Helper()
	dir := filepath.Join(enginesDir, name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("MkdirAll %s: %v", dir, err)
	}
	if err := os.WriteFile(filepath.Join(dir, "engine.yaml"), []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile engine.yaml: %v", err)
	}
}

func TestServerEndpoints(t *testing.T) {
	testCases := []struct {
		name            string
		engineYAML      string
		runtimeYAML     string
		wantEndpoints   map[string]string
		wantErrContains string
	}{
		{
			name: "multiple servers",
			engineYAML: `name: test-engine
runtime: test-runtime
`,
			runtimeYAML: `servers:
  openai:
    protocol: http
    base-path: /v1
  kserve:
    protocol: https
    base-path: /v2
`,
			wantEndpoints: map[string]string{
				"openai": "http://127.0.0.1:8080/v1",
				"kserve": "https://127.0.0.1:8080/v2",
				"webui":  "http://192.0.2.1:8080/",
			},
		},
		{
			name: "unsupported protocol",
			engineYAML: `name: test-engine
runtime: test-runtime
`,
			runtimeYAML: `servers:
  openai:
    protocol: ftp
`,
			wantErrContains: "unsupported protocol",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("ADDITIONAL_FEATURES", "webui")

			enginesDir := t.TempDir()
			runtimesDir := t.TempDir()

			const (
				engineName  = "test-engine"
				runtimeName = "test-runtime"
			)

			writeEngineYAML(t, enginesDir, engineName, tc.engineYAML)
			writeRuntimeYAML(t, runtimesDir, runtimeName, tc.runtimeYAML)

			cache := storage.NewMockCache()
			if err := cache.SetActiveEngine(engineName); err != nil {
				t.Fatalf("SetActiveEngine: %v", err)
			}

			config := storage.NewMockConfig()
			config.Set("http.port", "8080", storage.UserConfig)
			config.Set("http.host", "127.0.0.1", storage.UserConfig)
			config.Set("webui.http.port", "8080", storage.UserConfig)
			config.Set("webui.http.host", "192.0.2.1", storage.UserConfig)

			ctx := &Context{
				EnginesDir:  enginesDir,
				RuntimesDir: runtimesDir,
				Config:      config,
				Cache:       cache,
			}

			got, err := ServerEndpoints(ctx)
			if tc.wantErrContains != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tc.wantErrContains)
				}
				if !strings.Contains(err.Error(), tc.wantErrContains) {
					t.Fatalf("got error %q, want it to contain %q", err.Error(), tc.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			for name, wantURL := range tc.wantEndpoints {
				gotURL, found := got[name]
				if !found {
					t.Fatalf("missing endpoint %q", name)
				}
				if gotURL != wantURL {
					t.Fatalf("endpoint %q: got %q, want %q", name, gotURL, wantURL)
				}
			}
		})
	}
}

func TestServerHttpUrl(t *testing.T) {
	testCases := []struct {
		name    string
		server  runtimes.Server
		host    string
		setHost bool
		want    string
	}{
		{
			name: "default base path",
			server: runtimes.Server{
				Protocol: "http",
			},
			host:    "0.0.0.0",
			setHost: true,
			want:    "http://0.0.0.0:8080/",
		},
		{
			name: "custom base path",
			server: runtimes.Server{
				Protocol: "http",
				BasePath: "/v1",
			},
			host:    "127.0.0.1",
			setHost: true,
			want:    "http://127.0.0.1:8080/v1",
		},
		{
			name: "https protocol",
			server: runtimes.Server{
				Protocol: "https",
				BasePath: "/v3",
			},
			host:    "0.0.0.0",
			setHost: true,
			want:    "https://0.0.0.0:8080/v3",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := storage.NewMockConfig()
			config.Set("http.port", "8080", storage.UserConfig)
			if tc.setHost {
				config.Set("http.host", tc.host, storage.UserConfig)
			}
			ctx := &Context{
				Config: config,
			}

			got, err := serverHttpUrl(ctx, tc.server)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}
