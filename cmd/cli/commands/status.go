package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/canonical/go-snapctl"
	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type statusCommand struct {
	*common.Context

	// flags
	format string
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

	return cobraCmd
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

type status struct {
	Engine    string            `json:"engine" yaml:"engine"`
	Services  map[string]string `json:"services" yaml:"services"`
	Endpoints map[string]string `json:"endpoints" yaml:"endpoints"`
}

func (cmd *statusCommand) statusStruct() (*status, error) {
	var statusStr status

	activeEngineName, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		return nil, fmt.Errorf("error getting active engine: %v", err)
	}
	if activeEngineName == "" {
		return nil, fmt.Errorf("error no engine is active")
	}
	statusStr.Engine = activeEngineName

	services, err := snapctl.Services().Run()
	if err != nil {
		return nil, fmt.Errorf("error getting services: %v", err)
	}
	statusStr.Services = make(map[string]string)
	for name, service := range services {
		// The service name is in the format <snap-name>.<service-app>, we only want the service-app part.
		_, serviceApp, found := strings.Cut(name, ".")
		if !found {
			return nil, fmt.Errorf("error unexpected service name format: %q", name)
		}
		// Append the service status exactly as snapd reports it. Often this is in the host system language, see bug:
		// https://bugs.launchpad.net/snapd/+bug/2137543
		statusStr.Services[serviceApp] = service.Current
	}

	endpoints, err := common.ServerApiUrls(cmd.Context)
	if err != nil {
		return nil, fmt.Errorf("error getting server api endpoints: %v", err)
	}
	statusStr.Endpoints = endpoints

	return &statusStr, nil
}
