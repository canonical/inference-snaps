package storage

import (
	"strings"
)

type mockConfig struct {
	values map[string]any
}

func NewMockConfig() Config {
	configValues := make(map[string]any)

	return &mockConfig{values: configValues}
}

func (c *mockConfig) Set(key, value string, confType configType) error {
	scopedKey := string(confType) + "." + key
	c.values[scopedKey] = value
	return nil
}

func (c *mockConfig) SetDocument(key string, value any, confType configType) error {
	scopedKey := string(confType) + "." + key
	c.values[scopedKey] = value
	return nil
}

func (c *mockConfig) Get(key string) (map[string]any, error) {
	if value, found := c.values[key]; found {
		return map[string]any{key: value}, nil
	}

	for _, confType := range []configType{UserConfig, EngineConfig, PackageConfig} {
		scopedKey := string(confType) + "." + key
		if value, found := c.values[scopedKey]; found {
			return map[string]any{key: value}, nil
		}
	}

	result := make(map[string]any)
	for _, confType := range []configType{UserConfig, EngineConfig, PackageConfig} {
		prefix := string(confType) + "." + key + "."
		for fullKey, value := range c.values {
			if !strings.HasPrefix(fullKey, prefix) {
				continue
			}

			normalizedKey := strings.TrimPrefix(fullKey, string(confType)+".")
			if _, exists := result[normalizedKey]; !exists {
				result[normalizedKey] = value
			}
		}
	}

	if len(result) > 0 {
		return result, nil
	}

	return map[string]any{}, nil
}

func (c *mockConfig) GetAll() (map[string]any, error) {
	result := make(map[string]any)

	for _, confType := range []configType{PackageConfig, EngineConfig, UserConfig} {
		prefix := string(confType) + "."
		for fullKey, value := range c.values {
			if !strings.HasPrefix(fullKey, prefix) {
				continue
			}

			key := strings.TrimPrefix(fullKey, prefix)
			result[key] = value
		}
	}

	return result, nil
}

func (c *mockConfig) Unset(key string, confType configType) error {
	scopedKey := string(confType) + "." + key
	delete(c.values, scopedKey)
	return nil
}

// failConfig is a Config implementation whose mutating methods always return err.
type failConfig struct {
	err error
}

// NewFailingMockConfig returns a Config where Set, SetDocument, and Unset always
// return err. Get and GetAll return empty results without error.
func NewFailingMockConfig(err error) Config {
	return &failConfig{err: err}
}

func (c *failConfig) Set(_, _ string, _ configType) error              { return c.err }
func (c *failConfig) SetDocument(_ string, _ any, _ configType) error  { return c.err }
func (c *failConfig) Unset(_ string, _ configType) error               { return c.err }
func (c *failConfig) Get(_ string) (map[string]any, error)             { return map[string]any{}, nil }
func (c *failConfig) GetAll() (map[string]any, error)                  { return map[string]any{}, nil }

// selectiveFailConfig delegates to a normal mockConfig but returns failErr when
// Unset or SetDocument is called with failKey.
type selectiveFailConfig struct {
	base    *mockConfig
	failKey string
	failErr error
}

// NewSelectiveFailingMockConfig returns a Config that behaves like NewMockConfig
// except that Unset and SetDocument return failErr when called with failKey.
func NewSelectiveFailingMockConfig(failKey string, failErr error) Config {
	return &selectiveFailConfig{
		base:    &mockConfig{values: make(map[string]any)},
		failKey: failKey,
		failErr: failErr,
	}
}

func (c *selectiveFailConfig) Set(key, value string, confType configType) error {
	return c.base.Set(key, value, confType)
}
func (c *selectiveFailConfig) SetDocument(key string, value any, confType configType) error {
	if key == c.failKey {
		return c.failErr
	}
	return c.base.SetDocument(key, value, confType)
}
func (c *selectiveFailConfig) Unset(key string, confType configType) error {
	if key == c.failKey {
		return c.failErr
	}
	return c.base.Unset(key, confType)
}
func (c *selectiveFailConfig) Get(key string) (map[string]any, error) {
	return c.base.Get(key)
}
func (c *selectiveFailConfig) GetAll() (map[string]any, error) {
	return c.base.GetAll()
}

