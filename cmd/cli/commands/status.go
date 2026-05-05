package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type statusCommand struct {
	*common.Context

	// flags
	format            string
	waitForComponents bool
}

func Status(ctx *common.Context) *cobra.Command {
	var cmd statusCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:               "status",
		Short:             "Show the status",
		Long:              "Show the status of the inference snap",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.run,
	}

	// flags
	supportedFormats := []string{"json", "yaml"}
	cobraCmd.Flags().StringVar(
		&cmd.format,
		"format",
		"yaml",
		fmt.Sprintf("output format (%s)", strings.Join(supportedFormats, ", ")),
	)
	cobraCmd.Flags().BoolVar(&cmd.waitForComponents, "wait-for-components", false, "wait for engine components to be installed before reporting status")

	return cobraCmd
}

func (cmd *statusCommand) run(_ *cobra.Command, _ []string) error {
	var statusText string
	var err error

	if cmd.waitForComponents {
		if err := common.WaitForComponents(cmd.Context); err != nil {
			return fmt.Errorf("waiting for component: %s", err)
		}
	}

	stopProgress := common.StartProgressSpinner("Getting status")
	defer stopProgress()

	switch cmd.format {
	case "json":
		statusText, err = cmd.statusJson()
		if err != nil {
			return fmt.Errorf("getting json status: %v", err)
		}
		statusText += "\n"
	case "yaml":
		statusText, err = cmd.statusYaml()
		if err != nil {
			return fmt.Errorf("getting yaml status: %v", err)
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
		return "", fmt.Errorf("getting status: %v", err)
	}
	yamlStr, err := yaml.Marshal(statusStr)
	if err != nil {
		return "", fmt.Errorf("marshalling yaml: %v", err)
	}
	return string(yamlStr), nil
}

func (cmd *statusCommand) statusJson() (string, error) {
	statusStr, err := cmd.statusStruct()
	if err != nil {
		return "", fmt.Errorf("getting status: %v", err)
	}
	jsonStr, err := json.MarshalIndent(statusStr, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshalling json: %v", err)
	}
	return string(jsonStr), nil
}

type status struct {
	Engine    string            `json:"engine" yaml:"engine"`
	Services  map[string]string `json:"services" yaml:"services"`
	Endpoints map[string]string `json:"endpoints,omitempty" yaml:"endpoints,omitempty"`
	Model     map[string]string `json:"model,omitempty" yaml:"model,omitempty"`
}

func (cmd *statusCommand) statusStruct() (*status, error) {
	var statusStr status

	activeEngineName, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", common.LookingUpActiveEngine, err)
	}
	if activeEngineName == "" {
		return nil, common.ErrNoActiveEngine
	}
	statusStr.Engine = activeEngineName

	services, err := common.ServiceStatuses()
	if err != nil {
		return nil, fmt.Errorf("getting service statuses: %v", err)
	}
	statusStr.Services = services

	endpoints, err := common.ServerEndpoints(cmd.Context)
	if err != nil {
		return nil, fmt.Errorf("getting server api endpoints: %v", err)
	}
	statusStr.Endpoints = endpoints

	modelStatus, err := common.ModelStatus(cmd.Context)
	if err != nil {
		return nil, fmt.Errorf("getting model status: %v", err)
	}
	statusStr.Model = modelStatus

	return &statusStr, nil
}
