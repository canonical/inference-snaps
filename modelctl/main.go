package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/canonical/inference-snaps/modelctl/v2/commands"
	"github.com/canonical/inference-snaps/modelctl/v2/commands/common"
	"github.com/canonical/inference-snaps/modelctl/v2/internal/snap"
	"github.com/canonical/inference-snaps/modelctl/v2/internal/storage"
	"github.com/spf13/cobra"
)

// version is set at build time
var version string

func main() {
	s := snap.New()
	snapDir := s.Dir()
	ctx := &common.Context{
		EnginesDir:  filepath.Join(snapDir, "engines"),
		RuntimesDir: filepath.Join(snapDir, "runtimes"),
		ModelsDir:   filepath.Join(snapDir, "models"),
		Version:     version,
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
	rootCmd, err := commands.NewRootCmd(ctx, commands.RootOptions{
		Name:  instanceName,
		Chat:  common.ChatEnabled(),
		WebUi: common.WebUiEnabled(),
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

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

	// disable logging timestamps
	log.SetFlags(0)

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
