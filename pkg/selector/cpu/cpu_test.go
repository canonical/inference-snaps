package cpu

import (
	"testing"

	"github.com/canonical/ml-snap-utils/pkg/types"
)

func TestCheckCpuVendor(t *testing.T) {
	vendorId := "GenuineIntel"
	stackDevice := types.StackDevice{
		Type:     "cpu",
		Bus:      nil,
		VendorId: &vendorId,
	}

	hwInfoCpus := []types.CpuInfo{{
		Architecture: "",
		VendorId:     vendorId,
	}}

	result, err := checkCpus(stackDevice, hwInfoCpus)
	if err != nil {
		t.Error(err)
	}
	if result == 0 {
		t.Fatal("CPU vendor should match")
	}

	vendorId = "AuthenticAMD"

	result, err = checkCpus(stackDevice, hwInfoCpus)
	if err != nil {
		t.Error(err)
	}
	if result > 0 {
		t.Fatal("CPU vendor should NOT match")
	}

}

func TestCheckCpuFlags(t *testing.T) {
	vendorId := "GenuineIntel"
	stackDevice := types.StackDevice{
		Type:     "cpu",
		Bus:      nil,
		VendorId: &vendorId,
		Flags:    []string{"avx2"},
	}

	hwInfoCpus := []types.CpuInfo{{
		Architecture: "",
		VendorId:     vendorId,
		Flags:        []string{"avx2"},
	}}

	result, err := checkCpus(stackDevice, hwInfoCpus)
	if err != nil {
		t.Error(err)
	}
	if result == 0 {
		t.Fatal("CPU flags should match")
	}

	stackDevice.Flags = []string{"avx512"}

	result, err = checkCpus(stackDevice, hwInfoCpus)
	if err != nil {
		t.Error(err)
	}
	if result > 0 {
		t.Fatal("CPU flags should NOT match")
	}

}
