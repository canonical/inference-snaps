package hardware_info

import (
	"github.com/canonical/ml-snap-utils/pkg/hardware_info/cpu"
	"github.com/canonical/ml-snap-utils/pkg/hardware_info/disk"
	"github.com/canonical/ml-snap-utils/pkg/hardware_info/memory"
	"github.com/canonical/ml-snap-utils/pkg/hardware_info/pci"
	"github.com/canonical/ml-snap-utils/pkg/types"
)

func Get(friendlyNames bool) (types.HwInfo, error) {
	var hwInfo types.HwInfo

	memoryInfo, err := memory.Info()
	if err != nil {
		return hwInfo, err
	}
	hwInfo.Memory = memoryInfo

	cpus, err := cpu.Info()
	if err != nil {
		return hwInfo, err
	}
	hwInfo.Cpus = cpus

	diskInfo, err := disk.Info()
	if err != nil {
		return hwInfo, err
	}
	hwInfo.Disk = diskInfo

	pciDevices, err := pci.Devices(friendlyNames)
	if err != nil {
		return hwInfo, err
	}
	hwInfo.PciDevices = pciDevices

	return hwInfo, nil
}
