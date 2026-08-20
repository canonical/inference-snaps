package debug

import (
	"fmt"
	"path/filepath"

	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/common"
	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/canonical/inference-snaps-cli/v2/pkg/models"
	"github.com/canonical/inference-snaps-cli/v2/pkg/runtimes"
	"github.com/spf13/cobra"
)

type validateCommand struct {
	*common.Context
}

func ValidateCommand(ctx *common.Context) *cobra.Command {
	var cmd validateCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:               "lint-package",
		Short:             "Validate engines, models, and runtimes manifest files",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.run,
	}

	return cobraCmd
}

func (cmd *validateCommand) run(_ *cobra.Command, args []string) error {
	allManifestsValid := true
	snapDir := args[0]

	validators := []struct {
		subdir           string
		manifestFilename string
		validate         func(string) error
	}{
		{"engines", engines.ManifestFilename, engines.Validate},
		{"models", models.ManifestFilename, models.Validate},
		{"runtimes", runtimes.ManifestFilename, runtimes.Validate},
	}

	for _, v := range validators {
		pattern := filepath.Join(snapDir, v.subdir, "*", v.manifestFilename)
		manifestPaths, err := filepath.Glob(pattern)
		if err != nil {
			return fmt.Errorf("invalid glob pattern %s: %w", pattern, err)
		}

		if len(manifestPaths) == 0 {
			allManifestsValid = false
			fmt.Printf("❌ %s: no manifests found\n", filepath.Join(snapDir, v.subdir))
			continue
		}

		for _, manifestPath := range manifestPaths {
			if err := v.validate(manifestPath); err != nil {
				allManifestsValid = false
				fmt.Printf("❌ %s: %s\n", manifestPath, err)
			} else {
				fmt.Printf("✅ %s\n", manifestPath)
			}
		}
	}

	if err := models.ValidateAliases(filepath.Join(snapDir, "engines"), filepath.Join(snapDir, "models")); err != nil {
		allManifestsValid = false
		fmt.Printf("❌ model aliases: %s\n", err)
	} else {
		fmt.Printf("✅ model aliases are unique per engine\n")
	}

	if !allManifestsValid {
		return fmt.Errorf("some manifests are invalid")
	}
	return nil
}
