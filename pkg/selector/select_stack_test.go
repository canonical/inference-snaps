package selector

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/canonical/ml-snap-utils/pkg/types"
	"gopkg.in/yaml.v3"
)

type stackTestSet struct {
	ValidHw   []string
	InvalidHw []string
}

var stackTestSets = map[string]stackTestSet{
	"ampere": {
		ValidHw: []string{"ampere-one-x-mocked"},
		InvalidHw: []string{
			"amd-ryzen7-5700g",
			"amd-ryzen9-7900",
			"ampere-altra",
			"dell-r730xd",
			"hp-dl380p-gen8",
			"i7-2600k+arc-a580",
			"i7-2600k",
			"mustang",
			"nuc11-i5-1145G7",
			"raspberry-pi-5",
			"xps13-7390",
			"xps13-9350",
		},
	},

	"ampere-altra": {
		ValidHw: []string{
			"ampere-altra",
		},
		InvalidHw: []string{
			"ampere-one-x-mocked",
			"amd-ryzen7-5700g",
			"amd-ryzen9-7900",
			"dell-r730xd",
			"hp-dl380p-gen8",
			"i7-2600k+arc-a580",
			"i7-2600k",
			"mustang",
			"nuc11-i5-1145G7",
			"raspberry-pi-5",
			"xps13-7390",
			"xps13-9350",
		},
	},

	"example-cpu": {
		ValidHw: []string{
			"amd-ryzen7-5700g",
			"amd-ryzen9-7900",
			"dell-r730xd",
			"mustang",
			"nuc11-i5-1145G7",
			"xps13-7390",
			"xps13-9350",
		},
		InvalidHw: []string{
			"ampere-altra",
			"ampere-one-x-mocked",
			"hp-dl380p-gen8",
			"i7-2600k+arc-a580",
			"i7-2600k",
			"raspberry-pi-5",
		},
	},

	"example-cpu-avx512": {
		ValidHw: []string{
			"amd-ryzen9-7900",
			"nuc11-i5-1145G7",
		},
		InvalidHw: []string{
			"amd-ryzen7-5700g",
			"ampere-altra",
			"ampere-one-x-mocked",
			"dell-r730xd",
			"hp-dl380p-gen8",
			"i7-2600k+arc-a580",
			"i7-2600k",
			"mustang",
			"raspberry-pi-5",
			"xps13-7390",
			"xps13-9350",
		},
	},

	"example-memory": {
		ValidHw: []string{
			"dell-r730xd",
			"hp-dl380p-gen8",
		},
		InvalidHw: []string{
			"amd-ryzen7-5700g",
			"amd-ryzen9-7900",
			"ampere-altra",
			"ampere-one-x-mocked",
			"i7-2600k+arc-a580",
			"i7-2600k",
			"mustang",
			"nuc11-i5-1145G7",
			"raspberry-pi-5",
			"xps13-7390",
			"xps13-9350",
		},
	},

	"generic-cuda": {
		ValidHw: []string{},
		InvalidHw: []string{
			"amd-ryzen7-5700g",
			"amd-ryzen9-7900",
			"ampere-altra",
			"ampere-one-x-mocked",
			"dell-r730xd",
			"hp-dl380p-gen8",
			"i7-2600k+arc-a580",
			"i7-2600k",
			"mustang",
			"nuc11-i5-1145G7",
			"raspberry-pi-5",
			"xps13-7390",
			"xps13-9350",
		},
	},

	"intel-dgpu": {
		ValidHw: []string{
			"i7-2600k+arc-a580",
			"mustang",
		},
		InvalidHw: []string{
			"amd-ryzen7-5700g",
			"amd-ryzen9-7900",
			"ampere-altra",
			"ampere-one-x-mocked",
			"dell-r730xd",
			"hp-dl380p-gen8",
			"i7-2600k",
			"nuc11-i5-1145G7",
			"raspberry-pi-5",
			"xps13-7390",
			"xps13-9350",
		},
	},

	"intel-npu": {
		ValidHw: []string{
			"xps13-9350",
		},
		InvalidHw: []string{
			"amd-ryzen7-5700g",
			"amd-ryzen9-7900",
			"ampere-altra",
			"ampere-one-x-mocked",
			"dell-r730xd",
			"hp-dl380p-gen8",
			"i7-2600k+arc-a580",
			"i7-2600k",
			"mustang",
			"nuc11-i5-1145G7",
			"raspberry-pi-5",
			"xps13-7390",
		},
	},
}

func TestStack(t *testing.T) {
	for stackName, testSet := range stackTestSets {
		for _, hwName := range testSet.ValidHw {
			t.Run(stackName+" == "+hwName, func(t *testing.T) {
				testValidHw(t, stackName, hwName)
			})
		}

		for _, hwName := range testSet.InvalidHw {
			t.Run(stackName+" != "+hwName, func(t *testing.T) {
				testInvalidHw(t, stackName, hwName)
			})
		}
	}
}

func testValidHw(t *testing.T, stackName string, hwName string) {
	stackManifestFile := fmt.Sprintf("../../test_data/stacks/%s/stack.yaml", stackName)
	hwInfoFile := fmt.Sprintf("../../test_data/hardware_info/%s.json", hwName)

	file, err := os.Open(hwInfoFile)
	if err != nil {
		t.Fatal(err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	var hardwareInfo types.HwInfo
	err = json.Unmarshal(data, &hardwareInfo)
	if err != nil {
		t.Fatal(err)
	}

	data, err = os.ReadFile(stackManifestFile)
	if err != nil {
		t.Fatal(err)
	}

	var stack types.Stack
	err = yaml.Unmarshal(data, &stack)
	if err != nil {
		t.Fatal(err)
	}

	// Valid hardware for stack
	score, reasons, err := checkStack(hardwareInfo, stack)
	if err != nil {
		t.Fatal(err)
	}
	if score == 0 {
		t.Fatalf("Stack should match: %v", reasons)
	}
	t.Logf("Matching score: %d", score)

}

func testInvalidHw(t *testing.T, stackName string, hwName string) {
	stackManifestFile := fmt.Sprintf("../../test_data/stacks/%s/stack.yaml", stackName)
	hwInfoFile := fmt.Sprintf("../../test_data/hardware_info/%s.json", hwName)

	file, err := os.Open(hwInfoFile)
	if err != nil {
		t.Fatal(err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	var hardwareInfo types.HwInfo
	err = json.Unmarshal(data, &hardwareInfo)
	if err != nil {
		t.Fatal(err)
	}

	data, err = os.ReadFile(stackManifestFile)
	if err != nil {
		t.Fatal(err)
	}

	var stack types.Stack
	err = yaml.Unmarshal(data, &stack)
	if err != nil {
		t.Fatal(err)
	}

	score, _, err := checkStack(hardwareInfo, stack)
	if err != nil {
		t.Fatal(err)
	}
	if score != 0 {
		t.Fatalf("Stack should not match: %s", hwName)
	}
	t.Logf("Matching score: %d", score)
}

func TestFindStackEmpty(t *testing.T) {
	hwInfo := types.HwInfo{
		Memory: &types.MemoryInfo{
			TotalRam:  200000000,
			TotalSwap: 200000000,
		},
		Disk: map[string]*types.DirStats{
			"/var/lib/snapd/snaps": &types.DirStats{
				Total: 0,
				Avail: 400000000,
			},
		},
	}

	allStacks, err := LoadStacksFromDir("../../test_data/stacks")
	if err != nil {
		t.Fatal(err)
	}
	scoredStacks, err := ScoreStacks(hwInfo, allStacks)
	if err != nil {
		t.Fatal(err)
	}
	topStack, err := TopStack(scoredStacks)
	if err == nil {
		t.Fatal("Empty stack dir should return an error for top stack")
	}
	if topStack != nil {
		t.Fatal("No stack should be found in empty stacks dir")
	}
}

func TestDiskCheck(t *testing.T) {
	dirStat := types.DirStats{
		Total: 0,
		Avail: 400000000,
	}
	hwInfo := types.HwInfo{}
	hwInfo.Disk = make(map[string]*types.DirStats)
	hwInfo.Disk["/"] = &dirStat
	hwInfo.Disk["/var/lib/snapd/snaps"] = &dirStat

	stackDisk := "300M"
	stack := types.Stack{DiskSpace: &stackDisk}

	result, reasons, err := checkStack(hwInfo, stack)
	if err != nil {
		t.Fatal(err)
	}
	if result == 0 {
		t.Fatalf("disk should be enough: %v", reasons)
	}

	dirStat.Avail = 100000000
	result, reasons, err = checkStack(hwInfo, stack)
	if err != nil {
		t.Fatal(err)
	}
	if result > 0 {
		t.Fatalf("disk should NOT be enough: %v", reasons)
	}
}

func TestMemoryCheck(t *testing.T) {
	hwInfo := types.HwInfo{
		Memory: &types.MemoryInfo{
			TotalRam:  200000000,
			TotalSwap: 200000000,
		},
	}

	stackMemory := "300M"
	stack := types.Stack{Memory: &stackMemory}

	result, reasons, err := checkStack(hwInfo, stack)
	if err != nil {
		t.Fatal(err)
	}
	if result == 0 {
		t.Fatalf("memory should be enough: %v", reasons)
	}

	hwInfo.Memory.TotalRam = 100000000
	result, reasons, err = checkStack(hwInfo, stack)
	if err != nil {
		t.Fatal(err)
	}
	if result > 0 {
		t.Fatal("memory should NOT be enough")
	}
}

func TestNoCpuInHwInfo(t *testing.T) {
	hwInfo := types.HwInfo{
		// All fields are nil or zero
	}

	data, err := os.ReadFile("../../test_data/stacks/example-cpu-avx512/stack.yaml")
	if err != nil {
		t.Fatal(err)
	}

	var currentStack types.Stack
	err = yaml.Unmarshal(data, &currentStack)
	if err != nil {
		t.Fatal(err)
	}

	// No memory in hardware info
	_, _, err = checkStack(hwInfo, currentStack)
	if err == nil {
		t.Fatalf("No Memory in hardware_info should return err")
	}

	hwInfo.Memory = &types.MemoryInfo{
		TotalRam:  17000000000,
		TotalSwap: 2000000000,
	}

	// No disk space in hardware info
	_, _, err = checkStack(hwInfo, currentStack)
	if err == nil {
		t.Fatal("No Disk space in hardware_info should return err")
	}

	hwInfo.Disk = make(map[string]*types.DirStats)
	hwInfo.Disk["/"] = &types.DirStats{
		Avail: 6000000000,
	}

	// No CPU in hardware info
	_, _, err = checkStack(hwInfo, currentStack)
	if err == nil {
		t.Fatal("No CPU in hardware_info should return err")
	}
}
