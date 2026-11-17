package commands

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/canonical/inference-snaps/modelctl/v2/commands/common"
	"github.com/canonical/inference-snaps/modelctl/v2/internal/snap"
	"github.com/canonical/inference-snaps/modelctl/v2/internal/storage"
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
		Snap:        snap.Mock(),
	}
}

func activateTestModel(t *testing.T, ctx *common.Context) {
	t.Helper()

	modelsDir := t.TempDir()
	modelDir := filepath.Join(modelsDir, "test-model")
	if err := os.Mkdir(modelDir, 0o755); err != nil {
		t.Fatalf("creating model dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(modelDir, "model.yaml"), []byte("name: test-model\n"), 0o644); err != nil {
		t.Fatalf("writing model manifest: %v", err)
	}
	ctx.ModelsDir = modelsDir
	if err := ctx.Cache.SetActiveModel("test-model"); err != nil {
		t.Fatalf("setting active model: %v", err)
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
		expectedProviderDir := os.ExpandEnv("$SNAP_COMMON/share/provider")

		cmd := runCommand{Context: testRunContext(t, "name: test-runtime\nservers:\n  openai:\n    protocol: http\n    base-path: /v1\n")}
		cmd.shareProvider = cmd.defaultProviderDirectoryPath()
		if err := cmd.writeShareProviderEnv(); err != nil {
			t.Fatalf("writeShareProviderEnv() error = %v", err)
		}

		cachedShareProviderDirectory, err := cmd.Cache.GetSharedProviderDirectory()
		if err != nil {
			t.Fatalf("getting cached shared provider directory: %v", err)
		}
		if cachedShareProviderDirectory != expectedProviderDir {
			t.Fatalf("cached shared provider directory mismatch\nwant: %q\ngot:  %q", expectedProviderDir, cachedShareProviderDirectory)
		}

		path := filepath.Join(expectedProviderDir, "provider.env")
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("reading provider env file: %v", err)
		}

		want := "SNAP_NAME=gemma3-jane\nSNAP_INSTANCE_NAME=gemma3-jane\nOPENAI_BASE_URL=http://localhost:8080/v1\n"
		if string(content) != want {
			t.Fatalf("provider env contents mismatch\nwant: %q\ngot:  %q", want, string(content))
		}
	})

	t.Run("custom path", func(t *testing.T) {
		t.Setenv("SNAP_NAME", "gemma3-jane")
		t.Setenv("SNAP_INSTANCE_NAME", "gemma3-jane")

		path := filepath.Join(t.TempDir(), "custom")
		cmd := runCommand{Context: testRunContext(t, "name: test-runtime\nservers:\n  openai:\n    protocol: http\n    base-path: /v1\n"), shareProvider: path}
		if err := cmd.Config.Set("http.host", "0.0.0.0", storage.UserConfig); err != nil {
			t.Fatalf("setting HTTP host: %v", err)
		}
		if err := cmd.writeShareProviderEnv(); err != nil {
			t.Fatalf("writeShareProviderEnv() error = %v", err)
		}

		cachedShareProviderDirectory, err := cmd.Cache.GetSharedProviderDirectory()
		if err != nil {
			t.Fatalf("getting cached shared provider directory: %v", err)
		}
		if cachedShareProviderDirectory != path {
			t.Fatalf("cached shared provider directory mismatch\nwant: %q\ngot:  %q", path, cachedShareProviderDirectory)
		}

		content, err := os.ReadFile(filepath.Join(path, "provider.env"))
		if err != nil {
			t.Fatalf("reading provider env file: %v", err)
		}

		want := "SNAP_NAME=gemma3-jane\nSNAP_INSTANCE_NAME=gemma3-jane\nOPENAI_BASE_URL=http://localhost:8080/v1\n"
		if string(content) != want {
			t.Fatalf("provider env contents mismatch\nwant: %q\ngot:  %q", want, string(content))
		}
	})

	t.Run("runtime without openai entry", func(t *testing.T) {
		t.Setenv("SNAP_NAME", "gemma3-jane")
		t.Setenv("SNAP_INSTANCE_NAME", "gemma3-jane")

		path := t.TempDir()
		cmd := runCommand{Context: testRunContext(t, "name: test-runtime\nservers:\n  kserve:\n    protocol: http\n    base-path: /v2\n"), shareProvider: path}
		if err := cmd.writeShareProviderEnv(); err != nil {
			t.Fatalf("writeShareProviderEnv() error = %v", err)
		}

		cachedShareProviderDirectory, err := cmd.Cache.GetSharedProviderDirectory()
		if err != nil {
			t.Fatalf("getting cached shared provider directory: %v", err)
		}
		if cachedShareProviderDirectory != path {
			t.Fatalf("cached shared provider directory mismatch\nwant: %q\ngot:  %q", path, cachedShareProviderDirectory)
		}

		content, err := os.ReadFile(filepath.Join(path, "provider.env"))
		if err != nil {
			t.Fatalf("reading provider env file: %v", err)
		}

		want := "SNAP_NAME=gemma3-jane\nSNAP_INSTANCE_NAME=gemma3-jane\n"
		if string(content) != want {
			t.Fatalf("provider env contents mismatch\nwant: %q\ngot:  %q", want, string(content))
		}
	})

	t.Run("openai server over unix socket", func(t *testing.T) {
		t.Setenv("SNAP_NAME", "gemma3-jane")
		t.Setenv("SNAP_INSTANCE_NAME", "gemma3-jane")

		// OpenAiBaseUrl rejects an "openai" entrypoint that has no Url (e.g. a
		// Unix socket entrypoint), so writeShareProviderEnv must not let that
		// error abort the function before the UNIX_SOCKET entries are written.
		path := t.TempDir()
		cmd := runCommand{Context: testRunContext(t, "name: test-runtime\nservers:\n  openai:\n    protocol: http+unix\n    base-path: /v1\n"), shareProvider: path}
		if err := cmd.writeShareProviderEnv(); err != nil {
			t.Fatalf("writeShareProviderEnv() error = %v", err)
		}

		cachedShareProviderDirectory, err := cmd.Cache.GetSharedProviderDirectory()
		if err != nil {
			t.Fatalf("getting cached shared provider directory: %v", err)
		}
		if cachedShareProviderDirectory != path {
			t.Fatalf("cached shared provider directory mismatch\nwant: %q\ngot:  %q", path, cachedShareProviderDirectory)
		}

		content, err := os.ReadFile(filepath.Join(path, "provider.env"))
		if err != nil {
			t.Fatalf("reading provider env file: %v", err)
		}

		want := "SNAP_NAME=gemma3-jane\nSNAP_INSTANCE_NAME=gemma3-jane\nUNIX_SOCKET=server.sock\n"
		if string(content) != want {
			t.Fatalf("provider env contents mismatch\nwant: %q\ngot:  %q", want, string(content))
		}
	})

	tests := []struct {
		name       string
		serverName string
		protocol   string
		namespace  string
		wantSocket string
	}{
		{
			name:       "http unix socket",
			serverName: "test",
			protocol:   "http+unix",
			wantSocket: "UNIX_SOCKET=server.sock",
		},
		{
			name:       "https unix socket",
			serverName: "example",
			protocol:   "https+unix",
			wantSocket: "UNIX_SOCKET=server.sock",
		},
		{
			name:       "websocket unix socket",
			serverName: "server",
			protocol:   "ws+unix",
			wantSocket: "UNIX_SOCKET=server.sock",
		},
		{
			name:       "secure websocket unix socket",
			serverName: "server",
			protocol:   "wss+unix",
			wantSocket: "UNIX_SOCKET=server.sock",
		},
		{
			name:       "namespaced unix socket",
			serverName: "server",
			protocol:   "http+unix",
			namespace:  "example",
			wantSocket: "EXAMPLE_UNIX_SOCKET=example.sock",
		},
		{
			name:       "namespaced https unix socket",
			serverName: "example",
			protocol:   "https+unix",
			namespace:  "example",
			wantSocket: "EXAMPLE_UNIX_SOCKET=example.sock",
		},
		{
			name:       "namespaced websocket unix socket",
			serverName: "server",
			protocol:   "ws+unix",
			namespace:  "example",
			wantSocket: "EXAMPLE_UNIX_SOCKET=example.sock",
		},
		{
			name:       "secure websocket unix socket",
			serverName: "server",
			protocol:   "wss+unix",
			namespace:  "example",
			wantSocket: "EXAMPLE_UNIX_SOCKET=example.sock",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("SNAP_NAME", "gemma3-jane")
			t.Setenv("SNAP_INSTANCE_NAME", "gemma3-jane")

			runtimeYAML := fmt.Sprintf("name: test-runtime\nservers:\n  %s:\n    protocol: %s\n    base-path: /v1\n", tt.serverName, tt.protocol)
			if tt.namespace != "" {
				runtimeYAML += fmt.Sprintf("    namespace: %s\n", tt.namespace)
			}

			path := t.TempDir()
			cmd := runCommand{Context: testRunContext(t, runtimeYAML), shareProvider: path}
			if err := cmd.writeShareProviderEnv(); err != nil {
				t.Fatalf("writeShareProviderEnv() error = %v", err)
			}

			content, err := os.ReadFile(filepath.Join(path, "provider.env"))
			if err != nil {
				t.Fatalf("reading provider env file: %v", err)
			}

			want := "SNAP_NAME=gemma3-jane\nSNAP_INSTANCE_NAME=gemma3-jane\n" + tt.wantSocket + "\n"
			if string(content) != want {
				t.Fatalf("provider env contents mismatch\nwant: %q\ngot:  %q", want, string(content))
			}
		})
	}
}

func TestNoActiveModel(t *testing.T) {
	t.Run("no active model", func(t *testing.T) {
		cmd := runCommand{Context: testRunContext(t, "name: test-runtime\nservers:\n  openai:\n    protocol: http\n    base-path: /v1\n")}
		err := cmd.run(nil, []string{"echo", "Hello World!"})
		if err == nil || err.Error() != "no active model" {
			t.Fatalf("expected error 'no active model', got %v", err)
		}
	})
}

func TestRunCommandFailure(t *testing.T) {
	missingCommand := filepath.Join(t.TempDir(), "missing-command")
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{name: "non-zero exit", args: []string{"/bin/sh", "-c", "exit 7"}, wantErr: "exit status 7"},
		{name: "executable not found", args: []string{missingCommand}, wantErr: "fork/exec " + missingCommand + ": no such file or directory"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := testRunContext(t, "name: test-runtime\nservers:\n  openai:\n    protocol: http\n    base-path: /v1\n")
			activateTestModel(t, ctx)

			err := (&runCommand{Context: ctx}).run(nil, tt.args)
			if err == nil || err.Error() != tt.wantErr {
				t.Fatalf("expected error %q, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestRunDeprecatedFallbackFlag(t *testing.T) {
	missingCommand := filepath.Join(t.TempDir(), "missing-command")
	tests := []struct {
		name       string
		args       []string
		noModel    bool
		wantErr    string
		wantSignal syscall.Signal
	}{
		{name: "success", args: []string{"/bin/sh", "-c", "exit 0"}},
		{name: "failure", args: []string{"/bin/sh", "-c", "exit 7"}, wantErr: "exit status 7"},
		{name: "missing executable", args: []string{missingCommand}, wantErr: "fork/exec " + missingCommand + ": no such file or directory"},
		{name: "no model", args: []string{"/bin/true"}, noModel: true, wantErr: "no active model"},
		{name: "SIGTERM", args: []string{"/bin/sh", "-c", "kill -TERM $$"}, wantSignal: syscall.SIGTERM},
		{name: "SIGINT", args: []string{"/bin/sh", "-c", "kill -INT $$"}, wantSignal: syscall.SIGINT},
	}

	for _, flag := range []string{"--fallback-server", "--fallback-server=true", "--fallback-server=false"} {
		t.Run(flag, func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					ctx := testRunContext(t, diagnoseRuntime)
					if !tt.noModel {
						activateTestModel(t, ctx)
					}
					cmd := Run(ctx)
					if cmd.Flags().Lookup("fallback-server").Deprecated == "" {
						t.Fatal("fallback-server should be deprecated")
					}
					cmd.SetArgs(append([]string{flag, "--"}, tt.args...))
					output := &bytes.Buffer{}
					cmd.SetOut(output)
					cmd.SetErr(output)
					err := cmd.Execute()
					switch {
					case tt.wantSignal != 0:
						var exitErr *exec.ExitError
						if !errors.As(err, &exitErr) {
							t.Fatalf("expected child signal error, got %v", err)
						}
						status, ok := exitErr.ProcessState.Sys().(syscall.WaitStatus)
						if !ok || !status.Signaled() || status.Signal() != tt.wantSignal {
							t.Fatalf("expected signal %v, got %v", tt.wantSignal, err)
						}
					case tt.wantErr != "":
						if err == nil || err.Error() != tt.wantErr {
							t.Fatalf("error = %v, want %q", err, tt.wantErr)
						}
					case err != nil:
						t.Fatal(err)
					}
					if !strings.Contains(output.String(), "deprecated") || !strings.Contains(output.String(), "has no effect") {
						t.Fatalf("expected deprecation warning, got %q", output.String())
					}
				})
			}
		})
	}
}
