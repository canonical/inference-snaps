package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/canonical/stack-utils/pkg/engines"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
)

func addListCommand() {
	cmd := &cobra.Command{
		Use:   "list-engines",
		Short: "List available engines",
		// Long:  "",
		GroupID:           "engines",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              listEngines,
	}

	rootCmd.AddCommand(cmd)
}

func listEngines(_ *cobra.Command, _ []string) error {
	enginesJson, err := config.Get("engines")
	if err != nil {
		return fmt.Errorf("error loading engines: %v", err)
	}

	scoredEngines, err := parseEnginesJson(enginesJson)
	if err != nil {
		return fmt.Errorf("error parsing engines: %v", err)
	}

	err = printEnginesTable(scoredEngines)
	if err != nil {
		return fmt.Errorf("error printing list: %v", err)
	}

	return nil
}

func printEnginesTable(scoredEngines []engines.ScoredManifest) error {
	var headerRow = []string{"engine", "vendor", "description", "compat"}
	tableRows := [][]string{headerRow}

	// Sort by Score in descending order
	sort.Slice(scoredEngines, func(i, j int) bool {
		// Stable engines with equal score should be listed first
		if scoredEngines[i].Score == scoredEngines[j].Score {
			return scoredEngines[i].Grade == "stable"
		}
		return scoredEngines[i].Score > scoredEngines[j].Score
	})

	var engineNameMaxLen, engineVendorMaxLen int
	for _, engine := range scoredEngines {
		// Find max name and vendor lengths
		engineNameMaxLen = max(engineNameMaxLen, len(engine.Name))
		engineVendorMaxLen = max(engineVendorMaxLen, len(engine.Vendor))

		row := []string{engine.Name, engine.Vendor, engine.Description}

		compatibleStr := ""
		if engine.Compatible && engine.Grade == "stable" {
			compatibleStr = "yes"
		} else if engine.Compatible {
			compatibleStr = "beta"
		} else {
			compatibleStr = "no"
		}
		row = append(row, compatibleStr)

		tableRows = append(tableRows, row)
	}

	if len(tableRows) == 1 {
		_, err := fmt.Fprintln(os.Stderr, "No engines found.")
		return err
	}

	tableMaxWidth := 80

	// Increase column widths to account for paddings
	engineNameMaxLen += 2
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
