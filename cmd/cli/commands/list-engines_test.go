package commands

import (
	"fmt"
	"slices"
	"testing"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/hardware_info"
	"github.com/canonical/inference-snaps-cli/pkg/selector"
	"github.com/canonical/inference-snaps-cli/pkg/storage"
)

func prepareTestData() (*listEnginesCommand, *outputEngines, error) {
	cache := storage.NewMockCache()
	err := cache.SetActiveEngine("example-memory")
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

	expectedTable := `ENGINE                 VENDOR             DESCRIPTION                    COMPAT
intel-cpu              Intel Corporation  Use Intel CPUs                 yes   
cpu-avx2               Canonical Ltd      CPUs with AVX2                 yes   
cpu-avx1               Canonical Ltd      Legacy CPUs with only SSE4.2…  yes   
cpu                    Canonical Ltd      General CPU engine             yes   
cpu-devel              Canonical Ltd      Requires any CPU but is grad…  devel 
rocm-generic           Canonical Ltd      AMD GPUs using ROCm. All maj…  no    
not-compatible-engine  Canonical Ltd      This test engine is designed…  no    
intel-npu              Intel Corporation  Intel NPUs                     no    
intel-gpu              Intel Corporation  Modern Intel GPUs (>=gen 13)   no    
ampere                 Canonical Ltd      Test ampere selection          no    
example-memory*        Canonical Ltd      Legacy CPUs, offering full a…  no    
cuda-generic           Canonical Ltd      Nvidia GPUs using CUDA. All …  no    
amd-gpu                Canonical Ltd      AMD specific engine targetin…  no    
cpu-avx512             Canonical Ltd      CPUs with AVX512               no    
ampere-altra           Canonical Ltd      Test ampere selection          no    
arm-neon               Canonical Ltd      ARM CPUs with NEON instructi…  no    
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
	var engineWhitelist = []string{"amd-gpu", "example-memory"}
	filterEnginesByName(enginesList, engineWhitelist)

	err = cmd.printEnginesJson(*enginesList)
	if err != nil {
		panic(fmt.Sprintf("Error printing engines json: %v", err))
	}

	// Output:
	// {
	//   "active-engine": "example-memory",
	//   "engines": [
	//     {
	//       "name": "amd-gpu",
	//       "description": "AMD specific engine targeting only one microarchitecture.",
	//       "vendor": "Canonical Ltd",
	//       "grade": "stable",
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
	//       "memory": "2G",
	//       "disk-space": "5G",
	//       "components": [
	//         "dummy-component-1",
	//         "dummy-component-2"
	//       ],
	//       "configurations": null,
	//       "score": 0,
	//       "compatible": false,
	//       "compatibility-issues": [
	//         "required device not found"
	//       ]
	//     },
	//     {
	//       "name": "example-memory",
	//       "description": "Legacy CPUs, offering full accuracy but very high memory usage",
	//       "vendor": "Canonical Ltd",
	//       "grade": "stable",
	//       "devices": {
	//         "anyof": [
	//           {
	//             "type": "cpu",
	//             "architecture": "amd64",
	//             "manufacturer-id": "AuthenticAMD",
	//             "compatibility-issues": [
	//               "manufacturer id mismatch: GenuineIntel"
	//             ]
	//           },
	//           {
	//             "type": "cpu",
	//             "architecture": "amd64",
	//             "manufacturer-id": "GenuineIntel"
	//           }
	//         ],
	//         "allof": null
	//       },
	//       "memory": "35G",
	//       "disk-space": "29G",
	//       "components": [
	//         "dummy-component-3"
	//       ],
	//       "configurations": null,
	//       "score": 0,
	//       "compatible": false,
	//       "compatibility-issues": [
	//         "insufficient memory"
	//       ]
	//     }
	//   ]
	// }
}
