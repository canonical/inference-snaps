package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/canonical/ml-snap-utils/pkg/selector"
	"github.com/canonical/ml-snap-utils/pkg/types"
)

func main() {
	var stacksDir string
	var listAll bool
	var prettyOutput bool

	flag.StringVar(&stacksDir, "stacks", "stacks", "Override the path to the stacks directory.")
	flag.BoolVar(&listAll, "all", false, "List all available stacks.")
	flag.BoolVar(&prettyOutput, "pretty", false, "Pretty print JSON.")

	flag.Parse()

	// Read json piped in from the hardware-info app
	var hardwareInfo types.HwInfo

	err := json.NewDecoder(os.Stdin).Decode(&hardwareInfo)
	if err != nil {
		log.Fatal(err)
	}

	allStacks, err := selector.LoadStacksFromDir(stacksDir)
	if err != nil {
		log.Fatal(err)
	}
	scoredStacks, err := selector.ScoreStacks(hardwareInfo, allStacks)
	if err != nil {
		log.Fatal(err)
	}

	var stackSelection types.StackSelection

	// Print summary on STDERR
	for _, stack := range scoredStacks {
		stackSelection.Stacks = append(stackSelection.Stacks, stack)
		if stack.Score == 0 {
			log.Printf("Stack %s not selected: %s", stack.Name, strings.Join(stack.Notes, ", "))
		} else {
			log.Printf("Stack %s matches. Score = %d", stack.Name, stack.Score)
		}
	}

	topStack, err := selector.TopStack(scoredStacks)
	if err != nil {
		log.Fatal(err)
	}
	stackSelection.TopStack = topStack.Name
	log.Printf("Top stack: %s. Score = %d", topStack.Name, topStack.Score)

	var resultStr []byte
	if prettyOutput {
		resultStr, err = json.MarshalIndent(stackSelection, "", "  ")
	} else {
		resultStr, err = json.Marshal(stackSelection)
	}
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", resultStr)
}
