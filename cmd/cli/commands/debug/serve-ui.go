package debug

import (
	"encoding/json"
	"fmt"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/ui"
	"github.com/spf13/cobra"
)

type serveUICommand struct {
	*common.Context

	// flags
	baseUrl string
	port    int
	host    string
	htmlDir string
}

func ServeUICommand(ctx *common.Context) *cobra.Command {
	var cmd serveUICommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:               "serve-ui <static-files-dir>",
		Short:             "Run a debug web server to serve the Web UI",
		Hidden:            true,
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.serveUI,
	}

	// flags
	cobraCmd.Flags().StringVar(&cmd.baseUrl, "base-url", "http://localhost:8080/v1", "Base URL of the OpenAI-compatible server")
	cobraCmd.Flags().IntVar(&cmd.port, "port", 8081, "HTTP bind port of the web server")
	cmd.host = "localhost" // fixed to localhost as this is for debugging only

	return cobraCmd
}

func (cmd *serveUICommand) serveUI(_ *cobra.Command, args []string) error {
	staticDir := args[0]

	config := ui.Config{
		OpenAIBaseURL: cmd.baseUrl,
		Capabilities:  ui.SupportedCapabilities(), // set all capabilities for debugging
		InstanceName:  "debug",
		EngineName:    "unset",
	}

	j, _ := json.MarshalIndent(config, "", "  ")
	fmt.Printf("Config: %s\n", j)

	fmt.Printf("Serving %q on http://localhost:%d\n", staticDir, cmd.port)
	return ui.Serve(config, staticDir, cmd.port, cmd.host)
}
