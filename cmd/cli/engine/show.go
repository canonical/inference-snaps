package engine

import (
	"encoding/json"
	"fmt"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type showCommand struct {
	*common.Context

	// flags
	format string
}

func ShowCommand(ctx *common.Context) *cobra.Command {
	var cmd showCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:     "show-engine [<engine>]",
		Short:   "Print information about an engine",
		Long:    "Print information about the active engine, or the specified engine",
		GroupID: groupID,
		// Args
		// cli use-engine <engine> requires 1 argument
		// cli use-engine --auto does not support any arguments
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: cmd.validateArgs,
		RunE:              cmd.run,
	}

	// flags
	cobraCmd.Flags().StringVar(&cmd.format, "format", "yaml", "output format")

	return cobraCmd
}

func (cmd *showCommand) run(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.showCurrentEngine()
	} else if len(args) == 1 {
		return cmd.showEngine(args[0])

	} else {
		return fmt.Errorf("invalid number of arguments")
	}
}

func (cmd *showCommand) validateArgs(_ *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
	manifests, err := engines.LoadManifests(cmd.EnginesDir)
	if err != nil {
		fmt.Printf("Error loading engines: %v\n", err)
		return nil, cobra.ShellCompDirectiveError
	}

	var engineNames []cobra.Completion
	for i := range manifests {
		engineNames = append(engineNames, manifests[i].Name)
	}

	return engineNames, cobra.ShellCompDirectiveNoSpace
}

func (cmd *showCommand) showCurrentEngine() error {
	currentEngine, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		return fmt.Errorf("could not get the active engine: %v", err)
	}
	if currentEngine == "" {
		return fmt.Errorf("no active engine")
	}
	return cmd.showEngine(currentEngine)
}

func (cmd *showCommand) showEngine(engineName string) error {
	scoredEngines, err := scoreEngines(cmd.Context)
	if err != nil {
		return fmt.Errorf("error scoring engines: %v", err)
	}

	var scoredManifest engines.ScoredManifest
	for i := range scoredEngines {
		if scoredEngines[i].Name == engineName {
			scoredManifest = scoredEngines[i]
		}
	}
	if scoredManifest.Name != engineName {
		return fmt.Errorf(`engine "%s" does not exist`, engineName)
	}

	err = cmd.printEngineManifest(scoredManifest)
	if err != nil {
		return fmt.Errorf("error printing engine manifest: %v", err)
	}
	return nil
}

func (cmd *showCommand) printEngineManifest(engine engines.ScoredManifest) error {
	switch cmd.format {
	case "json":
		jsonString, err := json.MarshalIndent(engine, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal to JSON: %s", err)
		}
		fmt.Printf("%s\n", jsonString)
	case "yaml", "":
		engineYaml, err := yaml.Marshal(engine)
		if err != nil {
			return fmt.Errorf("failed to marshal to YAML: %s", err)
		}
		fmt.Print(string(engineYaml))
	default:
		return fmt.Errorf("unknown format %q", cmd.format)
	}

	return nil
}
