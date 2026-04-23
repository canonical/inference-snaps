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
		Use:   "run <command>",
		Short: "Run a subprocess",
		Long: "Run a command in the engine's environment\n\n" +
			"Use run to execute a program as a sub-process, within the active engine's environment.\n" +
			"To pass arguments to the program itself, separate the command and its arguments with\n" +
			"double dashes (--) from the run command and its flags. ",
		Example: "  cli run env\n" +
			"  cli run -- echo \"Hello World!\"\n" +
			"  cli run --wait-for-components -- python3 -m http.server",
		Hidden:            true,
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.run,
	}

	// flags
	cobraCmd.Flags().BoolVar(&cmd.waitForComponents, "wait-for-components", false, "wait for engine components to be installed before running")
	cobraCmd.Flags().MarkDeprecated("wait-for-components", "\"run\" always waits for components.")

	return cobraCmd
}

func (cmd *runCommand) run(_ *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("unexpected number of arguments, expected at least 1 got %d", len(args))
	}

	// Components are required for loading the engine environment
	if err := common.WaitForComponents(cmd.Context); err != nil {
		return fmt.Errorf("waiting for component: %s", err)
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

	command := args[0]

	execCmd := exec.Command(command, args[1:]...)
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
