package models

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const ManifestFilename = "model.yaml"

var ErrManifestNotFound = errors.New("engine manifest not found")

func LoadManifests(manifestsDir string) ([]Manifest, error) {
	var manifests []Manifest

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

		fileName := filepath.Join(manifestsDir, file.Name(), ManifestFilename)
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

func LoadManifest(manifestsDir, modelID string) (*Manifest, error) {

	fileName := filepath.Join(manifestsDir, modelID, ManifestFilename)
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
		return nil, fmt.Errorf("%s: %s", fileName, err)
	}

	return &manifest, nil
}
