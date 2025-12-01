package cpu

import (
	"fmt"
	"slices"
	"strings"

	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/selector/weights"
	"github.com/canonical/inference-snaps-cli/pkg/types"
)

/*
Match takes a Device with type CPU, and checks if it matches any of the CPU models reported for the system.
A score, a string slice with reasons and an error are returned. If there is a matching CPU on the system, the score will be positive and the error will be nil.
If no CPU is found, the score will be zero and there will be one or more reasons for the mismatch. In case of a runtime error, the error value will be non-nil.
*/
func Match(device engines.Device, cpus []types.CpuInfo) (int, error) {
	maxCpuScore := 0
	var reasons []string

	if cpus == nil {
		return 0, fmt.Errorf("no cpus on host system")
	}

	for _, cpu := range cpus {
		cpuScore, err := CheckCpu(device, cpu)

		if err != nil {
			reasons = append(reasons, err.Error())
		} else {
			if cpuScore > maxCpuScore {
				maxCpuScore = cpuScore
			}
		}
	}

	if maxCpuScore == 0 {
		return 0, fmt.Errorf("%s", strings.Join(reasons, ", "))
	}

	return maxCpuScore, nil
}

func CheckCpu(device engines.Device, cpu types.CpuInfo) (int, error) {
	cpuScore := weights.CpuDevice

	// architecture
	if device.Architecture != nil {
		if *device.Architecture == cpu.Architecture {
			// architecture matches - no additional weight
		} else {
			return 0, fmt.Errorf("cpu architecture mismatch: %s", cpu.Architecture)
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
			return 0, fmt.Errorf("cpu manufacturer id mismatch: %s", cpu.ManufacturerId)
		}
	}

	// amd64 flags
	for _, flag := range device.Flags {
		if !slices.Contains(cpu.Flags, flag) {
			return 0, fmt.Errorf("cpu flag not found: %s", flag)
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
			return 0, fmt.Errorf("cpu implementer id mismatch: %x", cpu.ImplementerId)
		}
	}

	// arm64 part number
	if device.PartNumber != nil {
		if *device.PartNumber == cpu.PartNumber {
			cpuScore += weights.CpuModel
		} else {
			return 0, fmt.Errorf("cpu part number mismatch: %x", cpu.PartNumber)
		}
	}

	// arm64 features
	for _, feature := range device.Features {
		if !slices.Contains(cpu.Features, feature) {
			return 0, fmt.Errorf("cpu feature not found: %s", feature)
		}
		cpuScore += weights.CpuFlag
	}

	return cpuScore, nil
}
