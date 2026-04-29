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
		Use:               "set <key=value>...",
		Short:             "Set configurations",
		Long:              "Set a configuration",
		Args:              cobra.MinimumNArgs(1),
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

	return cmd.set(args)
}

func (cmd *setCommand) set(keyValuePairs []string) error {
	keyValues, err := cmd.parseKeyValues(keyValuePairs)
	if err != nil {
		return err
	}

	switch {
	case cmd.packageConfig:
		return cmd.setPackageConfigs(keyValues)
	case cmd.engineConfig:
		return cmd.setEngineConfigs(keyValues)
	default:
		return cmd.setUserConfigs(keyValues)
	}
}

func (cmd *setCommand) setPackageConfigs(keyValues map[string]string) error {
	for k, v := range keyValues {
		if err := cmd.Config.Set(k, v, storage.PackageConfig); err != nil {
			return fmt.Errorf("setting %q to %q: %v", k, v, err)
		}
	}
	return nil
}

func (cmd *setCommand) setEngineConfigs(keyValues map[string]string) error {
	for k, v := range keyValues {
		if err := cmd.Config.Set(k, v, storage.EngineConfig); err != nil {
			return fmt.Errorf("setting %q to %q: %v", k, v, err)
		}
	}
	return nil
}

func (cmd *setCommand) setUserConfigs(keyValues map[string]string) error {

	currentValues := map[string]string{}
	currentKnown := map[string]bool{}

	// Validate key values
	for key := range keyValues {
		currentValue, found, err := cmd.getCurrentValue(key)
		if err != nil {
			return err
		}

		currentValues[key] = currentValue
		currentKnown[key] = found
	}

	// Apply configurations
	anyChange := false
	for k, v := range keyValues {
		if err := cmd.Config.Set(k, v, storage.UserConfig); err != nil {
			return fmt.Errorf("setting %q to %q: %v", k, v, err)
		}

		// User keys are known, except for passthrough keys
		if !currentKnown[k] || currentValues[k] != v {
			anyChange = true
		}
	}

	// Restart if configurations were changed
	if anyChange {
		if err := cmd.restartToApply(); err != nil {
			return err
		}
	}

	return nil
}

func (cmd *setCommand) parseKeyValues(keyValues []string) (map[string]string, error) {
	kvMap := map[string]string{}
	seenKeys := map[string]bool{}

	for _, keyValue := range keyValues {
		key, value, err := cmd.parseKeyValue(keyValue)
		if err != nil {
			return nil, err
		}
		kvMap[key] = value

		if seenKeys[key] {
			return nil, fmt.Errorf("duplicate key: %q", key)
		}
		seenKeys[key] = true
	}
	return kvMap, nil
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

func (cmd *setCommand) getCurrentValue(key string) (string, bool, error) {
	currValMap, err := cmd.Config.Get(key)
	if err != nil {
		return "", false, fmt.Errorf("checking existing keys: %s", err)
	}
	currVal, found := currValMap[key]
	if !found && !strings.HasPrefix(key, "passthrough.") {
		return "", false, fmt.Errorf("key %q is not found\n\n%s", key, common.SuggestKeyNotFound(key))
	}

	if !found {
		return "", false, nil
	}

	return fmt.Sprint(currVal), true, nil
}

func (cmd *setCommand) restartToApply() error {
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
