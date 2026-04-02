package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/canonical/go-snapctl/env"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
)

func ComponentInstalled(component string) (bool, error) {
	// Check in /snap/$SNAP_INSTANCE_NAME/components/$SNAP_REVISION if component is mounted
	directoryPath := fmt.Sprintf("/snap/%s/components/%s/%s", env.SnapInstanceName(), env.SnapRevision(), component)

	info, err := os.Stat(directoryPath)

	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, fmt.Errorf("error checking component directory %q: %v", component, err)
		}
	} else {
		if info.IsDir() {
			return true, nil
		} else {
			return false, fmt.Errorf("component %q exists but is not a directory", component)
		}
	}
}

func WaitForComponents(ctx *Context) error {
	const maxWait = 3600 // seconds
	const interval = 10  // seconds
	activeEngineName, err := ctx.Cache.GetActiveEngine()
	if err != nil {
		return fmt.Errorf("error looking up active engine: %v", err)
	}

	if activeEngineName == "" {
		return fmt.Errorf("no active engine")
	}

	manifest, err := engines.LoadManifest(ctx.EnginesDir, activeEngineName)
	if err != nil {
		return fmt.Errorf("error loading engine manifest: %v", err)
	}

	missing, err := checkMissingComponents(manifest)
	if err != nil {
		return err
	}

	for elapsed := 0; elapsed < maxWait && len(missing) > 0; elapsed += interval {
		fmt.Fprintf(os.Stderr, "Waiting for required snap components: %s (%d/%ds)\n",
			strings.Join(missing, ", "), elapsed, maxWait)

		time.Sleep(interval * time.Second)

		missing, err = checkMissingComponents(manifest)
		if err != nil {
			return err
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("timed out after %ds while waiting for required components: %s",
			maxWait, strings.Join(missing, ", "))
	}

	return nil
}

// TODO: unify with similar code in use-engine.go
func checkMissingComponents(manifest *engines.Manifest) ([]string, error) {
	componentsDir, found := os.LookupEnv("SNAP_COMPONENTS")
	if !found {
		return nil, fmt.Errorf("SNAP_COMPONENTS env var not set")
	}

	var missing []string
	for _, component := range manifest.Components {
		componentPath := filepath.Join(componentsDir, component)
		if _, err := os.Stat(componentPath); os.IsNotExist(err) {
			missing = append(missing, component)
		}
	}

	return missing, nil
}
