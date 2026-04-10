package memory

import (
	"fmt"

	"github.com/canonical/inference-snaps-cli/pkg/types"
)

func Info() (types.MemoryInfo, error) {
	hostProcMemInfoData, err := hostProcMemInfo()
	if err != nil {
		return types.MemoryInfo{}, fmt.Errorf("querying host meminfo: %v", err)
	}
	return InfoFromRawData(hostProcMemInfoData)
}

func InfoFromRawData(procMemInfoData string) (types.MemoryInfo, error) {
	machineMemInfo, err := parseProcMemInfo(procMemInfoData)
	if err != nil {
		return types.MemoryInfo{}, fmt.Errorf("parsing meminfo: %v", err)
	}
	return machineMemInfo, nil
}
