package cpu

import (
	"fmt"
	"slices"

	"github.com/canonical/inference-snaps-cli/pkg/constants"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/selector/weights"
	"github.com/canonical/inference-snaps-cli/pkg/types"
)

/*
Match takes a Device with type CPU, and checks if it matches any of the CPU models reported for the system.
A score, a string slice with reasons and an error are returned. If there is a matching CPU on the system, the score will be positive and the error will be nil.
If no CPU is found, the score will be zero and there will be one or more reasons for the mismatch. In case of a runtime error, the error value will be non-nil.
*/
func Match(manifestDevice engines.Device, hostCpus []types.CpuInfo) (maxCpuScore int, deviceIssues []string) {
	maxCpuScore = 0

	if hostCpus == nil {
		deviceIssues = append(deviceIssues, "no cpu found on host system")
	}

	for i, cpu := range hostCpus {
		cpuScore, cpuIssues := CheckCpu(manifestDevice, cpu)

		if len(cpuIssues) > 0 {
			if len(hostCpus) > 1 {
				for _, issue := range cpuIssues {
					deviceIssues = append(deviceIssues, fmt.Sprintf("cpu %d: %v", i, issue))
				}
			} else {
				deviceIssues = append(deviceIssues, cpuIssues...)
			}
		} else {
			if cpuScore > maxCpuScore {
				maxCpuScore = cpuScore
			}
		}
	}

	return
}

func CheckCpu(manifestDevice engines.Device, hostCpu types.CpuInfo) (cpuScore int, issues []string) {
	cpuScore = weights.CpuDevice

	// architecture
	if manifestDevice.Architecture != nil {
		if *manifestDevice.Architecture == hostCpu.Architecture {
			// architecture matches - no additional weight
		} else {
			issues = append(issues, fmt.Sprintf("architecture not %s", *manifestDevice.Architecture))
		}
	}

	/*
		amd64
	*/
	if hostCpu.Architecture == constants.Amd64 {
		// amd64 manufacturer ID
		if manifestDevice.ManufacturerId != nil {
			if *manifestDevice.ManufacturerId == hostCpu.ManufacturerId {
				cpuScore += weights.CpuVendor
			} else {
				issues = append(issues, fmt.Sprintf("manufacturer id mismatch: %s", hostCpu.ManufacturerId))
			}
		}

		// amd64 flags
		for _, flag := range manifestDevice.Flags {
			if slices.Contains(hostCpu.Flags, flag) {
				cpuScore += weights.CpuFlag
			} else {
				issues = append(issues, fmt.Sprintf("flag %s missing", flag))
			}
		}
	}

	/*
		arm64
	*/
	if hostCpu.Architecture == constants.Arm64 {
		// arm64 implementer ID
		if manifestDevice.ImplementerId != nil {
			if *manifestDevice.ImplementerId == hostCpu.ImplementerId {
				cpuScore += weights.CpuVendor
			} else {
				issues = append(issues, fmt.Sprintf("implementer id mismatch: %x", hostCpu.ImplementerId))
			}
		}

		// arm64 part number
		if manifestDevice.PartNumber != nil {
			if *manifestDevice.PartNumber == hostCpu.PartNumber {
				cpuScore += weights.CpuModel
			} else {
				issues = append(issues, fmt.Sprintf("part number mismatch: %x", hostCpu.PartNumber))
			}
		}

		// arm64 features
		for _, feature := range manifestDevice.Features {
			if slices.Contains(hostCpu.Features, feature) {
				cpuScore += weights.CpuFlag
			} else {
				issues = append(issues, fmt.Sprintf("feature not found: %s", feature))
			}
		}
	}

	if len(issues) > 0 {
		cpuScore = 0
	}

	return
}
