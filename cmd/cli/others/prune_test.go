package others

import (
	"os"
	"testing"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/storage"
	"github.com/creack/pty"
)

func TestPrune(t *testing.T) {
	// Make stdin look like a TTY for this test.
	ptyMaster, ptySlave, err := pty.Open()
	if err != nil {
		t.Fatalf("failed to open pty: %v", err)
	}
	defer ptyMaster.Close()
	defer ptySlave.Close()

	origStdin := os.Stdin
	os.Stdin = ptySlave
	t.Cleanup(func() { os.Stdin = origStdin })

	cache := storage.NewMockCache()
	err = cache.SetActiveEngine("example-memory")
	if err != nil {
		t.Fatalf("Error setting active engine name: %v", err)
	}

	ctx := &common.Context{
		EnginesDir: "../../../test_data/engines",
		Cache:      cache,
		Config:     nil,
	}

	allEngines, err := engines.LoadManifests(ctx.EnginesDir)
	if err != nil {
		t.Fatalf("error loading engines: %v", err)
	}

	cmd := pruneCommand{Context: ctx}
	activeEngine, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		t.Fatalf("error getting active engine: %v", err)
	}

	activeEngineManifest, err := engines.LoadManifest(ctx.EnginesDir, activeEngine)
	if err != nil {
		t.Fatalf("error loading active engine manifest: %v", err)
	}

	removableComponents, err := cmd.calculateRemovableComponents(allEngines, *activeEngineManifest)
	if err != nil {
		t.Fatalf("error calculating removable components: %v", err)
	}
	for component, engines := range removableComponents {
		t.Logf("Component '%s' has removable engines: %v", component, engines)
	}

	// prints "No components to remove."
	t.Logf("\n---------------------------------------------------------------------\n")
	go func() {
		_, _ = ptyMaster.Write([]byte("y\n"))
	}()
	cmd.printComponentsAndConfirm(removableComponents, false)
	t.Logf("\n---------------------------------------------------------------------\n")

	// Test with some removable engines. The user will be prompted to confirm pruning.
	removableComponents["componentA"] = []string{"TestEngine", "AnotherTestEngine"}
	removableComponents["componentB"] = []string{"YetAnotherTestEngine", "TestEngine"}
	go func() {
		_, _ = ptyMaster.Write([]byte("y\n"))
	}()
	cmd.printComponentsAndConfirm(removableComponents, false)
	t.Logf("\n---------------------------------------------------------------------\n")

	// Test with only one engine to prune. The user will be prompted to confirm pruning.
	removableComponents["componentA"] = []string{"TestEngine"}
	removableComponents["componentB"] = []string{"TestEngine"}
	go func() {
		_, _ = ptyMaster.Write([]byte("y\n"))
	}()
	cmd.printComponentsAndConfirm(removableComponents, true)
	t.Logf("\n---------------------------------------------------------------------\n")

	// Test when the user will decline pruning.
	go func() {
		_, _ = ptyMaster.Write([]byte("n\n"))
	}()
	cmd.printComponentsAndConfirm(removableComponents, true)
}
