package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/spf13/cobra"
)

type runCommand struct {
	*common.Context

	// flags
	waitForComponents bool
}

func Run(ctx *common.Context) *cobra.Command {
	var cmd runCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:               "run <path>",
		Short:             "Run a subprocess",
		Hidden:            true,
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.run,
	}

	// flags
	cobraCmd.Flags().BoolVar(&cmd.waitForComponents, "wait-for-components", false, "wait for engine components to be installed before running")

	return cobraCmd
}

func (cmd *runCommand) run(_ *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("unexpected number of arguments, expected 1 got %d", len(args))
	}
	if cmd.waitForComponents {
		if err := common.WaitForComponents(cmd.Context); err != nil {
			return fmt.Errorf("error waiting for component: %s", err)
		}
	}

	err := common.LoadEngineEnvironment(cmd.Context)
	if err != nil {
		return fmt.Errorf("error loading engine environment: %v", err)
	}

	path := args[0]

	execCmd := exec.Command(path)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	return execCmd.Run()
}
