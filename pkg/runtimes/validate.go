package runtimes

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/canonical/inference-snaps-cli/v2/pkg/utils"
	"go.yaml.in/yaml/v4"
)

func Validate(manifestFilePath string) error {

	if filepath.Base(manifestFilePath) != ManifestFilename {
		return fmt.Errorf("manifest file must be called %s: %s", ManifestFilename, manifestFilePath)
	}

	_, err := os.Stat(manifestFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("manifest file does not exist: %s", manifestFilePath)
		}
		return fmt.Errorf("getting file info: %v", err)
	}

	yamlData, err := os.ReadFile(manifestFilePath)
	if err != nil {
		return fmt.Errorf("reading file: %v", err)
	}

	// Get runtime name from path
	runtimeName := runtimeNameFromPath(manifestFilePath)

	manifest, err := parseManifest(yamlData)
	if err != nil {
		return err
	}

	if err := manifest.validate(runtimeName); err != nil {
		return err
	}

	// Cross-check the referenced components against snapcraft.yaml when it can
	// be located relative to the manifest. snapRoot is the package directory
	// containing the runtimes/ folder: <snapRoot>/runtimes/<name>/runtime.yaml.
	snapRoot := filepath.Dir(filepath.Dir(filepath.Dir(manifestFilePath)))
	knownComponents, err := engines.SnapcraftComponents(snapRoot)
	if err != nil {
		return err
	}
	return engines.ValidateComponents(manifest.Components, knownComponents)
}

func runtimeNameFromPath(manifestFilePath string) string {
	parts := utils.SplitPathIntoDirectories(manifestFilePath)
	if len(parts) < 2 {
		return ""
	}
	return parts[len(parts)-2] // second last part: runtime-name/runtime.yaml
}

func parseManifest(yamlData []byte) (Manifest, error) {
	yamlData = bytes.TrimSpace(yamlData)
	if len(yamlData) == 0 {
		return Manifest{}, errors.New("empty yaml data")
	}

	var manifest Manifest

	yamlDecoder := yaml.NewDecoder(bytes.NewReader(yamlData))

	// Error if there are unknown fields in the yaml
	yamlDecoder.KnownFields(true)

	// We depend on the yaml unmarshal to check field types
	if err := yamlDecoder.Decode(&manifest); err != nil {
		return Manifest{}, fmt.Errorf("decoding manifest: %v", err)
	}

	return manifest, nil
}

func (manifest Manifest) validate(expectedRuntimeName string) error {
	if manifest.Name == "" {
		return fmt.Errorf("required field is not set: name")
	}

	// Only do runtime name matching test if expected name is set
	if expectedRuntimeName != "" {
		if manifest.Name != expectedRuntimeName {
			return fmt.Errorf("runtime directory name should match name in manifest: %s != %s", expectedRuntimeName, manifest.Name)
		}
	}

	if len(manifest.Servers) == 0 {
		return fmt.Errorf("required field is not set: servers")
	}

	for name, server := range manifest.Servers {
		if err := server.validate(name); err != nil {
			return err
		}
	}

	// environment is optional; when set, validate the NAME=value syntax.
	if err := engines.ValidateEnvironment(manifest.Environment); err != nil {
		return err
	}

	// components are optional; when set, names must be non-empty. Cross-checking
	// against snapcraft.yaml is done by Validate when the file is available.
	if err := engines.ValidateComponents(manifest.Components, nil); err != nil {
		return err
	}

	if err := engines.ValidateLayout(manifest.Layout); err != nil {
		return err
	}

	return nil
}

func (server Server) validate(name string) error {
	if server.Protocol == "" {
		return fmt.Errorf("required field is not set for server %s: protocol", name)
	}

	if server.BasePath == "" {
		return fmt.Errorf("required field is not set for server %s: base-path", name)
	}

	return nil
}
