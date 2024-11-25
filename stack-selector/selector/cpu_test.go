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
		Models:       nil,
	}

	result := checkCpu(stackDevice, hwInfoCpu)
	if !result {
		t.Fatal("CPU vendor should match")
	}

	vendorId = "AuthenticAMD"

	result = checkCpu(stackDevice, hwInfoCpu)
	if result {
		t.Fatal("CPU vendor should NOT match")
	}

}
