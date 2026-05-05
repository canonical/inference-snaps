package common

import (
	"testing"
)

func TestModelStatus_NoComponents(t *testing.T) {
	status, err := modelStatus([]ComponentSettings{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(status) != 0 {
		t.Errorf("expected empty status, got: %v", status)
	}
}

func TestModelStatus_ModelNamePresent(t *testing.T) {
	settings := []ComponentSettings{
		{
			Environment: []string{
				"MODEL_NAME=my-model",
				"OTHER_VAR=value",
			},
		},
	}

	status, err := modelStatus(settings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status["name"] != "my-model" {
		t.Errorf("expected name %q, got %q", "my-model", status["name"])
	}
}

func TestModelStatus_NoModelName(t *testing.T) {
	settings := []ComponentSettings{
		{
			Environment: []string{
				"OTHER_VAR=value",
			},
		},
	}

	status, err := modelStatus(settings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := status["name"]; ok {
		t.Errorf("expected no 'name' key, but got: %q", status["name"])
	}
}

func TestModelStatus_MultipleComponents_LastWins(t *testing.T) {
	settings := []ComponentSettings{
		{
			Environment: []string{"MODEL_NAME=first-model"},
		},
		{
			Environment: []string{"MODEL_NAME=second-model"},
		},
	}

	status, err := modelStatus(settings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status["name"] != "second-model" {
		t.Errorf("expected name %q, got %q", "second-model", status["name"])
	}
}

func TestModelStatus_InvalidEnvVar(t *testing.T) {
	settings := []ComponentSettings{
		{
			Environment: []string{"INVALID_NO_EQUALS"},
		},
	}

	_, err := modelStatus(settings)
	if err == nil {
		t.Fatal("expected error for invalid env var, got nil")
	}
}

func TestModelStatus_ModelNameWithEqualsInValue(t *testing.T) {
	settings := []ComponentSettings{
		{
			Environment: []string{"MODEL_NAME=my=model"},
		},
	}

	status, err := modelStatus(settings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status["name"] != "my=model" {
		t.Errorf("expected name %q, got %q", "my=model", status["name"])
	}
}
