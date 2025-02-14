package memory

import (
	"github.com/canonical/ml-snap-utils/pkg/types"
)

func Info() (*types.MemoryInfo, error) {
	var memoryInfo types.MemoryInfo

	sysInfo, err := sysInfo()
	if err != nil {
		return nil, err
	}

	// The memory size fields need to be multiplied by the unit to get to bytes
	memoryInfo.TotalRam = sysInfo.Totalram * uint64(sysInfo.Unit)
	memoryInfo.TotalSwap = sysInfo.Totalswap * uint64(sysInfo.Unit)
	return &memoryInfo, nil
}
