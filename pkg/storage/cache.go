package storage

import (
	"errors"
	"fmt"
)

type Cache interface {
	SetActiveEngine(engine string) error
	GetActiveEngine() (string, error)
	SetActiveModel(model string) error
	GetActiveModel() (string, error)
}

type cache struct {
	storage             storage
	machineInfoTempFile string
}

func NewCache() Cache {
	return &cache{
		storage: NewSnapctlStorage(), // hardcoded since that's the only supported backend
	}
}

const (
	cacheKeyPrefix  = "cache."
	activeEngineKey = cacheKeyPrefix + "active-engine"
	activeModelKey  = cacheKeyPrefix + "active-model"
)

func (c *cache) SetActiveEngine(engine string) error {
	if engine == "" {
		return fmt.Errorf("engine name cannot be empty")
	}

	return c.storage.Set(activeEngineKey, engine)
}

// GetActiveEngine returns the currently active engine name, or an empty string if none is set
func (c *cache) GetActiveEngine() (string, error) {
	data, err := c.storage.Get(activeEngineKey)
	if err != nil {
		if errors.Is(err, ErrorNotFound) { // cache miss, no active engine set
			return "", nil
		}
		return "", err
	}

	return data[activeEngineKey].(string), nil
}

func (c *cache) SetActiveModel(model string) error {
	// An empty string is a valid model name for an engine that does not define a (default) model
	//if model == "" {
	//	return fmt.Errorf("model name cannot be empty")
	//}

	return c.storage.Set(activeModelKey, model)
}

// GetActiveModel returns the currently active model name, or an empty string if none is set
func (c *cache) GetActiveModel() (string, error) {
	data, err := c.storage.Get(activeModelKey)
	if err != nil {
		if errors.Is(err, ErrorNotFound) { // cache miss, no active model set
			return "", nil
		}
		return "", err
	}

	return data[activeModelKey].(string), nil
}
