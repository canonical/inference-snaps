package selector

import (
	"errors"
	"sort"

	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/selector/cpu"
	"github.com/canonical/inference-snaps-cli/pkg/selector/pci"
	"github.com/canonical/inference-snaps-cli/pkg/types"
)

var ErrorNoCompatibleEngine = errors.New("no compatible engines found")

func TopEngine(scoredEngines []engines.ScoredManifest) (*engines.ScoredManifest, error) {
	var compatibleEngines []engines.ScoredManifest

	for _, engine := range scoredEngines {
		if engine.Score > 0 && (engine.Experimental == nil || !*engine.Experimental) {
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
		score, compatibilityReport, err := checkEngine(hardwareInfo, currentManifest)
		if err != nil {
			return nil, err
		}

		scoredEngine := engines.ScoredManifest{
			Manifest:            currentManifest,
			Score:               score,
			CompatibilityReport: compatibilityReport,
		}

		scoredEngines = append(scoredEngines, scoredEngine)
	}

	return scoredEngines, nil
}

func checkEngine(hardwareInfo *types.HwInfo, manifest engines.Manifest) (int, engines.CompatibilityReport, error) {
	engineScore := 0
	compatibilityReport := engines.CompatibilityReport{
		CompatibleMemory:  true,
		CompatibleDisk:    true,
		CompatibleDevices: true,
	}

	// Devices

	// all
	if len(manifest.Devices.Allof) > 0 {
		deviceCompatibilityScore := scoreDevicesAll(hardwareInfo, manifest.Devices.Allof)
		if deviceCompatibilityScore == 0 {
			compatibilityReport.CompatibleDevices = false
		} else {
			engineScore += deviceCompatibilityScore
		}
	}

	// any
	if len(manifest.Devices.Anyof) > 0 {
		deviceCompatibilityScore := scoreDevicesAny(hardwareInfo, manifest.Devices.Anyof)
		if deviceCompatibilityScore == 0 {
			compatibilityReport.CompatibleDevices = false
		} else {
			engineScore += deviceCompatibilityScore
		}
	}

	if !compatibilityReport.EngineCompatible() {
		engineScore = 0
	}

	return engineScore, compatibilityReport, nil
}

func scoreDevicesAll(hardwareInfo *types.HwInfo, devices []engines.Device) int {
	compatible := true
	compatibilityScore := 0

	for i, _ := range devices {

		if devices[i].Type == "cpu" {
			cpuScore, deviceIssues := cpu.Match(devices[i], hardwareInfo.Cpus)
			if len(deviceIssues) > 0 {
				compatible = false
				devices[i].CompatibilityIssues = append(devices[i].CompatibilityIssues, deviceIssues...)
			} else {
				compatibilityScore += cpuScore
			}

		} else if devices[i].Bus == "usb" {
			// Not implemented
			compatible = false
			devices[i].CompatibilityIssues = append(devices[i].CompatibilityIssues, "usb device matching not implemented")

		} else if devices[i].Bus == "" || devices[i].Bus == "pci" {
			// Fallback to PCI as default bus
			pciScore, pciIssues := pci.Match(devices[i], hardwareInfo.PciDevices)
			if len(pciIssues) > 0 {
				compatible = false
				devices[i].CompatibilityIssues = append(devices[i].CompatibilityIssues, pciIssues...)
			} else {
				compatibilityScore += pciScore
			}
		}
	}

	if !compatible {
		compatibilityScore = 0
	}

	return compatibilityScore
}

func scoreDevicesAny(hardwareInfo *types.HwInfo, devices []engines.Device) int {
	compatible := true
	compatibilityScore := 0
	devicesFound := 0

	for i, device := range devices {

		if device.Type == "cpu" {
			cpuScore, deviceIssues := cpu.Match(device, hardwareInfo.Cpus)
			if len(deviceIssues) > 0 {
				devices[i].CompatibilityIssues = append(device.CompatibilityIssues, deviceIssues...)
			} else {
				devicesFound++
				compatibilityScore += cpuScore
			}

		} else if device.Bus == "usb" {
			compatible = false
			device.CompatibilityIssues = append(device.CompatibilityIssues, "usb device matching not implemented")

		} else if device.Bus == "" || device.Bus == "pci" {
			// Fallback to PCI as default bus
			pciScore, pciIssues := pci.Match(device, hardwareInfo.PciDevices)
			if len(pciIssues) > 0 {
				devices[i].CompatibilityIssues = append(device.CompatibilityIssues, pciIssues...)
			} else {
				devicesFound++
				compatibilityScore += pciScore
			}
		}
	}

	// If any-of devices are defined, we need to find at least one
	if len(devices) > 0 && devicesFound == 0 {
		compatible = false
	}

	if !compatible {
		compatibilityScore = 0
	}

	return compatibilityScore
}
