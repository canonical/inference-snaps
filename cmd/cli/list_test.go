package main

import (
	"testing"

	"github.com/canonical/stack-utils/pkg/hardware_info"
	"github.com/canonical/stack-utils/pkg/selector"
)

func TestList(t *testing.T) {
	allEngines, err := selector.LoadManifestsFromDir("../../test_data/engines")
	if err != nil {
		t.Fatalf("error loading engines: %v", err)
	}

	hardwareInfo, err := hardware_info.GetFromRawData(t, "xps13-7390", true)
	if err != nil {
		t.Fatalf("error getting hardware info: %v", err)
	}

	scoredEngines, err := selector.ScoreEngines(hardwareInfo, allEngines)
	if err != nil {
		t.Fatalf("error scoring engines: %v", err)
	}

	err = printEnginesTable(scoredEngines)
	if err != nil {
		t.Fatal(err)
	}
}
