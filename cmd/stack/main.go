package main

import (
	"log"

	"github.com/canonical/go-snapctl/env"
	"github.com/spf13/cobra"
)

var (
	stacksDir = env.Snap() + "/stacks"
	// rootCmd is the base command
	// It gets populated with subcommands via init functions
	rootCmd = &cobra.Command{
		Use:          env.SnapInstanceName(),
		SilenceUsage: true,
	}
)

func main() {
	// disable logging timestamps
	log.SetFlags(0)

	// set a dummy root command if not in a snap
	if rootCmd.Use == "" {
		rootCmd.Use = "app"
	}

	// Define groups for subcommands - used in usage help text
	rootCmd.AddGroup(
		&cobra.Group{ID: "stacks", Title: "Stack Management Commands:"},
		&cobra.Group{ID: "config", Title: "Configuration Commands:"},
	)
	// Hide the 'completion' command from help text
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	rootCmd.Execute()
}
