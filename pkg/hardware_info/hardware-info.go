package hardware_info

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/canonical/stack-utils/pkg/hardware_info/cpu"
	"github.com/canonical/stack-utils/pkg/hardware_info/disk"
	"github.com/canonical/stack-utils/pkg/hardware_info/memory"
	"github.com/canonical/stack-utils/pkg/hardware_info/pci"
	"github.com/canonical/stack-utils/pkg/types"
)

func Get(friendlyNames bool) (types.HwInfo, error) {
	var hwInfo types.HwInfo

	memoryInfo, err := memory.Info()
	if err != nil {
		return hwInfo, fmt.Errorf("error getting memory info: %v", err)
	}
	hwInfo.Memory = memoryInfo

	cpus, err := cpu.Info()
	if err != nil {
		return hwInfo, fmt.Errorf("error getting cpu info: %v", err)
	}
	hwInfo.Cpus = cpus

	diskInfo, err := disk.Info()
	if err != nil {
		return hwInfo, fmt.Errorf("error getting disk info: %v", err)
	}
	hwInfo.Disk = diskInfo

	pciDevices, err := pci.Devices(friendlyNames)
	if err != nil {
		return hwInfo, fmt.Errorf("error getting pci devices: %v", err)
	}
	hwInfo.PciDevices = pciDevices

	return hwInfo, nil
}

// GetFromRawData is mainly used during testing, but also from other packages, and therefore needs to be exported
func GetFromRawData(t *testing.T, device string, friendlyNames bool) (types.HwInfo, error) {
	var hwInfo types.HwInfo

	devicePath := "../../test_data/devices/" + device + "/"

	// memory
	procMemInfo, err := os.ReadFile(devicePath + "meminfo.txt")
	if err != nil {
		t.Fatal(err)
	}
	memInfo, err := memory.InfoFromRawData(string(procMemInfo))
	if err != nil {
		t.Fatal(err)
	}
	hwInfo.Memory = memInfo

	// disk
	dfInfo, err := os.ReadFile(devicePath + "disk.txt")
	if err != nil {
		t.Fatal(err)
	}
	diskInfo, err := disk.InfoFromRawData(string(dfInfo))
	if err != nil {
		t.Fatal(err)
	}
	hwInfo.Disk = diskInfo

	// cpu
	unameMachine, err := os.ReadFile(devicePath + "uname-m.txt")
	if err != nil {
		t.Fatal(err)
	}
	procCpuInfo, err := os.ReadFile(devicePath + "cpuinfo.txt")
	if err != nil {
		t.Fatal(err)
	}
	cpuInfo, err := cpu.InfoFromRawData(string(procCpuInfo), string(unameMachine))
	if err != nil {
		t.Fatal(err)
	}
	hwInfo.Cpus = cpuInfo

	// pci
	pciData, err := os.ReadFile(devicePath + "lspci.txt")
	if err != nil {
		t.Fatal(err)
	}
	pciDevices, err := pci.DevicesFromRawData(string(pciData), friendlyNames)
	if err != nil {
		t.Fatal(err)
	}
	hwInfo.PciDevices = pciDevices

	// Additional properties - we append these directly from a file, as we can not run the vendor specific tools on the machine
	addPropsFile := devicePath + "additional-properties.json"
	_, err = os.Stat(addPropsFile)
	if err != nil {
		if os.IsNotExist(err) {
			// File does not exist. Skipping additional properties
		} else {
			t.Fatalf("error checking file '%s': %v\n", addPropsFile, err)
		}
	} else {
		var addProps map[string]map[string]string
		addPropsData, err := os.ReadFile(devicePath + "additional-properties.json")
		if err != nil {
			t.Fatal(err)
		}
		err = json.Unmarshal(addPropsData, &addProps)
		if err != nil {
			t.Fatal(err)
		}
		for i, pciDevice := range hwInfo.PciDevices {
			if val, ok := addProps[pciDevice.Slot]; ok {
				hwInfo.PciDevices[i].AdditionalProperties = val
			}
		}
	}

	return hwInfo, nil
}
