package selector

import (
	"errors"
	"fmt"
	"sort"

	"github.com/canonical/inference-snaps-cli/pkg/constants"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/selector/cpu"
	"github.com/canonical/inference-snaps-cli/pkg/selector/pci"
	"github.com/canonical/inference-snaps-cli/pkg/types"
	"github.com/canonical/inference-snaps-cli/pkg/utils"
)

var ErrorNoCompatibleEngine = errors.New("no compatible engines found")

func TopEngine(scoredEngines []engines.ScoredManifest) (*engines.ScoredManifest, error) {
	var compatibleEngines []engines.ScoredManifest

	for _, engine := range scoredEngines {
		if engine.Score > 0 && engine.Grade == "stable" {
			compatibleEngines = append(compatibleEngines, engine)
		}
	}

	if len(compatibleEngines) == 0 {
		return nil, ErrorNoCompatibleEngine
	}

	// Sort by score (high to low) and return highest match
	sort.Slice(compatibleEngines, func(i, j int) bool {
		return compatibleEngines[i].Score > compatibleEngines[j].Score
	})

	// Top engine is the highest score
	return &compatibleEngines[0], nil
}

func ScoreEngines(hardwareInfo *types.HwInfo, manifests []engines.Manifest) ([]engines.ScoredManifest, error) {
	var scoredEngines []engines.ScoredManifest

	for _, currentManifest := range manifests {
		score, reasons, err := checkEngine(hardwareInfo, currentManifest)
		if err != nil {
			return nil, err
		}

		scoredEngine := engines.ScoredManifest{
			Manifest:   currentManifest,
			Score:      score,
			Compatible: true,
		}

		if score == 0 {
			scoredEngine.Compatible = false
		}
		scoredEngine.CompatibilityIssues = append(scoredEngine.CompatibilityIssues, reasons...)

		scoredEngines = append(scoredEngines, scoredEngine)
	}

	return scoredEngines, nil
}

func checkEngine(hardwareInfo *types.HwInfo, manifest engines.Manifest) (int, []string, error) {
	engineScore := 0
	var reasons []string
	compatible := true

	// Enough memory
	if manifest.Memory != nil {
		requiredMemory, err := utils.StringToBytes(*manifest.Memory)
		if err != nil {
			return 0, nil, fmt.Errorf("failed to parse required memory: %v", err)

		} else if hardwareInfo.Memory.TotalRam == 0 {
			// If the TotalRam field is the Go struct Zero value, it was never set.
			// We do not check swap for the Zero value, as swap can realistically be of size 0 bytes.
			return 0, nil, fmt.Errorf("total memory not reported by host system")

		} else if hardwareInfo.Memory.TotalRam+hardwareInfo.Memory.TotalSwap < requiredMemory {
			// Checking combination of ram and swap
			compatible = false
			reasons = append(reasons, fmt.Sprintf("host system memory too small"))

		} else {
			engineScore++
		}
	}

	// Enough disk space
	if manifest.DiskSpace != nil {
		requiredDisk, err := utils.StringToBytes(*manifest.DiskSpace)
		if err != nil {
			return 0, nil, fmt.Errorf("failed to parse required disk space: %v", err)

		} else if _, ok := hardwareInfo.Disk[constants.SnapStoragePath]; !ok {
			return 0, nil, fmt.Errorf("disk space not reported by host system")

		} else if hardwareInfo.Disk[constants.SnapStoragePath].Avail < requiredDisk {
			compatible = false
			reasons = append(reasons, "host system disk space too small")

		} else {
			engineScore++
		}
	}

	// Devices
	// all
	if len(manifest.Devices.Allof) > 0 {
		extraScore, issues := checkDevicesAll(hardwareInfo, manifest.Devices.Allof)
		if len(issues) > 0 {
			compatible = false
			reasons = append(reasons, issues...)
		} else {
			engineScore += extraScore
		}
	}

	// any
	if len(manifest.Devices.Anyof) > 0 {
		extraScore, issues := checkDevicesAny(hardwareInfo, manifest.Devices.Anyof)
		if len(issues) > 0 {
			compatible = false
			reasons = append(reasons, issues...)
		} else {
			engineScore += extraScore
		}
	}

	if !compatible {
		engineScore = 0
	}

	return engineScore, reasons, nil
}

func checkDevicesAll(hardwareInfo *types.HwInfo, devices []engines.Device) (int, []string) {
	var issues []string
	compatible := true
	extraScore := 0

	for i, _ := range devices {

		if devices[i].Type == "cpu" {
			cpuScore, deviceIssues := cpu.Match(devices[i], hardwareInfo.Cpus)
			if len(deviceIssues) > 0 {
				compatible = false
				devices[i].CompatibilityIssues = append(devices[i].CompatibilityIssues, deviceIssues...)
				issues = append(issues, "required cpu device not found")
			} else {
				extraScore += cpuScore
			}

		} else if devices[i].Bus == "usb" {
			// Not implemented
			compatible = false
			devices[i].CompatibilityIssues = append(devices[i].CompatibilityIssues, "usb device matching not implemented")
			issues = append(issues, "usb device matching not implemented")

		} else if devices[i].Bus == "" || devices[i].Bus == "pci" {
			// Fallback to PCI as default bus
			pciScore, pciIssues := pci.Match(devices[i], hardwareInfo.PciDevices)
			if len(pciIssues) > 0 {
				compatible = false
				devices[i].CompatibilityIssues = append(devices[i].CompatibilityIssues, pciIssues...)
				issues = append(issues, "required pci device not found")
			} else {
				extraScore += pciScore
			}
		}
	}

	if !compatible {
		extraScore = 0
	}

	return extraScore, issues
}

func checkDevicesAny(hardwareInfo *types.HwInfo, devices []engines.Device) (int, []string) {
	var issues []string
	compatible := true
	extraScore := 0

	devicesFound := 0

	for i, device := range devices {

		if device.Type == "cpu" {
			cpuScore, deviceIssues := cpu.Match(device, hardwareInfo.Cpus)
			if len(deviceIssues) > 0 {
				devices[i].CompatibilityIssues = append(device.CompatibilityIssues, deviceIssues...)
			} else {
				devicesFound++
				extraScore += cpuScore
			}

		} else if device.Bus == "usb" {
			compatible = false
			device.CompatibilityIssues = append(device.CompatibilityIssues, "device type usb not implemented")
			issues = append(issues, "usb device matching not implemented")

		} else if device.Bus == "" || device.Bus == "pci" {
			// Fallback to PCI as default bus
			pciScore, pciIssues := pci.Match(device, hardwareInfo.PciDevices)
			if len(pciIssues) > 0 {
				devices[i].CompatibilityIssues = append(device.CompatibilityIssues, pciIssues...)
			} else {
				devicesFound++
				extraScore += pciScore
			}
		}
	}

	// If any-of devices are defined, we need to find at least one
	if len(devices) > 0 && devicesFound == 0 {
		compatible = false
		issues = append(issues, "required device not found")
	}

	if !compatible {
		extraScore = 0
	}

	return extraScore, issues
}
