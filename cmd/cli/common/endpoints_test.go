package common

import (
	"strings"
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/storage"
)

func TestServerEndpoints(t *testing.T) {
	testCases := []struct {
		name             string
		componentConfigs []ComponentSettings
		want             map[string]string
		wantErrContains  string
	}{
		{
			name: "multiple components and servers",
			componentConfigs: []ComponentSettings{
				{
					Servers: map[string]map[string]string{
						"openai": {
							"protocol":  "http",
							"base-path": "/v1",
						},
					},
				},
				{
					Servers: map[string]map[string]string{
						"kserve": {
							"protocol":  "https",
							"base-path": "/v2",
						},
						"webui": {
							"protocol": "http",
						},
					},
				},
			},
			want: map[string]string{
				"openai": "http://127.0.0.1:8080/v1",
				"kserve": "https://127.0.0.1:8080/v2",
				"webui":  "http://192.0.2.1:8080/",
			},
		},
		{
			name: "unsupported protocol",
			componentConfigs: []ComponentSettings{
				{
					Servers: map[string]map[string]string{
						"openai": {
							"protocol": "ftp",
						},
					},
				},
			},
			wantErrContains: "unsupported protocol",
		},
		{
			name: "OPENAI_BASE_PATH deprecated",
			componentConfigs: []ComponentSettings{
				{
					Environment: []string{"OPENAI_BASE_PATH=/v1"},
				},
			},
			wantErrContains: "deprecated",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Setenv("ADDITIONAL_FEATURES", "webui")

			config := storage.NewMockConfig()
			config.Set("http.port", "8080", storage.UserConfig)
			config.Set("http.host", "127.0.0.1", storage.UserConfig)
			config.Set("webui.http.port", "8080", storage.UserConfig)
			config.Set("webui.http.host", "192.0.2.1", storage.UserConfig)
			ctx := &Context{
				Config: config,
			}

			got, err := serverEndpoints(ctx, testCase.componentConfigs)
			if testCase.wantErrContains != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", testCase.wantErrContains)
				}
				if !strings.Contains(err.Error(), testCase.wantErrContains) {
					t.Fatalf("got error %q, want it to contain %q", err.Error(), testCase.wantErrContains)
				}
				return
			}

			if err != nil {
				t.Fatal(err)
			}

			if len(got) != len(testCase.want) {
				t.Fatalf("got %d endpoints, want %d", len(got), len(testCase.want))
			}

			for endpointName, wantURL := range testCase.want {
				gotURL, found := got[endpointName]
				if !found {
					t.Fatalf("missing endpoint %q", endpointName)
				}
				if gotURL != wantURL {
					t.Fatalf("got %q: %q, want %q", endpointName, gotURL, wantURL)
				}
			}
		})
	}
}

func TestServerHttpUrl(t *testing.T) {
	testCases := []struct {
		name         string
		serverConfig map[string]string
		host         string
		setHost      bool
		want         string
	}{
		{
			name: "default base path",
			serverConfig: map[string]string{
				"protocol": "http",
			},
			host:    "0.0.0.0",
			setHost: true,
			want:    "http://0.0.0.0:8080/",
		},
		{
			name: "custom base path",
			serverConfig: map[string]string{
				"protocol":  "http",
				"base-path": "/v1",
			},
			host:    "127.0.0.1",
			setHost: true,
			want:    "http://127.0.0.1:8080/v1",
		},
		{
			name: "https protocol",
			serverConfig: map[string]string{
				"protocol":  "https",
				"base-path": "/v3",
			},
			host:    "0.0.0.0",
			setHost: true,
			want:    "https://0.0.0.0:8080/v3",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			config := storage.NewMockConfig()
			config.Set("http.port", "8080", storage.UserConfig)
			if testCase.setHost {
				config.Set("http.host", testCase.host, storage.UserConfig)
			}
			ctx := &Context{
				Config: config,
			}

			got, err := serverHttpUrl(ctx, testCase.serverConfig)
			if err != nil {
				t.Fatal(err)
			}

			if got != testCase.want {
				t.Fatalf("got %q, want %q", got, testCase.want)
			}
		})
	}
}
