package engine

import (
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/engines"
)

func TestInfoLong(t *testing.T) {
	engine, err := engines.LoadManifest("../../../test_data/engines", "intel-gpu")
	if err != nil {
		t.Fatal(err)
	}
	var scoredEngine = engines.ScoredManifest{Manifest: *engine}

	cmd := showCommand{
		format: "yaml",
	}
	err = cmd.printEngineManifest(scoredEngine)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInfoShort(t *testing.T) {
	engine, err := engines.LoadManifest("../../../test_data/engines", "cpu-avx1")
	if err != nil {
		t.Fatal(err)
	}
	var scoredEngine = engines.ScoredManifest{Manifest: *engine}

	cmd := showCommand{
		format: "yaml",
	}
	err = cmd.printEngineManifest(scoredEngine)
	if err != nil {
		t.Fatal(err)
	}
}
