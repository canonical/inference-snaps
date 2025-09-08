package main

import (
	"encoding/json"
	"fmt"

	"github.com/canonical/stack-utils/pkg/hardware_info"
	"github.com/canonical/stack-utils/pkg/validate"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	debugMachineInfoFormat string
)

func addDebugCommand() {
	debugCmd := &cobra.Command{
		Use:   "debug",
		Short: "Debugging commands",
		// Long:  "",
		// GroupID: "debug",
	}

	machineInfoCmd := &cobra.Command{
		Use:               "machine-info",
		Short:             "Print machine information",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              machineInfo,
	}
	machineInfoCmd.PersistentFlags().StringVar(&debugMachineInfoFormat, "format", "", "return the machine info as yaml or json")
	debugCmd.AddCommand(machineInfoCmd)

	validateCmd := &cobra.Command{
		Use:   "validate-engines",
		Short: "Validate engine manifest files",
		Args:  cobra.MinimumNArgs(1),
		RunE:  validateStackManifests,
	}
	debugCmd.AddCommand(validateCmd)

	rootCmd.AddCommand(debugCmd)
}

func machineInfo(_ *cobra.Command, args []string) error {
	hwInfo, err := hardware_info.Get(true)
	if err != nil {
		return fmt.Errorf("failed to get hardware info: %s", err)
	}

	var hwInfoStr string
	switch debugMachineInfoFormat {
	case "", "json": // Unset defaults to json
		jsonString, err := json.MarshalIndent(hwInfo, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal to JSON: %s", err)
		}
		hwInfoStr = string(jsonString)
	case "yaml":
		yamlString, err := yaml.Marshal(hwInfo)
		if err != nil {
			return fmt.Errorf("failed to marshal to YAML: %s", err)
		}
		hwInfoStr = string(yamlString)
	default:
		return fmt.Errorf("unknown format %q", debugMachineInfoFormat)
	}

	fmt.Println(hwInfoStr)

	return nil
}

func validateStackManifests(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no engine manifest specified")
	}

	allManifestsValid := true
	for _, manifestPath := range args {
		err := validate.Stack(manifestPath)
		if err != nil {
			allManifestsValid = false
			fmt.Printf("❌ %s: %s\n", manifestPath, err)
		} else {
			fmt.Printf("✅ %s\n", manifestPath)
		}
	}

	if !allManifestsValid {
		return fmt.Errorf("not all manifests are valid")
	}
	return nil
}
