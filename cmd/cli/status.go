package main

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	statusFormat string
)

func addStatusCommand() {
	cmd := &cobra.Command{
		Use:               "status",
		Short:             "Show the status",
		Long:              "Show the status of the inference snap",
		GroupID:           "basics",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              status,
	}

	// flags
	cmd.PersistentFlags().StringVar(&statusFormat, "format", "yaml", "output format")

	rootCmd.AddCommand(cmd)
}

func status(_ *cobra.Command, _ []string) error {
	var statusText string
	var err error

	stopProgress := startProgressSpinner("Getting status ")
	defer stopProgress()

	switch statusFormat {
	case "json":
		statusText, err = statusJson()
		if err != nil {
			return fmt.Errorf("error getting json status: %v", err)
		}
		statusText += "\n"
	case "yaml":
		statusText, err = statusYaml()
		if err != nil {
			return fmt.Errorf("error getting yaml status: %v", err)
		}
	default:
		return fmt.Errorf("unknown format %q", statusFormat)
	}

	stopProgress()

	fmt.Print(statusText)

	return nil
}

func statusYaml() (string, error) {
	statusStr, err := statusStruct()
	if err != nil {
		return "", fmt.Errorf("error getting status: %v", err)
	}
	yamlStr, err := yaml.Marshal(statusStr)
	if err != nil {
		return "", fmt.Errorf("error marshalling yaml: %v", err)
	}
	return string(yamlStr), nil
}

func statusJson() (string, error) {
	statusStr, err := statusStruct()
	if err != nil {
		return "", fmt.Errorf("error getting status: %v", err)
	}
	jsonStr, err := json.MarshalIndent(statusStr, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshalling json: %v", err)
	}
	return string(jsonStr), nil
}
