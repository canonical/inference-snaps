package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

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
		GroupID: "stacks",
		Args:    cobra.NoArgs,
		RunE:    list,
	}

	// flags
	cmd.PersistentFlags().BoolVar(&listAll, "all", false, "include incompatible stacks")

	rootCmd.AddCommand(cmd)
}

func list(_ *cobra.Command, _ []string) error {
	return listStacks(listAll)
}

func listStacks(includeIncompatible bool) error {
	stacksJson, err := snapctl.Get("stacks").Document().Run()
	if err != nil {
		return fmt.Errorf("error loading stacks: %v", err)
	}

	stacks, err := parseStacksJson(stacksJson)
	if err != nil {
		return fmt.Errorf("error parsing stacks: %v", err)
	}

	err = printStacks(stacks, includeIncompatible)
	if err != nil {
		return fmt.Errorf("error printing list: %v", err)
	}

	return nil
}

func printStacks(stacks []types.ScoredStack, includeIncompatible bool) error {

	var headers []string
	if includeIncompatible {
		headers = []string{"Stack Name", "Vendor", "Description", "Compatible"}
	} else {
		headers = []string{"Stack Name", "Vendor", "Description"}
	}
	data := [][]string{headers}

	// Sort by Score in descending order
	sort.Slice(stacks, func(i, j int) bool {
		// Stable stacks with equal score should be listed first
		if stacks[i].Score == stacks[j].Score {
			return stacks[i].Grade == "stable"
		}
		return stacks[i].Score > stacks[j].Score
	})

	for _, stack := range stacks {
		stackInfo := []string{stack.Name, stack.Vendor, stack.Description}

		if includeIncompatible {
			// Compatible column is: yes|no|grade
			compatibleStr := ""
			if stack.Compatible && stack.Grade == "stable" {
				compatibleStr = "Yes"
			} else if stack.Compatible {
				compatibleStr = cases.Title(language.Und).String(stack.Grade)
			} else {
				compatibleStr = "No"
			}
			if len(stack.Notes) > 0 {
				compatibleStr = compatibleStr + "\n" + strings.Join(stack.Notes, ", ")
			}
			stackInfo = append(stackInfo, compatibleStr)
			data = append(data, stackInfo)
		} else if stack.Compatible && stack.Grade == "stable" {
			data = append(data, stackInfo)
		}
	}

	if len(data) == 1 {
		if includeIncompatible {
			_, err := fmt.Fprintln(os.Stderr, "No stacks found.")
			return err
		} else {
			_, err := fmt.Fprintln(os.Stderr, "No compatible stacks found.")
			return err
		}
	}

	// Configure colors: green headers, cyan/magenta rows, yellow footer
	colorCfg := renderer.ColorizedConfig{
		Header: renderer.Tint{
			FG: renderer.Colors{color.Bold}, // Green bold headers
			BG: renderer.Colors{},
		},
		Column: renderer.Tint{
			FG: renderer.Colors{color.Reset},
			BG: renderer.Colors{color.Reset},
		},
		Footer: renderer.Tint{
			FG: renderer.Colors{color.Reset, color.Bold},
			BG: renderer.Colors{color.Reset},
		},
		//Border:    renderer.Tint{FG: renderer.Colors{color.Reset}}, // White borders
		//Separator: renderer.Tint{FG: renderer.Colors{color.Reset}}, // White separators
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
	}

	tableMaxWidth := 80
	if includeIncompatible {
		tableMaxWidth = 120
	}

	table := tablewriter.NewTable(os.Stdout,
		tablewriter.WithRenderer(renderer.NewColorized(colorCfg)),
		tablewriter.WithConfig(tablewriter.Config{
			MaxWidth: tableMaxWidth,
			Header: tw.CellConfig{
				Alignment: tw.CellAlignment{Global: tw.AlignLeft},
				Formatting: tw.CellFormatting{
					AutoWrap:   tw.WrapNone,
					MergeMode:  tw.MergeNone,
					AutoFormat: tw.On,
				},
				Padding: tw.CellPadding{
					Global: tw.Padding{
						Left:      tw.Space, // Bug: making this empty causes the last char in a field to be cut off
						Right:     tw.Space,
						Top:       tw.Empty,
						Bottom:    tw.Empty,
						Overwrite: true,
					},
				},
			},
			Row: tw.CellConfig{
				Formatting: tw.CellFormatting{AutoWrap: tw.WrapNormal}, // Wrap long content
				Alignment:  tw.CellAlignment{Global: tw.AlignLeft},     // Left-align rows
				Padding: tw.CellPadding{
					Global: tw.Padding{
						Left:      tw.Space, // Bug: making this empty causes the last char in a field to be cut off
						Right:     tw.Space,
						Top:       tw.Empty,
						Bottom:    tw.Space,
						Overwrite: true,
					},
				},
			},
			Footer: tw.CellConfig{
				Alignment: tw.CellAlignment{Global: tw.AlignRight},
			},
		}),
	)

	table.Header(data[0])
	err := table.Bulk(data[1:])
	if err != nil {
		return fmt.Errorf("error adding data to table: %v", err)
	}
	err = table.Render()
	if err != nil {
		return fmt.Errorf("error rendering table: %v", err)
	}
	return nil
}
