package commands

import (
	"fmt"

	"github.com/canonical/go-snapctl"
	"github.com/canonical/go-snapctl/env"
	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
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
	apiUrls, err := common.ServerApiUrls(cmd.Context)
	if err != nil {
		return fmt.Errorf("error getting server api urls: %v", err)
	}
	chatBaseUrl := apiUrls[common.OpenAiEndpointKey]

	if env.SnapInstanceName() != "" {
		// TODO: get app name dynamically
		serviceName := env.SnapInstanceName() + ".server"
		services, err := snapctl.Services(serviceName).Run()
		if err != nil {
			return fmt.Errorf("error getting services: %v", err)
		}
		if services[serviceName].Current == "inactive" {
			return fmt.Errorf("server not active\n\n%s",
				common.SuggestStartServer())
		}
	}

	chatClient := common.ChatClient(chatBaseUrl, "", cmd.Verbose)

	return chatClient.Start()
}
