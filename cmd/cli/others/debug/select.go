package debug

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/selector"
	"github.com/canonical/inference-snaps-cli/pkg/types"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type selectCommand struct {
	*common.Context

	// flags
	format     string
	enginesDir string
}

type EngineSelection struct {
	Engines   []engines.ScoredManifest `json:"engines"`
	TopEngine string                   `json:"top-engine"`
}

func SelectCommand(ctx *common.Context) *cobra.Command {
	var cmd selectCommand
	cmd.Context = ctx

	cobra := &cobra.Command{
		Use:               "select-engine",
		Short:             "Test which engine will be chosen",
		Long:              "Test which engine will be chosen from a directory of engines, given the machine information piped in via stdin",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.run,
	}

	// flags
	cobra.Flags().StringVar(&cmd.format, "format", "yaml", "engine selection results format")
	cobra.Flags().StringVar(&cmd.enginesDir, "engines", ctx.EnginesDir, "engine manifests directory")

	return cobra
}

func (cmd *selectCommand) run(_ *cobra.Command, args []string) error {
	// Read yaml piped in from the hardware-info app
	var hardwareInfo types.HwInfo

	err := yaml.NewDecoder(os.Stdin).Decode(&hardwareInfo)
	if err != nil {
		return fmt.Errorf("error decoding hardware info: %s", err)
	}

	allEngines, err := engines.LoadManifests(cmd.enginesDir)
	if err != nil {
		return fmt.Errorf("error loading engines from directory: %s", err)
	}
	scoredEngines, err := selector.ScoreEngines(&hardwareInfo, allEngines)
	if err != nil {
		return fmt.Errorf("error scoring engines: %s", err)
	}

	var engineSelection EngineSelection

	// Print summary on STDERR
	for _, engine := range scoredEngines {
		engineSelection.Engines = append(engineSelection.Engines, engine)

		if engine.Score == 0 {
			fmt.Fprintf(os.Stderr, "❌ %s - not compatible: %s\n", engine.Name, strings.Join(engine.CompatibilityIssues, ", "))
		} else if engine.Grade != "stable" {
			fmt.Fprintf(os.Stderr, "🟠 %s - score = %d, grade = %s\n", engine.Name, engine.Score, engine.Grade)
		} else {
			fmt.Fprintf(os.Stderr, "✅ %s - compatible, score = %d\n", engine.Name, engine.Score)
		}
	}

	selectedEngine, err := selector.TopEngine(scoredEngines)
	if err != nil {
		return fmt.Errorf("error finding top engine: %v", err)
	}
	engineSelection.TopEngine = selectedEngine.Name

	greenBold := color.New(color.FgGreen, color.Bold).SprintFunc()
	fmt.Fprintf(os.Stderr, greenBold("Selected engine for your hardware configuration: %s\n\n"), selectedEngine.Name)

	var resultStr string
	switch cmd.format {
	case "json":
		jsonString, err := json.MarshalIndent(engineSelection, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal to JSON: %s", err)
		}
		resultStr = string(jsonString)
	case "yaml":
		yamlString, err := yaml.Marshal(engineSelection)
		if err != nil {
			return fmt.Errorf("failed to marshal to YAML: %s", err)
		}
		resultStr = string(yamlString)
	default:
		return fmt.Errorf("unknown format %q", cmd.format)
	}

	fmt.Println(resultStr)
	return nil
}
