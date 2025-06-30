package cpu

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func hostProcCpuInfo() (string, error) {
	// cat /proc/cpuinfo
	cpuInfoBytes, err := os.ReadFile("/proc/cpuinfo")
	return string(cpuInfoBytes), err
}

func parseProcCpuInfo(cpuInfoString string, architecture string) ([]ProcCpuInfo, error) {
	switch architecture {
	case amd64:
		return parseProcCpuInfoAmd64(cpuInfoString)
	case arm64:
		return parseProcCpuInfoArm64(cpuInfoString)
	default:
		return nil, fmt.Errorf("can't parse /proc/cpuinfo. unsupported architecture: %s", architecture)
	}
}

func parseProcCpuInfoAmd64(cpuInfoString string) ([]ProcCpuInfo, error) {
	var parsedCpus []ProcCpuInfo

	lines := strings.Split(cpuInfoString, "\n")
	cpuIndex := 0

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.SplitN(line, ":", 2)
		key := strings.TrimSpace(fields[0]) // remove \t between key and colon
		value := strings.TrimSpace(fields[1])

		// New cpu block
		if key == "processor" {
			newCpu := ProcCpuInfo{}
			newCpu.Architecture = amd64
			parsedCpus = append(parsedCpus, newCpu)
			cpuIndex = len(parsedCpus) - 1
		}

		switch key {
		case "processor":
			processorIndex, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, err
			}
			parsedCpus[cpuIndex].Processor = processorIndex
		case "vendor_id":
			parsedCpus[cpuIndex].ManufacturerId = value

		case "flags":
			flags := strings.Split(value, " ")
			parsedCpus[cpuIndex].Flags = append(parsedCpus[cpuIndex].Flags, flags...)

		case "model name":
			parsedCpus[cpuIndex].BrandString = value
		}
	}

	return parsedCpus, nil
}

func parseProcCpuInfoArm64(cpuInfoString string) ([]ProcCpuInfo, error) {
	var parsedCpus []ProcCpuInfo

	lines := strings.Split(cpuInfoString, "\n")
	cpuIndex := 0

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.SplitN(line, ":", 2)
		key := strings.TrimSpace(fields[0]) // remove \t between key and colon
		value := strings.TrimSpace(fields[1])

		// New cpu block
		if key == "processor" {
			newCpu := ProcCpuInfo{}
			newCpu.Architecture = arm64
			parsedCpus = append(parsedCpus, newCpu)
			cpuIndex = len(parsedCpus) - 1
		}

		switch key {

		// Formatting strings above the following cases are from https://github.com/torvalds/linux/blob/master/arch/arm64/kernel/cpuinfo.c
		// "processor\t: %d\n"
		case "processor":
			processorIndex, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, err
			}
			parsedCpus[cpuIndex].Processor = processorIndex

		// "model name\t: ARMv8 Processor rev %d (%s)\n"
		case "model name":
			modelName := strings.TrimSpace(value)
			parsedCpus[cpuIndex].ModelName = &modelName

		// BogoMIPS\t: %lu.%02lu\n
		case "BogoMIPS", "bogomips":
			bogoMips, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, err
			}
			parsedCpus[cpuIndex].BogoMips = bogoMips

		// "Features\t:"+" %s"
		case "Features":
			flags := strings.Split(value, " ")
			parsedCpus[cpuIndex].Features = append(parsedCpus[cpuIndex].Features, flags...)

		// "CPU implementer\t: 0x%02x\n"
		case "CPU implementer":
			implementer, err := strconv.ParseUint(value, 0, 8) // use base 0 to allow parser to detect and remove 0x prefix
			if err != nil {
				return nil, err
			}
			parsedCpus[cpuIndex].ImplementerId = implementer

		// "CPU architecture: 8\n"
		case "CPU architecture":
			//architecture, err := strconv.ParseUint(value, 10, 64)
			//if err != nil {
			//	return nil, err
			//}
			parsedCpus[cpuIndex].Architecture = arm64

		// "CPU variant\t: 0x%x\n"
		case "CPU variant":
			variant, err := strconv.ParseUint(value, 0, 64)
			if err != nil {
				return nil, err
			}
			parsedCpus[cpuIndex].Variant = variant

		// "CPU part\t: 0x%03x\n"
		case "CPU part":
			part, err := strconv.ParseUint(value, 0, 16)
			if err != nil {
				return nil, err
			}
			parsedCpus[cpuIndex].PartNumber = part

		// "CPU revision\t: %d\n\n"
		case "CPU revision":
			revision, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return nil, err
			}
			parsedCpus[cpuIndex].Revision = revision
		}
	}

	return parsedCpus, nil
}
