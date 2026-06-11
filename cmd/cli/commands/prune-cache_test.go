package commands

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/canonical/inference-snaps-cli/v2/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/v2/pkg/storage"
	"github.com/creack/pty"
)

// pruneTestDirs holds paths to the crafted directory tree used by prune-cache tests.
type pruneTestDirs struct {
	enginesDir  string
	runtimesDir string
	modelsDir   string
}

// unsetenvForTest unsets the named environment variable for the duration of the
// test and restores its original value (or absence) on cleanup.
func unsetenvForTest(t *testing.T, key string) {
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

// setupPruneTestDirs creates a self-contained directory tree under t.TempDir()
// with minimal engine, runtime, and model YAML manifests.
//
// Layout:
//
//	engines/
//	  active-engine/engine.yaml   – runtime: active-runtime, model options: [active-model]
//	  inactive-engine/engine.yaml – runtime: inactive-runtime, model options: [model-a, model-b]
//	runtimes/
//	  active-runtime/runtime.yaml   – components: [active-runtime-comp]
//	  inactive-runtime/runtime.yaml – components: [inactive-runtime-comp, active-runtime-comp]
//	models/
//	  active-model/model.yaml – components: [active-model-comp]
//	  model-a/model.yaml      – components: [model-a-comp, shared-comp]
//	  model-b/model.yaml      – components: [model-b-comp, shared-comp]
//
// Note: inactive-runtime deliberately includes active-runtime-comp (also used by the active
// engine) to allow testing that required components are excluded from the prune list.
// model-a and model-b both include shared-comp to allow testing deduplication.
func setupPruneTestDirs(t *testing.T) pruneTestDirs {
	t.Helper()
	root := t.TempDir()

	writeFile := func(path, content string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatalf("MkdirAll %s: %v", filepath.Dir(path), err)
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("WriteFile %s: %v", path, err)
		}
	}

	writeFile(filepath.Join(root, "engines", "active-engine", "engine.yaml"), `name: active-engine
runtime: active-runtime
model:
  default: active-model
  options:
  - active-model
`)
	writeFile(filepath.Join(root, "engines", "inactive-engine", "engine.yaml"), `name: inactive-engine
runtime: inactive-runtime
model:
  default: model-a
  options:
  - model-a
  - model-b
`)
	writeFile(filepath.Join(root, "runtimes", "active-runtime", "runtime.yaml"), `components:
  - active-runtime-comp
`)
	// inactive-runtime also lists active-runtime-comp to exercise the
	// "required by active selection → must be excluded" path.
	writeFile(filepath.Join(root, "runtimes", "inactive-runtime", "runtime.yaml"), `components:
  - inactive-runtime-comp
  - active-runtime-comp
`)
	writeFile(filepath.Join(root, "models", "active-model", "model.yaml"), `components:
  - active-model-comp
`)
	// model-a and model-b both include shared-comp to exercise the deduplication path.
	writeFile(filepath.Join(root, "models", "model-a", "model.yaml"), `components:
  - model-a-comp
  - shared-comp
`)
	writeFile(filepath.Join(root, "models", "model-b", "model.yaml"), `components:
  - model-b-comp
  - shared-comp
`)

	return pruneTestDirs{
		enginesDir:  filepath.Join(root, "engines"),
		runtimesDir: filepath.Join(root, "runtimes"),
		modelsDir:   filepath.Join(root, "models"),
	}
}

// makePruneCmd builds a pruneCacheCommand backed by a mock cache with the given
// active engine and model.
func makePruneCmd(dirs pruneTestDirs, activeEngine, activeModel string) pruneCacheCommand {
	cache := storage.NewMockCache()
	_ = cache.SetActiveEngine(activeEngine)
	_ = cache.SetActiveModel(activeModel)
	return pruneCacheCommand{Context: &common.Context{
		EnginesDir:  dirs.enginesDir,
		RuntimesDir: dirs.runtimesDir,
		ModelsDir:   dirs.modelsDir,
		Cache:       cache,
	}}
}

// createInstalledComponents creates subdirectories under snapComponents to simulate
// installed snap components.
func createInstalledComponents(t *testing.T, snapComponents string, names ...string) {
	t.Helper()
	for _, name := range names {
		if err := os.MkdirAll(filepath.Join(snapComponents, name), 0755); err != nil {
			t.Fatalf("creating component dir %q: %v", name, err)
		}
	}
}

// ── unusedComponentsAll ──────────────────────────────────────────────────────

func TestUnusedComponentsAll(t *testing.T) {
	dirs := setupPruneTestDirs(t)

	t.Run("no installed components returns empty", func(t *testing.T) {
		t.Setenv("SNAP_COMPONENTS", t.TempDir())
		cmd := makePruneCmd(dirs, "active-engine", "active-model")

		result, err := cmd.unusedComponentsAll()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("expected empty result, got %v", result)
		}
	})

	t.Run("all installed components are required returns empty", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv("SNAP_COMPONENTS", tmpDir)
		createInstalledComponents(t, tmpDir, "active-runtime-comp", "active-model-comp")
		cmd := makePruneCmd(dirs, "active-engine", "active-model")

		result, err := cmd.unusedComponentsAll()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("expected empty result, got %v", result)
		}
	})

	t.Run("installed components not required by active selection are returned", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv("SNAP_COMPONENTS", tmpDir)
		// active-runtime-comp is required; inactive-runtime-comp and model-a-comp are not.
		createInstalledComponents(t, tmpDir, "active-runtime-comp", "inactive-runtime-comp", "model-a-comp")
		cmd := makePruneCmd(dirs, "active-engine", "active-model")

		result, err := cmd.unusedComponentsAll()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Fatalf("expected 2 unused components, got %v", result)
		}
		if slices.Contains(result, "active-runtime-comp") {
			t.Error("required component active-runtime-comp must not appear in unused list")
		}
		if !slices.Contains(result, "inactive-runtime-comp") || !slices.Contains(result, "model-a-comp") {
			t.Errorf("expected inactive-runtime-comp and model-a-comp in result, got %v", result)
		}
	})

	t.Run("SNAP_COMPONENTS not set returns error", func(t *testing.T) {
		unsetenvForTest(t, "SNAP_COMPONENTS")
		cmd := makePruneCmd(dirs, "active-engine", "active-model")

		_, err := cmd.unusedComponentsAll()
		if err == nil {
			t.Fatal("expected error when SNAP_COMPONENTS is not set, got nil")
		}
	})

	t.Run("invalid active engine returns error", func(t *testing.T) {
		t.Setenv("SNAP_COMPONENTS", t.TempDir())
		cmd := makePruneCmd(dirs, "non-existent-engine", "active-model")

		_, err := cmd.unusedComponentsAll()
		if err == nil {
			t.Fatal("expected error for non-existent active engine, got nil")
		}
	})

	t.Run("invalid active model returns error", func(t *testing.T) {
		t.Setenv("SNAP_COMPONENTS", t.TempDir())
		cmd := makePruneCmd(dirs, "active-engine", "non-existent-model")

		_, err := cmd.unusedComponentsAll()
		if err == nil {
			t.Fatal("expected error for non-existent active model, got nil")
		}
	})
}

// ── unusedComponentsEngine ───────────────────────────────────────────────────

func TestUnusedComponentsEngine(t *testing.T) {
	dirs := setupPruneTestDirs(t)

	t.Run("pruning the active engine returns error", func(t *testing.T) {
		t.Setenv("SNAP_COMPONENTS", t.TempDir())
		cmd := makePruneCmd(dirs, "active-engine", "active-model")

		_, err := cmd.unusedComponentsEngine("active-engine")
		if err == nil {
			t.Fatal("expected error when attempting to prune the active engine, got nil")
		}
	})

	t.Run("non-existent engine manifest returns error", func(t *testing.T) {
		t.Setenv("SNAP_COMPONENTS", t.TempDir())
		cmd := makePruneCmd(dirs, "active-engine", "active-model")

		_, err := cmd.unusedComponentsEngine("non-existent-engine")
		if err == nil {
			t.Fatal("expected error for non-existent engine manifest, got nil")
		}
	})

	t.Run("invalid active engine causes ComponentsRequiredByCurrentSelection to fail", func(t *testing.T) {
		t.Setenv("SNAP_COMPONENTS", t.TempDir())
		cmd := makePruneCmd(dirs, "non-existent-engine", "active-model")

		_, err := cmd.unusedComponentsEngine("inactive-engine")
		if err == nil {
			t.Fatal("expected error when active engine manifest is missing, got nil")
		}
	})

	t.Run("no components installed returns empty", func(t *testing.T) {
		t.Setenv("SNAP_COMPONENTS", t.TempDir())
		cmd := makePruneCmd(dirs, "active-engine", "active-model")

		result, err := cmd.unusedComponentsEngine("inactive-engine")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("expected empty when no components are installed, got %v", result)
		}
	})

	t.Run("installed components of engine are returned", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv("SNAP_COMPONENTS", tmpDir)
		createInstalledComponents(t, tmpDir, "inactive-runtime-comp", "model-a-comp")
		cmd := makePruneCmd(dirs, "active-engine", "active-model")

		result, err := cmd.unusedComponentsEngine("inactive-engine")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Fatalf("expected 2 components, got %v", result)
		}
		if !slices.Contains(result, "inactive-runtime-comp") || !slices.Contains(result, "model-a-comp") {
			t.Errorf("expected inactive-runtime-comp and model-a-comp, got %v", result)
		}
	})

	t.Run("component not installed is not returned", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv("SNAP_COMPONENTS", tmpDir)
		// Install only one of the engine's components.
		createInstalledComponents(t, tmpDir, "inactive-runtime-comp")
		cmd := makePruneCmd(dirs, "active-engine", "active-model")

		result, err := cmd.unusedComponentsEngine("inactive-engine")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !slices.Contains(result, "inactive-runtime-comp") {
			t.Errorf("expected inactive-runtime-comp in result, got %v", result)
		}
		if slices.Contains(result, "model-a-comp") {
			t.Error("uninstalled model-a-comp must not appear in result")
		}
	})

	t.Run("component required by active selection is excluded", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv("SNAP_COMPONENTS", tmpDir)
		// active-runtime-comp is listed in inactive-runtime but is also required by
		// the active engine's selection, so it must not appear in the prune list.
		createInstalledComponents(t, tmpDir, "active-runtime-comp", "inactive-runtime-comp")
		cmd := makePruneCmd(dirs, "active-engine", "active-model")

		result, err := cmd.unusedComponentsEngine("inactive-engine")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if slices.Contains(result, "active-runtime-comp") {
			t.Error("active-runtime-comp is required by the active selection and must not be in prune list")
		}
		if !slices.Contains(result, "inactive-runtime-comp") {
			t.Errorf("inactive-runtime-comp should be in result, got %v", result)
		}
	})

	t.Run("component shared across two models is deduplicated", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv("SNAP_COMPONENTS", tmpDir)
		// shared-comp appears in both model-a and model-b of inactive-engine.
		createInstalledComponents(t, tmpDir, "shared-comp")
		cmd := makePruneCmd(dirs, "active-engine", "active-model")

		result, err := cmd.unusedComponentsEngine("inactive-engine")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		count := 0
		for _, name := range result {
			if name == "shared-comp" {
				count++
			}
		}
		if count != 1 {
			t.Errorf("shared-comp should appear exactly once, got %d occurrence(s) in %v", count, result)
		}
	})

	t.Run("SNAP_COMPONENTS not set returns error", func(t *testing.T) {
		unsetenvForTest(t, "SNAP_COMPONENTS")
		cmd := makePruneCmd(dirs, "active-engine", "active-model")

		_, err := cmd.unusedComponentsEngine("inactive-engine")
		if err == nil {
			t.Fatal("expected error when SNAP_COMPONENTS is not set, got nil")
		}
	})
}

// ── printComponentsAndConfirm ────────────────────────────────────────────────

func TestPrintComponentsAndConfirm(t *testing.T) {
	dirs := setupPruneTestDirs(t)
	cmd := makePruneCmd(dirs, "active-engine", "active-model")

	t.Run("non-TTY auto-confirms without prompting", func(t *testing.T) {
		// In a test environment stdout is not a TTY, so the confirmation prompt is
		// skipped and the function returns confirmed=true immediately.
		confirmed, err := cmd.printComponentsAndConfirm([]string{"comp-a", "comp-b"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !confirmed {
			t.Error("expected confirmed=true in non-TTY mode")
		}
	})

	t.Run("TTY: user confirms with 'y' returns true", func(t *testing.T) {
		ptyMaster, ptySlave, err := pty.Open()
		if err != nil {
			t.Fatalf("failed to open pty: %v", err)
		}
		defer ptyMaster.Close()
		defer ptySlave.Close()

		// Redirect both stdout and stdin to the PTY slave so that
		// IsTerminalOutput() returns true and PromptYN reads from the slave.
		origStdout, origStdin := os.Stdout, os.Stdin
		os.Stdout = ptySlave
		os.Stdin = ptySlave
		t.Cleanup(func() { os.Stdout = origStdout; os.Stdin = origStdin })

		go func() { _, _ = ptyMaster.Write([]byte("y\n")) }()

		confirmed, err := cmd.printComponentsAndConfirm([]string{"comp-a"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !confirmed {
			t.Error("expected confirmed=true when user enters 'y'")
		}
	})

	t.Run("TTY: user declines with 'n' returns false", func(t *testing.T) {
		ptyMaster, ptySlave, err := pty.Open()
		if err != nil {
			t.Fatalf("failed to open pty: %v", err)
		}
		defer ptyMaster.Close()
		defer ptySlave.Close()

		origStdout, origStdin := os.Stdout, os.Stdin
		os.Stdout = ptySlave
		os.Stdin = ptySlave
		t.Cleanup(func() { os.Stdout = origStdout; os.Stdin = origStdin })

		go func() { _, _ = ptyMaster.Write([]byte("n\n")) }()

		confirmed, err := cmd.printComponentsAndConfirm([]string{"comp-a"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if confirmed {
			t.Error("expected confirmed=false when user enters 'n'")
		}
	})
}
