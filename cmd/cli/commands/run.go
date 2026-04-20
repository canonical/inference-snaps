package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/utils"
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
			return fmt.Errorf("waiting for component: %s", err)
		}
	}

	clean, err := common.LoadEngineEnvironment(cmd.Context)
	if err != nil {
		return fmt.Errorf("loading engine environment: %v", err)
	}

	// NOTE: defer does not run on SIGTERM or SIGKILL. It only runs when the child process exits.
	// TODO: add signal handling to intercept SIGTERM and invoke clean() before exiting.
	defer clean()

	if err := cmd.processPassthroughConfigs(); err != nil {
		return fmt.Errorf("processing passthrough configs: %v", err)
	}

	path := args[0]

	execCmd := exec.Command(path)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	return execCmd.Run()
}

func (cmd *runCommand) getEnvVarsFromPassthroughConfigs(envVars map[string]any) (map[string]any, error) {
	result := make(map[string]any)
	const keyPrefix = "passthrough.environment."
	for k, v := range envVars {
		if strings.HasPrefix(k, keyPrefix) {
			envVarName := strings.TrimPrefix(k, keyPrefix)
			envVarValue := fmt.Sprintf("%v", v)
			// Convert passthrough keys (my-key) to environment variables names (MY_KEY)
			envVarName = strings.ToUpper(strings.ReplaceAll(envVarName, "-", "_"))
			result[envVarName] = envVarValue
		}
	}

	return result, nil
}

func (cmd *runCommand) processPassthroughConfigs() error {
	passthroughConfigs, err := cmd.Config.Get("passthrough")
	if err != nil {
		return fmt.Errorf("getting configs: %v", err)
	}
	envVars, err := cmd.getEnvVarsFromPassthroughConfigs(passthroughConfigs)
	if err != nil {
		return err
	}
	err = utils.SetEnvironmentVariables(envVars)
	if err != nil {
		return fmt.Errorf("setting environment variables: %v", err)
	}
	return nil
}
