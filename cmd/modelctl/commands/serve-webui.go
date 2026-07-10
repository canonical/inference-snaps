package commands

import (
	"fmt"

	"github.com/canonical/go-snapctl/env"
	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/common"
	"github.com/canonical/inference-snaps-cli/v2/pkg/models"
	"github.com/canonical/inference-snaps-cli/v2/pkg/webui"
	"github.com/spf13/cobra"
)

type serveWebUiCommand struct {
	*common.Context

	// flags
	port int
	host string // bind address
}

func ServeWebUi(ctx *common.Context) *cobra.Command {
	var cmd serveWebUiCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:               "serve-webui <static-files-dir>",
		Short:             "Serve static files and configurations for the web UI",
		Hidden:            true,
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.serveWebUi,
	}

	// flags
	cobraCmd.Flags().IntVar(&cmd.port, "port", 8081, "HTTP bind port")
	cobraCmd.Flags().StringVar(&cmd.host, "host", "localhost", "HTTP bind address")

	return cobraCmd
}

func (cmd *serveWebUiCommand) serveWebUi(_ *cobra.Command, args []string) error {
	staticDir := args[0]

	// Components are required to get the OpenAI endpoint
	if err := common.WaitForComponents(cmd.Context); err != nil {
		return fmt.Errorf("waiting for component: %s", err)
	}

	baseURL, err := common.OpenAiEndpoint(cmd.Context)
	if err != nil {
		return fmt.Errorf("getting OpenAI base URL: %v", err)
	}

	activeEngineName, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		return fmt.Errorf("getting active engine: %v", err)
	}
	if activeEngineName == "" {
		return common.ErrNoActiveEngine
	}

	activeModelId, err := cmd.Cache.GetActiveModel()
	if err != nil {
		return fmt.Errorf("getting active model: %w", err)
	}

	// Always send an array (possibly empty) so the /config response has a
	// stable shape and never serializes capabilities as JSON null.
	capabilities := []string{}
	if activeModelId != "" {
		modelManifest, err := models.LoadManifest(cmd.ModelsDir, activeModelId)
		if err != nil {
			return fmt.Errorf("loading model manifest: %w", err)
		}
		if modelManifest.Capabilities != nil {
			capabilities = modelManifest.Capabilities
		}
	}

	config := webui.Config{
		OpenAIBaseURL: baseURL,
		Capabilities:  capabilities,
		InstanceName:  env.SnapInstanceName(),
		EngineName:    activeEngineName,
	}

	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %v", err)
	}

	if cmd.Verbose {
		fmt.Println("OpenAI base URL:", config.OpenAIBaseURL)
		fmt.Println("Capabilities:", config.Capabilities)
		fmt.Println("Instance name:", config.InstanceName)
		fmt.Println("Engine name:", config.EngineName)
	}

	fmt.Printf("Serving %q on http://localhost:%d\n", staticDir, cmd.port)
	return webui.Serve(config, staticDir, cmd.port, cmd.host)
}
