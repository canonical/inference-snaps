package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/canonical/go-snapctl"
	"github.com/canonical/stack-utils/pkg/types"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
)

var (
	listAll bool
)

func init() {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available stacks",
		// Long:  "",
		GroupID:           "stacks",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              list,
	}

	// flags
	cmd.PersistentFlags().BoolVar(&listAll, "all", false, "include incompatible stacks")

	rootCmd.AddCommand(cmd)
}

func list(_ *cobra.Command, _ []string) error {
	return listStacks(listAll)
}

func listStacks(all bool) error {
	stacksJson, err := snapctl.Get("stacks").Document().Run()
	if err != nil {
		return fmt.Errorf("error loading stacks: %v", err)
	}

	stacks, err := parseStacksJson(stacksJson)
	if err != nil {
		return fmt.Errorf("error parsing stacks: %v", err)
	}

	err = printStacks(stacks, all)
	if err != nil {
		return fmt.Errorf("error printing list: %v", err)
	}

	return nil
}

func printStacks(stacks []types.ScoredStack, all bool) error {

	var headerRow = []string{"stack", "vendor", "description"}
	if all {
		headerRow = append(headerRow, "compat")
	}
	tableRows := [][]string{headerRow}

	// Sort by Score in descending order
	sort.Slice(stacks, func(i, j int) bool {
		// Stable stacks with equal score should be listed first
		if stacks[i].Score == stacks[j].Score {
			return stacks[i].Grade == "stable"
		}
		return stacks[i].Score > stacks[j].Score
	})

	var stackNameMaxLen, stackVendorMaxLen int
	for _, stack := range stacks {
		row := []string{stack.Name, stack.Vendor, stack.Description}

		// Only for stacks that will be printed, find max name and vendor lengths
		if all || (stack.Compatible && stack.Grade == "stable") {
			stackNameMaxLen = max(stackNameMaxLen, len(stack.Name))
			stackVendorMaxLen = max(stackVendorMaxLen, len(stack.Vendor))
		}

		if all {
			compatibleStr := ""
			if stack.Compatible && stack.Grade == "stable" {
				compatibleStr = "yes"
			} else if stack.Compatible {
				compatibleStr = "beta"
			} else {
				compatibleStr = "no"
			}

			row = append(row, compatibleStr)
			tableRows = append(tableRows, row)
		} else if stack.Compatible && stack.Grade == "stable" {
			tableRows = append(tableRows, row)
		}
	}

	if len(tableRows) == 1 {
		if all {
			_, err := fmt.Fprintln(os.Stderr, "No stacks found.")
			return err
		} else {
			_, err := fmt.Fprintln(os.Stderr, "No compatible stacks found.")
			return err
		}
	}

	tableMaxWidth := 80

	// Increase column widths to account for paddings
	stackNameMaxLen += 2
	stackVendorMaxLen += 2
	// Description column fills the remaining space
	stackDescriptionMaxLen := tableMaxWidth - (stackNameMaxLen + stackVendorMaxLen)
	if all {
		// Reserve space for Compatible column if included
		stackDescriptionMaxLen -= len(headerRow[3]) + 2
	}

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
					0: stackNameMaxLen,        // Stack name
					1: stackVendorMaxLen,      // Vendor
					2: stackDescriptionMaxLen, // Description
					// 3:  0, // Compatible, not set because cell value is shorter than min width
				},
			},
			Header: tw.CellConfig{
				Alignment: tw.CellAlignment{Global: tw.AlignLeft},
			},
			Row: tw.CellConfig{
				Formatting: tw.CellFormatting{AutoWrap: tw.WrapTruncate},
				Alignment:  tw.CellAlignment{Global: tw.AlignLeft},
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
