package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func addGetCommand() {
	cmd := &cobra.Command{
		Use:   "get <key>",
		Short: "Print configuration option",
		// Long:  "",
		GroupID:           "config",
		Args:              cobra.ExactArgs(1),
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

	if value == "" {
		return fmt.Errorf("no value set for key %q", key)
	}

	// print config value
	fmt.Println(value)

	return nil
}

// not reachable currently due to arg requirement
func getValues() error {
	values, err := config.GetAll()
	if err != nil {
		return fmt.Errorf("error getting values: %v", err)
	}

	// print config value
	fmt.Println(values)

	return nil
}
