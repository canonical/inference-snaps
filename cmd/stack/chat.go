package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:     "chat",
		Short:   "Start the chat CLI",
		Long:    "Start the chat CLI for interacting with the server",
		GroupID: "basics",
		Args:    cobra.NoArgs,
		RunE:    chat,
	}
	rootCmd.AddCommand(cmd)
}

func chat(_ *cobra.Command, args []string) error {
	// Run the app at path set in CHAT environment variable
	chatPath := os.Getenv("CHAT")
	if chatPath == "" {
		return fmt.Errorf("CHAT environment variable is not set")
	}
	cmd := exec.Command(chatPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running chat command: %v", err)
	}
	return nil
}
