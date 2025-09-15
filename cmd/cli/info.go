package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/canonical/go-snapctl"
	"github.com/canonical/stack-utils/pkg/engines"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func addInfoCommand() {
	cmd := &cobra.Command{
		Use:   "show-engine <engine>",
		Short: "Print information about an engine",
		// Long:  "",
		GroupID:           "engines",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: infoValidArgs,
		RunE:              info,
	}
	rootCmd.AddCommand(cmd)
}

func info(_ *cobra.Command, args []string) error {
	return engineInfo(args[0])
}

func infoValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
	enginesJson, err := snapctl.Get("engines").Document().Run()
	if err != nil {
		fmt.Printf("Error loading engines: %v", err)
		return nil, cobra.ShellCompDirectiveError
	}

	engines, err := parseEnginesJson(enginesJson)
	if err != nil {
		fmt.Printf("Error parsing engines: %v", err)
		return nil, cobra.ShellCompDirectiveError
	}

	var engineNames []cobra.Completion
	for i := range engines {
		engineNames = append(engineNames, engines[i].Name)
	}

	return engineNames, cobra.ShellCompDirectiveNoSpace
}

func engineInfo(engineName string) error {
	enginesJson, err := snapctl.Get("engines." + engineName).Document().Run()
	if err != nil {
		return fmt.Errorf("error loading engine: %v", err)
	}

	engine, err := parseEngineJson(enginesJson)
	if err != nil {
		return fmt.Errorf("error parsing engine: %v", err)
	}

	err = printEngineInfo(engine)
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

	err = quick.Highlight(os.Stdout, string(engineYaml), "yaml", "terminal", "colorful")
	if err != nil {
		return fmt.Errorf("error formatting yaml: %v", err)
	}

	return nil
}
