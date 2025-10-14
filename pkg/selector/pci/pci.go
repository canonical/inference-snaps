package pci

import (
	"fmt"
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/selector/weights"
	"github.com/canonical/inference-snaps-cli/pkg/types"
	"github.com/canonical/go-snapctl"
)

func Match(device engines.Device, pcis []types.PciDevice) (int, []string, error) {
	var reasons []string
	maxDeviceScore := 0

	for _, pciDevice := range pcis {
		deviceScore, deviceReasons, err := checkPciDevice(device, pciDevice)
		reasons = append(reasons, deviceReasons...)
		if err != nil {
			return 0, reasons, err
		}

		if deviceScore > 0 {
			if deviceScore > maxDeviceScore {
				maxDeviceScore = deviceScore
			}
		}
	}

	return maxDeviceScore, reasons, nil
}

func checkPciDevice(device engines.Device, pciDevice types.PciDevice) (int, []string, error) {
	var reasons []string
	currentDeviceScore := 0

	// Device type: tpu, npu, gpu, etc
	if device.Type != "" {
		match := checkType(device.Type, pciDevice)
		if match {
			currentDeviceScore += weights.PciDeviceType
		} else {
			reasons = append(reasons, fmt.Sprintf("pci device type mismatch: %s", device.Type))
			return 0, reasons, nil
		}
	}

	// Prefer dGPU above iGPU
	// PCI devices on bus 0 are considered internal, and anything else external/discrete
	if pciDevice.BusNumber > 0 {
		currentDeviceScore += weights.PciDeviceExternal
	}

	if device.VendorId != nil {
		if *device.VendorId == pciDevice.VendorId {
			currentDeviceScore += weights.PciVendorId
		} else {
			reasons = append(reasons, fmt.Sprintf("pci vendor id mismatch: %04x", *device.VendorId))
			return 0, reasons, nil
		}

		// A model ID is only unique per vendor ID namespace. Only check it if the vendor is a match
		if device.DeviceId != nil {
			if *device.DeviceId == pciDevice.DeviceId {
				currentDeviceScore += weights.PciDeviceId
			} else {
				reasons = append(reasons, fmt.Sprintf("pci device id mismatch: %04x", *device.DeviceId))
				return 0, reasons, nil
			}
		}
	}

	// Check additional properties
	if hasAdditionalProperties(device) {
		propsScore, propReasons, err := checkProperties(device, pciDevice)
		reasons = append(reasons, propReasons...)
		if err != nil {
			return 0, reasons, err
		}
		if propsScore > 0 {
			currentDeviceScore += propsScore
		} else {
			return 0, reasons, nil
		}
	}

	// Check drivers
	for _, connection := range device.SnapConnections {
		connected, err := checkSnapConnection(connection)
		if err != nil {
			return 0, reasons, fmt.Errorf("error checking snap connection %q: %v", connection, err)
		}
		if !connected {
			reasons = append(reasons, fmt.Sprintf("%q is not connected", connection))
			return 0, reasons, nil
		}
	}

	return currentDeviceScore, reasons, nil
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
