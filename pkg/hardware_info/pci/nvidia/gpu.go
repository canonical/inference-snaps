package nvidia

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/canonical/stack-utils/pkg/types"
)

func gpuProperties(pciDevice types.PciDevice) map[string]string {
	properties := make(map[string]string)

	vRamVal, err := vRam(pciDevice)
	if err != nil {
		fmt.Fprintf(os.Stderr, "NVIDIA: error looking up vRAM: %v\n", err)
	}
	if vRamVal != nil {
		properties["vram"] = strconv.FormatUint(*vRamVal, 10)
	}

	ccVal, err := computeCapability(pciDevice)
	if err != nil {
		fmt.Fprintf(os.Stderr, "NVIDIA: error looking up compute capability: %v\n", err)
	}
	if ccVal != nil {
		properties["compute_capability"] = *ccVal
	}

	return properties
}

func vRam(device types.PciDevice) (*uint64, error) {
	/*
		Nvidia: LANG=C nvidia-smi --query-gpu=memory.total --format=csv,noheader,nounits

		$ nvidia-smi --id=00000000:01:00.0 --query-gpu=memory.total --format=csv,noheader
		4096 MiB
		$ nvidia-smi --id=00000000:02:00.0 --query-gpu=memory.total --format=csv,noheader
		No devices were found
	*/
	command := exec.Command("nvidia-smi", "--id="+device.Slot, "--query-gpu=memory.total", "--format=csv,noheader")
	command.Env = os.Environ()
	command.Env = append(command.Env, "LANG=C")
	data, err := command.Output()
	if err != nil {
		return nil, err
	} else {
		dataStr := string(data)
		dataStr = strings.TrimSpace(dataStr) // value ends in \n
		valueStr, unit, hasUnit := strings.Cut(dataStr, " ")
		vramValue, err := strconv.ParseUint(valueStr, 10, 64)
		if err != nil {
			return nil, err
		}

		if hasUnit {
			switch unit {
			case "KiB":
				vramValue = vramValue * 1024
			case "MiB":
				vramValue = vramValue * 1024 * 1024
			case "GiB":
				vramValue = vramValue * 1024 * 1024 * 1024
			}
		}

		return &vramValue, nil
	}
}

func computeCapability(device types.PciDevice) (*string, error) {
	// nvidia-smi --query-gpu=compute_cap --format=csv
	command := exec.Command("nvidia-smi", "--id="+device.Slot, "--query-gpu=compute_cap", "--format=csv,noheader")
	command.Env = os.Environ()
	command.Env = append(command.Env, "LANG=C")
	data, err := command.Output()
	if err != nil {
		return nil, err
	}

	ccValue := strings.TrimSpace(string(data))
	return &ccValue, nil
}
