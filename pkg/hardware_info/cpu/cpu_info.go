package cpu

import (
	"github.com/canonical/ml-snap-utils/pkg/types"
)

func Info() ([]types.CpuInfo, error) {
	hostLsCpu, err := hostLsCpu()
	if err != nil {
		return nil, err
	}

	cpus, err := parseLsCpu(hostLsCpu)
	if err != nil {
		return nil, err
	}

	return cpus, err
}
