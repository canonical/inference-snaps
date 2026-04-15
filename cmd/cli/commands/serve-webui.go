package commands

import (
	"fmt"
	"strings"

	"github.com/canonical/go-snapctl/env"
	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/webui"
	"github.com/spf13/cobra"
)

type serveWebUiCommand struct {
	*common.Context

	// flags
	port         int
	host         string // bind address
	capabilities string
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
	cobraCmd.Flags().StringVar(&cmd.capabilities, "capabilities", "text",
		fmt.Sprintf("Comma-separated list of capabilities (%s)", strings.Join(webui.SupportedCapabilities(), ", ")))

	return cobraCmd
}

func (cmd *serveWebUiCommand) serveWebUi(_ *cobra.Command, args []string) error {
	staticDir := args[0]

	baseURL, err := common.OpenAiEndpoint(cmd.Context)
	if err != nil {
		return fmt.Errorf("getting OpenAI base URL: %v", err)
	}

	activeEngineName, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		return fmt.Errorf("getting active engine: %v", err)
	}
	if activeEngineName == "" {
		return fmt.Errorf("no engine is active")
	}

	var capabilities []string
	for cap := range strings.SplitSeq(cmd.capabilities, ",") {
		capabilities = append(capabilities, strings.TrimSpace(cap))
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
