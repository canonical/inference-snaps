package storage

import (
	"fmt"
	"testing"
)

// TestMigratePreservesConfig tests that migration doesn't affect unrelated config
func TestMigratePreservesConfig(t *testing.T) {
	cfg := NewMockConfig()

	if err := cfg.Set("model", "mistral", UserConfig); err != nil {
		t.Fatalf("failed to set mock model config: %v", err)
	}
	if err := cfg.Set("env.var1", "value1", UserConfig); err != nil {
		t.Fatalf("failed to set mock env config: %v", err)
	}

	if err := cfg.Migrate(); err != nil {
		t.Fatalf("Migrate failed: %v", err)
	}

	if val, err := cfg.Get("model"); err != nil {
		t.Fatalf("Get model failed: %v", err)
	} else if fmt.Sprint(val["model"]) != "mistral" {
		t.Fatalf("expected model=mistral to be preserved, got %v", val)
	}

	if val, err := cfg.Get("env.var1"); err != nil {
		t.Fatalf("Get env.var1 failed: %v", err)
	} else if fmt.Sprint(val["env.var1"]) != "value1" {
		t.Fatalf("expected env.var1=value1 to be preserved, got %v", val)
	}
}
