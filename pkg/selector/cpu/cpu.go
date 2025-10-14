package cpu

import (
	"fmt"
	"slices"

	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/selector/weights"
	"github.com/canonical/inference-snaps-cli/pkg/types"
)

/*
Match takes a Device with type CPU, and checks if it matches any of the CPU models reported for the system.
A score, a string slice with reasons and an error are returned. If there is a matching CPU on the system, the score will be positive and the error will be nil.
If no CPU is found, the score will be zero and there will be one or more reasons for the mismatch. In case of a runtime error, the error value will be non-nil.
*/
func Match(device engines.Device, cpus []types.CpuInfo) (int, []string, error) {
	cpusScore := 0
	var reasons []string

iterateCpus:
	for _, cpu := range cpus {
		cpuScore := weights.CpuDevice

		// architecture
		if device.Architecture != nil {
			if *device.Architecture == cpu.Architecture {
				// architecture matches - no additional weight
			} else {
				reasons = append(reasons, fmt.Sprintf("cpu architecture mismatch: %s", *device.Architecture))
				continue
			}
		}

		/*
			amd64
		*/

		// amd64 manufacturer ID
		if device.ManufacturerId != nil {
			if *device.ManufacturerId == cpu.ManufacturerId {
				cpuScore += weights.CpuVendor
			} else {
				reasons = append(reasons, fmt.Sprintf("cpu manufacturer id mismatch: %s", *device.ManufacturerId))
				continue
			}
		}

		// amd64 flags
		for _, flag := range device.Flags {
			if !slices.Contains(cpu.Flags, flag) {
				reasons = append(reasons, fmt.Sprintf("cpu flag not found: %s", flag))
				continue iterateCpus
			}
			cpuScore += weights.CpuFlag
		}

		/*
			arm64
		*/

		// arm64 implementer ID
		if device.ImplementerId != nil {
			if *device.ImplementerId == cpu.ImplementerId {
				cpuScore += weights.CpuVendor
			} else {
				reasons = append(reasons, fmt.Sprintf("cpu implementer id mismatch: %x", *device.ImplementerId))
				continue
			}
		}

		// arm64 part number
		if device.PartNumber != nil {
			if *device.PartNumber == cpu.PartNumber {
				cpusScore += weights.CpuModel
			} else {
				reasons = append(reasons, fmt.Sprintf("cpu part number mismatch: %x", *device.PartNumber))
				continue
			}
		}

		// arm64 features
		for _, feature := range device.Features {
			if !slices.Contains(cpu.Features, feature) {
				reasons = append(reasons, fmt.Sprintf("cpu feature not found: %s", feature))
				continue iterateCpus
			}
			cpuScore += weights.CpuFlag
		}

		// Only add this CPU's score if it passed all the filters
		cpusScore += cpuScore
	}

	return cpusScore, reasons, nil
}
