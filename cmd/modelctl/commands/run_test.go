package commands

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"testing"
	"time"

	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/common"
	"github.com/canonical/inference-snaps-cli/v2/pkg/snap"
	"github.com/canonical/inference-snaps-cli/v2/pkg/storage"
)

const runFallbackHelperEnv = "GO_WANT_RUN_FALLBACK_HELPER"

func TestListenAddress(t *testing.T) {
	cmd := runCommand{}
	tests := []struct {
		name    string
		baseURL string
		want    string
		wantErr bool
	}{
		{name: "host port and path", baseURL: "http://127.0.0.1:8080/v1", want: "127.0.0.1:8080"},
		{name: "default port", baseURL: "http://localhost/v3", want: "localhost:80"},
		{name: "IPv6", baseURL: "http://[::1]:9000/api/v1", want: "[::1]:9000"},
		{name: "IPv6 default port", baseURL: "http://[::1]/v1", want: "[::1]:80"},
		{name: "no path", baseURL: "http://localhost:8080", want: "localhost:8080"},
		{name: "invalid URL", baseURL: "://localhost:8080/v1", wantErr: true},
		{name: "missing host", baseURL: "http:///v1", wantErr: true},
		{name: "invalid port", baseURL: "http://localhost:invalid/v1", wantErr: true},
		{name: "HTTPS", baseURL: "https://localhost:8080/v1", want: "localhost:8080"},
		{name: "user information", baseURL: "http://user@localhost:8080/v1", want: "localhost:8080"},
		{name: "query", baseURL: "http://localhost:8080/v1?key=value", want: "localhost:8080"},
		{name: "fragment", baseURL: "http://localhost:8080/v1#fragment", want: "localhost:8080"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.listenAddress(tt.baseURL)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("listenAddress(%q) unexpectedly succeeded with %q", tt.baseURL, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("listenAddress(%q) returned error: %v", tt.baseURL, err)
			}
			if got != tt.want {
				t.Fatalf("listenAddress(%q) = %q, want %q", tt.baseURL, got, tt.want)
			}
		})
	}
}

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

		want := "SNAP_NAME=gemma3-jane\nSNAP_INSTANCE_NAME=gemma3-jane\nOPENAI_BASE_URL=http://127.0.0.1:8080/v1\n"
		if string(content) != want {
			t.Fatalf("provider env contents mismatch\nwant: %q\ngot:  %q", want, string(content))
		}
	})

	t.Run("custom path", func(t *testing.T) {
		t.Setenv("SNAP_NAME", "gemma3-jane")
		t.Setenv("SNAP_INSTANCE_NAME", "gemma3-jane")

		path := filepath.Join(t.TempDir(), "custom")
		cmd := runCommand{Context: testRunContext(t, "name: test-runtime\nservers:\n  openai:\n    protocol: http\n    base-path: /v1\n"), shareProvider: path}
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

		want := "SNAP_NAME=gemma3-jane\nSNAP_INSTANCE_NAME=gemma3-jane\nOPENAI_BASE_URL=http://127.0.0.1:8080/v1\n"
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

func TestRunCommandFallbackServer(t *testing.T) {
	port := availableTCPPort(t)
	cmd, done, output := startRunFallbackHelper(t, "failure", port)
	t.Cleanup(func() {
		if cmd.ProcessState == nil {
			_ = cmd.Process.Kill()
			<-done
		}
	})

	url := fmt.Sprintf("http://127.0.0.1:%d/v1/chat/completions", port)
	client := &http.Client{Timeout: 100 * time.Millisecond}
	deadline := time.Now().Add(5 * time.Second)
	for {
		response, err := client.Get(url)
		if err == nil {
			_ = response.Body.Close()
			if response.StatusCode != http.StatusServiceUnavailable {
				t.Fatalf("fallback response status = %d, want %d", response.StatusCode, http.StatusServiceUnavailable)
			}
			return
		}

		select {
		case err := <-done:
			t.Fatalf("fallback helper exited before serving: %v\n%s", err, output.String())
		default:
		}
		if time.Now().After(deadline) {
			t.Fatalf("fallback server did not start at %s: %v", url, err)
		}
		time.Sleep(25 * time.Millisecond)
	}
}

func TestRunCommandStoppedChildDoesNotStartFallback(t *testing.T) {
	tests := []struct {
		name string
		mode string
	}{
		{name: "SIGTERM", mode: "sigterm"},
		{name: "SIGINT", mode: "sigint"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			port := availableTCPPort(t)
			cmd, done, output := startRunFallbackHelper(t, tt.mode, port)

			select {
			case err := <-done:
				if err != nil {
					t.Fatalf("fallback helper failed: %v\n%s", err, output.String())
				}
			case <-time.After(5 * time.Second):
				_ = cmd.Process.Kill()
				<-done
				t.Fatalf("run command did not return after child received %s", tt.name)
			}

			listener, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(port)))
			if err != nil {
				t.Fatalf("fallback server unexpectedly occupied configured endpoint: %v", err)
			}
			_ = listener.Close()
		})
	}
}

func TestRunCommandFallbackHelper(t *testing.T) {
	mode := os.Getenv(runFallbackHelperEnv)
	if mode == "" {
		return
	}

	port := os.Getenv("RUN_FALLBACK_HELPER_PORT")
	ctx := testRunContext(t, "name: test-runtime\nservers:\n  openai:\n    protocol: http\n    base-path: /v1\n")
	if err := ctx.Config.Set("http.port", port, storage.UserConfig); err != nil {
		t.Fatalf("setting fallback server port: %v", err)
	}
	activateTestModel(t, ctx)

	childCommand := "exit 7"
	switch mode {
	case "sigterm":
		childCommand = "kill -TERM $$"
	case "sigint":
		childCommand = "kill -INT $$"
	case "failure":
	default:
		t.Fatalf("unknown helper mode %q", mode)
	}

	err := (&runCommand{Context: ctx, fallbackServer: true}).run(nil, []string{"/bin/sh", "-c", childCommand})
	if mode == "sigterm" || mode == "sigint" {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			t.Fatalf("expected child signal error, got %v", err)
		}
		status, ok := exitErr.ProcessState.Sys().(syscall.WaitStatus)
		wantSignal := syscall.SIGTERM
		if mode == "sigint" {
			wantSignal = syscall.SIGINT
		}
		if !ok || !status.Signaled() || status.Signal() != wantSignal {
			t.Fatalf("expected child signal %v, got %v", wantSignal, err)
		}
		return
	}
	if err != nil {
		t.Fatalf("run command failed: %v", err)
	}
}

func availableTCPPort(t *testing.T) int {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserving TCP port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	if err := listener.Close(); err != nil {
		t.Fatalf("releasing TCP port: %v", err)
	}
	return port
}

func startRunFallbackHelper(t *testing.T, mode string, port int) (*exec.Cmd, <-chan error, *bytes.Buffer) {
	t.Helper()

	cmd := exec.Command(os.Args[0], "-test.run=^TestRunCommandFallbackHelper$")
	cmd.Env = append(os.Environ(),
		runFallbackHelperEnv+"="+mode,
		"RUN_FALLBACK_HELPER_PORT="+strconv.Itoa(port),
	)
	output := &bytes.Buffer{}
	cmd.Stdout = output
	cmd.Stderr = output
	if err := cmd.Start(); err != nil {
		t.Fatalf("starting fallback helper: %v", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	return cmd, done, output
}

func TestCommandStopped(t *testing.T) {
	cmd := runCommand{}
	tests := []struct {
		name    string
		command string
		want    bool
	}{
		{name: "SIGTERM", command: "kill -TERM $$", want: true},
		{name: "SIGINT", command: "kill -INT $$", want: true},
		{name: "SIGKILL", command: "kill -KILL $$", want: false},
		{name: "non-zero exit", command: "exit 7", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := exec.Command("/bin/sh", "-c", tt.command).Run()
			if err == nil {
				t.Fatal("expected command to fail")
			}
			if got := cmd.commandStopped(err); got != tt.want {
				var status syscall.WaitStatus
				var exitErr *exec.ExitError
				if errors.As(err, &exitErr) {
					status, _ = exitErr.ProcessState.Sys().(syscall.WaitStatus)
				}
				t.Fatalf("commandStopped() = %v, want %v (status %v)", got, tt.want, status)
			}
		})
	}
}
