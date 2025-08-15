package main

import (
	"encoding/json"
	"fmt"

	"github.com/canonical/go-snapctl"
	"github.com/canonical/stack-utils/pkg/selector"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "load",
		Short: "Initialize snap configurations",
		// Long:  "",
		Hidden:            true, // command for internal use
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              load,
	}
	rootCmd.AddCommand(cmd)
}

func load(_ *cobra.Command, _ []string) error {
	return loadStacksToSnapOptions()
}

func loadStacksToSnapOptions() error {
	fmt.Println("Loading stacks to snap options ...")

	allStacks, err := selector.LoadStacksFromDir(stacksDir)
	if err != nil {
		return fmt.Errorf("error loading stacks: %v", err)
	}

	// set all stacks as snap options
	// TODO: change to also handle stack deletions
	for _, stack := range allStacks {
		stackJson, err := json.Marshal(stack)
		if err != nil {
			return fmt.Errorf("error serializing stacks: %s", err)
		}

		err = snapctl.Set("stacks."+stack.Name, string(stackJson)).Document().Run()
		if err != nil {
			return fmt.Errorf("error setting stacks option: %s", err)
		}
	}

	return nil
}
