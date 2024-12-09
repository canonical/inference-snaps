package selector

import (
	"slices"

	"katemoss/common"
)

func checkCpus(stackDevice common.StackDevice, cpu common.CpuInfo) (float64, error) {
	cpuScore := 0.0

	// Vendor
	if stackDevice.VendorId != nil {
		if *stackDevice.VendorId == cpu.Vendor {
			cpuScore += 1.0 // vendor matched
		} else {
			return 0, nil
		}
	}

	// TODO
	// architecture
	// cpu count

	for _, cpuModel := range cpu.Models {
		modelScore, err := checkCpuModel(cpuModel, stackDevice)
		if err != nil {
			return 0, err
		}
		if modelScore > 0 {
			// At the first matching CPU model we stop and return
			return cpuScore + modelScore, nil
		}
	}

	// If we get here, we checked all the CPU models and none were matches
	return 0, nil
}

// Apply the same "filter" logic as we have for the GPUs. See checkGpus() and checkGpu().
func checkCpuModel(cpuModel common.CpuModel, stackDevice common.StackDevice) (float64, error) {
	// Each CPU that matches increases the score by 1
	score := 1.0

	// Flags
	for _, flag := range stackDevice.Flags {
		if !slices.Contains(cpuModel.Flags, flag) {
			return 0, nil
		}
		score += 0.1 // flag matched has a score of 0.1
	}

	// TODO
	// Family
	// CpuModel

	// If we get here, all the filters passed and the device is a match
	return score, nil
}
