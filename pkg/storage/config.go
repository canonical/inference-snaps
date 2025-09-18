package storage

import "encoding/json"

type config struct {
	storage storage
}

func NewConfig() *config {
	return &config{
		storage: NewSnapctlStorage(), // hardcoded since that's the only supported backend
	}
}

const configKeyPrefix = "config"

func (c config) Set(key, value string) error {
	return c.storage.Set(configKeyPrefix+"."+key, value)
}

// Deprecated
// Remove once migration from config to cache is complete
func (c config) SetDocument(key string, value any) error {
	return c.storage.SetDocument(configKeyPrefix+"."+key, value)
}

func (c *config) Get(key string) (string, error) {
	data, err := c.storage.Get(configKeyPrefix + "." + key)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (c *config) GetAll() (map[string]any, error) {
	data, err := c.storage.Get(configKeyPrefix)
	if err != nil {
		return nil, err
	}

	var valMap map[string]any
	err = json.Unmarshal(data, &valMap)
	if err != nil {
		return nil, err
	}

	return valMap, nil
}
