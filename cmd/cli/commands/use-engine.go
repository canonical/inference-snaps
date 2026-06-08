package commands

import (
	"errors"
	"fmt"
	"slices"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/models"
	"github.com/canonical/inference-snaps-cli/pkg/selector"
	"github.com/canonical/inference-snaps-cli/pkg/utils"
	"github.com/spf13/cobra"
)

type useEngineCommand struct {
	*common.Context

	// flags
	auto      bool
	fix       bool
	assumeYes bool
	noRestart bool
}

func UseEngine(ctx *common.Context) *cobra.Command {
	var cmd useEngineCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:   "use-engine [<engine>]",
		Short: "Select an engine",
		// Args
		// cli use-engine <engine> requires 1 argument
		// cli use-engine --auto does not support any arguments
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: cmd.validateArgs,
		RunE:              cmd.run,
	}

	// flags
	cobraCmd.Flags().BoolVar(&cmd.auto, "auto", false, "automatically select a compatible engine")
	cobraCmd.Flags().BoolVar(&cmd.fix, "fix", false, "fix issues with the currently active engine")
	cobraCmd.Flags().BoolVar(&cmd.assumeYes, "assume-yes", false, "assume yes for all prompts")
	cobraCmd.Flags().BoolVar(&cmd.noRestart, "no-restart", false, "do not restart the snap after changing engine")

	return cobraCmd
}

func (cmd *useEngineCommand) validateArgs(_ *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
	manifests, err := engines.LoadManifests(cmd.EnginesDir)
	if err != nil {
		fmt.Printf("Error loading engines: %v\n", err)
		return nil, cobra.ShellCompDirectiveError
	}

	var engineNames []cobra.Completion
	for i := range manifests {
		engineNames = append(engineNames, manifests[i].Name)
	}

	return engineNames, cobra.ShellCompDirectiveNoSpace
}

func (cmd *useEngineCommand) run(_ *cobra.Command, args []string) error {
	if !utils.IsRootUser() {
		return common.ErrPermissionDenied
	}

	if cmd.auto {
		if len(args) != 0 {
			return fmt.Errorf("cannot specify both engine name and --auto flag")
		}
		return cmd.autoSelectEngine()
	} else if cmd.fix {
		if len(args) != 0 {
			return fmt.Errorf("cannot specify both engine name and --fix flag")
		}
		// If no engine is active, there's nothing to fix.
		err := cmd.fixActiveEngine()
		if errors.Is(err, common.ErrNoActiveEngine) {
			return nil
		}
		return err
	} else {
		if len(args) == 1 {
			return cmd.switchEngine(args[0])
		} else {
			return fmt.Errorf("engine name not specified")
		}
	}
}

func (cmd *useEngineCommand) autoSelectEngine() error {
	scoredEngines, err := common.ScoreEnginesWithSpinner(cmd.Context)
	if err != nil {
		return fmt.Errorf("scoring engines: %v", err)
	}

	return cmd.autoSelectScoredEngine(scoredEngines)
}

func (cmd *useEngineCommand) autoSelectScoredEngine(scoredEngines []engines.ScoredManifest) error {

	fmt.Println("Evaluating engines for optimal hardware compatibility:")
	for _, engine := range scoredEngines {
		if engine.Score == 0 {
			fmt.Printf("✘ %s: not compatible\n", engine.Name)

			// Only print incompatibility reasons if verbose flag is set
			if cmd.Verbose {
				reasons := cmd.verboseIncompatibilityReasons(engine.CompatibilityReport)
				for _, reason := range reasons {
					fmt.Printf("  - %s\n", reason)
				}
			}
		} else if engine.IsExperimental() {
			fmt.Printf("• %s: experimental, score=%d\n", engine.Name, engine.Score)
		} else {
			fmt.Printf("✔ %s: compatible, score=%d\n", engine.Name, engine.Score)
		}
	}

	selectedEngine, err := selector.TopEngine(scoredEngines)
	if err != nil {
		return fmt.Errorf("finding top engine: %v", err)
	}

	fmt.Printf("Selected engine: %s\n", selectedEngine.Name)

	err = cmd.switchEngine(selectedEngine.Name)
	if err != nil {
		return fmt.Errorf("use engine: %s", err)
	}

	return nil
}

// switchEngine changes the engine that is used by the snap
func (cmd *useEngineCommand) switchEngine(engineName string) error {

	newEngineManifest, err := engines.LoadManifest(cmd.EnginesDir, engineName)
	if err != nil {
		if errors.Is(err, engines.ErrManifestNotFound) {
			if cmd.Verbose {
				fmt.Println(err)
			}
			return fmt.Errorf("%q not found", engineName)
		}
		return fmt.Errorf("loading engine manifest: %v", err)
	}

	// We need to check which components are required for the switch.
	// If the current model is supported by the new engine, we use the active model's components.
	// If the model is not supported, we need to use the components of the new engine's default model.

	activeModelName, err := cmd.Cache.GetActiveModel()
	if err != nil {
		return fmt.Errorf("getting active model name: %v", err)
	}

	// If the current active model is not supported by the new engine, switch to the engine's default model
	newModelName := activeModelName
	if !slices.Contains(newEngineManifest.Model.Options, activeModelName) {
		newModelName = newEngineManifest.Model.Default
	}

	var newModelManifest *models.Manifest
	if newModelName != "" {
		newModelManifest, err = models.LoadManifest(cmd.ModelsDir, newModelName)
		if err != nil {
			return fmt.Errorf("loading model manifest: %v", err)
		}
	}

	// Check for missing components
	cancelledByUser, err := common.InstallMissingComponents(cmd.Context, cmd.assumeYes, newEngineManifest, newModelManifest)
	if err != nil {
		return fmt.Errorf("installing missing components: %v", err)
	}
	if cancelledByUser {
		return nil
	}

	err = cmd.Cache.SetActiveModel(newModelName)
	if err != nil {
		return fmt.Errorf("setting active model: %v", err)
	}

	activeEngineName, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		return fmt.Errorf("%s: %w", common.LookingUpActiveEngine, err)
	}

	if activeEngineName == engineName {
		// Engine not changed, nothing left to do
		return nil
	}

	// Unset active engine's configurations
	if activeEngineName != "" {
		err = common.UnsetEngineConfig(activeEngineName, true, cmd.Context)
		if err != nil {
			return fmt.Errorf("un-setting engine configurations: %v", err)
		}
	}

	if err = cmd.Cache.SetActiveEngine(newEngineManifest.Name); err != nil {
		return fmt.Errorf("setting active engine: %v", err)
	}

	if err = common.SetEngineConfig(newEngineManifest, cmd.Context); err != nil {
		return fmt.Errorf("setting new engine configurations: %v", err)
	}

	fmt.Printf("Engine changed to %q.\n", engineName)

	// Ask if the user wants to restart
	if !cmd.noRestart {
		return common.PromptRestartToApplyChanges(cmd.Context, cmd.assumeYes)
	}

	return nil
}

// fixActiveEngine does the following:
// 1. auto selects an engine if the active engine no longer exists
// 2. verify that the active model is supported by the active engine, otherwise switches to the default model
// 2. if engine exists, make sure it is correctly installed and configured
func (cmd *useEngineCommand) fixActiveEngine() error {
	activeEngineName, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		return fmt.Errorf("%s: %w", common.LookingUpActiveEngine, err)
	}
	if activeEngineName == "" {
		return common.ErrNoActiveEngine
	}

	// If active engine no longer exists, auto select another one
	engineManifest, err := engines.LoadManifest(cmd.EnginesDir, activeEngineName)
	if errors.Is(err, engines.ErrManifestNotFound) {
		fmt.Printf("Active engine %q not found, performing auto selection instead.\n", activeEngineName)
		return cmd.autoSelectEngine()
	} else if err != nil {
		return fmt.Errorf("loading active engine manifest: %v", err)
	}

	// Check if the model is supported, otherwise switch to the default
	activeModelId, err := cmd.Cache.GetActiveModel()
	if err != nil {
		return fmt.Errorf("%s: %w", common.LookingUpActiveModel, err)
	}
	if !slices.Contains(engineManifest.Model.Options, activeModelId) {
		activeModelId = engineManifest.Model.Default
	}
	err = cmd.Cache.SetActiveModel(activeModelId)
	if err != nil {
		return fmt.Errorf("setting active model: %v", err)
	}

	var modelManifest *models.Manifest
	if activeModelId != "" {
		modelManifest, err = models.LoadManifest(cmd.ModelsDir, activeModelId)
		if err != nil {
			return fmt.Errorf("loading active model manifest: %v", err)
		}
	}

	// Make sure all components are correctly installed and engine is configured
	if _, err = common.InstallMissingComponents(cmd.Context, cmd.assumeYes, engineManifest, modelManifest); err != nil {
		return fmt.Errorf("installing missing components: %v", err)
	}

	if err = common.UnsetEngineConfig(activeEngineName, false, cmd.Context); err != nil {
		return fmt.Errorf("un-setting engine configurations: %v", err)
	}
	if err = common.SetEngineConfig(engineManifest, cmd.Context); err != nil {
		return fmt.Errorf("setting engine configurations: %v", err)
	}

	return nil
}

func (cmd *useEngineCommand) verboseIncompatibilityReasons(report engines.CompatibilityReport) []string {
	var reasons []string
	if !report.CompatibleMemory {
		reasons = append(reasons, fmt.Sprintf("requires %s memory, has %s (RAM + swap)", utils.FmtBytes(report.RequiredMemory), utils.FmtBytes(report.TotalRAM+report.TotalSwap)))
	}
	if !report.CompatibleDisk {
		reasons = append(reasons, fmt.Sprintf("requires %s disk space, has %s", utils.FmtBytes(report.RequiredDiskSpace), utils.FmtBytes(report.AvailableDiskSpace)))
	}
	if !report.CompatibleDevices {
		reasons = append(reasons, "required device not found")
	}
	return reasons
}
