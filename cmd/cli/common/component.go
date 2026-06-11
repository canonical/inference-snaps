package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/canonical/inference-snaps-cli/v2/pkg/models"
	"github.com/canonical/inference-snaps-cli/v2/pkg/runtimes"
	"github.com/canonical/inference-snaps-cli/v2/pkg/snap_store"
	"github.com/canonical/inference-snaps-cli/v2/pkg/utils"
)

// InstalledComponents returns the names of all currently installed components
// by listing subdirectories inside the SNAP_COMPONENTS directory.
func InstalledComponents() ([]string, error) {
	componentsDir, found := os.LookupEnv("SNAP_COMPONENTS")
	if !found {
		return nil, fmt.Errorf("SNAP_COMPONENTS env var not set")
	}

	entries, err := os.ReadDir(componentsDir)
	if err != nil {
		return nil, fmt.Errorf("reading components directory %q: %v", componentsDir, err)
	}

	var installed []string
	for _, entry := range entries {
		// Use os.Stat() as DirEntry.IsDir() does not resolve symlinks
		info, err := os.Stat(filepath.Join(componentsDir, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("getting file info for component %q: %v", entry.Name(), err)
		}
		if info.IsDir() {
			installed = append(installed, entry.Name())
		}
	}

	return installed, nil
}

// ComponentInstalled checks if a specific component is installed
func ComponentInstalled(component string) (bool, error) {
	componentsDir, found := os.LookupEnv("SNAP_COMPONENTS")
	if !found {
		return false, fmt.Errorf("SNAP_COMPONENTS env var not set")
	}

	// Check in $SNAP_COMPONENTS/ if component is mounted
	directoryPath := fmt.Sprintf("%s/%s", componentsDir, component)

	// os.Stat() resolves symlinks, so IsDir() will report true if target is a directory
	info, err := os.Stat(directoryPath)

	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, fmt.Errorf("checking component directory %q: %v", component, err)
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
	const maxWait = 1 * time.Hour
	const interval = 10 * time.Second

	return waitForComponentsWithTimeoutAndInterval(ctx, maxWait, interval)
}

func waitForComponentsWithTimeoutAndInterval(ctx *Context, timeout time.Duration, interval time.Duration) error {
	required, err := ComponentsRequiredByCurrentSelection(ctx)
	if err != nil {
		return fmt.Errorf("determining required components: %v", err)
	}

	missing, err := MissingComponents(required)
	if err != nil {
		return err
	}

	var elapsed time.Duration
	for elapsed = 0; elapsed < timeout && len(missing) > 0; elapsed += interval {
		fmt.Fprintf(os.Stderr, "Waiting for required snap components: %s (%.0f/%.0fs)\n",
			strings.Join(missing, ", "), elapsed.Seconds(), timeout.Seconds())

		time.Sleep(interval)

		missing, err = MissingComponents(required)
		if err != nil {
			return err
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("timeout after waiting %.0fs for required components: %s",
			timeout.Seconds(), strings.Join(missing, ", "))
	}

	return nil
}

func ComponentsRequiredByCurrentSelection(ctx *Context) ([]string, error) {
	var requiredComponents []string

	activeEngineName, err := ctx.Cache.GetActiveEngine()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", LookingUpActiveEngine, err)
	}
	engineManifest, err := engines.LoadManifest(ctx.EnginesDir, activeEngineName)
	if err != nil {
		return nil, fmt.Errorf("loading engine manifest: %v", err)
	}

	// Components required by active engine's runtime
	if engineManifest.Runtime != "" {
		runtimeRequiredComponents, err := ComponentsRequiredByRuntime(ctx, engineManifest.Runtime)
		if err != nil {
			return nil, fmt.Errorf("getting components required by runtime: %v", err)
		}
		requiredComponents = append(requiredComponents, runtimeRequiredComponents...)
	}

	// Components required by active model
	activeModelId, err := ctx.Cache.GetActiveModel()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", LookingUpActiveModel, err)
	}
	if activeModelId != "" {
		modelManifest, err := models.LoadManifest(ctx.ModelsDir, activeModelId)
		if err != nil {
			return nil, fmt.Errorf("loading model manifest: %v", err)
		}
		requiredComponents = append(requiredComponents, modelManifest.Components...)
	}

	return requiredComponents, nil
}

func ComponentsRequiredByRuntime(ctx *Context, runtimeName string) ([]string, error) {
	var requiredComponents []string

	// Components required by runtime
	runtimeManifest, err := runtimes.LoadManifest(ctx.RuntimesDir, runtimeName)
	if err != nil {
		return nil, fmt.Errorf("loading runtime manifest: %v", err)
	}
	requiredComponents = append(requiredComponents, runtimeManifest.Components...)

	return requiredComponents, nil
}

func MissingComponents(requiredComponents []string) ([]string, error) {
	if len(requiredComponents) == 0 {
		return nil, nil
	}

	componentsDir, found := os.LookupEnv("SNAP_COMPONENTS")
	if !found {
		return nil, fmt.Errorf("SNAP_COMPONENTS env var not set")
	}

	var missing []string
	for _, component := range requiredComponents {
		componentPath := filepath.Join(componentsDir, component)
		if _, err := os.Stat(componentPath); os.IsNotExist(err) {
			missing = append(missing, component)
		}
	}

	return missing, nil
}

// InstallMissingComponents determines which components are required by the given engine and model,
// prompts the user if needed, and installs any that are missing.
// It returns cancelledByUser=true if the user declined the installation prompt.
func InstallMissingComponents(ctx *Context, assumeYes bool, engineManifest *engines.Manifest, modelManifest *models.Manifest) (cancelledByUser bool, err error) {
	var requiredComponents []string

	if engineManifest.Runtime != "" {
		runtimeComponents, err := ComponentsRequiredByRuntime(ctx, engineManifest.Runtime)
		if err != nil {
			return false, fmt.Errorf("getting components required by runtime: %v", err)
		}
		requiredComponents = append(requiredComponents, runtimeComponents...)
	}

	if modelManifest != nil {
		requiredComponents = append(requiredComponents, modelManifest.Components...)
	}

	missingComponents, err := MissingComponents(requiredComponents)
	if err != nil {
		return false, fmt.Errorf("checking installed components: %v", err)
	}
	if len(missingComponents) == 0 {
		return false, nil
	}

	componentSizes, err := snap_store.ComponentSizes()
	if err != nil && ctx.Verbose {
		fmt.Fprintf(os.Stderr, "Warning: unable to query component sizes: %v\n", err)
	}

	// Format list of components, adding size if it is known
	fmt.Println("Need to install the following components:")
	for _, componentName := range missingComponents {
		line := fmt.Sprintf("- %s", componentName)
		if size, found := componentSizes[componentName]; found {
			line += fmt.Sprintf(" (%s)", utils.FmtBytes(uint64(size)))
		}
		fmt.Println(line)
	}

	// Only ask for confirmation if it is an interactive terminal
	if !assumeYes && utils.IsTerminalOutput() {
		fmt.Println()
		if !PromptYN("Do you want to continue?", true) {
			fmt.Println("Cancelled. No changes applied.")
			return true, nil
		}
	}

	// Leave a blank line after printing component list and optional confirmation, before printing component installation progress
	fmt.Println()

	// This is blocking, but there is a timeout bug:
	// https://github.com/canonical/inference-snaps-cli/issues/122
	err = InstallComponents(ctx, missingComponents)
	if err != nil {
		return false, fmt.Errorf("installing components: %v", err)
	}

	return false, nil
}

func InstallComponents(ctx *Context, components []string) error {
	return installComponents(ctx, components, 60*time.Minute, 10*time.Second)
}

func installComponents(ctx *Context, components []string, installTimeout time.Duration, retryDelay time.Duration) error {
	const (
		snapdAlreadyInstalledError = "already installed"
		snapdUnknownSnapError      = "cannot install components for a snap that is unknown to the store"
		snapdTimeoutError          = "timeout exceeded while waiting for response"
		snapdChangeInProgressError = "change in progress"
	)
	startTime := time.Now()

	for _, component := range components {
		stopProgress := StartProgressSpinner("Installing " + component)
		err := ctx.Snap.InstallComponent(component)

		for err != nil {
			// Only retry up to the set timeout
			if time.Since(startTime) > installTimeout {
				stopProgress()
				return fmt.Errorf("timed out while installing %q:"+
					"\nMonitor the installation progress with \"snap changes\""+
					"\n\nRerun this command once the installation is complete",
					component)
			}

			if strings.Contains(err.Error(), snapdAlreadyInstalledError) {
				// All good. Continue installing next component.
				break

			} else if strings.Contains(err.Error(), snapdUnknownSnapError) {
				// Install component manually
				stopProgress()
				return fmt.Errorf("snap not known to the store:"+
					"\nRerun this command after manually installing %q",
					component)

			} else if strings.Contains(err.Error(), snapdTimeoutError) {
				// Snapd timed out while installing this component
				time.Sleep(retryDelay)
				err = ctx.Snap.InstallComponent(component)

			} else if strings.Contains(err.Error(), snapdChangeInProgressError) {
				// Snapd is busy with installing this component or busy with an unrelated change
				time.Sleep(retryDelay)
				err = ctx.Snap.InstallComponent(component)

			} else {
				// Any other error we do not specifically handle will stop installing components
				stopProgress()
				return fmt.Errorf("installing %q: %s", component, err)
			}
		}

		stopProgress()
		fmt.Println("Installed " + component)
	}

	return nil
}
