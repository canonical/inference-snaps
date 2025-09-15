package main

import (
	"encoding/json"
	"fmt"

	"github.com/canonical/stack-utils/pkg/engines"
)

func parseEnginesJson(enginesJson string) ([]engines.ScoredManifest, error) {
	var enginesOption map[string]map[string]engines.ScoredManifest
	err := json.Unmarshal([]byte(enginesJson), &enginesOption)
	if err != nil {
		return nil, fmt.Errorf("error parsing json: %v", err)
	}
	if enginesMap, ok := enginesOption["engines"]; ok {
		var enginesSlice []engines.ScoredManifest
		for _, engine := range enginesMap {
			enginesSlice = append(enginesSlice, engine)
		}
		return enginesSlice, nil
	}
	return nil, fmt.Errorf("no engines found")
}

func parseEngineJson(engineJson string) (engines.ScoredManifest, error) {
	var engineOption map[string]engines.ScoredManifest

	err := json.Unmarshal([]byte(engineJson), &engineOption)
	if err != nil {
		return engines.ScoredManifest{}, fmt.Errorf("error parsing json: %v", err)
	}

	if len(engineOption) == 0 {
		return engines.ScoredManifest{}, fmt.Errorf("engine not found")
	}

	if len(engineOption) > 1 {
		return engines.ScoredManifest{}, fmt.Errorf("only one engine expected in json")
	}

	for _, engine := range engineOption {
		return engine, nil
	}

	return engines.ScoredManifest{}, fmt.Errorf("unexpected error occurred")
}
