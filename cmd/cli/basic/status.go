package basic

import (
	"encoding/json"
	"fmt"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type statusCommand struct {
	*common.Context

	// flags
	format string
}

func StatusCommand(ctx *common.Context) *cobra.Command {
	var cmd statusCommand
	cmd.Context = ctx

	cobra := &cobra.Command{
		Use:               "status",
		Short:             "Show the status",
		Long:              "Show the status of the inference snap",
		GroupID:           groupID,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.run,
	}

	// flags
	cobra.Flags().StringVar(&cmd.format, "format", "yaml", "output format")

	return cobra
}

func (cmd *statusCommand) run(_ *cobra.Command, _ []string) error {
	var statusText string
	var err error

	stopProgress := common.StartProgressSpinner("Getting status")
	defer stopProgress()

	switch cmd.format {
	case "json":
		statusText, err = cmd.statusJson()
		if err != nil {
			return fmt.Errorf("error getting json status: %v", err)
		}
		statusText += "\n"
	case "yaml":
		statusText, err = cmd.statusYaml()
		if err != nil {
			return fmt.Errorf("error getting yaml status: %v", err)
		}
	default:
		return fmt.Errorf("unknown format %q", cmd.format)
	}

	stopProgress()

	fmt.Print(statusText)

	return nil
}

func (cmd *statusCommand) statusYaml() (string, error) {
	statusStr, err := cmd.statusStruct()
	if err != nil {
		return "", fmt.Errorf("error getting status: %v", err)
	}
	yamlStr, err := yaml.Marshal(statusStr)
	if err != nil {
		return "", fmt.Errorf("error marshalling yaml: %v", err)
	}
	return string(yamlStr), nil
}

func (cmd *statusCommand) statusJson() (string, error) {
	statusStr, err := cmd.statusStruct()
	if err != nil {
		return "", fmt.Errorf("error getting status: %v", err)
	}
	jsonStr, err := json.MarshalIndent(statusStr, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshalling json: %v", err)
	}
	return string(jsonStr), nil
}

type Status struct {
	Engine    string            `json:"engine" yaml:"engine"`
	Endpoints map[string]string `json:"endpoints" yaml:"endpoints"`
}

func (cmd *statusCommand) statusStruct() (*Status, error) {
	var statusStr Status

	activeEngineName, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		return nil, fmt.Errorf("error getting active engine: %v", err)
	}
	if activeEngineName == "" {
		return nil, fmt.Errorf("error no engine is active")
	}
	statusStr.Engine = activeEngineName

	endpoints, err := serverApiUrls(cmd.Context)
	if err != nil {
		return nil, fmt.Errorf("error getting server api endpoints: %v", err)
	}
	statusStr.Endpoints = endpoints

	return &statusStr, nil
}
