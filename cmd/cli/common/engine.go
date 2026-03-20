package common

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/hardware_info"
	"github.com/canonical/inference-snaps-cli/pkg/selector"
	"github.com/canonical/inference-snaps-cli/pkg/storage"
	"gopkg.in/yaml.v3"
)

const (
	componentEnv    = "COMPONENT"
	ProgressScoring = "Checking engines"
)

type ComponentSettings struct {
	componentName string
	Servers       map[string]map[string]string `yaml:"servers"`
	Environment   []string                     `yaml:"environment"`
}

func EngineComponentSettings(ctx *Context) ([]ComponentSettings, error) {
	activeEngineName, err := ctx.Cache.GetActiveEngine()
	if err != nil {
		return nil, fmt.Errorf("error looking up active engine: %v", err)
	}

	if activeEngineName == "" {
		return nil, fmt.Errorf("no active engine")
	}

	manifest, err := engines.LoadManifest(ctx.EnginesDir, activeEngineName)
	if err != nil {
		return nil, fmt.Errorf("error loading engine manifest: %v", err)
	}

	componentsDir, found := os.LookupEnv("SNAP_COMPONENTS")
	if !found {
		return nil, fmt.Errorf("SNAP_COMPONENTS env var not set")
	}

	var settingsCollection []ComponentSettings
	for _, componentName := range manifest.Components {
		componentPath := filepath.Join(componentsDir, componentName)
		componentYamlFile := filepath.Join(componentPath, "component.yaml")

		data, err := os.ReadFile(componentYamlFile)
		if err != nil {
			return nil, fmt.Errorf("error reading %s: %v", componentYamlFile, err)
		}

		var settings ComponentSettings
		err = yaml.Unmarshal(data, &settings)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling %s: %v", componentYamlFile, err)
		}

		settings.componentName = componentName

		settingsCollection = append(settingsCollection, settings)
	}

	return settingsCollection, nil
}

// LoadEngineEnvironment sets env vars of the active engine's components for the current process
func LoadEngineEnvironment(ctx *Context) error {
	settingsCollection, err := EngineComponentSettings(ctx)
	if err != nil {
		return fmt.Errorf("error loading engine component settings: %v", err)
	}

	componentsDir, found := os.LookupEnv("SNAP_COMPONENTS")
	if !found {
		return fmt.Errorf("SNAP_COMPONENTS env var not set")
	}

	for _, settings := range settingsCollection {
		for i := range settings.Environment {
			// Split into key/value
			kv := settings.Environment[i]
			parts := strings.SplitN(kv, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid env var %q", kv)
			}
			k, v := parts[0], parts[1]

			// Set component path env var for expansion
			componentPath := filepath.Join(componentsDir, settings.componentName)
			if err := os.Setenv(componentEnv, componentPath); err != nil {
				return fmt.Errorf("error setting %q: %v", componentEnv, err)
			}

			// Expand all env vars in value
			v = os.ExpandEnv(v)

			// Unset the component path
			if err := os.Unsetenv(componentEnv); err != nil {
				return fmt.Errorf("error unsetting %q: %v", componentEnv, err)
			}

			err = os.Setenv(k, v)
			if err != nil {
				return fmt.Errorf("error setting %q: %v", k, err)
			}
		}

	}

	return nil
}

// SetEngineConfig sets configurations of the given engine.
// It does not unset previous engine configurations.
func SetEngineConfig(engine *engines.Manifest, ctx *Context) error {
	for confKey, confVal := range engine.Configurations {
		err := ctx.Config.SetDocument(confKey, confVal, storage.EngineConfig)
		if err != nil {
			return fmt.Errorf("error setting engine configuration %q: %v", confKey, err)
		}
	}
	return nil
}

func UnsetEngineConfig(engineName string, unsetUserOverrides bool, ctx *Context) error {
	// Unset all engine configurations
	err := ctx.Config.Unset(".", storage.EngineConfig)
	if err != nil {
		return fmt.Errorf("error un-setting engine configurations: %v", err)
	}

	if unsetUserOverrides {
		engine, err := engines.LoadManifest(ctx.EnginesDir, engineName)
		if err != nil {
			if errors.Is(err, engines.ErrManifestNotFound) {
				// TODO: remove this when implementing per-engine configuration
				// We can't know what user overrides were set if the manifest is missing
				fmt.Fprintf(os.Stderr, "Warning: previously active engine %q not found; skipping user configuration cleanup.\n", engineName)
				return nil
			}
			return fmt.Errorf("error loading engine manifest: %v", err)
		} else {
			// Unset any user overrides
			for k := range engine.Configurations {
				err = ctx.Config.Unset(k, storage.UserConfig)
				if err != nil {
					return fmt.Errorf("error un-setting configuration %q: %v", k, err)
				}
			}
		}
	}

	return nil
}

/*
ScoreEngines loads all engine manifests, looks up the host machine information,
and scores the engines according to their compatibility with the host.

Warning: calls to this function can block for a number of seconds while the host machine information is being looked up.
*/
func ScoreEngines(ctx *Context) ([]engines.ScoredManifest, error) {
	allEngines, err := engines.LoadManifests(ctx.EnginesDir)
	if err != nil {
		return nil, fmt.Errorf("error loading engines: %v", err)
	}

	machineInfo, err := hardware_info.Get(false)
	if err != nil {
		return nil, fmt.Errorf("error getting machine info: %v", err)
	}

	scoredEngines, err := selector.ScoreEngines(machineInfo, allEngines)
	if err != nil {
		return nil, fmt.Errorf("error scoring engines: %v", err)
	}

	return scoredEngines, nil
}
