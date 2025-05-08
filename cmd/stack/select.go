package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/canonical/go-snapctl"
	"github.com/canonical/ml-snap-utils/pkg/hardware_info"
	"github.com/canonical/ml-snap-utils/pkg/selector"
	"github.com/canonical/ml-snap-utils/pkg/types"
	"gopkg.in/yaml.v3"
)

func autoSelectStacks() {
	fmt.Println("Automatically selecting a compatible stack ...")

	allStacks, err := selector.LoadStacksFromDir(stacksDir)
	if err != nil {
		fmt.Println("Error loading stacks:", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d stacks\n", len(allStacks))

	// get hardware info
	hardwareInfo, err := hardware_info.Get(false)
	if err != nil {
		fmt.Println("Error getting hardware info:", err)
		os.Exit(1)
	}

	// score stacks
	scoredStacks, err := selector.ScoreStacks(hardwareInfo, allStacks)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, stack := range scoredStacks {
		if stack.Score == 0 {
			fmt.Printf("Stack %s not selected: %s\n", stack.Name, strings.Join(stack.Notes, ", "))
		} else {
			fmt.Printf("Stack %s matches. Score = %d\n", stack.Name, stack.Score)
		}
	}

	// set all scored stacks as snap options
	for _, stack := range scoredStacks {
		stackJson, err := json.Marshal(stack)
		if err != nil {
			fmt.Println("Error serializing stacks:", err)
			os.Exit(1)
		}

		err = snapctl.Set("stacks."+stack.Name, string(stackJson)).Document().Run()
		if err != nil {
			fmt.Println("Error setting stacks option:", err)
			os.Exit(1)
		}
	}

	// find top stack
	topStack, err := selector.TopStack(scoredStacks)
	if err != nil {
		fmt.Println("Error selecting a stack:", err)
		os.Exit(1)
	}

	// set top stack name as snap option
	err = snapctl.Set("stack", topStack.Name).String().Run()
	if err != nil {
		fmt.Println("Error setting stack:", err)
		os.Exit(1)
	}

	// set snap options from stack configurations
	for confKey, confVal := range topStack.Configurations {
		valJson, err := json.Marshal(confVal)
		if err != nil {
			fmt.Printf("Error serializing configuration %s: %v - %v\n", confKey, confVal, err)
			os.Exit(1)
		}
		err = snapctl.Set(confKey, string(valJson)).Document().Run()
		if err != nil {
			fmt.Println("Error setting snap option:", err)
			os.Exit(1)
		}
	}

	fmt.Println("Selected stack for your hardware configuration:", topStack.Name)
}

func selectStack(stackName string) {

	stack, err := loadStacksFromDir(stacksDir, stackName)
	if err != nil {
		fmt.Println("Error loading stack:", err)
		os.Exit(1)
	}

	for confKey, confVal := range stack.Configurations {
		valJson, err := json.Marshal(confVal)
		if err != nil {
			fmt.Printf("Error serializing configuration %s: %v - %v\n", confKey, confVal, err)
			os.Exit(1)
		}
		err = snapctl.Set(confKey, string(valJson)).Document().Run()
		if err != nil {
			fmt.Println("Error setting snap option:", err)
			os.Exit(1)
		}
	}
}

func loadStacksFromDir(stacksDir, stackName string) (*types.Stack, error) {

	// Sanitise stack dir path
	if !strings.HasSuffix(stacksDir, "/") {
		stacksDir += "/"
	}

	stackFilename := stacksDir + stackName + "/stack.yaml"
	data, err := os.ReadFile(stackFilename)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %s", stackFilename, err)
	}

	var stack types.Stack
	err = yaml.Unmarshal(data, &stack)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", stacksDir, err)
	}

	return &stack, nil
}
