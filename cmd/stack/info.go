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

func init() {
	cmd := &cobra.Command{
		Use:   "info <stack>",
		Short: "Print information about a stack",
		// Long:  "",
		Args: cobra.ExactArgs(1),
		RunE: info,
	}
	rootCmd.AddCommand(cmd)
}

func info(_ *cobra.Command, args []string) error {
	return stackInfo(args[0])
}
func stackInfo(stackName string) error {
	stackJson, err := snapctl.Get("stacks." + stackName).Document().Run()
	if err != nil {
		return fmt.Errorf("error loading stack: %v", err)
	}

	stack, err := parseStackJson(stackJson)
	if err != nil {
		return fmt.Errorf("error parsing stack: %v", err)
	}

	err = printStackInfo(stack)
	if err != nil {
		return fmt.Errorf("error printing stack info: %v", err)
	}
	return nil
}

func printStackInfo(stack types.ScoredStack) error {
	stackYaml, err := yaml.Marshal(stack)
	if err != nil {
		return fmt.Errorf("error converting stack to yaml: %v", err)
	}

	err = quick.Highlight(os.Stdout, string(stackYaml), "yaml", "terminal", "colorful")
	if err != nil {
		return fmt.Errorf("error formatting yaml: %v", err)
	}

	return nil
}
