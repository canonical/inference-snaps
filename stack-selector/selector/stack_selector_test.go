package selector

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"katemoss/common"
)

var hwInfoFiles = []string{
	"../test_data/hardware_info/xps13-gen10.json",
	"../test_data/hardware_info/hp-dl380p-gen8.json",
}

func TestFindStack(t *testing.T) {
	for _, hwInfoFile := range hwInfoFiles {
		t.Run(hwInfoFile, func(t *testing.T) {
			file, err := os.Open(hwInfoFile)
			if err != nil {
				t.Fatal(err)
			}

			data, err := io.ReadAll(file)
			if err != nil {
				t.Fatal(err)
			}

			var hardwareInfo common.HwInfo
			err = json.Unmarshal(data, &hardwareInfo)
			if err != nil {
				t.Fatal(err)
			}

			result, err := FindStack(hardwareInfo, "../test_data/stacks")
			if err != nil {
				t.Fatal(err)
			}

			t.Logf("Found stack %s which installs %v", result.Name, result.Components)
		})
	}
}

func TestDiskCheck(t *testing.T) {
	dirStat := common.DirStats{
		Total: 0,
		Used:  0,
		Free:  0,
		Avail: 400000000,
	}
	hwInfo := common.HwInfo{}
	hwInfo.Disk = make(map[string]*common.DirStats)
	hwInfo.Disk["/"] = &dirStat
	hwInfo.Disk["/var/lib/snapd/snaps"] = &dirStat

	stackDisk := "300M"
	stack := common.Stack{DiskSpace: &stackDisk}

	result, err := checkStack(hwInfo, stack)
	if err != nil {
		t.Fatal(err)
	}
	if result == 0 {
		t.Fatal("disk should be enough")
	}

	dirStat.Avail = 100000000
	result, err = checkStack(hwInfo, stack)
	if err == nil {
		t.Fatal("Not enough disk should return err")
	}
	if result > 0 {
		t.Fatal("disk should NOT be enough")
	}
}

func TestMemoryCheck(t *testing.T) {
	hwInfo := common.HwInfo{
		Memory: &common.MemoryInfo{
			RamTotal:  200000000,
			SwapTotal: 200000000,
		},
	}

	stackMemory := "300M"
	stack := common.Stack{Memory: &stackMemory}

	result, err := checkStack(hwInfo, stack)
	if err != nil {
		t.Fatal(err)
	}
	if result == 0 {
		t.Fatal("memory should be enough")
	}

	hwInfo.Memory.RamTotal = 100000000
	result, err = checkStack(hwInfo, stack)
	if err == nil {
		t.Fatal("Not enough memory should return err")
	}
	if result > 0 {
		t.Fatal("memory should NOT be enough")
	}
}
