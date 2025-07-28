package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/canonical/go-snapctl"
	"github.com/canonical/go-snapctl/env"
	"github.com/canonical/stack-utils/pkg/hardware_info"
	"github.com/canonical/stack-utils/pkg/selector"
	"github.com/canonical/stack-utils/pkg/types"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	useAuto      bool
	useAssumeYes bool
)

func init() {
	cmd := &cobra.Command{
		Use:   "use [<stack>]",
		Short: "Select a stack",
		// Long:  "",
		// stack use <stack> requires 1 argument
		// stack use --auto does not support any arguments
		Args: cobra.MaximumNArgs(1),
		RunE: use,
	}

	// flags
	cmd.PersistentFlags().BoolVar(&useAuto, "auto", false, "automatically select a compatible stack")
	cmd.PersistentFlags().BoolVar(&useAssumeYes, "assume-yes", false, "assume yes for downloading new components")

	rootCmd.AddCommand(cmd)
}

func use(_ *cobra.Command, args []string) error {

	if useAuto {
		if len(args) != 0 {
			return fmt.Errorf("cannot specify stack with --auto flag")
		}
		err := autoSelectStacks(useAssumeYes)
		if err != nil {
			return fmt.Errorf("failed to automatically set used stack: %s", err)
		}
	} else {
		if len(args) == 1 {
			err := useStack(args[0], useAssumeYes)
			if err != nil {
				return fmt.Errorf("failed to use stack: %s", err)
			}
		} else {
			return fmt.Errorf("stack name not specified")
		}
	}
	return nil
}

func autoSelectStacks(assumeYes bool) error {
	fmt.Println("Automatically selecting a compatible stack ...")

	allStacks, err := selector.LoadStacksFromDir(stacksDir)
	if err != nil {
		return fmt.Errorf("error loading stacks: %v", err)
	}

	// get hardware info
	hardwareInfo, err := hardware_info.Get(false)
	if err != nil {
		return fmt.Errorf("error getting hardware info: %v", err)
	}

	// score stacks
	scoredStacks, err := selector.ScoreStacks(hardwareInfo, allStacks)
	if err != nil {
		return fmt.Errorf("error scoring stacks: %v", err)
	}

	for _, stack := range scoredStacks {
		if stack.Score == 0 {
			fmt.Printf("❌ %s - not compatible: %s\n", stack.Name, strings.Join(stack.Notes, ", "))
		} else if stack.Grade != "stable" {
			fmt.Printf("⏺️ %s - score = %d, grade = %s\n", stack.Name, stack.Score, stack.Grade)
		} else {
			fmt.Printf("✅ %s - compatible, score = %d\n", stack.Name, stack.Score)
		}
	}

	// set all scored stacks as snap options
	for _, stack := range scoredStacks {
		stackJson, err := json.Marshal(stack)
		if err != nil {
			return fmt.Errorf("error serializing stacks: %v", err)
		}

		err = snapctl.Set("stacks."+stack.Name, string(stackJson)).Document().Run()
		if err != nil {
			return fmt.Errorf("error setting stacks option: %v", err)
		}
	}

	// find top stack
	topStack, err := selector.TopStack(scoredStacks)
	if err != nil {
		return fmt.Errorf("error selecting a stack: %v", err)
	}

	fmt.Printf("Selected stack for your hardware configuration: %s\n\n", topStack.Name)

	return useStack(topStack.Name, assumeYes)
}

/*
useStack changes the stack that is used by the snap
*/
func useStack(stackName string, assumeYes bool) error {
	stackJson, err := snapctl.Get("stacks." + stackName).Document().Run()
	if err != nil {
		return fmt.Errorf("error loading stack: %v", err)
	}

	stack, err := parseStackJson(stackJson)
	if err != nil {
		return fmt.Errorf("error parsing stack: %v", err)
	}

	components, err := missingComponents(stack.Components)
	if err != nil {
		return fmt.Errorf("error checking installed components: %v", err)
	}
	if len(components) > 0 {
		// ask user if they want to continue
		fmt.Println("Need to download and install the following components:")
		for _, component := range components {
			fmt.Printf("\t%s\n", component)
		}
		fmt.Println("This can take a long time to complete.")

		// Only ask for confirmation of download if it is an interactive terminal
		if !assumeYes && term.IsTerminal(int(os.Stdin.Fd())) {
			if !confirmationPrompt("Are you sure you want to continue?") {
				fmt.Println("Exiting. No changes applied.")
				return nil
			}
		}

		// Leave a blank line after printing component list and optional confirmation, before printing component installation progress
		fmt.Println()
	}

	// First change the stack, then download the components.
	// Even if a timeout occurs, the download is expected to complete in the background.
	err = setStackOptions(stack)
	if err != nil {
		return fmt.Errorf("error setting stack options: %v", err)
	}

	if len(components) > 0 {
		// This is blocking, but there is a timeout
		downloadComponents(stack.Components)
	}

	// TODO restart service

	return nil
}

func missingComponents(components []string) ([]string, error) {
	var missing []string
	for _, component := range components {
		isInstalled, err := componentInstalled(component)
		if err != nil {
			return missing, err
		}
		if !isInstalled {
			missing = append(missing, component)
		}
	}
	return missing, nil
}

func componentInstalled(component string) (bool, error) {
	// Check in /snap/$SNAP_INSTANCE_NAME/components/$SNAP_REVISION if component is mounted
	directoryPath := fmt.Sprintf("/snap/%s/components/%s/%s", env.SnapInstanceName(), env.SnapRevision(), component)

	info, err := os.Stat(directoryPath)

	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, fmt.Errorf("error checking component directory %q: %v", component, err)
		}
	} else {
		if info.IsDir() {
			return true, nil
		} else {
			return false, fmt.Errorf("component %q exists but is not a directory", component)
		}
	}
}

func setStackOptions(stack types.ScoredStack) error {
	// set stack config option
	err := snapctl.Set("stack", stack.Name).Run()
	if err != nil {
		return fmt.Errorf(`error setting snap option "stack": %v`, err)
	}

	// set other config options
	// TODO: clear beforehand
	for confKey, confVal := range stack.Configurations {
		valJson, err := json.Marshal(confVal)
		if err != nil {
			return fmt.Errorf("error serializing configuration %q: %v - %v", confKey, confVal, err)
		}
		err = snapctl.Set(confKey, string(valJson)).Document().Run()
		if err != nil {
			return fmt.Errorf("error setting snap option %q: %v", confKey, err)
		}
	}

	return nil
}

// confirmationPrompt prompts the user for a yes/no answer and returns true for 'y', false for 'n'.
func confirmationPrompt(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n] ", prompt)
		input, _ := reader.ReadString('\n')
		input = strings.ToLower(strings.TrimSpace(input))

		if input == "y" || input == "yes" {
			return true
		} else if input == "n" || input == "no" {
			return false
		} else {
			fmt.Println(`Invalid input. Please enter "y" or "n".`)
		}
	}
}
