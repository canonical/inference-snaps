package common

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/storage"
	"github.com/canonical/inference-snaps-cli/pkg/types"
)

// errCache is a storage.Cache that returns errors from the specified methods.
type errCache struct {
	failGetEngine bool
	failGetModel  bool
}

func (c *errCache) SetActiveEngine(string) error { return nil }
func (c *errCache) SetActiveModel(string) error  { return nil }
func (c *errCache) GetActiveEngine() (string, error) {
	if c.failGetEngine {
		return "", errors.New("cache error: GetActiveEngine")
	}
	return "test-engine", nil
}
func (c *errCache) GetActiveModel() (string, error) {
	if c.failGetModel {
		return "", errors.New("cache error: GetActiveModel")
	}
	return "test-model", nil
}

// writeFile writes content to path, creating all parent directories.
func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("creating directories for %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writing file %s: %v", path, err)
	}
}

// setupTestComponent creates a component directory with a target file on disk,
// sets SNAP_COMPONENTS, and returns the path where the test symlink should be created.
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

// setupEngineContext creates temp directories with the minimal manifest YAML files
// required by EngineSettings and returns a populated Context.
// runtimeLayoutYAML and modelLayoutYAML are raw YAML fragments placed under "layout:" in each manifest.
func setupEngineContext(t *testing.T, runtimeEnv, modelEnv []string, runtimeLayoutYAML, modelLayoutYAML string) *Context {
	t.Helper()

	base := t.TempDir()
	enginesDir := filepath.Join(base, "engines")
	runtimesDir := filepath.Join(base, "runtimes")
	modelsDir := filepath.Join(base, "models")

	// Engine manifest points at "test-runtime"
	writeFile(t, filepath.Join(enginesDir, "test-engine", "engine.yaml"),
		"name: test-engine\nruntime: test-runtime\n")

	// Runtime manifest
	runtimeYAML := "servers: {}\n"
	if len(runtimeEnv) > 0 {
		runtimeYAML += "environment:\n"
		for _, e := range runtimeEnv {
			runtimeYAML += "  - " + e + "\n"
		}
	}
	if runtimeLayoutYAML != "" {
		runtimeYAML += "layout:\n" + runtimeLayoutYAML
	}
	writeFile(t, filepath.Join(runtimesDir, "test-runtime", "runtime.yaml"), runtimeYAML)

	// Model manifest
	modelYAML := ""
	if len(modelEnv) > 0 {
		modelYAML += "environment:\n"
		for _, e := range modelEnv {
			modelYAML += "  - " + e + "\n"
		}
	}
	if modelLayoutYAML != "" {
		modelYAML += "layout:\n" + modelLayoutYAML
	}
	writeFile(t, filepath.Join(modelsDir, "test-model", "model.yaml"), modelYAML)

	cache := storage.NewMockCache()
	_ = cache.SetActiveEngine("test-engine")
	_ = cache.SetActiveModel("test-model")

	return &Context{
		EnginesDir:  enginesDir,
		RuntimesDir: runtimesDir,
		ModelsDir:   modelsDir,
		Cache:       cache,
		Config:      storage.NewMockConfig(),
	}
}

// ---- loadEngineEnvironmentFromSettings ----

func TestLoadEngineEnvironmentFromSettings(t *testing.T) {
	symlinkPath := setupTestComponent(t)
	settings := &Settings{
		Layout: map[string]types.Layout{
			symlinkPath: {Symlink: "$SNAP_COMPONENTS/dummy-component-2/test_file.txt"},
		},
		Environment: []string{"TEST_ENV_VAR=test"},
	}
	t.Cleanup(func() { os.Unsetenv("TEST_ENV_VAR") })

	if err := loadEngineEnvironmentFromSettings(settings); err != nil {
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
}

func TestLoadEngineEnvironmentEnvVarExpansion(t *testing.T) {
	t.Setenv("BASE_VAR", "/some/base")
	settings := &Settings{
		Environment: []string{"DERIVED_VAR=$BASE_VAR/sub"},
	}
	t.Cleanup(func() { os.Unsetenv("DERIVED_VAR") })

	if err := loadEngineEnvironmentFromSettings(settings); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := os.Getenv("DERIVED_VAR"); got != "/some/base/sub" {
		t.Errorf("expected /some/base/sub, got %q", got)
	}
}

func TestLoadEngineEnvironmentSkipsEmptySymlink(t *testing.T) {
	// A layout entry with an empty Symlink value should not create any file
	linkPath := filepath.Join(t.TempDir(), "should-not-exist")
	settings := &Settings{
		Layout: map[string]types.Layout{
			linkPath: {Symlink: ""},
		},
	}

	if err := loadEngineEnvironmentFromSettings(settings); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Lstat(linkPath); !os.IsNotExist(err) {
		t.Errorf("expected no symlink to be created for empty Symlink value")
	}
}

func TestInvalidEnvVarFormat(t *testing.T) {
	settings := &Settings{
		Environment: []string{"INVALID_NO_EQUALS"},
	}
	err := loadEngineEnvironmentFromSettings(settings)
	if err == nil || !strings.Contains(err.Error(), "invalid env var") {
		t.Fatalf("expected 'invalid env var' error, got: %v", err)
	}
}

func TestRejectsLayoutOutsideTmp(t *testing.T) {
	setupTestComponent(t)
	settings := &Settings{
		Layout: map[string]types.Layout{
			"/not/tmp": {Symlink: "$SNAP_COMPONENTS/dummy-component-2/non_existent_file.txt"},
		},
	}

	err := loadEngineEnvironmentFromSettings(settings)
	if err == nil {
		t.Fatal("expected error for layout path outside /tmp, got nil")
	}
	if !strings.Contains(err.Error(), "layout path outside of /tmp") {
		t.Fatalf("expected outside /tmp error, got: %v", err)
	}
}

// ---- unloadEngineEnvironmentFromSettings ----

func TestUnloadEngineEnvironmentFromSettings(t *testing.T) {
	symlinkPath := setupTestComponent(t)
	settings := &Settings{
		Layout: map[string]types.Layout{
			symlinkPath: {Symlink: "$SNAP_COMPONENTS/dummy-component-2/test_file.txt"},
		},
		Environment: []string{"TEST_ENV_VAR=test"},
	}
	t.Cleanup(func() { os.Unsetenv("TEST_ENV_VAR") })

	// Load first so there is something to unload
	if err := loadEngineEnvironmentFromSettings(settings); err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if _, err := os.Lstat(symlinkPath); err != nil {
		t.Fatalf("precondition: symlink missing: %v", err)
	}

	if err := unloadEngineEnvironmentFromSettings(settings); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify symlink was removed
	if _, err := os.Lstat(symlinkPath); !os.IsNotExist(err) {
		t.Errorf("expected symlink to be removed, but it still exists")
	}
}

func TestUnloadEngineEnvironmentEmptyLayout(t *testing.T) {
	// Unloading with no expandedLayout should be a no-op
	if err := unloadEngineEnvironmentFromSettings(&Settings{}); err != nil {
		t.Fatalf("expected no error for empty layout, got: %v", err)
	}
}

// ---- EngineSettings ----

func TestEngineSettingsNoActiveEngine(t *testing.T) {
	ctx := &Context{Cache: storage.NewMockCache()}
	_, err := EngineSettings(ctx)
	if !errors.Is(err, ErrNoActiveEngine) {
		t.Fatalf("expected ErrNoActiveEngine, got: %v", err)
	}
}

func TestEngineSettingsNoActiveModel(t *testing.T) {
	base := t.TempDir()
	enginesDir := filepath.Join(base, "engines")
	runtimesDir := filepath.Join(base, "runtimes")

	writeFile(t, filepath.Join(enginesDir, "test-engine", "engine.yaml"),
		"name: test-engine\nruntime: test-runtime\n")
	writeFile(t, filepath.Join(runtimesDir, "test-runtime", "runtime.yaml"),
		"servers: {}\n")

	cache := storage.NewMockCache()
	_ = cache.SetActiveEngine("test-engine")
	// no active model set

	ctx := &Context{
		EnginesDir:  enginesDir,
		RuntimesDir: runtimesDir,
		ModelsDir:   filepath.Join(base, "models"),
		Cache:       cache,
		Config:      storage.NewMockConfig(),
	}

	_, err := EngineSettings(ctx)
	if !errors.Is(err, ErrNoActiveModel) {
		t.Fatalf("expected ErrNoActiveModel, got: %v", err)
	}
}

func TestEngineSettingsMergesRuntimeAndModelEnv(t *testing.T) {
	ctx := setupEngineContext(t,
		[]string{"RUNTIME_VAR=runtime_val"},
		[]string{"MODEL_VAR=model_val"},
		"", "",
	)

	settings, err := EngineSettings(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(settings.Environment) != 2 {
		t.Fatalf("expected 2 env vars, got %d: %v", len(settings.Environment), settings.Environment)
	}
	// Runtime env comes first, model env appended after
	if settings.Environment[0] != "RUNTIME_VAR=runtime_val" {
		t.Errorf("expected runtime env first, got %q", settings.Environment[0])
	}
	if settings.Environment[1] != "MODEL_VAR=model_val" {
		t.Errorf("expected model env second, got %q", settings.Environment[1])
	}
}

func TestEngineSettingsModelOverridesRuntimeLayout(t *testing.T) {
	sharedKey := filepath.Join(t.TempDir(), "shared-link")

	ctx := setupEngineContext(t, nil, nil,
		fmt.Sprintf("  %s:\n    symlink: /runtime/target\n", sharedKey),
		fmt.Sprintf("  %s:\n    symlink: /model/target\n", sharedKey),
	)

	settings, err := EngineSettings(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := settings.Layout[sharedKey].Symlink; got != "/model/target" {
		t.Errorf("expected model to override runtime layout, got %q", got)
	}
}

// ---- LoadEngineEnvironment ----

func TestLoadEngineEnvironmentSucceeds(t *testing.T) {
	symlinkPath := setupTestComponent(t)
	ctx := setupEngineContext(t,
		[]string{"ENGINE_TEST_VAR=hello"},
		nil,
		fmt.Sprintf("  %s:\n    symlink: $SNAP_COMPONENTS/dummy-component-2/test_file.txt\n", symlinkPath),
		"",
	)
	t.Cleanup(func() { os.Unsetenv("ENGINE_TEST_VAR") })

	clean, err := LoadEngineEnvironment(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if clean == nil {
		t.Fatal("expected non-nil cleanup func")
	}

	if got := os.Getenv("ENGINE_TEST_VAR"); got != "hello" {
		t.Errorf("expected ENGINE_TEST_VAR=hello, got %q", got)
	}
	if _, err := os.Lstat(symlinkPath); err != nil {
		t.Fatalf("expected symlink to exist: %v", err)
	}

	// Calling the cleanup func should remove the symlink
	clean()
	if _, err := os.Lstat(symlinkPath); !os.IsNotExist(err) {
		t.Error("expected symlink to be removed after cleanup")
	}
}

func TestLoadEngineEnvironmentReturnsNilCleanupOnError(t *testing.T) {
	// No active engine → EngineSettings fails → must return nil cleanup
	ctx := &Context{Cache: storage.NewMockCache(), Config: storage.NewMockConfig()}
	clean, err := LoadEngineEnvironment(ctx)
	if err == nil {
		t.Fatal("expected error")
	}
	if clean != nil {
		t.Error("expected nil cleanup func when LoadEngineEnvironment fails")
	}
}

// ---- SetEngineConfig ----

func TestSetEngineConfig(t *testing.T) {
	cfg := storage.NewMockConfig()
	ctx := &Context{Config: cfg}
	manifest := &engines.Manifest{
		Configurations: engines.Configurations{
			"key1": "value1",
			"key2": 42,
		},
	}

	if err := SetEngineConfig(manifest, ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	all, _ := cfg.GetAll()
	if all["key1"] != "value1" {
		t.Errorf("expected key1=value1, got %v", all["key1"])
	}
	if all["key2"] != 42 {
		t.Errorf("expected key2=42, got %v", all["key2"])
	}
}

func TestSetEngineConfigEmpty(t *testing.T) {
	ctx := &Context{Config: storage.NewMockConfig()}
	if err := SetEngineConfig(&engines.Manifest{}, ctx); err != nil {
		t.Fatalf("unexpected error for empty configurations: %v", err)
	}
}

func TestSetEngineConfigDocumentError(t *testing.T) {
	ctx := &Context{Config: storage.NewFailingMockConfig(errors.New("set doc error"))}
	manifest := &engines.Manifest{
		Configurations: engines.Configurations{"key": "value"},
	}
	err := SetEngineConfig(manifest, ctx)
	if err == nil || !strings.Contains(err.Error(), "set doc error") {
		t.Fatalf("expected set doc error, got: %v", err)
	}
}

// ---- UnsetEngineConfig ----

func TestUnsetEngineConfigSkipsUserOverridesWhenFlagFalse(t *testing.T) {
	// When unsetUserOverrides=false no manifest is loaded; only engine config namespace is cleared
	ctx := &Context{Config: storage.NewMockConfig()}
	if err := UnsetEngineConfig("nonexistent-engine", false, ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnsetEngineConfigUnsetError(t *testing.T) {
	// Unset(".", EngineConfig) returns an error
	ctx := &Context{Config: storage.NewFailingMockConfig(errors.New("unset error"))}
	err := UnsetEngineConfig("any-engine", false, ctx)
	if err == nil || !strings.Contains(err.Error(), "unset error") {
		t.Fatalf("expected unset error, got: %v", err)
	}
}

func TestUnsetEngineConfigMissingManifestVerbose(t *testing.T) {
	// When the manifest is missing and Verbose=true, a warning is printed but nil is returned
	ctx := &Context{
		EnginesDir: t.TempDir(),
		Config:     storage.NewMockConfig(),
		Verbose:    true,
	}
	if err := UnsetEngineConfig("missing-engine", true, ctx); err != nil {
		t.Fatalf("expected nil for missing manifest, got: %v", err)
	}
}

func TestUnsetEngineConfigUnsetUserOverrideError(t *testing.T) {
	// Unset(k, UserConfig) returns an error when removing a user override.
	// Use a selective-failing config so the initial Unset(".", EngineConfig) succeeds
	// but Unset("mykey", UserConfig) fails.
	base := t.TempDir()
	enginesDir := filepath.Join(base, "engines")
	writeFile(t, filepath.Join(enginesDir, "my-engine", "engine.yaml"),
		"name: my-engine\nruntime: rt\nconfigurations:\n  mykey: default\n")

	ctx := &Context{
		EnginesDir: enginesDir,
		Config:     storage.NewSelectiveFailingMockConfig("mykey", errors.New("unset key error")),
	}

	err := UnsetEngineConfig("my-engine", true, ctx)
	if err == nil || !strings.Contains(err.Error(), "unset key error") {
		t.Fatalf("expected unset key error, got: %v", err)
	}
}

func TestUnsetEngineConfigManifestReadError(t *testing.T) {
	// Make engine.yaml a directory so os.ReadFile returns an error other than ErrNotExist,
	// hitting the non-ErrManifestNotFound error branch in UnsetEngineConfig.
	base := t.TempDir()
	enginesDir := filepath.Join(base, "engines")
	if err := os.MkdirAll(filepath.Join(enginesDir, "bad-engine", "engine.yaml"), 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	ctx := &Context{EnginesDir: enginesDir, Config: storage.NewMockConfig()}
	err := UnsetEngineConfig("bad-engine", true, ctx)
	if err == nil {
		t.Fatal("expected error for unreadable manifest")
	}
}

func TestUnsetEngineConfigMissingManifest(t *testing.T) {
	// When unsetUserOverrides=true but the manifest is gone, should return nil gracefully
	ctx := &Context{
		EnginesDir: t.TempDir(), // empty — no manifest files
		Config:     storage.NewMockConfig(),
	}
	if err := UnsetEngineConfig("missing-engine", true, ctx); err != nil {
		t.Fatalf("expected nil for missing manifest, got: %v", err)
	}
}

func TestUnsetEngineConfigUnsetsUserOverrides(t *testing.T) {
	base := t.TempDir()
	enginesDir := filepath.Join(base, "engines")
	writeFile(t, filepath.Join(enginesDir, "my-engine", "engine.yaml"),
		"name: my-engine\nruntime: rt\nconfigurations:\n  mykey: default\n")

	cfg := storage.NewMockConfig()
	_ = cfg.Set("mykey", "user-override", storage.UserConfig)

	ctx := &Context{EnginesDir: enginesDir, Config: cfg}

	if err := UnsetEngineConfig("my-engine", true, ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	all, _ := cfg.GetAll()
	if _, exists := all["user.mykey"]; exists {
		t.Errorf("expected user override to be unset, but it still exists")
	}
}

// ---- EngineSettings additional error paths ----

func TestEngineSettingsCacheErrorOnGetEngine(t *testing.T) {
	ctx := &Context{Cache: &errCache{failGetEngine: true}}
	_, err := EngineSettings(ctx)
	if err == nil || !strings.Contains(err.Error(), "cache error: GetActiveEngine") {
		t.Fatalf("expected cache error for GetActiveEngine, got: %v", err)
	}
}

func TestEngineSettingsCacheErrorOnGetModel(t *testing.T) {
	base := t.TempDir()
	enginesDir := filepath.Join(base, "engines")
	runtimesDir := filepath.Join(base, "runtimes")
	writeFile(t, filepath.Join(enginesDir, "test-engine", "engine.yaml"),
		"name: test-engine\nruntime: test-runtime\n")
	writeFile(t, filepath.Join(runtimesDir, "test-runtime", "runtime.yaml"),
		"servers: {}\n")

	ctx := &Context{
		EnginesDir:  enginesDir,
		RuntimesDir: runtimesDir,
		ModelsDir:   filepath.Join(base, "models"),
		Cache:       &errCache{failGetModel: true},
		Config:      storage.NewMockConfig(),
	}
	_, err := EngineSettings(ctx)
	if err == nil || !strings.Contains(err.Error(), "cache error: GetActiveModel") {
		t.Fatalf("expected cache error for GetActiveModel, got: %v", err)
	}
}

func TestEngineSettingsMissingEngineManifest(t *testing.T) {
	cache := storage.NewMockCache()
	_ = cache.SetActiveEngine("nonexistent-engine")
	ctx := &Context{
		EnginesDir: t.TempDir(), // empty — no engine manifest
		Cache:      cache,
		Config:     storage.NewMockConfig(),
	}
	_, err := EngineSettings(ctx)
	if err == nil {
		t.Fatal("expected error for missing engine manifest")
	}
}

func TestEngineSettingsMissingRuntimeManifest(t *testing.T) {
	base := t.TempDir()
	enginesDir := filepath.Join(base, "engines")
	// Engine manifest references "missing-runtime" which has no directory
	writeFile(t, filepath.Join(enginesDir, "test-engine", "engine.yaml"),
		"name: test-engine\nruntime: missing-runtime\n")

	cache := storage.NewMockCache()
	_ = cache.SetActiveEngine("test-engine")
	ctx := &Context{
		EnginesDir:  enginesDir,
		RuntimesDir: t.TempDir(), // empty — no runtime manifest
		Cache:       cache,
		Config:      storage.NewMockConfig(),
	}
	_, err := EngineSettings(ctx)
	if err == nil {
		t.Fatal("expected error for missing runtime manifest")
	}
}

func TestEngineSettingsMissingModelManifest(t *testing.T) {
	base := t.TempDir()
	enginesDir := filepath.Join(base, "engines")
	runtimesDir := filepath.Join(base, "runtimes")
	writeFile(t, filepath.Join(enginesDir, "test-engine", "engine.yaml"),
		"name: test-engine\nruntime: test-runtime\n")
	writeFile(t, filepath.Join(runtimesDir, "test-runtime", "runtime.yaml"),
		"servers: {}\n")

	cache := storage.NewMockCache()
	_ = cache.SetActiveEngine("test-engine")
	_ = cache.SetActiveModel("nonexistent-model")
	ctx := &Context{
		EnginesDir:  enginesDir,
		RuntimesDir: runtimesDir,
		ModelsDir:   t.TempDir(), // empty — no model manifest
		Cache:       cache,
		Config:      storage.NewMockConfig(),
	}
	_, err := EngineSettings(ctx)
	if err == nil {
		t.Fatal("expected error for missing model manifest")
	}
}

// ---- unloadEngineEnvironmentFromSettings error path ----

func TestUnloadEngineEnvironmentErrorPath(t *testing.T) {
	// A regular file (not a symlink) in /tmp causes RemoveTempSymlink to return an error.
	regularFile := filepath.Join(t.TempDir(), "not-a-symlink")
	if err := os.WriteFile(regularFile, []byte("data"), 0644); err != nil {
		t.Fatalf("creating regular file: %v", err)
	}

	// Set expandedLayout directly (same package access)
	settings := &Settings{
		expandedLayout: map[string]types.Layout{
			regularFile: {Symlink: "/some/target"},
		},
	}

	err := unloadEngineEnvironmentFromSettings(settings)
	if err == nil {
		t.Fatal("expected error when removing a non-symlink path")
	}
	if !strings.Contains(err.Error(), "not a symlink") {
		t.Errorf("expected 'not a symlink' error, got: %v", err)
	}
}

// ---- LoadEngineEnvironment verbose warning paths ----

func TestLoadEngineEnvironmentVerboseWarningOnCleanupError(t *testing.T) {
	// Pre-create a regular file at the link path so CreateTempSymlink fails,
	// and the subsequent cleanup also fails (RemoveTempSymlink rejects non-symlinks).
	linkPath := filepath.Join(t.TempDir(), "blocking-file")
	if err := os.WriteFile(linkPath, []byte("block"), 0644); err != nil {
		t.Fatalf("creating blocking file: %v", err)
	}

	ctx := setupEngineContext(t, nil, nil,
		fmt.Sprintf("  %s:\n    symlink: /some/target\n", linkPath),
		"",
	)
	ctx.Verbose = true

	clean, err := LoadEngineEnvironment(ctx)
	// Load must fail (can't create symlink over a regular file)
	if err == nil {
		t.Fatal("expected error")
	}
	// Cleanup func must be nil on error
	if clean != nil {
		t.Error("expected nil cleanup func")
	}
}

func TestLoadEngineEnvironmentCleanupFuncVerboseWarning(t *testing.T) {
	// Load succeeds, then replace the symlink with a regular file so the
	// cleanup func encounters a "not a symlink" error and prints a verbose warning.
	symlinkPath := setupTestComponent(t)
	ctx := setupEngineContext(t, nil, nil,
		fmt.Sprintf("  %s:\n    symlink: $SNAP_COMPONENTS/dummy-component-2/test_file.txt\n", symlinkPath),
		"",
	)
	ctx.Verbose = true

	clean, err := LoadEngineEnvironment(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Remove the symlink and put a regular file in its place so cleanup fails
	if err := os.Remove(symlinkPath); err != nil {
		t.Fatalf("removing symlink: %v", err)
	}
	if err := os.WriteFile(symlinkPath, []byte("block"), 0644); err != nil {
		t.Fatalf("creating blocking file: %v", err)
	}

	// cleanup func must not panic; it only prints a warning to stderr
	clean()
}

// ---- ScoreEngines ----

func TestScoreEnginesLoadManifestsError(t *testing.T) {
	ctx := &Context{
		EnginesDir: filepath.Join(t.TempDir(), "nonexistent"), // does not exist
		Config:     storage.NewMockConfig(),
		Cache:      storage.NewMockCache(),
	}
	_, _, err := ScoreEngines(ctx)
	if err == nil {
		t.Fatal("expected error from missing engines dir")
	}
}

func TestScoreEnginesHardwareInfoError(t *testing.T) {
	// Set up an engines dir so LoadManifests succeeds (returns empty slice)
	orig := hardwareInfoGet
	t.Cleanup(func() { hardwareInfoGet = orig })
	hardwareInfoGet = func(bool) (*types.HwInfo, []string, error) {
		return nil, nil, errors.New("hw error")
	}

	ctx := &Context{
		EnginesDir: t.TempDir(), // empty but valid dir
		Config:     storage.NewMockConfig(),
		Cache:      storage.NewMockCache(),
	}
	_, _, err := ScoreEngines(ctx)
	if err == nil || !strings.Contains(err.Error(), "hw error") {
		t.Fatalf("expected hw error, got: %v", err)
	}
}

func TestScoreEnginesScorerError(t *testing.T) {
	origGet := hardwareInfoGet
	origScorer := engineScorer
	t.Cleanup(func() {
		hardwareInfoGet = origGet
		engineScorer = origScorer
	})
	hardwareInfoGet = func(bool) (*types.HwInfo, []string, error) {
		return &types.HwInfo{}, nil, nil
	}
	engineScorer = func(*types.HwInfo, []engines.Manifest) ([]engines.ScoredManifest, error) {
		return nil, errors.New("scorer error")
	}

	ctx := &Context{
		EnginesDir: t.TempDir(),
		Config:     storage.NewMockConfig(),
		Cache:      storage.NewMockCache(),
	}
	_, _, err := ScoreEngines(ctx)
	if err == nil || !strings.Contains(err.Error(), "scorer error") {
		t.Fatalf("expected scorer error, got: %v", err)
	}
}

func TestScoreEnginesSuccess(t *testing.T) {
	origGet := hardwareInfoGet
	origScorer := engineScorer
	t.Cleanup(func() {
		hardwareInfoGet = origGet
		engineScorer = origScorer
	})
	hardwareInfoGet = func(bool) (*types.HwInfo, []string, error) {
		return &types.HwInfo{}, []string{"a warning"}, nil
	}
	want := []engines.ScoredManifest{{Manifest: engines.Manifest{Name: "mock-engine"}}}
	engineScorer = func(*types.HwInfo, []engines.Manifest) ([]engines.ScoredManifest, error) {
		return want, nil
	}

	enginesDir := t.TempDir()
	writeFile(t, filepath.Join(enginesDir, "mock-engine", "engine.yaml"),
		"name: mock-engine\nruntime: rt\n")

	ctx := &Context{
		EnginesDir: enginesDir,
		Config:     storage.NewMockConfig(),
		Cache:      storage.NewMockCache(),
	}
	scored, warnings, err := ScoreEngines(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(scored) != 1 || scored[0].Name != "mock-engine" {
		t.Errorf("unexpected scored engines: %v", scored)
	}
	if len(warnings) != 1 || warnings[0] != "a warning" {
		t.Errorf("unexpected warnings: %v", warnings)
	}
}

// ---- ScoreEnginesWithSpinner ----

func TestScoreEnginesWithSpinnerSuccess(t *testing.T) {
	origGet := hardwareInfoGet
	origScorer := engineScorer
	t.Cleanup(func() {
		hardwareInfoGet = origGet
		engineScorer = origScorer
	})
	hardwareInfoGet = func(bool) (*types.HwInfo, []string, error) {
		return &types.HwInfo{}, nil, nil
	}
	engineScorer = func(*types.HwInfo, []engines.Manifest) ([]engines.ScoredManifest, error) {
		return []engines.ScoredManifest{}, nil
	}

	ctx := &Context{
		EnginesDir: t.TempDir(),
		Config:     storage.NewMockConfig(),
		Cache:      storage.NewMockCache(),
	}
	_, err := ScoreEnginesWithSpinner(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestScoreEnginesWithSpinnerVerboseWarnings(t *testing.T) {
	origGet := hardwareInfoGet
	origScorer := engineScorer
	t.Cleanup(func() {
		hardwareInfoGet = origGet
		engineScorer = origScorer
	})
	hardwareInfoGet = func(bool) (*types.HwInfo, []string, error) {
		return &types.HwInfo{}, []string{"warning1", "warning2"}, nil
	}
	engineScorer = func(*types.HwInfo, []engines.Manifest) ([]engines.ScoredManifest, error) {
		return []engines.ScoredManifest{}, nil
	}

	ctx := &Context{
		EnginesDir: t.TempDir(),
		Config:     storage.NewMockConfig(),
		Cache:      storage.NewMockCache(),
		Verbose:    true,
	}
	_, err := ScoreEnginesWithSpinner(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
