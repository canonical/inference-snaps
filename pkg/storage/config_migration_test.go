package storage

import (
	"fmt"
	"testing"
)

func TestMigratePassthroughEnv(t *testing.T) {
	cfg := NewMockConfig()

	if err := cfg.Set("passthrough.environment.var1", "value1", UserConfig); err != nil {
		t.Fatalf("failed to set passthrough.environment.var1: %v", err)
	}
	if err := cfg.Set("passthrough.environment.var2", "value2", UserConfig); err != nil {
		t.Fatalf("failed to set passthrough.environment.var2: %v", err)
	}

	// Run migration
	err := cfg.Migrate()
	if err != nil {
		t.Fatalf("Migrate failed: %v", err)
	}

	// Verify passthrough.environment values were migrated to env
	if val, err := cfg.Get("env.var1"); err != nil {
		t.Fatalf("Get env failed: %v", err)
	} else if val["env.var1"] != "value1" {
		t.Fatalf("expected env.var1=value1, got %v", val)
	}
	if val, err := cfg.Get("env.var2"); err != nil {
		t.Fatalf("Get env failed: %v", err)
	} else if val["env.var2"] != "value2" {
		t.Fatalf("expected env.var2=value2, got %v", val)
	}

	// Verify old passthrough values were removed
	if val, err := cfg.Get("passthrough.environment.var1"); err != nil {
		t.Fatalf("Get passthrough.environment failed: %v", err)
	} else if len(val) > 0 {
		t.Fatalf("expected passthrough.environment.var1 to be unset, got %v", val)
	}
	if val, err := cfg.Get("passthrough.environment.var2"); err != nil {
		t.Fatalf("Get passthrough.environment failed: %v", err)
	} else if len(val) > 0 {
		t.Fatalf("expected passthrough.environment.var2 to be unset, got %v", val)
	}
}

// TestMigratePassthroughEnvPreservesOtherConfig tests that migration doesn't affect other config
func TestMigratePassthroughEnvPreservesOtherConfig(t *testing.T) {
	cfg := NewMockConfig()

	// Set up passthrough config
	if err := cfg.Set("passthrough.environment.var1", "value1", UserConfig); err != nil {
		t.Fatalf("failed to set mock passthrough config: %v", err)
	}
	// Set other config values
	if err := cfg.Set("model", "mistral", UserConfig); err != nil {
		t.Fatalf("failed to set mock model config: %v", err)
	}
	if err := cfg.Set("engine", "llama", UserConfig); err != nil {
		t.Fatalf("failed to set mock engine config: %v", err)
	}

	// Run migration
	err := cfg.Migrate()
	if err != nil {
		t.Fatalf("Migrate failed: %v", err)
	}

	// Verify non-passthrough config is preserved
	if val, err := cfg.Get("model"); err != nil {
		t.Fatalf("Get model failed: %v", err)
	} else if fmt.Sprint(val["model"]) != "mistral" {
		t.Fatalf("expected model=mistral to be preserved, got %v", val)
	}

	if val, err := cfg.Get("engine"); err != nil {
		t.Fatalf("Get engine failed: %v", err)
	} else if fmt.Sprint(val["engine"]) != "llama" {
		t.Fatalf("expected engine=llama to be preserved, got %v", val)
	}
}

// TestMigratePassthroughEnvNoPassthroughConfig tests migration when no passthrough config exists
func TestMigratePassthroughEnvNoPassthroughConfig(t *testing.T) {
	cfg := NewMockConfig()

	// Set other config values
	if err := cfg.Set("model", "mistral", UserConfig); err != nil {
		t.Fatalf("failed to set mock model config: %v", err)
	}

	// Run migration with no passthrough config
	err := cfg.Migrate()
	if err != nil {
		t.Fatalf("Migrate with no passthrough config should not error: %v", err)
	}

	// Verify non-passthrough config is still there
	if modelVars, err := cfg.Get("model"); err != nil {
		t.Fatalf("Get model failed: %v", err)
	} else if fmt.Sprint(modelVars["model"]) != "mistral" {
		t.Fatalf("expected model=mistral to be preserved, got %v", modelVars)
	}
}
