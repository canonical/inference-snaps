package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/canonical/inference-snaps/modelctl/v2/commands/common"
	"github.com/canonical/inference-snaps/modelctl/v2/internal/models"
	"github.com/canonical/inference-snaps/modelctl/v2/internal/storage"
)

func TestModelUnsupportedFormatResultsInError(t *testing.T) {
	manifest, err := models.LoadManifest("../test_data/models", "4b-it-int4-fq-ov")
	if err != nil {
		t.Fatalf("could not load model manifest: %v", err)
	}
	modelDetails, err := common.NewModelDetails(manifest)
	if err != nil {
		t.Fatalf("could not build model details: %v", err)
	}

	cmd := modelCommand{
		Context: &common.Context{EnginesDir: "../test_data/engines"},
		format:  "invalid-format",
	}
	err = cmd.printModelDetails(&modelDetails)

	if err == nil {
		t.Fatalf("expected unsupported format to error out, got nil error")
	}
}

func Example_modelCommand_printModelDetailsYaml() {
	// Create a temporary directory for SNAP_COMPONENTS
	tmpComponentsDir, err := ioutil.TempDir("", "snap-components-")
	if err != nil {
		panic(fmt.Sprintf("failed to create temp dir: %v", err))
	}
	defer os.RemoveAll(tmpComponentsDir)

	// Set SNAP_COMPONENTS env var for InstalledComponents() calls
	oldSnapComponents := os.Getenv("SNAP_COMPONENTS")
	defer func() {
		if oldSnapComponents != "" {
			os.Setenv("SNAP_COMPONENTS", oldSnapComponents)
		} else {
			os.Unsetenv("SNAP_COMPONENTS")
		}
	}()
	os.Setenv("SNAP_COMPONENTS", tmpComponentsDir)

	cache := storage.NewMockCache()
	if err := cache.SetActiveEngine("intel-gpu"); err != nil {
		panic(fmt.Sprintf("failed to set active engine: %v", err))
	}

	ctx := &common.Context{
		ModelsDir:   "../test_data/models",
		EnginesDir:  "../test_data/engines",
		RuntimesDir: "../test_data/runtimes",
		Cache:       cache,
	}
	machine, err := machineFixture("i7-1165G7")
	if err != nil {
		panic(fmt.Sprintf("could not create machine fixture: %v", err))
	}
	modelDetails, err := common.GetModelDetailsByNameOrAlias(ctx, "4b-it-int4-fq-ov", machine)
	if err != nil {
		panic(fmt.Sprintf("failed to load model details: %v", err))
	}

	cmd := modelCommand{Context: ctx, format: "yaml"}
	if err := cmd.printModelDetails(modelDetails); err != nil {
		panic(fmt.Sprintf("failed to print model manifest: %v", err))
	}

	// Output:
	// name: 4b-it-int4-fq-ov
	// alias: 4b-it
	// description: OpenVino 4b test model
	// model-card-url: https://example.com/model-card
	// format: OpenVINO_IR
	// quantization: int4-fq
	// capabilities:
	//     - text
	// disk-size: 6G
	// components:
	//     - model-4b-it-int4-fq-ov
	// compatible-engines:
	//     - intel-cpu
	//     - intel-gpu
	//     - intel-npu
}

func Example_modelCommand_printModelDetailsJson() {
	// Create a temporary directory for SNAP_COMPONENTS
	tmpComponentsDir, err := ioutil.TempDir("", "snap-components-")
	if err != nil {
		panic(fmt.Sprintf("failed to create temp dir: %v", err))
	}
	defer os.RemoveAll(tmpComponentsDir)

	// Set SNAP_COMPONENTS env var for InstalledComponents() calls
	oldSnapComponents := os.Getenv("SNAP_COMPONENTS")
	defer func() {
		if oldSnapComponents != "" {
			os.Setenv("SNAP_COMPONENTS", oldSnapComponents)
		} else {
			os.Unsetenv("SNAP_COMPONENTS")
		}
	}()
	os.Setenv("SNAP_COMPONENTS", tmpComponentsDir)

	cache := storage.NewMockCache()
	if err := cache.SetActiveEngine("intel-gpu"); err != nil {
		panic(fmt.Sprintf("failed to set active engine: %v", err))
	}

	ctx := &common.Context{
		ModelsDir:  "../test_data/models",
		EnginesDir: "../test_data/engines",
		Cache:      cache,
	}
	machine, err := machineFixture("i7-1165G7")
	if err != nil {
		panic(fmt.Sprintf("could not create machine fixture: %v", err))
	}
	modelDetails, err := common.GetModelDetailsByNameOrAlias(ctx, "4b-it-int4-fq-ov", machine)
	if err != nil {
		panic(fmt.Sprintf("failed to load model details: %v", err))
	}

	cmd := modelCommand{Context: ctx, format: "json"}
	if err := cmd.printModelDetails(modelDetails); err != nil {
		panic(fmt.Sprintf("failed to print model manifest: %v", err))
	}

	// Output:
	// {
	//   "name": "4b-it-int4-fq-ov",
	//   "alias": "4b-it",
	//   "description": "OpenVino 4b test model",
	//   "model-card-url": "https://example.com/model-card",
	//   "format": "OpenVINO_IR",
	//   "quantization": "int4-fq",
	//   "capabilities": [
	//     "text"
	//   ],
	//   "disk-size": "6G",
	//   "components": [
	//     "model-4b-it-int4-fq-ov"
	//   ],
	//   "compatible-engines": [
	//     "intel-cpu",
	//     "intel-gpu",
	//     "intel-npu"
	//   ]
	// }
}
