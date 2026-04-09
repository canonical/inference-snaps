package ui

import (
	"testing"
)

func TestSupportedCapabilities(t *testing.T) {
	caps := SupportedCapabilities()
	if len(caps) == 0 {
		t.Fatal("expected non-empty capabilities list")
	}
	found := map[string]bool{}
	for _, c := range caps {
		found[c] = true
	}
	if !found["text"] {
		t.Error("expected 'text' in supported capabilities")
	}
	if !found["vision"] {
		t.Error("expected 'vision' in supported capabilities")
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config with no capabilities",
			config: Config{
				OpenAIBaseURL: "http://localhost:11434",
				Capabilities:  []string{},
				InstanceName:  "my-instance",
				EngineName:    "llama3",
			},
			wantErr: false,
		},
		{
			name: "valid config with text capability",
			config: Config{
				OpenAIBaseURL: "http://localhost:11434/v1",
				Capabilities:  []string{"text"},
			},
			wantErr: false,
		},
		{
			name: "valid config with vision capability",
			config: Config{
				OpenAIBaseURL: "http://localhost:11434/v1",
				Capabilities:  []string{"vision"},
			},
			wantErr: false,
		},
		{
			name: "valid config with both capabilities",
			config: Config{
				OpenAIBaseURL: "http://localhost:11434/v1",
				Capabilities:  []string{"text", "vision"},
			},
			wantErr: false,
		},
		{
			name: "empty OpenAIBaseURL is valid (url.Parse accepts empty string)",
			config: Config{
				OpenAIBaseURL: "",
				Capabilities:  []string{"text"},
			},
			wantErr: false,
		},
		{
			name: "unsupported capability",
			config: Config{
				OpenAIBaseURL: "http://localhost:11434/v1",
				Capabilities:  []string{"audio"},
			},
			wantErr: true,
		},
		{
			name: "one valid one invalid capability",
			config: Config{
				OpenAIBaseURL: "http://localhost:11434/v1",
				Capabilities:  []string{"text", "unknown"},
			},
			wantErr: true,
		},
		{
			name: "nil capabilities",
			config: Config{
				OpenAIBaseURL: "http://localhost:11434/v1",
				Capabilities:  nil,
			},
			wantErr: false,
		},
		{
			name: "missing scheme in URL",
			config: Config{
				OpenAIBaseURL: "://localhost:11434/v1",
				Capabilities:  nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
