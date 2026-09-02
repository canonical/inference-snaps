package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/common"
	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type engineCommand struct {
	*common.Context

	// flags
	format string
}

func Engine(ctx *common.Context) *cobra.Command {
	return newEngineCmd(ctx, "engine [<engine>]", "")
}

// TODO: remove when we fully migrate to "engine" command
func ShowEngine(ctx *common.Context) *cobra.Command {
	return newEngineCmd(ctx, "show-engine", `use "engine" instead`)
}

func newEngineCmd(ctx *common.Context, use, deprecated string) *cobra.Command {
	var cmd engineCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:   use,
		Short: "Print information about an engine",
		Long:  "Print information about the active engine, or the specified engine",
		// Args
		// modelctl use-engine <engine> requires 1 argument
		// modelctl use-engine --auto does not support any arguments
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: cmd.validateArgs,
		RunE:              cmd.run,
		Deprecated:        deprecated,
	}

	// flags
	supportedFormats := []string{"json", "yaml"}
	cobraCmd.Flags().StringVar(
		&cmd.format,
		"format",
		"yaml",
		fmt.Sprintf("output format (%s)", strings.Join(supportedFormats, ", ")),
	)

	return cobraCmd
}

func (cmd *engineCommand) run(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.showCurrentEngine()
	} else if len(args) == 1 {
		return cmd.engine(args[0])

	} else {
		return fmt.Errorf("invalid number of arguments")
	}
}

func (cmd *engineCommand) validateArgs(_ *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
	manifests, err := engines.LoadManifests(cmd.EnginesDir)
	if err != nil {
		fmt.Printf("Error: loading engines: %v\n", err)
		return nil, cobra.ShellCompDirectiveError
	}

	var engineNames []cobra.Completion
	for i := range manifests {
		engineNames = append(engineNames, manifests[i].Name)
	}

	return engineNames, cobra.ShellCompDirectiveNoSpace
}

func (cmd *engineCommand) showCurrentEngine() error {
	currentEngine, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		return fmt.Errorf("%s: %w", common.LookingUpActiveEngine, err)
	}
	if currentEngine == "" {
		return common.ErrNoActiveEngine
	}
	return cmd.engine(currentEngine)
}

func (cmd *engineCommand) engine(engineName string) error {
	scoredEngines, err := common.ScoreEnginesWithSpinner(cmd.Context)
	if err != nil {
		return fmt.Errorf("scoring engines: %v", err)
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
		return fmt.Errorf("printing engine manifest: %v", err)
	}
	return nil
}

func (cmd *engineCommand) printEngineManifest(engine engines.ScoredManifest) error {
	var output common.EngineDetails = common.NewEngineDetails(engine)

	switch cmd.format {
	case "json":
		jsonString, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			return fmt.Errorf("json: %s", err)
		}
		fmt.Printf("%s\n", jsonString)
	case "yaml", "":
		engineYaml, err := yaml.Marshal(output)
		if err != nil {
			return fmt.Errorf("yaml: %s", err)
		}
		fmt.Print(string(engineYaml))
	default:
		return fmt.Errorf("unknown format %q", cmd.format)
	}

	return nil
}
