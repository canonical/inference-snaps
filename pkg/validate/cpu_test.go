package validate

import (
	"testing"

	"github.com/canonical/stack-utils/pkg/constants"
	"github.com/canonical/stack-utils/pkg/types"
)

func TestCpuArchitecture(t *testing.T) {
	device := types.StackDevice{Type: "cpu"}

	t.Run("cpu arch amd64", func(t *testing.T) {
		architecture := constants.Amd64
		device.Architecture = &architecture

		err := cpu(device)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("cpu arch arm64", func(t *testing.T) {
		architecture := constants.Arm64
		device.Architecture = &architecture

		err := cpu(device)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("cpu arch invalid", func(t *testing.T) {
		architecture := "invalid-arch"
		device.Architecture = &architecture

		err := cpu(device)
		if err == nil {
			t.Fatal("CPU architecture should be invalid")
		}
		t.Log(err)
	})
}

func TestCpuAmd64Fields(t *testing.T) {
	architecture := constants.Amd64
	device := types.StackDevice{Type: "cpu", Architecture: &architecture}
	manufacturer := "My Manufacturer"
	device.ManufacturerId = &manufacturer

	t.Run("cpu amd64 valid fields", func(t *testing.T) {
		device.Flags = []string{"one", "two"}

		err := cpu(device)
		if err != nil {
			t.Fatalf("amd64 cpu fields should be valid: %v", err)
		}
	})

	t.Run("cpu amd64 invalid fields", func(t *testing.T) {
		device.Features = []string{"one", "two"}

		err := cpu(device)
		if err == nil {
			t.Fatal("amd64 cpu should not have features")
		}
		t.Log(err)
	})
}

func TestCpuArm64Fields(t *testing.T) {
	architecture := constants.Arm64
	device := types.StackDevice{Type: "cpu", Architecture: &architecture}
	implementer := types.HexInt(0x41)
	device.ImplementerId = &implementer

	t.Run("cpu arm64 valid fields", func(t *testing.T) {
		device.Features = []string{"one", "two"}

		err := cpu(device)
		if err != nil {
			t.Fatalf("arm64 cpu fields should be valid: %v", err)
		}
	})

	t.Run("cpu arm64 invalid fields", func(t *testing.T) {
		device.Flags = []string{"one", "two"}

		err := cpu(device)
		if err == nil {
			t.Fatal("arm64 cpu should not have flags")
		}
		t.Log(err)
	})
}
