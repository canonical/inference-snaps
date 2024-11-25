package selector

import (
	"katemoss/common"
)

func checkCpus(device common.StackDevice, cpu common.CpuInfo) bool {
	// When the model changes to support multiple CPUs, iterate over them here.
	// for _, cpu := range cpus {

	if checkCpu(device, cpu) {
		return true
	}

	// If we get here none of the CPUs on the system is a match
	return false
}

// Apply the same "filter" logic as we have for the GPUs. See checkGpus() and checkGpu().
func checkCpu(device common.StackDevice, cpu common.CpuInfo) bool {
	// Architecture

	// Vendor
	if device.VendorId != nil && *device.VendorId != cpu.Vendor {
		return false
	}

	// TODO
	// Family
	// CpuModel
	// Flags

	// If we get here, all the filters passed and the device is a match
	return true
}
