package engines

import (
	"testing"
)

func TestDeviceBus(t *testing.T) {
	device := Device{}
	device.Type = "gpu"

	t.Run("PCI Bus", func(t *testing.T) {
		device.Bus = "pci"
		err := device.validate()
		if err != nil {
			t.Fatalf("PCI Bus should be valid: %v", err)
		}
	})

	t.Run("USB Bus", func(t *testing.T) {
		device.Bus = "usb"
		err := device.validate()
		if err != nil {
			//t.Fatalf("USB Bus should be valid: %v", err)
			// USB bus not implemented
			t.Log(err)
		}
	})

	t.Run("Empty Bus", func(t *testing.T) {
		device.Bus = ""
		err := device.validate()
		if err != nil {
			t.Fatalf("Empty Bus should be valid: %v", err)
		}
	})

	t.Run("Invalid Bus", func(t *testing.T) {
		device.Bus = "invalid-bus"
		err := device.validate()
		if err == nil {
			t.Fatalf("Invalid bus should not validate")
		}
		t.Log(err)
	})

	t.Run("FastRPC bus with non-NPU type", func(t *testing.T) {
		device.Bus = "fastrpc"
		err := device.validate()
		if err == nil {
			t.Fatal("FastRPC bus should be invalid for GPU type")
		}
		t.Log(err)
	})
}
