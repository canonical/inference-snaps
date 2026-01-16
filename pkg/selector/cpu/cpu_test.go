package cpu

import (
	"strings"
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/constants"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/types"
)

func TestCheckCpuVendor(t *testing.T) {
	manufacturerId := "GenuineIntel"
	architecture := constants.Amd64
	device := engines.Device{
		Type:           "cpu",
		Bus:            "",
		Architecture:   &architecture,
		ManufacturerId: &manufacturerId,
	}

	hwInfoCpus := []types.CpuInfo{{
		Architecture:   constants.Amd64,
		ManufacturerId: manufacturerId,
	}}

	score, issues := Match(device, hwInfoCpus)
	if len(issues) != 0 {
		t.Fatalf("CPU vendor should match: %v", strings.Join(issues, ","))
	}

	manufacturerId = "AuthenticAMD"

	score, issues = Match(device, hwInfoCpus)
	if len(issues) == 0 || score > 0 {
		t.Fatal("CPU vendor should NOT match")
	}

}

func TestCheckCpuFlags(t *testing.T) {
	manufacturerId := "GenuineIntel"
	architecture := constants.Amd64
	device := engines.Device{
		Type:           "cpu",
		Bus:            "",
		Architecture:   &architecture,
		ManufacturerId: &manufacturerId,
		Flags:          []string{"avx2"},
	}

	hwInfoCpus := []types.CpuInfo{{
		Architecture:   constants.Amd64,
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
