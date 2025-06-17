package selector

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/canonical/stack-utils/pkg/selector/cpu"
	"github.com/canonical/stack-utils/pkg/selector/pci"
	"github.com/canonical/stack-utils/pkg/types"
	"github.com/canonical/stack-utils/pkg/utils"
	"gopkg.in/yaml.v3"
)

func TopStack(scoredStacks []types.ScoredStack) (*types.ScoredStack, error) {
	var compatibleStacks []types.ScoredStack

	for _, stack := range scoredStacks {
		if stack.Score > 0 && stack.Grade == "stable" {
			compatibleStacks = append(compatibleStacks, stack)
		}
	}

	if len(compatibleStacks) == 0 {
		return nil, errors.New("no compatible stacks found")
	}

	// Sort by score (high to low) and return highest match
	sort.Slice(compatibleStacks, func(i, j int) bool {
		return compatibleStacks[i].Score > compatibleStacks[j].Score
	})

	// Top stack is highest score
	return &compatibleStacks[0], nil
}

func LoadStacksFromDir(stacksDir string) ([]types.Stack, error) {
	var stacks []types.Stack

	// Sanitise stack dir path
	if !strings.HasSuffix(stacksDir, "/") {
		stacksDir += "/"
	}

	// Iterate stacks
	files, err := os.ReadDir(stacksDir)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", stacksDir, err)
	}

	for _, file := range files {
		// Stacks dir should contain a dir per stack
		if !file.IsDir() {
			continue
		}

		data, err := os.ReadFile(stacksDir + file.Name() + "/stack.yaml")
		if err != nil {
			return nil, fmt.Errorf("%s: %s", stacksDir+file.Name(), err)
		}

		var currentStack types.Stack
		err = yaml.Unmarshal(data, &currentStack)
		if err != nil {
			return nil, fmt.Errorf("%s: %s", stacksDir, err)
		}

		stacks = append(stacks, currentStack)
	}
	return stacks, nil
}

func ScoreStacks(hardwareInfo types.HwInfo, stacks []types.Stack) ([]types.ScoredStack, error) {
	var scoredStacks []types.ScoredStack

	for _, currentStack := range stacks {
		score, reasons, err := checkStack(hardwareInfo, currentStack)
		if err != nil {
			return nil, err
		}

		scoredStack := types.ScoredStack{
			Stack:      currentStack,
			Score:      score,
			Compatible: true,
		}

		if score == 0 {
			scoredStack.Compatible = false
		}
		scoredStack.Notes = append(scoredStack.Notes, reasons...)

		scoredStacks = append(scoredStacks, scoredStack)
	}

	return scoredStacks, nil
}

func checkStack(hardwareInfo types.HwInfo, stack types.Stack) (int, []string, error) {
	stackScore := 0
	var reasons []string

	// Enough memory
	if stack.Memory != nil {
		requiredMemory, err := utils.StringToBytes(*stack.Memory)
		if err != nil {
			return 0, reasons, err
		}

		if hardwareInfo.Memory == nil {
			return 0, reasons, fmt.Errorf("no memory in hardware info")
		}

		// Checking combination of ram and swap
		if hardwareInfo.Memory.TotalRam+hardwareInfo.Memory.TotalSwap < requiredMemory {
			reasons = append(reasons, fmt.Sprintf("memory: system memory too small"))
			return 0, reasons, nil
		}
		stackScore++
	}

	// Enough disk space
	if stack.DiskSpace != nil {
		requiredDisk, err := utils.StringToBytes(*stack.DiskSpace)
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
		stackScore++
	}

	// Devices
	// all
	if len(stack.Devices.All) > 0 {
		extraScore, reasonsAll, err := checkDevicesAll(hardwareInfo, stack.Devices.All)
		for _, reason := range reasonsAll {
			reasons = append(reasons, "all: "+reason)
		}
		if err != nil {
			return 0, reasons, err
		}
		if extraScore == 0 {
			return 0, reasons, nil
		}
		stackScore += extraScore
	}

	// any
	if len(stack.Devices.Any) > 0 {
		extraScore, reasonsAny, err := checkDevicesAny(hardwareInfo, stack.Devices.Any)
		for _, reason := range reasonsAny {
			reasons = append(reasons, "any: "+reason)
		}
		if err != nil {
			return 0, reasons, err
		}
		if extraScore == 0 {
			return 0, reasons, nil
		}
		stackScore += extraScore
	}

	return stackScore, reasons, nil
}

func checkDevicesAll(hardwareInfo types.HwInfo, stackDevices []types.StackDevice) (int, []string, error) {
	devicesFound := 0
	extraScore := 0
	var reasons []string

	for _, device := range stackDevices {

		if device.Type == "cpu" {
			if hardwareInfo.Cpus == nil {
				reasons = append(reasons, "cpu device is required but host reported none")
				return 0, reasons, nil
			}
			cpuScore, cpuReasons, err := cpu.Match(device, hardwareInfo.Cpus)
			if err != nil {
				return 0, reasons, fmt.Errorf("cpu: %v", err)
			}
			if cpuScore == 0 {
				for _, reason := range cpuReasons {
					reasons = append(reasons, "cpu: "+reason)
				}
				reasons = append(reasons, "required cpu device not found")
				return 0, reasons, nil
			}
			extraScore += cpuScore
			devicesFound++

		} else if device.Bus == "usb" {
			// Not implemented

		} else if device.Bus == "" || device.Bus == "pci" {
			// Fallback to PCI as default bus
			if len(hardwareInfo.PciDevices) == 0 {
				reasons = append(reasons, "pci device is required but none found")
				return 0, reasons, nil
			}
			pciScore, pciReasons, err := pci.Match(device, hardwareInfo.PciDevices)
			if err != nil {
				return 0, reasons, fmt.Errorf("pci: %v", err)
			}
			if pciScore == 0 {
				for _, reason := range pciReasons {
					reasons = append(reasons, "pci: "+reason)
				}
				reasons = append(reasons, "required pci device not found")
				return 0, reasons, nil
			}
			extraScore += pciScore
			devicesFound++
		}
	}

	if len(stackDevices) > 0 && devicesFound != len(stackDevices) {
		reasons = append(reasons, "could not find a required device")
		return 0, reasons, nil
	}

	return extraScore, reasons, nil
}

func checkDevicesAny(hardwareInfo types.HwInfo, stackDevices []types.StackDevice) (int, []string, error) {
	devicesFound := 0
	extraScore := 0
	var reasons []string

	for _, device := range stackDevices {

		if device.Type == "cpu" {
			if hardwareInfo.Cpus == nil {
				continue
			}
			cpuScore, cpuReasons, err := cpu.Match(device, hardwareInfo.Cpus)
			if err != nil {
				return 0, reasons, err
			}
			if cpuScore > 0 {
				devicesFound++
				extraScore += cpuScore
			} else {
				for _, reason := range cpuReasons {
					reasons = append(reasons, "cpu: "+reason)
				}
			}

		} else if device.Bus == "usb" {
			reasons = append(reasons, "usb: not implemented")
			return 0, reasons, nil

		} else if device.Bus == "" || device.Bus == "pci" {
			// Fallback to PCI as default bus
			if hardwareInfo.PciDevices == nil {
				continue
			}
			pciScore, pciReasons, err := pci.Match(device, hardwareInfo.PciDevices)
			if err != nil {
				return 0, reasons, err
			}
			if pciScore > 0 {
				devicesFound++
				extraScore += pciScore
			} else {
				for _, reason := range pciReasons {
					reasons = append(reasons, "pci: "+reason)
				}
			}
		}
	}

	// If any-of devices are defined, we need to find at least one
	if len(stackDevices) > 0 && devicesFound == 0 {
		reasons = append(reasons, "could not find a required device")
		return 0, reasons, nil
	}

	return extraScore, reasons, nil
}
