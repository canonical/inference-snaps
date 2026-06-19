package common

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"strings"

	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/canonical/inference-snaps-cli/v2/pkg/hardware_info"
	"github.com/canonical/inference-snaps-cli/v2/pkg/models"
	"github.com/canonical/inference-snaps-cli/v2/pkg/runtimes"
	"github.com/canonical/inference-snaps-cli/v2/pkg/selector"
	"github.com/canonical/inference-snaps-cli/v2/pkg/storage"
	"github.com/canonical/inference-snaps-cli/v2/pkg/types"
	"github.com/canonical/inference-snaps-cli/v2/pkg/utils"
)

type Settings struct {
	Environment    []string                `yaml:"environment"`
	Layout         map[string]types.Layout `yaml:"layout"`
	expandedLayout map[string]types.Layout
}

func EngineSettings(ctx *Context) (*Settings, error) {
	activeEngineName, err := ctx.Cache.GetActiveEngine()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", LookingUpActiveEngine, err)
	}

	if activeEngineName == "" {
		return nil, ErrNoActiveEngine
	}

	engineManifest, err := engines.LoadManifest(ctx.EnginesDir, activeEngineName)
	if err != nil {
		return nil, fmt.Errorf("loading engine manifest: %w", err)
	}

	// Load runtime settings
	runtimeManifest, err := runtimes.LoadManifest(ctx.RuntimesDir, engineManifest.Runtime)
	if err != nil {
		return nil, fmt.Errorf("loading runtime manifest: %w", err)
	}

	// Load active model settings
	activeModelId, err := ctx.Cache.GetActiveModel()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", LookingUpActiveModel, err)
	}
	if activeModelId == "" {
		return nil, ErrNoActiveModel
	}

	modelManifest, err := models.LoadManifest(ctx.ModelsDir, activeModelId)
	if err != nil {
		return nil, fmt.Errorf("loading model manifest: %w", err)
	}

	var engineSettings Settings
	engineSettings.Environment = append(engineSettings.Environment, runtimeManifest.Environment...)
	engineSettings.Environment = append(engineSettings.Environment, modelManifest.Environment...)
	engineSettings.Layout = make(map[string]types.Layout)
	maps.Copy(engineSettings.Layout, runtimeManifest.Layout)
	maps.Copy(engineSettings.Layout, modelManifest.Layout)

	return &engineSettings, nil
}

func loadEngineEnvironmentFromSettings(settings *Settings) error {

	for _, kv := range settings.Environment {
		// Split into key/value
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid env var %q", kv)
		}
		k, v := parts[0], parts[1]

		// Expand all env vars in value
		v = os.ExpandEnv(v)

		err := os.Setenv(k, v)
		if err != nil {
			return fmt.Errorf("setting env var %q: %w", k, err)
		}
	}

	settings.expandedLayout = make(map[string]types.Layout, len(settings.Layout))
	for k, v := range settings.Layout {
		engineLayout := types.Layout{
			Symlink: os.ExpandEnv(v.Symlink),
		}
		settings.expandedLayout[os.ExpandEnv(k)] = engineLayout
	}

	for layoutPath, layout := range settings.expandedLayout {
		if layout.Symlink != "" {
			if err := utils.CreateTempSymlink(layout.Symlink, layoutPath); err != nil {
				return fmt.Errorf("creating tmp symlink %s -> %s: %w", layoutPath, layout.Symlink, err)
			}
		}
	}

	return nil
}

func unloadEngineEnvironmentFromSettings(settings *Settings) error {
	// remove the symlinks created for the engine components
	var errs []error
	for layoutPath := range settings.expandedLayout {
		if _, err := utils.RemoveTempSymlink(layoutPath); err != nil {
			errs = append(errs, fmt.Errorf("removing symlink %q: %w", layoutPath, err))
		}
	}
	return errors.Join(errs...)
}

// LoadEngineEnvironment sets env vars of the active engine's components for the current process
// and creates any necessary symlinks
func LoadEngineEnvironment(ctx *Context) (func(), error) {
	engineSettings, err := EngineSettings(ctx)
	if err != nil {
		return nil, fmt.Errorf("error loading engine component settings: %w", err)
	}

	if err = loadEngineEnvironmentFromSettings(engineSettings); err != nil {
		// Clean up any symlinks that were partially created before the failure
		if cleanErr := unloadEngineEnvironmentFromSettings(engineSettings); cleanErr != nil && ctx.Verbose {
			fmt.Fprintf(os.Stderr, "Warning: failed to unload engine environment after load error: %v\n", cleanErr)
		}
		return nil, err
	}

	return func() {
		if err := unloadEngineEnvironmentFromSettings(engineSettings); err != nil {
			if ctx.Verbose {
				fmt.Fprintf(os.Stderr, "Warning: failed to unload engine environment: %v\n", err)
			}
		}
	}, nil
}

// SetEngineConfig sets configurations of the given engine.
// It does not unset previous engine configurations.
func SetEngineConfig(engine *engines.Manifest, ctx *Context) error {
	for confKey, confVal := range engine.Configurations {
		err := ctx.Config.SetDocument(confKey, confVal, storage.EngineConfig)
		if err != nil {
			return fmt.Errorf("setting engine configuration %q: %w", confKey, err)
		}
	}
	return nil
}

func UnsetEngineConfig(engineName string, unsetUserOverrides bool, ctx *Context) error {
	// Unset all engine configurations
	err := ctx.Config.Unset(".", storage.EngineConfig)
	if err != nil {
		return fmt.Errorf("un-setting engine configurations: %w", err)
	}

	if unsetUserOverrides {
		engine, err := engines.LoadManifest(ctx.EnginesDir, engineName)
		if err != nil {
			if errors.Is(err, engines.ErrManifestNotFound) {
				// TODO: remove this when implementing per-engine configuration
				// We can't know what user overrides were set if the manifest is missing
				if ctx.Verbose {
					fmt.Fprintf(os.Stderr, "Warning: previously active engine %q not found; skipping user configuration cleanup.\n", engineName)
				}
				return nil
			}
			return fmt.Errorf("loading engine manifest: %w", err)
		}
		// Unset any user overrides
		for k := range engine.Configurations {
			err = ctx.Config.Unset(k, storage.UserConfig)
			if err != nil {
				return fmt.Errorf("un-setting configuration %q: %w", k, err)
			}
		}
	}

	return nil
}

// hardwareInfoGet and engineScorer are package-level variables so tests can
// inject fakes without changing any production behaviour.
var (
	hardwareInfoGet = hardware_info.Get
	engineScorer    = selector.ScoreEngines
)

/*
ScoreEngines loads all engine manifests, looks up the host machine information,
and scores the engines according to their compatibility with the host.

Warning: calls to this function can block for a number of seconds while the host machine information is being looked up.
*/
func ScoreEngines(ctx *Context) ([]engines.ScoredManifest, []string, error) {
	allEngines, err := engines.LoadManifests(ctx.EnginesDir)
	if err != nil {
		return nil, nil, fmt.Errorf("loading engines: %w", err)
	}

	machineInfo, warnings, err := hardwareInfoGet(false)
	if err != nil {
		return nil, nil, fmt.Errorf("getting machine info: %w", err)
	}

	scoredEngines, err := engineScorer(machineInfo, allEngines)
	if err != nil {
		return nil, nil, fmt.Errorf("scoring engines: %w", err)
	}

	return scoredEngines, warnings, nil
}

// ScoreEnginesWithSpinner is same as ScoreEngines but with a progress spinner.
// It prints the warnings to stderr.
func ScoreEnginesWithSpinner(ctx *Context) ([]engines.ScoredManifest, error) {
	stopProgress := StartProgressSpinner("Checking engine compatibility")
	scoredEngines, warnings, err := ScoreEngines(ctx)
	stopProgress()

	if len(warnings) > 0 && ctx.Verbose {
		for _, warning := range warnings {
			fmt.Fprintf(os.Stderr, "Warning: %s\n", warning)
		}
	}

	return scoredEngines, err
}
