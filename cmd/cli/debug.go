package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/canonical/famous-models-cli/pkg/engines"
	"github.com/canonical/famous-models-cli/pkg/selector"
	"github.com/canonical/famous-models-cli/pkg/types"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	debugMachineInfoFormat string
	debugSelectFormat      string
	debugEnginesDir        string
)

func addDebugCommand() {
	debugCmd := &cobra.Command{
		Use:    "debug",
		Long:   "Developer/debugging commands",
		Hidden: true,
	}

	validateCmd := &cobra.Command{
		Use:   "validate-engines",
		Short: "Validate engine manifest files",
		Args:  cobra.MinimumNArgs(1),
		RunE:  validateEngineManifests,
	}
	debugCmd.AddCommand(validateCmd)

	selectCmd := &cobra.Command{
		Use:   "select-engine",
		Short: "Test which engine will be chosen",
		Long:  "Test which engine will be chosen from a directory of engines, given the machine information piped in via stdin",
		RunE:  debugSelectEngine,
	}
	selectCmd.PersistentFlags().StringVar(&debugSelectFormat, "format", "yaml", "engine selection results format")
	// If engines flag is set, override the globally defined engines directory
	selectCmd.PersistentFlags().StringVar(&debugEnginesDir, "engines", enginesDir, "engine manifests directory")
	debugCmd.AddCommand(selectCmd)

	rootCmd.AddCommand(debugCmd)
}

func validateEngineManifests(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no engine manifest specified")
	}

	allManifestsValid := true
	for _, manifestPath := range args {
		err := engines.Validate(manifestPath)
		if err != nil {
			allManifestsValid = false
			fmt.Printf("❌ %s: %s\n", manifestPath, err)
		} else {
			fmt.Printf("✅ %s\n", manifestPath)
		}
	}

	if !allManifestsValid {
		return fmt.Errorf("not all manifests are valid")
	}
	return nil
}

type EngineSelection struct {
	Engines   []engines.ScoredManifest `json:"engines"`
	TopEngine string                   `json:"top-engine"`
}

func debugSelectEngine(_ *cobra.Command, args []string) error {

	// Read json piped in from the hardware-info app
	var hardwareInfo types.HwInfo

	err := json.NewDecoder(os.Stdin).Decode(&hardwareInfo)
	if err != nil {
		return fmt.Errorf("error decoding hardware info: %s", err)
	}

	allEngines, err := selector.LoadManifestsFromDir(debugEnginesDir)
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
			fmt.Fprintf(os.Stderr, "❌ %s - not compatible: %s\n", engine.Name, strings.Join(engine.Notes, ", "))
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
	switch debugSelectFormat {
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
		return fmt.Errorf("unknown format %q", debugSelectFormat)
	}

	fmt.Println(resultStr)
	return nil
}
