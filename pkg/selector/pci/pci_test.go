package pci

import (
	"strings"
	"testing"

	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/canonical/inference-snaps-cli/v2/pkg/utils"
	"github.com/canonical/lscompute/pkg/machine/device/pci"
)

func TestCheckGpuVendor(t *testing.T) {
	gpuVendorId := utils.HexInt(0xb33f)

	hwInfoGpu := pci.Device{
		DeviceClass:          0x0300,
		VendorId:             uint16(gpuVendorId),
		DeviceId:             0,
		SubvendorId:          nil,
		SubdeviceId:          nil,
		AdditionalProperties: map[string]string{
			//VRam:              nil,
			//ComputeCapability: nil,
		},
	}

	testPciDevices := []pciDevice{{Device: hwInfoGpu}}

	device := engines.Device{
		Type:     "gpu",
		Bus:      "pci",
		VendorId: &gpuVendorId,
	}

	availableDevices := filterPciDevices(testPciDevices, device.VendorId, device.DeviceId)
	_, scoreIssues := scorePciDevices(device, availableDevices)
	if len(scoreIssues) != 0 {
		t.Fatalf("GPU vendor should match: %s", strings.Join(scoreIssues, ", "))
	}

	// Same value, upper case string
	gpuVendorId = utils.HexInt(0xB33F)
	availableDevices = filterPciDevices(testPciDevices, device.VendorId, device.DeviceId)
	_, scoreIssues = scorePciDevices(device, availableDevices)
	if len(scoreIssues) != 0 {
		t.Fatalf("GPU vendor should match: %s", strings.Join(scoreIssues, ", "))
	}

	gpuVendorId = utils.HexInt(0x1337)
	availableDevices = filterPciDevices(testPciDevices, device.VendorId, device.DeviceId)
	_, scoreIssues = scorePciDevices(device, availableDevices)
	if len(scoreIssues) == 0 {
		t.Fatalf("GPU vendor should NOT match")
	}
}

func TestCheckGpuVram(t *testing.T) {

	hwInfoGpu := pci.Device{
		DeviceClass: 0x0300,
		VendorId:    0x0,
		DeviceId:    0x0,
		SubvendorId: nil,
		SubdeviceId: nil,
		AdditionalProperties: map[string]string{
			"vram": "5000000000",
		},
	}

	testPciDevices := []pciDevice{{Device: hwInfoGpu}}

	requiredVram := "4G"
	device := engines.Device{
		Type:     "gpu",
		Bus:      "pci",
		VendorId: nil,
		VRam:     &requiredVram,
	}

	availableDevices := filterPciDevices(testPciDevices, device.VendorId, device.DeviceId)
	scoredDevices, scoreIssues := scorePciDevices(device, availableDevices)
	if len(scoreIssues) != 0 {
		t.Fatalf("GPU vram should be enough: %s", strings.Join(scoreIssues, ", "))
	}

	requiredVram = "24G"
	availableDevices = filterPciDevices(testPciDevices, device.VendorId, device.DeviceId)
	scoredDevices, scoreIssues = scorePciDevices(device, availableDevices)
	if len(scoreIssues) == 0 || scoredDevices[0].Score != 0 {
		t.Fatalf("GPU vram should NOT be enough")
	}
}

func TestCheckComputeCapability(t *testing.T) {

	hwInfoGpu := pci.Device{
		DeviceClass: 0x0300,
		AdditionalProperties: map[string]string{
			"compute-capability": "6.1",
		},
	}

	testPciDevices := []pciDevice{{Device: hwInfoGpu}}

	requiredComputeCapability := ">=6.0, <7.0"
	device := engines.Device{
		Type:              "gpu",
		Bus:               "pci",
		VendorId:          nil,
		ComputeCapability: &requiredComputeCapability,
	}

	availableDevices := filterPciDevices(testPciDevices, device.VendorId, device.DeviceId)
	scoredDevices, scoreIssues := scorePciDevices(device, availableDevices)
	if len(scoreIssues) != 0 {
		t.Fatalf("GPU compute-capability should match: %s", strings.Join(scoreIssues, ", "))
	}

	requiredComputeCapability = ">=7.0"
	availableDevices = filterPciDevices(testPciDevices, device.VendorId, device.DeviceId)
	scoredDevices, scoreIssues = scorePciDevices(device, availableDevices)
	if len(scoreIssues) == 0 || scoredDevices[0].Score != 0 {
		t.Fatalf("GPU compute-capability should NOT match")
	}
}

func TestCheckNpuDriver(t *testing.T) {
	npuVendorId := utils.HexInt(0x8086)
	npuDeviceId := utils.HexInt(0x643e)

	hwInfo := pci.Device{
		DeviceClass: 0x1200,
		VendorId:    uint16(npuVendorId),
		DeviceId:    uint16(npuDeviceId),
		SubvendorId: nil,
		SubdeviceId: nil,
	}

	testPciDevices := []pciDevice{{Device: hwInfo}}

	device := engines.Device{
		Bus:             "pci",
		VendorId:        &npuVendorId,
		DeviceId:        &npuDeviceId,
		SnapConnections: []string{"intel-npu", "npu-libs"},
	}

	availableDevices := filterPciDevices(testPciDevices, device.VendorId, device.DeviceId)
	_, scoreIssues := scorePciDevices(device, availableDevices)
	if len(scoreIssues) != 0 {
		t.Fatalf("NPU with driver should match: %s", strings.Join(scoreIssues, ", "))
	}

	// TODO test the negative case
}

func TestCheckMicroarchitecture(t *testing.T) {

	hwInfoGpu := pci.Device{
		DeviceClass: 0x0300,
		AdditionalProperties: map[string]string{
			"microarchitecture": "gfx1152",
		},
	}

	testPciDevices := []pciDevice{{Device: hwInfoGpu}}

	requiredMicroarchitecture := "gfx1152"
	device := engines.Device{
		Type:              "gpu",
		Bus:               "pci",
		VendorId:          nil,
		Microarchitecture: &requiredMicroarchitecture,
	}

	availableDevices := filterPciDevices(testPciDevices, device.VendorId, device.DeviceId)
	scoredDevices, scoreIssues := scorePciDevices(device, availableDevices)
	if len(scoreIssues) != 0 {
		t.Fatalf("GPU microarchitecture should match: %s", strings.Join(scoreIssues, ", "))
	}

	requiredMicroarchitecture = "gfx2200"
	availableDevices = filterPciDevices(testPciDevices, device.VendorId, device.DeviceId)
	scoredDevices, scoreIssues = scorePciDevices(device, availableDevices)
	if len(scoreIssues) == 0 || scoredDevices[0].Score != 0 {
		t.Fatalf("GPU microarchitecture should NOT match")
	}
}
