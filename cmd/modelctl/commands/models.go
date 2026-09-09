package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/common"
	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
)

type modelsCommand struct {
	*common.Context

	// flags
	format string
	all    bool
}

type outputModels struct {
	ActiveModel string                `json:"active-model"`
	Models      []common.ModelDetails `json:"models"`
}

func Models(ctx *common.Context) *cobra.Command {
	return newModelsCmd(ctx, "models", "")
}

// TODO: remove when we fully migrate to "models" command
func ListModels(ctx *common.Context) *cobra.Command {
	return newModelsCmd(ctx, "list-models", `use "models" instead`)
}

func newModelsCmd(ctx *common.Context, use, deprecated string) *cobra.Command {
	var cmd modelsCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:               use,
		Short:             "List available models",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.run,
		Deprecated:        deprecated,
	}

	// flags
	supportedFormats := []string{"table", "json"}
	cobraCmd.Flags().StringVar(
		&cmd.format,
		"format",
		"table",
		fmt.Sprintf("output format (%s)", strings.Join(supportedFormats, ", ")),
	)
	cobraCmd.Flags().BoolVar(
		&cmd.all,
		"all",
		false,
		"list all models, including those not supported by the active engine",
	)
	return cobraCmd
}

func (cmd *modelsCommand) run(_ *cobra.Command, _ []string) error {
	activeEngine, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		return fmt.Errorf("%s: %w", common.LookingUpActiveEngine, err)
	}

	engineManifest, err := engines.LoadManifest(cmd.EnginesDir, activeEngine)
	if err != nil {
		return fmt.Errorf("%s: %w", common.LoadingEngineManifest, err)
	}

	allModels, err := common.GetAllModels(cmd.Context)
	if err != nil {
		return fmt.Errorf("getting all models: %w", err)
	}

	var modelsList outputModels
	if cmd.all {
		modelsList.Models = allModels
	} else {
		allModelsByName := make(map[string]common.ModelDetails, len(allModels))
		for _, model := range allModels {
			allModelsByName[model.Name] = model
		}

		for _, model := range engineManifest.Model.Options {
			outputModel, ok := allModelsByName[model]
			if !ok {
				return fmt.Errorf("model %q referenced by engine %q does not exist", model, activeEngine)
			}
			modelsList.Models = append(modelsList.Models, outputModel)
		}
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

func (cmd *modelsCommand) printModelsJson(modelsList outputModels) error {
	jsonString, err := json.MarshalIndent(modelsList, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling models: %v", err)
	}
	fmt.Printf("%s\n", jsonString)
	return nil
}

func (cmd *modelsCommand) getModelsTable(modelsList outputModels) (string, error) {
	headerRow := []string{"name", "capabilities", "disk"}
	if cmd.all {
		headerRow = append(headerRow, "engines")
	}
	tableRows := [][]string{headerRow}

	var modelNameMaxLen, modelCapabilitiesMaxLen, modelDiskMaxLen int

	for _, model := range modelsList.Models {
		name := model.Name
		// Mark active model with "*"
		if model.Name == modelsList.ActiveModel {
			name = name + "*"
		}

		capabilities := strings.Join(model.Capabilities, ", ")
		diskSize := model.DiskSize
		var engines string
		if cmd.all {
			engines = strings.Join(model.CompatibleEngines, ", ")
		}
		// Find max name and capabilities lengths
		modelNameMaxLen = max(modelNameMaxLen, len(name), len(headerRow[0]))
		modelCapabilitiesMaxLen = max(modelCapabilitiesMaxLen, len(capabilities), len(headerRow[1]))
		modelDiskMaxLen = max(modelDiskMaxLen, len(diskSize), len(headerRow[2]))

		row := []string{name, capabilities, diskSize}
		if cmd.all {
			row = append(row, engines)
		}
		tableRows = append(tableRows, row)
	}

	var tableMaxWidth int
	if cmd.all {
		tableMaxWidth = 120
	} else {
		tableMaxWidth = 80
	}
	// Increase column widths to account for paddings
	modelNameMaxLen += 1
	modelCapabilitiesMaxLen += 2
	modelDiskMaxLen += 2
	modelEnginesMaxLen := 0
	if cmd.all {
		modelDiskMaxLen += 1
		// Engines column fills the remaining space
		modelEnginesMaxLen = tableMaxWidth - (modelNameMaxLen + modelCapabilitiesMaxLen + modelDiskMaxLen)
	}

	widths := tw.Mapper[int, int]{
		0: modelNameMaxLen,         // Model name
		1: modelCapabilitiesMaxLen, // Capabilities
		2: modelDiskMaxLen,         // Disk
	}
	headerPadding := []tw.Padding{
		{Overwrite: true, Right: " "},
		{Overwrite: true, Left: " ", Right: " "},
		{Overwrite: true, Left: " "},
	}
	rowPadding := []tw.Padding{
		{Overwrite: true, Right: " "},
		{Overwrite: true, Left: " ", Right: " "},
		{Overwrite: true, Left: " "},
	}
	if cmd.all {
		widths[3] = modelEnginesMaxLen // Engines
		headerPadding = append(headerPadding, tw.Padding{Overwrite: true, Left: " "})
		rowPadding = append(rowPadding, tw.Padding{Overwrite: true, Left: " "})
	}

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
				PerColumn: widths,
			},
			Header: tw.CellConfig{
				Alignment: tw.CellAlignment{Global: tw.AlignLeft},
				Padding: tw.CellPadding{
					PerColumn: headerPadding,
				},
			},
			Row: tw.CellConfig{
				Formatting: tw.CellFormatting{AutoWrap: tw.WrapTruncate},
				Alignment:  tw.CellAlignment{Global: tw.AlignLeft},
				Padding: tw.CellPadding{
					PerColumn: rowPadding,
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

	if cmd.all {
		return tableOutput.String(), nil
	}

	activeEngine, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		return "", fmt.Errorf("%s: %w", common.LookingUpActiveEngine, err)
	}

	allModels, err := common.GetAllModels(cmd.Context)
	if err != nil {
		return "", fmt.Errorf("getting all models: %w", err)
	}

	incompatibleModelsCount := len(allModels) - len(tableRows[1:])
	if incompatibleModelsCount != 0 {
		hint := common.SuggestListModels(incompatibleModelsCount, activeEngine)
		return tableOutput.String() + "\n" + hint + "\n", nil
	}
	return tableOutput.String(), nil
}

func (cmd *modelsCommand) printModelsTable(modelsList outputModels) error {
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
