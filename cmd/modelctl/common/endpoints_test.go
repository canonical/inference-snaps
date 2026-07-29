package common

import (
	"strings"
	"testing"

	"github.com/canonical/inference-snaps-cli/v2/pkg/runtimes"
	"github.com/canonical/inference-snaps-cli/v2/pkg/storage"
)

func TestServerEntrypoints(t *testing.T) {
	testCases := []struct {
		name               string
		engineYAML         string
		runtimeYAML        string
		wantEntrypointURLs map[string]string
		wantUnixSockets    map[string]string
		wantUnixSocketUrls map[string]string
		wantErrContains    string
	}{
		{
			name: "multiple servers",
			engineYAML: `
name: test-engine
runtime: test-runtime
`,
			runtimeYAML: `
servers:
  openai:
    protocol: http
    base-path: /v1
  kserve:
    protocol: https
    base-path: /v2
`,
			wantEntrypointURLs: map[string]string{
				"openai": "http://127.0.0.1:8080/v1",
				"kserve": "https://127.0.0.1:8080/v2",
				"webui":  "http://192.0.2.1:8080",
			},
		},
		{
			name: "http unix servers",
			engineYAML: `
name: test-engine
runtime: test-runtime
`,
			runtimeYAML: `
servers:
  openai-unix:
    protocol: http+unix
    base-path: /v1
  kserve-unix:
    protocol: https+unix
    base-path: /v2
    namespace: kserve
`,
			wantEntrypointURLs: map[string]string{
				"openai-unix": "",
				"kserve-unix": "",
				"webui":       "http://192.0.2.1:8080",
			},
			wantUnixSockets: map[string]string{
				"openai-unix": "/run/openai.sock",
				"kserve-unix": "/run/kserve.sock",
			},
			wantUnixSocketUrls: map[string]string{
				"openai-unix": "http://unix/v1",
				"kserve-unix": "https://unix/v2",
			},
		},
		{
			name: "unsupported protocol",
			engineYAML: `
name: test-engine
runtime: test-runtime
`,
			runtimeYAML: `
servers:
  openai:
    protocol: ftp
`,
			wantErrContains: "unsupported protocol",
		},
		{
			name: "websocket servers",
			engineYAML: `
name: test-engine
runtime: test-runtime
`,
			runtimeYAML: `
servers:
  openai-ws:
    protocol: ws
    base-path: /v1
  openai-unix:
    protocol: ws+unix
    base-path: /v1
`,
			wantEntrypointURLs: map[string]string{
				"openai-ws":   "ws://127.0.0.1:8081/v1",
				"openai-unix": "",
			},
			wantUnixSockets: map[string]string{
				"openai-unix": "/run/openai.sock",
			},
			wantUnixSocketUrls: map[string]string{
				"openai-unix": "ws://unix/v1",
			},
		},
		{
			name: "websocket namespaced server",
			engineYAML: `
name: test-engine
runtime: test-runtime
`,
			runtimeYAML: `
servers:
  logger:
    protocol: ws
    base-path: /stream
    namespace: logger
`,
			wantEntrypointURLs: map[string]string{
				"logger": "ws://localhost:9091/stream",
				"webui":  "http://192.0.2.1:8080",
			},
		},
		{
			name: "secure websocket servers",
			engineYAML: `
name: test-engine
runtime: test-runtime
`,
			runtimeYAML: `
servers:
  openai-wss:
    protocol: wss
    base-path: /v1
  openai-unix:
    protocol: wss+unix
    base-path: /v1
`,
			wantEntrypointURLs: map[string]string{
				"openai-wss":  "wss://127.0.0.1:8081/v1",
				"openai-unix": "",
			},
			wantUnixSockets: map[string]string{
				"openai-unix": "/run/openai.sock",
			},
			wantUnixSocketUrls: map[string]string{
				"openai-unix": "wss://unix/v1",
			},
		},
		{
			name: "secure websocket namespaced server",
			engineYAML: `
name: test-engine
runtime: test-runtime
`,
			runtimeYAML: `
servers:
  logger:
    protocol: wss
    base-path: /stream
    namespace: logger
`,
			wantEntrypointURLs: map[string]string{
				"logger": "wss://localhost:9091/stream",
				"webui":  "http://192.0.2.1:8080",
			},
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

			configs := map[string]string{
				"http.port":        "8080",
				"http.host":        "127.0.0.1",
				"http.unix-socket": "/run/openai.sock",
				"ws.port":          "8081",
				"ws.host":          "127.0.0.1",
				"ws.unix-socket":   "/run/openai.sock",
				// namespaced configurations
				"webui.http.port":         "8080",
				"webui.http.host":         "192.0.2.1",
				"logger.ws.port":          "9091",
				"logger.ws.host":          "localhost",
				"kserve.http.unix-socket": "/run/kserve.sock",
			}

			for key, value := range configs {
				if err := config.Set(key, value, storage.UserConfig); err != nil {
					t.Fatalf("Set(%q): %v", key, err)
				}
			}

			ctx := &Context{
				EnginesDir:  enginesDir,
				RuntimesDir: runtimesDir,
				Config:      config,
				Cache:       cache,
			}

			got, err := ServerEntrypoints(ctx)
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
			for name, wantURL := range tc.wantEntrypointURLs {
				entrypoint, found := got[name]
				if !found {
					t.Fatalf("missing entrypoint %q", name)
				}
				if entrypoint.Url != wantURL {
					t.Fatalf("entrypoint %q: got %q, want %q", name, entrypoint.Url, wantURL)
				}
			}
			for name, wantSocket := range tc.wantUnixSockets {
				entrypoint, found := got[name]
				if !found {
					t.Fatalf("missing entrypoint %q", name)
				}
				if entrypoint.UnixSocket != wantSocket {
					t.Fatalf("entrypoint %q unix socket: got %q, want %q", name, entrypoint.UnixSocket, wantSocket)
				}
			}
			for name, wantSocketURL := range tc.wantUnixSocketUrls {
				entrypoint, found := got[name]
				if !found {
					t.Fatalf("missing entrypoint %q", name)
				}
				if entrypoint.UnixSocketUrl != wantSocketURL {
					t.Fatalf("entrypoint %q unix socket URL: got %q, want %q", name, entrypoint.UnixSocketUrl, wantSocketURL)
				}
			}
		})
	}
}

func TestServerHTTPEntrypoint(t *testing.T) {
	testCases := []struct {
		name    string
		server  runtimes.Server
		host    string
		want    string
		wantErr bool
	}{
		{
			name: "default base path",
			server: runtimes.Server{
				Protocol: "http",
			},
			host: "0.0.0.0",
			want: "http://0.0.0.0:8080",
		},
		{
			name: "custom base path",
			server: runtimes.Server{
				Protocol: "http",
				BasePath: "/v1",
			},
			host: "127.0.0.1",
			want: "http://127.0.0.1:8080/v1",
		},
		{
			name: "https protocol",
			server: runtimes.Server{
				Protocol: "https",
				BasePath: "/v3",
			},
			host: "0.0.0.0",
			want: "https://0.0.0.0:8080/v3",
		},
		{
			name: "missing host",
			server: runtimes.Server{
				Protocol: "http",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := storage.NewMockConfig()
			configs := map[string]string{
				"http.port": "8080",
				"http.host": tc.host,
			}
			for key, value := range configs {
				if err := config.Set(key, value, storage.UserConfig); err != nil {
					t.Fatalf("Set(%q): %v", key, err)
				}
			}
			ctx := &Context{
				Config: config,
			}

			got, err := serverHttpEntrypoint(ctx, tc.server)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got == nil {
				t.Fatal("expected non-nil entrypoint")
			}
			if got.Url != tc.want {
				t.Fatalf("got %q, want %q", got.Url, tc.want)
			}
		})
	}
}

func TestServerHttpUnixSocketEntrypoint(t *testing.T) {
	testCases := []struct {
		name          string
		server        runtimes.Server
		socketPath    string
		wantSocketURL string
	}{
		{
			name: "http unix default namespace",
			server: runtimes.Server{
				Protocol: "http+unix",
				BasePath: "/v1",
			},
			socketPath:    "/run/test.sock",
			wantSocketURL: "http://unix/v1",
		},
		{
			name: "http unix with namespace",
			server: runtimes.Server{
				Protocol:  "http+unix",
				BasePath:  "/api/v2",
				Namespace: "proxy",
			},
			socketPath:    "/run/proxy.sock",
			wantSocketURL: "http://unix/api/v2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := storage.NewMockConfig()
			configs := map[string]string{}
			if tc.server.Namespace != "" {
				configs[tc.server.Namespace+".http.unix-socket"] = tc.socketPath
			} else {
				configs["http.unix-socket"] = tc.socketPath
			}
			for key, value := range configs {
				if err := config.Set(key, value, storage.UserConfig); err != nil {
					t.Fatalf("Set(%q): %v", key, err)
				}
			}
			ctx := &Context{
				Config: config,
			}

			got, err := serverHttpOverUnixSocketEntrypoint(ctx, tc.server)
			if err != nil {
				t.Fatal(err)
			}
			if got == nil {
				t.Fatal("expected non-nil entrypoint")
			}
			if got.Url != "" {
				t.Fatalf("URL: got %q, want empty", got.Url)
			}
			if got.UnixSocket != tc.socketPath {
				t.Fatalf("UnixSocket: got %q, want %q", got.UnixSocket, tc.socketPath)
			}
			if got.UnixSocketUrl != tc.wantSocketURL {
				t.Fatalf("UnixSocketUrl: got %q, want %q", got.UnixSocketUrl, tc.wantSocketURL)
			}
		})
	}
}

func TestServerWSEntrypoint(t *testing.T) {
	testCases := []struct {
		name    string
		server  runtimes.Server
		wsHost  string
		wsPort  string
		wantURL string
		wantErr bool
	}{
		{
			name: "empty base path",
			server: runtimes.Server{
				Protocol: "ws",
			},
			wsHost:  "0.0.0.0",
			wsPort:  "8080",
			wantURL: "ws://0.0.0.0:8080",
		},
		{
			name: "custom base path",
			server: runtimes.Server{
				Protocol: "ws",
				BasePath: "/v1",
			},
			wsHost:  "127.0.0.1",
			wsPort:  "8081",
			wantURL: "ws://127.0.0.1:8081/v1",
		},
		{
			name: "complex base path",
			server: runtimes.Server{
				Protocol: "ws",
				BasePath: "/api/v2/stream",
			},
			wsHost:  "localhost",
			wsPort:  "9090",
			wantURL: "ws://localhost:9090/api/v2/stream",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := storage.NewMockConfig()
			configs := map[string]string{
				"ws.port": tc.wsPort,
				"ws.host": tc.wsHost,
			}
			for key, value := range configs {
				if err := config.Set(key, value, storage.UserConfig); err != nil {
					t.Fatalf("Set(%q): %v", key, err)
				}
			}
			ctx := &Context{Config: config}

			got, err := serverWsEntrypoint(ctx, tc.server)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got == nil {
				t.Fatal("expected non-nil entrypoint")
			}
			if got.Url != tc.wantURL {
				t.Fatalf("URL: got %q, want %q", got.Url, tc.wantURL)
			}
		})
	}
}

func TestServerWsUnixSocketEntrypoint(t *testing.T) {
	testCases := []struct {
		name          string
		server        runtimes.Server
		socketPath    string
		wantSocketURL string
	}{
		{
			name: "ws unix default namespace",
			server: runtimes.Server{
				Protocol: "ws+unix",
				BasePath: "/v1",
			},
			socketPath:    "/run/ws.sock",
			wantSocketURL: "ws://unix/v1",
		},
		{
			name: "ws unix with namespace",
			server: runtimes.Server{
				Protocol:  "ws+unix",
				BasePath:  "/api/stream",
				Namespace: "proxy",
			},
			socketPath:    "/run/proxy-ws.sock",
			wantSocketURL: "ws://unix/api/stream",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := storage.NewMockConfig()
			configs := map[string]string{}
			if tc.server.Namespace != "" {
				configs[tc.server.Namespace+".ws.unix-socket"] = tc.socketPath
			} else {
				configs["ws.unix-socket"] = tc.socketPath
			}
			for key, value := range configs {
				if err := config.Set(key, value, storage.UserConfig); err != nil {
					t.Fatalf("Set(%q): %v", key, err)
				}
			}
			ctx := &Context{Config: config}

			got, err := serverWsOverUnixSocketEntrypoint(ctx, tc.server)
			if err != nil {
				t.Fatal(err)
			}
			if got == nil {
				t.Fatal("expected non-nil entrypoint")
			}
			if got.Url != "" {
				t.Fatalf("URL: got %q, want empty", got.Url)
			}
			if got.UnixSocket != tc.socketPath {
				t.Fatalf("UnixSocket: got %q, want %q", got.UnixSocket, tc.socketPath)
			}
			if got.UnixSocketUrl != tc.wantSocketURL {
				t.Fatalf("UnixSocketUrl: got %q, want %q", got.UnixSocketUrl, tc.wantSocketURL)
			}
		})
	}
}
