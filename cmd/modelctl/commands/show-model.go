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

type showModelCommand struct {
	*common.Context

	// flags
	format string
}

func ShowModel(ctx *common.Context) *cobra.Command {
	var cmd showModelCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:   "show-model [<model>]",
		Short: "Print information about a model",
		Long:  "Print information about the active model, or the specified model",
		// Args
		// modelctl show-model <model> requires 0 or 1 argument
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: cmd.validateArgs,
		RunE:              cmd.run,
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

func (cmd *showModelCommand) run(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.showCurrentModel()
	} else if len(args) == 1 {
		return cmd.showModel(args[0])
	} else {
		return fmt.Errorf("invalid number of arguments")
	}
}

// validateArgs returns a list of model names supported by the currently active engine
func (cmd *showModelCommand) validateArgs(_ *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
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

func (cmd *showModelCommand) showCurrentModel() error {
	currentModel, err := cmd.Cache.GetActiveModel()
	if err != nil {
		return fmt.Errorf("%s: %w", common.LookingUpActiveModel, err)
	}
	if currentModel == "" {
		return common.ErrNoActiveModel
	}
	return cmd.showModel(currentModel)
}

func (cmd *showModelCommand) showModel(modelNameOrID string) error {
	modelManifest, err := common.GetModelByNameOrId(cmd.Context, modelNameOrID)
	if err != nil {
		return err
	}

	err = cmd.printModelManifest(modelManifest)
	if err != nil {
		return fmt.Errorf("printing model manifest: %v", err)
	}
	return nil
}

func (cmd *showModelCommand) printModelManifest(manifest *models.Manifest) error {
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
