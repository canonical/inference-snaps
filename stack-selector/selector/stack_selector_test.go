package selector

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"gopkg.in/yaml.v2"
	"katemoss/common"
)

var hwInfoFiles = []string{
	"../test_data/hardware_info/amd-ryzen7-5700g.json",
	"../test_data/hardware_info/amd-ryzen9-7900.json",
	"../test_data/hardware_info/cbrd-i5-1350pe.json",
	"../test_data/hardware_info/hp-dl380p-gen8.json",
	"../test_data/hardware_info/xeon-6138.json",
	"../test_data/hardware_info/xps13-gen10.json",
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

func TestCpuFlagsAvx2(t *testing.T) {
	file, err := os.Open("../test_data/hardware_info/amd-ryzen7-5700g.json")
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

	data, err = os.ReadFile("../test_data/stacks/llamacpp-avx2/stack.yaml")
	if err != nil {
		t.Fatal(err)
	}

	var stack common.Stack
	err = yaml.Unmarshal(data, &stack)
	if err != nil {
		t.Fatal(err)
	}

	// Valid hardware for stack
	result, err := checkStack(hardwareInfo, stack)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", result)

	file, err = os.Open("../test_data/hardware_info/hp-dl380p-gen8.json")
	if err != nil {
		t.Fatal(err)
	}

	data, err = io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(data, &hardwareInfo)
	if err != nil {
		t.Fatal(err)
	}

	// Invalid hardware for stack
	result, err = checkStack(hardwareInfo, stack)
	if err == nil {
		t.Fatal("Stack should not match if avx2 is not available")
	}
}

func TestCpuFlagsAvx512(t *testing.T) {
	file, err := os.Open("../test_data/hardware_info/amd-ryzen9-7900.json")
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

	data, err = os.ReadFile("../test_data/stacks/llamacpp-avx512/stack.yaml")
	if err != nil {
		t.Fatal(err)
	}

	var stack common.Stack
	err = yaml.Unmarshal(data, &stack)
	if err != nil {
		t.Fatal(err)
	}

	// Valid hardware for stack
	_, err = checkStack(hardwareInfo, stack)
	if err != nil {
		t.Fatal(err)
	}

	file, err = os.Open("../test_data/hardware_info/hp-dl380p-gen8.json")
	if err != nil {
		t.Fatal(err)
	}

	data, err = io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(data, &hardwareInfo)
	if err != nil {
		t.Fatal(err)
	}

	// Invalid hardware for stack
	_, err = checkStack(hardwareInfo, stack)
	if err == nil {
		t.Fatal("Stack should not match if avx512 is not available")
	}
}

func TestNoCpuInHwInfo(t *testing.T) {
	hwInfo := common.HwInfo{
		// All fields are nil or zero
	}

	data, err := os.ReadFile("../test_data/stacks/llamacpp-avx512/stack.yaml")
	if err != nil {
		t.Fatal(err)
	}

	var stack common.Stack
	err = yaml.Unmarshal(data, &stack)
	if err != nil {
		t.Fatal(err)
	}

	// No memory in hardware info
	_, err = checkStack(hwInfo, stack)
	if err == nil {
		t.Fatal("No Memory in hardware_info should return err")
	}
	t.Log(err)

	hwInfo.Memory = &common.MemoryInfo{
		RamTotal:  17000000000,
		SwapTotal: 2000000000,
	}

	// No disk space in hardware info
	_, err = checkStack(hwInfo, stack)
	if err == nil {
		t.Fatal("No Disk space in hardware_info should return err")
	}
	t.Log(err)

	hwInfo.Disk = make(map[string]*common.DirStats)
	hwInfo.Disk["/"] = &common.DirStats{
		Avail: 6000000000,
	}

	// No CPU in hardware info
	_, err = checkStack(hwInfo, stack)
	if err == nil {
		t.Fatal("No CPU in hardware_info should return err")
	}
	t.Log(err)
}
