package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/canonical/inference-snaps/modelctl/v2/commands/common"
	"github.com/canonical/inference-snaps/modelctl/v2/commands/debug"
	"github.com/spf13/cobra"
)

// RootOptions controls how the command tree is built
type RootOptions struct {
	// Name is the root command name, normally the snap instance name
	Name string
	// Chat adds the chat command
	Chat bool
	// WebUi adds the webui command
	WebUi bool
}

// NewRootCmd builds the root command with all subcommands registered
func NewRootCmd(ctx *common.Context, opts RootOptions) (*cobra.Command, error) {
	rootCmd := &cobra.Command{
		SilenceUsage: true,
		Long: opts.Name + " runs an engine that is optimized for your host machine,\n" +
			"providing a local service endpoint.\n\n" +
			"Use this command to configure the active engine, or switch to an alternative engine.",
		PersistentPreRunE: persistentPreRunE,
		Use:               opts.Name,
	}

	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&ctx.Verbose, "verbose", "v", false, "Enable verbose logging")

	// Disable command sorting to keep commands sorted as added below
	cobra.EnableCommandSorting = false

	if err := registerCommands(rootCmd, ctx, opts); err != nil {
		return nil, err
	}

	// Hide the 'completion' command from help text
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	return rootCmd, nil
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

func registerCommands(rootCmd *cobra.Command, ctx *common.Context, opts RootOptions) error {
	addCommandGroup(rootCmd, "basic", "Basic Commands:",
		Status(ctx),
		// Chat is added conditionally
	)
	if opts.Chat {
		err := appendCommandToGroup(rootCmd, "basic", Chat(ctx))
		if err != nil {
			return err
		}
	}
	if opts.WebUi {
		err := appendCommandToGroup(rootCmd, "basic", WebUi(ctx))
		if err != nil {
			return err
		}
	}

	addCommandGroup(rootCmd, "config", "Configuration Commands:",
		Get(ctx),
		Set(ctx),
		Unset(ctx),
	)

	addCommandGroup(rootCmd, "engine", "Management Commands:",
		Engines(ctx),
		Engine(ctx),
		UseEngine(ctx),
		Models(ctx),
		Model(ctx),
		UseModel(ctx),
	)

	addCommands(rootCmd,
		Machine(ctx),
		PruneCache(ctx),
		Version(ctx),
	)

	// Hidden commands
	addCommands(rootCmd,
		Run(ctx),
		ServeWebUi(ctx),
		debug.DebugCommand(ctx),
		ListEngines(ctx), // deprecated alias for "engines"
		ListModels(ctx),  // deprecated alias for "models"
		ShowEngine(ctx),  // deprecated alias for "engine"
		ShowMachine(ctx), // deprecated alias for "machine"
		ShowModel(ctx),   // deprecated alias for "model"
	)

	return nil
}
