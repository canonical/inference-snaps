package commands

import (
	"errors"
	"fmt"
	"slices"

	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/common"
	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/canonical/inference-snaps-cli/v2/pkg/models"
	"github.com/canonical/inference-snaps-cli/v2/pkg/runtimes"
	"github.com/canonical/inference-snaps-cli/v2/pkg/selector"
	"github.com/canonical/inference-snaps-cli/v2/pkg/utils"
	"github.com/spf13/cobra"
)

type useEngineCommand struct {
	*common.Context

	// flags
	auto               bool
	fix                bool
	fallback           string
	assumeYes          bool
	noRestart          bool
	considerComponents bool
}

func UseEngine(ctx *common.Context) *cobra.Command {
	var cmd useEngineCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:   "use-engine [<engine>]",
		Short: "Select an engine",
		// Args
		// modelctl use-engine <engine> requires 1 argument
		// modelctl use-engine --auto does not support any arguments
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: cmd.validateArgs,
		RunE:              cmd.run,
	}

	// flags
	cobraCmd.Flags().BoolVar(&cmd.auto, "auto", false, "automatically select a compatible engine")
	cobraCmd.Flags().BoolVar(&cmd.fix, "fix", false, "fix issues with the currently active engine")
	cobraCmd.Flags().StringVar(&cmd.fallback, "fallback", "", "fallback engine to use when hardware information is unavailable (requires --auto or --fix)")
	cobraCmd.Flags().BoolVar(&cmd.assumeYes, "assume-yes", false, "assume yes for all prompts")
	cobraCmd.Flags().BoolVar(&cmd.noRestart, "no-restart", false, "do not restart the snap after changing engine")
	cobraCmd.Flags().BoolVar(&cmd.considerComponents, "components", false, "consider pre-installed components (requires --auto)")

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

	if cmd.fallback != "" && !cmd.auto && !cmd.fix {
		return fmt.Errorf("--fallback must be used together with --auto or --fix")
	}

	if cmd.considerComponents && !cmd.auto {
		return fmt.Errorf("--components must be used together with --auto")
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

		if err := cmd.migrateConfig(); err != nil {
			return err
		}

		err := cmd.fixActiveEngine()
		if errors.Is(err, common.ErrNoActiveEngine) { // If no engine is active, there's nothing to fix
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
	observable, err := cmd.Snap.HardwareObservable()
	if err != nil {
		return fmt.Errorf("checking hardware observability: %v", err)
	}

	if !observable && cmd.fallback != "" {
		fmt.Printf("Hardware information is unavailable; falling back to engine %q.\n", cmd.fallback)
		return cmd.switchEngine(cmd.fallback)
	}

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

	if cmd.considerComponents {
		ok, err := selectEngineForSeededComponents(cmd, scoredEngines)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
		// otherwise continue with standard auto selection
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
// - auto selects an engine if the active engine no longer exists
// - verify that the active model is supported by the active engine, otherwise switches to the default model
// - if engine exists, make sure it is correctly installed and configured
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

func (cmd *useEngineCommand) migrateConfig() error {
	if err := cmd.Config.Migrate(); err != nil {
		return fmt.Errorf("migrating config: %v", err)
	}
	return nil
}

// engineNames extracts the Name field from a slice of engine manifests using the provided accessor.
func engineNames[T any](items []T, getName func(T) string) []string {
	names := make([]string, len(items))
	for i, item := range items {
		names[i] = getName(item)
	}
	return names
}

// selectEngineForSeededComponents looks at components that are currently installed,
// tries to match these to an engine and model that are compatible, and switches to it.
func selectEngineForSeededComponents(cmd *useEngineCommand, scoredEngines []engines.ScoredManifest) (bool, error) {
	fmt.Println("Checking preinstalled components to influence engine and model selection")

	allEngines, err := engines.LoadManifests(cmd.EnginesDir)
	if err != nil {
		return false, err
	}
	allRuntimes, err := runtimes.LoadManifests(cmd.RuntimesDir)
	if err != nil {
		return false, err
	}
	allModels, err := models.LoadManifests(cmd.ModelsDir)
	if err != nil {
		return false, err
	}

	installedComponents, err := common.InstalledComponents()
	if err != nil {
		return false, err
	}

	fmt.Printf("Installed components: %v\n", installedComponents)

	var seededRuntimes []string
	var seededModels []string

	// A runtime or model is considered seeded if any of its components are currently installed

	for _, runtime := range allRuntimes {
		for _, component := range runtime.Components {
			if slices.Contains(installedComponents, component) {
				seededRuntimes = append(seededRuntimes, runtime.Name)
				break
			}
		}
	}
	fmt.Printf("Seeded runtimes: %v\n", seededRuntimes)

	for _, model := range allModels {
		for _, component := range model.Components {
			if slices.Contains(installedComponents, component) {
				seededModels = append(seededModels, model.ID)
				break
			}
		}
	}
	fmt.Printf("Seeded models: %v\n", seededModels)

	// Check which engines have a seeded runtime and/or a seeded model
	var fullySeededEngines []engines.Manifest
	var partiallySeededEngines []engines.Manifest
	for _, engine := range allEngines {
		hasSeededRuntime := slices.Contains(seededRuntimes, engine.Runtime)
		hasSeededModel := false
		for _, option := range engine.Model.Options {
			if slices.Contains(seededModels, option) {
				hasSeededModel = true
				break
			}
		}
		if hasSeededRuntime && hasSeededModel {
			fullySeededEngines = append(fullySeededEngines, engine)
		} else if hasSeededRuntime || hasSeededModel {
			partiallySeededEngines = append(partiallySeededEngines, engine)
		}
	}

	fmt.Printf("Partially seeded engines: %v\n",
		engineNames(partiallySeededEngines,
			func(e engines.Manifest) string {
				return e.Name
			},
		),
	)

	fmt.Printf("Fully seeded engines: %v\n",
		engineNames(fullySeededEngines,
			func(e engines.Manifest) string {
				return e.Name
			},
		),
	)

	// filterCompatible returns the subset of engines that are compatible
	filterCompatible := func(candidates []engines.Manifest) []engines.ScoredManifest {
		var compatible []engines.ScoredManifest
		for _, candidate := range candidates {
			for _, scoredEngine := range scoredEngines {
				if scoredEngine.Name == candidate.Name && scoredEngine.Score > 0 {
					compatible = append(compatible, scoredEngine)
					break
				}
			}
		}
		return compatible
	}

	// Prefer engines with both a seeded runtime and a seeded model
	compatibleSeededEngines := filterCompatible(fullySeededEngines)
	if len(compatibleSeededEngines) == 0 {
		compatibleSeededEngines = filterCompatible(partiallySeededEngines)
	}
	fmt.Printf("Compatible seeded engines: %v\n", engineNames(compatibleSeededEngines, func(e engines.ScoredManifest) string { return e.Name }))

	// Seeded components do not target any compatible engine
	if len(compatibleSeededEngines) == 0 {
		fmt.Printf("No compatible seeded engines found; falling back to standard auto selection.\n")
		return false, nil
	}

	// If multiple seeded engines, find the top one
	topEngine, err := selector.TopEngine(compatibleSeededEngines)
	if err != nil {
		return false, fmt.Errorf("finding top engine: %v", err)
	}
	fmt.Printf("Top engine: %v\n", topEngine.Name)

	// If multiple models were seeded, prefer the engine's default
	seededModelForEngine := ""
	for _, option := range topEngine.Model.Options {
		if slices.Contains(seededModels, option) {
			seededModelForEngine = option
			if option == topEngine.Model.Default {
				break
			}
		}
	}

	// If a model was seeded, switch to it. Otherwise, switchEngine() will use the default model.
	if seededModelForEngine != "" {
		fmt.Printf("Seeded model for engine: %v\n", seededModelForEngine)
		err = cmd.Cache.SetActiveModel(seededModelForEngine)
		if err != nil {
			return false, fmt.Errorf("setting active model: %v", err)
		}
	}

	err = cmd.switchEngine(topEngine.Name)
	if err != nil {
		return false, fmt.Errorf("switching engine: %v", err)
	}

	return true, nil
}
