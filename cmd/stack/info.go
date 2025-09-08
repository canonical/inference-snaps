package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/canonical/go-snapctl"
	"github.com/canonical/stack-utils/pkg/types"
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
	return stackInfo(args[0])
}

func infoValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
	stacksJson, err := snapctl.Get("engines").Document().Run()
	if err != nil {
		fmt.Printf("Error loading engines: %v", err)
		return nil, cobra.ShellCompDirectiveError
	}

	stacks, err := parseStacksJson(stacksJson)
	if err != nil {
		fmt.Printf("Error parsing engines: %v", err)
		return nil, cobra.ShellCompDirectiveError
	}

	var stackNames []cobra.Completion
	for i := range stacks {
		stackNames = append(stackNames, stacks[i].Name)
	}

	return stackNames, cobra.ShellCompDirectiveNoSpace
}

func stackInfo(stackName string) error {
	stackJson, err := snapctl.Get("engines." + stackName).Document().Run()
	if err != nil {
		return fmt.Errorf("error loading engine: %v", err)
	}

	stack, err := parseStackJson(stackJson)
	if err != nil {
		return fmt.Errorf("error parsing engine: %v", err)
	}

	err = printStackInfo(stack)
	if err != nil {
		return fmt.Errorf("error printing engine info: %v", err)
	}
	return nil
}

func printStackInfo(stack types.ScoredStack) error {
	stackYaml, err := yaml.Marshal(stack)
	if err != nil {
		return fmt.Errorf("error converting engine to yaml: %v", err)
	}

	err = quick.Highlight(os.Stdout, string(stackYaml), "yaml", "terminal", "colorful")
	if err != nil {
		return fmt.Errorf("error formatting yaml: %v", err)
	}

	return nil
}
