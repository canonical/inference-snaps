package storage

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/canonical/stack-utils/pkg/hardware_info"
	"github.com/canonical/stack-utils/pkg/types"
)

type cache struct {
	storage storage
}

func NewCache() *cache {
	return &cache{
		storage: NewSnapctlStorage(), // hardcoded since that's the only supported backend
	}
}

const (
	cacheKeyPrefix  = "cache."
	activeEngineKey = cacheKeyPrefix + "active-engine"
	machineInfoKey  = cacheKeyPrefix + "machine-info"
)

var ErrNoCache = errors.New("no cache")

func (c cache) SetActiveEngine(engine string) error {
	if engine == "" {
		return fmt.Errorf("engine name cannot be empty")
	}

	return c.storage.Set(activeEngineKey, engine)
}

func (c *cache) GetActiveEngine() (string, error) {
	data, err := c.storage.Get(activeEngineKey)
	if err != nil {
		return "", err
	}

	if len(data) == 0 { // cache miss
		return "", ErrNoCache
	}

	return string(data), nil
}

func (c *cache) SetMachineInfo(info types.HwInfo) error {
	return c.storage.SetDocument(machineInfoKey, info)
}

func (c *cache) GetMachineInfo() (*types.HwInfo, error) {

	data, err := c.storage.Get(machineInfoKey)
	if err != nil {
		return nil, err
	}

	var machine *types.HwInfo
	if len(data) == 0 { // cache miss
		machine, err = c.loadMachineInfo()
		if err != nil {
			return nil, err
		}
		err = c.SetMachineInfo(*machine)
		if err != nil {
			return nil, fmt.Errorf("error caching machine info: %v", err)
		}
	} else {
		err = json.Unmarshal(data, machine)
		if err != nil {
			return nil, err
		}
	}

	return machine, err
}

func (c *cache) loadMachineInfo() (*types.HwInfo, error) {
	machine, err := hardware_info.Get(true)
	if err != nil {
		return nil, fmt.Errorf("error getting machine info: %v", err)
	}

	err = c.SetMachineInfo(machine)
	if err != nil {
		return nil, fmt.Errorf("error caching machine info: %v", err)
	}

	return &machine, nil
}
