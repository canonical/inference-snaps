package selector

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/selector/cpu"
	"github.com/canonical/inference-snaps-cli/pkg/selector/pci"
	"github.com/canonical/inference-snaps-cli/pkg/types"
	"github.com/canonical/inference-snaps-cli/pkg/utils"
	"gopkg.in/yaml.v3"
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

func LoadManifestsFromDir(manifestsDir string) ([]engines.Manifest, error) {
	var manifests []engines.Manifest

	// Sanitize dir path
	if !strings.HasSuffix(manifestsDir, "/") {
		manifestsDir += "/"
	}

	// Iterate engines
	files, err := os.ReadDir(manifestsDir)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", manifestsDir, err)
	}

	for _, file := range files {
		// Engines dir should contain a dir per engine
		if !file.IsDir() {
			continue
		}

		fileName := manifestsDir + file.Name() + "/engine.yaml"
		data, err := os.ReadFile(fileName)
		if err != nil {
			return nil, fmt.Errorf("%s: %s", fileName, err)
		}

		var manifest engines.Manifest
		err = yaml.Unmarshal(data, &manifest)
		if err != nil {
			return nil, fmt.Errorf("%s: %s", manifestsDir, err)
		}

		manifests = append(manifests, manifest)
	}
	return manifests, nil
}

func LoadManifestFromDir(manifestsDir, engineName string) (*engines.Manifest, error) {
	// Sanitize dir path
	if !strings.HasSuffix(manifestsDir, "/") {
		manifestsDir += "/"
	}

	fileName := manifestsDir + engineName + "/engine.yaml"
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", fileName, err)
	}

	var manifest engines.Manifest
	err = yaml.Unmarshal(data, &manifest)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", manifestsDir, err)
	}

	return &manifest, nil
}

func ScoreEngines(hardwareInfo *types.HwInfo, manifests []engines.Manifest) ([]engines.ScoredManifest, error) {
	var scoredEngines []engines.ScoredManifest

	for _, currentManifest := range manifests {
		if os.Getenv("VERBOSE") == "true" {
			fmt.Printf("Checking engine: %s\n", currentManifest.Name)
		}
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
			if os.Getenv("VERBOSE") == "true" {
				fmt.Printf("Engine not compatible: %s\n", strings.Join(reasons, ", "))
			}
			scoredEngine.Compatible = false
		} else {
			if os.Getenv("VERBOSE") == "true" {
				fmt.Printf("Engine compatible\n")
			}
		}
		scoredEngine.Notes = append(scoredEngine.Notes, reasons...)

		scoredEngines = append(scoredEngines, scoredEngine)
	}

	return scoredEngines, nil
}

func checkEngine(hardwareInfo *types.HwInfo, manifest engines.Manifest) (int, []string, error) {
	engineScore := 0
	var reasons []string

	// Enough memory
	if manifest.Memory != nil {
		requiredMemory, err := utils.StringToBytes(*manifest.Memory)
		if err != nil {
			return 0, reasons, err
		}

		if hardwareInfo.Memory.TotalRam == 0 {
			return 0, reasons, fmt.Errorf("system can't have zero ram")
		}

		// Checking combination of ram and swap
		if hardwareInfo.Memory.TotalRam+hardwareInfo.Memory.TotalSwap < requiredMemory {
			reasons = append(reasons, fmt.Sprintf("memory: system memory too small"))
			return 0, reasons, nil
		}
		engineScore++
	}

	// Enough disk space
	if manifest.DiskSpace != nil {
		requiredDisk, err := utils.StringToBytes(*manifest.DiskSpace)
		if err != nil {
			return 0, reasons, err
		}
		if _, ok := hardwareInfo.Disk["/var/lib/snapd/snaps"]; !ok {
			return 0, reasons, fmt.Errorf("disk space not reported by hardware info")
		}
		if hardwareInfo.Disk["/var/lib/snapd/snaps"].Avail < requiredDisk {
			reasons = append(reasons, fmt.Sprintf("disk: system disk space too small"))
			return 0, reasons, nil
		}
		engineScore++
	}

	// Devices
	// all
	if len(manifest.Devices.Allof) > 0 {

		extraScore, err := checkDevicesAll(hardwareInfo, manifest.Devices.Allof)
		if err != nil {
			reasons = append(reasons, err.Error())
			return 0, reasons, nil
		} else {
			engineScore += extraScore
		}
	}

	// any
	if len(manifest.Devices.Anyof) > 0 {
		extraScore, err := checkDevicesAny(hardwareInfo, manifest.Devices.Anyof)
		if err != nil {
			reasons = append(reasons, err.Error())
			return 0, reasons, nil
		} else {
			engineScore += extraScore
		}
	}

	return engineScore, reasons, nil
}

func checkDevicesAll(hardwareInfo *types.HwInfo, devices []engines.Device) (int, error) {
	devicesFound := 0
	extraScore := 0

	for _, device := range devices {
		if os.Getenv("VERBOSE") == "true" {
			jsonBytes, _ := json.Marshal(device)
			fmt.Printf("  Checking for all-of required device: %s\n", string(jsonBytes))
		}

		if device.Type == "cpu" {
			cpuScore, err := cpu.Match(device, hardwareInfo.Cpus)
			if err != nil {
				return 0, fmt.Errorf("required cpu device not found: %v", err)
			}
			extraScore += cpuScore
			devicesFound++

		} else if device.Bus == "usb" {
			// Not implemented
			return 0, fmt.Errorf("usb device matching not implemented")

		} else if device.Bus == "" || device.Bus == "pci" {
			// Fallback to PCI as default bus
			pciScore, err := pci.Match(device, hardwareInfo.PciDevices)
			if err != nil {
				return 0, fmt.Errorf("required pci device under all-of not found: %v", err)
			}
			extraScore += pciScore
			devicesFound++
		}

		if os.Getenv("VERBOSE") == "true" {
			fmt.Printf("  Device found\n")
		}
	}

	return extraScore, nil
}

func checkDevicesAny(hardwareInfo *types.HwInfo, devices []engines.Device) (int, error) {
	devicesFound := 0
	extraScore := 0
	var reasons []string

	for _, device := range devices {
		if os.Getenv("VERBOSE") == "true" {
			jsonBytes, _ := json.Marshal(device)
			fmt.Printf("  Checking for any-of required device: %s\n", string(jsonBytes))
		}

		if device.Type == "cpu" {
			cpuScore, err := cpu.Match(device, hardwareInfo.Cpus)
			if err != nil {
				if os.Getenv("VERBOSE") == "true" {
					fmt.Println("  " + err.Error())
				}
				reasons = append(reasons, err.Error())
			} else {
				if os.Getenv("VERBOSE") == "true" {
					fmt.Printf("  Device found\n")
				}
				devicesFound++
				extraScore += cpuScore
			}

		} else if device.Bus == "usb" {
			return 0, fmt.Errorf("devices any-of: device type usb not implemented")

		} else if device.Bus == "" || device.Bus == "pci" {
			// Fallback to PCI as default bus
			pciScore, err := pci.Match(device, hardwareInfo.PciDevices)
			if err != nil {
				if os.Getenv("VERBOSE") == "true" {
					fmt.Println("  " + err.Error())
				}
				reasons = append(reasons, err.Error())
			} else {
				if os.Getenv("VERBOSE") == "true" {
					fmt.Printf("  Device found\n")
				}
				devicesFound++
				extraScore += pciScore
			}
		}

	}

	// If any-of devices are defined, we need to find at least one
	if len(devices) > 0 && devicesFound == 0 {

		if os.Getenv("VERBOSE") == "true" {
			for _, reason := range reasons {
				fmt.Println("  " + reason)
			}
		}
		return 0, fmt.Errorf("no devices under anyof found")
	}

	return extraScore, nil
}
