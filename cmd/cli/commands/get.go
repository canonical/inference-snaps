package commands

import (
	"fmt"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type getCommand struct {
	*common.Context
}

func Get(ctx *common.Context) *cobra.Command {
	var cmd getCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:               "get [<key>]",
		Short:             "Print configurations",
		Long:              "Print one or more configurations",
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
		return fmt.Errorf("getting value of %q: %v", key, err)
	}

	if len(value) == 0 {
		return fmt.Errorf("key %q is not found\n\n%s", key, common.SuggestKeyNotFound(key))
	}

	if len(value) == 1 {
		fmt.Println(value[key])
	} else {
		// print as yaml
		yamlOutput, err := yaml.Marshal(value)
		if err != nil {
			return fmt.Errorf("serializing value: %v", err)
		}
		fmt.Printf("%s", yamlOutput) // the yaml output ends with a newline
	}

	return nil
}

func (cmd *getCommand) getValues() error {
	values, err := cmd.Config.GetAll()
	if err != nil {
		return fmt.Errorf("getting values: %v", err)
	}

	// print config value
	yamlOutput, err := yaml.Marshal(values)
	if err != nil {
		return fmt.Errorf("serializing values: %v", err)
	}
	fmt.Printf("%s", yamlOutput) // the yaml output ends with a newline

	return nil
}
