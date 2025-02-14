package selector

import (
	"testing"

	"github.com/canonical/ml-snap-utils/pkg/types"
)

func TestCheckGpuVendor(t *testing.T) {

	gpuVendorId := "b33f"

	hwInfoGpu := types.Gpu{
		VendorId:          gpuVendorId,
		VendorName:        nil,
		DeviceId:          "",
		DeviceName:        nil,
		SubvendorId:       nil,
		SubvendorName:     nil,
		SubdeviceId:       nil,
		SubdeviceName:     nil,
		VRam:              nil,
		ComputeCapability: nil,
	}

	stackDevice := types.StackDevice{
		Type:     "gpu",
		Bus:      nil,
		VendorId: &gpuVendorId,
	}

	result, err := gpuMatchesStack(hwInfoGpu, stackDevice)
	if err != nil {
		t.Error(err)
	}
	if !result {
		t.Fatal("GPU vendor should match")
	}

	// Same value, different case
	gpuVendorId = "B33F"
	result, err = gpuMatchesStack(hwInfoGpu, stackDevice)
	if err != nil {
		t.Error(err)
	}
	if !result {
		t.Fatal("GPU vendor should match")
	}

	gpuVendorId = "1337"
	result, err = gpuMatchesStack(hwInfoGpu, stackDevice)
	if err != nil {
		t.Error(err)
	}
	if result {
		t.Fatal("GPU vendor should NOT match")
	}
}

func TestCheckGpuVram(t *testing.T) {

	var vram uint64 = 5000000000

	hwInfoGpu := types.Gpu{
		VendorId:          "",
		VendorName:        nil,
		DeviceId:          "",
		DeviceName:        nil,
		SubvendorId:       nil,
		SubvendorName:     nil,
		SubdeviceId:       nil,
		SubdeviceName:     nil,
		VRam:              &vram,
		ComputeCapability: nil,
	}

	stackVram := "4G"
	stackDevice := types.StackDevice{
		Type:     "gpu",
		Bus:      nil,
		VendorId: nil,
		VRam:     &stackVram,
	}

	result, err := gpuMatchesStack(hwInfoGpu, stackDevice)
	if err != nil {
		t.Error(err)
	}
	if !result {
		t.Fatal("GPU vram should be enough")
	}

	stackVram = "24G"
	result, err = gpuMatchesStack(hwInfoGpu, stackDevice)
	if err != nil {
		t.Error(err)
	}
	if result {
		t.Fatal("GPU vram should NOT be enough")
	}
}
