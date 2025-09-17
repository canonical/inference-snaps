package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/canonical/go-snapctl"
	"github.com/canonical/go-snapctl/env"
	"github.com/canonical/stack-utils/pkg/engines"
	"github.com/canonical/stack-utils/pkg/hardware_info"
	"github.com/canonical/stack-utils/pkg/selector"
	"github.com/canonical/stack-utils/pkg/snap_store"
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
		// cli use-engine <engine> requires 1 argument
		// cli use-engine --auto does not support any arguments
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
	enginesJson, err := snapctl.Get("engines").Document().Run()
	if err != nil {
		fmt.Printf("Error loading engines: %v", err)
		return nil, cobra.ShellCompDirectiveError
	}

	engines, err := parseEnginesJson(enginesJson)
	if err != nil {
		fmt.Printf("Error parsing engines: %v", err)
		return nil, cobra.ShellCompDirectiveError
	}

	var engineNames []cobra.Completion
	for i := range engines {
		if engines[i].Compatible {
			engineNames = append(engineNames, engines[i].Name)
		}
	}
	if len(engineNames) == 0 {
		// No engines flagged as compatible
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return engineNames, cobra.ShellCompDirectiveNoFileComp
}

func use(_ *cobra.Command, args []string) error {
	if !utils.IsRootUser() {
		return ErrPermissionDenied
	}

	if useAuto {
		if len(args) != 0 {
			return fmt.Errorf("cannot specify both engine name and --auto flag")
		}

		scoredEngines, err := scoreEngines()
		if err != nil {
			return fmt.Errorf("error scoring engines: %v", err)
		}

		for _, engine := range scoredEngines {
			if engine.Score == 0 {
				fmt.Printf("❌ %s - not compatible: %s\n", engine.Name, strings.Join(engine.Notes, ", "))
			} else if engine.Grade != "stable" {
				fmt.Printf("🟠 %s - score = %d, grade = %s\n", engine.Name, engine.Score, engine.Grade)
			} else {
				fmt.Printf("✅ %s - compatible, score = %d\n", engine.Name, engine.Score)
			}
		}

		err = enginesToSnapOptions(scoredEngines)
		if err != nil {
			return fmt.Errorf("error saving scored engines: %v", err)
		}

		fmt.Println("Automatically selecting a compatible engine ...")

		selectedEngine, err := selector.TopEngine(scoredEngines)
		if err != nil {
			return fmt.Errorf("error finding top engine: %v", err)
		}

		fmt.Printf("Selected engine for your hardware configuration: %s\n\n", selectedEngine.Name)

		err = useEngine(selectedEngine.Name, useAssumeYes)
		if err != nil {
			return fmt.Errorf("failed to use engine: %s", err)
		}

	} else {
		if len(args) == 1 {
			err := useEngine(args[0], useAssumeYes)
			if err != nil {
				return fmt.Errorf("failed to use engine: %s", err)
			}
		} else {
			return fmt.Errorf("engine name not specified")
		}
	}
	return nil
}

func scoreEngines() ([]engines.ScoredManifest, error) {
	allEngines, err := selector.LoadManifestsFromDir(enginesDir)
	if err != nil {
		return nil, fmt.Errorf("error loading engines: %v", err)
	}

	// get hardware info
	hardwareInfo, err := hardware_info.Get(false)
	if err != nil {
		return nil, fmt.Errorf("error getting hardware info: %v", err)
	}

	// score engines
	scoredEngines, err := selector.ScoreEngines(hardwareInfo, allEngines)
	if err != nil {
		return nil, fmt.Errorf("error scoring engines: %v", err)
	}

	return scoredEngines, nil
}

func enginesToSnapOptions(scoredEngines []engines.ScoredManifest) error {
	// set all scored engines as snap options
	for _, engine := range scoredEngines {
		engineJson, err := json.Marshal(engine)
		if err != nil {
			return fmt.Errorf("error serializing engines: %v", err)
		}

		err = snapctl.Set("engines."+engine.Name, string(engineJson)).Document().Run()
		if err != nil {
			return fmt.Errorf("error setting engine option: %v", err)
		}
	}
	return nil
}

/*
useEngine changes the engine that is used by the snap
*/
func useEngine(engineName string, assumeYes bool) error {
	engineJson, err := snapctl.Get("engines." + engineName).Document().Run()
	if err != nil {
		return fmt.Errorf("error loading engine: %v", err)
	}

	engine, err := parseEngineJson(engineJson)
	if err != nil {
		return fmt.Errorf("error parsing engine: %v", err)
	}

	components, err := missingComponents(engine.Components)
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
		err = installComponents(engine.Components)
		if err != nil {
			return fmt.Errorf("error installing components: %v", err)
		}
	}

	currentEngine, err := snapctl.Get("engine").Run()
	if err != nil {
		return fmt.Errorf("error getting current engine: %v", err)
	}
	if currentEngine == engineName {
		// Nothing left to do
		return nil
	}

	if len(components) > 0 {
		// Leave a blank line if components were installed, before continuing
		fmt.Println()
	}

	err = setEngineOptions(engine)
	if err != nil {
		return fmt.Errorf("error setting engine options: %v", err)
	}

	fmt.Println("Restarting the snap service ...")
	err = snapctl.Restart(snapInstanceName).Run()
	if err != nil {
		return fmt.Errorf("error restarting snap service: %v", err)
	}

	fmt.Printf("Engine successfully changed to %q\n", engineName)

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

func setEngineOptions(engine engines.ScoredManifest) error {
	// set engine config option
	err := snapctl.Set("engine", engine.Name).Run()
	if err != nil {
		return fmt.Errorf(`error setting snap option "engine": %v`, err)
	}

	// set other config options
	// TODO: clear beforehand
	for confKey, confVal := range engine.Configurations {
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
