package storage

import (
	"fmt"
	"maps"
	"strings"
)

type Config interface {
	Set(key, value string, confType configType) error
	SetDocument(key string, value any, confType configType) error
	Get(key string) (map[string]any, error)
	GetAll() (map[string]any, error)
	Unset(key string, confType configType) error
}

type config struct {
	storage storage
}

func NewConfig() Config {
	return &config{
		storage: NewSnapctlStorage(), // hardcoded since that's the only supported backend
	}
}

const configKeyPrefix = "config"

type configType string

// config precedence, from lowest to highest
var confPrecedence = []configType{
	PackageConfig, // values set by the package
	EngineConfig,  // values set by the active engine, overriding package values
	UserConfig,    // values set by the user, overriding all others
}

// config types
const (
	PackageConfig configType = "package"
	EngineConfig  configType = "engine"
	UserConfig    configType = "user"
)

// Set sets a configuration value
func (c *config) Set(key, value string, confType configType) error {
	// User configs are overrides, reject unknown keys
	if confType == UserConfig {
		valMap, err := c.Get(key)
		if err != nil {
			return fmt.Errorf("error checking existing keys: %s", err)
		}
		if len(valMap) == 0 {
			return fmt.Errorf("unknown key")
		}
	}

	return c.storage.Set(c.nestKeys(confType, key), value)
}

// SetDocument sets a configuration value that is primitive or an object
func (c *config) SetDocument(key string, value any, confType configType) error {
	return c.storage.SetDocument(c.nestKeys(confType, key), value)
}

// Get returns one or more configuration fields in as a flat map, after applying precedence rules
// If the value is a single primitive value, the map will have one entry with the full key
func (c *config) Get(key string) (map[string]any, error) {
	configs, err := c.loadConfigs()
	if err != nil {
		return nil, err
	}

	// Filter to needed keys
	for k := range configs {
		// Only keep exact key matches for both primitives and objects
		// e.g. model and model.source
		if k != key && !strings.HasPrefix(k, key+".") {
			delete(configs, k)
		}
	}

	return configs, nil
}

// GetAll returns all configurations as a flattened map
func (c *config) GetAll() (map[string]any, error) {
	return c.loadConfigs()
}

func (c *config) Unset(key string, confType configType) error {
	return c.storage.Unset(c.nestKeys(confType, key))
}

// loadConfigs loads all configurations as a flattened map, after applying precedence rules
func (c *config) loadConfigs() (map[string]any, error) {
	values, err := c.storage.Get(configKeyPrefix)
	if err != nil {
		return nil, err
	}

	// Load configurations in the order of precedence
	var finalMap = make(map[string]any)
	for _, k := range confPrecedence {
		if v, found := values[string(k)]; found {
			maps.Copy(
				finalMap,
				c.flattenMap(v.(map[string]any)),
			)
		}
	}

	return finalMap, nil
}

// flattenMap creates a single-level map with dot-separated keys
func (c *config) flattenMap(input map[string]any) map[string]any {
	flatMap := make(map[string]any)

	var recurse func(map[string]any, string)
	recurse = func(m map[string]any, prefix string) {
		for k, v := range m {
			fullKey := k
			if prefix != "" {
				fullKey = prefix + "." + k
			}
			switch val := v.(type) {
			case map[string]any:
				recurse(val, fullKey)
			default:
				flatMap[fullKey] = val
			}
		}
	}
	recurse(input, "")

	return flatMap
}

// nestKeys creates a dot-separated key with the expected prefix
func (c *config) nestKeys(confType configType, key string) string {
	if key == "." { // special case, referencing the parent
		return strings.Join([]string{configKeyPrefix, string(confType)}, ".")
	} else {
		return strings.Join([]string{configKeyPrefix, string(confType), key}, ".")
	}
}
