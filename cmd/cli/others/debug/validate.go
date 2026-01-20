package debug

import (
	"fmt"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/spf13/cobra"
)

type validateCommand struct {
	*common.Context
}

func ValidateCommand(ctx *common.Context) *cobra.Command {
	var cmd validateCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:               "validate-engines",
		Short:             "Validate engine manifest files",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.run,
	}

	return cobraCmd
}

func (cmd *validateCommand) run(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no engine manifest specified")
	}

	allManifestsValid := true
	for _, manifestPath := range args {
		err := engines.Validate(manifestPath)
		if err != nil {
			allManifestsValid = false
			fmt.Printf("❌ %s: %s\n", manifestPath, err)
		} else {
			fmt.Printf("✅ %s\n", manifestPath)
		}
	}

	if !allManifestsValid {
		return fmt.Errorf("not all manifests are valid")
	}
	return nil
}
