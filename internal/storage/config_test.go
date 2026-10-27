package storage

import (
	"testing"
)

func TestSetEngineOverridesPackage(t *testing.T) {
	config := NewMockConfig()

	// Set package config
	err := config.Set("model", "llama", PackageConfig)
	if err != nil {
		t.Fatalf("Set package config returned unexpected error: %v", err)
	}

	// Set engine config with different value
	err = config.Set("model", "mistral", EngineConfig)
	if err != nil {
		t.Fatalf("Set engine config returned unexpected error: %v", err)
	}

	// Verify engine value overrides package value
	values, err := config.Get("model")
	if err != nil {
		t.Fatalf("Get returned unexpected error: %v", err)
	}

	if value, found := values["model"]; !found || value != "mistral" {
		t.Fatalf("expected engine config to override package config, got %#v", values)
	}
}

func TestSetUserOverridesEngine(t *testing.T) {
	config := NewMockConfig()

	// Set engine config
	err := config.Set("model", "mistral", EngineConfig)
	if err != nil {
		t.Fatalf("Set engine config returned unexpected error: %v", err)
	}

	// Set user config with different value
	err = config.Set("model", "custom-model", UserConfig)
	if err != nil {
		t.Fatalf("Set user config returned unexpected error: %v", err)
	}

	// Verify user value overrides engine value
	values, err := config.Get("model")
	if err != nil {
		t.Fatalf("Get returned unexpected error: %v", err)
	}

	if value, found := values["model"]; !found || value != "custom-model" {
		t.Fatalf("expected user config to override engine config, got %#v", values)
	}
}

func TestSetUserOverridesPackage(t *testing.T) {
	config := NewMockConfig()

	// Set package config
	err := config.Set("api.endpoint", "https://package.example.com", PackageConfig)
	if err != nil {
		t.Fatalf("Set package config returned unexpected error: %v", err)
	}

	// Set user config with different value (skipping engine level)
	err = config.Set("api.endpoint", "https://user.example.com", UserConfig)
	if err != nil {
		t.Fatalf("Set user config returned unexpected error: %v", err)
	}

	// Verify user value overrides package value
	values, err := config.Get("api.endpoint")
	if err != nil {
		t.Fatalf("Get returned unexpected error: %v", err)
	}

	if value, found := values["api.endpoint"]; !found || value != "https://user.example.com" {
		t.Fatalf("expected user config to override package config, got %#v", values)
	}
}
