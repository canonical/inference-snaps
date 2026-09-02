package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/commands"
	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/commands/debug"
	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/common"
	"github.com/canonical/inference-snaps-cli/v2/pkg/snap"
	"github.com/canonical/inference-snaps-cli/v2/pkg/storage"
	"github.com/spf13/cobra"
)

func main() {
	s := snap.New()
	snapDir := s.Dir()
	ctx := &common.Context{
		EnginesDir:  filepath.Join(snapDir, "engines"),
		RuntimesDir: filepath.Join(snapDir, "runtimes"),
		ModelsDir:   filepath.Join(snapDir, "models"),
		Cache:       storage.NewCache(),
		Config:      storage.NewConfig(),
		Snap:        s,
	}
	// Get snap name for dynamic commands
	instanceName := s.InstanceName()
	if instanceName == "" {
		instanceName = "cli"
	}

	// rootCmd is the base command
	// It gets populated with subcommands
	rootCmd := newRootCmd(instanceName)

	// Add custom text after the help message - only show service management for top-level help
	if s.Dir() != "" {
		services, err := ctx.Snap.ServiceStatuses()
		if err != nil {
			fmt.Printf("Error: retrieving snap services: %v\n", err)
			return
		}
		if len(services) > 0 {
			helpFunc := rootCmd.HelpFunc()
			rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
				helpFunc(cmd, args)
				if cmd == rootCmd {
					fmt.Printf("\n%s\n", common.SuggestServiceManagement())
				}
			})
		}
	}

	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&ctx.Verbose, "verbose", "v", false, "Enable verbose logging")

	// Disable command sorting to keep commands sorted as added below
	cobra.EnableCommandSorting = false

	if err := registerCommands(rootCmd, ctx); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

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

// addCommandGroup adds a group of commands to the root command
func addCommandGroup(rootCmd *cobra.Command, groupID, groupTitle string, commands ...*cobra.Command) {
	group := &cobra.Group{
		ID:    groupID,
		Title: groupTitle,
	}
	rootCmd.AddGroup(group)
	for _, cmd := range commands {
		cmd.GroupID = groupID
		rootCmd.AddCommand(cmd)
	}
}

// appendCommandToGroup adds a command to an existing group in the root command
func appendCommandToGroup(rootCmd *cobra.Command, groupID string, cmd *cobra.Command) error {
	// Make sure the group exists
	for _, existingCmd := range rootCmd.Commands() {
		if existingCmd.GroupID == groupID {
			// Add command to the root command with the specified group ID
			cmd.GroupID = groupID
			rootCmd.AddCommand(cmd)
			return nil
		}
	}
	return fmt.Errorf("group with ID %q not found", groupID)
}

// addCommands adds commands to the root command without a group
// These commands will be shown in the "Additional Commands" section of the help text
func addCommands(rootCmd *cobra.Command, commands ...*cobra.Command) {
	for _, cmd := range commands {
		rootCmd.AddCommand(cmd)
	}
}

func newRootCmd(instanceName string) *cobra.Command {
	rootCmd := &cobra.Command{
		SilenceUsage: true,
		Long: instanceName + " runs an engine that is optimized for your host machine,\n" +
			"providing a local service endpoint.\n\n" +
			"Use this command to configure the active engine, or switch to an alternative engine.",
		PersistentPreRunE: persistentPreRunE,
		Use:               instanceName,
	}

	return rootCmd
}

func registerCommands(rootCmd *cobra.Command, ctx *common.Context) error {
	addCommandGroup(rootCmd, "basic", "Basic Commands:",
		commands.Status(ctx),
		// Chat is added conditionally
	)
	if common.ChatEnabled() {
		err := appendCommandToGroup(rootCmd, "basic", commands.Chat(ctx))
		if err != nil {
			return err
		}
	}
	if common.WebUiEnabled() {
		err := appendCommandToGroup(rootCmd, "basic", commands.WebUi(ctx))
		if err != nil {
			return err
		}
	}

	addCommandGroup(rootCmd, "config", "Configuration Commands:",
		commands.Get(ctx),
		commands.Set(ctx),
		commands.Unset(ctx),
	)

	addCommandGroup(rootCmd, "engine", "Management Commands:",
		commands.Engines(ctx),
		commands.Engine(ctx),
		commands.UseEngine(ctx),
		commands.Models(ctx),
		commands.Model(ctx),
		commands.UseModel(ctx),
	)

	addCommands(rootCmd,
		commands.Machine(ctx),
		commands.PruneCache(ctx),
		commands.Version(ctx),
	)

	// Hidden commands
	addCommands(rootCmd,
		commands.ExportStatus(ctx),
		commands.Run(ctx),
		commands.ServeWebUi(ctx),
		debug.DebugCommand(ctx),
		commands.ListEngines(ctx), // deprecated alias for "engines"
		commands.ListModels(ctx),  // deprecated alias for "models"
		commands.ShowEngine(ctx),  // deprecated alias for "engine"
		commands.ShowMachine(ctx), // deprecated alias for "machine"
		commands.ShowModel(ctx),   // deprecated alias for "model"
	)

	return nil
}
