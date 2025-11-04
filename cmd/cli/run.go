package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/canonical/inference-snaps-cli/pkg/selector"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func addRunCommand() {
	cmd := &cobra.Command{
		Use:    "run <path>",
		Short:  "Run a subprocess",
		Hidden: true,
		Args:   cobra.MaximumNArgs(1),
		RunE:   run,
	}
	rootCmd.AddCommand(cmd)
}

func run(_ *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("unexpected number of arguments, expected 1 got %d", len(args))
	}

	// TODO
	// --wait-for-components

	err := loadEngineEnvironment()
	if err != nil {
		return fmt.Errorf("error loading engine environment: %v", err)
	}

	path := args[0]

	cmd := exec.Command(path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func loadEngineEnvironment() error {
	activeEngineName, err := cache.GetActiveEngine()
	if err != nil {
		return fmt.Errorf("error looking up active engine: %v", err)
	}

	if activeEngineName == "" {
		return nil
	}

	manifest, err := selector.LoadManifestFromDir(enginesDir, activeEngineName)
	if err != nil {
		return fmt.Errorf("error loading active engine manifest: %v", err)
	}

	componentsDir, found := os.LookupEnv("SNAP_COMPONENTS")
	if !found {
		return fmt.Errorf("SNAP_COMPONENTS env var not set")
	}

	type comp struct {
		Environment []string `yaml:"environment"`
	}

	for _, componentName := range manifest.Components {
		componentPath := filepath.Join(componentsDir, componentName)
		componentYamlFile := filepath.Join(componentPath, "component.yaml")

		data, err := os.ReadFile(componentYamlFile)
		if err != nil {
			return fmt.Errorf("error reading %s: %v", componentYamlFile, err)
		}

		var component comp
		err = yaml.Unmarshal(data, &component)
		if err != nil {
			return fmt.Errorf("error unmarshaling %s: %v", componentYamlFile, err)
		}

		for i := range component.Environment {
			// Split into key/value
			kv := component.Environment[i]
			parts := strings.SplitN(kv, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid env var %q", kv)
			}
			k, v := parts[0], parts[1]

			// Set component path env var for expansion
			if err := os.Setenv(envComponent, componentPath); err != nil {
				return fmt.Errorf("error setting %q: %v", envComponent, err)
			}

			// Expand all env vars in value
			v = os.ExpandEnv(v)

			// Unset the component path
			if err := os.Unsetenv(envComponent); err != nil {
				return fmt.Errorf("error unsetting %q: %v", envComponent, err)
			}

			err = os.Setenv(k, v)
			if err != nil {
				return fmt.Errorf("error setting %q: %v", k, err)
			}
		}

	}

	return nil
}
