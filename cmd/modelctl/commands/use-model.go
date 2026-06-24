package commands

import (
	"fmt"
	"slices"

	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/common"
	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/canonical/inference-snaps-cli/v2/pkg/models"
	"github.com/canonical/inference-snaps-cli/v2/pkg/utils"
	"github.com/spf13/cobra"
)

type useModelCommand struct {
	*common.Context

	// flags
	assumeYes bool
	noRestart bool
}

func UseModel(ctx *common.Context) *cobra.Command {
	var cmd useModelCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:   "use-model [<model>]",
		Short: "Select a model",
		// Args
		// modelctl use-model <model> requires 1 argument
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: cmd.validateArgs,
		RunE:              cmd.run,
	}

	// flags
	cobraCmd.Flags().BoolVar(&cmd.assumeYes, "assume-yes", false, "assume yes for all prompts")
	cobraCmd.Flags().BoolVar(&cmd.noRestart, "no-restart", false, "do not restart the snap after changing model")

	return cobraCmd
}

// validateArgs returns a list of model names supported by the currently active engine
func (cmd *useModelCommand) validateArgs(_ *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
	activeEngine, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	if activeEngine == "" {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	engineManifest, err := engines.LoadManifest(cmd.EnginesDir, activeEngine)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	supportedModels := engineManifest.Model.Options

	modelManifests, err := models.LoadManifests(cmd.ModelsDir)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var completions []cobra.Completion
	for _, manifest := range modelManifests {
		if slices.Contains(supportedModels, manifest.ID) {
			completions = append(completions, manifest.Name)
		}
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

func (cmd *useModelCommand) run(_ *cobra.Command, args []string) error {
	if !utils.IsRootUser() {
		return common.ErrPermissionDenied
	}

	if len(args) == 1 {
		return cmd.switchModel(args[0])
	} else {
		return fmt.Errorf("model name not specified")
	}
}

func (cmd *useModelCommand) switchModel(modelNameOrID string) error {

	modelManifest, err := common.GetModelByNameOrId(cmd.Context, modelNameOrID)
	if err != nil {
		return err
	}

	activeEngine, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		return fmt.Errorf("%s: %w", common.LookingUpActiveEngine, err)
	}

	engineManifest, err := engines.LoadManifest(cmd.EnginesDir, activeEngine)
	if err != nil {
		return fmt.Errorf("%s: %w", "loading engine manifest", err)
	}

	cancelledByUser, err := common.InstallMissingComponents(cmd.Context, cmd.assumeYes, engineManifest, modelManifest)
	if err != nil {
		return fmt.Errorf("installing missing components: %v", err)
	}

	if cancelledByUser {
		return nil
	}

	activeModelId, err := cmd.Cache.GetActiveModel()
	if err != nil {
		return fmt.Errorf("%s: %w", common.LookingUpActiveModel, err)
	}

	if activeModelId == modelManifest.ID {
		// Model not changed, nothing left to do
		return nil
	}

	if err = cmd.Cache.SetActiveModel(modelManifest.ID); err != nil {
		return fmt.Errorf("setting active model: %v", err)
	}

	fmt.Printf("Model changed to %q.\n", modelManifest.Name)

	// Ask if the user wants to restart
	if !cmd.noRestart {
		return common.PromptRestartToApplyChanges(cmd.Context, cmd.assumeYes)
	}

	return nil
}
