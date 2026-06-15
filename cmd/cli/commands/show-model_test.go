package commands

import (
	"fmt"
	"testing"

	"github.com/canonical/inference-snaps-cli/v2/pkg/models"
)

func TestShowModelUnsupportedFormatResultsInError(t *testing.T) {
	manifest, err := models.LoadManifest("../../../test_data/models", "4b-it-int4-fq-ov")
	if err != nil {
		t.Fatalf("could not load model manifest: %v", err)
	}

	cmd := showModelCommand{format: "invalid-format"}
	err = cmd.printModelManifest(manifest)

	if err == nil {
		t.Fatalf("expected unsupported format to error out, got nil error")
	}
}

func Example_showModelCommand_printModelManifestYaml() {
	manifest, err := models.LoadManifest("../../../test_data/models", "4b-it-int4-fq-ov")
	if err != nil {
		panic(fmt.Sprintf("failed to load model manifest: %v", err))
	}

	cmd := showModelCommand{format: "yaml"}
	if err := cmd.printModelManifest(manifest); err != nil {
		panic(fmt.Sprintf("failed to print model manifest: %v", err))
	}

	// Output:
	// id: 4b-it-int4-fq-ov
	// name: 4b
	// description: OpenVino 4b test model
	// model-card-url: https://example.com/model-card
	// quantization: int4-fq
	// capabilities:
	//     - text
	// disk-size: 6G
	// components:
	//     - model-4b-it-int4-fq-ov
}

func Example_showModelCommand_printModelManifestJson() {
	manifest, err := models.LoadManifest("../../../test_data/models", "4b-it-int4-fq-ov")
	if err != nil {
		panic(fmt.Sprintf("failed to load model manifest: %v", err))
	}

	cmd := showModelCommand{format: "json"}
	if err := cmd.printModelManifest(manifest); err != nil {
		panic(fmt.Sprintf("failed to print model manifest: %v", err))
	}

	// Output:
	// {
	//   "id": "4b-it-int4-fq-ov",
	//   "name": "4b",
	//   "description": "OpenVino 4b test model",
	//   "model-card-url": "https://example.com/model-card",
	//   "quantization": "int4-fq",
	//   "capabilities": [
	//     "text"
	//   ],
	//   "disk-size": "6G",
	//   "components": [
	//     "model-4b-it-int4-fq-ov"
	//   ]
	// }
}
