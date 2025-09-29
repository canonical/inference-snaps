package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func addGetCommand() {
	cmd := &cobra.Command{
		Use:   "get <key>",
		Short: "Print configuration option",
		// Long:  "",
		GroupID:           "config",
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: cobra.NoFileCompletions, // To do
		RunE:              get,
	}
	rootCmd.AddCommand(cmd)
}

func get(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		return getValues()
	} else {
		return getValue(args[0])
	}
}

func getValue(key string) error {
	value, err := config.Get(key)
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

	return nil
}

// not reachable currently due to arg requirement
func getValues() error {
	values, err := config.GetAll()
	if err != nil {
		return fmt.Errorf("error getting values: %v", err)
	}

	// print config value
	yamlOutput, err := yaml.Marshal(values)
	if err != nil {
		return fmt.Errorf("error serializing values: %v", err)
	}
	fmt.Printf("%s", yamlOutput) // the yaml output ends with a newline

	return nil
}
