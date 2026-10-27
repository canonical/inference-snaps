package engines

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"go.yaml.in/yaml/v4"
)

// envVarNamePattern matches a valid shell environment variable name.
var envVarNamePattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// ValidateEnvironment checks that each entry has the form NAME=value with a
// syntactically valid variable name. The value part may be empty and may
// reference other variables (e.g. "$FOO").
func ValidateEnvironment(environment []string) error {
	for _, entry := range environment {
		name, _, found := strings.Cut(entry, "=")
		if !found {
			return fmt.Errorf("invalid environment entry, expected NAME=value: %q", entry)
		}
		if !envVarNamePattern.MatchString(name) {
			return fmt.Errorf("invalid environment variable name: %q", name)
		}
	}
	return nil
}

// ValidateComponents checks that component names are non-empty and, when
// knownComponents is not nil, that each referenced component is declared
// (e.g. in snapcraft.yaml).
func ValidateComponents(components []string, knownComponents map[string]bool) error {
	for _, name := range components {
		if name == "" {
			return fmt.Errorf("component name is empty")
		}
		if knownComponents != nil && !knownComponents[name] {
			return fmt.Errorf("component not declared in snapcraft.yaml: %s", name)
		}
	}
	return nil
}

// SnapcraftComponents returns the set of component names declared in the
// snapcraft.yaml located under snapRoot. It looks for snap/snapcraft.yaml first
// and falls back to snapcraft.yaml. It returns nil (without error) when no
// snapcraft.yaml is found.
func SnapcraftComponents(snapRoot string) (map[string]bool, error) {
	candidates := []string{
		filepath.Join(snapRoot, "snap", "snapcraft.yaml"),
		filepath.Join(snapRoot, "snapcraft.yaml"),
	}

	var snapcraftPath string
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			snapcraftPath = candidate
			break
		}
	}
	if snapcraftPath == "" {
		return nil, nil
	}

	data, err := os.ReadFile(snapcraftPath)
	if err != nil {
		return nil, fmt.Errorf("reading snapcraft.yaml: %v", err)
	}

	var snapcraft struct {
		Components map[string]yaml.Node `yaml:"components"`
	}
	if err := yaml.Unmarshal(data, &snapcraft); err != nil {
		return nil, fmt.Errorf("parsing snapcraft.yaml: %v", err)
	}

	known := make(map[string]bool, len(snapcraft.Components))
	for name := range snapcraft.Components {
		known[name] = true
	}
	return known, nil
}
