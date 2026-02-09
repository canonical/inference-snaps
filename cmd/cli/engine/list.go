package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
)

type listCommand struct {
	*common.Context

	// flags
	format string
}

type outputEngines struct {
	ActiveEngine string                   `json:"active-engine"`
	Engines      []engines.ScoredManifest `json:"engines"`
}

func ListCommand(ctx *common.Context) *cobra.Command {
	var cmd listCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:               "list-engines",
		Short:             "List available engines",
		GroupID:           groupID,
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

func (cmd *listCommand) run(_ *cobra.Command, _ []string) error {
	scoredEngines, err := scoreEngines(cmd.Context)
	if err != nil {
		return fmt.Errorf("error scoring engines: %v", err)
	}

	activeEngine, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		return fmt.Errorf("could not determine active engine: %v", err)
	}

	enginesList := outputEngines{
		ActiveEngine: activeEngine,
		Engines:      scoredEngines,
	}

	switch cmd.format {
	case "table":
		err = cmd.printEnginesTable(enginesList)
		if err != nil {
			return fmt.Errorf("error printing list: %v", err)
		}
	case "json":
		err = cmd.printEnginesJson(enginesList)
		if err != nil {
			return fmt.Errorf("error printing list: %v", err)
		}
	default:
		return fmt.Errorf("unknown format %q", cmd.format)
	}

	return nil
}

func (cmd *listCommand) printEnginesJson(enginesList outputEngines) error {
	jsonString, err := json.MarshalIndent(enginesList, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling engines to json: %v", err)
	}
	fmt.Printf("%s\n", jsonString)
	return nil
}

func (cmd *listCommand) printEnginesTable(enginesList outputEngines) error {
	var headerRow = []string{"engine", "vendor", "description", "compat"}
	tableRows := [][]string{headerRow}

	// Sort by Score in descending order
	sort.Slice(enginesList.Engines, func(i, j int) bool {
		// Stable engines with equal score should be listed first
		if enginesList.Engines[i].Score == enginesList.Engines[j].Score {
			return enginesList.Engines[i].Grade == "stable"
		}
		return enginesList.Engines[i].Score > enginesList.Engines[j].Score
	})

	var engineNameMaxLen, engineVendorMaxLen int

	for _, engine := range enginesList.Engines {
		// Mark active engine with "*"
		if engine.Name == enginesList.ActiveEngine {
			engine.Name = engine.Name + "*"
		}

		// Find max name and vendor lengths
		engineNameMaxLen = max(engineNameMaxLen, len(engine.Name))
		engineVendorMaxLen = max(engineVendorMaxLen, len(engine.Vendor))

		row := []string{engine.Name, engine.Vendor, engine.Description}

		compatibleStr := ""
		if engine.Compatible && engine.Grade == "stable" {
			compatibleStr = "yes"
		} else if engine.Compatible {
			compatibleStr = "devel"
		} else {
			compatibleStr = "no"
		}
		row = append(row, compatibleStr)

		tableRows = append(tableRows, row)
	}

	if len(tableRows) == 1 {
		fmt.Fprintln(os.Stderr, "No engines found.")
		return nil
	}

	tableMaxWidth := 80
	// Increase column widths to account for paddings
	engineNameMaxLen += 1
	engineVendorMaxLen += 2
	// Description column fills the remaining space
	engineDescriptionMaxLen := tableMaxWidth - (engineNameMaxLen + engineVendorMaxLen)
	// Reserve space for Compatible column
	engineDescriptionMaxLen -= len(headerRow[3]) + 1
	options := []tablewriter.Option{
		tablewriter.WithRenderer(renderer.NewColorized(renderer.ColorizedConfig{
			Header: renderer.Tint{
				FG: renderer.Colors{color.Bold}, // Bold headers
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
					0: engineNameMaxLen,        // Engine name
					1: engineVendorMaxLen,      // Vendor
					2: engineDescriptionMaxLen, // Description
					// 3:  0, // Compatible, not set because cell value is shorter than min width
				},
			},
			Header: tw.CellConfig{
				Alignment: tw.CellAlignment{Global: tw.AlignLeft},
				Padding: tw.CellPadding{
					PerColumn: []tw.Padding{
						{Overwrite: true, Right: " "},
						{Overwrite: true, Left: " ", Right: " "},
						{Overwrite: true, Left: " ", Right: " "},
						{Overwrite: true},
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
						{Overwrite: true, Left: " ", Right: " "},
						{Overwrite: true},
					},
				},
			},
		}),
	}

	table := tablewriter.NewTable(os.Stdout, options...)
	table.Header(tableRows[0])
	err := table.Bulk(tableRows[1:])
	if err != nil {
		return fmt.Errorf("error adding data to table: %v", err)
	}
	err = table.Render()
	if err != nil {
		return fmt.Errorf("error rendering table: %v", err)
	}
	return nil
}
