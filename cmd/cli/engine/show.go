package main

import (
	"encoding/json"
	"fmt"

	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/selector"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	showEngineFormat string
)

func addInfoCommand() {
	cmd := &cobra.Command{
		Use:               "show-engine [<engine>]",
		Short:             "Print information about an engine",
		Long:              "Print information about the active engine, or the specified engine",
		GroupID:           "engines",
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: showEngineValidArgs,
		RunE:              showEngine,
	}
	cmd.PersistentFlags().StringVar(&showEngineFormat, "format", "yaml", "output format")
	rootCmd.AddCommand(cmd)
}

func showEngine(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		currentEngine, err := cache.GetActiveEngine()
		if err != nil {
			return fmt.Errorf("could not get the active engine: %v", err)
		}
		if currentEngine == "" {
			return fmt.Errorf("no active engine")
		}
		return engineInfo(currentEngine)

	} else if len(args) == 1 {
		return engineInfo(args[0])

	} else {
		return fmt.Errorf("invalid number of arguments")
	}
}

func showEngineValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
	manifests, err := selector.LoadManifestsFromDir(enginesDir)
	if err != nil {
		fmt.Printf("Error loading engines: %v\n", err)
		return nil, cobra.ShellCompDirectiveError
	}

	var engineNames []cobra.Completion
	for i := range manifests {
		engineNames = append(engineNames, manifests[i].Name)
	}

	return engineNames, cobra.ShellCompDirectiveNoSpace
}

func engineInfo(engineName string) error {
	scoredEngines, err := scoreEngines()
	if err != nil {
		return fmt.Errorf("error scoring engines: %v", err)
	}

	var scoredManifest engines.ScoredManifest
	for i := range scoredEngines {
		if scoredEngines[i].Name == engineName {
			scoredManifest = scoredEngines[i]
		}
	}
	if scoredManifest.Name != engineName {
		return fmt.Errorf(`engine "%s" does not exist`, engineName)
	}

	err = printEngineManifest(scoredManifest)
	if err != nil {
		return fmt.Errorf("error printing engine manifest: %v", err)
	}
	return nil
}

func printEngineManifest(engine engines.ScoredManifest) error {
	switch showEngineFormat {
	case "json":
		jsonString, err := json.MarshalIndent(engine, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal to JSON: %s", err)
		}
		fmt.Printf("%s\n", jsonString)
	case "yaml", "":
		engineYaml, err := yaml.Marshal(engine)
		if err != nil {
			return fmt.Errorf("failed to marshal to YAML: %s", err)
		}
		fmt.Print(string(engineYaml))
	default:
		return fmt.Errorf("unknown format %q", showEngineFormat)
	}

	return nil
}
