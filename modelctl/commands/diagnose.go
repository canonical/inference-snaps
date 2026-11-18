package commands

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/canonical/inference-snaps/modelctl/v2/commands/common"
	"github.com/canonical/inference-snaps/modelctl/v2/internal/fallbackserver"
	"github.com/spf13/cobra"
)

type diagnoseCommand struct {
	*common.Context

	serveFallback string
}

func Diagnose(ctx *common.Context) *cobra.Command {
	cmd := diagnoseCommand{Context: ctx}
	cobraCmd := &cobra.Command{
		Use:               "diagnose",
		Short:             "Check that an active engine and model are available",
		Hidden:            true,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.run,
	}
	cobraCmd.Flags().StringVar(&cmd.serveFallback, "serve-fallback", "", "serve an error response on failure (openai)")
	return cobraCmd
}

func (cmd *diagnoseCommand) run(cobraCmd *cobra.Command, _ []string) error {
	_, diagnosisErr := common.EngineSettings(cmd.Context)
	switch cmd.serveFallback {
	case "":
		return diagnosisErr
	case "openai":
		if diagnosisErr != nil {
			fmt.Fprintf(cobraCmd.ErrOrStderr(), "Error: %v\n", diagnosisErr)
			address, err := cmd.fallbackListenAddress()
			if err != nil {
				return err
			}
			return cmd.startFallbackServer(address, diagnosisErr)
		}
	default:
		return fmt.Errorf("unsupported fallback protocol %q", cmd.serveFallback)
	}
	return diagnosisErr
}

func (cmd *diagnoseCommand) fallbackListenAddress() (string, error) {
	config, err := cmd.Config.Get("http")
	if err != nil {
		return "", fmt.Errorf("getting fallback HTTP configuration: %w", err)
	}
	var values [2]string
	for i, key := range []string{"http.host", "http.port"} {
		value, ok := config[key]
		if !ok || value == nil || strings.TrimSpace(fmt.Sprint(value)) == "" {
			return "", fmt.Errorf("config %q is not set", key)
		}
		values[i] = strings.TrimSpace(fmt.Sprint(value))
	}
	return net.JoinHostPort(values[0], values[1]), nil
}

func (cmd *diagnoseCommand) startFallbackServer(address string, diagnosisErr error) error {
	message := "Inference server failed to start. Please check the logs for more information."
	switch {
	case errors.Is(diagnosisErr, common.ErrNoActiveEngine):
		message = "No active engine is set. Please set an active engine and try again."
	case errors.Is(diagnosisErr, common.ErrNoActiveModel):
		message = "No active model is set. Please set an active model and try again."
	}

	fmt.Println("Starting fallback server...")
	if err := fallbackserver.Run(address, message); err != nil {
		return fmt.Errorf("running fallback server: %w", err)
	}
	return nil
}
