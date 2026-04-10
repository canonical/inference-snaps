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

type ComponentLayout struct {
	Symlink string `yaml:"symlink"`
}

type ComponentSettings struct {
	componentName  string
	Servers        map[string]map[string]string `yaml:"servers"`
	Environment    []string                     `yaml:"environment"`
	Layout         map[string]ComponentLayout   `yaml:"layout"`
	expandedLayout map[string]ComponentLayout
}

func EngineComponentSettings(ctx *Context) ([]ComponentSettings, error) {
	activeEngineName, err := ctx.Cache.GetActiveEngine()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", LookingUpActiveEngine, err)
	}

	if activeEngineName == "" {
		return nil, ErrNoActiveEngine
	}

	manifest, err := engines.LoadManifest(ctx.EnginesDir, activeEngineName)
	if err != nil {
		return nil, fmt.Errorf("loading engine manifest: %v", err)
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
			return nil, fmt.Errorf("reading %s: %v", componentYamlFile, err)
		}

		var settings ComponentSettings
		err = yaml.Unmarshal(data, &settings)
		if err != nil {
			return nil, fmt.Errorf("unmarshaling %s: %v", componentYamlFile, err)
		}

		settings.componentName = componentName

		settingsCollection = append(settingsCollection, settings)
	}

	return settingsCollection, nil
}

func loadEngineEnvironmentFromSettingsCollection(settingsCollection []ComponentSettings) error {
	componentsDir, found := os.LookupEnv("SNAP_COMPONENTS")
	if !found {
		return fmt.Errorf("SNAP_COMPONENTS env var not set")
	}
	for i, settings := range settingsCollection {
		// Set component path env var for expansion
		componentPath := filepath.Join(componentsDir, settings.componentName)
		if err := os.Setenv(componentEnv, componentPath); err != nil {
			return fmt.Errorf("setting env %q: %v", componentEnv, err)
		}

		for i := range settings.Environment {
			// Split into key/value
			kv := settings.Environment[i]
			parts := strings.SplitN(kv, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid env var %q", kv)
			}
			k, v := parts[0], parts[1]

			// Expand all env vars in value
			v = os.ExpandEnv(v)

			err := os.Setenv(k, v)
			if err != nil {
				return fmt.Errorf("setting env var %q: %v", k, err)
			}
		}

		settingsCollection[i].expandedLayout = make(map[string]ComponentLayout, len(settings.Layout))
		for k, v := range settings.Layout {
			ComponentLayout := ComponentLayout{
				Symlink: os.ExpandEnv(v.Symlink),
			}
			settingsCollection[i].expandedLayout[os.ExpandEnv(k)] = ComponentLayout
		}
		for layoutPath, layout := range settingsCollection[i].expandedLayout {
			if layout.Symlink != "" {
				// Assigning variables for better readability
				target := layout.Symlink
				link := layoutPath
				if !strings.HasPrefix(link, "/tmp") {
					fmt.Fprintf(os.Stderr, "Warning: skipping symlink %q for component %q: path is not in /tmp\n", link, settings.componentName)
					continue
				}

				// Ensure the parent directory for the link exists
				if err := os.MkdirAll(filepath.Dir(link), 0755); err != nil {
					return fmt.Errorf("error creating directory for symlink %q: %v", link, err)
				}

				// Remove existing file, symlink, or directory if it exists
				if err := os.RemoveAll(link); err != nil {
					return fmt.Errorf("error removing existing path %q: %v", link, err)
				}
				err := os.Symlink(target, link)
				if err != nil {
					return fmt.Errorf("error creating symlink from %q to %q: %v", link, target, err)
				}
			}
		}
	}
	if err := os.Unsetenv(componentEnv); err != nil {
		return fmt.Errorf("error unsetting %q: %v", componentEnv, err)
	}
	return nil
}

func unloadEngineEnvironmentFromSettingsCollection(settingsCollection []ComponentSettings) error {

	// remove the symlinks created for the engine components
	for _, settings := range settingsCollection {
		for layoutPath := range settings.expandedLayout {
			if err := os.RemoveAll(layoutPath); err != nil {
				return fmt.Errorf("removing symlink %q: %v", layoutPath, err)
			}
		}
	}

	return nil
}

// LoadEngineEnvironment sets env vars of the active engine's components for the current process
// and creates any necessary symlinks
func LoadEngineEnvironment(ctx *Context) (func(), error) {
	settingsCollection, err := EngineComponentSettings(ctx)
	if err != nil {
		return nil, fmt.Errorf("error loading engine component settings: %v", err)
	}
	err = loadEngineEnvironmentFromSettingsCollection(settingsCollection)
	return func() {
		if err := unloadEngineEnvironmentFromSettingsCollection(settingsCollection); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to unload engine environment: %v\n", err)
		}
		return
	}, err
}

// SetEngineConfig sets configurations of the given engine.
// It does not unset previous engine configurations.
func SetEngineConfig(engine *engines.Manifest, ctx *Context) error {
	for confKey, confVal := range engine.Configurations {
		err := ctx.Config.SetDocument(confKey, confVal, storage.EngineConfig)
		if err != nil {
			return fmt.Errorf("setting engine configuration %q: %v", confKey, err)
		}
	}
	return nil
}

func UnsetEngineConfig(engineName string, unsetUserOverrides bool, ctx *Context) error {
	// Unset all engine configurations
	err := ctx.Config.Unset(".", storage.EngineConfig)
	if err != nil {
		return fmt.Errorf("un-setting engine configurations: %v", err)
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
			return fmt.Errorf("loading engine manifest: %v", err)
		} else {
			// Unset any user overrides
			for k := range engine.Configurations {
				err = ctx.Config.Unset(k, storage.UserConfig)
				if err != nil {
					return fmt.Errorf("un-setting configuration %q: %v", k, err)
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
		return nil, fmt.Errorf("loading engines: %v", err)
	}

	machineInfo, err := hardware_info.Get(false)
	if err != nil {
		return nil, fmt.Errorf("getting machine info: %v", err)
	}

	scoredEngines, err := selector.ScoreEngines(machineInfo, allEngines)
	if err != nil {
		return nil, fmt.Errorf("scoring engines: %v", err)
	}

	return scoredEngines, nil
}
