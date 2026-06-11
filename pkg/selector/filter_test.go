package selector

import (
	"os"
	"testing"

	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/canonical/inference-snaps-cli/v2/pkg/types"
	"gopkg.in/yaml.v3"
)

/*
If the model snap has no engines defined, scoring should pass, but finding a top engine should not be possible.
*/
func TestFindTopEngineFromNone(t *testing.T) {
	hwInfo := types.HwInfo{}

	allEngines, err := engines.LoadManifests("../../test_data/engines")
	if err != nil {
		t.Fatal(err)
	}
	scoredEngines, err := ScoreEngines(&hwInfo, allEngines)
	if err != nil {
		t.Fatal(err)
	}
	topEngine, err := TopEngine(scoredEngines)
	if err == nil {
		t.Fatal("TopEngine should return an error if no engines are provided")
	}
	if topEngine != nil {
		t.Fatal("No top engine should be returned if no engines are provided")
	}
}

func TestNoCpuInHwInfo(t *testing.T) {
	hwInfo := types.HwInfo{
		// All fields are nil or zero
	}

	data, err := os.ReadFile("../../test_data/engines/cpu-avx512/" + engines.ManifestFilename)
	if err != nil {
		t.Fatal(err)
	}

	var currentEngine engines.Manifest
	err = yaml.Unmarshal(data, &currentEngine)
	if err != nil {
		t.Fatal(err)
	}


	// No CPU in hardware info
	_, report, err := checkEngine(&hwInfo, currentEngine)
	if err != nil {
		t.Fatal(err)
	}
	if report.EngineCompatible() {
		t.Fatal("Missing CPU info in hardware_info should result in an incompatible engine")
	}
}
