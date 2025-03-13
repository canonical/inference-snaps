package main

import (
	"encoding/json"
	"strings"

	"github.com/canonical/go-snapctl"
	"github.com/canonical/go-snapctl/env"
	slog "github.com/canonical/go-snapctl/log"
	"github.com/canonical/ml-snap-utils/pkg/hardware_info"
	"github.com/canonical/ml-snap-utils/pkg/selector"
)

func main() {
	slog.SetComponentName("hook.install")

	// get all stacks from "$SNAP"/stacks/*/stack.yaml and append to a single json array
	stacksDir := env.Snap + "/stacks"
	allStacks, err := selector.LoadStacksFromDir(stacksDir)
	if err != nil {
		slog.Fatalf("Error loading stacks: %v", err)
	}

	slog.Infof("Found %d stacks", len(allStacks))

	// get hardware info
	hardwareInfo, err := hardware_info.Get(false)
	if err != nil {
		slog.Fatalf("Error getting hardware info: %v", err)
	}

	// score stacks
	scoredStacks, err := selector.ScoreStacks(hardwareInfo, allStacks)
	if err != nil {
		slog.Fatal(err)
	}

	for _, stack := range scoredStacks {
		if stack.Score == 0 {
			slog.Infof("Stack %s not selected: %s", stack.Name, strings.Join(stack.Notes, ", "))
		} else {
			slog.Infof("Stack %s matches. Score = %d", stack.Name, stack.Score)
		}
	}

	// set all stacks as snap options under `stacks`
	for _, stack := range scoredStacks {
		stackJson, err := json.Marshal(stack)
		if err != nil {
			slog.Fatalf("Error serializing stacks: %v", err)
		}

		err = snapctl.Set("stacks."+stack.Name, string(stackJson)).Document().Run()
		if err != nil {
			slog.Fatalf("Error setting stacks option: %v", err)
		}
	}

	// find top stack
	topStack, err := selector.TopStack(scoredStacks)
	if err != nil {
		// FIXME: If no matching stacks are found, installation of this snap will fail
		slog.Fatal(err)
	}

	// set top stack name in `stack` snap option
	err = snapctl.Set("stack", topStack.Name).String().Run()
	if err != nil {
		slog.Fatalf("Error setting stack: %v", err)
	}

	slog.Infof("Selected stack for your hardware configuration: %s", topStack.Name)

	// Installation of topStack's components and snap options are done in configure hook, which runs after install hook

	// set generic configurations
	err = snapctl.Set("http.port", "8080").Run()
	if err != nil {
		slog.Fatalf("can't set snap http.port: %v", err)
	}

	err = snapctl.Set("http.host", "127.0.0.1").Run()
	if err != nil {
		slog.Fatalf("can't set snap http.host: %v", err)
	}
}
