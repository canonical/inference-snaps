package basic

import (
	"fmt"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/chat"
	"github.com/spf13/cobra"
)

type chatCommand struct {
	*common.Context
}

func ChatCommand(ctx *common.Context) *cobra.Command {
	var cmd chatCommand
	cmd.Context = ctx

	cobra := &cobra.Command{
		Use:               "chat",
		Short:             "Start the chat CLI",
		Long:              "Chat with the server via its OpenAI API.\nThis CLI supports text-based prompting only.",
		GroupID:           groupID,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.run,
	}

	return cobra
}

func (cmd *chatCommand) run(_ *cobra.Command, _ []string) error {
	apiUrls, err := serverApiUrls(cmd.Context)
	if err != nil {
		return fmt.Errorf("error getting server api urls: %v", err)
	}
	chatBaseUrl := apiUrls[openAi]

	return chat.Client(chatBaseUrl, "")
}
