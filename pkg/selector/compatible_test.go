package selector

import (
	"fmt"
	"os"
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/hardware_info"
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
			"lenovo-thinkpad-p16s",
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
			"lenovo-thinkpad-p16s",
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
			"lenovo-thinkpad-p16s",
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
			"lenovo-thinkpad-p16s",
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
			"lenovo-thinkpad-p16s",
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
			"lenovo-thinkpad-p16s",
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
			"lenovo-thinkpad-p16s",
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
			"lenovo-thinkpad-p16s",
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
			"lenovo-thinkpad-p16s",
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
			"lenovo-thinkpad-p16s",
		},
	},

	"amd-gpu": {
		ValidMachines: []string{
			"hp-zbook-i712850HX+RadeonPROW6600M",
		},
		InvalidMachines: []string{
			"ampere-one-m-banshee-12",
			"ampere-one-siryn",
			"ampere-one-x-banshee-8",
			"asus-ux301l",
			"hp-pavilion-15-cs-3037nl",
			"hp-proliant-rl300-gen11-altra",
			"hp-proliant-rl300-gen11-altra-max",
			"hp-zbook-power-16-inch-g11",
			"i5-3570k+arc-a580+gtx1080ti",
			"i7-1165G7",
			"i7-2600k+arc-a580",
			"i7-10510U",
			"mustang",
			"raspberry-pi-5",
			"raspberry-pi-5+hailo-8",
			"system76-addw4",
			"xps13-7390",
			"xps13-9350",
			"lenovo-thinkpad-p16s",
		},
	},
	"rocm-generic": {
		ValidMachines: []string{
			"lenovo-thinkpad-p16s",
		},
		InvalidMachines: []string{
			"ampere-one-m-banshee-12",
			"ampere-one-siryn",
			"ampere-one-x-banshee-8",
			"asus-ux301l",
			"hp-pavilion-15-cs-3037nl",
			"hp-proliant-rl300-gen11-altra",
			"hp-proliant-rl300-gen11-altra-max",
			"hp-zbook-power-16-inch-g11",
			"i5-3570k+arc-a580+gtx1080ti",
			"i7-1165G7",
			"i7-2600k+arc-a580",
			"i7-10510U",
			"mustang",
			"raspberry-pi-5",
			"raspberry-pi-5+hailo-8",
			"system76-addw4",
			"xps13-7390",
			"xps13-9350",
			"hp-zbook-i712850HX+RadeonPROW6600M",
		},
	},
}

func TestEngine(t *testing.T) {
	for engineName, testSet := range validInvalidSets {
		for _, hwName := range testSet.ValidMachines {
			t.Run(engineName+" == "+hwName, func(t *testing.T) {
				testValidHw(t, engineName, hwName)
			})
		}

		for _, hwName := range testSet.InvalidMachines {
			t.Run(engineName+" != "+hwName, func(t *testing.T) {
				testInvalidHw(t, engineName, hwName)
			})
		}
	}
}

func testValidHw(t *testing.T, engineName string, hwName string) {
	manifestFile := fmt.Sprintf("../../test_data/engines/%s/%s", engineName, engines.ManifestFilename)

	hardwareInfo, err := hardware_info.GetFromRawData(hwName, true, "../../test_data")
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(manifestFile)
	if err != nil {
		t.Fatal(err)
	}

	var manifest engines.Manifest
	err = yaml.Unmarshal(data, &manifest)
	if err != nil {
		t.Fatal(err)
	}

	// Valid hardware for engine
	score, report, err := checkEngine(hardwareInfo, manifest)
	if err != nil {
		t.Fatal(err)
	}
	if !report.EngineCompatible() {
		t.Fatalf("Engine should match: %+v", report)
	} else if score == 0 {
		t.Fatalf("A compatible engine should have a non-zero score")
	}

	t.Logf("Matching score: %d", score)
}

func testInvalidHw(t *testing.T, engineName string, hwName string) {
	manifestFile := fmt.Sprintf("../../test_data/engines/%s/%s", engineName, engines.ManifestFilename)

	hardwareInfo, err := hardware_info.GetFromRawData(hwName, true, "../../test_data")
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(manifestFile)
	if err != nil {
		t.Fatal(err)
	}

	var manifest engines.Manifest
	err = yaml.Unmarshal(data, &manifest)
	if err != nil {
		t.Fatal(err)
	}

	score, report, err := checkEngine(hardwareInfo, manifest)
	if err != nil {
		t.Fatal(err)
	}
	if report.EngineCompatible() {
		t.Fatalf("Engine should not match: %s", hwName)
	} else if score != 0 {
		t.Fatalf("An incompatible engine should have a score of 0")
	}

	t.Logf("Matching score: %d", score)
}
