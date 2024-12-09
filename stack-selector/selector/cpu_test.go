package selector

import (
	"testing"

	"katemoss/common"
)

func TestCheckCpu(t *testing.T) {
	vendorId := "GenuineIntel"
	stackDevice := common.StackDevice{
		Type:     "cpu",
		Bus:      nil,
		VendorId: &vendorId,
	}

	hwInfoCpu := common.CpuInfo{
		Architecture: "",
		CpuCount:     0,
		Vendor:       vendorId,
		Models: []common.CpuModel{
			common.CpuModel{},
		},
	}

	result, err := checkCpus(stackDevice, hwInfoCpu)
	if err != nil {
		t.Error(err)
	}
	if result == 0 {
		t.Fatal("CPU vendor should match")
	}

	vendorId = "AuthenticAMD"

	result, err = checkCpus(stackDevice, hwInfoCpu)
	if err != nil {
		t.Error(err)
	}
	if result > 0 {
		t.Fatal("CPU vendor should NOT match")
	}

}
