package others

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/spf13/cobra"
)

type runCommand struct {
	*common.Context

	// flags
	waitForComponentsFlag bool
}

func RunCommand(ctx *common.Context) *cobra.Command {
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
	cobraCmd.Flags().BoolVar(&cmd.waitForComponentsFlag, "wait-for-components", false, "wait for engine components to be installed before running")

	return cobraCmd
}

func (cmd *runCommand) run(_ *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("unexpected number of arguments, expected 1 got %d", len(args))
	}

	if cmd.waitForComponentsFlag {
		if err := cmd.waitForComponents(); err != nil {
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

// TODO: unify with similar code in use.go
func (cmd *runCommand) checkMissingComponents(manifest *engines.Manifest) ([]string, error) {
	componentsDir, found := os.LookupEnv("SNAP_COMPONENTS")
	if !found {
		return nil, fmt.Errorf("SNAP_COMPONENTS env var not set")
	}

	var missing []string
	for _, component := range manifest.Components {
		componentPath := filepath.Join(componentsDir, component)
		if _, err := os.Stat(componentPath); os.IsNotExist(err) {
			missing = append(missing, component)
		}
	}

	return missing, nil
}

func (cmd *runCommand) waitForComponents() error {
	const maxWait = 3600 // seconds
	const interval = 10  // seconds

	activeEngineName, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		return fmt.Errorf("error looking up active engine: %v", err)
	}

	if activeEngineName == "" {
		return fmt.Errorf("no active engine")
	}

	manifest, err := engines.LoadManifest(cmd.EnginesDir, activeEngineName)
	if err != nil {
		return fmt.Errorf("error loading engine manifest: %v", err)
	}

	missing, err := cmd.checkMissingComponents(manifest)
	if err != nil {
		return err
	}

	for elapsed := 0; elapsed < maxWait && len(missing) > 0; elapsed += interval {
		fmt.Printf("Waiting for required snap components: %s (%d/%ds)\n",
			strings.Join(missing, ", "), elapsed, maxWait)

		time.Sleep(interval * time.Second)

		missing, err = cmd.checkMissingComponents(manifest)
		if err != nil {
			return err
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("timed out after %ds while waiting for required components: %s",
			maxWait, strings.Join(missing, ", "))
	}

	return nil
}
