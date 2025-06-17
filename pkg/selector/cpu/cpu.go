package cpu

import (
	"fmt"
	"slices"

	"github.com/canonical/stack-utils/pkg/selector/weights"
	"github.com/canonical/stack-utils/pkg/types"
)

/*
Match takes a Stack Manifest Device with type CPU, and checks if it matches any of the CPU models reported for the system.
A score, a string slice with reasons and an error are returned. If there is a matching CPU on the system, the score will be positive and the error will be nil.
If no CPU is found, the score will be zero and there will be one or more reasons for the mismatch. In case of a runtime error, the error value will be non-nil.
*/
func Match(stackDevice types.StackDevice, cpus []types.CpuInfo) (int, []string, error) {
	cpusScore := 0
	var reasons []string

iterateCpus:
	for _, cpu := range cpus {
		cpuScore := weights.CpuDevice

		// amd64 manufacturer ID
		if stackDevice.ManufacturerId != nil {
			if *stackDevice.ManufacturerId == cpu.ManufacturerId {
				cpuScore += weights.CpuVendor // vendor matched
			} else {
				reasons = append(reasons, fmt.Sprintf("Manufacturer ID does not match: %s = %s", *stackDevice.ManufacturerId, cpu.ManufacturerId))
				continue
			}
		}

		// arm64 implementer ID
		if stackDevice.ImplementerId != nil {
			if *stackDevice.ImplementerId == cpu.ImplementerId {
				cpuScore += weights.CpuVendor
			} else {
				reasons = append(reasons, fmt.Sprintf("Implementer ID does not match: %x = %x", *stackDevice.ImplementerId, cpu.ImplementerId))
				continue
			}
		}

		// TODO
		// architecture
		// cpu count
		// Family list
		// CpuModel list
		// TODO stackDevice.ModelName - see #48

		// Flags
		for _, flag := range stackDevice.Flags {
			if !slices.Contains(cpu.Flags, flag) {
				reasons = append(reasons, "Required flag not found: %s", flag)
				continue iterateCpus
			}
			cpuScore += weights.CpuFlag
		}

		// Only add this CPU's score if it passed all the filters
		cpusScore += cpuScore
	}

	return cpusScore, reasons, nil
}
