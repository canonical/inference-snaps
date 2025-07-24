package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/canonical/go-snapctl"
	"github.com/canonical/stack-utils/pkg/types"
	"gopkg.in/yaml.v3"
)

func stackInfo(stackName string) {
	stackJson, err := snapctl.Get("stacks." + stackName).Document().Run()
	if err != nil {
		log.Fatalf("Error loading stack: %v\n", err)
	}

	stack, err := parseStackJson(stackJson)
	if err != nil {
		log.Fatalf("Error parsing stack: %v\n", err)
	}
	err = printStackInfo(stack)
	if err != nil {
		log.Fatalf("Error printing stack info: %v\n", err)
	}
}

func parseStackJson(stackJson string) (types.ScoredStack, error) {
	var stackOption map[string]types.ScoredStack

	err := json.Unmarshal([]byte(stackJson), &stackOption)
	if err != nil {
		return types.ScoredStack{}, fmt.Errorf("error parsing json: %v", err)
	}

	if len(stackOption) == 0 {
		return types.ScoredStack{}, fmt.Errorf("stack not found")
	}

	if len(stackOption) > 1 {
		return types.ScoredStack{}, fmt.Errorf("only one stack expected in json")
	}

	for _, stack := range stackOption {
		return stack, nil
	}

	return types.ScoredStack{}, fmt.Errorf("unexpected error occurred")
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
