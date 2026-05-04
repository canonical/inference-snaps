package commands

import (
	"fmt"
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/hardware_info"
	"github.com/canonical/inference-snaps-cli/pkg/selector"
)

func TestInfoLong(t *testing.T) {
	engine, err := engines.LoadManifest("../../../test_data/engines", "intel-gpu")
	if err != nil {
		t.Fatal(err)
	}
	var scoredEngine = engines.ScoredManifest{Manifest: *engine}

	cmd := showEngineCommand{
		format: "yaml",
	}
	err = cmd.printEngineManifest(scoredEngine)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInfoShort(t *testing.T) {
	engine, err := engines.LoadManifest("../../../test_data/engines", "cpu-avx1")
	if err != nil {
		t.Fatal(err)
	}
	var scoredEngine = engines.ScoredManifest{Manifest: *engine}

	cmd := showEngineCommand{
		format: "yaml",
	}
	err = cmd.printEngineManifest(scoredEngine)
	if err != nil {
		t.Fatal(err)
	}
}

func scoreEngineAgainstMachine(engineName string, machineName string) (*engines.ScoredManifest, error) {
	engineManifest, err := engines.LoadManifest("../../../test_data/engines", engineName)
	if err != nil {
		return nil, fmt.Errorf("failed to load engine manifest: %v", err)
	}
	info, err := hardware_info.GetFromRawData(machineName, true, "../../../test_data")
	if err != nil {
		return nil, fmt.Errorf("failed to get hardware info: %v", err)
	}
	scoredEngines, err := selector.ScoreEngines(info, []engines.Manifest{*engineManifest})
	if err != nil {
		return nil, fmt.Errorf("failed to score engines: %v", err)
	}

	if len(scoredEngines) != 1 {
		return nil, fmt.Errorf("invalid scored engines count: %d, expected 1", len(scoredEngines))
	}

	return &scoredEngines[0], nil
}

func TestUnsupportedFormatResultsInError(t *testing.T) {
	engineManifest, err := scoreEngineAgainstMachine("cpu-avx1", "dummy-machine")
	if err != nil {
		t.Fatalf("could not score manifest: %v", err)
	}

	cmd := showEngineCommand{format: "invalid-format"}
	err = cmd.printEngineManifest(*engineManifest)

	if err == nil {
		t.Fatalf("expected unsupported format to error out, got nil error")
	}
}

func Example_showEngineCommand_printEngineManifestYaml() {
	engineManifest, err := scoreEngineAgainstMachine("cuda-generic", "dummy-machine")
	if err != nil {
		panic(fmt.Sprintf("failed to score engine against machine: %v", err))
	}

	cmd := showEngineCommand{format: "yaml"}
	if err := cmd.printEngineManifest(*engineManifest); err != nil {
		panic(fmt.Sprintf("failed to print engine manifest: %v", err))
	}

	// Output:
	// name: cuda-generic
	// description: Nvidia GPUs using CUDA. All major CUDA versions are targeted.
	// vendor: Canonical Ltd
	// grade: stable
	// devices:
	//     allof:
	//         - type: cpu
	//           architecture: amd64
	//           flags:
	//             - sse4_2
	//             - f16c
	//             - fma
	//             - avx
	//             - avx2
	//           compatibility-issues:
	//             - flag sse4_2 missing
	//             - flag f16c missing
	//             - flag fma missing
	//             - flag avx missing
	//             - flag avx2 missing
	//         - type: gpu
	//           bus: pci
	//           vendor-id: "0x10DE"
	//           vram: 5G
	//           compatibility-issues:
	//             - device not found
	// memory: 2G
	// disk-space: 5G
	// components:
	//     - dummy-component-1
	//     - dummy-component-2
	// configurations: {}
	// score: 0
	// compatible: false
	// compatibility-issues:
	//     - required device not found
}

func Example_showEngineCommand_printEngineManifestJson() {
	engineManifest, err := scoreEngineAgainstMachine("cuda-generic", "dummy-machine")
	if err != nil {
		panic(fmt.Sprintf("failed to score engine against machine: %v", err))
	}

	cmd := showEngineCommand{format: "json"}
	if err := cmd.printEngineManifest(*engineManifest); err != nil {
		panic(fmt.Sprintf("failed to print engine manifest: %v", err))
	}

	// Output:
	// {
	//   "name": "cuda-generic",
	//   "description": "Nvidia GPUs using CUDA. All major CUDA versions are targeted.",
	//   "vendor": "Canonical Ltd",
	//   "grade": "stable",
	//   "devices": {
	//     "anyof": null,
	//     "allof": [
	//       {
	//         "type": "cpu",
	//         "architecture": "amd64",
	//         "flags": [
	//           "sse4_2",
	//           "f16c",
	//           "fma",
	//           "avx",
	//           "avx2"
	//         ],
	//         "compatibility-issues": [
	//           "flag sse4_2 missing",
	//           "flag f16c missing",
	//           "flag fma missing",
	//           "flag avx missing",
	//           "flag avx2 missing"
	//         ]
	//       },
	//       {
	//         "type": "gpu",
	//         "bus": "pci",
	//         "vendor-id": "0x10DE",
	//         "vram": "5G",
	//         "compatibility-issues": [
	//           "device not found"
	//         ]
	//       }
	//     ]
	//   },
	//   "memory": "2G",
	//   "disk-space": "5G",
	//   "components": [
	//     "dummy-component-1",
	//     "dummy-component-2"
	//   ],
	//   "configurations": null,
	//   "score": 0,
	//   "compatible": false,
	//   "compatibility-issues": [
	//     "required device not found"
	//   ]
	// }
}

func Example_showEngineCommand_printHappyEngineManifestYaml() {
	engineManifest, err := scoreEngineAgainstMachine("intel-cpu", "i7-1165G7")
	if err != nil {
		panic(fmt.Sprintf("failed to score engine against machine: %v", err))
	}

	cmd := showEngineCommand{format: "yaml"}
	if err := cmd.printEngineManifest(*engineManifest); err != nil {
		panic(fmt.Sprintf("failed to print engine manifest: %v", err))
	}

	// Output:
	// name: intel-cpu
	// description: Use Intel CPUs
	// vendor: Intel Corporation
	// grade: stable
	// devices:
	//     allof:
	//         - type: cpu
	//           architecture: amd64
	//           manufacturer-id: GenuineIntel
	// memory: 5G
	// disk-space: 10G
	// components:
	//     - dummy-component-1
	// configurations: {}
	// score: 18
	// compatible: true
}

func Example_showEngineCommand_printHappyEngineManifestJson() {
	engineManifest, err := scoreEngineAgainstMachine("intel-cpu", "i7-1165G7")
	if err != nil {
		panic(fmt.Sprintf("failed to score engine against machine: %v", err))
	}

	cmd := showEngineCommand{format: "json"}
	if err := cmd.printEngineManifest(*engineManifest); err != nil {
		panic(fmt.Sprintf("failed to print engine manifest: %v", err))
	}

	// Output:
	// {
	//   "name": "intel-cpu",
	//   "description": "Use Intel CPUs",
	//   "vendor": "Intel Corporation",
	//   "grade": "stable",
	//   "devices": {
	//     "anyof": null,
	//     "allof": [
	//       {
	//         "type": "cpu",
	//         "architecture": "amd64",
	//         "manufacturer-id": "GenuineIntel"
	//       }
	//     ]
	//   },
	//   "memory": "5G",
	//   "disk-space": "10G",
	//   "components": [
	//     "dummy-component-1"
	//   ],
	//   "configurations": null,
	//   "score": 18,
	//   "compatible": true
	// }
}
