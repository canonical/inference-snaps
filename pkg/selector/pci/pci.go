package pci

import (
	"fmt"
	"testing"

	"github.com/canonical/go-snapctl"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/selector/weights"
	"github.com/canonical/inference-snaps-cli/pkg/types"
)

func Match(manifestDevice engines.Device, hostPciDevices []types.PciDevice) (maxDeviceScore int, deviceIssues []string) {
	maxDeviceScore = 0

	if len(hostPciDevices) == 0 {
		deviceIssues = append(deviceIssues, "no pci devices on host system")
		return
	}

	availableDevices := filterPciDevices(hostPciDevices, manifestDevice.VendorId, manifestDevice.DeviceId)
	scoredDevices, scoreIssues := scorePciDevices(manifestDevice, availableDevices)

	for _, pci := range scoredDevices {
		if pci.Score > maxDeviceScore {
			maxDeviceScore = pci.Score
		}
	}
	if maxDeviceScore == 0 {
		deviceIssues = append(deviceIssues, scoreIssues...)
		return
	}

	return
}

// filterPciDevices returns all PCI devices from the provided list, where the Vendor ID and the Device ID match.
//
// Filtering does not return compatibility issues. If we did, an engine with N device on a machine with M pci devices,
// would print NxM issues. These will all read "vendor id mismatch" or "device id mismatch" for each NxM combination.
// In the end the reason is just "device not found".
func filterPciDevices(pciDevices []types.PciDevice, vendorId *types.HexInt, deviceId *types.HexInt) []types.PciDevice {
	var foundDevices []types.PciDevice
	for _, pciDevice := range pciDevices {
		include := true

		if vendorId != nil {
			if *vendorId != pciDevice.VendorId {
				include = false
			} else {
				// A model ID is only unique per vendor ID namespace. Only check it if the vendor is a match
				if deviceId != nil {
					if *deviceId != pciDevice.DeviceId {
						include = false
					}
				}
			}
		}

		if include {
			foundDevices = append(foundDevices, pciDevice)
		}
	}
	return foundDevices
}

// scorePciDevices takes a list of host pci devices, which should already be filtered based on Vendor ID and Device ID,
// performs scoring on all the devices, and returns a list of scored devices
func scorePciDevices(manifestDevice engines.Device, hostPciDevices []types.PciDevice) ([]types.PciDevice, []string) {
	var issues []string

	if len(hostPciDevices) == 0 {
		issues = append(issues, "device not found")
	}

	for i, pciDevice := range hostPciDevices {
		deviceScore, deviceIssues := scorePciDevice(manifestDevice, pciDevice)

		hostPciDevices[i].Score = deviceScore
		for _, issue := range deviceIssues {
			issues = append(issues, fmt.Sprintf("pci %s: %s", pciDevice.Slot, issue))
		}
	}
	return hostPciDevices, issues
}

func scorePciDevice(manifestDevice engines.Device, hostPciDevice types.PciDevice) (deviceScore int, issues []string) {
	deviceScore = 0

	// Device type: tpu, npu, gpu, etc
	if manifestDevice.Type != "" {
		match := checkType(manifestDevice.Type, hostPciDevice)
		if match {
			deviceScore += weights.PciDeviceType
		} else {
			deviceScore = 0
			// The device type does not map directly to a device class. We use a decision tree to check if the device class
			// and subclass fall in known ranges per device type. This makes printing a direct reason here difficult.
			// The message "device class 0x%04x not of required type %s" was chosen as best compromise here.
			issues = append(issues,
				fmt.Sprintf("device class 0x%04x not of required type %s",
					hostPciDevice.DeviceClass, manifestDevice.Type))
			return
		}
	}

	// Prefer dGPU over iGPU
	// PCI devices on bus 0 are considered internal, and anything else external/discrete
	if hostPciDevice.BusNumber > 0 {
		deviceScore += weights.PciDeviceExternal
	}

	// Check additional properties
	if hasAdditionalProperties(manifestDevice) {
		propsScore, err := checkProperties(manifestDevice, hostPciDevice)
		if err != nil {
			deviceScore = 0
			issues = append(issues, err.Error())
			return
		}
		deviceScore += propsScore
	}

	// Check drivers
	for _, connection := range manifestDevice.SnapConnections {
		connected, err := checkSnapConnection(connection)
		if err != nil {
			deviceScore = 0
			issues = append(issues, fmt.Sprintf("error checking snap connection %q: %v", connection, err))
			return
		}
		if !connected {
			deviceScore = 0
			issues = append(issues, fmt.Sprintf("%q is not connected", connection))
			return
		}
	}

	return deviceScore, nil
}

func checkType(requiredType string, pciDevice types.PciDevice) bool {
	if requiredType == "gpu" {
		// 00 01 - legacy VGA devices
		// 03 xx - display controllers
		if pciDevice.DeviceClass == 0x0001 || pciDevice.DeviceClass&0xFF00 == 0x0300 {
			return true
		}
	}

	/*
		Base class 0x12 = Processing Accelerator - Intel Lunar Lake NPU identifies as this class
		Base class 0x0B = Processor, Sub class 0x40 = Co-Processor - Hailo PCI devices identify as this class
	*/
	if requiredType == "npu" || requiredType == "tpu" {
		if pciDevice.DeviceClass&0xFF00 == 0x1200 {
			// Processing accelerator
			return true
		}
		if pciDevice.DeviceClass == 0x0B40 {
			// Coprocessor
			return true
		}
	}

	return false
}

func checkSnapConnection(connection string) (bool, error) {
	if testing.Testing() {
		// Tests do not necessarily run inside a snap
		// Stub out and always return true for all connections
		return true, nil
	}
	return snapctl.IsConnected(connection).Run()
}
