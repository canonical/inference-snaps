package commands

import (
	"fmt"
	"strings"

	"github.com/canonical/inference-snaps-cli/v2/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/spf13/cobra"
)

type listModelsCommand struct {
	*common.Context

	// flags
	format string
}

func ListModels(ctx *common.Context) *cobra.Command {
	var cmd listModelsCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:               "list-models",
		Short:             "List available models",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.run,
	}

	// flags
	supportedFormats := []string{"table", "json"}
	cobraCmd.Flags().StringVar(
		&cmd.format,
		"format",
		"table",
		fmt.Sprintf("output format (%s)", strings.Join(supportedFormats, ", ")),
	)

	return cobraCmd
}

func (cmd *listModelsCommand) run(_ *cobra.Command, _ []string) error {
	activeEngine, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		return fmt.Errorf("%s: %w", common.LookingUpActiveEngine, err)
	}

	engineManifest, err := engines.LoadManifest(cmd.EnginesDir, activeEngine)
	if err != nil {
		return fmt.Errorf("%s: %w", common.LoadingEngineManifest, err)
	}

	for _, model := range engineManifest.Model.Options {
		// TODO in IENG-2398: handle json and table formats, including model metadata
		// For now just print a list of model IDs
		fmt.Println(model)
	}

	return nil
}
