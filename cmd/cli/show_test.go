package main

import (
	"testing"

	"github.com/canonical/stack-utils/pkg/engines"
	"github.com/canonical/stack-utils/pkg/selector"
)

func TestInfoLong(t *testing.T) {
	engine, err := selector.LoadManifestFromDir("../../test_data/engines", "intel-gpu")
	if err != nil {
		t.Fatal(err)
	}
	var scoredEngine = engines.ScoredManifest{Manifest: *engine}

	err = printEngineInfo(scoredEngine)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInfoShort(t *testing.T) {
	engine, err := selector.LoadManifestFromDir("../../test_data/engines", "cpu-avx1")
	if err != nil {
		t.Fatal(err)
	}
	var scoredEngine = engines.ScoredManifest{Manifest: *engine}

	err = printEngineInfo(scoredEngine)
	if err != nil {
		t.Fatal(err)
	}
}
