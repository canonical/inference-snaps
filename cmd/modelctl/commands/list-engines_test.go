package commands

import (
	"fmt"
	"slices"
	"testing"

	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/common"
	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/canonical/inference-snaps-cli/v2/pkg/hardware_info"
	"github.com/canonical/inference-snaps-cli/v2/pkg/selector"
	"github.com/canonical/inference-snaps-cli/v2/pkg/storage"
)

func prepareTestData() (*listEnginesCommand, *outputEngines, error) {
	cache := storage.NewMockCache()
	err := cache.SetActiveEngine("intel-cpu")
	if err != nil {
		return nil, nil, fmt.Errorf("Error setting active engine name: %v", err)
	}

	allEngines, err := engines.LoadManifests("../../../test_data/engines")
	if err != nil {
		return nil, nil, fmt.Errorf("error loading engines: %v", err)
	}

	hardwareInfo, err := hardware_info.GetFromRawData("xps13-7390", true, "../../../test_data")
	if err != nil {
		return nil, nil, fmt.Errorf("error getting hardware info: %v", err)
	}

	scoredEngines, err := selector.ScoreEngines(hardwareInfo, allEngines)
	if err != nil {
		return nil, nil, fmt.Errorf("error scoring engines: %v", err)
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
		return nil, nil, fmt.Errorf("%s: %v", common.LookingUpActiveEngine, err)
	}

	enginesList := outputEngines{
		ActiveEngine: activeEngine,
	}

	for _, se := range scoredEngines {
		enginesList.Engines = append(enginesList.Engines, common.NewEngineDetails(se))
	}

	return &cmd, &enginesList, nil
}

func filterEnginesByName(engines *outputEngines, nameWhitelist []string) {
	var filteredEngines = make([]common.EngineDetails, 0, len(nameWhitelist))
	for _, ed := range engines.Engines {
		if slices.Contains(nameWhitelist, ed.Name) {
			filteredEngines = append(filteredEngines, ed)
		}
	}
	engines.Engines = filteredEngines
}

func TestList(t *testing.T) {
	cmd, enginesList, err := prepareTestData()
	if err != nil {
		t.Fatalf("Error preparing test data: %v", err)
	}

	err = cmd.printEnginesJson(*enginesList)
	if err != nil {
		t.Fatal(err)
	}

	err = cmd.printEnginesTable(*enginesList)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetEnginesTable(t *testing.T) {
	cmd, enginesList, err := prepareTestData()
	if err != nil {
		t.Fatalf("Error preparing test data: %v", err)
	}

	tableStr, err := cmd.getEnginesTable(*enginesList)
	if err != nil {
		t.Fatalf("Error getting engines table: %v", err)
	}

	expectedTable := `ENGINE                 VENDOR             SUMMARY                        COMPAT
intel-cpu*             Intel Corporation  Intel CPU engine               yes   
cpu-avx2               Canonical Ltd      CPU AVX2 engine                yes   
cpu-avx1               Canonical Ltd      CPU AVX1 engine                yes   
cpu                    Canonical Ltd      General CPU engine             yes   
cpu-exptl              Canonical Ltd      CPU development engine         exptl 
intel-gpu              Intel Corporation  Intel GPU engine               no    
not-compatible-engine  Canonical Ltd      This summary is too long and…  no    
cpu-avx512             Canonical Ltd      CPU AVX512 engine              no    
arm-neon               Canonical Ltd      ARM NEON engine                no    
ampere-altra           Canonical Ltd      Ampere Altra engine            no    
ampere                 Canonical Ltd      Ampere ARM64 engine            no    
intel-npu              Intel Corporation  Intel NPU engine               no    
rocm-generic           Canonical Ltd      ROCm generic engine            no    
amd-gpu                Canonical Ltd      AMD GPU engine                 no    
cuda-generic           Canonical Ltd      CUDA generic engine            no    
`

	if tableStr != expectedTable {
		t.Errorf("Engine table not as expected.\n\nGot:\n\n%s\n\nWant:\n\n%s", tableStr, expectedTable)
	}
}

func Example_printEnginesJson() {
	cmd, enginesList, err := prepareTestData()
	if err != nil {
		panic(fmt.Sprintf("Error preparing test data: %v", err))
	}

	// Reduce available engines to make the output more concise for this example test
	var engineWhitelist = []string{"amd-gpu", "intel-cpu"}
	filterEnginesByName(enginesList, engineWhitelist)

	err = cmd.printEnginesJson(*enginesList)
	if err != nil {
		panic(fmt.Sprintf("Error printing engines json: %v", err))
	}

	// Output:
	// {
	//   "active-engine": "intel-cpu",
	//   "engines": [
	//     {
	//       "name": "amd-gpu",
	//       "summary": "AMD GPU engine",
	//       "description": "AMD specific engine targeting one microarchitecture:\n  - gfx1032\n",
	//       "vendor": "Canonical Ltd",
	//       "devices": {
	//         "anyof": null,
	//         "allof": [
	//           {
	//             "type": "cpu",
	//             "architecture": "amd64"
	//           },
	//           {
	//             "type": "gpu",
	//             "vendor-id": "0x1002",
	//             "microarchitecture": "gfx1032",
	//             "compatibility-issues": [
	//               "device not found"
	//             ]
	//           }
	//         ]
	//       },
	//       "runtime": "",
	//       "model": {
	//         "default": "",
	//         "options": null
	//       },
	//       "configurations": null,
	//       "score": 0,
	//       "compatible": false,
	//       "compatibility-issues": [
	//         "required device not found"
	//       ]
	//     },
	//     {
	//       "name": "intel-cpu",
	//       "summary": "Intel CPU engine",
	//       "description": "Use Intel CPUs",
	//       "vendor": "Intel Corporation",
	//       "devices": {
	//         "anyof": null,
	//         "allof": [
	//           {
	//             "type": "cpu",
	//             "architecture": "amd64",
	//             "manufacturer-id": "GenuineIntel"
	//           }
	//         ]
	//       },
	//       "runtime": "openvino-model-server",
	//       "model": {
	//         "default": "4b-it-int4-fq-ov",
	//         "options": [
	//           "4b-it-int4-fq-ov"
	//         ]
	//       },
	//       "configurations": {
	//         "target-device": "CPU"
	//       },
	//       "score": 16,
	//       "compatible": true
	//     }
	//   ]
	// }
}
