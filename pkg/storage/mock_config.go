package storage

import "maps"

type mockConfig struct {
	values map[string]any
}

func NewMockConfig(values map[string]any) Config {
	configValues := make(map[string]any)
	maps.Copy(configValues, values)

	return &mockConfig{values: configValues}
}

func (c *mockConfig) Set(key, value string, confType configType) error {
	c.values[key] = value
	return nil
}

func (c *mockConfig) SetDocument(key string, value any, confType configType) error {
	c.values[key] = value
	return nil
}

func (c *mockConfig) Get(key string) (map[string]any, error) {
	if value, found := c.values[key]; found {
		return map[string]any{key: value}, nil
	}

	return map[string]any{}, nil
}

func (c *mockConfig) GetAll() (map[string]any, error) {
	allValues := make(map[string]any)
	maps.Copy(allValues, c.values)
	return allValues, nil
}

func (c *mockConfig) Unset(key string, confType configType) error {
	delete(c.values, key)
	return nil
}
