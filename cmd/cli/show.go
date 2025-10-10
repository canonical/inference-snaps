package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/canonical/famous-models-cli/pkg/engines"
	"github.com/canonical/famous-models-cli/pkg/selector"
	"github.com/canonical/famous-models-cli/pkg/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func addInfoCommand() {
	cmd := &cobra.Command{
		Use:               "show-engine [<engine>]",
		Short:             "Print information about an engine",
		Long:              "Print information about the currently selected engine, or the specified engine",
		GroupID:           "engines",
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: showEngineValidArgs,
		RunE:              showEngine,
	}
	rootCmd.AddCommand(cmd)
}

func showEngine(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		currentEngine, err := cache.GetActiveEngine()
		if err != nil {
			return fmt.Errorf("could not get currently selected engine: %v", err)
		}
		if currentEngine == "" {
			return fmt.Errorf(`No active engine. Select one with "use-engine".`)
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

	err = printEngineInfo(scoredManifest)
	if err != nil {
		return fmt.Errorf("error printing engine info: %v", err)
	}
	return nil
}

func printEngineInfo(engine engines.ScoredManifest) error {
	engineYaml, err := yaml.Marshal(engine)
	if err != nil {
		return fmt.Errorf("error converting engine to yaml: %v", err)
	}

	if utils.IsTerminalOutput() {
		err = quick.Highlight(os.Stdout, string(engineYaml), "yaml", "terminal", "colorful")
		if err != nil {
			return fmt.Errorf("error formatting yaml: %v", err)
		}
	} else {
		fmt.Print(string(engineYaml))
		return nil
	}

	return nil
}
