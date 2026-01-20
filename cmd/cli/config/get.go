package config

import (
	"fmt"
	"os"
	"slices"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// Deprecated configurations from the user
var deprecatedConfig = []string{
	"model",
	"model-name",
	"multimodel-projector",
	"server",
	"target-device",
	"http.base-path",
}

type getCommand struct {
	*common.Context
}

func GetCommand(ctx *common.Context) *cobra.Command {
	var cmd getCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:               "get [<key>]",
		Short:             "Print configurations",
		Long:              "Print one or more configurations",
		GroupID:           groupID,
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.run,
	}

	return cobraCmd
}

func (cmd *getCommand) run(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.getValues()
	} else {
		return cmd.getValue(args[0])
	}
}

func (cmd *getCommand) getValue(key string) error {
	value, err := cmd.Config.Get(key)
	if err != nil {
		return fmt.Errorf("error getting value of %q: %v", key, err)
	}

	if len(value) == 0 {
		return fmt.Errorf("no value set for key %q", key)
	}

	if len(value) == 1 {
		fmt.Println(value[key])
	} else {
		// print as yaml
		yamlOutput, err := yaml.Marshal(value)
		if err != nil {
			return fmt.Errorf("error serializing value: %v", err)
		}
		fmt.Printf("%s", yamlOutput) // the yaml output ends with a newline
	}

	// Warn the user about deprecated fields. These are still consumed by the engines.
	if slices.Contains(deprecatedConfig, key) && utils.IsTerminalOutput() {
		fmt.Fprintf(os.Stderr, "Note: %q configuration field is deprecated!\n", key)
	}

	return nil
}

func (cmd *getCommand) getValues() error {
	values, err := cmd.Config.GetAll()
	if err != nil {
		return fmt.Errorf("error getting values: %v", err)
	}

	// Drop deprecated configurations. The user doesn't need to see them.
	for k := range values {
		if slices.Contains(deprecatedConfig, k) {
			delete(values, k)
		}
	}

	// print config value
	yamlOutput, err := yaml.Marshal(values)
	if err != nil {
		return fmt.Errorf("error serializing values: %v", err)
	}
	fmt.Printf("%s", yamlOutput) // the yaml output ends with a newline

	return nil
}
