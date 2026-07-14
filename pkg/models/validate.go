package models

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"slices"

	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/canonical/inference-snaps-cli/v2/pkg/utils"
	"go.yaml.in/yaml/v4"
)

// diskSizePattern matches a human-readable size such as "6G" or "512M".
var diskSizePattern = regexp.MustCompile(`^\d+(\.\d+)?(K|M|G|T)i?B?$`)

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

	// Get model name from the directory name
	modelName := modelNameFromPath(manifestFilePath)

	manifest, err := parseManifest(yamlData)
	if err != nil {
		return err
	}

	if err := manifest.validate(modelName); err != nil {
		return err
	}

	// Cross-check the referenced components against snapcraft.yaml when it can
	// be located relative to the manifest. snapRoot is the package directory
	// containing the models/ folder: <snapRoot>/models/<model-name>/model.yaml.
	snapRoot := filepath.Dir(filepath.Dir(filepath.Dir(manifestFilePath)))
	knownComponents, err := engines.SnapcraftComponents(snapRoot)
	if err != nil {
		return err
	}
	return engines.ValidateComponents(manifest.Components, knownComponents)
}

func modelNameFromPath(manifestFilePath string) string {
	parts := utils.SplitPathIntoDirectories(manifestFilePath)
	if len(parts) < 2 {
		return ""
	}
	return parts[len(parts)-2] // second last part: model-name/model.yaml
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

func (manifest Manifest) validate(expectedModelName string) error {
	if manifest.Name == "" {
		return fmt.Errorf("required field is not set: name")
	}

	if expectedModelName != "" {
		if manifest.Name != expectedModelName {
			return fmt.Errorf("model directory name should match name in manifest: %s != %s", expectedModelName, manifest.Name)
		}
	}

	if manifest.Description == "" {
		return fmt.Errorf("required field is not set: description")
	}

	if manifest.ModelCardUrl != "" {
		if u, err := url.ParseRequestURI(manifest.ModelCardUrl); err != nil || u.Scheme == "" || u.Host == "" {
			return fmt.Errorf("invalid model-card-url: %s", manifest.ModelCardUrl)
		}
	}

	for _, cap := range manifest.Capabilities {
		if !slices.Contains(SupportedCapabilities(), cap) {
			return fmt.Errorf("unsupported capability: %q", cap)
		}
	}

	if manifest.DiskSize == "" {
		return fmt.Errorf("required field is not set: disk-size")
	}
	if !diskSizePattern.MatchString(manifest.DiskSize) {
		return fmt.Errorf("invalid disk-size format: %s", manifest.DiskSize)
	}

	// components are optional; when set, names must be non-empty. Cross-checking
	// against snapcraft.yaml is done by Validate when the file is available.
	if err := engines.ValidateComponents(manifest.Components, nil); err != nil {
		return err
	}

	// environment is optional; when set, validate the NAME=value syntax.
	if err := engines.ValidateEnvironment(manifest.Environment); err != nil {
		return err
	}

	if err := engines.ValidateLayout(manifest.Layout); err != nil {
		return err
	}

	return nil
}
