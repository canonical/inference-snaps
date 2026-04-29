package storage

import (
	"maps"
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
	allValues := make(map[string]any)
	maps.Copy(allValues, c.values)
	return allValues, nil
}

func (c *mockConfig) Unset(key string, confType configType) error {
	scopedKey := string(confType) + "." + key
	delete(c.values, scopedKey)
	return nil
}
