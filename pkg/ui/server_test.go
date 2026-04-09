package ui

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestVerifyStaticContent(t *testing.T) {
	t.Run("missing index.html", func(t *testing.T) {
		dir := t.TempDir()
		if err := verifyStaticContent(dir); err == nil {
			t.Error("expected error when index.html is missing")
		}
	})

	t.Run("present index.html", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("<html/>"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := verifyStaticContent(dir); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("non-existent directory", func(t *testing.T) {
		if err := verifyStaticContent("/does/not/exist"); err == nil {
			t.Error("expected error for non-existent directory")
		}
	})
}

func TestConfigEndpoint(t *testing.T) {
	conf := Config{
		OpenAIBaseURL: "http://localhost:11434/v1",
		Capabilities:  []string{"text", "vision"},
		InstanceName:  "test-instance",
		EngineName:    "llama3",
	}

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("<html/>"), 0644); err != nil {
		t.Fatal(err)
	}
	mux := serverMux(conf, dir)

	t.Run("returns 200 with JSON config", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/config", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}

		var got Config
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if got.OpenAIBaseURL != conf.OpenAIBaseURL {
			t.Errorf("openAIBaseURL: got %q, want %q", got.OpenAIBaseURL, conf.OpenAIBaseURL)
		}
		if got.InstanceName != conf.InstanceName {
			t.Errorf("instanceName: got %q, want %q", got.InstanceName, conf.InstanceName)
		}
		if got.EngineName != conf.EngineName {
			t.Errorf("engineName: got %q, want %q", got.EngineName, conf.EngineName)
		}
		if len(got.Capabilities) != len(conf.Capabilities) {
			t.Errorf("capabilities: got %v, want %v", got.Capabilities, conf.Capabilities)
		}
	})

	t.Run("sets Content-Type and Cache-Control headers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/config", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type: got %q, want %q", ct, "application/json")
		}
		if cc := rec.Header().Get("Cache-Control"); cc != "no-store" {
			t.Errorf("Cache-Control: got %q, want %q", cc, "no-store")
		}
	})
}

func TestStaticFileServing(t *testing.T) {
	dir := t.TempDir()
	content := []byte("<html><body>hello</body></html>")
	if err := os.WriteFile(filepath.Join(dir, "index.html"), content, 0644); err != nil {
		t.Fatal(err)
	}

	mux := serverMux(Config{}, dir)

	t.Run("serves index.html", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
		if body := rec.Body.String(); body == "" {
			t.Error("expected non-empty response body")
		}
	})

	t.Run("returns 404 for missing file", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/missing.js", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", rec.Code)
		}
	})
}
