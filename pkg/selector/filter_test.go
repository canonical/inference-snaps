package selector

import (
	"os"
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/constants"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/types"
	"gopkg.in/yaml.v3"
)

/*
If the model snap has no engines defined, scoring should pass, but finding a top engine should not be possible.
*/
func TestFindTopEngineFromNone(t *testing.T) {
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

func TestDiskCheck(t *testing.T) {
	dirStat := types.DirStats{
		Total: 0,
		Avail: 400000000,
	}
	hwInfo := types.HwInfo{}
	hwInfo.Disk = make(map[string]types.DirStats)
	hwInfo.Disk["/"] = dirStat
	hwInfo.Disk["/var/lib/snapd/snaps"] = dirStat

	manifestDisk := "300M"
	engine := engines.Manifest{DiskSpace: &manifestDisk}

	_, report, err := checkEngine(&hwInfo, engine)
	if err != nil {
		t.Fatal(err)
	}
	if !report.EngineCompatible() {
		t.Fatalf("engine should be compatible: %+v", report)
	}

	dirStat = types.DirStats{
		Total: 0,
		Avail: 100000000,
	}
	hwInfo.Disk["/var/lib/snapd/snaps"] = dirStat
	_, report, err = checkEngine(&hwInfo, engine)
	if err != nil {
		t.Fatal(err)
	}
	if report.EngineCompatible() {
		t.Fatalf("engine should NOT be compatible: %+v", report)
	}
}

func TestMemoryCheck(t *testing.T) {
	hwInfo := types.HwInfo{
		Memory: types.MemoryInfo{
			TotalRam:  200000000,
			TotalSwap: 200000000,
		},
	}

	engineMemory := "300M"
	engine := engines.Manifest{Memory: &engineMemory}

	_, report, err := checkEngine(&hwInfo, engine)
	if err != nil {
		t.Fatal(err)
	}
	if !report.EngineCompatible() {
		t.Fatalf("engine should be compatible: %+v", report)
	}

	hwInfo.Memory.TotalRam = 100000000
	_, report, err = checkEngine(&hwInfo, engine)
	if err != nil {
		t.Fatal(err)
	}
	if report.EngineCompatible() {
		t.Fatalf("engine should NOT be compatible: %+v", report)
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

	// No memory in hardware info
	_, _, err = checkEngine(&hwInfo, currentEngine)
	if err == nil {
		t.Fatal("Missing Memory info in hardware_info should return an error")
	}

	hwInfo.Memory = types.MemoryInfo{
		TotalRam:  17000000000,
		TotalSwap: 2000000000,
	}

	// No disk space in hardware info
	_, _, err = checkEngine(&hwInfo, currentEngine)
	if err == nil {
		t.Fatal("Missing Disk space info in hardware_info should return an error")
	}

	hwInfo.Disk = make(map[string]types.DirStats)
	hwInfo.Disk[constants.SnapStoragePath] = types.DirStats{
		Avail: 6000000000,
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
