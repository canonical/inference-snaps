package commands

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/common"
	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/canonical/inference-snaps-cli/v2/pkg/models"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type modelCommand struct {
	*common.Context

	// flags
	format string
}

func Model(ctx *common.Context) *cobra.Command {
	return newModelCmd(ctx, "model [<model>]", "")
}

// TODO: remove when we fully migrate to "model" command
func ShowModel(ctx *common.Context) *cobra.Command {
	return newModelCmd(ctx, "show-model", `use "model" instead`)
}

func newModelCmd(ctx *common.Context, use, deprecated string) *cobra.Command {
	var cmd modelCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:   use,
		Short: "Print information about a model",
		Long:  "Print information about the active model, or the specified model",
		// Args
		// modelctl show-model <model> requires 0 or 1 argument
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: cmd.validateArgs,
		RunE:              cmd.run,
		Deprecated:        deprecated,
	}

	// flags
	supportedFormats := []string{"json", "yaml"}
	cobraCmd.Flags().StringVar(
		&cmd.format,
		"format",
		"yaml",
		fmt.Sprintf("output format (%s)", strings.Join(supportedFormats, ", ")),
	)

	return cobraCmd
}

func (cmd *modelCommand) run(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.showCurrentModel()
	} else if len(args) == 1 {
		return cmd.model(args[0])
	} else {
		return fmt.Errorf("invalid number of arguments")
	}
}

// validateArgs returns a list of model names supported by the currently active engine
func (cmd *modelCommand) validateArgs(_ *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
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
		if slices.Contains(supportedModels, manifest.Name) {
			completions = append(completions, manifest.Name)
		}
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

func (cmd *modelCommand) showCurrentModel() error {
	currentModel, err := cmd.Cache.GetActiveModel()
	if err != nil {
		return fmt.Errorf("%s: %w", common.LookingUpActiveModel, err)
	}
	if currentModel == "" {
		return common.ErrNoActiveModel
	}
	return cmd.model(currentModel)
}

func (cmd *modelCommand) model(modelNameOrAlias string) error {
	modelManifest, err := common.GetModelByNameOrAlias(cmd.Context, modelNameOrAlias)
	if err != nil {
		return err
	}

	err = cmd.printModelManifest(modelManifest)
	if err != nil {
		return fmt.Errorf("printing model manifest: %v", err)
	}
	return nil
}

func (cmd *modelCommand) printModelManifest(manifest *models.Manifest) error {
	output, err := common.NewModelDetails(manifest)
	if err != nil {
		return fmt.Errorf("creating model details: %v", err)
	}

	switch cmd.format {
	case "json":
		jsonString, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			return fmt.Errorf("json: %s", err)
		}
		fmt.Printf("%s\n", jsonString)
	case "yaml", "":
		modelYaml, err := yaml.Marshal(output)
		if err != nil {
			return fmt.Errorf("yaml: %s", err)
		}
		fmt.Print(string(modelYaml))
	default:
		return fmt.Errorf("unknown format %q", cmd.format)
	}

	return nil
}
