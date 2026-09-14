package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/common"
	"github.com/canonical/inference-snaps-cli/v2/pkg/storage"
)

func testRunContext(t *testing.T, runtimeYAML string) *common.Context {
	t.Helper()

	base := t.TempDir()
	enginesDir := filepath.Join(base, "engines")
	runtimesDir := filepath.Join(base, "runtimes")
	if err := os.MkdirAll(filepath.Join(enginesDir, "test-engine"), 0o755); err != nil {
		t.Fatalf("creating engine dir: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(runtimesDir, "test-runtime"), 0o755); err != nil {
		t.Fatalf("creating runtime dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(enginesDir, "test-engine", "engine.yaml"), []byte("name: test-engine\nruntime: test-runtime\n"), 0o644); err != nil {
		t.Fatalf("writing engine manifest: %v", err)
	}
	if err := os.WriteFile(filepath.Join(runtimesDir, "test-runtime", "runtime.yaml"), []byte(runtimeYAML), 0o644); err != nil {
		t.Fatalf("writing runtime manifest: %v", err)
	}

	cfg := storage.NewMockConfig()
	if err := cfg.Set("http.host", "127.0.0.1", storage.UserConfig); err != nil {
		t.Fatalf("setting http.host: %v", err)
	}
	if err := cfg.Set("http.port", "8080", storage.UserConfig); err != nil {
		t.Fatalf("setting http.port: %v", err)
	}

	cache := storage.NewMockCache()
	if err := cache.SetActiveEngine("test-engine"); err != nil {
		t.Fatalf("setting active engine: %v", err)
	}

	return &common.Context{
		EnginesDir:  enginesDir,
		RuntimesDir: runtimesDir,
		Cache:       cache,
		Config:      cfg,
	}
}

func TestProcessEnvConfigs(t *testing.T) {
	mockConfig := storage.NewMockConfig()
	mockConfig.Set("env.my-key", "value", storage.UserConfig)
	mockConfig.Set("env.other", "123", storage.UserConfig)
	mockConfig.Set("other.ignored", "ignored", storage.UserConfig)
	cmd := runCommand{
		Context: &common.Context{
			Config: mockConfig,
		},
	}

	err := cmd.processEnvConfigs()
	if err != nil {
		t.Fatalf("processEnvConfigs returned error: %v", err)
	}

	if got := os.Getenv("MY_KEY"); got != "value" {
		t.Fatalf("expected MY_KEY to be %q, got %q", "value", got)
	}

	if got := os.Getenv("OTHER"); got != "123" {
		t.Fatalf("expected OTHER to be %q, got %q", "123", got)
	}
}

func TestWriteShareProviderEnv(t *testing.T) {
	t.Run("default path", func(t *testing.T) {
		t.Setenv("SNAP_COMMON", t.TempDir())
		t.Setenv("SNAP_NAME", "gemma3-jane")
		t.Setenv("SNAP_INSTANCE_NAME", "gemma3-jane")

		cmd := runCommand{Context: testRunContext(t, "servers:\n  openai:\n    protocol: http\n    base-path: /v1\n"), shareProvider: defaultProviderFilePath}
		if err := cmd.writeShareProviderEnv(); err != nil {
			t.Fatalf("writeShareProviderEnv() error = %v", err)
		}

		path := filepath.Join(os.Getenv("SNAP_COMMON"), "share", "provider", "provider.env")
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("reading provider env file: %v", err)
		}

		want := "SNAP_NAME=gemma3-jane\nSNAP_INSTANCE_NAME=gemma3-jane\nOPENAI_BASE_URL=http://127.0.0.1:8080/v1\n"
		if string(content) != want {
			t.Fatalf("provider env contents mismatch\nwant: %q\ngot:  %q", want, string(content))
		}
	})

	t.Run("custom path", func(t *testing.T) {
		t.Setenv("SNAP_NAME", "gemma3-jane")
		t.Setenv("SNAP_INSTANCE_NAME", "gemma3-jane")

		path := filepath.Join(t.TempDir(), "custom", "provider.env")
		cmd := runCommand{Context: testRunContext(t, "servers:\n  openai:\n    protocol: http\n    base-path: /v1\n"), shareProvider: path}
		if err := cmd.writeShareProviderEnv(); err != nil {
			t.Fatalf("writeShareProviderEnv() error = %v", err)
		}

		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("reading provider env file: %v", err)
		}

		want := "SNAP_NAME=gemma3-jane\nSNAP_INSTANCE_NAME=gemma3-jane\nOPENAI_BASE_URL=http://127.0.0.1:8080/v1\n"
		if string(content) != want {
			t.Fatalf("provider env contents mismatch\nwant: %q\ngot:  %q", want, string(content))
		}
	})

	t.Run("runtime without openai entry", func(t *testing.T) {
		t.Setenv("SNAP_NAME", "gemma3-jane")
		t.Setenv("SNAP_INSTANCE_NAME", "gemma3-jane")

		path := filepath.Join(t.TempDir(), "provider.env")
		cmd := runCommand{Context: testRunContext(t, "servers:\n  kserve:\n    protocol: http\n    base-path: /v2\n"), shareProvider: path}
		if err := cmd.writeShareProviderEnv(); err != nil {
			t.Fatalf("writeShareProviderEnv() error = %v", err)
		}

		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("reading provider env file: %v", err)
		}

		want := "SNAP_NAME=gemma3-jane\nSNAP_INSTANCE_NAME=gemma3-jane\n"
		if string(content) != want {
			t.Fatalf("provider env contents mismatch\nwant: %q\ngot:  %q", want, string(content))
		}
	})

}
