package cpu

import (
	"strings"
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

	result, reasons, err := Match(device, hwInfoCpus)
	if err != nil {
		t.Error(err)
	}
	if result == 0 {
		t.Fatalf("CPU vendor should match: %s", strings.Join(reasons, ","))
	}

	manufacturerId = "AuthenticAMD"

	result, reasons, err = Match(device, hwInfoCpus)
	if err != nil {
		t.Error(err)
	}
	if result > 0 {
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

	result, reasons, err := Match(device, hwInfoCpus)
	if err != nil {
		t.Error(err)
	}
	if result == 0 {
		t.Fatalf("CPU flags should match: %s", strings.Join(reasons, ","))
	}

	device.Flags = []string{"avx512"}

	result, reasons, err = Match(device, hwInfoCpus)
	if err != nil {
		t.Error(err)
	}
	if result > 0 {
		t.Fatal("CPU flags should NOT match")
	}

}
