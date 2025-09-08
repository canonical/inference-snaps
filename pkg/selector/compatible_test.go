package selector

import (
	"fmt"
	"os"
	"testing"

	"github.com/canonical/stack-utils/pkg/hardware_info"
	"github.com/canonical/stack-utils/pkg/types"
	"gopkg.in/yaml.v3"
)

type testValidInvalid struct {
	ValidMachines   []string
	InvalidMachines []string
}

var validInvalidSets = map[string]testValidInvalid{
	"ampere": {
		ValidMachines: []string{
			"ampere-one-m-banshee-12",
			"ampere-one-siryn",
			"ampere-one-x-banshee-8",
		},
		InvalidMachines: []string{
			"asus-ux301l",
			"hp-pavilion-15-cs-3037nl",
			"hp-proliant-rl300-gen11-altra",
			"hp-proliant-rl300-gen11-altra-max",
			"hp-zbook-power-16-inch-g11",
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
		ValidMachines: []string{
			"hp-proliant-rl300-gen11-altra",
			"hp-proliant-rl300-gen11-altra-max",
		},
		InvalidMachines: []string{
			"ampere-one-m-banshee-12",
			"ampere-one-siryn",
			"ampere-one-x-banshee-8",
			"asus-ux301l",
			"hp-pavilion-15-cs-3037nl",
			"hp-zbook-power-16-inch-g11",
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
		ValidMachines: []string{
			"ampere-one-m-banshee-12",
			"ampere-one-siryn",
			"ampere-one-x-banshee-8",
			"raspberry-pi-5",
			"raspberry-pi-5+hailo-8",
			"hp-proliant-rl300-gen11-altra",
			"hp-proliant-rl300-gen11-altra-max",
		},
		InvalidMachines: []string{
			"asus-ux301l",
			"hp-pavilion-15-cs-3037nl",
			"hp-zbook-power-16-inch-g11",
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
		ValidMachines: []string{
			"asus-ux301l",
			"hp-pavilion-15-cs-3037nl",
			"hp-zbook-power-16-inch-g11",
			"i7-1165G7",
			"i7-2600k+arc-a580",
			"i7-10510U",
			"mustang",
			"system76-addw4",
			"xps13-7390",
			"xps13-9350",
		},
		InvalidMachines: []string{
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
		ValidMachines: []string{
			"asus-ux301l",
			"hp-pavilion-15-cs-3037nl",
			"hp-zbook-power-16-inch-g11",
			"i7-1165G7",
			"i7-10510U",
			"mustang",
			"system76-addw4",
			"xps13-7390",
			"xps13-9350",
		},
		InvalidMachines: []string{
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
		ValidMachines: []string{
			"hp-pavilion-15-cs-3037nl",
			"i7-1165G7",
		},
		InvalidMachines: []string{
			"ampere-one-m-banshee-12",
			"ampere-one-siryn",
			"ampere-one-x-banshee-8",
			"asus-ux301l",
			"hp-proliant-rl300-gen11-altra",
			"hp-proliant-rl300-gen11-altra-max",
			"hp-zbook-power-16-inch-g11",
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
		ValidMachines: []string{
			"hp-zbook-power-16-inch-g11",
			"mustang",
			"system76-addw4",
			"xps13-9350",
		},
		InvalidMachines: []string{
			"ampere-one-m-banshee-12",
			"ampere-one-siryn",
			"ampere-one-x-banshee-8",
			"asus-ux301l",
			"hp-pavilion-15-cs-3037nl",
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

	"cuda-generic": {
		ValidMachines: []string{
			"system76-addw4",
		},
		InvalidMachines: []string{
			"ampere-one-m-banshee-12",
			"ampere-one-siryn",
			"ampere-one-x-banshee-8",
			"asus-ux301l",
			"hp-pavilion-15-cs-3037nl", // Not enough vram
			"hp-proliant-rl300-gen11-altra",
			"hp-proliant-rl300-gen11-altra-max",
			"hp-zbook-power-16-inch-g11", // nvidia drivers not installed
			"i7-1165G7",
			"i7-2600k+arc-a580",
			"i7-10510U",
			"mustang",
			//"orange-pi-rv2",
			"raspberry-pi-5",
			"raspberry-pi-5+hailo-8",
			"xps13-7390",
			"xps13-9350",
		},
	},

	"intel-gpu": {
		ValidMachines: []string{
			"hp-zbook-power-16-inch-g11",
			"i7-2600k+arc-a580",
			"mustang",
			"system76-addw4",
			"xps13-9350",
		},
		InvalidMachines: []string{
			"ampere-one-m-banshee-12",
			"ampere-one-siryn",
			"ampere-one-x-banshee-8",
			"asus-ux301l", // has intel gpu, but clinfo not working
			"hp-pavilion-15-cs-3037nl",
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
		ValidMachines: []string{
			"hp-zbook-power-16-inch-g11",
			"xps13-9350",
		},
		InvalidMachines: []string{
			"ampere-one-m-banshee-12",
			"ampere-one-siryn",
			"ampere-one-x-banshee-8",
			"asus-ux301l",
			"hp-pavilion-15-cs-3037nl",
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
	for stackName, testSet := range validInvalidSets {
		for _, hwName := range testSet.ValidMachines {
			t.Run(stackName+" == "+hwName, func(t *testing.T) {
				testValidHw(t, stackName, hwName)
			})
		}

		for _, hwName := range testSet.InvalidMachines {
			t.Run(stackName+" != "+hwName, func(t *testing.T) {
				testInvalidHw(t, stackName, hwName)
			})
		}
	}
}

func testValidHw(t *testing.T, stackName string, hwName string) {
	stackManifestFile := fmt.Sprintf("../../test_data/engines/%s/engine.yaml", stackName)

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
	stackManifestFile := fmt.Sprintf("../../test_data/engines/%s/engine.yaml", stackName)

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
