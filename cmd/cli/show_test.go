package main

import (
	"os"
	"testing"
)

func TestInfoLong(t *testing.T) {
	data, err := os.ReadFile("../../test_data/snap-options/engines.intel-gpu.json")
	if err != nil {
		t.Fatal(err)
	}

	engine, err := parseEngineJson(string(data))
	if err != nil {
		t.Fatal(err)
	}

	err = printEngineInfo(engine)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInfoShort(t *testing.T) {
	data, err := os.ReadFile("../../test_data/snap-options/engines.cpu.json")
	if err != nil {
		t.Fatal(err)
	}

	engine, err := parseEngineJson(string(data))
	if err != nil {
		t.Fatal(err)
	}

	err = printEngineInfo(engine)
	if err != nil {
		t.Fatal(err)
	}
}
