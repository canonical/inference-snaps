package main

import (
	"encoding/json"
	"fmt"

	"github.com/canonical/stack-utils/pkg/types"
)

func parseStacksJson(stacksJson string) ([]types.ScoredStack, error) {
	var stacksOption map[string]map[string]types.ScoredStack
	err := json.Unmarshal([]byte(stacksJson), &stacksOption)
	if err != nil {
		return nil, fmt.Errorf("error parsing json: %v", err)
	}
	if stacksMap, ok := stacksOption["stacks"]; ok {
		var stacksSlice []types.ScoredStack
		for _, stack := range stacksMap {
			stacksSlice = append(stacksSlice, stack)
		}
		return stacksSlice, nil
	}
	return nil, fmt.Errorf("no stacks found")
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
