package main

import (
	"encoding/json"
	"fmt"

	"github.com/canonical/stack-utils/pkg/hardware_info"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:               "system",
		Short:             "Query system info",
		Long:              "Query information about the hardware and available resources",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              system,
	}

	rootCmd.AddCommand(cmd)
}

func system(_ *cobra.Command, args []string) error {
	hwInfo, err := hardware_info.Get(true)
	if err != nil {
		return fmt.Errorf("failed to get hardware info: %s", err)
	}

	jsonString, err := json.MarshalIndent(hwInfo, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal to JSON: %s", err)
	}

	// print the JSON output
	fmt.Println(string(jsonString))

	return nil
}
