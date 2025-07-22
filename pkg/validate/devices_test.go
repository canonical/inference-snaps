package validate

import (
	"testing"

	"github.com/canonical/stack-utils/pkg/constants"
	"github.com/canonical/stack-utils/pkg/types"
)

func TestDeviceType(t *testing.T) {

	t.Run("Type CPU", func(t *testing.T) {
		arch := constants.Amd64
		device := types.StackDevice{
			Type:         "cpu",
			Architecture: &arch,
		}

		err := stackDevice(device)
		if err != nil {
			t.Fatalf("Type cpu should be valid: %v", err)
		}
	})
	t.Run("Type GPU", func(t *testing.T) {
		device := types.StackDevice{Type: "gpu"}

		err := stackDevice(device)
		if err != nil {
			t.Fatalf("Type gpu should be valid: %v", err)
		}
	})
	t.Run("Type NPU", func(t *testing.T) {
		device := types.StackDevice{Type: "npu"}

		err := stackDevice(device)
		if err != nil {
			t.Fatalf("Type npu should be valid: %v", err)
		}
	})
	t.Run("Type empty", func(t *testing.T) {
		device := types.StackDevice{Type: ""}

		err := stackDevice(device)
		if err != nil {
			t.Fatalf("Empty type should be valid: %v", err)
		}
	})
	t.Run("Type invalid", func(t *testing.T) {
		device := types.StackDevice{Type: "test"}

		err := stackDevice(device)
		if err == nil {
			t.Fatalf("Invalid type should be invalid: %v", err)
		}
		t.Log(err)
	})
}

func TestDeviceGpu(t *testing.T) {
	device := types.StackDevice{}
	device.Type = "gpu"
	device.Bus = ""

	t.Run("GPU valid fields", func(t *testing.T) {
		hexValue := types.HexInt(0xAA)
		device.VendorId = &hexValue
		device.DeviceId = &hexValue

		vram := "1G"
		device.VRam = &vram

		computeCap := "12.4"
		device.ComputeCapability = &computeCap

		err := stackDevice(device)
		if err != nil {
			t.Fatalf("GPU fields should be valid: %v", err)
		}
	})

	t.Run("GPU invalid fields", func(t *testing.T) {
		hexValue := types.HexInt(0xAA)
		device.VendorId = &hexValue
		device.DeviceId = &hexValue

		vram := "1G"
		device.VRam = &vram

		manufacturer := "test"
		device.ManufacturerId = &manufacturer

		err := stackDevice(device)
		if err == nil {
			t.Fatal("GPU fields should be invalid")
		}
		t.Log(err)
	})
}

func TestDeviceNpu(t *testing.T) {
	device := types.StackDevice{}
	device.Type = "npu"
	device.Bus = ""

	t.Run("NPU valid fields", func(t *testing.T) {
		hexValue := types.HexInt(0xAA)
		device.VendorId = &hexValue
		device.DeviceId = &hexValue

		err := stackDevice(device)
		if err != nil {
			t.Fatalf("NPU fields should be valid: %v", err)
		}
	})

	t.Run("NPU invalid fields", func(t *testing.T) {
		hexValue := types.HexInt(0xAA)
		device.VendorId = &hexValue
		device.DeviceId = &hexValue

		vram := "1G"
		device.VRam = &vram

		computeCap := "12.4"
		device.ComputeCapability = &computeCap

		err := stackDevice(device)
		if err == nil {
			t.Fatal("NPU fields should be invalid")
		}
		t.Log(err)
	})
}

func TestDeviceTypeless(t *testing.T) {
	device := types.StackDevice{}
	device.Type = ""
	device.Bus = "pci"

	t.Run("PCI valid fields", func(t *testing.T) {
		hexValue := types.HexInt(0xAA)
		device.VendorId = &hexValue
		device.DeviceId = &hexValue
		err := stackDevice(device)
		if err != nil {
			t.Fatalf("PCI device fields should be valid: %v", err)
		}
	})

	t.Run("PCI invalid fields", func(t *testing.T) {
		hexValue := types.HexInt(0xAA)
		device.VendorId = &hexValue
		device.DeviceId = &hexValue
		device.Features = []string{"one", "two"}
		err := stackDevice(device)
		if err == nil {
			t.Fatal("PCI device fields should be invalid")
		}
		t.Log(err)
	})
}
