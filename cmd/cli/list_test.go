package main

import (
	"os"
	"testing"
)

func TestList(t *testing.T) {
	data, err := os.ReadFile("../../test_data/snap-options/engines.json")
	if err != nil {
		t.Fatal(err)
	}

	engines, err := parseEnginesJson(string(data))
	if err != nil {
		t.Fatal(err)
	}

	err = printEnginesTable(engines)
	if err != nil {
		t.Fatal(err)
	}
}
