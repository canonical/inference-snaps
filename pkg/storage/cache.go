package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/canonical/go-snapctl/env"
	"github.com/canonical/inference-snaps-cli/pkg/hardware_info"
	"github.com/canonical/inference-snaps-cli/pkg/types"
)

type cache struct {
	storage             storage
	machineInfoTempFile string
}

func NewCache() *cache {
	return &cache{
		storage:             NewSnapctlStorage(), // hardcoded since that's the only supported backend
		machineInfoTempFile: "/tmp/machine-info-" + env.SnapRevision() + ".json",
	}
}

const (
	cacheKeyPrefix  = "cache."
	activeEngineKey = cacheKeyPrefix + "active-engine"
)

func (c cache) SetActiveEngine(engine string) error {
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

func (c *cache) setMachineInfo(machine types.HwInfo) error {

	b, err := json.Marshal(machine)
	if err != nil {
		return fmt.Errorf("error marshalling machine info to json: %v", err)
	}

	err = os.WriteFile(c.machineInfoTempFile, b, 0644)
	if err != nil {
		return fmt.Errorf("error writing machine info to temp file: %v", err)
	}

	return nil
}

func (c *cache) GetMachineInfo() (*types.HwInfo, error) {

	b, err := os.ReadFile(c.machineInfoTempFile)
	if err != nil {
		if os.IsNotExist(err) { // cache miss
			return c.loadMachineInfo()
		}

		return nil, fmt.Errorf("error reading machine info from temp file: %v", err)
	}

	var machine types.HwInfo
	err = json.Unmarshal(b, &machine)
	if err != nil {
		return nil, err
	}

	return &machine, err
}

func (c *cache) loadMachineInfo() (*types.HwInfo, error) {
	machine, err := hardware_info.Get(true)
	if err != nil {
		return nil, fmt.Errorf("error getting machine info: %v", err)
	}

	err = c.setMachineInfo(*machine)
	if err != nil {
		return nil, fmt.Errorf("error caching machine info: %v", err)
	}

	return machine, nil
}
