package common

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// setupTestComponent creates the component directory and target file on disk,
// sets SNAP_COMPONENTS, and returns the resolved symlink path.
func setupTestComponent(t *testing.T) (symlinkPath string) {
	t.Helper()

	componentsDir := t.TempDir()
	t.Setenv("SNAP_COMPONENTS", componentsDir)

	componentPath := filepath.Join(componentsDir, "dummy-component-2")
	if err := os.MkdirAll(componentPath, 0755); err != nil {
		t.Fatalf("failed to create component dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(componentPath, "test_file.txt"), []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to create target file: %v", err)
	}

	return filepath.Join(t.TempDir(), "test_folder", "test_symlink")
}

func TestLoadEngineEnvironmentFromSettingsCollection(t *testing.T) {
	symlinkPath := setupTestComponent(t)
	settings := []ComponentSettings{
		{
			componentName: "dummy-component-2",
			Layout: map[string]ComponentLayout{
				symlinkPath: {
					Symlink: "$SNAP_COMPONENTS/dummy-component-2/test_file.txt",
				},
			},
			Environment: []string{
				"TEST_ENV_VAR=test",
			},
		},
	}

	err := loadEngineEnvironmentFromSettingsCollection(settings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify environment variable was set
	if got := os.Getenv("TEST_ENV_VAR"); got != "test" {
		t.Errorf("expected TEST_ENV_VAR=test, got %q", got)
	}

	// Verify symlink was created and points to the expanded target
	linkTarget, err := os.Readlink(symlinkPath)
	if err != nil {
		t.Fatalf("expected symlink at %s: %v", symlinkPath, err)
	}
	expectedTarget := filepath.Join(os.Getenv("SNAP_COMPONENTS"), "dummy-component-2", "test_file.txt")
	if linkTarget != expectedTarget {
		t.Errorf("expected symlink target %q, got %q", expectedTarget, linkTarget)
	}

	// Verify the file is reachable through the symlink
	content, err := os.ReadFile(symlinkPath)
	if err != nil {
		t.Fatalf("failed to read through symlink: %v", err)
	}
	if string(content) != "hello" {
		t.Errorf("expected content %q, got %q", "hello", string(content))
	}

	os.Unsetenv("TEST_ENV_VAR")
}

func TestUnloadEngineEnvironmentFromSettingsCollection(t *testing.T) {
	symlinkPath := setupTestComponent(t)
	settings := []ComponentSettings{
		{
			componentName: "dummy-component-2",
			Layout: map[string]ComponentLayout{
				symlinkPath: {
					Symlink: "$SNAP_COMPONENTS/dummy-component-2/test_file.txt",
				},
			},
			Environment: []string{
				"TEST_ENV_VAR=test",
			},
		},
	}

	// Load first so there is something to unload
	if err := loadEngineEnvironmentFromSettingsCollection(settings); err != nil {
		t.Fatalf("load failed: %v", err)
	}
	// Verify preconditions
	if os.Getenv("TEST_ENV_VAR") != "test" {
		t.Fatal("precondition: TEST_ENV_VAR not set")
	}
	if _, err := os.Lstat(symlinkPath); err != nil {
		t.Fatalf("precondition: symlink missing: %v", err)
	}

	// Unload
	if err := unloadEngineEnvironmentFromSettingsCollection(settings); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify symlink was removed
	if _, err := os.Lstat(symlinkPath); !os.IsNotExist(err) {
		t.Errorf("expected symlink to be removed, but it still exists")
	}
}

func TestSnapComponentsNotSet(t *testing.T) {
	settings := []ComponentSettings{
		{
			componentName: "dummy-component-2",
			Layout: map[string]ComponentLayout{
				"": {
					Symlink: "$SNAP_COMPONENTS/dummy-component-2/test_file.txt",
				},
			},
			Environment: []string{
				"TEST_ENV_VARtest",
			},
		},
	}
	err := loadEngineEnvironmentFromSettingsCollection(settings)
	if err.Error() != "SNAP_COMPONENTS env var not set" {
		t.Fatalf("expected error about SNAP_COMPONENTS, got: %v", err)
	}
}

func TestCannotCreateDir(t *testing.T) {
	_ = setupTestComponent(t)
	settings := []ComponentSettings{
		{
			componentName: "dummy-component-2",
			Layout: map[string]ComponentLayout{
				"/NOT/EXISTENT": {
					Symlink: "$SNAP_COMPONENTS/dummy-component-2/non_existent_file.txt",
				},
			},
			Environment: []string{
				"TEST_ENV_VAR=test",
			},
		},
	}

	err := loadEngineEnvironmentFromSettingsCollection(settings)
	if err != nil && !strings.Contains(err.Error(), "error creating directory for symlink") {
		t.Fatalf("expected skipping symlink creation due to path different from /tmp, got: %v", err)
	}
}
