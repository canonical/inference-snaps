package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/canonical/famous-models-cli/pkg/engines"
	"github.com/canonical/famous-models-cli/pkg/selector"
	"github.com/fatih/color"
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
		Long:              "Show the status of the model snap",
		GroupID:           "basics",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              status,
	}

	// flags
	cmd.PersistentFlags().StringVar(&statusFormat, "format", "", "return the status as yaml or json")

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
	case "yaml":
		statusText, err = statusYaml()
		if err != nil {
			return fmt.Errorf("error getting yaml status: %v", err)
		}
	default:
		statusText, err = statusHuman()
		if err != nil {
			return fmt.Errorf("error getting status: %v", err)
		}
	}

	stopProgress()

	fmt.Println(statusText)

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

func statusHuman() (string, error) {
	// Find the selected engine
	engine, err := activeEngine()
	if err != nil {
		return "", fmt.Errorf("error loading active engine: %v", err)
	}

	// Get all engines with scores
	compatibleEngines := true
	scoredEngines, err := scoreEngines()
	if err != nil {
		return "", fmt.Errorf("error scoring engines: %v", err)
	}

	// Find top engine
	autoEngine, err := selector.TopEngine(scoredEngines)
	if err != nil {
		if errors.Is(err, selector.ErrorNoCompatibleEngine) {
			compatibleEngines = false
		} else {
			return "", fmt.Errorf("error loading top engine: %v", err)
		}
	}

	engineStatusText := statusHumanEngine(engine, compatibleEngines && engine.Name == autoEngine.Name)

	// TODO check if all required snap components are available, otherwise print "Downloading resources..."

	serverStatusText, err := statusHumanServer(engine)
	if err != nil {
		return "", fmt.Errorf("error getting server status: %v", err)
	}

	return fmt.Sprintf("%s\n\n%s", engineStatusText, serverStatusText), nil
}

func statusHumanEngine(engine *engines.ScoredManifest, auto bool) string {
	bold := color.New(color.Bold).SprintFunc()
	engineString := fmt.Sprintf("Using %s", bold(engine.Name))
	if auto {
		engineString += " (automatically selected)"
	}
	return engineString
}

func statusHumanServer(engine *engines.ScoredManifest) (string, error) {
	// Start, stop, log commands
	startCmd := fmt.Sprintf(`Run "sudo snap start %s" to start the server.`, snapInstanceName)
	stopCmd := fmt.Sprintf(`Run "sudo snap stop %s" to stop the server.`, snapInstanceName)
	logsCmd := fmt.Sprintf(`Run "sudo snap logs %s" to view the server logs.`, snapInstanceName)

	apiUrls, err := serverApiUrls(engine)
	if err != nil {
		return "", fmt.Errorf("error getting api urls: %v", err)
	}

	checkExitCode, err := serverStatusCode(engine.Name)
	if err != nil {
		return "", fmt.Errorf("error checking server status: %v", err)
	}

	statusText := ""
	switch checkExitCode {
	case 0:
		if apiUrl, ok := apiUrls[openAi]; ok {
			statusText += fmt.Sprintf("OpenAI API at %s\n", apiUrl)
		}
		// TODO if we can detect that the server is OVMS, we can add "OpenVINO API at http://localhost:8080/v1"
		statusText += "\n"
		statusText += fmt.Sprintf("%s", stopCmd)

	case 1:
		statusText = "Starting runtime...\n"
	case 2:
		statusText = "Runtime stopped.\n"
		statusText += "\n"
		statusText += fmt.Sprintf("%s", startCmd)
	default:
		statusText = fmt.Sprintf("Runtime error: unknown exit code %d\n", checkExitCode)
		statusText += "\n"
		statusText += fmt.Sprintf("%s", logsCmd)
	}

	return statusText, nil
}
