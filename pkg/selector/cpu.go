package selector

import (
	"slices"

	"github.com/canonical/ml-snap-utils/pkg/types"
)

func checkCpus(stackDevice types.StackDevice, cpus []types.CpuInfo) (int, error) {
	cpusScore := 0

iterateCpus:
	for _, cpu := range cpus {
		cpuScore := WeightCpu

		// Vendor
		if stackDevice.VendorId != nil {
			if *stackDevice.VendorId == cpu.VendorId {
				cpuScore += WeightCpuVendor // vendor matched
			} else {
				continue
			}
		}

		// TODO
		// architecture
		// cpu count
		// Family
		// CpuModel

		// Flags
		for _, flag := range stackDevice.Flags {
			if !slices.Contains(cpu.Flags, flag) {
				continue iterateCpus
			}
			cpuScore += WeightCpuFlag
		}

		// Only add this CPU's score if it passed all the filters
		cpusScore += cpuScore
	}

	return cpusScore, nil
}
