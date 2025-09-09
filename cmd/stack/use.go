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
	"github.com/canonical/stack-utils/pkg/snap_store"
	"github.com/canonical/stack-utils/pkg/types"
	"github.com/canonical/stack-utils/pkg/utils"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	useAuto      bool
	useAssumeYes bool
)

func addUseCommand() {
	cmd := &cobra.Command{
		Use:   "use-engine [<engine>]",
		Short: "Select an engine",
		// Long:  "",
		GroupID: "engines",
		// stack use <stack> requires 1 argument
		// stack use --auto does not support any arguments
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: useValidArgs,
		RunE:              use,
	}

	// flags
	cmd.PersistentFlags().BoolVar(&useAuto, "auto", false, "automatically select a compatible engine")
	cmd.PersistentFlags().BoolVar(&useAssumeYes, "assume-yes", false, "assume yes for downloading new components")

	rootCmd.AddCommand(cmd)
}

func useValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
	stacksJson, err := snapctl.Get("engines").Document().Run()
	if err != nil {
		fmt.Printf("Error loading engines: %v", err)
		return nil, cobra.ShellCompDirectiveError
	}

	stacks, err := parseStacksJson(stacksJson)
	if err != nil {
		fmt.Printf("Error parsing engines: %v", err)
		return nil, cobra.ShellCompDirectiveError
	}

	var stackNames []cobra.Completion
	for i := range stacks {
		if stacks[i].Compatible {
			stackNames = append(stackNames, stacks[i].Name)
		}
	}
	if len(stackNames) == 0 {
		// No stacks flagged as compatible
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return stackNames, cobra.ShellCompDirectiveNoFileComp
}

func use(_ *cobra.Command, args []string) error {

	if useAuto {
		if len(args) != 0 {
			return fmt.Errorf("cannot specify both engine name and --auto flag")
		}

		scoredStacks, err := scoreStacks()
		if err != nil {
			return fmt.Errorf("error scoring engines: %v", err)
		}

		for _, stack := range scoredStacks {
			if stack.Score == 0 {
				fmt.Printf("❌ %s - not compatible: %s\n", stack.Name, strings.Join(stack.Notes, ", "))
			} else if stack.Grade != "stable" {
				fmt.Printf("🟠 %s - score = %d, grade = %s\n", stack.Name, stack.Score, stack.Grade)
			} else {
				fmt.Printf("✅ %s - compatible, score = %d\n", stack.Name, stack.Score)
			}
		}

		err = stacksToSnapOptions(scoredStacks)
		if err != nil {
			return fmt.Errorf("error saving scored engines: %v", err)
		}

		fmt.Println("Automatically selecting a compatible engine ...")

		selectedStack, err := selector.TopStack(scoredStacks)
		if err != nil {
			return fmt.Errorf("error finding top engine: %v", err)
		}

		fmt.Printf("Selected engine for your hardware configuration: %s\n\n", selectedStack.Name)

		err = useStack(selectedStack.Name, useAssumeYes)
		if err != nil {
			return fmt.Errorf("failed to use engine: %s", err)
		}

	} else {
		if len(args) == 1 {
			err := useStack(args[0], useAssumeYes)
			if err != nil {
				return fmt.Errorf("failed to use engine: %s", err)
			}
		} else {
			return fmt.Errorf("engine name not specified")
		}
	}
	return nil
}

func scoreStacks() ([]types.ScoredStack, error) {
	allStacks, err := selector.LoadManifestsFromDir(enginesDir)
	if err != nil {
		return nil, fmt.Errorf("error loading engines: %v", err)
	}

	// get hardware info
	hardwareInfo, err := hardware_info.Get(false)
	if err != nil {
		return nil, fmt.Errorf("error getting hardware info: %v", err)
	}

	// score stacks
	scoredStacks, err := selector.ScoreStacks(hardwareInfo, allStacks)
	if err != nil {
		return nil, fmt.Errorf("error scoring engines: %v", err)
	}

	return scoredStacks, nil
}

func stacksToSnapOptions(scoredStacks []types.ScoredStack) error {
	// set all scored stacks as snap options
	for _, stack := range scoredStacks {
		stackJson, err := json.Marshal(stack)
		if err != nil {
			return fmt.Errorf("error serializing engines: %v", err)
		}

		err = snapctl.Set("engines."+stack.Name, string(stackJson)).Document().Run()
		if err != nil {
			return fmt.Errorf("error setting engine option: %v", err)
		}
	}
	return nil
}

/*
useStack changes the stack that is used by the snap
*/
func useStack(stackName string, assumeYes bool) error {
	stackJson, err := snapctl.Get("engines." + stackName).Document().Run()
	if err != nil {
		return fmt.Errorf("error loading engine: %v", err)
	}

	stack, err := parseStackJson(stackJson)
	if err != nil {
		return fmt.Errorf("error parsing engine: %v", err)
	}

	components, err := missingComponents(stack.Components)
	if err != nil {
		return fmt.Errorf("error checking installed components: %v", err)
	}
	if len(components) > 0 {
		// Look up component sizes from the snap store
		componentSizes, err := snap_store.ComponentSizes()
		if err != nil {
			// If component size lookup failed, continue but log the error
			fmt.Fprintf(os.Stderr, "Error getting component sizes: %v\n", err)
		}

		// Format list of components, adding size if it is known
		var componentList []string
		for _, componentName := range components {
			line := fmt.Sprintf("- %s", componentName)
			if size, ok := componentSizes[componentName]; ok {
				line += fmt.Sprintf(" (%s)", utils.FmtBytes(uint64(size)))
			}
			componentList = append(componentList, line)
		}

		fmt.Println("Need to install the following components:")
		for _, component := range componentList {
			fmt.Println(component)
		}

		// Only ask for confirmation of download if it is an interactive terminal
		if !assumeYes && term.IsTerminal(int(os.Stdin.Fd())) {
			fmt.Println()
			if !confirmationPrompt("Do you want to continue?") {
				fmt.Println("Exiting. No changes applied.")
				return nil
			}
		}

		// Leave a blank line after printing component list and optional confirmation, before printing component installation progress
		fmt.Println()

		// This is blocking, but there is a timeout bug:
		// https://github.com/canonical/stack-utils/issues/122
		err = installComponents(stack.Components)
		if err != nil {
			return fmt.Errorf("error installing components: %v", err)
		}
	}

	currentEngine, err := snapctl.Get("engine").Run()
	if err != nil {
		return fmt.Errorf("error getting current engine: %v", err)
	}
	if currentEngine == stackName {
		// Nothing left to do
		return nil
	}

	if len(components) > 0 {
		// Leave a blank line if components were installed, before continuing
		fmt.Println()
	}

	err = setStackOptions(stack)
	if err != nil {
		return fmt.Errorf("error setting engine options: %v", err)
	}

	fmt.Println("Restarting the snap service ...")
	err = snapctl.Restart(snapInstanceName).Run()
	if err != nil {
		return fmt.Errorf("error restarting snap service: %v", err)
	}

	fmt.Printf("Engine successfully changed to %q\n", stackName)

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
	err := snapctl.Set("engine", stack.Name).Run()
	if err != nil {
		return fmt.Errorf(`error setting snap option "engine": %v`, err)
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
