package main

import (
	"os"
	"testing"
)

func TestListCompatible(t *testing.T) {
	data, err := os.ReadFile("../../test_data/snap-options/engines.json")
	if err != nil {
		t.Fatal(err)
	}

	engines, err := parseEnginesJson(string(data))
	if err != nil {
		t.Fatal(err)
	}

	err = printEngines(engines, false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestListAll(t *testing.T) {
	data, err := os.ReadFile("../../test_data/snap-options/engines.json")
	if err != nil {
		t.Fatal(err)
	}

	engines, err := parseEnginesJson(string(data))
	if err != nil {
		t.Fatal(err)
	}

	err = printEngines(engines, true)
	if err != nil {
		t.Fatal(err)
	}
}
