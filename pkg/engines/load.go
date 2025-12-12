package engines

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func LoadManifests(manifestsDir string) ([]Manifest, error) {
	var manifests []Manifest

	// Sanitize dir path
	if !strings.HasSuffix(manifestsDir, "/") {
		manifestsDir += "/"
	}

	// Iterate engines
	files, err := os.ReadDir(manifestsDir)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", manifestsDir, err)
	}

	for _, file := range files {
		// Engines dir should contain a dir per engine
		if !file.IsDir() {
			continue
		}

		fileName := manifestsDir + file.Name() + "/engine.yaml"
		data, err := os.ReadFile(fileName)
		if err != nil {
			return nil, fmt.Errorf("%s: %s", fileName, err)
		}

		var manifest Manifest
		err = yaml.Unmarshal(data, &manifest)
		if err != nil {
			return nil, fmt.Errorf("%s: %s", manifestsDir, err)
		}

		manifests = append(manifests, manifest)
	}
	return manifests, nil
}

var ErrManifestNotFound = errors.New("engine manifest not found")

func LoadManifest(manifestsDir, engineName string) (*Manifest, error) {
	// Sanitize dir path
	if !strings.HasSuffix(manifestsDir, "/") {
		manifestsDir += "/"
	}

	fileName := manifestsDir + engineName + "/engine.yaml"
	data, err := os.ReadFile(fileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("%w: %s", ErrManifestNotFound, err)
		}
		return nil, fmt.Errorf("%s: %s", fileName, err)
	}

	var manifest Manifest
	err = yaml.Unmarshal(data, &manifest)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", manifestsDir, err)
	}

	return &manifest, nil
}
