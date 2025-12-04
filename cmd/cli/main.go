package main

import (
	"log"
	"os"

	"github.com/canonical/go-snapctl/env"
	"github.com/canonical/inference-snaps-cli/cmd/cli/basic"
	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/cmd/cli/config"
	"github.com/canonical/inference-snaps-cli/cmd/cli/engine"
	"github.com/canonical/inference-snaps-cli/cmd/cli/others"
	"github.com/canonical/inference-snaps-cli/cmd/cli/others/debug"
	"github.com/canonical/inference-snaps-cli/pkg/storage"
	"github.com/spf13/cobra"
)

func main() {
	ctx := &common.Context{
		EnginesDir: env.Snap() + "/engines",
		Cache:      storage.NewCache(),
		Config:     storage.NewConfig(),
	}

	// rootCmd is the base command
	// It gets populated with subcommands
	rootCmd := &cobra.Command{
		SilenceUsage:      true,
		Long:              "", // Base command description TBA
		PersistentPreRunE: persistentPreRunE,
	}

	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&ctx.Verbose, "verbose", "v", false, "Enable verbose logging")

	// Use snap instance name in a snap
	if instanceName := env.SnapInstanceName(); instanceName != "" {
		rootCmd.Use = instanceName
	} else {
		rootCmd.Use = "cli"
	}

	// Disable command sorting to keep commands sorted as added below
	cobra.EnableCommandSorting = false

	rootCmd.AddGroup(basic.Group("Basic Commands:"))
	rootCmd.AddCommand(
		basic.StatusCommand(ctx),
		basic.ChatCommand(ctx),
	)

	rootCmd.AddGroup(config.Group("Configuration Commands:"))
	rootCmd.AddCommand(
		config.GetCommand(ctx),
		config.SetCommand(ctx),
	)

	rootCmd.AddGroup(engine.Group("Management Commands:"))
	rootCmd.AddCommand(
		engine.ListCommand(ctx),
		engine.ShowCommand(ctx),
		engine.UseCommand(ctx),
	)

	// other commands (help is added by default)
	rootCmd.AddCommand(
		others.ShowMachineCommand(ctx),
		others.RunCommand(ctx),
		debug.DebugCommand(ctx),
	)

	// disable logging timestamps
	log.SetFlags(0)

	// Hide the 'completion' command from help text
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func persistentPreRunE(cmd *cobra.Command, args []string) error {
	// get value of verbose flag
	verbose := cmd.Flags().Lookup("verbose").Value.String() == "true"
	if verbose {
		log.Println("Verbose output enabled globally.")
		return os.Setenv("VERBOSE", "true")
	}
	return nil
}
