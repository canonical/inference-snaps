package commands

import (
	"fmt"

	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/common"
	"github.com/spf13/cobra"
)

type chatCommand struct {
	*common.Context
}

func Chat(ctx *common.Context) *cobra.Command {
	var cmd chatCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:               "chat",
		Short:             "Start the chat CLI",
		Long:              "Chat with the server via its OpenAI API.\nThis CLI supports text-based prompting only.",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.run,
	}

	return cobraCmd
}

func (cmd *chatCommand) run(_ *cobra.Command, _ []string) error {
	chatBaseUrl, err := common.OpenAiEndpoint(cmd.Context)
	if err != nil {
		return fmt.Errorf("getting OpenAI base URL: %v", err)
	}

	if cmd.Snap.InstanceName() != "" {
		// TODO: get app name dynamically
		services, err := cmd.Snap.ServiceStatuses()
		if err != nil {
			return fmt.Errorf("getting services: %v", err)
		}
		status, ok := services["server"]
		if !ok {
			return fmt.Errorf("server service not found")
		}
		if status == "inactive" {
			return fmt.Errorf("server not active\n\n%s",
				common.SuggestStartServer())
		}
	}

	chatClient := common.ChatClient(chatBaseUrl, "", cmd.Verbose)

	return chatClient.Start()
}
