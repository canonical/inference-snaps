package commands

import (
	"fmt"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/storage"
	"github.com/canonical/inference-snaps-cli/pkg/utils"
	"github.com/spf13/cobra"
)

type unsetCommand struct {
	*common.Context

	// flags
	packageConfig bool
	engineConfig  bool
	assumeYes     bool
	noRestart     bool
}

func Unset(ctx *common.Context) *cobra.Command {
	var cmd unsetCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:               "unset <key>",
		Short:             "Unset configurations",
		Long:              "Unset a user configuration, reverting to package or engine default. If no default exists for the key, it will be removed entirely.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.run,
	}

	cobraCmd.Flags().BoolVar(&cmd.assumeYes, "assume-yes", false, "assume yes for all prompts")
	cobraCmd.Flags().BoolVar(&cmd.noRestart, "no-restart", false, "do not restart the snap after setting the configuration")

	return cobraCmd
}

func (cmd *unsetCommand) run(_ *cobra.Command, args []string) error {
	if !utils.IsRootUser() {
		return common.ErrPermissionDenied
	}

	return cmd.unsetValue(args[0])
}

func (cmd *unsetCommand) unsetValue(key string) error {
	currValMap, err := cmd.Config.Get(key)
	if err != nil {
		return fmt.Errorf("checking existing keys: %s", err)
	}
	currVal, found := currValMap[key]
	if !found {
		return fmt.Errorf("key %q is not found\n\n%s", key, common.SuggestKeyNotFound(key))
	}

	err = cmd.Config.Unset(key, storage.UserConfig)
	if err != nil {
		return fmt.Errorf("unsetting %q: %v", key, err)
	}

	newValMap, err := cmd.Config.Get(key)
	if err != nil {
		return fmt.Errorf("checking existing keys: %s", err)
	}
	newVal := newValMap[key]

	if fmt.Sprint(currVal) == fmt.Sprint(newVal) {
		return nil // value not changed
	}

	if !cmd.noRestart {
		msg := fmt.Sprintf("Restart %s to apply the changes?", cmd.Snap.InstanceName())
		if cmd.assumeYes || common.PromptYN(msg, true) {
			if err := cmd.Snap.Restart(); err != nil {
				return fmt.Errorf("restarting snap: %v", err)
			}
		}
	}
	return nil
}
