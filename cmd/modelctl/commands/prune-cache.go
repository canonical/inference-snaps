package commands

import (
	"fmt"
	"os"
	"slices"

	"github.com/canonical/go-snapctl"
	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/common"
	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/canonical/inference-snaps-cli/v2/pkg/models"
	"github.com/canonical/inference-snaps-cli/v2/pkg/runtimes"
	"github.com/canonical/inference-snaps-cli/v2/pkg/snap_store"
	"github.com/canonical/inference-snaps-cli/v2/pkg/utils"
	"github.com/spf13/cobra"
)

type pruneCacheCommand struct {
	*common.Context

	// flags
	engine string
}

func PruneCache(ctx *common.Context) *cobra.Command {
	var cmd pruneCacheCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:   "prune-cache",
		Short: "Remove cached data",
		RunE:  cmd.run,
	}

	// flags
	cobraCmd.Flags().StringVar(&cmd.engine, "engine", "", "Remove caches of an engine")

	return cobraCmd
}

func (cmd *pruneCacheCommand) run(_ *cobra.Command, _ []string) error {
	if !utils.IsRootUser() {
		return common.ErrPermissionDenied
	}

	var err error
	var componentsToRemove []string

	if cmd.engine == "" {
		componentsToRemove, err = cmd.unusedComponentsAll()
		if err != nil {
			return fmt.Errorf("finding all unused components: %w", err)
		}
	} else {
		componentsToRemove, err = cmd.unusedComponentsEngine(cmd.engine)
		if err != nil {
			return fmt.Errorf("finding unused engine components: %w", err)
		}
	}

	if len(componentsToRemove) == 0 {
		fmt.Println("No components to remove.")
		return nil
	}

	if confirmed, err := cmd.printComponentsAndConfirm(componentsToRemove); err != nil {
		return err
	} else if !confirmed {
		return nil
	}
	return snapctl.RemoveComponents(componentsToRemove...).Run()
}

func (cmd *pruneCacheCommand) unusedComponentsAll() ([]string, error) {
	allInstalledComponents, err := common.InstalledComponents()
	if err != nil {
		return nil, fmt.Errorf("getting installed components: %w", err)
	}

	requiredComponents, err := common.ComponentsRequiredByCurrentSelection(cmd.Context)
	if err != nil {
		return nil, fmt.Errorf("getting required components: %w", err)
	}

	var unusedComponents []string

	for _, component := range allInstalledComponents {
		if !slices.Contains(requiredComponents, component) {
			unusedComponents = append(unusedComponents, component)
		}
	}

	return unusedComponents, nil
}

func (cmd *pruneCacheCommand) unusedComponentsEngine(engineName string) ([]string, error) {
	activeEngine, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", common.LookingUpActiveEngine, err)
	}

	if activeEngine == engineName {
		return nil, fmt.Errorf("cannot prune active engine")
	}

	requiredComponents, err := common.ComponentsRequiredByCurrentSelection(cmd.Context)
	if err != nil {
		return nil, fmt.Errorf("getting required components: %w", err)
	}

	var engineComponents []string

	engineManifest, err := engines.LoadManifest(cmd.EnginesDir, engineName)
	if err != nil {
		return nil, fmt.Errorf("loading engine manifest: %w", err)
	}

	// Include runtime components if engine has a runtime
	if engineManifest.Runtime != "" {
		runtimeManifest, err := runtimes.LoadManifest(cmd.RuntimesDir, engineManifest.Runtime)
		if err != nil {
			return nil, fmt.Errorf("loading runtimes manifest: %w", err)
		}
		engineComponents = append(engineComponents, runtimeManifest.Components...)
	}

	// Include model components of all models compatible with this engine
	for _, modelId := range engineManifest.Model.Options {
		modelManifest, err := models.LoadManifest(cmd.ModelsDir, modelId)
		if err != nil {
			return nil, fmt.Errorf("loading model manifest for model %q: %w", modelId, err)
		}
		engineComponents = append(engineComponents, modelManifest.Components...)
	}

	var unusedComponents []string

	for _, engineComponent := range engineComponents {
		if !slices.Contains(requiredComponents, engineComponent) && // only remove if not required by current active engine and model
			!slices.Contains(unusedComponents, engineComponent) { // prevent same components from being listed multiple times
			componentInstalled, err := common.ComponentInstalled(engineComponent)
			if err != nil {
				return nil, err
			}
			if componentInstalled {
				unusedComponents = append(unusedComponents, engineComponent)
			}
		}
	}

	return unusedComponents, nil
}

func (cmd *pruneCacheCommand) printComponentsAndConfirm(componentsToRemove []string) (bool, error) {

	componentSizes, err := snap_store.ComponentSizes()
	if err != nil && cmd.Verbose {
		_, _ = fmt.Fprintf(os.Stderr, "Warning: unable to query component sizes: %v\n", err)
	}

	fmt.Println("Removing components:")
	for _, componentName := range componentsToRemove {
		componentLine := componentName
		if size, ok := componentSizes[componentName]; ok {
			componentLine += fmt.Sprintf(" (%s)", utils.FmtBytes(uint64(size)))
		}

		fmt.Printf("- %s\n", componentLine)
	}

	if utils.IsTerminalOutput() {
		fmt.Println()
		if !common.PromptYN("Continue removing components?", false) {
			fmt.Println("Cancelled. No changes applied.")
			return false, nil
		}
	}

	return true, nil
}
