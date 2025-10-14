package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/canonical/inference-snaps-cli/pkg/storage"
	"github.com/canonical/inference-snaps-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	setPackageConfig bool
	setEngineConfig  bool
)

func addSetCommand() {
	cmd := &cobra.Command{
		Use:   "set <key>",
		Short: "Set configuration option",
		// Long:  "",
		GroupID:           "config",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              set,
	}

	// flags
	cmd.PersistentFlags().BoolVar(&setPackageConfig, "package", false, "set package configurations")
	cmd.PersistentFlags().MarkHidden("package")
	cmd.PersistentFlags().BoolVar(&setEngineConfig, "engine", false, "set engine configuration")
	cmd.PersistentFlags().MarkHidden("engine")

	rootCmd.AddCommand(cmd)
}

func set(_ *cobra.Command, args []string) error {
	if !utils.IsRootUser() {
		return ErrPermissionDenied
	}
	return setValue(args[0])
}

func setValue(keyValue string) error {
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
	if setPackageConfig {
		err = config.Set(key, value, storage.PackageConfig)
	} else if setEngineConfig {
		err = config.Set(key, value, storage.EngineConfig)
	} else {
		// Reject use of internal keys by the user
		if slices.Contains(deprecatedConfig, key) {
			return fmt.Errorf("%q is read-only", key)
		}
		err = config.Set(key, value, storage.UserConfig)
	}
	if err != nil {
		return fmt.Errorf("error setting value %q for %q: %v", value, key, err)
	}

	return nil
}
