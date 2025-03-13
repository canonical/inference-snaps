package selector

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/canonical/ml-snap-utils/pkg/types"
	"github.com/canonical/ml-snap-utils/pkg/utils"
	"gopkg.in/yaml.v3"
)

func TopStack(scoredStacks []types.ScoredStack) (*types.ScoredStack, error) {
	var compatibleStacks []types.ScoredStack

	for _, stack := range scoredStacks {
		if stack.Score == 0 {
			log.Printf("Skipping stack %s because it is not compatible", stack.Name)
		} else if stack.Grade != "stable" {
			log.Printf("Skipping stack %s because it is not stable", stack.Name)
		} else {
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
		score, err := checkStack(hardwareInfo, currentStack)

		scoredStack := types.ScoredStack{
			Stack:      currentStack,
			Score:      score,
			Compatible: true,
		}

		if score == 0 {
			scoredStack.Compatible = false
		}

		if err != nil {
			scoredStack.Notes = append(scoredStack.Notes, err.Error())
		}

		scoredStacks = append(scoredStacks, scoredStack)
	}

	return scoredStacks, nil
}

func checkStack(hardwareInfo types.HwInfo, stack types.Stack) (int, error) {
	stackScore := 0

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
		if hardwareInfo.Memory.TotalRam+hardwareInfo.Memory.TotalSwap < requiredMemory {
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
			if hardwareInfo.Cpus == nil {
				return 0, fmt.Errorf("cpu device is required but none found")
			}
			cpuScore, err := checkCpus(device, hardwareInfo.Cpus)
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
			if hardwareInfo.Cpus == nil {
				continue
			}
			cpuScore, err := checkCpus(device, hardwareInfo.Cpus)
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
