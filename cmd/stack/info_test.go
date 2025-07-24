package main

import (
	"os"
	"testing"
)

func TestInfoLong(t *testing.T) {
	data, err := os.ReadFile("../../test_data/snap-options/stacks.intel-gpu.json")
	if err != nil {
		t.Fatal(err)
	}

	stack, err := parseStackJson(string(data))
	if err != nil {
		t.Fatal(err)
	}

	err = printStackInfo(stack)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInfoShort(t *testing.T) {
	data, err := os.ReadFile("../../test_data/snap-options/stacks.cpu.json")
	if err != nil {
		t.Fatal(err)
	}

	stack, err := parseStackJson(string(data))
	if err != nil {
		t.Fatal(err)
	}

	err = printStackInfo(stack)
	if err != nil {
		t.Fatal(err)
	}
}
