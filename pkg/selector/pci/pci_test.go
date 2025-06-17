package pci

import (
	"testing"

	"github.com/canonical/stack-utils/pkg/types"
)

func TestCheckGpuVendor(t *testing.T) {
	gpuVendorId := "0xb33f" // test lower case string

	hwInfoGpu := types.PciDevice{
		DeviceClass:          0x0300,
		VendorId:             0xb33f,
		DeviceId:             0,
		SubvendorId:          nil,
		SubdeviceId:          nil,
		AdditionalProperties: map[string]string{
			//VRam:              nil,
			//ComputeCapability: nil,
		},
	}

	stackDevice := types.StackDevice{
		Type:        "gpu",
		Bus:         "pci",
		PciVendorId: &gpuVendorId,
	}

	score, reasons, err := checkPciDevice(stackDevice, hwInfoGpu)
	if err != nil {
		t.Error(err)
	}
	if score == 0 {
		t.Fatalf("GPU vendor should match: %v", reasons)
	}

	// Same value, upper case string
	gpuVendorId = "0xB33F"
	score, reasons, err = checkPciDevice(stackDevice, hwInfoGpu)
	if err != nil {
		t.Error(err)
	}
	if score == 0 {
		t.Fatalf("GPU vendor should match: %v", reasons)
	}

	gpuVendorId = "0x1337"
	score, reasons, err = checkPciDevice(stackDevice, hwInfoGpu)
	if err != nil {
		t.Error(err)
	}
	if score > 0 {
		t.Fatal("GPU vendor should NOT match")
	}
}

func TestCheckGpuVram(t *testing.T) {

	hwInfoGpu := types.PciDevice{
		VendorId:    0x0,
		DeviceId:    0x0,
		SubvendorId: nil,
		SubdeviceId: nil,
		AdditionalProperties: map[string]string{
			"vram": "5000000000",
		},
	}

	stackVram := "4G"
	stackDevice := types.StackDevice{
		Type:        "gpu",
		Bus:         "pci",
		PciVendorId: nil,
		VRam:        &stackVram,
	}

	score, reasons, err := checkPciDevice(stackDevice, hwInfoGpu)
	if err != nil {
		t.Error(err)
	}
	if score == 0 {
		t.Fatalf("GPU vram should be enough: %v", reasons)
	}

	stackVram = "24G"
	score, reasons, err = checkPciDevice(stackDevice, hwInfoGpu)
	if err != nil {
		t.Error(err)
	}
	if score > 0 {
		t.Fatal("GPU vram should NOT be enough")
	}
}
