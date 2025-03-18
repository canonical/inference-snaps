package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/canonical/go-snapctl"
	"github.com/canonical/ml-snap-utils/pkg/selector"
)

func loadStacksToSnapOptions() {
	fmt.Println("Loading stacks to snap options ...")

	allStacks, err := selector.LoadStacksFromDir(stacksDir)
	if err != nil {
		fmt.Printf("Error loading stacks: %v", err)
		os.Exit(1)
	}

	// set all stacks as snap options
	// TODO: change to also handle stack deletions
	for _, stack := range allStacks {
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
}
