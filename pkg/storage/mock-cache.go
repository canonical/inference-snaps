package storage

import (
	"fmt"

	"github.com/canonical/inference-snaps-cli/pkg/hardware_info"
	"github.com/canonical/inference-snaps-cli/pkg/types"
)

type mockCache struct {
	activeEngine string
	machineInfo  *types.HwInfo
}

func NewMockCache() Cache {
	return &mockCache{}
}

func (c *mockCache) SetActiveEngine(engine string) error {
	c.activeEngine = engine
	return nil
}

// GetActiveEngine returns the currently active engine name, or an empty string if none is set
func (c *mockCache) GetActiveEngine() (string, error) {
	return c.activeEngine, nil
}

func (c *mockCache) setMachineInfo(machine types.HwInfo) error {
	c.machineInfo = &machine
	return nil
}

func (c *mockCache) GetMachineInfo() (*types.HwInfo, error) {
	if c.machineInfo == nil {
		machineInfo, err := c.loadMachineInfo()
		if err != nil {
			return nil, err
		}
		c.machineInfo = machineInfo
	}
	return c.machineInfo, nil
}

func (c *mockCache) loadMachineInfo() (*types.HwInfo, error) {
	machine, err := hardware_info.Get(false)
	if err != nil {
		return nil, fmt.Errorf("error getting machine info: %v", err)
	}

	err = c.setMachineInfo(*machine)
	if err != nil {
		return nil, fmt.Errorf("error caching machine info: %v", err)
	}

	return machine, nil
}
