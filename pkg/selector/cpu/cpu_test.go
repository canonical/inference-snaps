package cpu

import (
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/types"
)

func TestCheckCpuVendor(t *testing.T) {
	manufacturerId := "GenuineIntel"
	device := engines.Device{
		Type:           "cpu",
		Bus:            "",
		ManufacturerId: &manufacturerId,
	}

	hwInfoCpus := []types.CpuInfo{{
		Architecture:   "",
		ManufacturerId: manufacturerId,
	}}

	result, err := Match(device, hwInfoCpus)
	if err != nil {
		t.Fatalf("CPU vendor should match: %v", err)
	}

	manufacturerId = "AuthenticAMD"

	result, err = Match(device, hwInfoCpus)
	if err == nil || result > 0 {
		t.Fatal("CPU vendor should NOT match")
	}

}

func TestCheckCpuFlags(t *testing.T) {
	manufacturerId := "GenuineIntel"
	device := engines.Device{
		Type:           "cpu",
		Bus:            "",
		ManufacturerId: &manufacturerId,
		Flags:          []string{"avx2"},
	}

	hwInfoCpus := []types.CpuInfo{{
		Architecture:   "",
		ManufacturerId: manufacturerId,
		Flags:          []string{"avx2"},
	}}

	result, err := Match(device, hwInfoCpus)
	if err != nil {
		t.Fatalf("CPU flags should match: %v", err)
	}

	device.Flags = []string{"avx512"}

	result, err = Match(device, hwInfoCpus)
	if err == nil || result > 0 {
		t.Fatal("CPU flags should NOT match")
	}

}
