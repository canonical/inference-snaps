package debug

import (
	"fmt"

	"github.com/canonical/inference-snaps-cli/cmd/cli/basic/chat"
	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/spf13/cobra"
)

type chatCommand struct {
	*common.Context

	// flags
	baseUrl   string
	modelName string
}

func ChatCommand(ctx *common.Context) *cobra.Command {
	var cmd chatCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:               "chat",
		Short:             "Start the chat CLI providing connection parameters",
		Long:              "Open a text-only chat session to the OpenAI-compatible server at the provided URL, requesting the model as specified.",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.run,
	}

	// flags
	cobraCmd.Flags().StringVar(&cmd.baseUrl, "base-url", "", "Base URL of the OpenAI-compatible server")
	cobraCmd.Flags().StringVar(&cmd.modelName, "model", "", "Name of the model to use")

	return cobraCmd
}

func (cmd *chatCommand) run(_ *cobra.Command, args []string) error {
	if cmd.baseUrl == "" {
		return fmt.Errorf("the --base-url parameter is required")
	}

	return chat.Client(cmd.baseUrl, cmd.modelName, cmd.Verbose)
}
