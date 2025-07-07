package selector

import (
	"fmt"
	"os"
	"testing"

	"github.com/canonical/stack-utils/pkg/hardware_info"
	"github.com/canonical/stack-utils/pkg/types"
	"gopkg.in/yaml.v3"
)

type stackTestSet struct {
	ValidDevices   []string
	InvalidDevices []string
}

var stackTestSets = map[string]stackTestSet{
	"ampere": {
		ValidDevices: []string{
			"ampere-one-m-banshee-12",
			"ampere-one-siryn",
			"ampere-one-x-banshee-8",
		},
		InvalidDevices: []string{
			"asus-ux301l",
			"hp-proliant-rl300-gen11-altra",
			"hp-proliant-rl300-gen11-altra-max",
			"i7-1165G7",
			"i7-2600k+arc-a580",
			"i7-10510U",
			"mustang",
			//"orange-pi-rv2",
			"raspberry-pi-5",
			"raspberry-pi-5+hailo-8",
			"system76-addw4",
			"xps13-7390",
			"xps13-9350",
		},
	},

	"ampere-altra": {
		ValidDevices: []string{
			"hp-proliant-rl300-gen11-altra",
			"hp-proliant-rl300-gen11-altra-max",
		},
		InvalidDevices: []string{
			"ampere-one-m-banshee-12",
			"ampere-one-siryn",
			"ampere-one-x-banshee-8",
			"asus-ux301l",
			"i7-1165G7",
			"i7-2600k+arc-a580",
			"i7-10510U",
			"mustang",
			//"orange-pi-rv2",
			"raspberry-pi-5",
			"raspberry-pi-5+hailo-8",
			"system76-addw4",
			"xps13-7390",
			"xps13-9350",
		},
	},

	"arm-neon": {
		ValidDevices: []string{
			"ampere-one-m-banshee-12",
			"ampere-one-siryn",
			"ampere-one-x-banshee-8",
			"raspberry-pi-5",
			"raspberry-pi-5+hailo-8",
			"hp-proliant-rl300-gen11-altra",
			"hp-proliant-rl300-gen11-altra-max",
		},
		InvalidDevices: []string{
			"asus-ux301l",
			"i7-1165G7",
			"i7-2600k+arc-a580",
			"i7-10510U",
			"mustang",
			//"orange-pi-rv2",
			"system76-addw4",
			"xps13-7390",
			"xps13-9350",
		},
	},

	"cpu-avx1": {
		ValidDevices: []string{
			"asus-ux301l",
			"i7-1165G7",
			"i7-2600k+arc-a580",
			"i7-10510U",
			"mustang",
			"system76-addw4",
			"xps13-7390",
			"xps13-9350",
		},
		InvalidDevices: []string{
			"ampere-one-m-banshee-12",
			"ampere-one-siryn",
			"ampere-one-x-banshee-8",
			"hp-proliant-rl300-gen11-altra",
			"hp-proliant-rl300-gen11-altra-max",
			//"orange-pi-rv2",
			"raspberry-pi-5",
			"raspberry-pi-5+hailo-8",
		},
	},

	"cpu-avx2": {
		ValidDevices: []string{
			"asus-ux301l",
			"i7-1165G7",
			"i7-10510U",
			"mustang",
			"system76-addw4",
			"xps13-7390",
			"xps13-9350",
		},
		InvalidDevices: []string{
			"ampere-one-m-banshee-12",
			"ampere-one-siryn",
			"ampere-one-x-banshee-8",
			"hp-proliant-rl300-gen11-altra",
			"hp-proliant-rl300-gen11-altra-max",
			"i7-2600k+arc-a580",
			//"orange-pi-rv2",
			"raspberry-pi-5",
			"raspberry-pi-5+hailo-8",
		},
	},

	"cpu-avx512": {
		ValidDevices: []string{
			"i7-1165G7",
		},
		InvalidDevices: []string{
			"ampere-one-m-banshee-12",
			"ampere-one-siryn",
			"ampere-one-x-banshee-8",
			"asus-ux301l",
			"hp-proliant-rl300-gen11-altra",
			"hp-proliant-rl300-gen11-altra-max",
			"i7-2600k+arc-a580",
			"i7-10510U",
			"mustang",
			//"orange-pi-rv2",
			"raspberry-pi-5",
			"raspberry-pi-5+hailo-8",
			"system76-addw4",
			"xps13-7390",
			"xps13-9350",
		},
	},

	"example-memory": {
		ValidDevices: []string{
			"mustang",
			"system76-addw4",
			"xps13-9350",
		},
		InvalidDevices: []string{
			"ampere-one-m-banshee-12",
			"ampere-one-siryn",
			"ampere-one-x-banshee-8",
			"asus-ux301l",
			"hp-proliant-rl300-gen11-altra",
			"hp-proliant-rl300-gen11-altra-max",
			"i7-1165G7",
			"i7-2600k+arc-a580",
			"i7-10510U",
			//"orange-pi-rv2",
			"raspberry-pi-5",
			"raspberry-pi-5+hailo-8",
			"xps13-7390",
		},
	},

	"generic-cuda": {
		ValidDevices: []string{},
		InvalidDevices: []string{
			"ampere-one-m-banshee-12",
			"ampere-one-siryn",
			"ampere-one-x-banshee-8",
			"asus-ux301l",
			"hp-proliant-rl300-gen11-altra",
			"hp-proliant-rl300-gen11-altra-max",
			"i7-1165G7",
			"i7-2600k+arc-a580",
			"i7-10510U",
			"mustang",
			//"orange-pi-rv2",
			"raspberry-pi-5",
			"raspberry-pi-5+hailo-8",
			"system76-addw4",
			"xps13-7390",
			"xps13-9350",
		},
	},

	"intel-gpu": {
		ValidDevices: []string{
			"i7-2600k+arc-a580",
			"mustang",
			"system76-addw4",
			"xps13-9350",
		},
		InvalidDevices: []string{
			"ampere-one-m-banshee-12",
			"ampere-one-siryn",
			"ampere-one-x-banshee-8",
			"asus-ux301l", // has intel gpu, but clinfo not working
			"hp-proliant-rl300-gen11-altra",
			"hp-proliant-rl300-gen11-altra-max",
			"i7-1165G7", // 9a49 TigerLake-LP GT2 [Iris Xe Graphics]
			"i7-10510U",
			//"orange-pi-rv2",
			"raspberry-pi-5",
			"raspberry-pi-5+hailo-8",
			"xps13-7390",
		},
	},

	"intel-npu": {
		ValidDevices: []string{
			"xps13-9350",
		},
		InvalidDevices: []string{
			"ampere-one-m-banshee-12",
			"ampere-one-siryn",
			"ampere-one-x-banshee-8",
			"asus-ux301l",
			"hp-proliant-rl300-gen11-altra",
			"hp-proliant-rl300-gen11-altra-max",
			"i7-1165G7",
			"i7-2600k+arc-a580",
			"i7-10510U",
			"mustang",
			//"orange-pi-rv2",
			"raspberry-pi-5",
			"raspberry-pi-5+hailo-8",
			"system76-addw4",
			"xps13-7390",
		},
	},
}

func TestStack(t *testing.T) {
	for stackName, testSet := range stackTestSets {
		for _, hwName := range testSet.ValidDevices {
			t.Run(stackName+" == "+hwName, func(t *testing.T) {
				testValidHw(t, stackName, hwName)
			})
		}

		for _, hwName := range testSet.InvalidDevices {
			t.Run(stackName+" != "+hwName, func(t *testing.T) {
				testInvalidHw(t, stackName, hwName)
			})
		}
	}
}

func testValidHw(t *testing.T, stackName string, hwName string) {
	stackManifestFile := fmt.Sprintf("../../test_data/stacks/%s/stack.yaml", stackName)

	hardwareInfo, err := hardware_info.GetFromRawData(t, hwName, true)
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(stackManifestFile)
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

	hardwareInfo, err := hardware_info.GetFromRawData(t, hwName, true)
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(stackManifestFile)
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
		Memory: types.MemoryInfo{
			TotalRam:  200000000,
			TotalSwap: 200000000,
		},
		Disk: map[string]types.DirStats{
			"/var/lib/snapd/snaps": {
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
	hwInfo.Disk = make(map[string]types.DirStats)
	hwInfo.Disk["/"] = dirStat
	hwInfo.Disk["/var/lib/snapd/snaps"] = dirStat

	stackDisk := "300M"
	stack := types.Stack{DiskSpace: &stackDisk}

	result, reasons, err := checkStack(hwInfo, stack)
	if err != nil {
		t.Fatal(err)
	}
	if result == 0 {
		t.Fatalf("disk should be enough: %v", reasons)
	}

	dirStat = types.DirStats{
		Total: 0,
		Avail: 100000000,
	}
	hwInfo.Disk["/var/lib/snapd/snaps"] = dirStat
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
		Memory: types.MemoryInfo{
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

	data, err := os.ReadFile("../../test_data/stacks/cpu-avx512/stack.yaml")
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

	hwInfo.Memory = types.MemoryInfo{
		TotalRam:  17000000000,
		TotalSwap: 2000000000,
	}

	// No disk space in hardware info
	_, _, err = checkStack(hwInfo, currentStack)
	if err == nil {
		t.Fatal("No Disk space in hardware_info should return err")
	}

	hwInfo.Disk = make(map[string]types.DirStats)
	hwInfo.Disk["/"] = types.DirStats{
		Avail: 6000000000,
	}

	// No CPU in hardware info
	_, _, err = checkStack(hwInfo, currentStack)
	if err == nil {
		t.Fatal("No CPU in hardware_info should return err")
	}
}
