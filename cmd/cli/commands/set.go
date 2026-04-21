package commands

import (
	"fmt"
	"strings"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/storage"
	"github.com/canonical/inference-snaps-cli/pkg/utils"
	"github.com/spf13/cobra"
)

type setCommand struct {
	*common.Context

	// flags
	packageConfig bool
	engineConfig  bool
	assumeYes     bool
	noRestart     bool
}

func Set(ctx *common.Context) *cobra.Command {
	var cmd setCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:               "set <key=value>",
		Short:             "Set configurations",
		Long:              "Set a configuration",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.run,
	}

	// flags
	cobraCmd.Flags().BoolVar(&cmd.packageConfig, "package", false, "set package configurations")
	if err := cobraCmd.Flags().MarkHidden("package"); err != nil {
		panic(err)
	}
	cobraCmd.Flags().BoolVar(&cmd.engineConfig, "engine", false, "set engine configuration")
	if err := cobraCmd.Flags().MarkHidden("engine"); err != nil {
		panic(err)
	}
	cobraCmd.Flags().BoolVar(&cmd.assumeYes, "assume-yes", false, "assume yes for all prompts")
	cobraCmd.Flags().BoolVar(&cmd.noRestart, "no-restart", false, "do not restart the snap after setting the configuration")

	return cobraCmd
}

func (cmd *setCommand) run(_ *cobra.Command, args []string) error {
	if !utils.IsRootUser() {
		return common.ErrPermissionDenied
	}

	return cmd.setValue(args[0])
}

func (cmd *setCommand) setValue(keyValue string) error {
	key, value, err := cmd.parseKeyValue(keyValue)
	if err != nil {
		return err
	}

	if cmd.packageConfig {
		if err := cmd.Config.Set(key, value, storage.PackageConfig); err != nil {
			return fmt.Errorf("setting %q to %q: %v", key, value, err)
		}
	} else if cmd.engineConfig {
		if err := cmd.Config.Set(key, value, storage.EngineConfig); err != nil {
			return fmt.Errorf("setting %q to %q: %v", key, value, err)
		}
	}

	return cmd.setUserConfig(key, value)
}

func (cmd *setCommand) parseKeyValue(keyValue string) (key, value string, err error) {
	if keyValue == "" {
		return "", "", fmt.Errorf("expected key=value, got %q", keyValue)
	}

	if keyValue[0] == '=' {
		return "", "", fmt.Errorf("key must not start with an equal sign")
	}

	// The value itself can contain an equal sign, so we split only on the first occurrence
	parts := strings.SplitN(keyValue, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("expected key=value, got %q", keyValue)
	}
	return parts[0], parts[1], nil
}

// setUserConfig overrides existing package and engine configurations.
// Unknown keys are rejected, except for passthrough keys.
func (cmd *setCommand) setUserConfig(key, value string) error {

	currValMap, err := cmd.Config.Get(key)
	if err != nil {
		return fmt.Errorf("checking existing keys: %s", err)
	}
	currVal, found := currValMap[key]
	if !found && !strings.HasPrefix(key, "passthrough.") {
		return fmt.Errorf("unknown key: %q", key)
	}

	if fmt.Sprint(currVal) == value {
		return nil // no change needed
	}

	if err := cmd.Config.Set(key, value, storage.UserConfig); err != nil {
		return fmt.Errorf("setting %q to %q: %v", key, value, err)
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
