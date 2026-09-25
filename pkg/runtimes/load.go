package runtimes

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const ManifestFilename = "runtime.yaml"

func LoadManifests(manifestsDir string) ([]Manifest, error) {
	var manifests []Manifest

	// Iterate runtimes
	runtimeDirs, err := os.ReadDir(manifestsDir)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", manifestsDir, err)
	}

	for _, runtimeDir := range runtimeDirs {
		// Runtimes dir should contain a dir per runtime
		if !runtimeDir.IsDir() {
			continue
		}

		fileName := filepath.Join(manifestsDir, runtimeDir.Name(), ManifestFilename)
		data, err := os.ReadFile(fileName)
		if err != nil {
			return nil, fmt.Errorf("%s: %s", fileName, err)
		}

		var manifest Manifest
		err = yaml.Unmarshal(data, &manifest)
		if err != nil {
			return nil, fmt.Errorf("malformed manifest %s: %v", fileName, err)
		}

		if err := manifest.validate(runtimeDir.Name()); err != nil {
			return nil, fmt.Errorf("invalid manifest %s: %v", fileName, err)
		}

		manifests = append(manifests, manifest)
	}
	return manifests, nil
}

var ErrManifestNotFound = errors.New("runtime manifest not found")

func LoadManifest(manifestsDir, runtimeName string) (*Manifest, error) {

	fileName := filepath.Join(manifestsDir, runtimeName, ManifestFilename)
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
		return nil, fmt.Errorf("malformed manifest %s: %v", fileName, err)
	}

	if err := manifest.validate(runtimeName); err != nil {
		return nil, fmt.Errorf("invalid manifest %s: %v", fileName, err)
	}

	return &manifest, nil
}
