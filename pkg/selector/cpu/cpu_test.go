package cpu

import (
	"strings"
	"testing"

	"github.com/canonical/inference-snaps-cli/v2/pkg/constants"
	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/canonical/lscompute/pkg/machine"
	"github.com/canonical/lscompute/pkg/machine/cpu"
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

	hwInfoCpus := []cpu.CpuInfo{{
		Architecture:   constants.Amd64,
		ManufacturerId: manufacturerId,
	}}

	score, issues := Match(device, &machine.MachineInfo{Cpus: hwInfoCpus})
	if len(issues) != 0 {
		t.Fatalf("CPU vendor should match: %v", strings.Join(issues, ","))
	}

	manufacturerId = "AuthenticAMD"

	score, issues = Match(device, &machine.MachineInfo{Cpus: hwInfoCpus})
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

	hwInfoCpus := []cpu.CpuInfo{{
		Architecture:   constants.Amd64,
		ManufacturerId: manufacturerId,
		Flags:          []string{"avx2"},
	}}

	score, issues := Match(device, &machine.MachineInfo{Cpus: hwInfoCpus})
	if len(issues) != 0 {
		t.Fatalf("CPU flags should match: %v", strings.Join(issues, ","))
	}

	device.Flags = []string{"avx512"}

	score, issues = Match(device, &machine.MachineInfo{Cpus: hwInfoCpus})
	if len(issues) == 0 || score > 0 {
		t.Fatal("CPU flags should NOT match")
	}

}
