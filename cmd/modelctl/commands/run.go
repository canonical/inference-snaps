package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/common"
	"github.com/canonical/inference-snaps-cli/v2/pkg/snap"
	"github.com/canonical/inference-snaps-cli/v2/pkg/storage"
	"github.com/canonical/inference-snaps-cli/v2/pkg/utils"
	"github.com/spf13/cobra"
)

const defaultProviderFilePath = "$SNAP_COMMON/share/provider/provider.env"

type runCommand struct {
	*common.Context

	// flags
	waitForComponents bool
	shareProvider     string
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
		Example: "  modelctl run env\n" +
			"  modelctl run -- echo \"Hello World!\"\n" +
			"  modelctl run --wait-for-components -- python3 -m http.server",
		Hidden:            true,
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.run,
	}

	// flags
	// --wait-for-components
	cobraCmd.Flags().BoolVar(&cmd.waitForComponents, "wait-for-components", false, "wait for engine components to be installed before running")
	cobraCmd.Flags().MarkDeprecated("wait-for-components", "\"run\" always waits for components.")
	// --share-provider [path]
	cobraCmd.Flags().StringVar(&cmd.shareProvider, "share-provider", "", "write provider env file to a shared path")
	cobraCmd.Flags().Lookup("share-provider").NoOptDefVal = defaultProviderFilePath

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

	if err := cmd.processEnvConfigs(); err != nil {
		return fmt.Errorf("processing env configs: %v", err)
	}
	if err := cmd.writeShareProviderEnv(); err != nil {
		return fmt.Errorf("writing share provider env: %v", err)
	}

	command := args[0]

	execCmd := exec.Command(command, args[1:]...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	return execCmd.Run()
}

func (cmd *runCommand) processEnvConfigs() error {
	envConfigs, err := cmd.Config.Get("env")
	if err != nil {
		return fmt.Errorf("getting configs: %v", err)
	}

	envVars := make(map[string]any, len(envConfigs))
	for k, v := range envConfigs {
		// Convert env keys (my-key) to environment variable names (MY_KEY)
		name, ok := strings.CutPrefix(k, storage.EnvKeyPrefix)
		if !ok {
			return fmt.Errorf("unexpected config key %q: expected prefix %q", k, storage.EnvKeyPrefix)
		}
		name = strings.ToUpper(strings.ReplaceAll(name, "-", "_"))
		envVars[name] = fmt.Sprintf("%v", v)
	}

	err = utils.SetEnvironmentVariables(envVars)
	if err != nil {
		return fmt.Errorf("setting environment variables: %v", err)
	}
	return nil
}

func (cmd *runCommand) writeShareProviderEnv() error {
	if cmd.shareProvider == "" {
		return nil
	}

	path := cmd.shareProvider
	if strings.Contains(path, "$SNAP_COMMON") || strings.Contains(path, "$SNAP_INSTANCE_NAME") {
		path = os.ExpandEnv(path)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("creating provider env directory: %v", err)
	}

	content := "SNAP_NAME=" + snap.SnapName() + "\n"
	content += "SNAP_INSTANCE_NAME=" + snap.InstanceName() + "\n"

	baseURL, err := common.OpenAiBaseUrl(cmd.Context)
	if err != nil && !errors.Is(err, common.ErrNoOpenAiServer) {
		return fmt.Errorf("getting OpenAI base URL: %v", err)
	} else if err == nil {
		content += "OPENAI_BASE_URL=" + baseURL + "\n"
	}

	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing provider env file: %v", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("renaming provider env file: %v", err)
	}

	return nil
}
