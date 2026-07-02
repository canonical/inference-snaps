package commands

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteShareFiles_writesStatusJson(t *testing.T) {
	dir := t.TempDir()

	status := &exportedStatus{
		Endpoints: map[string]string{"openai": "http://localhost:8080/v1"},
	}

	cmd := &exportStatusCommand{}
	if err := cmd.writeShareFiles(status, dir); err != nil {
		t.Fatalf("writeShareFiles returned error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "status.json"))
	if err != nil {
		t.Fatalf("reading status.json: %v", err)
	}

	var got exportedStatus
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshalling status.json: %v", err)
	}

	if got.Endpoints["openai"] != "http://localhost:8080/v1" {
		t.Errorf("endpoints[openai] = %q, want %q", got.Endpoints["openai"], "http://localhost:8080/v1")
	}
}

func TestWriteShareFiles_writesOpenaiJsonWhenEndpointPresent(t *testing.T) {
	dir := t.TempDir()

	status := &exportedStatus{
		Endpoints: map[string]string{"openai": "http://localhost:8080/v1"},
	}

	cmd := &exportStatusCommand{}
	if err := cmd.writeShareFiles(status, dir); err != nil {
		t.Fatalf("writeShareFiles returned error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "openai.json"))
	if err != nil {
		t.Fatalf("reading openai.json: %v", err)
	}

	var got exportedOpenai
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshalling openai.json: %v", err)
	}

	if got.BaseUrl != "http://localhost:8080/v1" {
		t.Errorf("base_url = %q, want %q", got.BaseUrl, "http://localhost:8080/v1")
	}
}

func TestWriteShareFiles_noOpenaiJsonWhenEndpointAbsent(t *testing.T) {
	dir := t.TempDir()

	status := &exportedStatus{}

	cmd := &exportStatusCommand{}
	if err := cmd.writeShareFiles(status, dir); err != nil {
		t.Fatalf("writeShareFiles returned error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "openai.json")); !os.IsNotExist(err) {
		t.Errorf("expected openai.json to be absent, but it exists (or unexpected error: %v)", err)
	}
}

func TestWriteShareFiles_removesStaleOpenaiJsonWhenEndpointDropped(t *testing.T) {
	dir := t.TempDir()

	// First call: endpoint present → openai.json is written.
	withEndpoint := &exportedStatus{
		Endpoints: map[string]string{"openai": "http://localhost:8080/v1"},
	}
	cmd := &exportStatusCommand{}
	if err := cmd.writeShareFiles(withEndpoint, dir); err != nil {
		t.Fatalf("first writeShareFiles returned error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "openai.json")); err != nil {
		t.Fatalf("expected openai.json to exist after first call, got: %v", err)
	}

	// Second call: endpoint gone → openai.json should be removed.
	withoutEndpoint := &exportedStatus{}
	if err := cmd.writeShareFiles(withoutEndpoint, dir); err != nil {
		t.Fatalf("second writeShareFiles returned error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "openai.json")); !os.IsNotExist(err) {
		t.Errorf("expected openai.json to be removed after second call, got: %v", err)
	}
}
