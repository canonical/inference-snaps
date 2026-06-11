package common

import (
	"reflect"
	"testing"

	"github.com/canonical/inference-snaps-cli/v2/pkg/storage"
)

func TestCompleteConfigKeysSuggestsKnownKeys(t *testing.T) {
	cfg := storage.NewMockConfig()
	cfg.Set("api.port", "8080", storage.UserConfig)
	cfg.Set("api.endpoint", "https://example.com", storage.UserConfig)
	cfg.Set("model", "foo", storage.UserConfig)

	got := CompleteConfigKeys(cfg, "api.", false, nil)
	want := []string{"api.endpoint", "api.port"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestCompleteConfigKeysAppendsEquals(t *testing.T) {
	cfg := storage.NewMockConfig()
	cfg.Set("api.port", "8080", storage.UserConfig)

	got := CompleteConfigKeys(cfg, "api.", true, nil)
	want := []string{"api.port="}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestCompleteConfigKeysExcludesKeys(t *testing.T) {
	cfg := storage.NewMockConfig()
	cfg.Set("api.port", "8080", storage.UserConfig)
	cfg.Set("api.endpoint", "https://example.com", storage.UserConfig)

	excluded := map[string]struct{}{"api.endpoint": {}}
	got := CompleteConfigKeys(cfg, "api.", false, excluded)
	want := []string{"api.port"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}
