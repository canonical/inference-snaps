package cpu

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/canonical/stack-utils/pkg/constants"
	"github.com/canonical/stack-utils/pkg/types"
)

func Info() ([]types.CpuInfo, error) {
	hostProcCpu, err := hostProcCpuInfo()
	if err != nil {
		return nil, err
	}

	hostUname, err := hostUnameMachine()
	if err != nil {
		return []types.CpuInfo{}, err
	}

	cpus, err := InfoFromRawData(hostProcCpu, hostUname)

	return cpus, nil
}

func InfoFromRawData(procCpuInfoData string, uname string) ([]types.CpuInfo, error) {
	architecture, err := debianArchitecture(uname)

	machineProcCpuInfo, err := parseProcCpuInfo(procCpuInfoData, architecture)
	if err != nil {
		return nil, err
	}

	cpus, err := uniqueCpuInfo(machineProcCpuInfo)

	return cpus, nil
}

func uniqueCpuInfo(procCpus []ProcCpuInfo) ([]types.CpuInfo, error) {
	// Set processor index to 0 to only check other fields for uniqueness
	for i := range procCpus {
		procCpus[i].Processor = 0
	}

	procCpus = slices.CompactFunc(procCpus, isDuplicate)

	cpuInfos, err := cpuInfoFromProc(procCpus)
	if err != nil {
		return nil, err
	}
	return cpuInfos, nil
}

func isDuplicate(a ProcCpuInfo, b ProcCpuInfo) bool {
	return reflect.DeepEqual(a, b)
}

func cpuInfoFromProc(procCpus []ProcCpuInfo) ([]types.CpuInfo, error) {
	var cpuInfos []types.CpuInfo
	for _, procCpu := range procCpus {
		var cpuInfo types.CpuInfo
		if procCpu.Architecture == constants.Amd64 {
			cpuInfo.Architecture = procCpu.Architecture
			cpuInfo.ManufacturerId = procCpu.ManufacturerId
			cpuInfo.Flags = procCpu.Flags
		} else if procCpu.Architecture == constants.Arm64 {
			cpuInfo.Architecture = procCpu.Architecture
			cpuInfo.ImplementerId = types.HexInt(procCpu.ImplementerId)
			cpuInfo.PartNumber = types.HexInt(procCpu.PartNumber)
			cpuInfo.Features = procCpu.Features
		} else {
			return nil, fmt.Errorf("unsupported architecture: %s", procCpu.Architecture)
		}
		cpuInfos = append(cpuInfos, cpuInfo)
	}
	return cpuInfos, nil
}
