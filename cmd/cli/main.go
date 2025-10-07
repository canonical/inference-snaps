package main

import (
	"errors"
	"log"
	"os"

	"github.com/canonical/famous-models-cli/pkg/storage"
	"github.com/canonical/go-snapctl/env"
	"github.com/spf13/cobra"
)

var (
	enginesDir       = env.Snap() + "/engines"
	snapInstanceName = env.SnapInstanceName()
	// rootCmd is the base command
	// It gets populated with subcommands via init functions
	rootCmd = &cobra.Command{
		Use:          snapInstanceName,
		SilenceUsage: true,
	}

	cache  = storage.NewCache()
	config = storage.NewConfig()

	// Error types
	ErrPermissionDenied = errors.New("permission denied, try again with sudo")
)

func main() {
	cobra.EnableCommandSorting = false

	rootCmd.AddGroup(&cobra.Group{ID: "basics", Title: "Basic Commands:"})
	addStatusCommand()
	addChatCommand()

	rootCmd.AddGroup(&cobra.Group{ID: "config", Title: "Configuration Commands:"})
	addGetCommand()
	addSetCommand()

	rootCmd.AddGroup(&cobra.Group{ID: "engines", Title: "Management Commands:"})
	addListCommand()
	addInfoCommand()
	addUseCommand()

	// other commands (help is added by default)
	addDebugCommand()

	// disable logging timestamps
	log.SetFlags(0)

	// set a dummy root command if not in a snap
	if rootCmd.Use == "" {
		rootCmd.Use = "cli"
	}

	// Hide the 'completion' command from help text
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
