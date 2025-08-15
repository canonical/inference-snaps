package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/canonical/go-snapctl"
	"github.com/canonical/stack-utils/pkg/selector"
	"github.com/canonical/stack-utils/pkg/types"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:               "status",
		Short:             "Show the status",
		Long:              "Show the status of the model snap",
		GroupID:           "basics",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              status,
	}

	rootCmd.AddCommand(cmd)
}

func status(_ *cobra.Command, _ []string) error {
	return snapStatus()
}

func snapStatus() error {
	// Find the top stack
	compatibleStacks := true
	scoredStacks, err := scoredStacksFromOptions()
	if err != nil {
		return fmt.Errorf("error loading scored stacks: %v", err)
	}

	autoStack, err := selector.TopStack(scoredStacks)
	if err != nil {
		if errors.Is(err, selector.ErrorNoCompatibleStacks) {
			compatibleStacks = false
		} else {
			return fmt.Errorf("error loading top stack: %v", err)
		}
	}

	// Find the selected stack
	stack, err := selectedStackFromOptions()
	if err != nil {
		return fmt.Errorf("error loading selected stack: %v", err)
	}

	printStack(stack, compatibleStacks && stack.Name == autoStack.Name)
	fmt.Println("")
	err = printServer(stack)
	if err != nil {
		return fmt.Errorf("error showing server status: %v", err)
	}
	fmt.Println("")

	return nil
}

func scoredStacksFromOptions() ([]types.ScoredStack, error) {
	stacksJson, err := snapctl.Get("stacks").Document().Run()
	if err != nil {
		return nil, fmt.Errorf("error loading stacks: %v", err)
	}

	stacksMap, err := parseStacksJson(stacksJson)
	if err != nil {
		return nil, fmt.Errorf("error parsing stacks: %v", err)
	}

	// map to slice
	var stacks []types.ScoredStack
	for _, stack := range stacksMap {
		stacks = append(stacks, stack)
	}

	return stacks, nil
}

func selectedStackFromOptions() (types.ScoredStack, error) {
	selectedStackName, err := snapctl.Get("stack").Run()
	if err != nil {
		return types.ScoredStack{}, fmt.Errorf("error loading selected stack: %v", err)
	}

	stackJson, err := snapctl.Get("stacks." + selectedStackName).Document().Run()
	if err != nil {
		return types.ScoredStack{}, fmt.Errorf("error loading stack: %v", err)
	}

	stack, err := parseStackJson(stackJson)
	if err != nil {
		return types.ScoredStack{}, fmt.Errorf("error parsing stack: %v", err)
	}

	return stack, nil
}

func printStack(stack types.ScoredStack, auto bool) {
	autoString := ""
	if auto {
		autoString = " (auto)"
	}
	fmt.Printf("Stack: %s%s\n", stack.Name, autoString)

	if val, ok := stack.Configurations["model"]; ok {
		fmt.Printf("  Model: %s\n", val)
	}
	if val, ok := stack.Configurations["engine"]; ok {
		fmt.Printf("  Engine: %s\n", val)
	}
	if val, ok := stack.Configurations["multimodel-projector"]; ok {
		fmt.Printf("  Multimodal projector: %s\n", val)
	}
}

func printServer(stack types.ScoredStack) error {
	apiBasePath := "v1"
	if val, ok := stack.Configurations["http.base-path"]; ok {
		apiBasePath, ok = val.(string)
		if !ok {
			return fmt.Errorf("unexpected type for base path: %v", val)
		}

	}
	httpPort, err := snapctl.Get("http.port").Run()
	if err != nil {
		return fmt.Errorf("error getting http port: %v", err)
	}

	// Depend on existing check server scripts for status
	checkScript := os.ExpandEnv("$SNAP/stacks/" + stack.Name + "/check-server")
	cmd := exec.Command(checkScript)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error checking server: %v", err)
	}

	checkExitCode := 0
	if err := cmd.Wait(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			checkExitCode = exitError.ExitCode()
		}
	}

	statusText := "online"
	switch checkExitCode {
	case 0:
		statusText = "online"
	case 1:
		statusText = "starting"
	case 2:
		statusText = "offline"
	default:
		statusText = fmt.Sprintf("unknown (exit code %d)", checkExitCode)
	}

	fmt.Printf("Server:\n")
	fmt.Printf("  Status: %s\n", statusText)
	fmt.Printf("  OpenAI endpoint: http://localhost:%s/%s\n", httpPort, apiBasePath)

	return nil
}
