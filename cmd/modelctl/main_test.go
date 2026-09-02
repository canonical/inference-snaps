package main

import (
	"testing"

	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/common"
	"github.com/canonical/inference-snaps-cli/v2/pkg/storage"
	"github.com/spf13/cobra"
)

func testCtx() *common.Context {
	return &common.Context{
		Cache:  storage.NewMockCache(),
		Config: storage.NewMockConfig(),
	}
}

func setupRoot(t *testing.T) *cobra.Command {
	t.Helper()
	root := newRootCmd("modelctl")
	if err := registerCommands(root, testCtx()); err != nil {
		t.Fatalf("registering commands: %v", err)
	}
	return root
}

func TestCommandNames(t *testing.T) {
	root := setupRoot(t)

	want := []string{
		"status", "get", "set", "unset",
		"engines", "engine", "use-engine",
		"models", "model", "use-model",
		"machine", "prune-cache", "version",
	}

	got := map[string]bool{}

	for _, cmd := range root.Commands() {
		got[cmd.Name()] = true
	}

	for _, name := range want {
		if !got[name] {
			t.Errorf("expected command %q to be registered", name)
		}
	}
}

func TestDeprecatedCommandNames(t *testing.T) {
	root := setupRoot(t)

	deprecated := map[string]bool{
		"list-engines": true, "list-models": true,
		"show-engine": true, "show-machine": true,
		"show-model": true,
	}

	for _, cmd := range root.Commands() {
		isOld := deprecated[cmd.Name()]
		if isOld && cmd.Deprecated == "" {
			t.Errorf("%q should be marked deprecated", cmd.Name())
		}

		if !isOld && cmd.Deprecated != "" {
			t.Errorf("%q should not be deprecated: %q", cmd.Name(), cmd.Deprecated)
		}
	}
}

func TestCommandGroups(t *testing.T) {
	root := setupRoot(t)

	commands := map[string]*cobra.Command{}
	for _, c := range root.Commands() {
		commands[c.Name()] = c
	}

	commandGroups := map[string]string{
		// commandName -> GroupID
		"status":      "basic",
		"get":         "config",
		"set":         "config",
		"unset":       "config",
		"engines":     "engine",
		"engine":      "engine",
		"use-engine":  "engine",
		"models":      "engine",
		"model":       "engine",
		"use-model":   "engine",
		"machine":     "",
		"prune-cache": "",
		"version":     "",
	}

	for cmdName, expectedGroup := range commandGroups {
		t.Run(cmdName, func(t *testing.T) {
			cmd, ok := commands[cmdName]
			if !ok {
				t.Fatalf("command %q is not registered", cmdName)
			}

			if cmd.GroupID != expectedGroup {
				t.Errorf("got GroupID %q, want %q", cmd.GroupID, expectedGroup)
			}
		})
	}
}
