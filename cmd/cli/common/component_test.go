package common

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/canonical/inference-snaps-cli/pkg/snap"
	"github.com/canonical/inference-snaps-cli/pkg/storage"
)

// testDirs holds paths to the crafted directory tree used by component tests.
type testDirs struct {
	enginesDir         string
	runtimesDir        string
	modelsDir          string
	requiredComponents []string
}

// unSetEnvForTest unsets the named environment variable for the duration of the
// test and restores its original value (or absence) on cleanup.
func unSetEnvForTest(t *testing.T, key string) {
	t.Helper()
	prev, wasSet := os.LookupEnv(key)
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("os.Unsetenv(%q): %v", key, err)
	}
	t.Cleanup(func() {
		if wasSet {
			os.Setenv(key, prev)
		} else {
			os.Unsetenv(key)
		}
	})
}

// setupComponentTestDirs creates a self-contained directory tree under t.TempDir()
// with minimal engine, runtime, and model YAML manifests. No test_data is used.
//
// Layout:
//
//	<root>/
//	  engines/test-engine/engine.yaml   – references runtime "test-runtime" and model "test-model"
//	  runtimes/test-runtime/runtime.yaml – requires component "test-runtime-component"
//	  models/test-model/model.yaml       – requires component "test-model-component"
func setupComponentTestDirs(t *testing.T) testDirs {
	t.Helper()
	root := t.TempDir()

	engineYAML := `name: test-engine
runtime: test-runtime
model:
  default: test-model
  options:
  - test-model
`
	runtimeYAML := `components:
  - test-runtime-component
`
	modelYAML := `components:
  - test-model-component
`

	writeFile := func(path, content string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatalf("MkdirAll %s: %v", filepath.Dir(path), err)
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("WriteFile %s: %v", path, err)
		}
	}

	writeFile(filepath.Join(root, "engines", "test-engine", "engine.yaml"), engineYAML)
	writeFile(filepath.Join(root, "runtimes", "test-runtime", "runtime.yaml"), runtimeYAML)
	writeFile(filepath.Join(root, "models", "test-model", "model.yaml"), modelYAML)

	return testDirs{
		enginesDir:         filepath.Join(root, "engines"),
		runtimesDir:        filepath.Join(root, "runtimes"),
		modelsDir:          filepath.Join(root, "models"),
		requiredComponents: []string{"test-runtime-component", "test-model-component"},
	}
}

func TestComponentInstalled(t *testing.T) {
	t.Run("SNAP_COMPONENTS not set returns error", func(t *testing.T) {
		unSetEnvForTest(t, "SNAP_COMPONENTS")

		installed, err := ComponentInstalled("my-component")
		if err == nil {
			t.Fatal("expected error when SNAP_COMPONENTS is not set, got nil")
		}
		if installed {
			t.Error("expected false when SNAP_COMPONENTS is not set")
		}
	})

	t.Run("component directory exists returns true", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv("SNAP_COMPONENTS", tmpDir)
		if err := os.MkdirAll(filepath.Join(tmpDir, "my-component"), 0755); err != nil {
			t.Fatalf("failed to create component dir: %v", err)
		}

		installed, err := ComponentInstalled("my-component")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !installed {
			t.Error("expected true for existing component directory")
		}
	})

	t.Run("component directory does not exist returns false", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv("SNAP_COMPONENTS", tmpDir)

		installed, err := ComponentInstalled("non-existent-component")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if installed {
			t.Error("expected false for non-existent component")
		}
	})

	t.Run("component path exists but is a file returns error", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv("SNAP_COMPONENTS", tmpDir)

		filePath := filepath.Join(tmpDir, "not-a-dir-component")
		if err := os.WriteFile(filePath, []byte("not a directory"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		installed, err := ComponentInstalled("not-a-dir-component")
		if err == nil {
			t.Fatal("expected error when component path is a file, got nil")
		}
		if installed {
			t.Error("expected false when component path is a file")
		}
	})
}

// makeTestCtx builds a Context backed by a mock cache with the given active engine and model.
func makeTestCtx(t *testing.T, dirs testDirs, engine, model string) *Context {
	t.Helper()
	cache := storage.NewMockCache()
	if err := cache.SetActiveEngine(engine); err != nil {
		t.Fatalf("SetActiveEngine: %v", err)
	}
	if err := cache.SetActiveModel(model); err != nil {
		t.Fatalf("SetActiveModel: %v", err)
	}
	return &Context{
		EnginesDir:  dirs.enginesDir,
		RuntimesDir: dirs.runtimesDir,
		ModelsDir:   dirs.modelsDir,
		Cache:       cache,
	}
}

func TestWaitForComponentsWithTimeoutAndInterval(t *testing.T) {
	dirs := setupComponentTestDirs(t)

	t.Run("all components installed returns nil immediately", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv("SNAP_COMPONENTS", tmpDir)
		for _, comp := range dirs.requiredComponents {
			if err := os.MkdirAll(filepath.Join(tmpDir, comp), 0755); err != nil {
				t.Fatalf("failed to create component dir %s: %v", comp, err)
			}
		}

		ctx := makeTestCtx(t, dirs, "test-engine", "test-model")
		err := waitForComponentsWithTimeoutAndInterval(ctx, 5*time.Second, 100*time.Millisecond)
		if err != nil {
			t.Errorf("expected nil error when all components are installed, got: %v", err)
		}
	})

	t.Run("timeout when components remain missing", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv("SNAP_COMPONENTS", tmpDir)
		// No component directories created.

		ctx := makeTestCtx(t, dirs, "test-engine", "test-model")
		err := waitForComponentsWithTimeoutAndInterval(ctx, 1*time.Millisecond, 1*time.Millisecond)
		if err == nil {
			t.Fatal("expected timeout error, got nil")
		}
		if !strings.Contains(err.Error(), "timeout") {
			t.Errorf("expected error message to contain 'timeout', got: %v", err)
		}
	})

	t.Run("SNAP_COMPONENTS not set returns error", func(t *testing.T) {
		unSetEnvForTest(t, "SNAP_COMPONENTS")

		ctx := makeTestCtx(t, dirs, "test-engine", "test-model")
		err := waitForComponentsWithTimeoutAndInterval(ctx, 5*time.Second, 100*time.Millisecond)
		if err == nil {
			t.Fatal("expected error when SNAP_COMPONENTS is not set, got nil")
		}
		if !strings.Contains(err.Error(), "SNAP_COMPONENTS") {
			t.Errorf("expected error message to mention SNAP_COMPONENTS, got: %v", err)
		}
	})

	t.Run("invalid engine returns error from ComponentsRequiredByCurrentSelection", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv("SNAP_COMPONENTS", tmpDir)

		ctx := makeTestCtx(t, dirs, "non-existent-engine", "test-model")
		err := waitForComponentsWithTimeoutAndInterval(ctx, 5*time.Second, 100*time.Millisecond)
		if err == nil {
			t.Fatal("expected error for non-existent engine, got nil")
		}
		if !strings.Contains(err.Error(), "determining required components") {
			t.Errorf("expected error about determining required components, got: %v", err)
		}
	})
}

func TestComponentsRequiredByRuntime(t *testing.T) {
	dirs := setupComponentTestDirs(t)

	t.Run("returns components listed in runtime manifest", func(t *testing.T) {
		ctx := makeTestCtx(t, dirs, "test-engine", "test-model")
		got, err := ComponentsRequiredByRuntime(ctx, "test-runtime")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := []string{"test-runtime-component"}
		if len(got) != len(want) || got[0] != want[0] {
			t.Errorf("expected %v, got %v", want, got)
		}
	})

	t.Run("non-existent runtime returns error", func(t *testing.T) {
		ctx := makeTestCtx(t, dirs, "test-engine", "test-model")
		_, err := ComponentsRequiredByRuntime(ctx, "non-existent-runtime")
		if err == nil {
			t.Fatal("expected error for non-existent runtime, got nil")
		}
	})
}

func TestComponentsRequiredByCurrentSelection(t *testing.T) {
	dirs := setupComponentTestDirs(t)

	t.Run("returns combined runtime and model components", func(t *testing.T) {
		ctx := makeTestCtx(t, dirs, "test-engine", "test-model")
		got, err := ComponentsRequiredByCurrentSelection(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Expect runtime component followed by model component.
		want := []string{"test-runtime-component", "test-model-component"}
		if len(got) != len(want) {
			t.Fatalf("expected %v, got %v", want, got)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("index %d: expected %q, got %q", i, want[i], got[i])
			}
		}
	})

	t.Run("invalid engine in cache returns error", func(t *testing.T) {
		ctx := makeTestCtx(t, dirs, "non-existent-engine", "test-model")
		_, err := ComponentsRequiredByCurrentSelection(ctx)
		if err == nil {
			t.Fatal("expected error for non-existent engine, got nil")
		}
	})

	t.Run("invalid model in cache returns error", func(t *testing.T) {
		ctx := makeTestCtx(t, dirs, "test-engine", "non-existent-model")
		_, err := ComponentsRequiredByCurrentSelection(ctx)
		if err == nil {
			t.Fatal("expected error for non-existent model, got nil")
		}
	})
}

func TestInstalledComponents(t *testing.T) {
	t.Run("SNAP_COMPONENTS not set returns error", func(t *testing.T) {
		unSetEnvForTest(t, "SNAP_COMPONENTS")

		_, err := InstalledComponents()
		if err == nil {
			t.Fatal("expected error when SNAP_COMPONENTS is not set, got nil")
		}
		if !strings.Contains(err.Error(), "SNAP_COMPONENTS") {
			t.Errorf("expected error to mention SNAP_COMPONENTS, got: %v", err)
		}
	})

	t.Run("empty components directory returns empty slice", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv("SNAP_COMPONENTS", tmpDir)

		installed, err := InstalledComponents()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(installed) != 0 {
			t.Errorf("expected empty slice, got %v", installed)
		}
	})

	t.Run("returns names of installed component directories", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv("SNAP_COMPONENTS", tmpDir)
		for _, comp := range []string{"comp-a", "comp-b", "comp-c"} {
			if err := os.MkdirAll(filepath.Join(tmpDir, comp), 0755); err != nil {
				t.Fatalf("MkdirAll: %v", err)
			}
		}

		installed, err := InstalledComponents()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(installed) != 3 {
			t.Fatalf("expected 3 installed components, got %v", installed)
		}
		want := map[string]bool{"comp-a": true, "comp-b": true, "comp-c": true}
		for _, name := range installed {
			if !want[name] {
				t.Errorf("unexpected component name %q in result", name)
			}
		}
	})

	t.Run("files in components directory are not included", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv("SNAP_COMPONENTS", tmpDir)
		if err := os.MkdirAll(filepath.Join(tmpDir, "real-component"), 0755); err != nil {
			t.Fatalf("MkdirAll: %v", err)
		}
		if err := os.WriteFile(filepath.Join(tmpDir, "not-a-component"), []byte("file"), 0644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}

		installed, err := InstalledComponents()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(installed) != 1 || installed[0] != "real-component" {
			t.Errorf("expected [real-component], got %v", installed)
		}
	})

	t.Run("non-existent components directory returns error", func(t *testing.T) {
		t.Setenv("SNAP_COMPONENTS", "/this/path/does/not/exist")

		_, err := InstalledComponents()
		if err == nil {
			t.Fatal("expected error for non-existent components directory, got nil")
		}
	})
}

func TestMissingComponents(t *testing.T) {
	t.Run("SNAP_COMPONENTS not set returns error", func(t *testing.T) {
		unSetEnvForTest(t, "SNAP_COMPONENTS")

		_, err := MissingComponents([]string{"any-component"})
		if err == nil {
			t.Fatal("expected error when SNAP_COMPONENTS is not set, got nil")
		}
		if !strings.Contains(err.Error(), "SNAP_COMPONENTS") {
			t.Errorf("expected error to mention SNAP_COMPONENTS, got: %v", err)
		}
	})

	t.Run("all components present returns empty slice", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv("SNAP_COMPONENTS", tmpDir)
		for _, comp := range []string{"comp-a", "comp-b"} {
			if err := os.MkdirAll(filepath.Join(tmpDir, comp), 0755); err != nil {
				t.Fatalf("MkdirAll: %v", err)
			}
		}

		missing, err := MissingComponents([]string{"comp-a", "comp-b"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(missing) != 0 {
			t.Errorf("expected no missing components, got %v", missing)
		}
	})

	t.Run("some components missing returns only missing ones", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv("SNAP_COMPONENTS", tmpDir)
		if err := os.MkdirAll(filepath.Join(tmpDir, "comp-present"), 0755); err != nil {
			t.Fatalf("MkdirAll: %v", err)
		}
		// "comp-absent" is intentionally not created.

		missing, err := MissingComponents([]string{"comp-present", "comp-absent"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(missing) != 1 || missing[0] != "comp-absent" {
			t.Errorf("expected [comp-absent], got %v", missing)
		}
	})

	t.Run("all components missing returns all", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv("SNAP_COMPONENTS", tmpDir)

		required := []string{"comp-x", "comp-y"}
		missing, err := MissingComponents(required)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(missing) != len(required) {
			t.Errorf("expected %v, got %v", required, missing)
		}
	})

	t.Run("empty required list returns empty slice", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv("SNAP_COMPONENTS", tmpDir)

		missing, err := MissingComponents([]string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(missing) != 0 {
			t.Errorf("expected empty slice, got %v", missing)
		}
	})
}

func TestInstallComponents(t *testing.T) {
	dirs := setupComponentTestDirs(t)

	// makeInstallCtx builds a Context with the given snap mock.
	makeInstallCtx := func(s snap.Snap) *Context {
		return &Context{
			EnginesDir:  dirs.enginesDir,
			RuntimesDir: dirs.runtimesDir,
			ModelsDir:   dirs.modelsDir,
			Cache:       storage.NewMockCache(),
			Snap:        s,
		}
	}

	// callCount returns a mock InstallComponent function that tracks how many
	// times it has been called and returns errors from the provided sequence,
	// returning nil once the sequence is exhausted.
	errorSequence := func(errs ...error) func(string) error {
		i := 0
		return func(_ string) error {
			if i < len(errs) {
				err := errs[i]
				i++
				return err
			}
			return nil
		}
	}

	t.Run("empty component list returns nil without calling snap", func(t *testing.T) {
		called := false
		ctx := makeInstallCtx(snap.MockWithInstall(func(_ string) error {
			called = true
			return nil
		}))
		if err := InstallComponents(ctx, []string{}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if called {
			t.Error("expected InstallComponent to not be called for empty list")
		}
	})

	t.Run("single component installs successfully", func(t *testing.T) {
		ctx := makeInstallCtx(snap.Mock())
		if err := InstallComponents(ctx, []string{"comp-a"}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("multiple components all install successfully", func(t *testing.T) {
		var installed []string
		ctx := makeInstallCtx(snap.MockWithInstall(func(name string) error {
			installed = append(installed, name)
			return nil
		}))
		if err := InstallComponents(ctx, []string{"comp-a", "comp-b"}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(installed) != 2 {
			t.Errorf("expected 2 installs, got %v", installed)
		}
	})

	t.Run("already installed error is treated as success", func(t *testing.T) {
		ctx := makeInstallCtx(snap.MockWithInstall(errorSequence(
			errors.New("already installed"),
		)))
		if err := InstallComponents(ctx, []string{"comp-a"}); err != nil {
			t.Fatalf("expected success for 'already installed', got: %v", err)
		}
	})

	t.Run("unknown snap error returns user-facing error immediately", func(t *testing.T) {
		ctx := makeInstallCtx(snap.MockWithInstall(func(_ string) error {
			return errors.New("cannot install components for a snap that is unknown to the store")
		}))
		err := InstallComponents(ctx, []string{"comp-a"})
		if err == nil {
			t.Fatal("expected error for unknown snap, got nil")
		}
		if !strings.Contains(err.Error(), "snap not known to the store") {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("snapd timeout error retries and eventually succeeds", func(t *testing.T) {
		ctx := makeInstallCtx(snap.MockWithInstall(errorSequence(
			errors.New("timeout exceeded while waiting for response"),
			nil, // succeeds on second attempt
		)))
		if err := installComponents(ctx, []string{"comp-a"}, time.Hour, 0); err != nil {
			t.Fatalf("expected success after retry, got: %v", err)
		}
	})

	t.Run("change in progress error retries and eventually succeeds", func(t *testing.T) {
		ctx := makeInstallCtx(snap.MockWithInstall(errorSequence(
			errors.New("change in progress"),
			nil, // succeeds on second attempt
		)))
		if err := installComponents(ctx, []string{"comp-a"}, time.Hour, 0); err != nil {
			t.Fatalf("expected success after retry, got: %v", err)
		}
	})

	t.Run("unhandled error is returned immediately", func(t *testing.T) {
		ctx := makeInstallCtx(snap.MockWithInstall(func(_ string) error {
			return errors.New("some unexpected snapd error")
		}))
		err := InstallComponents(ctx, []string{"comp-a"})
		if err == nil {
			t.Fatal("expected error for unhandled snapd error, got nil")
		}
		if !strings.Contains(err.Error(), "installing") {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("overall timeout exceeded during retries returns timeout error", func(t *testing.T) {
		// Always return a retryable error so the overall-timeout branch is hit.
		ctx := makeInstallCtx(snap.MockWithInstall(func(_ string) error {
			return errors.New("timeout exceeded while waiting for response")
		}))
		err := installComponents(ctx, []string{"comp-a"}, 0, 0)
		if err == nil {
			t.Fatal("expected timeout error, got nil")
		}
		if !strings.Contains(err.Error(), "timed out while installing") {
			t.Errorf("unexpected error message: %v", err)
		}
	})
}
