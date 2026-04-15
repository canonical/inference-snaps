package debug

import (
	"encoding/json"
	"fmt"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/webui"
	"github.com/spf13/cobra"
)

type serveWebUiCommand struct {
	*common.Context

	// flags
	baseUrl string
	port    int
	host    string
	htmlDir string
}

func ServeWebUiCommand(ctx *common.Context) *cobra.Command {
	var cmd serveWebUiCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:               "serve-webui <static-files-dir>",
		Short:             "Serve web UI static files and configurations for debugging",
		Hidden:            true,
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.serveWebUi,
	}

	// flags
	cobraCmd.Flags().StringVar(&cmd.baseUrl, "base-url", "http://localhost:8080/v1", "Base URL of the OpenAI-compatible server")
	cobraCmd.Flags().IntVar(&cmd.port, "port", 8081, "HTTP bind port of the web server")
	cmd.host = "localhost" // fixed to localhost as this is for debugging only

	return cobraCmd
}

func (cmd *serveWebUiCommand) serveWebUi(_ *cobra.Command, args []string) error {
	staticDir := args[0]

	config := webui.Config{
		OpenAIBaseURL: cmd.baseUrl,
		Capabilities:  webui.SupportedCapabilities(), // set all capabilities for debugging
		InstanceName:  "debug",
		EngineName:    "unset",
	}

	j, _ := json.MarshalIndent(config, "", "  ")
	fmt.Printf("Config: %s\n", j)

	fmt.Printf("Serving %q on http://localhost:%d\n", staticDir, cmd.port)
	return webui.Serve(config, staticDir, cmd.port, cmd.host)
}
