package main

import (
	"encoding/json"
	"fmt"

	"github.com/canonical/inference-snaps-cli/pkg/hardware_info"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func addShowMachineCommand() {
	cmd := &cobra.Command{
		Use:               "show-machine",
		Short:             "Print information about the host machine",
		Long:              "Print information about the host machine, including hardware and compute resources",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              showMachine,
	}
	cmd.PersistentFlags().StringVar(&debugMachineInfoFormat, "format", "yaml", "output format")
	rootCmd.AddCommand(cmd)
}

func showMachine(_ *cobra.Command, args []string) error {
	hwInfo, err := hardware_info.Get(true)
	if err != nil {
		return fmt.Errorf("failed to get machine info: %s", err)
	}

	switch debugMachineInfoFormat {
	case "json":
		jsonString, err := json.MarshalIndent(hwInfo, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal to JSON: %s", err)
		}
		fmt.Printf("%s\n", jsonString)
	case "yaml":
		yamlString, err := yaml.Marshal(hwInfo)
		if err != nil {
			return fmt.Errorf("failed to marshal to YAML: %s", err)
		}
		fmt.Printf("%s", yamlString)
	default:
		return fmt.Errorf("unknown format %q", debugMachineInfoFormat)
	}

	return nil
}
