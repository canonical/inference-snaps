package selector

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/canonical/ml-snap-utils/pkg/types"
	"github.com/canonical/ml-snap-utils/pkg/utils"
	"gopkg.in/yaml.v3"
)

func FindStack(hardwareInfo types.HwInfo, stacksDir string) (*types.StackResult, error) {
	var foundStacks []types.StackResult

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

		score, err := checkStack(hardwareInfo, currentStack)
		if err != nil {
			log.Printf("Stack %s not selected: %s", currentStack.Name, err)
			continue
		}

		if score > 0 {
			foundStack := types.StackResult{
				Name:           currentStack.Name,
				Components:     currentStack.Components,
				Configurations: currentStack.Configurations,
				Score:          score,
			}
			foundStacks = append(foundStacks, foundStack)
			log.Printf("Stack %s matches. Score = %f", currentStack.Name, score)
		}
	}

	// If none found, return err
	if len(foundStacks) == 0 {
		return nil, fmt.Errorf("no stack found matching this hardware")
	}

	// Sort by score (high to low) and return best match
	sort.Slice(foundStacks, func(i, j int) bool {
		return foundStacks[i].Score > foundStacks[j].Score
	})

	// TODO find duplicate scores, use a different metric to choose one of them

	return &foundStacks[0], nil
}

func checkStack(hardwareInfo types.HwInfo, stack types.Stack) (float64, error) {
	stackScore := 0.0

	// Enough memory
	if stack.Memory != nil {
		requiredMemory, err := utils.StringToBytes(*stack.Memory)
		if err != nil {
			return 0, err
		}

		if hardwareInfo.Memory == nil {
			return 0, fmt.Errorf("no memory in hardware info")
		}

		// Checking combination of ram and swap
		if hardwareInfo.Memory.RamTotal+hardwareInfo.Memory.SwapTotal < requiredMemory {
			return 0, fmt.Errorf("not enough memory")
		}
		stackScore++
	}

	// Enough disk space
	if stack.DiskSpace != nil {
		requiredDisk, err := utils.StringToBytes(*stack.DiskSpace)
		if err != nil {
			return 0, err
		}
		if _, ok := hardwareInfo.Disk["/var/lib/snapd/snaps"]; !ok {
			return 0, fmt.Errorf("disk space not provided by hardware info")
		}
		if hardwareInfo.Disk["/var/lib/snapd/snaps"].Avail < requiredDisk {
			return 0, fmt.Errorf("not enough free disk space")
		}
		stackScore++
	}

	// Devices
	// all
	allOfDevicesFound := 0
	for _, device := range stack.Devices.All {
		switch device.Type {
		case "cpu":
			if hardwareInfo.Cpu == nil {
				return 0, fmt.Errorf("cpu device is required but none found")
			}
			cpuScore, err := checkCpus(device, *hardwareInfo.Cpu)
			if err != nil {
				return 0, err
			}
			if cpuScore == 0 {
				return 0, fmt.Errorf("required cpu device not found")
			}
			stackScore += cpuScore
			allOfDevicesFound++

		case "gpu":
			if len(hardwareInfo.Gpus) == 0 {
				return 0, fmt.Errorf("gpu device is required but none found")
			}
			gpuScore, err := checkGpus(hardwareInfo.Gpus, device)
			if err != nil {
				return 0, err
			}
			if gpuScore == 0 {
				return 0, fmt.Errorf("required gpu device not found")
			}
			stackScore += gpuScore
			allOfDevicesFound++
		}
	}

	if len(stack.Devices.All) > 0 && allOfDevicesFound != len(stack.Devices.All) {
		return 0, fmt.Errorf("all: could not find a required device")
	}

	// any
	anyOfDevicesFound := 0
	for _, device := range stack.Devices.Any {
		switch device.Type {
		case "cpu":
			if hardwareInfo.Cpu == nil {
				continue
			}
			cpuScore, err := checkCpus(device, *hardwareInfo.Cpu)
			if err != nil {
				return 0, err
			}
			if cpuScore > 0 {
				anyOfDevicesFound++
			}
			stackScore += cpuScore

		case "gpu":
			if hardwareInfo.Gpus == nil {
				continue
			}
			gpuScore, err := checkGpus(hardwareInfo.Gpus, device)
			if err != nil {
				return 0, err
			}
			if gpuScore > 0 {
				anyOfDevicesFound++
			}
			stackScore += gpuScore
		}
	}

	// If any-of devices are defined, we need to find at least one
	if len(stack.Devices.Any) > 0 && anyOfDevicesFound == 0 {
		return 0, fmt.Errorf("any: could not find a required device")
	}

	return stackScore, nil
}
