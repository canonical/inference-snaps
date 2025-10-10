package pci

import (
	"testing"

	"github.com/canonical/famous-models-cli/pkg/engines"
	"github.com/canonical/famous-models-cli/pkg/types"
)

func TestCheckGpuVendor(t *testing.T) {
	gpuVendorId := types.HexInt(0xb33f)

	hwInfoGpu := types.PciDevice{
		DeviceClass:          0x0300,
		VendorId:             gpuVendorId,
		DeviceId:             0,
		SubvendorId:          nil,
		SubdeviceId:          nil,
		AdditionalProperties: map[string]string{
			//VRam:              nil,
			//ComputeCapability: nil,
		},
	}

	device := engines.Device{
		Type:     "gpu",
		Bus:      "pci",
		VendorId: &gpuVendorId,
	}

	score, reasons, err := checkPciDevice(device, hwInfoGpu)
	if err != nil {
		t.Error(err)
	}
	if score == 0 {
		t.Fatalf("GPU vendor should match: %v", reasons)
	}

	// Same value, upper case string
	gpuVendorId = types.HexInt(0xB33F)
	score, reasons, err = checkPciDevice(device, hwInfoGpu)
	if err != nil {
		t.Error(err)
	}
	if score == 0 {
		t.Fatalf("GPU vendor should match: %v", reasons)
	}

	gpuVendorId = types.HexInt(0x1337)
	score, reasons, err = checkPciDevice(device, hwInfoGpu)
	if err != nil {
		t.Error(err)
	}
	if score > 0 {
		t.Fatal("GPU vendor should NOT match")
	}
}

func TestCheckGpuVram(t *testing.T) {

	hwInfoGpu := types.PciDevice{
		DeviceClass: 0x0300,
		VendorId:    0x0,
		DeviceId:    0x0,
		SubvendorId: nil,
		SubdeviceId: nil,
		AdditionalProperties: map[string]string{
			"vram": "5000000000",
		},
	}

	requiredVram := "4G"
	device := engines.Device{
		Type:     "gpu",
		Bus:      "pci",
		VendorId: nil,
		VRam:     &requiredVram,
	}

	score, reasons, err := checkPciDevice(device, hwInfoGpu)
	if err != nil {
		t.Error(err)
	}
	if score == 0 {
		t.Fatalf("GPU vram should be enough: %v", reasons)
	}

	requiredVram = "24G"
	score, reasons, err = checkPciDevice(device, hwInfoGpu)
	if err != nil {
		t.Error(err)
	}
	if score > 0 {
		t.Fatal("GPU vram should NOT be enough")
	}
}

func TestCheckNpuDriver(t *testing.T) {
	npuVendorId := types.HexInt(0x8086)
	npuDeviceId := types.HexInt(0x643e)

	hwInfo := types.PciDevice{
		DeviceClass: 0x1200,
		VendorId:    npuVendorId,
		DeviceId:    npuDeviceId,
		SubvendorId: nil,
		SubdeviceId: nil,
	}

	device := engines.Device{
		Bus:             "pci",
		VendorId:        &npuVendorId,
		DeviceId:        &npuDeviceId,
		SnapConnections: []string{"intel-npu", "npu-libs"},
	}

	score, reasons, err := checkPciDevice(device, hwInfo)
	if err != nil {
		t.Error(err)
	}
	if score == 0 {
		t.Fatalf("NPU with driver should match: %v", reasons)
	}

	// TODO test the negative case
}
