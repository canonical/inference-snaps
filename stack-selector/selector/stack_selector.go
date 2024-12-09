package selector

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
	"katemoss/common"
)

func FindStack(hardwareInfo common.HwInfo, stacksDir string) (*common.StackResult, error) {
	var foundStacks []common.StackResult

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

		var stack common.Stack
		err = yaml.Unmarshal(data, &stack)
		if err != nil {
			return nil, fmt.Errorf("%s: %s", stacksDir, err)
		}

		score, err := checkStack(hardwareInfo, stack)
		if err != nil {
			log.Printf("Stack %s not selected: %s", stack.Name, err)
			continue
		}

		if score > 0 {
			foundStack := common.StackResult{
				Name:       stack.Name,
				Components: stack.Components,
				Score:      score,
			}
			foundStacks = append(foundStacks, foundStack)
			log.Printf("Stack %s matches. Score = %d", stack.Name, score)
		}
	}

	// If none found, return err
	if len(foundStacks) == 0 {
		return nil, fmt.Errorf("no stack found matching hardware")
	}

	// Sort by score (high to low) and return best match
	sort.Slice(foundStacks, func(i, j int) bool {
		return foundStacks[i].Score > foundStacks[j].Score
	})

	// TODO find duplicate scores, use a different metric to choose one of them

	return &foundStacks[0], nil
}

func checkStack(hardwareInfo common.HwInfo, stack common.Stack) (int, error) {
	score := 0

	// Enough memory
	if stack.Memory != nil {
		requiredMemory, err := StringToBytes(*stack.Memory)
		if err != nil {
			return 0, err
		}
		// Checking combination of ram and swap
		if hardwareInfo.Memory.RamTotal+hardwareInfo.Memory.SwapTotal < requiredMemory {
			return 0, fmt.Errorf("not enough memory")
		}
		score++
	}

	// Enough disk space
	if stack.DiskSpace != nil {
		requiredDisk, err := StringToBytes(*stack.DiskSpace)
		if err != nil {
			return 0, err
		}
		if _, ok := hardwareInfo.Disk["/var/lib/snapd/snaps"]; !ok {
			return 0, fmt.Errorf("disk space not provided by hardware info")
		}
		if hardwareInfo.Disk["/var/lib/snapd/snaps"].Avail < requiredDisk {
			return 0, fmt.Errorf("not enough free disk space")
		}
		score++
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
			if !checkCpus(device, *hardwareInfo.Cpu) {
				return 0, fmt.Errorf("required cpu device not found")
			}
			allOfDevicesFound++

		case "gpu":
			if len(hardwareInfo.Gpus) == 0 {
				return 0, fmt.Errorf("gpu device is required but none found")
			}
			result, err := checkGpus(hardwareInfo.Gpus, device)
			if err != nil {
				return 0, err
			}
			if !result {
				return 0, fmt.Errorf("required gpu device not found")
			}
			allOfDevicesFound++
		}
	}

	if len(stack.Devices.All) > 0 && allOfDevicesFound != len(stack.Devices.All) {
		return 0, fmt.Errorf("all: could not find a required device")
	}
	score += allOfDevicesFound

	// any
	anyOfDevicesFound := 0
	for _, device := range stack.Devices.Any {
		switch device.Type {
		case "cpu":
			if checkCpus(device, *hardwareInfo.Cpu) {
				anyOfDevicesFound++
			}

		case "gpu":
			result, err := checkGpus(hardwareInfo.Gpus, device)
			if err != nil {
				return 0, err
			}
			if result {
				anyOfDevicesFound++
			}
		}
	}

	// If any-of devices are defined, we need to find at least one
	if len(stack.Devices.Any) > 0 && anyOfDevicesFound == 0 {
		return 0, fmt.Errorf("any: could not find a required device")
	}
	score += anyOfDevicesFound

	return score, nil
}
