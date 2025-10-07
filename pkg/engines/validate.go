package engines

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/canonical/famous-models-cli/pkg/utils"
	"gopkg.in/yaml.v3"
)

func Validate(manifestFilePath string) error {

	if !strings.HasSuffix(manifestFilePath, "engine.yaml") {
		return fmt.Errorf("manifest file must be called engine.yaml: %s", manifestFilePath)
	}

	_, err := os.Stat(manifestFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("manifest file does not exist: %s", manifestFilePath)
		}
		return fmt.Errorf("error getting file info: %v", err)
	}

	yamlData, err := os.ReadFile(manifestFilePath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Get engine name from path
	engineName := engineNameFromPath(manifestFilePath)

	return validateManifestYaml(engineName, yamlData)
}

func engineNameFromPath(manifestFilePath string) string {
	parts := utils.SplitPathIntoDirectories(manifestFilePath)
	if len(parts) < 2 {
		return ""
	}
	return parts[len(parts)-2] // second last part: engine-name/engine.yaml
}

func validateManifestYaml(expectedName string, yamlData []byte) error {
	yamlData = bytes.TrimSpace(yamlData)
	if len(yamlData) == 0 {
		return errors.New("empty yaml data")
	}

	var manifest Manifest

	yamlDecoder := yaml.NewDecoder(bytes.NewReader(yamlData))

	// Error if there are unknown fields in the yaml
	yamlDecoder.KnownFields(true)

	// We depend on the yaml unmarshal to check field types
	if err := yamlDecoder.Decode(&manifest); err != nil {
		return fmt.Errorf("error decoding: %v", err)
	}

	return manifest.validate(expectedName)
}

func (manifest Manifest) validate(expectedEngineName string) error {
	if manifest.Name == "" {
		return fmt.Errorf("required field is not set: name")
	}

	// Only do engine name matching test if expected name is set
	if expectedEngineName != "" {
		if manifest.Name != expectedEngineName {
			return fmt.Errorf("engine directory name should match name in manifest: %s != %s", expectedEngineName, manifest.Name)
		}
	}

	if manifest.Description == "" {
		return fmt.Errorf("required field is not set: description")
	}

	if manifest.Vendor == "" {
		return fmt.Errorf("required field is not set: vendor")
	}

	if manifest.Grade == "" {
		return fmt.Errorf("required field is not set: grade")
	}
	if manifest.Grade != "stable" && manifest.Grade != "devel" {
		return fmt.Errorf("grade should be 'stable' or 'devel'")
	}

	if manifest.Memory != nil {
		_, err := utils.StringToBytes(*manifest.Memory)
		if err != nil {
			return fmt.Errorf("error parsing memory: %v", err)
		}
	}

	if manifest.DiskSpace != nil {
		_, err := utils.StringToBytes(*manifest.DiskSpace)
		if err != nil {
			return fmt.Errorf("error parsing disk space: %v", err)
		}
	}

	for key, val := range manifest.Configurations {
		if !utils.IsPrimitive(val) {
			return fmt.Errorf("configuration field %s is not a primitive value: %v", key, val)
		}
	}

	return manifest.Devices.validate()
}
