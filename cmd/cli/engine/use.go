package engine

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/canonical/go-snapctl"
	"github.com/canonical/go-snapctl/env"
	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/selector"
	"github.com/canonical/inference-snaps-cli/pkg/snap_store"
	"github.com/canonical/inference-snaps-cli/pkg/storage"
	"github.com/canonical/inference-snaps-cli/pkg/utils"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

type useCommand struct {
	*common.Context

	// flags
	auto      bool
	assumeYes bool
}

func UseCommand(ctx *common.Context) *cobra.Command {
	var cmd useCommand
	cmd.Context = ctx

	cobra := &cobra.Command{
		Use:     "use-engine [<engine>]",
		Short:   "Select an engine",
		GroupID: groupID,
		// Args
		// cli use-engine <engine> requires 1 argument
		// cli use-engine --auto does not support any arguments
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: cmd.validateArgs,
		RunE:              cmd.run,
	}

	// flags
	cobra.Flags().BoolVar(&cmd.auto, "auto", false, "automatically select a compatible engine")
	cobra.Flags().BoolVar(&cmd.assumeYes, "assume-yes", false, "assume yes for downloading new components")

	return cobra
}

func (cmd *useCommand) validateArgs(_ *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
	scoredEngines, err := scoreEngines(cmd.Context)
	if err != nil {
		fmt.Printf("Error scoring engines: %v\n", err)
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var engineNames []cobra.Completion
	for i := range scoredEngines {
		if scoredEngines[i].Compatible {
			engineNames = append(engineNames, scoredEngines[i].Name)
		}
	}
	if len(engineNames) == 0 {
		// No engines flagged as compatible
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return engineNames, cobra.ShellCompDirectiveNoFileComp
}

func (cmd *useCommand) run(_ *cobra.Command, args []string) error {
	if !utils.IsRootUser() {
		return common.ErrPermissionDenied
	}

	if cmd.auto {
		if len(args) != 0 {
			return fmt.Errorf("cannot specify both engine name and --auto flag")
		}
		return cmd.autoSelectEngine()
	} else {
		if len(args) == 1 {
			return cmd.switchEngine(args[0])
		} else {
			return fmt.Errorf("engine name not specified")
		}
	}
}

func (cmd *useCommand) autoSelectEngine() error {
	scoredEngines, err := scoreEngines(cmd.Context)
	if err != nil {
		return fmt.Errorf("error scoring engines: %v", err)
	}

	fmt.Println("Evaluating engines for optimal hardware compatibility:")
	for _, engine := range scoredEngines {
		if engine.Score == 0 {
			fmt.Printf("✘ %s: not compatible: %s\n", engine.Name, strings.Join(engine.Notes, ", "))
		} else if engine.Grade != "stable" {
			fmt.Printf("− %s: devel, score=%d\n", engine.Name, engine.Score)
		} else {
			fmt.Printf("✔ %s: compatible, score=%d\n", engine.Name, engine.Score)
		}
	}

	selectedEngine, err := selector.TopEngine(scoredEngines)
	if err != nil {
		return fmt.Errorf("error finding top engine: %v", err)
	}

	fmt.Printf("Selected engine: %s\n", selectedEngine.Name)

	err = cmd.switchEngine(selectedEngine.Name)
	if err != nil {
		return fmt.Errorf("failed to use engine: %s", err)
	}

	return nil
}

// switchEngine changes the engine that is used by the snap
func (cmd *useCommand) switchEngine(engineName string) error {

	engine, err := engines.LoadManifest(cmd.EnginesDir, engineName)
	if err != nil {
		if errors.Is(err, engines.ErrManifestNotFound) {
			if cmd.Verbose {
				fmt.Println(err)
			}
			return fmt.Errorf("%q not found", engineName)
		}
		return fmt.Errorf("error loading engine manifest: %v", err)
	}

	components, err := cmd.missingComponents(engine.Components)
	if err != nil {
		return fmt.Errorf("error checking installed components: %v", err)
	}
	if len(components) > 0 {
		// Look up component sizes from the snap store
		componentSizes, err := snap_store.ComponentSizes()
		if err != nil {
			// If component size lookup failed, continue but log the error
			fmt.Fprintf(os.Stderr, "Error getting component sizes: %v\n", err)
		}

		// Format list of components, adding size if it is known
		var componentList []string
		for _, componentName := range components {
			line := fmt.Sprintf("- %s", componentName)
			if size, ok := componentSizes[componentName]; ok {
				line += fmt.Sprintf(" (%s)", utils.FmtBytes(uint64(size)))
			}
			componentList = append(componentList, line)
		}

		fmt.Println("Need to install the following components:")
		for _, component := range componentList {
			fmt.Println(component)
		}

		// Only ask for confirmation of download if it is an interactive terminal
		if !cmd.assumeYes && term.IsTerminal(int(os.Stdin.Fd())) {
			fmt.Println()
			if !common.ConfirmationPrompt("Do you want to continue?") {
				fmt.Println("Exiting. No changes applied.")
				return nil
			}
		}

		// Leave a blank line after printing component list and optional confirmation, before printing component installation progress
		fmt.Println()

		// This is blocking, but there is a timeout bug:
		// https://github.com/canonical/inference-snaps-cli/issues/122
		err = cmd.installComponents(engine.Components)
		if err != nil {
			return fmt.Errorf("error installing components: %v", err)
		}
	}

	activeEngineName, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		return fmt.Errorf("error getting active engine: %v", err)
	}

	if activeEngineName == engineName {
		// Engine not changed, nothing left to do
		return nil
	}

	// Unset active engine's configurations
	if activeEngineName != "" {
		err = cmd.unsetEngineConfig(activeEngineName)
		if err != nil {
			return fmt.Errorf("error un-setting engine configurations: %v", err)
		}
	}

	if len(components) > 0 {
		// Leave a blank line if components were installed, before continuing
		fmt.Println()
	}

	err = cmd.setEngineConfig(engine)
	if err != nil {
		return fmt.Errorf("error setting new engine configurations: %v", err)
	}

	// Restart if any of the services are active
	// TODO: get this from an env var instead (e.g. ENGINE_SERVICES=server,proxy)
	serviceName := env.SnapInstanceName() + ".server"

	// TODO: Only perform a restart if the service is active
	// Currently we can not reliable determine if the service is active due to a snapd bug:
	// https://bugs.launchpad.net/snapd/+bug/2137543
	// For now, rely on systemd to restart the service, which will only be done if the service is active
	fmt.Printf("Restarting %q (if active) ...\n", serviceName)
	err = snapctl.Restart(serviceName).Run()
	if err != nil {
		return fmt.Errorf("error restarting service: %v", err)
	}

	fmt.Printf("Engine successfully changed to %q\n", engineName)

	return nil
}

func (cmd *useCommand) setEngineConfig(engine *engines.Manifest) error {
	// set engine config option
	err := cmd.Cache.SetActiveEngine(engine.Name)
	if err != nil {
		return fmt.Errorf("error setting active engine: %v", err)
	}

	// set other config options
	// TODO: clear beforehand
	for confKey, confVal := range engine.Configurations {
		err = cmd.Config.SetDocument(confKey, confVal, storage.EngineConfig)
		if err != nil {
			return fmt.Errorf("error setting engine configuration %q: %v", confKey, err)
		}
	}

	return nil
}

func (cmd *useCommand) unsetEngineConfig(engineName string) error {
	// Unset all engine configurations
	err := cmd.Config.Unset(".", storage.EngineConfig)
	if err != nil {
		return fmt.Errorf("error un-setting engine configurations: %v", err)
	}

	engine, err := engines.LoadManifest(cmd.EnginesDir, engineName)
	if err != nil {
		return fmt.Errorf("error loading engine manifest: %v", err)
	}

	// Unset any user overrides
	for k := range engine.Configurations {
		err = cmd.Config.Unset(k, storage.UserConfig)
		if err != nil {
			return fmt.Errorf("error un-setting configuration %q: %v", k, err)
		}
	}

	return nil
}

// TODO: unify with similar code in run.go
func (cmd *useCommand) missingComponents(components []string) ([]string, error) {
	var missing []string
	for _, component := range components {
		isInstalled, err := cmd.componentInstalled(component)
		if err != nil {
			return missing, err
		}
		if !isInstalled {
			missing = append(missing, component)
		}
	}
	return missing, nil
}

func (*useCommand) componentInstalled(component string) (bool, error) {
	// Check in /snap/$SNAP_INSTANCE_NAME/components/$SNAP_REVISION if component is mounted
	directoryPath := fmt.Sprintf("/snap/%s/components/%s/%s", env.SnapInstanceName(), env.SnapRevision(), component)

	info, err := os.Stat(directoryPath)

	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, fmt.Errorf("error checking component directory %q: %v", component, err)
		}
	} else {
		if info.IsDir() {
			return true, nil
		} else {
			return false, fmt.Errorf("component %q exists but is not a directory", component)
		}
	}
}

func (*useCommand) installComponents(components []string) error {
	const (
		snapdUnknownSnapError = "cannot install components for a snap that is unknown to the store"
		snapdTimeoutError     = "timeout exceeded while waiting for response"
	)

	for _, component := range components {
		stopProgress := common.StartProgressSpinner("Installing " + component)
		err := snapctl.InstallComponents(component).Run()
		stopProgress()
		if err != nil {
			if strings.Contains(err.Error(), snapdUnknownSnapError) {
				return fmt.Errorf("snap not known to the store:"+
					"\nRerun this command after manually installing %q",
					component)
			} else if strings.Contains(err.Error(), snapdTimeoutError) {
				return fmt.Errorf("timed out while installing %q:"+
					"\nMonitor the installation progress with \"snap changes\""+
					"\n\nRerun this command once the installation is complete",
					component)
			} else if strings.Contains(err.Error(), "already installed") {
				continue
			} else {
				return fmt.Errorf("error installing %q: %s", component, err)
			}
		}
		fmt.Println("Installed " + component)
	}

	return nil
}
