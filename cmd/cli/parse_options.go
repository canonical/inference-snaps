package main

import (
	"encoding/json"
	"fmt"

	"github.com/canonical/stack-utils/pkg/engines"
)

func parseEnginesJson(enginesJson string) ([]engines.ScoredManifest, error) {
	var enginesOption map[string]engines.ScoredManifest
	err := json.Unmarshal([]byte(enginesJson), &enginesOption)
	if err != nil {
		return nil, fmt.Errorf("error parsing json: %v", err)
	}

	var enginesSlice []engines.ScoredManifest
	for _, engine := range enginesOption {
		enginesSlice = append(enginesSlice, engine)
	}

	return enginesSlice, nil
}

func parseEngineJson(engineJson string) (engines.ScoredManifest, error) {
	var scoredManifest engines.ScoredManifest

	err := json.Unmarshal([]byte(engineJson), &scoredManifest)
	if err != nil {
		return engines.ScoredManifest{}, fmt.Errorf("error parsing json: %v", err)
	}

	return scoredManifest, nil
}
