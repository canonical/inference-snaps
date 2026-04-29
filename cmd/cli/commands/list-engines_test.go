package commands

import (
	"testing"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/hardware_info"
	"github.com/canonical/inference-snaps-cli/pkg/selector"
	"github.com/canonical/inference-snaps-cli/pkg/storage"
)

func TestList(t *testing.T) {
	cache := storage.NewMockCache()
	err := cache.SetActiveEngine("example-memory")
	if err != nil {
		t.Fatalf("Error setting active engine name: %v", err)
	}

	allEngines, err := engines.LoadManifests("../../../test_data/engines")
	if err != nil {
		t.Fatalf("error loading engines: %v", err)
	}

	hardwareInfo, err := hardware_info.GetFromRawData("xps13-7390", true, "../../../test_data")
	if err != nil {
		t.Fatalf("error getting hardware info: %v", err)
	}

	scoredEngines, err := selector.ScoreEngines(hardwareInfo, allEngines)
	if err != nil {
		t.Fatalf("error scoring engines: %v", err)
	}

	// cmd.printEnginesTable needs to call `cmd.Cache.GetActiveEngine()` to get the current active engine
	// We therefore need to pass in the cache as context to `cmd`
	ctx := &common.Context{
		EnginesDir: "",
		Cache:      cache,
		Config:     nil,
	}
	cmd := listEnginesCommand{Context: ctx}

	activeEngine, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		t.Fatalf("%s: %v", common.LookingUpActiveEngine, err)
	}

	enginesList := outputEngines{
		ActiveEngine: activeEngine,
	}

	for _, se := range scoredEngines {
		enginesList.Engines = append(enginesList.Engines, common.NewEngineDetails(se))
	}

	err = cmd.printEnginesJson(enginesList)
	if err != nil {
		t.Fatal(err)
	}

	err = cmd.printEnginesTable(enginesList)
	if err != nil {
		t.Fatal(err)
	}
}
