package pci

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/canonical/go-snapctl"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/selector/weights"
	"github.com/canonical/inference-snaps-cli/pkg/types"
)

func Match(device engines.Device, pcis []types.PciDevice) (int, error) {
	maxDeviceScore := 0

	if len(pcis) == 0 {
		return 0, fmt.Errorf("no pci device on host system")
	}

	for _, pciDevice := range pcis {
		deviceScore, err := checkPciDevice(device, pciDevice)
		if err != nil {
			if os.Getenv("VERBOSE") == "true" {
				fmt.Printf("    %v\n", err)
			}
		}

		if deviceScore > 0 {
			if deviceScore > maxDeviceScore {
				maxDeviceScore = deviceScore
			}
		}
	}

	if maxDeviceScore == 0 {
		deviceJson, _ := json.Marshal(device)
		return 0, fmt.Errorf("device not found: %v", string(deviceJson))
	}
	return maxDeviceScore, nil
}

func checkPciDevice(device engines.Device, pciDevice types.PciDevice) (int, error) {
	currentDeviceScore := 0

	// Device type: tpu, npu, gpu, etc
	if device.Type != "" {
		match := checkType(device.Type, pciDevice)
		if match {
			currentDeviceScore += weights.PciDeviceType
		} else {
			return 0, fmt.Errorf("device class 0x%04x not of required type %s", pciDevice.DeviceClass, device.Type)
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
			return 0, fmt.Errorf("vendor id mismatch: 0x%04x", pciDevice.VendorId)
		}

		// A model ID is only unique per vendor ID namespace. Only check it if the vendor is a match
		if device.DeviceId != nil {
			if *device.DeviceId == pciDevice.DeviceId {
				currentDeviceScore += weights.PciDeviceId
			} else {
				return 0, fmt.Errorf("device id mismatch: 0x%04x", pciDevice.DeviceId)
			}
		}
	}

	// Check additional properties
	if hasAdditionalProperties(device) {
		propsScore, err := checkProperties(device, pciDevice)
		if err != nil {
			return 0, err
		}
		if propsScore > 0 {
			currentDeviceScore += propsScore
		}
	}

	// Check drivers
	for _, connection := range device.SnapConnections {
		connected, err := checkSnapConnection(connection)
		if err != nil {
			return 0, fmt.Errorf("error checking snap connection %q: %v", connection, err)
		}
		if !connected {
			return 0, fmt.Errorf("%q is not connected", connection)
		}
	}

	return currentDeviceScore, nil
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
