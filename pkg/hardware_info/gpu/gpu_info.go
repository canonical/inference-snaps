package gpu

import (
	"errors"
	"fmt"
	"os"

	"github.com/canonical/ml-snap-utils/pkg/constants"
	"github.com/canonical/ml-snap-utils/pkg/hardware_info/pci"
	"github.com/canonical/ml-snap-utils/pkg/types"
)

func Info(friendlyNames bool) ([]types.Gpu, error) {
	pciDevices, err := pci.PciDevices(friendlyNames)
	if err != nil {
		return nil, err
	}

	return pciGpus(pciDevices)
}

func pciGpus(pciDevices []pci.PciDevice) ([]types.Gpu, error) {
	var gpus []types.Gpu

	for _, device := range pciDevices {
		// 00 01 - legacy VGA devices
		// 03 xx - display controllers
		if device.DeviceClass == 0x0001 || device.DeviceClass&0xFF00 == 0x0300 {
			var gpu types.Gpu
			gpu.VendorId = fmt.Sprintf("0x%04x", device.VendorId)
			gpu.DeviceId = fmt.Sprintf("0x%04x", device.DeviceId)
			if device.SubvendorId != nil {
				subVendorId := fmt.Sprintf("0x%04x", *device.SubvendorId)
				gpu.SubvendorId = &subVendorId
			}
			if device.SubdeviceId != nil {
				subDeviceId := fmt.Sprintf("0x%04x", *device.SubdeviceId)
				gpu.SubdeviceId = &subDeviceId
			}

			gpu.VendorName = device.VendorName
			gpu.DeviceName = device.DeviceName
			gpu.SubvendorName = device.SubvendorName
			gpu.SubdeviceName = device.SubdeviceName

			vram, err := getVRam(device)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error getting VRAM info for GPU:", err)
			}
			gpu.VRam = vram

			computeCapability, err := getComputeCapability(device)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error getting compute capability for GPU:", err)
			}
			gpu.ComputeCapability = computeCapability

			gpus = append(gpus, gpu)
		}
	}
	return gpus, nil
}

func getVRam(pciDevice pci.PciDevice) (*uint64, error) {
	switch pciDevice.VendorId {
	case constants.PciVendorAmd:
		return amdVram(pciDevice)
	case constants.PciVendorNvidia:
		return nvidiaVram(pciDevice)
	case constants.PciVendorIntel:
		return nil, errors.New("vram lookup for Intel GPU not implemented")
	default:
		return nil, errors.New("unknown GPU, not looking up vram")
	}
}

func getComputeCapability(pciDevice pci.PciDevice) (*string, error) {
	if pciDevice.VendorId == constants.PciVendorNvidia {
		return nvidiaComputeCapability(pciDevice)
	}
	// For other vendors we do not look up the Compute Capability
	return nil, nil
}
