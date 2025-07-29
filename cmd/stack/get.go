package main

import (
	"fmt"

	"github.com/canonical/go-snapctl"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "get <key>",
		Short: "Print configuration option",
		// Long:  "",
		GroupID: "config",
		Args:    cobra.ExactArgs(1),
		RunE:    get,
	}
	rootCmd.AddCommand(cmd)
}

func get(_ *cobra.Command, args []string) error {
	return getValue(args[0])
}

func getValue(key string) error {
	value, err := snapctl.Get(key).Run()
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
