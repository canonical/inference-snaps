package config

import (
	"fmt"
	"slices"
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
}

func SetCommand(ctx *common.Context) *cobra.Command {
	var cmd setCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:               "set <key=value>",
		Short:             "Set configurations",
		Long:              "Set a configuration",
		GroupID:           groupID,
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.run,
	}

	// flags
	cobraCmd.Flags().BoolVar(&cmd.packageConfig, "package", false, "set package configurations")
	err := cobraCmd.Flags().MarkHidden("package")
	if err != nil {
		panic(err)
	}
	cobraCmd.Flags().BoolVar(&cmd.engineConfig, "engine", false, "set engine configuration")
	err = cobraCmd.Flags().MarkHidden("engine")
	if err != nil {
		panic(err)
	}

	return cobraCmd
}

func (cmd *setCommand) run(_ *cobra.Command, args []string) error {
	if !utils.IsRootUser() {
		return common.ErrPermissionDenied
	}
	return cmd.setValue(args[0])
}

func (cmd *setCommand) setValue(keyValue string) error {
	if keyValue[0] == '=' {
		return fmt.Errorf("key must not start with an equal sign")
	}

	// The value itself can contain an equal sign, so we split only on the first occurrence
	parts := strings.SplitN(keyValue, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("expected key=value, got %q", keyValue)
	}
	key, value := parts[0], parts[1]

	var err error
	if cmd.packageConfig {
		err = cmd.Config.Set(key, value, storage.PackageConfig)
	} else if cmd.engineConfig {
		err = cmd.Config.Set(key, value, storage.EngineConfig)
	} else {
		// Reject use of internal keys by the user
		if slices.Contains(deprecatedConfig, key) {
			return fmt.Errorf("%q is read-only", key)
		}
		err = cmd.Config.Set(key, value, storage.UserConfig)
	}
	if err != nil {
		return fmt.Errorf("error setting value %q for %q: %v", value, key, err)
	}

	return nil
}
