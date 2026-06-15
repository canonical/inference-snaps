package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/canonical/inference-snaps-cli/v2/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/canonical/inference-snaps-cli/v2/pkg/models"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
)

type listModelsCommand struct {
	*common.Context

	// flags
	format string
}

type outputModels struct {
	ActiveModel string                `json:"active-model"`
	Models      []common.ModelDetails `json:"models"`
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

	var modelsList outputModels
	for _, model := range engineManifest.Model.Options {
		modelManifest, err := models.LoadManifest(cmd.ModelsDir, model)
		if err != nil {
			return fmt.Errorf("loading model manifest for model %s: %v", model, err)
		}
		outputModel, err := common.NewModelDetails(modelManifest)
		if err != nil {
			return fmt.Errorf("creating model details for model %s: %v", model, err)
		}
		modelsList.Models = append(modelsList.Models, outputModel)
	}

	activeModel, err := cmd.Cache.GetActiveModel()
	if err != nil {
		return fmt.Errorf("%s: %w", common.LookingUpActiveModel, err)
	}
	modelsList.ActiveModel = activeModel

	switch cmd.format {
	case "table", "":
		if err := cmd.printModelsTable(modelsList); err != nil {
			return fmt.Errorf("table: %w", err)
		}
	case "json":
		if err := cmd.printModelsJson(modelsList); err != nil {
			return fmt.Errorf("json: %w", err)
		}
	default:
		return fmt.Errorf("unknown format %q", cmd.format)
	}

	return nil
}

func (cmd *listModelsCommand) printModelsJson(modelsList outputModels) error {
	jsonString, err := json.MarshalIndent(modelsList, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling models: %v", err)
	}
	fmt.Printf("%s\n", jsonString)
	return nil
}

func (cmd *listModelsCommand) getModelsTable(modelsList outputModels) (string, error) {
	var headerRow = []string{"name", "capabilities", "disk size"}
	tableRows := [][]string{headerRow}

	var modelNameMaxLen, modelCapabilitiesMaxLen int

	for _, model := range modelsList.Models {
		name := model.Name
		// Mark active model with "*"
		if model.ID == modelsList.ActiveModel {
			name = name + "*"
		}

		capabilities := strings.Join(model.Capabilities, ", ")
		diskSize := model.DiskSize

		// Find max name and capabilities lengths
		modelNameMaxLen = max(modelNameMaxLen, len(name), len(headerRow[0]))
		modelCapabilitiesMaxLen = max(modelCapabilitiesMaxLen, len(capabilities), len(headerRow[1]))

		row := []string{name, capabilities, diskSize}
		tableRows = append(tableRows, row)
	}

	tableMaxWidth := 80
	// Increase column widths to account for paddings
	modelNameMaxLen += 1
	modelCapabilitiesMaxLen += 2
	// Disk size column fills the remaining space
	modelDiskSizeMaxLen := tableMaxWidth - (modelNameMaxLen + modelCapabilitiesMaxLen)

	options := []tablewriter.Option{
		tablewriter.WithRenderer(renderer.NewColorized(renderer.ColorizedConfig{
			Header: renderer.Tint{
				FG: renderer.Colors{color.Bold},
			},
			Column: renderer.Tint{
				FG: renderer.Colors{color.Reset},
				BG: renderer.Colors{color.Reset},
			},
			Borders: tw.BorderNone,
			Settings: tw.Settings{
				Separators: tw.Separators{ShowHeader: tw.Off, ShowFooter: tw.Off, BetweenRows: tw.Off, BetweenColumns: tw.Off},
				Lines: tw.Lines{
					ShowTop:        tw.Off,
					ShowBottom:     tw.Off,
					ShowHeaderLine: tw.Off,
					ShowFooterLine: tw.Off,
				},
				CompactMode: tw.On,
			},
		})),
		tablewriter.WithConfig(tablewriter.Config{
			MaxWidth: tableMaxWidth,
			Widths: tw.CellWidth{
				PerColumn: tw.Mapper[int, int]{
					0: modelNameMaxLen,         // Model name
					1: modelCapabilitiesMaxLen, // Capabilities
					2: modelDiskSizeMaxLen,     // Disk size
				},
			},
			Header: tw.CellConfig{
				Alignment: tw.CellAlignment{Global: tw.AlignLeft},
				Padding: tw.CellPadding{
					PerColumn: []tw.Padding{
						{Overwrite: true, Right: " "},
						{Overwrite: true, Left: " ", Right: " "},
						{Overwrite: true, Left: " "},
					},
				},
			},
			Row: tw.CellConfig{
				Formatting: tw.CellFormatting{AutoWrap: tw.WrapTruncate},
				Alignment:  tw.CellAlignment{Global: tw.AlignLeft},
				Padding: tw.CellPadding{
					PerColumn: []tw.Padding{
						{Overwrite: true, Right: " "},
						{Overwrite: true, Left: " ", Right: " "},
						{Overwrite: true, Left: " "},
					},
				},
			},
		}),
	}

	var tableOutput bytes.Buffer
	table := tablewriter.NewTable(&tableOutput, options...)
	table.Header(tableRows[0])
	err := table.Bulk(tableRows[1:])
	if err != nil {
		return "", fmt.Errorf("adding data: %v", err)
	}
	err = table.Render()
	if err != nil {
		return "", fmt.Errorf("rendering: %v", err)
	}
	return tableOutput.String(), nil
}

func (cmd *listModelsCommand) printModelsTable(modelsList outputModels) error {
	if len(modelsList.Models) == 0 {
		fmt.Fprintln(os.Stderr, "No models found.")
		return nil
	}

	tableOutput, err := cmd.getModelsTable(modelsList)
	if err != nil {
		return fmt.Errorf("generating table: %v", err)
	}

	fmt.Print(tableOutput)
	return nil
}
