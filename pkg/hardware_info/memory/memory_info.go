package memory

import (
	"fmt"

	"github.com/canonical/stack-utils/pkg/types"
)

func Info() (types.MemoryInfo, error) {
	hostProcMemInfoData, err := hostProcMemInfo()
	if err != nil {
		return types.MemoryInfo{}, fmt.Errorf("failed to look up host /proc/meminfo: %v", err)
	}
	return InfoFromRawData(hostProcMemInfoData)
}

func InfoFromRawData(procMemInfoData string) (types.MemoryInfo, error) {
	machineMemInfo, err := parseProcMemInfo(procMemInfoData)
	if err != nil {
		return types.MemoryInfo{}, fmt.Errorf("failed to parse /proc/meminfo data: %v", err)
	}
	return machineMemInfo, nil
}
