package main

import (
	"os"
	"testing"
)

func TestListCompatible(t *testing.T) {
	data, err := os.ReadFile("../../test_data/snap-options/stacks.json")
	if err != nil {
		t.Fatal(err)
	}

	stacks, err := parseStacksJson(string(data))
	if err != nil {
		t.Fatal(err)
	}

	err = printStacks(stacks, false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestListAll(t *testing.T) {
	data, err := os.ReadFile("../../test_data/snap-options/stacks.json")
	if err != nil {
		t.Fatal(err)
	}

	stacks, err := parseStacksJson(string(data))
	if err != nil {
		t.Fatal(err)
	}

	err = printStacks(stacks, true)
	if err != nil {
		t.Fatal(err)
	}
}
